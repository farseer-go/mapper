package mapper

import (
	"reflect"
	"strconv"

	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
)

type analysisOjb struct {
	source     []*valueMeta  // 分析后的结果
	sourcePool *[]*valueMeta // 归还给 pool 用的指针
	fromMeta   *valueMeta    // 当前元数据
	usePool    bool          // 是否使用对象池
	arena      *metaArena    // 局部分配器
}

// entry 初始化分析对象
func (receiver *analysisOjb) entry(fromVal reflect.Value) []*valueMeta {
	receiver.usePool = true
	receiver.arena = getArena()

	sp := getSliceFromPool()
	receiver.source = *sp
	receiver.sourcePool = sp

	receiver.fromMeta = receiver.arena.alloc()
	receiver.fromMeta.ReflectValue = fromVal
	receiver.fromMeta.IsNil = false
	receiver.fromMeta.setReflectValue(fromVal)
	receiver.fromMeta.PointerMeta = fastReflect.PointerOfValue(receiver.fromMeta.ReflectValue)

	switch receiver.fromMeta.Type {
	case fastReflect.Map:
		receiver.analysisMap()
	case fastReflect.Struct:
		receiver.analysisStruct()
	default:
		flog.Warningf("mapper未知的类型解析：%s", receiver.fromMeta.ReflectTypeString)
	}

	return receiver.source
}

// 解析结构体
func (receiver *analysisOjb) analysisStruct() {
	if receiver.fromMeta.Level >= 20 {
		flog.Warningf("解析对象时，超过了20层深度，将停止解析:%s", receiver.fromMeta.ReflectTypeString)
		return
	}
	parent := receiver.fromMeta
	exportedField := parent.ExportedField
	parentReflectValue := parent.ReflectValue

	// 结构体：用 ByIndex 快速路径，利用 FieldTypeMetas 跳过 sync.Map 查表
	for _, i := range exportedField {
		numFieldValue := parentReflectValue.Field(i)
		receiver.fromMeta = receiver.allocStructFieldByIndex(numFieldValue, i, parent)
		receiver.analysisField()
	}
	receiver.fromMeta = parent
}

// allocStructFieldByIndex 结构体字段快速路径：利用 FieldTypeMetas 缓存，避免 PointerOfValue 查 sync.Map
func (receiver *analysisOjb) allocStructFieldByIndex(reflectValue reflect.Value, fieldIdx int, parent *valueMeta) *valueMeta {
	mt := receiver.arena.alloc()
	fillStructFieldByIndex(mt, reflectValue, fieldIdx, parent)
	return mt
}

// allocStructField 使用 arena 分配一个子元数据（替代全局 newStructField）
func (receiver *analysisOjb) allocStructField(reflectValue reflect.Value, field reflect.StructField, parent *valueMeta) *valueMeta {
	mt := receiver.arena.alloc()
	fillStructField(mt, reflectValue, field, parent)
	return mt
}

// allocItemField 使用 arena 分配 map/slice 元素元数据：利用 parent.GetItemMeta() 跳过 sync.Map 查表
func (receiver *analysisOjb) allocItemField(reflectValue reflect.Value, fieldName string, parent *valueMeta) *valueMeta {
	mt := receiver.arena.alloc()
	mt.Parent = parent
	mt.Level = parent.Level + 1
	mt.IsAnonymous = false
	mt.FieldName = fieldName

	mt.setReflectValue(reflectValue)

	// 利用 parent 预缓存的 item TypeMeta，完全跳过 sync.Map 查找
	// 当 itemMeta 是 interface{} 类型时（如 map[string]any），setReflectValue 已经取出了实际值，
	// 不能强制用 itemMeta，否则 mt.Type 会是 Interface 导致 analysisField 跳过。
	if itemMeta := parent.GetItemMeta(); itemMeta != nil && itemMeta.Type != fastReflect.Interface {
		mt.PointerMeta = fastReflect.PointerOfValueWithMeta(mt.ReflectValue, itemMeta)
	} else {
		mt.PointerMeta = fastReflect.PointerOfValue(mt.ReflectValue)
	}

	fillFullName(mt, fieldName, parent)
	return mt
}

// 解析map
func (receiver *analysisOjb) analysisMap() {
	parent := receiver.fromMeta
	// 遍历map
	miter := parent.ReflectValue.MapRange()
	for miter.Next() {
		mapKey := miter.Key()
		mapValue := miter.Value()
		fieldName := parse.ToString(mapKey.Interface())

		// 利用 parent 的 item TypeMeta 快速路径
		receiver.fromMeta = receiver.allocItemField(mapValue, fieldName, parent)
		receiver.analysisField()
	}

	receiver.fromMeta = parent
}

// 解析字典
func (receiver *analysisOjb) analysisDic() {
	// 转成map
	m := types.GetDictionaryToMap(receiver.fromMeta.ReflectValue)
	// 这里不能用receiver.valueMeta作为父级传入，而必须传receiver.valueMeta.Parent
	// 因为receiver.ReflectStructField是同一个，否则会出现Name名称重复，如：原来是A，变成AA
	receiver.fromMeta = receiver.allocStructField(m, reflect.StructField{Name: receiver.fromMeta.FieldName, Anonymous: receiver.fromMeta.IsAnonymous}, receiver.fromMeta.Parent)
	// 解析map
	receiver.analysisMap()
	// Dic统一将转成map类型，方便赋值时直接取map，而不用区分类型
	receiver.fromMeta.Type = fastReflect.Map
	receiver.source = append(receiver.source, receiver.fromMeta)
}

// 解析字段
func (receiver *analysisOjb) analysisField() {
	// 不可导出类型，则退出（此行必须先执行，否则会出现指针字段不需要赋值时，被赋值了）
	if receiver.fromMeta.IsNil || receiver.fromMeta.Type == fastReflect.Interface {
		return
	}

	metaType := receiver.fromMeta.Type
	isMap := receiver.fromMeta.IsMap

	// 先完整的赋值（如果目标类型一致，则可以直接取出来，不用分析）
	// List因为会转成切片，所以这里不能append
	if !isMap && metaType != fastReflect.List {
		receiver.source = append(receiver.source, receiver.fromMeta)
	}

	switch metaType {
	case fastReflect.GoBasicType, fastReflect.Interface:
	// 不需要处理
	case fastReflect.Slice:
		receiver.analysisSlice()
	case fastReflect.List:
		receiver.analysisList()
	case fastReflect.CustomList:
		receiver.analysisList()
	case fastReflect.Struct:
		receiver.analysisStruct()
	case fastReflect.Dic:
		receiver.analysisDic()
	case fastReflect.Map:
		// 解析map
		receiver.analysisMap()
		receiver.source = append(receiver.source, receiver.fromMeta)
		return
	default:
	}
}

// 解析切片
func (receiver *analysisOjb) analysisSlice() {
	parent := receiver.fromMeta

	length := parent.ReflectValue.Len()
	if length == 0 {
		return
	}

	for i := 0; i < length; i++ {
		sVal := parent.ReflectValue.Index(i)
		fieldName := strconv.Itoa(i)

		// 利用 parent 的 item TypeMeta 快速路径
		receiver.fromMeta = receiver.allocItemField(sVal, fieldName, parent)
		receiver.analysisField()
	}

	receiver.fromMeta = parent
}

// 解析List
func (receiver *analysisOjb) analysisList() {
	// 获取List中的数组元数
	array := types.GetListToArrayValue(receiver.fromMeta.ReflectValue)
	if array.Len() > 0 {
		parent := receiver.fromMeta

		receiver.fromMeta = receiver.allocStructField(array, reflect.StructField{}, parent)
		// 分析List中的切片
		receiver.analysisField()

		receiver.fromMeta = parent
	}
}

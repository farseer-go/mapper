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
	source []*valueMeta // 分析后的结果
	//sourceMap map[string]*valueMeta // 分析后的结果
	fromMeta *valueMeta // 当前元数据
}

// init 初始化分析对象
func (receiver *analysisOjb) entry(fromVal reflect.Value) []*valueMeta {
	// 定义存储map ,保存解析出来的字段和值
	//receiver.sourceMap = sourceMap
	// 定义存储map ,保存解析出来的字段和值
	receiver.fromMeta = &valueMeta{
		//Id:           1,
		ReflectValue: fromVal,
		IsNil:        false,
		//PointerMeta:  fastReflect.PointerOfValue(fromVal),
	}
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
	// 结构体
	for _, i := range parent.ExportedField {
		numFieldValue := parent.ReflectValue.Field(i)
		// 先分析元数据
		// 10ms
		receiver.fromMeta = newStructField(numFieldValue, parent.StructField[i], parent)
		// 12ms
		receiver.analysisField()
	}
	receiver.fromMeta = parent
}

// 解析map
func (receiver *analysisOjb) analysisMap() {
	parent := receiver.fromMeta
	// 遍历map
	miter := parent.ReflectValue.MapRange()
	for miter.Next() {
		mapKey := miter.Key()
		mapValue := miter.Value()
		field := reflect.StructField{Name: parse.ToString(mapKey.Interface())}

		// 先分析元数据
		receiver.fromMeta = newStructField(mapValue, field, parent)
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
	receiver.fromMeta = newStructField(m, reflect.StructField{Name: receiver.fromMeta.FieldName, Anonymous: receiver.fromMeta.IsAnonymous}, receiver.fromMeta.Parent)
	// 解析map
	receiver.analysisMap()
	// Dic统一将转成map类型，方便赋值时直接取map，而不用区分类型
	receiver.fromMeta.Type = fastReflect.Map
	//receiver.sourceMap[receiver.fromMeta.FullName] = &receiver.fromMeta
	receiver.source = append(receiver.source, receiver.fromMeta)
}

// 解析字段
func (receiver *analysisOjb) analysisField() {
	// 不可导出类型，则退出（此行必须先执行，否则会出现指针字段不需要赋值时，被赋值了）
	if receiver.fromMeta.IsNil || receiver.fromMeta.Type == fastReflect.Interface {
		return
	}

	// 先完整的赋值（如果目标类型一致，则可以直接取出来，不用分析）
	// List因为会转成切片，所以这里不能append
	if !receiver.fromMeta.IsMap && receiver.fromMeta.Type != fastReflect.List {
		receiver.source = append(receiver.source, receiver.fromMeta)
	}

	switch receiver.fromMeta.Type {
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
		//receiver.sourceMap[receiver.fromMeta.FullName] = &receiver.fromMeta
		receiver.source = append(receiver.source, receiver.fromMeta)
		return
	default:
	}
}

// 解析切片
func (receiver *analysisOjb) analysisSlice() {
	parent := receiver.fromMeta

	length := parent.ReflectValue.Len()
	for i := 0; i < length; i++ {
		sVal := parent.ReflectValue.Index(i)

		field := reflect.StructField{Name: strconv.Itoa(i)}

		// 先分析元数据
		valMeta := newStructField(sVal, field, parent)
		receiver.fromMeta = valMeta
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

		receiver.fromMeta = newStructField(array, reflect.StructField{}, parent)
		// 分析List中的切片
		receiver.analysisField()

		receiver.fromMeta = parent
	}
}

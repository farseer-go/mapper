package mapper

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
)

// 解析要赋值的对象
type assignObj struct {
	*valueMeta                         // 当前元数据
	sourceSlice []*valueMeta          // 分析后的结果
	sourceMap   map[uint64]*valueMeta // FullHash → meta，O(1) 查找
	usePool     bool                  // 是否使用对象池
	arena       *metaArena            // 局部分配器
}

// entry 赋值操作
func (receiver *assignObj) entry(targetVal reflect.Value, sourceSlice []*valueMeta) error {
	receiver.sourceSlice = sourceSlice
	receiver.usePool = true

	// 获取 arena 分配器
	receiver.arena = getArena()

	// 从 pool 取 map，O(1) 查找替代 O(n) 线性扫描
	receiver.sourceMap = getSourceMap()
	for _, meta := range sourceSlice {
		if meta.FullHash != 0 {
			receiver.sourceMap[meta.FullHash] = meta
		}
	}

	// 从 arena 分配根 valueMeta
	receiver.valueMeta = receiver.arena.alloc()
	receiver.valueMeta.IsNil = false
	receiver.valueMeta.ReflectValue = targetVal
	receiver.valueMeta.PointerMeta = fastReflect.PointerOfValue(targetVal)

	if receiver.valueMeta.Type != fastReflect.Struct {
		return fmt.Errorf("mapper赋值类型，必须是Struct：%s", receiver.valueMeta.ReflectTypeString)
	}

	receiver.assembleStruct(nil)

	return nil
}

// allocStructField 用 arena 分配一个子 valueMeta（赋值阶段：只计算 hash，不分配 FullName）
func (receiver *assignObj) allocStructField(reflectValue reflect.Value, field reflect.StructField, parent *valueMeta) *valueMeta {
	mt := receiver.arena.alloc()
	fillStructFieldAssign(mt, reflectValue, field, parent)
	return mt
}

// allocStructFieldByIndex 结构体字段快速路径（赋值阶段：只计算 hash，不分配 FullName）
func (receiver *assignObj) allocStructFieldByIndex(reflectValue reflect.Value, fieldIdx int, parent *valueMeta) *valueMeta {
	mt := receiver.arena.alloc()
	fillStructFieldByIndexAssign(mt, reflectValue, fieldIdx, parent)
	return mt
}

// 赋值结构体
func (receiver *assignObj) assembleStruct(sourceMeta *valueMeta) {
	// BenchmarkSample-12    	     902	   1212457 ns/op	  960268 B/op	   10007 allocs/op
	// 目标是否为指针
	// 指针类型，只有在源值存在的情况下，才赋值。否则跳过
	if sourceMeta != nil || (receiver.valueMeta.IsAddr && receiver.valueMeta.IsNil) {
		// 如果是指针，且值为nil时。receiver.ReflectType得到的是指针类型。
		// 所以这里必须使用去指针的原始类型：receiver.RealReflectType
		// 结构内字段转换 赋值。（目标字段是指针结构体，需要先初始化）

		// 当目标A.B 为指针时，如果找到源值的A.BC
		// 满足receiver.valueMeta.IsAddr && receiver.valueMeta.IsNil时也要执行，否则遍历时会异常
		receiver.NewReflectValue()
	}

	parent := receiver.valueMeta
	exportedField := parent.ExportedField
	parentReflectValue := parent.ReflectValue

	for _, i := range exportedField {
		numFieldValue := parentReflectValue.Field(i)
		// 先分析元数据：用 ByIndex 快速路径，跳过 sync.Map 查表
		valMeta := receiver.allocStructFieldByIndex(numFieldValue, i, parent)
		receiver.valueMeta = valMeta
		receiver.assignField()
	}
	receiver.valueMeta = parent
	receiver.Addr()
}

func (receiver *assignObj) assignField() {
	// 忽略未导出的字段、忽略字段
	if receiver.IsIgnore {
		return
	}

	// 7ms
	sourceValue := receiver.getSourceValue()
	if sourceValue != nil && sourceValue.Type == fastReflect.Invalid {
		return
	}

	receiverType := receiver.Type
	// 源值为nil，且不是结构、字典时，不需要继续往下走。没有意义
	// 结构体跳过，因为需要支持：Client.Id = ClientId 这种格式 && receiver.Type != Struct
	if sourceValue == nil && !receiver.IsAnonymous && (receiverType != fastReflect.Struct || !receiver.ContainsSourceKey()) { //
		//  && (receiver.Type != Struct && receiver.Type != Map && receiver.Type != Dic)
		return
	}

	// 类型完全相等时，直接赋值
	if sourceValue != nil && receiverType == sourceValue.Type && receiver.ReflectTypeString == sourceValue.ReflectTypeString {
		// 左值是指针类型，且为nil，需要先初始化
		receiver.NewReflectValue()
		receiver.ReflectValue.Set(sourceValue.ReflectValue)
		receiver.Addr()
		return
	}

	switch receiverType {
	case fastReflect.List:
		// 只处理数据源是切片的（源数据List也会转成切片）
		if sourceValue.Type != fastReflect.Slice || sourceValue == nil {
			return
		}
		receiver.assembleList(sourceValue)
	case fastReflect.Slice:
		if sourceValue.Type != fastReflect.Slice || sourceValue == nil {
			return
		}
		receiver.assembleSlice(sourceValue)
	case fastReflect.PageList:
	case fastReflect.CustomList:
		receiver.assembleCustomList(sourceValue)
	case fastReflect.Array:
	case fastReflect.Interface:
		if sourceValue.ReflectValue.CanAddr() {
			sourceValue.ReflectValue = sourceValue.ReflectValue.Addr()
		}
		receiver.ReflectValue.Set(sourceValue.ReflectValue)
	case fastReflect.GoBasicType:
		val := sourceValue.ReflectValue.Interface()
		if receiver.ReflectTypeString != sourceValue.ReflectTypeString {
			val = parse.ConvertValue(val, receiver.ReflectType)
		}
		switch {
		case receiver.Kind == reflect.String:
			*(*string)(receiver.PointerValue) = val.(string)
		case receiver.Kind == reflect.Bool:
			*(*bool)(receiver.PointerValue) = val.(bool)
		case !receiver.IsEmum && receiver.Kind == reflect.Int:
			*(*int)(receiver.PointerValue) = val.(int)
		case !receiver.IsEmum && receiver.Kind == reflect.Int8:
			*(*int8)(receiver.PointerValue) = val.(int8)
		case !receiver.IsEmum && receiver.Kind == reflect.Int16:
			*(*int16)(receiver.PointerValue) = val.(int16)
		case !receiver.IsEmum && receiver.Kind == reflect.Int32:
			*(*int32)(receiver.PointerValue) = val.(int32)
		case !receiver.IsEmum && receiver.Kind == reflect.Int64:
			*(*int64)(receiver.PointerValue) = val.(int64)
		case !receiver.IsEmum && receiver.Kind == reflect.Uint:
			*(*uint)(receiver.PointerValue) = val.(uint)
		case !receiver.IsEmum && receiver.Kind == reflect.Uint8:
			*(*uint8)(receiver.PointerValue) = val.(uint8)
		case !receiver.IsEmum && receiver.Kind == reflect.Uint16:
			*(*uint16)(receiver.PointerValue) = val.(uint16)
		case !receiver.IsEmum && receiver.Kind == reflect.Uint32:
			*(*uint32)(receiver.PointerValue) = val.(uint32)
		case !receiver.IsEmum && receiver.Kind == reflect.Uint64:
			*(*uint64)(receiver.PointerValue) = val.(uint64)
		case receiver.Kind == reflect.Float32:
			*(*float32)(receiver.PointerValue) = val.(float32)
		case receiver.Kind == reflect.Float64:
			*(*float64)(receiver.PointerValue) = val.(float64)
		case receiver.IsTime:
			*(*time.Time)(receiver.PointerValue) = val.(time.Time)
		case receiver.IsDateTime:
			*(*dateTime.DateTime)(receiver.PointerValue) = val.(dateTime.DateTime)
		default:
			receiver.ReflectValue.Set(reflect.ValueOf(val))
		}
	case fastReflect.Struct:
		receiver.assembleStruct(sourceValue)
	case fastReflect.Map:
		receiver.assembleMap(sourceValue)
	case fastReflect.Dic:
		receiver.assembleDic(sourceValue)
	default:
	}
}

// 组装List[T]
func (receiver *assignObj) assembleList(sourceMeta *valueMeta) {
	parent := receiver.valueMeta
	// 从List类型中得item类型：T
	//itemType := types.GetListItemType(receiver.ReflectType)
	// 组装[]T 元数据
	//receiver.valueMeta = NewMetaByType(reflect.SliceOf(itemType), receiver.valueMeta)
	valMeta := receiver.allocStructField(reflect.New(receiver.SliceType).Elem(), reflect.StructField{}, receiver.valueMeta)
	receiver.valueMeta = valMeta

	// 赋值组装的字段
	receiver.assembleSlice(sourceMeta)

	// new List[T]
	sliceLen := receiver.ReflectValue.Len()
	if sliceLen == 0 {
		receiver.valueMeta = parent
		return
	}

	toList := types.ListNew(parent.ReflectType, sliceLen)
	method := types.GetAddMethod(toList)

	for i := 0; i < sliceLen; i++ {
		//获取数组内的元素
		structObj := receiver.ReflectValue.Index(i)
		method.Call([]reflect.Value{structObj})
	}

	receiver.valueMeta = parent
	receiver.ReflectValue.Set(reflect.Indirect(toList))
}

// 组装List[T]
func (receiver *assignObj) assembleCustomList(sourceMeta *valueMeta) {
	parent := receiver.valueMeta
	// 从List类型中得item类型：T
	//itemType := types.GetListItemType(receiver.ReflectType)
	// 组装[]T 元数据
	//receiver.valueMeta = NewMetaByType(reflect.SliceOf(itemType), receiver.valueMeta)
	valMeta := receiver.allocStructField(reflect.New(receiver.SliceType).Elem(), reflect.StructField{}, receiver.valueMeta)
	receiver.valueMeta = valMeta
	// 赋值组装的字段
	receiver.assembleSlice(sourceMeta)

	sliceLen := receiver.ReflectValue.Len()
	if sliceLen == 0 {
		receiver.valueMeta = parent
		return
	}

	// 得到类型：List[T]
	lstType := reflect.ValueOf(reflect.New(parent.ReflectType).Elem().MethodByName("ToList").Call([]reflect.Value{})[0].Interface()).Type()
	// new List[T]
	toList := types.ListNew(lstType, sliceLen)
	method := types.GetAddMethod(toList)

	for i := 0; i < sliceLen; i++ {
		//获取数组内的元素
		structObj := receiver.ReflectValue.Index(i)
		method.Call([]reflect.Value{structObj})
	}
	receiver.valueMeta = parent
	// 转换成自定义类型
	toList = reflect.Indirect(toList).Convert(receiver.ReflectType)
	receiver.ReflectValue.Set(toList)
}

func (receiver *assignObj) assembleSlice(sourceMeta *valueMeta) {
	parent := receiver.valueMeta

	itemMeta := receiver.GetItemMeta()
	sourceItemMeta := sourceMeta.GetItemMeta()

	// T
	targetItemType := itemMeta.ReflectType

	sourceSliceCount := 0
	// 遍历源数组（前面已经判断这里一定是切片类型）
	if sourceMeta.Type == fastReflect.Slice {
		sourceSliceCount = sourceMeta.ReflectValue.Len()
	}

	// 没有元素，直接返回
	if sourceSliceCount == 0 {
		return
	}

	// 预分配切片容量
	newArr := reflect.MakeSlice(receiver.SliceType, 0, sourceSliceCount)

	// item类型一致，直接赋值
	if itemMeta.ReflectTypeString == sourceItemMeta.ReflectTypeString {
		for i := 0; i < sourceSliceCount; i++ {
			// 获取数组内的元素
			sourceItemValue := sourceMeta.ReflectValue.Index(i)
			newArr = reflect.Append(newArr, sourceItemValue)
		}
	} else {
		// 预分配 StructField，避免循环内重复创建
		field := reflect.StructField{}
		for i := 0; i < sourceSliceCount; i++ {
			// 转成切片的索引字段
			field.Name = strconv.Itoa(i)
			valMeta := receiver.allocStructField(reflect.New(targetItemType).Elem(), field, parent)
			receiver.valueMeta = valMeta
			receiver.assignField()
			newArr = reflect.Append(newArr, receiver.ReflectValue)
			// 这里改变了层级，需要恢复
			receiver.valueMeta = parent
		}
	}
	receiver.valueMeta = parent

	// 有值，才要赋值，不然会出现没意义的实例化
	receiver.NewReflectValue()
	receiver.ReflectValue.Set(newArr)
	receiver.Addr()
}

func (receiver *assignObj) assembleDic(sourceMeta *valueMeta) {
	parent := receiver.valueMeta

	// 从Dictionary类型中得source类型：map[K]V
	// new map[K]V
	newMap := reflect.New(receiver.MapType).Elem()

	// 组装map[K]V 元数据
	receiver.valueMeta = receiver.allocStructField(newMap, reflect.StructField{}, receiver.valueMeta)
	// 赋值组装的字段
	receiver.assembleMap(sourceMeta)

	// new Dictionary[K,V]
	newDictionary := types.DictionaryNew(parent.ReflectType)
	types.DictionaryAddMap(newDictionary, receiver.ReflectValue.Interface())

	receiver.valueMeta = parent
	receiver.ReflectValue.Set(newDictionary.Elem()) // newDictionary是指针类型，所以要取址
}

func (receiver *assignObj) assembleMap(sourceValue *valueMeta) {
	// 如果两边类型相待，则直接赋值
	if sourceValue != nil && receiver.Type == sourceValue.Type && receiver.ReflectTypeString == sourceValue.ReflectTypeString {
		// 左值是指针类型，且为nil，需要先初始化
		receiver.NewReflectValue()
		receiver.ReflectValue.Set(sourceValue.ReflectValue)
		receiver.Addr()
		return
	}

	parent := receiver.valueMeta
	// 遍历
	if sourceValue != nil && sourceValue.Type == fastReflect.Map {
		if receiver.IsNil {
			// 如果字段map为nil，则需要初始化
			receiver.NewReflectValue()
		}

		itemMeta := parent.GetItemMeta()
		// 是否为指针类型
		itemIsAddr := itemMeta.IsAddr

		iter := sourceValue.ReflectValue.MapRange()
		for iter.Next() {
			// 转成Map的索引字段
			mapKey := iter.Key()
			field := reflect.StructField{Name: parse.ToString(mapKey.Interface())}
			// 先分析元数据
			value := reflect.New(itemMeta.ReflectType).Elem()
			receiver.valueMeta = receiver.allocStructField(value, field, parent)
			receiver.assignField()

			// 如果左边的item是指针，则要转成指针类型
			if itemIsAddr {
				receiver.ReflectValue = receiver.ReflectValue.Addr()
			}
			parent.ReflectValue.SetMapIndex(mapKey, receiver.ReflectValue)
		}
	}

	receiver.valueMeta = parent
	receiver.Addr()
}

// 查找源字段的值
func (receiver *assignObj) getSourceValue() *valueMeta {
	if receiver.FullHash == 0 {
		return nil
	}
	return receiver.sourceMap[receiver.FullHash]
}

// ContainsSourceKey 因为需要支持：当目标字段Client为指定类型时，源值有：ClientId字段，则要支持Client.Id = ClientId 这种格式，所以使用包含的方式来判断
// 赋值阶段不预计算 FullName，此处按需从 FieldName 链构建，仅在 sourceValue==nil && Struct 时触发
func (receiver *assignObj) ContainsSourceKey() bool {
	// 从当前节点往上走，拼出 FullName（FieldName 链）
	fullName := receiver.buildFullName()
	for _, k := range receiver.sourceSlice {
		if strings.HasPrefix(k.FullName, fullName) {
			return true
		}
	}
	return false
}

// buildFullName 从 FieldName 父链重建 FullName（仅在 ContainsSourceKey 中按需调用）
func (receiver *assignObj) buildFullName() string {
	// 收集路径段（逆序）
	var parts [20]string
	n := 0
	cur := receiver.valueMeta
	for cur != nil && cur.FieldName != "" && n < 20 {
		parts[n] = cur.FieldName
		n++
		cur = cur.Parent
	}
	if n == 0 {
		return ""
	}
	// 反转并拼接
	result := parts[n-1]
	for i := n - 2; i >= 0; i-- {
		result += parts[i]
	}
	return result
}

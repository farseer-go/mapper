package mapper

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type assignObj struct {
	*valueMeta // 当前元数据
	sourceMap  map[string]*valueMeta
}

// 赋值操作
func (receiver *assignObj) assignment(targetVal reflect.Value, sourceMap map[string]*valueMeta) error {
	// 获取指针之后的值
	metaVal := newMetaVal(targetVal)
	receiver.valueMeta = metaVal
	receiver.sourceMap = sourceMap

	if receiver.valueMeta.Type != fastReflect.Struct {
		return fmt.Errorf("mapper赋值类型，必须是Struct：%s", receiver.valueMeta.ReflectTypeString)
	}

	receiver.assembleStruct(nil)

	return nil
}

// 赋值结构体
func (receiver *assignObj) assembleStruct(sourceMeta *valueMeta) {
	// 目标是否为指针
	// 指针类型，只有在源值存在的情况下，才赋值。否则跳过
	if sourceMeta != nil || (receiver.valueMeta.IsAddr && receiver.valueMeta.IsNil) {
		// 如果是指针，且值为nil时。receiver.ReflectType得到的是指针类型。
		// 所以这里必须使用去指针的原始类型：receiver.RealReflectType
		// 结构内字段转换 赋值。（目标字段是指针结构体，需要先初始化）

		// 满足receiver.valueMeta.IsAddr && receiver.valueMeta.IsNil时也要执行，否则遍历时会异常
		receiver.NewReflectValue()
	}

	parent := receiver.valueMeta
	for _, i := range parent.ExportedField {
		numFieldValue := parent.ReflectValue.Field(i)
		// 先分析元数据
		valMeta := newStructField(numFieldValue, parent.StructField[i], parent, false)
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
	sourceValue := receiver.getSourceValue()
	if sourceValue != nil && sourceValue.Type == fastReflect.Invalid {
		return
	}

	// 源值为nil，且不是结构、字典时，不需要继续往下走。没有意义
	// 结构体跳过，因为需要支持：Client.Id = ClientId 这种格式 && receiver.Type != Struct
	if sourceValue == nil && !receiver.IsAnonymous && (receiver.Type != fastReflect.Struct || !receiver.ContainsSourceKey()) { //  && (receiver.Type != Struct && receiver.Type != Map && receiver.Type != Dic)
		return
	}

	// 类型完全相等时，直接赋值
	if sourceValue != nil && receiver.Type == sourceValue.Type && receiver.ReflectTypeString == sourceValue.ReflectTypeString {
		// 左值是指针类型，且为nil，需要先初始化
		receiver.NewReflectValue()
		receiver.ReflectValue.Set(sourceValue.ReflectValue)
		receiver.Addr()
		return
	}

	switch receiver.Type {
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
	case fastReflect.GoBasicType, fastReflect.Interface:
		val := sourceValue.ReflectValue.Interface()
		if receiver.ReflectTypeString != sourceValue.ReflectTypeString {
			val = parse.ConvertValue(sourceValue.ReflectValue.Interface(), receiver.ReflectType)
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
	valMeta := newStructField(reflect.New(receiver.SliceType).Elem(), reflect.StructField{}, receiver.valueMeta, false)
	receiver.valueMeta = valMeta

	// 赋值组装的字段
	receiver.assembleSlice(sourceMeta)

	// new List[T]
	toList := types.ListNew(parent.ReflectType)
	method := types.GetAddMethod(toList)

	for i := 0; i < receiver.ReflectValue.Len(); i++ {
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
	valMeta := newStructField(reflect.New(receiver.SliceType).Elem(), reflect.StructField{}, receiver.valueMeta, false)
	receiver.valueMeta = valMeta
	// 赋值组装的字段
	receiver.assembleSlice(sourceMeta)

	// 得到类型：List[T]
	lstType := reflect.ValueOf(reflect.New(parent.ReflectType).Elem().MethodByName("ToList").Call([]reflect.Value{})[0].Interface()).Type()
	// new List[T]
	toList := types.ListNew(lstType)
	method := types.GetAddMethod(toList)

	for i := 0; i < receiver.ReflectValue.Len(); i++ {
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
	// New []T
	newArr := receiver.ZeroReflectValue
	//newArr := reflect.MakeSlice(receiver.SliceType, 0, 0)

	// 遍历源数组（前面已经判断这里一定是切片类型）
	sourceSliceCount := sourceMeta.ReflectValue.Len()
	// item类型一致，直接赋值
	if itemMeta.ReflectTypeString == sourceItemMeta.ReflectTypeString {
		for i := 0; i < sourceSliceCount; i++ {
			// 获取数组内的元素
			sourceItemValue := sourceMeta.ReflectValue.Index(i)
			newArr = reflect.Append(newArr, sourceItemValue)
		}
	} else {
		for i := 0; i < sourceSliceCount; i++ {
			// 转成切片的索引字段
			field := reflect.StructField{
				Name: parse.ToString(i),
			}
			valMeta := newStructField(reflect.New(targetItemType).Elem(), field, parent, false)
			receiver.valueMeta = valMeta
			receiver.assignField()
			newArr = reflect.Append(newArr, receiver.ReflectValue)
			// 这里改变了层级，需要恢复
			receiver.valueMeta = parent
		}
	}
	receiver.valueMeta = parent

	// 有值，才要赋值，不然会出现没意义的实例化
	if sourceSliceCount > 0 {
		receiver.NewReflectValue()
		receiver.ReflectValue.Set(newArr)
		receiver.Addr()
	}
}

func (receiver *assignObj) assembleMap(sourceMeta *valueMeta) {
	parent := receiver.valueMeta
	if sourceMeta != nil {
		// 如果字段map为nil，则需要初始化
		receiver.NewReflectValue()
	}

	// 遍历
	if sourceMeta.Type == fastReflect.Map { // sourceMeta != nil &&
		itemMeta := parent.GetItemMeta()
		// 创建一个共享的左字段（没必要每次遍历时创建一个新的）
		value := reflect.New(itemMeta.ReflectType).Elem()

		iter := sourceMeta.ReflectValue.MapRange()
		for iter.Next() {
			// 转成Map的索引字段
			mapKey := iter.Key()
			field := reflect.StructField{Name: parse.ToString(mapKey.Interface())}
			// 先分析元数据
			receiver.valueMeta = newStructField(value, field, parent, false)
			receiver.assignField()

			// 如果左边的item是指针，则要转成指针类型
			if itemMeta.IsAddr {
				receiver.ReflectValue = receiver.ReflectValue.Addr()
			}
			parent.ReflectValue.SetMapIndex(mapKey, receiver.ReflectValue)
		}
	}

	receiver.valueMeta = parent
	receiver.Addr()
}

func (receiver *assignObj) assembleDic(sourceMeta *valueMeta) {
	parent := receiver.valueMeta

	// 从Dictionary类型中得source类型：map[K]V
	//mapType := types.GetDictionaryMapType(receiver.ReflectType)
	// new map[K]V
	newMap := reflect.MakeMap(receiver.MapType)
	//newMap := receiver.ZeroReflectValue
	// 组装map[K]V 元数据
	//receiver.valueMeta = newMetaVal(newMap, receiver.valueMeta)
	valMeta := newStructField(newMap, reflect.StructField{}, receiver.valueMeta, false)
	receiver.valueMeta = valMeta
	// 赋值组装的字段
	receiver.assembleMap(sourceMeta)

	// new Dictionary[K,V]
	newDictionary := types.DictionaryNew(parent.ReflectType)
	types.DictionaryAddMap(newDictionary, receiver.ReflectValue.Interface())

	receiver.valueMeta = parent
	receiver.ReflectValue.Set(newDictionary.Elem()) // newDictionary是指针类型，所以要取址
}

func (receiver *assignObj) getSourceValue() *valueMeta {
	if receiver.Name == "" {
		return nil
	}

	// 找到源字段
	sourceValue := receiver.sourceMap[receiver.Name]
	if sourceValue != nil {
		return sourceValue
	}

	// 使用正则
	lst := collections.NewList[*valueMeta]()
	for _, v := range receiver.sourceMap {
		// 跳过没有设置正则的源
		if len(v.RegexPattern) == 0 {
			continue
		}
		if v.Regexp == nil {
			v.Regexp = regexp.MustCompile("^" + v.RegexPattern + "$")
		}
		if v.Regexp.MatchString(receiver.Name) {
			lst.Add(v)
		}
	}

	sourceValue = lst.OrderByDescending(func(item *valueMeta) any {
		return len(item.Name)
	}).First()

	return sourceValue
}

// ContainsSourceKey 因为需要支持：Client.Id = ClientId 这种格式，所以使用包含的方式来判断
func (receiver *assignObj) ContainsSourceKey() bool {
	for k := range receiver.sourceMap {
		if strings.Contains(k, receiver.Name) {
			return true
		}
	}
	return false
}

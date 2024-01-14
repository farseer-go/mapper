package mapper

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"reflect"
	"regexp"
	"strings"
)

type assignObj struct {
	*valueMeta // 当前元数据
	sourceMap  map[string]*valueMeta
}

// 赋值操作
func (receiver *assignObj) assignment(targetVal reflect.Value, sourceMap map[string]*valueMeta) error {
	// 解析父元素
	targetVal = targetVal.Elem()
	receiver.valueMeta = NewMeta(targetVal, nil)
	receiver.sourceMap = sourceMap

	if receiver.valueMeta.Type != Struct {
		return fmt.Errorf("mapper赋值类型，必须是Struct：%s", receiver.valueMeta.ReflectType.String())
	}

	receiver.assembleStruct(nil)

	return nil
}

func (receiver *assignObj) assignField() {
	// 忽略未导出的字段、忽略字段
	if !receiver.IsExported || receiver.IsIgnore {
		return
	}
	sourceValue := receiver.getSourceValue()
	if sourceValue != nil && sourceValue.Type == Invalid {
		return
	}

	if sourceValue == nil && (receiver.Type != Struct && receiver.Type != Map && receiver.Type != Dic) {
		return
	}

	switch receiver.Type {
	case List:
		// 只处理数据源是切片的（源数据List也会转成切片）
		if sourceValue.Type != Slice {
			return
		}
		receiver.assembleList(sourceValue)
	case PageList:
	case CustomList:
	case Slice:
		if sourceValue.Type != Slice {
			return
		}
		receiver.assembleSlice(sourceValue)
	case ArrayType:
	case GoBasicType:
		val := parse.ConvertValue(sourceValue.ValueAny, receiver.ReflectType)
		receiver.ReflectValue.Set(val)
	case Struct:
		receiver.assembleStruct(sourceValue)
	case Map:
		receiver.assembleMap(sourceValue)
	case Dic:
		receiver.assembleDic(sourceValue)
	default:
	}
}

// 组装List[T]
func (receiver *assignObj) assembleList(sourceMeta *valueMeta) {
	parent := receiver.valueMeta

	// 从List类型中得source类型：T
	arrType := types.GetListItemArrayType(receiver.ReflectType)
	// new []T
	newArr := reflect.MakeSlice(arrType, 0, 0)
	// 组装[]T 元数据
	receiver.valueMeta = NewMeta(newArr, receiver.valueMeta)
	// 赋值组装的字段
	receiver.assembleSlice(sourceMeta)

	// new List[T]
	toList := types.ListNew(receiver.ReflectType)
	for i := 0; i < receiver.ReflectValue.Len(); i++ {
		//获取数组内的元素
		structObj := receiver.ReflectValue.Index(i)
		types.ListAdd(toList, structObj.Interface())
	}

	receiver.valueMeta = parent
	receiver.ReflectValue.Set(toList)
}

// 数组设置值
func (receiver *assignObj) assembleSlice(sourceMeta *valueMeta) {
	parent := receiver.valueMeta

	// T
	targetItemType := receiver.ReflectType.Elem()

	// New []T
	newArr := reflect.MakeSlice(receiver.ReflectType, 0, 0)

	// 遍历源数组（前面已经判断这里一定是切片类型）
	for i := 0; i < sourceMeta.ReflectValue.Len(); i++ {
		// 获取数组内的元素
		sourceItemValue := sourceMeta.ReflectValue.Index(i)
		sourceItemMeta := NewMeta(sourceItemValue, sourceMeta)

		// item类型一致，直接赋值
		if targetItemType.String() == sourceItemMeta.ReflectType.String() {
			receiver.ReflectValue = reflect.Append(receiver.ReflectValue, sourceItemValue)
			continue
		}

		switch receiver.Type {
		case GoBasicType:
			val := parse.ConvertValue(sourceItemMeta.ValueAny, targetItemType)
			receiver.ReflectValue = reflect.Append(receiver.ReflectValue, val)
		default:
			panic("未知类型：" + receiver.ReflectType.String())
		}
	}
	receiver.valueMeta = parent

	// 赋值
	receiver.ReflectValue.Set(newArr)
}

// 集合中的Item赋值
func (receiver *assignObj) assembleStruct(sourceMeta *valueMeta) {
	// 目标是否为指针
	if receiver.IsNil {
		// 判断源值是否为nil
		sourceHaveVal := false
		for k, _ := range receiver.sourceMap {
			if strings.HasPrefix(k, receiver.Name) {
				sourceHaveVal = true
				break
			}
		}
		// 指针类型，只有在源值存在的情况下，才赋值。否则跳过
		if !sourceHaveVal {
			return
		}
		// 如果是指针，且值为nil时。receiver.ReflectType得到的是指针类型。
		// 所以这里必须使用去指针的原始类型：receiver.RealReflectType
		// 结构内字段转换 赋值。（目标字段是指针结构体，需要先初始化）
		receiver.ReflectValue.Set(reflect.New(receiver.RealReflectType))
		receiver.ReflectValue = receiver.ReflectValue.Elem()
	}

	parent := receiver.valueMeta
	for i := 0; i < parent.ReflectValue.NumField(); i++ {
		numFieldValue := parent.ReflectValue.Field(i)
		numFieldType := parent.RealReflectType.Field(i)

		// 先分析元数据
		receiver.valueMeta = newStructField(numFieldValue, numFieldType, parent)
		receiver.assignField()
	}
}

func (receiver *assignObj) assembleMap(sourceMeta *valueMeta) {
	// 类型完全相等时，直接赋值
	if m, exists := receiver.sourceMap[receiver.Name]; exists {
		if m.ReflectType.String() == receiver.ReflectType.String() {
			receiver.ReflectValue.Set(m.ReflectValue)
			return
		}
	}

	// 如果字段map为nil，则需要初始化
	if receiver.ReflectValue.IsNil() {
		receiver.ReflectValue.Set(reflect.MakeMap(receiver.ReflectType))
	}

	// 遍历
	valType := receiver.ReflectType.Elem()
	for k, v := range receiver.sourceMap {
		if strings.HasPrefix(k, receiver.Name+"{") {
			val := parse.ConvertValue(v.ReflectValue.Interface(), valType)
			receiver.ReflectValue.SetMapIndex(v.MapKey, val)
		}
	}
}

func (receiver *assignObj) assembleDic(sourceMeta *valueMeta) {
	parent := receiver.valueMeta

	// 从Dictionary类型中得source类型：map[K]V
	mapType := types.GetDictionaryMapType(receiver.ReflectType)
	// new map[K]V
	newMap := reflect.MakeMap(mapType)
	// 组装map[K]V 元数据
	receiver.valueMeta = NewMeta(newMap, receiver.valueMeta)
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
		re := regexp.MustCompile(v.RegexPattern)
		if re.MatchString(receiver.Name) {
			lst.Add(v)
		}
	}

	sourceValue = lst.OrderByDescending(func(item *valueMeta) any {
		return len(item.Name)
	}).First()

	// 没有直接匹配到
	//name := receiver.Name
	// [\{|]Exception[\}|][\{|]Age[\}|]
	// ExceptionAge
	return sourceValue
}

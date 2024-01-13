package mapper

import (
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

type analysisOjb struct {
	*valueMeta                       // 当前元数据
	sourceMap  map[string]*valueMeta // 分析后的结果
}

func (receiver *analysisOjb) analysis(from any) {
	// 定义存储map ,保存解析出来的字段和值
	receiver.sourceMap = make(map[string]*valueMeta)
	// 解析from元数据
	fromValue := reflect.Indirect(reflect.ValueOf(from))
	receiver.valueMeta = NewMeta(fromValue, nil)

	switch receiver.Type {
	case Map:
		receiver.analysisMap()
	default:
		receiver.analysisStruct()
	}
}

// 解析结构体
func (receiver *analysisOjb) analysisStruct() {
	parent := receiver.valueMeta
	// 结构体
	for i := 0; i < receiver.ReflectValue.NumField(); i++ {
		numFieldValue := receiver.ReflectValue.Field(i)
		numFieldType := receiver.RealReflectType.Field(i)

		// 先分析元数据
		receiver.NewStructField(numFieldValue, numFieldType, parent)
		receiver.analysisField()
	}
}

// 解析map
func (receiver *analysisOjb) analysisMap() {
	parent := receiver.valueMeta
	for _, mapKey := range receiver.ReflectValue.MapKeys() {
		mapValue := receiver.ReflectValue.MapIndex(mapKey)
		keyName := mapKey.String()

		if types.IsGoBasicType(mapKey.Type()) {
			keyName = parse.ToString(mapKey.Interface())
		}

		field := reflect.StructField{
			Name:    keyName,
			PkgPath: mapValue.Type().PkgPath(),
		}

		// 先分析元数据
		receiver.NewStructField(mapValue, field, parent)
		receiver.analysisField()
	}
}

// 解析map/struct的字段
func (receiver *analysisOjb) analysisField() {
	// 不可导出类型，则退出
	if receiver.IsNil || !receiver.IsExported || receiver.Type == Interface {
		return
	}
	itemName := receiver.Name
	if receiver.Parent.IsAnonymous {
		itemName = "anonymous_" + itemName
	}

	switch receiver.Type {
	case GoBasicType, CustomList, Slice:
		if receiver.valueMeta.CanInterface {
			receiver.sourceMap[itemName] = receiver.valueMeta
		}
	case List:
		array := types.ListToArray(receiver.ReflectValue)
		receiver.sourceMap[itemName] = NewMeta(reflect.ValueOf(array), receiver.valueMeta)
		return
	case Struct:
		if strings.Contains(receiver.ParentName, receiver.Name) {
			receiver.Name = receiver.ParentName
		} else {
			receiver.ParentName = receiver.Name
		}
		receiver.analysisStruct()
		return
	case Map:
		// 解析map
		receiver.analysisMap()
		if strings.Contains(receiver.ParentName, receiver.Name) {
			receiver.sourceMap[receiver.ParentName] = receiver.valueMeta
		} else {
			receiver.sourceMap[receiver.Name] = receiver.valueMeta
		}
		return
	default:
		// 非结构体遍历
		if strings.Contains(receiver.ParentName, receiver.Name) {
			receiver.sourceMap[receiver.ParentName] = receiver.valueMeta
		} else {
			receiver.sourceMap[receiver.Name] = receiver.valueMeta
		}
	}
}

// NewStructField 创建子元数据
func (receiver *analysisOjb) NewStructField(value reflect.Value, field reflect.StructField, parent *valueMeta) {
	mt := NewMeta(value, parent)
	mt.ReflectStructField = field
	mt.Name = parent.Name + field.Name
	mt.IsExported = field.IsExported()
	mt.IsAnonymous = field.Anonymous

	// 使用字段内的类型
	mt.ReflectType = field.Type
	mt.parseType()
	receiver.valueMeta = mt
}

// NewMeta 创建子元数据
func (receiver *analysisOjb) NewMeta(value reflect.Value, parent *valueMeta) {
	mt := NewMeta(value, parent)
	receiver.valueMeta = mt
}

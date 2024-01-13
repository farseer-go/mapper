package mapper

import (
	"github.com/farseer-go/fs/types"
	"reflect"
)

// 元数据
type valueMeta struct {
	Parent             *valueMeta          // 上级元数据
	ParentName         string              // 上级字段名称
	Name               string              // 字段名称
	ValueAny           any                 // 值
	ReflectValue       reflect.Value       // 值
	ReflectType        reflect.Type        // 字段类型
	RealReflectType    reflect.Type        // 字段真实类型
	ReflectStructField reflect.StructField // 字段类型
	Type               FieldType           // 集合类型
	IsNil              bool                // 是否为nil
	IsAnonymous        bool                // 是否为内嵌类型
	IsExported         bool                // 是否为可导出类型
	CanInterface       bool                // 是否可以转成Any类型
}

// NewMeta 得到类型的元数据
func NewMeta(from reflect.Value, parent *valueMeta) *valueMeta {
	meta := &valueMeta{}
	meta.Parent = parent
	meta.ReflectValue = reflect.Indirect(from)
	meta.ReflectType = meta.ReflectValue.Type()
	meta.Name = meta.ReflectType.Name()
	meta.IsNil = types.IsNil(meta.ReflectValue)
	meta.IsAnonymous = false

	// 取出实际值
	if meta.ReflectValue.CanInterface() {
		meta.CanInterface = true
		meta.ValueAny = meta.ReflectValue.Interface()
	}

	// 解析类型
	meta.parseType()

	return meta
}

func (receiver *valueMeta) parseType() {
	// 指针类型，需要取出指针指向的类型
	if receiver.ReflectType.Kind() == reflect.Pointer {
		receiver.RealReflectType = receiver.ReflectType.Elem()
	} else {
		receiver.RealReflectType = receiver.ReflectType
	}

	// 取真实的类型
	if receiver.ReflectType.Kind() == reflect.Interface && receiver.ReflectValue.CanInterface() {
		receiver.ReflectValue = receiver.ReflectValue.Elem()
		receiver.RealReflectType = receiver.ReflectValue.Type()
	}

	switch receiver.RealReflectType.Kind() {
	case reflect.Slice:
		receiver.Type = Slice
	case reflect.Array:
		receiver.Type = ArrayType
	case reflect.Map:
		receiver.Type = Map
	case reflect.Chan:
		receiver.Type = Chan
	case reflect.Func:
		receiver.Type = Func
	case reflect.Invalid:
		receiver.Type = Invalid
	case reflect.Interface:
		receiver.Type = Interface
	default:
		// 基础类型
		if types.IsGoBasicType(receiver.RealReflectType) {
			receiver.Type = GoBasicType
			return
		}

		// List类型
		if _, isTrue := types.IsListByType(receiver.RealReflectType); isTrue {
			receiver.Type = List
			return
		}

		// Dictionary类型
		if isTrue := types.IsDictionaryByType(receiver.RealReflectType); isTrue {
			receiver.Type = Dic
			return
		}

		// PageList类型
		if isTrue := types.IsPageListByType(receiver.RealReflectType); isTrue {
			receiver.Type = PageList
			return
		}

		// 自定义集合类型
		if len(receiver.RealReflectType.String()) > 8 && receiver.RealReflectType.String()[len(receiver.RealReflectType.String())-8:] == "ListType" {
			receiver.Type = CustomList
			return
		}

		if types.IsStruct(receiver.RealReflectType) {
			receiver.Type = Struct
			return
		}
		receiver.Type = Unknown
	}
}

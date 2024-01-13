package mapper

import (
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

const (
	mapSplitTag       = ":-:"
	anonymousSplitTag = ":anonymous:"
)

// 元数据
type valueMeta struct {
	Parent             *valueMeta          // 上级元数据
	ParentName         string              // 上级字段名称
	Name               string              // 字段名称
	ValueAny           any                 // 值
	ReflectValue       reflect.Value       // 值
	ReflectType        reflect.Type        // 字段类型
	RealReflectType    reflect.Type        // 字段去指针后的类型
	ReflectStructField reflect.StructField // 字段类型
	Type               FieldType           // 集合类型
	IsNil              bool                // 是否为nil
	IsAnonymous        bool                // 是否为内嵌类型
	IsExported         bool                // 是否为可导出类型
	IsIgnore           bool                // 是否为忽略字段
	CanInterface       bool                // 是否可以转成Any类型
	Level              int                 // 当前解析的层数（默认为第0层）
	MapKey             reflect.Value       // map key
}

// NewMeta 得到类型的元数据
func NewMeta(reflectValue reflect.Value, parent *valueMeta) *valueMeta {
	if reflectValue.Kind() == reflect.Pointer && !reflectValue.IsNil() {
		reflectValue = reflect.Indirect(reflectValue)
	}
	reflectType := reflectValue.Type()
	meta := NewMetaByType(reflectType, parent)
	meta.setReflectValue(reflectValue)
	return meta
}

// NewMetaByType 得到类型的元数据
func NewMetaByType(reflectType reflect.Type, parent *valueMeta) *valueMeta {
	meta := &valueMeta{}
	if parent != nil {
		meta.Parent = parent
		meta.ParentName = parent.Name
		meta.Level = parent.Level + 1
		meta.Name = parent.Name + reflectType.Name()
	}
	meta.ReflectType = reflectType
	meta.IsNil = true
	meta.IsAnonymous = false
	meta.IsExported = true

	// 解析类型
	meta.parseType()

	return meta
}

// newStructField 创建子元数据
func newStructField(value reflect.Value, field reflect.StructField, parent *valueMeta) *valueMeta {
	mt := NewMeta(value, parent)
	mt.ReflectStructField = field
	mt.IsExported = field.IsExported()
	mt.IsAnonymous = field.Anonymous

	// 内嵌字段类型的Name为类型名称，这里不需要
	if field.Anonymous {
		mt.Name = parent.Name + anonymousSplitTag
	} else {
		mt.Name = parent.Name + field.Name
	}

	// 定义的标签
	tags := strings.Split(field.Tag.Get("mapper"), ";")
	for _, tag := range tags {
		if tag == "ignore" {
			mt.IsIgnore = true
			break
		}
	}

	// 使用字段内的类型
	if field.Type != nil {
		mt.ReflectType = field.Type
		mt.parseType()
	}
	return mt
}

func (receiver *valueMeta) parseType() {
	// 指针类型，需要取出指针指向的类型
	if receiver.ReflectType.Kind() == reflect.Pointer {
		receiver.RealReflectType = receiver.ReflectType.Elem()
	} else {
		receiver.RealReflectType = receiver.ReflectType
	}

	// 取真实的类型
	if receiver.ReflectType.Kind() == reflect.Interface && !receiver.IsNil && receiver.ReflectValue.CanInterface() {
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

func (receiver *valueMeta) setReflectValue(reflectValue reflect.Value) {
	receiver.ReflectValue = reflectValue
	receiver.CanInterface = receiver.ReflectValue.CanInterface()
	receiver.IsNil = types.IsNil(receiver.ReflectValue)

	// 取出实际值
	if receiver.CanInterface && !receiver.IsNil {
		receiver.ValueAny = receiver.ReflectValue.Interface()
	}

	// 解析类型
	receiver.parseType()
}

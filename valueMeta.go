package mapper

import (
	"fmt"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

const (
	mapSplitTag           = ":-:"
	anonymousSplitTag     = "" // :anonymous:
	collectionsTypeString = "github.com/farseer-go/collections"
)

// 元数据
type valueMeta struct {
	fastReflect.PointerMeta
	Parent             *valueMeta          // 上级元数据
	ParentName         string              // 上级字段名称
	Name               string              // 字段名称
	RegexPattern       string              // 字段名称匹配规则
	ReflectValue       reflect.Value       // 值
	ReflectStructField reflect.StructField // 字段类型
	IsNil              bool                // 是否为nil
	IsAnonymous        bool                // 是否为内嵌类型
	IsIgnore           bool                // 是否为忽略字段
	Level              int                 // 当前解析的层数（默认为第0层）
	MapKey             reflect.Value       // map key
	//ValueAny           any                 // 值
	//CanInterface       bool                // 是否可以转成Any类型
}

// newMeta 得到类型的元数据
func newMetaVal(value reflect.Value, parent *valueMeta) *valueMeta {
	pointerMeta := fastReflect.PointerOf(value.Interface())
	meta := &valueMeta{}
	if parent != nil {
		meta.Parent = parent
		meta.ParentName = parent.Name
		meta.Level = parent.Level + 1
		meta.Name = parent.Name + pointerMeta.Name
	}
	meta.PointerMeta = pointerMeta
	meta.IsNil = true
	meta.IsAnonymous = false
	meta.setReflectValue(value)
	return meta
}

// newStructField 创建子元数据
func newStructField(value reflect.Value, field reflect.StructField, parent *valueMeta) *valueMeta {
	mt := newMetaVal(value, parent)
	mt.ReflectStructField = field
	mt.IsAnonymous = field.Anonymous

	// 定义的标签
	tags := strings.Split(field.Tag.Get("mapper"), ";")
	for _, tag := range tags {
		if tag == "ignore" {
			mt.IsIgnore = true
			break
		}
	}

	switch parent.Type {
	case fastReflect.Slice:
		mt.Name = parent.Name + fmt.Sprintf("[%s]", field.Name)
	case fastReflect.Map, fastReflect.Dic:
		mt.Name = parent.Name + fmt.Sprintf("{%s}", field.Name)
	default:
		// 内嵌字段类型的Name为类型名称，这里用标记代替
		if field.Anonymous {
			mt.Name = parent.Name + anonymousSplitTag
		} else {
			mt.Name = parent.Name + field.Name
		}
	}
	if mt.Name != "" {
		mt.RegexPattern = fmt.Sprintf("^%s$", mt.Name)
		mt.RegexPattern = strings.ReplaceAll(mt.RegexPattern, "{", "(\\{|)")
		mt.RegexPattern = strings.ReplaceAll(mt.RegexPattern, "}", "(\\}|)")
	}
	//mt.RegexPattern = strings.ReplaceAll(mt.RegexPattern, "[", "(\\[|)")
	//mt.RegexPattern = strings.ReplaceAll(mt.RegexPattern, "]", "(\\]|)")
	return mt
}

func (receiver *valueMeta) setReflectValue(reflectValue reflect.Value) {
	if reflectValue.Kind() == reflect.Pointer && !reflectValue.IsNil() {
		reflectValue = reflectValue.Elem()
	}

	receiver.ReflectValue = reflectValue
	//receiver.CanInterface = receiver.ReflectValue.CanInterface()
	receiver.IsNil = types.IsNil(receiver.ReflectValue)

	// 取真实的类型
	if receiver.ReflectValue.Kind() == reflect.Interface && !receiver.IsNil { // && receiver.CanInterface
		receiver.ReflectValue = receiver.ReflectValue.Elem()
	}

	// 解析类型
	//receiver.parseType()
}

// NewReflectValue 左值为指针类型时，需要先初始化
func (receiver *valueMeta) NewReflectValue() {
	if !receiver.ReflectValue.IsValid() {
		// 只能使用reflect.New,否则会出现无法寻址的问题
		receiver.ReflectValue = reflect.New(receiver.ReflectType).Elem()
		receiver.setReflectValue(receiver.ReflectValue)
		return
	}

	if types.IsNil(receiver.ReflectValue) {
		switch receiver.Type {
		case fastReflect.Slice:
			receiver.ReflectValue.Set(reflect.MakeSlice(receiver.ReflectType, 0, 0))
		case fastReflect.Map:
			receiver.ReflectValue.Set(reflect.MakeMap(receiver.ReflectType))
		default:
			receiver.ReflectValue.Set(reflect.New(receiver.ReflectType))
		}
		receiver.setReflectValue(reflect.Indirect(receiver.ReflectValue))
	}
}

// Addr 如果之前是指针，则赋值完后恢复回指针类型
// 否则如果将当前字段做为其它字段的值进行赋值时，就会出现 指针 = 非指针 赋值时的异常
func (receiver *valueMeta) Addr() {
	if receiver.IsAddr {
		receiver.ReflectValue = receiver.ReflectValue.Addr()
	}
}

/*
func (receiver *valueMeta) parseType() {
	// 取真实的类型
	if receiver.ReflectType.Kind() == reflect.Interface && !receiver.IsNil && receiver.ReflectValue.CanInterface() {
		receiver.RealReflectType = receiver.ReflectValue.Type()
		receiver.ReflectType = receiver.RealReflectType
	}

	// 指针类型，需要取出指针指向的类型
	if receiver.ReflectType.Kind() == reflect.Pointer {
		receiver.RealReflectType = receiver.ReflectType.Elem()
	} else {
		receiver.RealReflectType = receiver.ReflectType
	}

	receiver.ReflectTypeString = receiver.RealReflectType.String()

	switch receiver.RealReflectType.Kind() {
	case reflect.Slice:
		receiver.Type = fastReflect.Slice
		receiver.ItemMeta.ReflectType = receiver.RealReflectType.Elem()
		receiver.SliceType = reflect.SliceOf(receiver.ItemMeta.ReflectType)
	case reflect.Array:
		receiver.Type = fastReflect.Array
		receiver.ItemMeta.ReflectType = receiver.RealReflectType.Elem()
	case reflect.Map:
		receiver.Type = fastReflect.Map
	case reflect.Chan:
		receiver.Type = fastReflect.Chan
	case reflect.Func:
		receiver.Type = fastReflect.Func
	case reflect.Invalid:
		receiver.Type = fastReflect.Invalid
	case reflect.Interface:
		receiver.Type = fastReflect.Interface
	default:
		// 基础类型
		if types.IsGoBasicType(receiver.RealReflectType) {
			receiver.Type = fastReflect.GoBasicType
			break
		}

		// List类型
		if _, isTrue := types.IsListByType(receiver.RealReflectType); isTrue {
			receiver.Type = fastReflect.List
			receiver.ItemMeta.ReflectType = types.GetListItemType(receiver.ReflectType)
			receiver.SliceType = reflect.SliceOf(receiver.ItemMeta.ReflectType)
			break
		}

		// Dictionary类型
		if isTrue := types.IsDictionaryByType(receiver.RealReflectType); isTrue {
			receiver.Type = fastReflect.Dic
			break
		}

		// PageList类型
		if isTrue := types.IsPageListByType(receiver.RealReflectType); isTrue {
			receiver.Type = fastReflect.PageList
			break
		}

		// 自定义集合类型
		numField := receiver.RealReflectType.NumField()
		if numField > 0 && receiver.RealReflectType.Field(0).PkgPath == collectionsTypeString {
			receiver.Type = fastReflect.CustomList
			receiver.ItemMeta.ReflectType = types.GetListItemType(receiver.ReflectType)
			receiver.SliceType = reflect.SliceOf(receiver.ItemMeta.ReflectType)
			break
		}

		// 结构体
		if types.IsStruct(receiver.RealReflectType) {
			receiver.Type = fastReflect.Struct
			receiver.NumField = receiver.RealReflectType.NumField()
			break
		}
		receiver.Type = fastReflect.Unknown
	}
}
*/

//// newMetaByType 得到类型的元数据
//func newMetaByType(pointerMeta fastReflect.PointerMeta, parent *valueMeta) *valueMeta {
//	meta := &valueMeta{}
//	if parent != nil {
//		meta.Parent = parent
//		meta.ParentName = parent.Name
//		meta.Level = parent.Level + 1
//		meta.Name = parent.Name + pointerMeta.Name
//	}
//	meta.PointerMeta = pointerMeta
//	meta.IsNil = true
//	meta.IsAnonymous = false
//	meta.IsExported = true
//
//	// 解析类型
//	return meta
//}

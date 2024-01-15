package mapper

import (
	"fmt"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

const (
	mapSplitTag       = ":-:"
	anonymousSplitTag = "" // :anonymous:
)

// 元数据
type valueMeta struct {
	Parent             *valueMeta          // 上级元数据
	ParentName         string              // 上级字段名称
	Name               string              // 字段名称
	RegexPattern       string              // 字段名称匹配规则
	ValueAny           any                 // 值
	ReflectValue       reflect.Value       // 值
	ReflectType        reflect.Type        // 字段类型
	ReflectTypeString  string              // 类型
	RealReflectType    reflect.Type        // 字段去指针后的类型
	ReflectStructField reflect.StructField // 字段类型
	Type               FieldType           // 集合类型
	IsNil              bool                // 是否为nil
	IsAnonymous        bool                // 是否为内嵌类型
	IsExported         bool                // 是否为可导出类型
	IsIgnore           bool                // 是否为忽略字段
	CanInterface       bool                // 是否可以转成Any类型
	IsAddr             bool                // 原类型是否带指针
	Level              int                 // 当前解析的层数（默认为第0层）
	MapKey             reflect.Value       // map key
}

// NewMeta 得到类型的元数据
func NewMeta(reflectValue reflect.Value, parent *valueMeta) *valueMeta {
	isAddr := reflectValue.Kind() == reflect.Pointer
	if isAddr && !reflectValue.IsNil() {
		reflectValue = reflect.Indirect(reflectValue)
	}
	reflectType := reflectValue.Type()
	meta := NewMetaByType(reflectType, parent)
	meta.setReflectValue(reflectValue)
	meta.IsAddr = isAddr
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
	meta.IsAddr = reflectType.Kind() == reflect.Pointer

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

	switch parent.Type {
	case Slice:
		mt.Name = parent.Name + fmt.Sprintf("[%s]", field.Name)
	case Map, Dic:
		mt.Name = parent.Name + fmt.Sprintf("{%s}", field.Name)
	default:
		// 内嵌字段类型的Name为类型名称，这里用标记代替
		if field.Anonymous {
			mt.Name = parent.Name + anonymousSplitTag
		} else {
			mt.Name = parent.Name + field.Name
		}
	}
	mt.RegexPattern = mt.Name
	mt.RegexPattern = strings.ReplaceAll(mt.RegexPattern, "{", "(\\{|)")
	mt.RegexPattern = strings.ReplaceAll(mt.RegexPattern, "}", "(\\}|)")

	//mt.RegexPattern = strings.ReplaceAll(mt.RegexPattern, "[", "(\\[|)")
	//mt.RegexPattern = strings.ReplaceAll(mt.RegexPattern, "]", "(\\]|)")
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

	receiver.ReflectTypeString = receiver.RealReflectType.String()

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
		numField := receiver.RealReflectType.NumField()
		if numField > 0 && receiver.RealReflectType.Field(0).PkgPath == "github.com/farseer-go/collections" {
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

// NewReflectValue 左值为指针类型时，需要先初始化
func (receiver *valueMeta) NewReflectValue() {
	if !receiver.ReflectValue.IsValid() {
		// 只能使用reflect.New,否则会出现无法寻址的问题
		receiver.ReflectValue = reflect.New(receiver.RealReflectType).Elem()
		//switch receiver.Type {
		//case Slice:
		//	// 只能使用reflect.New,否则会出现无法寻址的问题
		//	//receiver.ReflectValue = reflect.MakeSlice(receiver.RealReflectType, 0, 0)
		//	receiver.ReflectValue = reflect.New(receiver.RealReflectType).Elem()
		//case Map:
		//	receiver.ReflectValue = reflect.MakeMap(receiver.RealReflectType)
		//default:
		//	receiver.ReflectValue = reflect.New(receiver.RealReflectType).Elem()
		//}
		receiver.setReflectValue(receiver.ReflectValue)
		return
	}

	if types.IsNil(receiver.ReflectValue) {
		switch receiver.Type {
		case Slice:
			receiver.ReflectValue.Set(reflect.MakeSlice(receiver.RealReflectType, 0, 0))
		case Map:
			receiver.ReflectValue.Set(reflect.MakeMap(receiver.RealReflectType))
		default:
			receiver.ReflectValue.Set(reflect.New(receiver.RealReflectType))
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

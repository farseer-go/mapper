package mapper

import (
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/types"
	"reflect"
	"regexp"
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
	useRegex           bool                // 使用了正则
	Regexp             *regexp.Regexp      // 正则
	RegexPattern       string              // 字段名称匹配规则
	ReflectValue       reflect.Value       // 值
	ReflectStructField reflect.StructField // 字段类型
	IsNil              bool                // 是否为nil
	IsAnonymous        bool                // 是否为内嵌类型
	IsIgnore           bool                // 是否为忽略字段
	Level              int                 // 当前解析的层数（默认为第0层）
	MapKey             reflect.Value       // map key
}

// newMeta 得到类型的元数据
func newMetaVal(value reflect.Value) *valueMeta {
	pointerMeta := fastReflect.PointerOf(value.Interface())
	meta := &valueMeta{
		PointerMeta: pointerMeta,
		IsNil:       true,
	}
	meta.setReflectValue(value)
	return meta
}

// newStructField 创建子元数据
func newStructField(value reflect.Value, field reflect.StructField, parent *valueMeta, isBuildRegex bool) *valueMeta {
	pointerMeta := fastReflect.PointerOf(value.Interface())
	mt := &valueMeta{
		Parent:             parent,
		ParentName:         parent.Name,
		Level:              parent.Level + 1,
		PointerMeta:        pointerMeta,
		IsNil:              true,
		ReflectStructField: field,
		IsAnonymous:        field.Anonymous,
	}
	if parent.useRegex {
		mt.useRegex = true
	}
	mt.setReflectValue(value)

	// 定义的标签
	//tags := strings.Split(field.Tag.Get("mapper"), ";")
	//for _, tag := range tags {
	//	if tag == "ignore" {
	//		mt.IsIgnore = true
	//		break
	//	}
	//}

	switch parent.Type {
	case fastReflect.Slice:
		var str strings.Builder
		str.WriteString(parent.Name)
		if len(field.Name) > 0 {
			str.WriteString("[")
			str.WriteString(field.Name)
			str.WriteString("]")
		}
		mt.Name = str.String()
	case fastReflect.Map, fastReflect.Dic:
		var str strings.Builder
		str.WriteString(parent.Name)
		if len(field.Name) > 0 {
			str.WriteString("{")
			str.WriteString(field.Name)
			str.WriteString("}")
			// 简写：if isBuildRegex { mt.useRegex = true }
			mt.useRegex = isBuildRegex
		}
		mt.Name = str.String()
	default:
		var str strings.Builder
		str.WriteString(parent.Name)
		// 内嵌字段类型的Name为类型名称，这里用标记代替
		if field.Anonymous {
			str.WriteString(anonymousSplitTag)
		} else {
			str.WriteString(field.Name)
		}
		mt.Name = str.String()
	}

	// 正则
	if mt.useRegex {
		mt.setRegex()
	} else {
		mt.RegexPattern = mt.Name
	}
	return mt
}

// 设置正则规则
func (receiver *valueMeta) setRegex() {
	switch receiver.Parent.Type {
	case fastReflect.Slice:
		if len(receiver.ReflectStructField.Name) > 0 {
			var str strings.Builder
			str.WriteString(receiver.Parent.RegexPattern)
			str.WriteString("[")
			str.WriteString(receiver.ReflectStructField.Name)
			str.WriteString("]")
			receiver.RegexPattern = str.String()
		} else {
			receiver.RegexPattern = receiver.Name
		}
	case fastReflect.Map, fastReflect.Dic:
		if len(receiver.ReflectStructField.Name) > 0 {
			var str strings.Builder
			str.WriteString(receiver.Parent.RegexPattern)
			str.WriteString("(\\{|)")
			str.WriteString(receiver.ReflectStructField.Name)
			str.WriteString("(\\}|)")
			receiver.RegexPattern = str.String()
		} else {
			receiver.RegexPattern = receiver.Name
		}
	default:
		var str strings.Builder
		str.WriteString(receiver.Parent.RegexPattern)
		if receiver.ReflectStructField.Anonymous {
			str.WriteString(anonymousSplitTag)
		} else {
			str.WriteString(receiver.ReflectStructField.Name)
		}
		receiver.RegexPattern = str.String()
	}
}

func (receiver *valueMeta) setReflectValue(reflectValue reflect.Value) {
	if reflectValue.Kind() == reflect.Pointer && !reflectValue.IsNil() {
		reflectValue = reflectValue.Elem()
	}

	receiver.ReflectValue = reflectValue
	receiver.IsNil = types.IsNil(receiver.ReflectValue)

	// 取真实的类型
	if receiver.ReflectValue.Kind() == reflect.Interface && !receiver.IsNil { // && receiver.CanInterface
		receiver.ReflectValue = receiver.ReflectValue.Elem()
	}
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

package mapper

import (
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/types"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

const (
	mapSplitTag           = ":-:"
	anonymousSplitTag     = "" // :anonymous:
	collectionsTypeString = "github.com/farseer-go/collections"
)

// var cacheRegexp = make(sync.Map[string]*regexp.Regexp)
var cacheRegexp = sync.Map{}

// 元数据
type valueMeta struct {
	fastReflect.PointerMeta
	Parent       *valueMeta     // 上级元数据
	ParentName   string         // 上级字段名称
	Name         string         // 字段名称（由ParentName + 当前名称）
	FieldName    string         // 字段名称（StructField）
	useRegex     bool           // 使用了正则
	Regexp       *regexp.Regexp // 正则
	RegexPattern string         // 字段名称匹配规则
	ReflectValue reflect.Value  // 值
	IsNil        bool           // 是否为nil
	IsAnonymous  bool           // 是否为内嵌类型
	IsIgnore     bool           // 是否为忽略字段
	Level        int            // 当前解析的层数（默认为第0层）
}

// newMeta 得到类型的元数据
func newMetaVal(value reflect.Value) *valueMeta {
	meta := &valueMeta{
		IsNil: true,
	}
	meta.setReflectValue(value)
	meta.PointerMeta = fastReflect.PointerOfValue(meta.ReflectValue)
	return meta
}

// newStructField 创建子元数据
func newStructField(value reflect.Value, field reflect.StructField, parent *valueMeta, isBuildRegex bool) *valueMeta {
	mt := &valueMeta{
		Parent:      parent,
		ParentName:  parent.Name,
		Level:       parent.Level + 1,
		IsAnonymous: field.Anonymous,
		FieldName:   field.Name,
		useRegex:    parent.useRegex,
	}
	mt.setReflectValue(value)
	mt.PointerMeta = fastReflect.PointerOfValue(mt.ReflectValue)

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
		if len(field.Name) > 0 {
			mt.Name = parent.Name + "[" + field.Name + "]"
			mt.useRegex = isBuildRegex
		} else {
			mt.Name = parent.Name
		}

	case fastReflect.Map, fastReflect.Dic:
		if len(field.Name) > 0 {
			mt.Name = parent.Name + "{" + field.Name + "}"
			mt.useRegex = isBuildRegex
		} else {
			mt.Name = parent.Name
		}
	default:
		// 内嵌字段类型的Name为类型名称，这里用标记代替
		if mt.IsAnonymous {
			mt.Name = parent.Name + anonymousSplitTag
		} else {
			mt.Name = parent.Name + field.Name
		}
	}

	// 正则
	if mt.useRegex {
		mt.setRegex()
	}
	return mt
}

// 设置正则规则
func (receiver *valueMeta) setRegex() {
	switch receiver.Parent.Type {
	case fastReflect.Slice:
		if len(receiver.FieldName) > 0 {
			receiver.RegexPattern = receiver.Parent.RegexPattern + "[" + receiver.FieldName + "]"
		} else {
			receiver.RegexPattern = receiver.Name
		}
	case fastReflect.Map, fastReflect.Dic:
		if len(receiver.FieldName) > 0 {
			receiver.RegexPattern = receiver.Parent.RegexPattern + "(\\{|)" + receiver.FieldName + "(\\}|)"
		} else {
			receiver.RegexPattern = receiver.Name
		}
	default:
		if receiver.IsAnonymous {
			receiver.RegexPattern = receiver.Parent.RegexPattern + anonymousSplitTag
		} else {
			receiver.RegexPattern = receiver.Parent.RegexPattern + receiver.FieldName
		}
	}

	// 将正则表达式缓存起来
	if receiver.RegexPattern != "" && strings.Contains(receiver.RegexPattern, "|") {
		expr := "^" + receiver.RegexPattern + "$"
		if reg, isOk := cacheRegexp.Load(expr); isOk {
			receiver.Regexp = reg.(*regexp.Regexp)
			return
		}
		// 正则编译
		receiver.Regexp = regexp.MustCompile(expr)
		cacheRegexp.Store(expr, receiver.Regexp)
	}
}

func (receiver *valueMeta) setReflectValue(reflectValue reflect.Value) {
	if reflectValue.Kind() == reflect.Pointer && !reflectValue.IsNil() {
		reflectValue = reflectValue.Elem()
	}
	receiver.IsNil = types.IsNil(reflectValue)

	// 取真实的类型
	if reflectValue.Kind() == reflect.Interface && !receiver.IsNil { // && receiver.CanInterface
		reflectValue = reflectValue.Elem()
	}
	receiver.ReflectValue = reflectValue
}

// NewReflectValue 左值为指针类型时，需要先初始化
func (receiver *valueMeta) NewReflectValue() {
	if !receiver.ReflectValue.IsValid() {
		//	// 只能使用reflect.New,否则会出现无法寻址的问题
		//	receiver.ReflectValue = reflect.New(receiver.ReflectType).Elem()
		receiver.setReflectValue(receiver.ZeroReflectValueElem)
		return
	}

	if types.IsNil(receiver.ReflectValue) {
		if receiver.IsAddr {
			receiver.ReflectValue.Set(receiver.ZeroReflectValue)
		} else {
			receiver.ReflectValue.Set(receiver.ZeroReflectValueElem)
		}
		//switch receiver.Type {
		//case fastReflect.Slice:
		//	//receiver.ReflectValue.Set(reflect.MakeSlice(receiver.ReflectType, 0, 0))
		//case fastReflect.Map:
		//	//receiver.ReflectValue.Set(reflect.MakeMap(receiver.ReflectType))
		//default:
		//	//receiver.ReflectValue.Set(reflect.New(receiver.ReflectType))
		//}
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

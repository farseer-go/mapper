package mapper

import (
	"github.com/farseer-go/fs/fastReflect"
	"reflect"
	"strings"
)

// 元数据
type valueMeta struct {
	fastReflect.PointerMeta // 类型
	//Id                      uint64         // 唯一ID
	Parent       *valueMeta    // 上级元数据
	FullName     string        // 字段名称（由ParentName + 当前名称）
	FieldName    string        // 字段名称（StructField）
	ReflectValue reflect.Value // 值
	IsNil        bool          // 是否为nil
	IsAnonymous  bool          // 是否为内嵌类型
	IsIgnore     bool          // 是否为忽略字段
	Level        int           // 当前解析的层数（默认为第0层）
}

// newMeta 得到类型的元数据
//
//	func newMetaVal(fromVal reflect.Value) valueMeta {
//		return valueMeta{
//			//Id:           1,
//			ReflectValue: fromVal,
//			IsNil:        false,
//			PointerMeta:  fastReflect.PointerOfValue(fromVal),
//		}
//	}
//
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
func newStructField(reflectValue reflect.Value, field reflect.StructField, parent *valueMeta) valueMeta {
	mt := valueMeta{
		//		Id:          parent.Id*10 + uint64(field.Index[0]),
		Parent:      parent,
		Level:       parent.Level + 1,
		IsAnonymous: field.Anonymous,
		FieldName:   field.Name,
	}
	// 4ms
	mt.setReflectValue(reflectValue)
	// 4ms
	mt.PointerMeta = fastReflect.PointerOfValue(mt.ReflectValue)
	//return valueMeta{}
	// 定义的标签
	//tags := strings.Split(field.Tag.Get("mapper"), ";")
	//for _, tag := range tags {
	//	if tag == "ignore" {
	//		mt.IsIgnore = true
	//		break
	//	}
	//}

	// 1ms
	switch parent.Type {
	case fastReflect.Slice:
		if len(field.Name) > 0 {
			//mt.FullName = parent.FullName + "[" + field.Name + "]"
			var strBuilder strings.Builder
			strBuilder.WriteString(parent.FullName)
			strBuilder.WriteString("[")
			strBuilder.WriteString(field.Name)
			strBuilder.WriteString("]")
			mt.FullName = strBuilder.String()
		} else {
			mt.FullName = parent.FullName
		}

	case fastReflect.Map, fastReflect.Dic:
		if len(field.Name) > 0 {
			mt.FullName = parent.FullName + field.Name
			//mt.FullName = parent.FullName + "{" + field.Name + "}"
			//if parent.FullName != "" {
			//	var strBuilder strings.Builder
			//	strBuilder.WriteString(parent.FullName)
			//	strBuilder.WriteString("{")
			//	strBuilder.WriteString(field.Name)
			//	strBuilder.WriteString("}")
			//	mt.FullName = strBuilder.String()
			//} else {
			//	mt.FullName = parent.FullName + field.Name
			//}
		} else {
			mt.FullName = parent.FullName
		}
	default:
		// 内嵌字段类型的Name为类型名称，这里用标记代替
		if mt.IsAnonymous {
			//mt.FullName = parent.FullName + anonymousSplitTag
			mt.FullName = parent.FullName
		} else {
			mt.FullName = parent.FullName + field.Name
		}
	}

	return mt
}

func (receiver *valueMeta) setReflectValue(reflectValue reflect.Value) {
	kind := reflectValue.Kind()
	if kind == reflect.Pointer {
		if receiver.IsNil = reflectValue.IsNil(); !receiver.IsNil {
			reflectValue = reflectValue.Elem()
			receiver.setReflectValue(reflectValue)
			return
			//kind = reflectValue.Kind()
		}
	}
	switch kind {
	case reflect.Pointer:
		if !receiver.IsNil {
			receiver.IsNil = reflectValue.IsNil()
		}
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Slice:
		receiver.IsNil = reflectValue.IsNil()
	case reflect.Interface:
		receiver.IsNil = reflectValue.IsNil()
		// 取真实的类型
		if !receiver.IsNil { // && receiver.CanInterface
			reflectValue = reflectValue.Elem()
			receiver.setReflectValue(reflectValue)
			return
		}
	default:
		receiver.IsNil = false
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

	//if types.IsNil(receiver.ReflectValue) {
	if receiver.IsNil {
		// 不能使用此缓存的对象，会出现目标结构有同样结构体类型时，出现同样的指针地址
		//if receiver.IsAddr {
		//	receiver.ReflectValue.Set(receiver.ZeroReflectValue)
		//} else {
		//	receiver.ReflectValue.Set(receiver.ZeroReflectValueElem)
		//}
		switch receiver.Type {
		case fastReflect.Slice:
			receiver.ReflectValue.Set(reflect.MakeSlice(receiver.ReflectType, 0, 0))
		case fastReflect.Map:
			receiver.ReflectValue.Set(reflect.MakeMap(receiver.ReflectType))
		default:
			if receiver.IsAddr {
				receiver.ReflectValue.Set(reflect.New(receiver.ReflectType))
			} else {
				receiver.ReflectValue.Set(reflect.New(receiver.ReflectType).Elem())
			}
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

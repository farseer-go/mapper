package mapper

import (
	"reflect"

	"github.com/farseer-go/fs/fastReflect"
)

// 元数据
type valueMeta struct {
	fastReflect.PointerMeta // 类型
	//Id                      uint64         // 唯一ID
	Parent       *valueMeta    // 上级元数据
	FullName     string        // 字段名称（由ParentName + 当前名称）
	FullHash     uint64        // FullName 的 FNV-1a hash，用于 O(1) 查找，避免字符串分配
	FieldName    string        // 字段名称（StructField）
	ReflectValue reflect.Value // 值
	IsNil        bool          // 是否为nil
	IsAnonymous  bool          // 是否为内嵌类型
	IsIgnore     bool          // 是否为忽略字段
	Level        int           // 当前解析的层数（默认为第0层）
}

const fnvOffset64 uint64 = 14695981039346656037
const fnvPrime64 uint64 = 1099511628211

// hashString 对字符串做 FNV-1a hash（inline 友好）
func hashString(s string) uint64 {
	h := fnvOffset64
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= fnvPrime64
	}
	return h
}

// mixHash 将父 hash 与当前字段名混合，生成子路径的 hash
func mixHash(parentHash uint64, sep byte, name string) uint64 {
	h := parentHash
	if sep != 0 {
		h ^= uint64(sep)
		h *= fnvPrime64
	}
	for i := 0; i < len(name); i++ {
		h ^= uint64(name[i])
		h *= fnvPrime64
	}
	return h
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

// fillStructField 将 reflectValue/field/parent 的信息填入已有的 mt（arena 分配后调用）
func fillStructField(mt *valueMeta, reflectValue reflect.Value, field reflect.StructField, parent *valueMeta) {
	mt.Parent = parent
	mt.Level = parent.Level + 1
	mt.IsAnonymous = field.Anonymous
	mt.FieldName = field.Name

	mt.setReflectValue(reflectValue)
	mt.PointerMeta = fastReflect.PointerOfValue(mt.ReflectValue)

	fillFullName(mt, field.Name, parent)
}

// fillStructFieldByIndex 结构体字段的快速路径：利用预缓存的 FieldTypeMetas 跳过 sync.Map 查表
func fillStructFieldByIndex(mt *valueMeta, reflectValue reflect.Value, fieldIdx int, parent *valueMeta) {
	sf := parent.StructField[fieldIdx]
	mt.Parent = parent
	mt.Level = parent.Level + 1
	mt.IsAnonymous = sf.Anonymous
	mt.FieldName = sf.Name

	mt.setReflectValue(reflectValue)

	// 利用 FieldTypeMetas 预缓存的 TypeMeta，完全跳过 sync.Map 查找
	if cachedTypeMeta := parent.FieldTypeMetas[fieldIdx]; cachedTypeMeta != nil {
		mt.PointerMeta = fastReflect.PointerOfValueWithMeta(mt.ReflectValue, cachedTypeMeta)
	} else {
		mt.PointerMeta = fastReflect.PointerOfValue(mt.ReflectValue)
	}

	fillFullName(mt, sf.Name, parent)
}

// fillStructFieldAssign 赋值阶段的 fillStructField：只计算 hash，不分配 FullName 字符串
func fillStructFieldAssign(mt *valueMeta, reflectValue reflect.Value, field reflect.StructField, parent *valueMeta) {
	mt.Parent = parent
	mt.Level = parent.Level + 1
	mt.IsAnonymous = field.Anonymous
	mt.FieldName = field.Name

	mt.setReflectValue(reflectValue)
	mt.PointerMeta = fastReflect.PointerOfValue(mt.ReflectValue)

	fillHashOnly(mt, field.Name, parent)
}

// fillStructFieldByIndexAssign 赋值阶段的快速路径：只计算 hash，不分配 FullName 字符串
func fillStructFieldByIndexAssign(mt *valueMeta, reflectValue reflect.Value, fieldIdx int, parent *valueMeta) {
	sf := parent.StructField[fieldIdx]
	mt.Parent = parent
	mt.Level = parent.Level + 1
	mt.IsAnonymous = sf.Anonymous
	mt.FieldName = sf.Name

	mt.setReflectValue(reflectValue)

	if cachedTypeMeta := parent.FieldTypeMetas[fieldIdx]; cachedTypeMeta != nil {
		mt.PointerMeta = fastReflect.PointerOfValueWithMeta(mt.ReflectValue, cachedTypeMeta)
	} else {
		mt.PointerMeta = fastReflect.PointerOfValue(mt.ReflectValue)
	}

	fillHashOnly(mt, sf.Name, parent)
}

// fillFullName 计算并设置 FullHash（FNV-1a 增量混合，无字符串分配）和 FullName（供 ContainsSourceKey 使用）
func fillFullName(mt *valueMeta, fieldName string, parent *valueMeta) {
	parentFullName := parent.FullName
	parentHash := parent.FullHash
	parentType := parent.Type

	switch parentType {
	case fastReflect.Slice:
		if len(fieldName) > 0 {
			mt.FullName = parentFullName + "[" + fieldName + "]"
			// hash: parent + '[' + name + ']'
			mt.FullHash = mixHash(mixHash(parentHash, '[', fieldName), ']', "")
		} else {
			mt.FullName = parentFullName
			mt.FullHash = parentHash
		}
	case fastReflect.Map, fastReflect.Dic:
		if len(fieldName) > 0 {
			mt.FullName = parentFullName + fieldName
			if parentHash == 0 {
				mt.FullHash = hashString(fieldName)
			} else {
				mt.FullHash = mixHash(parentHash, 0, fieldName)
			}
		} else {
			mt.FullName = parentFullName
			mt.FullHash = parentHash
		}
	default:
		if mt.IsAnonymous {
			mt.FullName = parentFullName
			mt.FullHash = parentHash
		} else if parentFullName == "" {
			// 根节点直接用 fieldName（无分配，共享字符串常量）
			mt.FullName = fieldName
			mt.FullHash = hashString(fieldName)
		} else {
			mt.FullName = parentFullName + fieldName
			mt.FullHash = mixHash(parentHash, 0, fieldName)
		}
	}
}

// fillHashOnly 仅计算 FullHash，不分配 FullName 字符串
// 用于赋值阶段（assign phase）—— getSourceValue 只需 hash 查找，不需要 FullName
func fillHashOnly(mt *valueMeta, fieldName string, parent *valueMeta) {
	parentHash := parent.FullHash
	parentType := parent.Type

	switch parentType {
	case fastReflect.Slice:
		if len(fieldName) > 0 {
			mt.FullHash = mixHash(mixHash(parentHash, '[', fieldName), ']', "")
		} else {
			mt.FullHash = parentHash
		}
	case fastReflect.Map, fastReflect.Dic:
		if len(fieldName) > 0 {
			if parentHash == 0 {
				mt.FullHash = hashString(fieldName)
			} else {
				mt.FullHash = mixHash(parentHash, 0, fieldName)
			}
		} else {
			mt.FullHash = parentHash
		}
	default:
		if mt.IsAnonymous {
			mt.FullHash = parentHash
		} else if parentHash == 0 {
			mt.FullHash = hashString(fieldName)
		} else {
			mt.FullHash = mixHash(parentHash, 0, fieldName)
		}
	}
}

// newStructField 创建子元数据（供 assignObj 使用，从 heap 分配）
func newStructField(reflectValue reflect.Value, field reflect.StructField, parent *valueMeta) *valueMeta {
	mt := &valueMeta{}
	fillStructField(mt, reflectValue, field, parent)
	return mt
}

func (receiver *valueMeta) setReflectValue(reflectValue reflect.Value) {
	kind := reflectValue.Kind()
	if kind == reflect.Pointer {
		if receiver.IsNil = reflectValue.IsNil(); !receiver.IsNil {
			reflectValue = reflectValue.Elem()
			receiver.setReflectValue(reflectValue)
			return
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
		if !receiver.IsNil {
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
		receiver.setReflectValue(receiver.ZeroReflectValueElem)
		return
	}

	//if types.IsNil(receiver.ReflectValue) {
	if receiver.IsNil {
		// 不能使用此缓存的对象，会出现目标结构有同样结构体类型时，出现同样的指针地址
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

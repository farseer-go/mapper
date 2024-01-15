package mapper

import (
	"github.com/farseer-go/fs/flog"
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
	case Struct:
		receiver.analysisStruct()
	default:
		flog.Warningf("mapper未知的类型解析：%s", receiver.ReflectType.String())
	}
}

// 解析结构体
func (receiver *analysisOjb) analysisStruct() {
	parent := receiver.valueMeta
	// 结构体
	for i := 0; i < parent.ReflectValue.NumField(); i++ {
		numFieldValue := parent.ReflectValue.Field(i)
		numFieldType := parent.RealReflectType.Field(i)

		// 先分析元数据
		receiver.valueMeta = newStructField(numFieldValue, numFieldType, parent)
		receiver.analysisField()
	}
	receiver.valueMeta = parent
}

// 解析map
func (receiver *analysisOjb) analysisMap() {
	parent := receiver.valueMeta
	keyIsGoBasicType := types.IsGoBasicType(receiver.ReflectValue.Type().Key())

	// 遍历map
	miter := parent.ReflectValue.MapRange()
	for miter.Next() {
		mapKey := miter.Key()
		mapValue := miter.Value()
		keyName := mapKey.String()

		// keyName有可能出现<int64>这种值，所以如果是基础类型，再取一次。
		if keyIsGoBasicType {
			keyName = parse.ToString(mapKey.Interface())
		}

		field := reflect.StructField{
			Name:    keyName,
			PkgPath: "",
		}

		// 先分析元数据
		receiver.valueMeta = newStructField(mapValue, field, parent)
		receiver.valueMeta.MapKey = mapKey // 设置MapKey
		receiver.analysisField()
	}

	receiver.valueMeta = parent
}

// 解析字典
func (receiver *analysisOjb) analysisDic() {
	// 转成map
	m := types.GetDictionaryToMap(receiver.ReflectValue)
	// 这里不能用receiver.valueMeta作为父级传入，而必须传receiver.valueMeta.Parent
	// 因为receiver.ReflectStructField是同一个，否则会出现Name名称重复，如：原来是A，变成AA
	receiver.valueMeta = newStructField(m, receiver.ReflectStructField, receiver.valueMeta.Parent)
	// 解析map
	receiver.analysisMap()
	// Dic统一将转成map类型，方便赋值时直接取map，而不用区分类型
	receiver.Type = Map
	receiver.sourceMap[receiver.Name] = receiver.valueMeta
}

// 解析字段
func (receiver *analysisOjb) analysisField() {
	// 不可导出类型，则退出
	if receiver.IsNil || !receiver.IsExported || receiver.Type == Interface {
		return
	}

	// 先完整的赋值（如果目标类型一致，则可以直接取出来，不用分析）
	if receiver.valueMeta.CanInterface {
		receiver.sourceMap[receiver.Name] = receiver.valueMeta
	}

	switch receiver.Type {
	case GoBasicType, CustomList, Interface:
		if receiver.valueMeta.CanInterface {
			receiver.sourceMap[receiver.Name] = receiver.valueMeta
		}
	case Slice:
		receiver.analysisSlice()
	case List:
		receiver.analysisList()
	case Struct:
		receiver.analysisStruct()
	case Dic:
		receiver.analysisDic()
	case Map:
		// 解析map
		receiver.analysisMap()
		receiver.sourceMap[receiver.Name] = receiver.valueMeta
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

// 解析切片
func (receiver *analysisOjb) analysisSlice() {
	parent := receiver.valueMeta

	for i := 0; i < parent.ReflectValue.Len(); i++ {
		sVal := parent.ReflectValue.Index(i)

		field := reflect.StructField{
			Name:    parse.ToString(i),
			PkgPath: "",
		}

		// 先分析元数据
		receiver.valueMeta = newStructField(sVal, field, parent)
		receiver.analysisField()
	}

	receiver.valueMeta = parent
}

// 解析List
func (receiver *analysisOjb) analysisList() {
	// 获取List中的数组元数
	array := types.GetListToArray(receiver.ReflectValue)
	if array != nil {
		receiver.sourceMap[receiver.Name] = NewMeta(reflect.ValueOf(array), receiver.valueMeta)
	}
}

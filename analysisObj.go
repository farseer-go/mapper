package mapper

import (
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

type AnalysisOjb struct {
	*valueMeta                       // 当前元数据
	sourceMap  map[string]*valueMeta // 分析后的结果
}

func (receiver *AnalysisOjb) Analysis(from any) {
	// 定义存储map ,保存解析出来的字段和值
	receiver.sourceMap = make(map[string]*valueMeta)
	// 解析from元数据
	fromValue := reflect.ValueOf(from)
	valMeta := newMetaVal(fromValue)
	receiver.valueMeta = &valMeta
	switch receiver.Type {
	case fastReflect.Map:
		receiver.analysisMap()
	case fastReflect.Struct:
		receiver.analysisStruct()
	default:
		flog.Warningf("mapper未知的类型解析：%s", receiver.ReflectTypeString)
	}
}

// 解析结构体
func (receiver *AnalysisOjb) analysisStruct() {
	parent := receiver.valueMeta
	// 结构体
	for i := 0; i < parent.NumField; i++ {
		if parent.StructField[i].IsExported() {
			// todo 这里的性能要做测试
			numFieldValue := parent.ReflectValue.Field(i)
			// 先分析元数据
			valMeta := newStructField(numFieldValue, parent.StructField[i], parent, true)
			receiver.valueMeta = &valMeta
			receiver.analysisField()
		}
	}
	receiver.valueMeta = parent
}

// 解析map
func (receiver *AnalysisOjb) analysisMap() {
	parent := receiver.valueMeta
	keyIsGoBasicType := receiver.GetKeyMeta().Type == fastReflect.GoBasicType

	// 遍历map
	miter := parent.ReflectValue.MapRange()
	for miter.Next() {
		mapKey := miter.Key()
		mapValue := miter.Value()
		var field reflect.StructField
		// keyName有可能出现<int64>这种值，所以如果是基础类型，再取一次。
		if keyIsGoBasicType {
			field.Name = parse.ToString(mapKey.Interface())
		} else {
			field.Name = mapKey.String()
		}

		// 先分析元数据
		valMeta := newStructField(mapValue, field, parent, true)
		receiver.valueMeta = &valMeta
		receiver.analysisField()
	}

	receiver.valueMeta = parent
}

// 解析字典
func (receiver *AnalysisOjb) analysisDic() {
	// 转成map
	m := types.GetDictionaryToMap(receiver.ReflectValue)
	// 这里不能用receiver.valueMeta作为父级传入，而必须传receiver.valueMeta.Parent
	// 因为receiver.ReflectStructField是同一个，否则会出现Name名称重复，如：原来是A，变成AA
	valMeta := newStructField(m, reflect.StructField{Name: receiver.FieldName, Anonymous: receiver.IsAnonymous}, receiver.valueMeta.Parent, true)
	receiver.valueMeta = &valMeta
	// 解析map
	receiver.analysisMap()
	// Dic统一将转成map类型，方便赋值时直接取map，而不用区分类型
	receiver.Type = fastReflect.Map
	receiver.sourceMap[receiver.Name] = receiver.valueMeta
}

// 解析字段
func (receiver *AnalysisOjb) analysisField() {
	// 不可导出类型，则退出（此行必须先执行，否则会出现指针字段不需要赋值时，被赋值了）
	if receiver.IsNil || receiver.Type == fastReflect.Interface {
		return
	}

	// 先完整的赋值（如果目标类型一致，则可以直接取出来，不用分析）
	receiver.sourceMap[receiver.Name] = receiver.valueMeta

	switch receiver.Type {
	case fastReflect.GoBasicType, fastReflect.Interface:
	// 不需要处理
	case fastReflect.Slice:
		receiver.analysisSlice()
	case fastReflect.List:
		receiver.analysisList()
	case fastReflect.CustomList:
		receiver.analysisList()
	case fastReflect.Struct:
		receiver.analysisStruct()
	case fastReflect.Dic:
		receiver.analysisDic()
	case fastReflect.Map:
		// 解析map
		receiver.analysisMap()
		receiver.sourceMap[receiver.Name] = receiver.valueMeta
		return
	default:
		// 非结构体遍历（好像没用到）
		if strings.Contains(receiver.ParentName, receiver.Name) {
			receiver.sourceMap[receiver.ParentName] = receiver.valueMeta
		} else {
			receiver.sourceMap[receiver.Name] = receiver.valueMeta
		}
	}
}

// 解析切片
func (receiver *AnalysisOjb) analysisSlice() {
	parent := receiver.valueMeta

	length := parent.ReflectValue.Len()
	for i := 0; i < length; i++ {
		sVal := parent.ReflectValue.Index(i)

		field := reflect.StructField{
			Name: parse.ToString(i),
		}

		// 先分析元数据
		valMeta := newStructField(sVal, field, parent, true)
		receiver.valueMeta = &valMeta
		receiver.analysisField()
	}

	receiver.valueMeta = parent
}

// 解析List
func (receiver *AnalysisOjb) analysisList() {
	// 获取List中的数组元数
	array := types.GetListToArrayValue(receiver.ReflectValue)
	if array.Len() > 0 {
		parent := receiver.valueMeta

		valMeta := newStructField(array, reflect.StructField{}, parent, true)
		receiver.valueMeta = &valMeta
		// 分析List中的切片
		receiver.analysisField()

		receiver.valueMeta = parent
	}
}

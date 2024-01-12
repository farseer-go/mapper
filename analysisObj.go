package mapper

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

type analysisOjb struct {
	sourceMap map[string]any // 分析后的结果
}

func (receiver *analysisOjb) analysis(sourceVal reflect.Value) {
	// 定义存储map ,保存解析出来的字段和值
	receiver.sourceMap = make(map[string]any)

	switch sourceVal.Kind() {
	case reflect.Map:
		//mapAnalysis(parentName, fieldName, sourceVal, sourceMap)
		receiver.analysisMap("", sourceVal)
	default:
		// 结构体
		for i := 0; i < sourceVal.NumField(); i++ {
			sourceNumFieldValue := sourceVal.Field(i)
			sourceNumFieldType := sourceVal.Type().Field(i)
			receiver.analysisField(sourceNumFieldType.Name, sourceNumFieldValue, sourceNumFieldType)
		}
	}
}

// 解析map
func (receiver *analysisOjb) analysisMap(parentName string, sourceVal reflect.Value) {
	for _, key := range sourceVal.MapKeys() {
		sourceMapValue := sourceVal.MapIndex(key)
		keyName := key.String()
		if types.IsGoBasicType(key.Type()) {
			keyName = parse.ToString(key.Interface())
		}
		field := reflect.StructField{
			Name:    keyName,
			PkgPath: sourceMapValue.Type().PkgPath(),
		}
		receiver.analysisField(parentName+keyName, sourceMapValue, field)
	}
}

// 解析map/struct的字段
func (receiver *analysisOjb) analysisField(parentName string, sourceFieldValue reflect.Value, sourceFieldType reflect.StructField) {
	if types.IsNil(sourceFieldValue) {
		return
	}

	sourceFieldValue = reflect.Indirect(sourceFieldValue)
	sourceFieldValueType := sourceFieldValue.Type()

	// 取真实的类型
	if sourceFieldValueType.Kind() == reflect.Interface && sourceFieldValue.CanInterface() {
		sourceFieldValue = sourceFieldValue.Elem()
		sourceFieldValueType = sourceFieldValue.Type()
	}

	// 不可导出类型，则退出
	if sourceFieldValueType.Kind() == reflect.Interface || !sourceFieldType.IsExported() {
		return
	}

	// 是否为集合
	if _, isList := types.IsList(sourceFieldValue); isList {
		array := types.ListToArray(sourceFieldValue)
		receiver.sourceMap[sourceFieldType.Name] = array
		return
	}

	// 自定义类型
	if len(sourceFieldValueType.String()) > 8 && sourceFieldValueType.String()[len(sourceFieldValueType.String())-8:] == "ListType" {
		receiver.sourceMap[sourceFieldType.Name] = sourceFieldValue.Interface()
		return
	}

	// 结构体
	if types.IsStruct(sourceFieldValueType) {
		if strings.Contains(parentName, sourceFieldType.Name) {
			receiver.analysisStruct(sourceFieldType.Anonymous, parentName, parentName, sourceFieldValue, sourceFieldType.Type)
		} else {
			receiver.analysisStruct(sourceFieldType.Anonymous, sourceFieldType.Name, sourceFieldType.Name, sourceFieldValue, sourceFieldType.Type)
		}
		return
	}

	// map
	if sourceFieldValueType.Kind() == reflect.Map {
		// 解析map
		receiver.analysisMap(parentName, sourceFieldValue)

		if strings.Contains(parentName, sourceFieldType.Name) {
			receiver.sourceMap[parentName] = sourceFieldValue.Interface()
		} else {
			receiver.sourceMap[sourceFieldType.Name] = sourceFieldValue.Interface()
		}
		//mapAnalysis(parentName, sourceFieldType.Name, sourceFieldValue, sourceMap)
		return
	}

	// 非结构体遍历
	itemValue := sourceFieldValue.Interface()
	if strings.Contains(parentName, sourceFieldType.Name) {
		receiver.sourceMap[parentName] = itemValue
	} else {
		receiver.sourceMap[sourceFieldType.Name] = itemValue
	}
}

// 解析结构体
func (receiver *analysisOjb) analysisStruct(anonymous bool, parentName string, fieldName string, fromStructVal reflect.Value, fromStructType reflect.Type) {
	if fromStructVal.Kind() == reflect.Pointer {
		fromStructVal = fromStructVal.Elem()
	}
	// 转换完成之后 执行初始化MapperInit方法
	for i := 0; i < fromStructVal.NumField(); i++ {
		fieldVal := fromStructVal.Field(i)
		itemType := fieldVal.Type()
		sourceFieldName := fromStructVal.Type().Field(i).Name
		itemName := fieldName + sourceFieldName
		if anonymous {
			itemName = "anonymous_" + sourceFieldName
		}
		// 跳过字段值为nil的字段
		if types.IsNil(fieldVal) {
			continue
		}
		if types.IsDateTime(itemType) {
			receiver.sourceMap[itemName] = fieldVal.Interface()
		} else if types.IsGoBasicType(itemType) {
			if fieldVal.CanInterface() {
				receiver.sourceMap[itemName] = fieldVal.Interface()
			}
			// map
		} else if itemType.Kind() == reflect.Map {
			receiver.mapAnalysis(parentName, fieldName, fieldVal)
			// struct
		} else if types.IsStruct(itemType) {
			receiver.analysisStruct(anonymous, parentName, sourceFieldName, fieldVal, fromStructType.Field(i).Type)
		} else if itemType.Kind() == reflect.Slice {
			if fieldVal.CanInterface() {
				receiver.sourceMap[itemName] = fieldVal.Interface()
			}
		} else {
			if fieldVal.CanInterface() || itemType.String() == "decimal.Decimal" {
				receiver.sourceMap[itemName] = fieldVal.Interface()
			}
		}
	}
}

// map 解析
func (receiver *analysisOjb) mapAnalysis(parentName string, fieldName string, fieldVal reflect.Value) {
	newMaps := make(map[string]string)
	maps := fieldVal.MapRange()
	for maps.Next() {
		str := fmt.Sprintf("%v=%v", maps.Key(), maps.Value())
		array := strings.Split(str, "=")
		newMaps[array[0]] = array[1]
	}
	dic := collections.NewDictionaryFromMap(newMaps)
	receiver.sourceMap[parentName] = dic
}

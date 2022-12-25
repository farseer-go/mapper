package mapper

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

// Auto 对象相互转换
func Auto(from, to any) error {
	ts := reflect.TypeOf(to)
	//判断是否指针
	if ts.Kind() != reflect.Ptr {
		return fmt.Errorf("toDTO must be a struct pointer")
	}
	// 反射来源对象
	fsVal := reflect.Indirect(reflect.ValueOf(from))
	// 定义存储map ,保存解析出来的字段和值
	objMap := make(map[string]any)
	// 遍历来源对象
	analysis(fsVal, objMap)
	//转换对象赋值操作
	//反射转换对象 to 指针使用Elem 获取具体值
	tsVal := reflect.ValueOf(to).Elem()
	//赋值操作
	assignment(tsVal, objMap)
	return nil
}

// 赋值操作
func assignment(tsVal reflect.Value, objMap map[string]interface{}) {
	for i := 0; i < tsVal.NumField(); i++ {
		//获取单个字段类型
		item := tsVal.Field(i).Type()
		f := tsVal.Type().Field(i)
		name := f.Name
		objVal := objMap[name]
		//结构体赋值
		if item.Kind() == reflect.Struct {
			var structObj = tsVal.Field(i)
			//list ,pagelist ,dic 转换 ，直接赋值
			if types.IsCollections(structObj.Type()) {
				setVal(objVal, tsVal, f, i)
			} else if types.IsGoBasicType(item) {
				setVal(objVal, tsVal, f, i)
			} else {
				//结构内字段转换 赋值
				setStructVal(structObj, f, tsVal, objMap, i)
			}
		} else {
			//正常字段转换
			setVal(objVal, tsVal, f, i)
		}
	}
}

// 设置值
func setVal(objVal interface{}, tsVal reflect.Value, f reflect.StructField, i int) {
	if objVal != nil {
		objType := reflect.TypeOf(objVal)
		if f.Type.String() == objType.String() {
			tsVal.Field(i).Set(reflect.ValueOf(objVal))
		}
	}
}

// 结构赋值
func setStructVal(structObj reflect.Value, f reflect.StructField, tsVal reflect.Value, objMap map[string]interface{}, i int) {
	for j := 0; j < structObj.NumField(); j++ {
		itemType := structObj.Field(j).Type()
		name := f.Name + f.Type.Field(j).Name
		objVal := objMap[name]
		if objVal == nil {
			continue
		}
		objType := reflect.TypeOf(objVal)
		if itemType.Kind() == objType.Kind() {
			tsVal.Field(i).Field(j).Set(reflect.ValueOf(objVal))
		}
	}
}

// 解析结构体
func analysis(fsVal reflect.Value, objMap map[string]interface{}) {
	for i := 0; i < fsVal.NumField(); i++ {
		itemType := fsVal.Field(i).Type()
		fieldName := fsVal.Type().Field(i).Name

		// 结构体遍历
		if itemType.Kind() == reflect.Struct && !types.IsGoBasicType(itemType) {
			structAnalysis(fieldName, fsVal.Field(i), fsVal.Type().Field(i).Type, objMap)
		} else {
			// 非结构体遍历
			itemValue := fsVal.Field(i).Interface()
			objMap[fieldName] = itemValue
		}
	}
}

// 结构体递归解析
func structAnalysis(fieldName string, fromStructVal reflect.Value, fromStructType reflect.Type, objMap map[string]interface{}) {
	for i := 0; i < fromStructVal.NumField(); i++ {
		fieldVal := fromStructVal.Field(i)
		itemType := fieldVal.Type()
		// go 基础类型
		if types.IsGoBasicType(itemType) {
			itemName := fieldName + fromStructType.Field(i).Name
			itemValue := fieldVal.Interface()
			objMap[itemName] = itemValue
			// map
		} else if itemType.Kind() == reflect.Map {
			mapAnalysis(fieldName, fieldVal, objMap)
			// struct
		} else if itemType.Kind() == reflect.Struct {
			structAnalysis(fromStructType.Field(i).Name, fieldVal, fromStructType.Field(i).Type, objMap)
		}
	}
}

// map 解析
func mapAnalysis(fieldName string, fieldVal reflect.Value, objMap map[string]interface{}) {
	newMaps := make(map[string]string)
	maps := fieldVal.MapRange()
	for maps.Next() {
		str := fmt.Sprintf("%v=%v", maps.Key(), maps.Value())
		array := strings.Split(str, "=")
		newMaps[array[0]] = array[1]
	}
	dic := collections.NewDictionaryFromMap(newMaps)
	objMap[fieldName] = dic
}

// StructToMap 结构转map
func StructToMap(fromObjPtr any, dic any) error {
	ts := reflect.TypeOf(fromObjPtr)
	if ts.Kind() != reflect.Ptr {
		return fmt.Errorf("toDTO must be a struct pointer")
	}
	fsVal := reflect.ValueOf(fromObjPtr).Elem()
	for i := 0; i < fsVal.NumField(); i++ {
		itemName := fsVal.Type().Field(i).Name
		itemValue := fsVal.Field(i)
		reflect.ValueOf(dic).SetMapIndex(reflect.ValueOf(itemName), itemValue)
	}
	return nil
}

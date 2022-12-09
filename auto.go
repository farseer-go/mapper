package mapper

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

// 对象相互转换
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
	for i := 0; i < fsVal.NumField(); i++ {
		itemType := fsVal.Field(i).Type()
		fieldName := fsVal.Type().Field(i).Name

		// 结构体遍历
		if itemType.Kind() == reflect.Struct && !types.IsGoBasicType(itemType) {
			mapRecursion(fieldName, fsVal.Field(i), fsVal.Type().Field(i).Type, objMap)
		} else {
			// 非结构体遍历
			itemValue := fsVal.Field(i).Interface()
			objMap[fieldName] = itemValue
		}
	}
	//转换对象赋值操作
	//反射转换对象 to 指针使用Elem 获取具体值
	tsVal := reflect.ValueOf(to).Elem()
	for i := 0; i < tsVal.NumField(); i++ {
		//获取单个字段类型
		item := tsVal.Field(i).Type()
		//结构体赋值
		if item.Kind() == reflect.Struct {
			f := tsVal.Type().Field(i)
			var structObj = tsVal.Field(i)
			//list ,pagelist ,dic 转换 ，直接赋值
			if types.IsCollections(structObj.Type()) {
				f := tsVal.Type().Field(i)
				name := f.Name
				objVal := objMap[name]
				if objVal == nil {
					continue
				}
				objType := reflect.TypeOf(objVal)
				if f.Type.String() == objType.String() {
					tsVal.Field(i).Set(reflect.ValueOf(objVal))
				}
			} else {
				//结构内字段转换 赋值
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
		} else {
			//正常字段转换
			f := tsVal.Type().Field(i)
			name := f.Name
			objVal := objMap[name]
			if objVal == nil {
				continue
			}
			objType := reflect.TypeOf(objVal)
			if f.Type.String() == objType.String() {
				tsVal.Field(i).Set(reflect.ValueOf(objVal))
			}
		}
	}
	return nil
}

// 结构体递归取值
func mapRecursion(fieldName string, fromStructVal reflect.Value, fromStructType reflect.Type, objMap map[string]interface{}) {
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
			newMaps := make(map[string]string)
			maps := fieldVal.MapRange()
			for maps.Next() {
				str := fmt.Sprintf("%v=%v", maps.Key(), maps.Value())
				array := strings.Split(str, "=")
				newMaps[array[0]] = array[1]
			}
			dic := collections.NewDictionaryFromMap(newMaps)
			objMap[fieldName] = dic

			// struct
		} else if itemType.Kind() == reflect.Struct {
			mapRecursion(fromStructType.Field(i).Name, fieldVal, fromStructType.Field(i).Type, objMap)
		}
	}
}

// 结构转map
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

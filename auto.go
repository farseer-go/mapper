package mapper

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

// 单例实现相互转换
func Auto(from, to any) error {
	ts := reflect.TypeOf(to)
	if ts.Kind() != reflect.Ptr {
		return fmt.Errorf("toDTO must be a struct pointer")
	}
	fsVal := reflect.ValueOf(from)
	objMap := make(map[string]any)
	// 切片类型
	for i := 0; i < fsVal.NumField(); i++ {
		itemType := fsVal.Field(i).Type()
		if itemType.Kind() == reflect.Struct {
			mapRecursion(fsVal.Type().Field(i).Name, fsVal.Field(i), fsVal.Type().Field(i).Type, objMap)
		} else {
			itemName := fsVal.Type().Field(i).Name
			itemValue := fsVal.Field(i).Interface()
			objMap[itemName] = itemValue
		}
	}
	//赋值toStruct
	tsVal := reflect.ValueOf(to).Elem()
	for i := 0; i < tsVal.NumField(); i++ {
		// 在源结构体中查询有数据结构体中相同属性和类型的字段，有则修改其值
		item := tsVal.Field(i).Type()
		if item.Kind() == reflect.Struct {
			f := tsVal.Type().Field(i)
			var structObj = tsVal.Field(i)
			if types.IsCollections(structObj.Type()) {
				f := tsVal.Type().Field(i)
				name := f.Name
				objVal := objMap[name]
				if objVal == nil {
					continue
				}
				objType := reflect.TypeOf(objVal)
				//fmt.Println(f.Type.Kind(), objType.Kind())
				if f.Type.String() == objType.String() {
					tsVal.Field(i).Set(reflect.ValueOf(objVal))
				}
			} else {
				for j := 0; j < structObj.NumField(); j++ {
					itemType := structObj.Field(j).Type()
					name := f.Name + f.Type.Field(j).Name
					objVal := objMap[name]
					if objVal == nil {
						continue
					}
					objType := reflect.TypeOf(objVal)
					//fmt.Println(f.Type.Kind(), objType.Kind())
					if itemType.Kind() == objType.Kind() {
						tsVal.Field(i).Field(j).Set(reflect.ValueOf(objVal))
					}
				}
			}
		} else {
			f := tsVal.Type().Field(i)
			name := f.Name
			objVal := objMap[name]
			if objVal == nil {
				continue
			}
			objType := reflect.TypeOf(objVal)
			//fmt.Println(f.Type.Kind(), objType.Kind())
			if f.Type.String() == objType.String() {
				tsVal.Field(i).Set(reflect.ValueOf(objVal))
			}
		}

	}
	//fmt.Println(objMap)
	return nil
}
func mapRecursion(fieldName string, fromStructVal reflect.Value, fromStructType reflect.Type, objMap map[string]interface{}) {
	for i := 0; i < fromStructVal.NumField(); i++ {
		itemType := fromStructVal.Field(i).Type()
		if itemType.Kind() == reflect.Struct {
			mapRecursion(fromStructType.Field(i).Name, fromStructVal.Field(i), fromStructType.Field(i).Type, objMap)
		} else if itemType.Kind() == reflect.Map {

			newMaps := make(map[string]string)
			maps := fromStructVal.Field(i).MapRange()
			for maps.Next() {
				str := fmt.Sprintf("%v=%v", maps.Key(), maps.Value())
				array := strings.Split(str, "=")
				newMaps[array[0]] = array[1]
			}
			dic := collections.NewDictionaryFromMap(newMaps)
			objMap[fieldName] = dic
		} else {
			itemName := fieldName + fromStructType.Field(i).Name
			itemValue := fromStructVal.Field(i).Interface()
			objMap[itemName] = itemValue
		}
	}
}

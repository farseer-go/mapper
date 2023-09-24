package mapper

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/parse"
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

	// 转换完成之后 执行初始化MapperInit方法
	defer execInitFunc(reflect.ValueOf(to))

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
func assignment(tsVal reflect.Value, objMap map[string]any) {
	for i := 0; i < tsVal.NumField(); i++ {
		//获取单个字段类型
		fieldType := tsVal.Type().Field(i)
		fieldVal := tsVal.Field(i)
		item := fieldVal.Type()
		objVal := objMap[fieldType.Name]
		//结构体赋值
		if item.Kind() == reflect.Struct && item.String() != "dateTime.DateTime" {
			//list ,pagelist ,dic 转换 ，直接赋值
			if types.IsCollections(fieldVal.Type()) {
				setVal(objVal, fieldVal, fieldType)
			} else if types.IsGoBasicType(item) {
				setVal(objVal, fieldVal, fieldType)
			} else {
				//结构内字段转换 赋值
				setStructVal(fieldType, fieldVal, objMap)
			}

		} else if item.Kind() == reflect.Slice {
			setSliceVal(objVal, fieldVal)
		} else {
			//正常字段转换
			setVal(objVal, fieldVal, fieldType)
		}
	}
}

// 数组设置值
func setSliceVal(objVal any, fieldVal reflect.Value) {
	//获取到具体的值信息
	if objVal != nil {
		//数组对象
		tsValType := fieldVal.Type()
		//要转换的类型
		itemType := tsValType.Elem()
		// 取得数组中元素的类型
		newArr := reflect.MakeSlice(tsValType, 0, 0)
		sliArray := reflect.Indirect(reflect.ValueOf(objVal))
		for i := 0; i < sliArray.Len(); i++ {
			//获取数组内的元素
			structObj := sliArray.Index(i)
			if structObj.Kind() == reflect.Struct {
				newItem := reflect.New(itemType)
				// 要转换对象的值
				newItemField := newItem.Elem()
				for j := 0; j < structObj.NumField(); j++ {
					itemSubValue := structObj.Field(j)
					itemSubType := structObj.Type().Field(j)
					name := itemSubType.Name
					//相同字段赋值
					field := newItemField.FieldByName(name)
					//没有发现相同字段的直接跳过
					if !field.IsValid() {
						continue
					}
					if itemSubValue.Kind() == reflect.Struct {
						setSliceValStruct(itemSubValue, field)
					} else if itemSubValue.Kind() == reflect.Slice {
						setSliceVal(itemSubValue.Interface(), field)
					} else if itemSubValue.Kind() == reflect.Map {
						setSliceValMap(itemSubValue.Interface(), field)
					} else {
						field.Set(itemSubValue)
					}
				}
				newArr = reflect.Append(newArr, newItem.Elem())
			} else {
				newArr = reflect.Append(newArr, structObj)
			}
		}
		fieldVal.Set(newArr)
	}
}

// map赋值
func setSliceValMap(objVal any, fieldVal reflect.Value) {
	fsVal := reflect.ValueOf(objVal).MapRange()
	fieldVal.Set(reflect.MakeMap(fieldVal.Type()))
	for fsVal.Next() {
		k := fsVal.Key()
		v := fsVal.Value()
		item := v.Type()
		value := v.Interface()
		//key := k.Interface()
		//指针类型的值对象
		if v.Type().Kind() == reflect.Pointer {
			vKind := reflect.TypeOf(v.Elem().Interface()).Kind()
			if vKind == reflect.Struct {
				newObj := reflect.New(fieldVal.Type().Elem().Elem())
				setMapValStruct(v.Elem(), newObj)
				fieldVal.SetMapIndex(k, newObj)
			} else if vKind == reflect.Slice {
				setSliceVal(v.Elem().Interface(), fieldVal)
			} else if vKind == reflect.Map {
				setSliceValMap(v.Elem().Interface(), fieldVal)
			} else {
				fieldVal.Set(v.Elem())
			}
		} else {
			//非指针类型的
			if item.Kind() == reflect.Struct {
				setSliceValStruct(reflect.ValueOf(value), fieldVal)
			} else if item.Kind() == reflect.Slice {
				setSliceVal(value, fieldVal)
			} else if item.Kind() == reflect.Map {
				setSliceValMap(value, fieldVal)
			} else {
				fieldVal.Set(v)
			}
		}

	}
}

// Map结构赋值
func setMapValStruct(fieldVal reflect.Value, fields reflect.Value) {
	for j := 0; j < fieldVal.NumField(); j++ {
		itemSubValue := fieldVal.Field(j)
		itemSubType := fieldVal.Type().Field(j)
		name := itemSubType.Name
		//相同字段赋值
		field := fields.Elem().FieldByName(name)
		if itemSubValue.Kind() == reflect.Struct {
			setMapValStruct(itemSubValue, fields)
		} else if itemSubValue.Kind() == reflect.Slice {
			setSliceVal(itemSubValue.Elem(), fields)
		} else {
			field.Set(itemSubValue)
		}
	}
	// 转换完成之后 执行初始化MapperInit方法

	defer execInitFunc(fieldVal.Addr())
	//defer execInitFunc(reflect.ValueOf(tsVal.Field(i).Interface()))
}

// 数组结构赋值
func setSliceValStruct(fieldVal reflect.Value, fields reflect.Value) {
	for j := 0; j < fieldVal.NumField(); j++ {
		itemSubValue := fieldVal.Field(j)
		itemSubType := fieldVal.Type().Field(j)
		name := itemSubType.Name
		//相同字段赋值
		field := fields.FieldByName(name)
		if itemSubValue.Kind() == reflect.Struct {
			setSliceValStruct(itemSubValue, fields)
		} else if itemSubValue.Kind() == reflect.Slice {
			setSliceVal(itemSubValue.Elem(), fields)
		} else {
			field.Set(itemSubValue)
		}
	}
	// 转换完成之后 执行初始化MapperInit方法

	defer execInitFunc(fieldVal.Addr())
	//defer execInitFunc(reflect.ValueOf(tsVal.Field(i).Interface()))
}

// 设置值
func setVal(objVal any, fieldVal reflect.Value, fieldType reflect.StructField) {

	if objVal != nil {
		objType := reflect.TypeOf(objVal)
		if fieldType.Type.String() == objType.String() {
			fieldVal.Set(reflect.ValueOf(objVal))
		} else {
			convert := parse.ConvertValue(objVal, fieldType.Type)
			fieldVal.Set(convert)
		}
	}
}

// 结构赋值
func setStructVal(fieldType reflect.StructField, fieldVal reflect.Value, objMap map[string]any) {
	for j := 0; j < fieldVal.NumField(); j++ {
		itemType := fieldVal.Field(j).Type()
		name := fieldType.Name + fieldType.Type.Field(j).Name
		objVal := objMap[name]
		if objVal == nil {
			continue
		}
		objType := reflect.TypeOf(objVal)
		if itemType.Kind() == objType.Kind() {
			fieldVal.Field(j).Set(reflect.ValueOf(objVal))
		}
	}
	// 转换完成之后 执行初始化MapperInit方法

	defer execInitFunc(fieldVal.Addr())
	//defer execInitFunc(reflect.ValueOf(tsVal.Field(i).Interface()))
}

// 解析结构体
func analysis(fsVal reflect.Value, objMap map[string]any) {
	for i := 0; i < fsVal.NumField(); i++ {
		itemType := fsVal.Field(i).Type()
		field := fsVal.Type().Field(i)

		// 结构体遍历
		if itemType.Kind() == reflect.Interface || !field.IsExported() {
			continue
		}
		if itemType.Kind() == reflect.Struct && !types.IsGoBasicType(itemType) && itemType.String() != "dateTime.DateTime" {
			structAnalysis(field.Name, field.Name, fsVal.Field(i), field.Type, objMap)
		} else {
			// 非结构体遍历
			itemValue := fsVal.Field(i).Interface()
			objMap[field.Name] = itemValue
		}
	}
}

// 结构体递归解析
func structAnalysis(parentName string, fieldName string, fromStructVal reflect.Value, fromStructType reflect.Type, objMap map[string]any) {
	// 转换完成之后 执行初始化MapperInit方法
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
			mapAnalysis(parentName, fieldName, fieldVal, objMap)
			// struct
		} else if itemType.Kind() == reflect.Struct {
			structAnalysis(parentName, fromStructType.Field(i).Name, fieldVal, fromStructType.Field(i).Type, objMap)
		}
	}
}

// map 解析
func mapAnalysis(parentName string, fieldName string, fieldVal reflect.Value, objMap map[string]any) {
	newMaps := make(map[string]string)
	maps := fieldVal.MapRange()
	for maps.Next() {
		str := fmt.Sprintf("%v=%v", maps.Key(), maps.Value())
		array := strings.Split(str, "=")
		newMaps[array[0]] = array[1]
	}
	dic := collections.NewDictionaryFromMap(newMaps)
	objMap[parentName] = dic
}

// StructToMap 结构转map
func StructToMap(fromObjPtr any, dic any) error {
	//ts := reflect.TypeOf(fromObjPtr)
	//if ts.Kind() != reflect.Ptr {
	//	return fmt.Errorf("toDTO must be a struct pointer")
	//}
	fsVal := reflect.Indirect(reflect.ValueOf(fromObjPtr))
	for i := 0; i < fsVal.NumField(); i++ {
		itemName := fsVal.Type().Field(i).Name
		itemValue := fsVal.Field(i)
		reflect.ValueOf(dic).SetMapIndex(reflect.ValueOf(itemName), itemValue)
	}
	return nil
}

// execInitFunc map转换完成之后执行 初始化方法
func execInitFunc(cVal reflect.Value) {
	// 是否实现了IMapperInit
	var actionMapperInit = reflect.TypeOf((*core.IMapperInit)(nil)).Elem()
	if actionMapperInit != nil {
		isImplActionMapperInit := cVal.Type().Implements(actionMapperInit)
		if isImplActionMapperInit {
			//执行方法
			cVal.MethodByName("MapperInit").Call([]reflect.Value{})
			return
		}
	}
	actionMapperInit = reflect.TypeOf((core.IMapperInit)(nil))
	if actionMapperInit != nil {
		isImplActionMapperInit := cVal.Type().Implements(actionMapperInit)
		if isImplActionMapperInit {
			//执行方法
			cVal.MethodByName("MapperInit").Call([]reflect.Value{})
		}
	}
}

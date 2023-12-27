package mapper

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
	"time"
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

// 解析结构体
func analysis(fsVal reflect.Value, objMap map[string]any) {
	for i := 0; i < fsVal.NumField(); i++ {
		itemType := fsVal.Field(i).Type()
		field := fsVal.Type().Field(i)

		// 结构体遍历
		if itemType.Kind() == reflect.Interface || !field.IsExported() {
			continue
		}
		if _, isList := types.IsList(fsVal.Field(i)); isList {
			array := types.ListToArray(fsVal.Field(i))
			//toArrayType := types.GetListItemArrayType(itemType)
			//newArr := reflect.MakeSlice(toArrayType, 0, 0)
			//for i := 0; i < len(array); i++ {
			//	item := array[i]
			//	newArr = reflect.AppendSlice(newArr, reflect.ValueOf(item))
			//}
			objMap[field.Name] = array
		} else if len(itemType.String()) > 8 && itemType.String()[len(itemType.String())-8:] == "ListType" {
			objMap[field.Name] = fsVal.Field(i).Interface()
		} else if itemType.Kind() == reflect.Struct && !types.IsGoBasicType(itemType) && itemType.String() != "dateTime.DateTime" && itemType.String() != "decimal.Decimal" {
			structAnalysis(field.Anonymous, field.Name, field.Name, fsVal.Field(i), field.Type, objMap)
		} else {
			// 非结构体遍历
			itemValue := fsVal.Field(i).Interface()
			objMap[field.Name] = itemValue
		}
	}
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
		if _, isList := types.IsList(fieldVal); isList {
			if reflect.ValueOf(objVal).Kind() == reflect.Invalid {
				continue
			}
			sourceType := types.GetListItemArrayType(item)
			newArr := reflect.New(sourceType)
			toList := types.ListNew(item)
			setSliceVal(reflect.ValueOf(objVal), newArr)
			sliArray := reflect.Indirect(newArr)
			for i := 0; i < sliArray.Len(); i++ {
				//获取数组内的元素
				structObj := sliArray.Index(i)
				types.ListAdd(toList, structObj.Interface())
			}
			fieldVal.Set(toList.Elem())
		} else if len(item.String()) > 8 && item.String()[len(item.String())-8:] == "ListType" {
			if reflect.ValueOf(objVal).Kind() == reflect.Invalid {
				continue
			}
			fieldVal.Set(reflect.ValueOf(objVal))
		} else if item.Kind() == reflect.Struct && item.String() != "dateTime.DateTime" && item.String() != "decimal.Decimal" {
			//list ,pagelist ,dic 转换 ，直接赋值
			if types.IsCollections(fieldVal.Type()) {
				setVal(objVal, fieldVal, fieldType)
			} else if types.IsGoBasicType(item) {
				setVal(objVal, fieldVal, fieldType)
			} else {
				//结构内字段转换 赋值
				setStructVal(fieldType.Anonymous, fieldType, fieldVal, objMap)
			}

		} else if item.Kind() == reflect.Slice {
			if reflect.ValueOf(objVal).Kind() == reflect.Invalid {
				continue
			}
			setSliceVal(reflect.ValueOf(objVal), fieldVal)
		} else {
			//正常字段转换
			setVal(objVal, fieldVal, fieldType)
		}
	}
}

// 数组设置值
func setSliceVal(objVal reflect.Value, fieldVal reflect.Value) {
	//获取到具体的值信息
	//if objVal != nil {
	//数组对象
	tsValType := fieldVal.Type()
	if tsValType.Kind() == reflect.Pointer {
		tsValType = tsValType.Elem()
	}
	if objVal.Kind() == reflect.Pointer {
		objVal = objVal.Elem()
	}
	//要转换的类型
	itemType := tsValType.Elem()
	// 取得数组中元素的类型
	newArr := reflect.MakeSlice(tsValType, 0, 0)
	sliArray := reflect.Indirect(objVal)
	for i := 0; i < sliArray.Len(); i++ {
		//获取数组内的元素
		structObj := sliArray.Index(i)
		if structObj.Type().Kind() == reflect.Interface {
			structObj = structObj.Elem()
		}
		if structObj.Kind() == reflect.Struct && itemType.String() != structObj.Type().String() {
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
					setStruct(itemSubValue, field)
				} else if itemSubValue.Kind() == reflect.Slice {
					setSliceVal(itemSubValue, field)
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
	if fieldVal.Kind() == reflect.Pointer {
		fieldVal.Elem().Set(newArr)
	} else {
		fieldVal.Set(newArr)
	}

	//}
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
				setStruct(v.Elem(), newObj)
				fieldVal.SetMapIndex(k, newObj)
			} else if vKind == reflect.Slice {
				setSliceVal(v.Elem(), fieldVal)
			} else if vKind == reflect.Map {
				setSliceValMap(v.Elem().Interface(), fieldVal)
			} else {
				fieldVal.Set(v.Elem())
			}
		} else {
			//非指针类型的
			if item.Kind() == reflect.Struct {
				newObj := reflect.New(fieldVal.Type().Elem())
				setStruct(v, newObj)
				fieldVal.SetMapIndex(k, newObj.Elem())
			} else if item.Kind() == reflect.Slice {
				setSliceVal(v, fieldVal)
			} else if item.Kind() == reflect.Map {
				setSliceValMap(value, fieldVal)
			} else {
				fieldVal.Set(v)
			}
		}

	}
}

// 数组结构赋值
func setStruct(fieldVal reflect.Value, fields reflect.Value) {
	//list ,pagelist ,dic 转换 ，直接赋值
	if types.IsCollections(fieldVal.Type()) {
		setListVal(fieldVal.Interface(), fields)
	} else {
		for j := 0; j < fieldVal.NumField(); j++ {
			itemSubValue := fieldVal.Field(j)
			itemSubType := fieldVal.Type().Field(j)
			name := itemSubType.Name
			// 指针类型
			if fields.Type().Kind() == reflect.Pointer {
				field := fields.Elem().FieldByName(name)
				if !field.IsValid() {
					continue
				}
				if itemSubValue.Kind() == reflect.Struct {
					setStruct(itemSubValue, fields)
				} else if itemSubValue.Kind() == reflect.Slice {
					setSliceVal(itemSubValue.Elem(), fields)
				} else {
					field.Set(itemSubValue)
				}
			} else {
				// 非指针类型
				field := fields.FieldByName(name)
				if !field.IsValid() {
					continue
				}
				if itemSubValue.Kind() == reflect.Struct {
					setStruct(itemSubValue, fields)
				} else if itemSubValue.Kind() == reflect.Slice {
					setSliceVal(itemSubValue.Elem(), fields)
				} else {
					field.Set(itemSubValue)
				}
			}

		}
	}

	// 转换完成之后 执行初始化MapperInit方法
	if fieldVal.CanAddr() {
		defer execInitFunc(fieldVal.Addr())
	}

	//defer execInitFunc(reflect.ValueOf(tsVal.Field(i).Interface()))
}

// 设置值
func setVal(objVal any, fieldVal reflect.Value, fieldType reflect.StructField) {

	if objVal != nil {
		objType := reflect.TypeOf(objVal)
		if fieldType.Type.String() == objType.String() {
			fieldVal.Set(reflect.ValueOf(objVal))
		} else if fieldType.Type.String() == "string" && objType.String() == "time.Time" {
			stringValue := reflect.ValueOf(objVal).Interface().(time.Time).Format("2006-01-02 15:04:05")
			fieldVal.Set(reflect.ValueOf(stringValue))
		} else if fieldType.Type.String() == "time.Time" && objType.String() == "string" {
			// 将字符串转换为时间类型并赋值给time.Time类型字段
			timeValue, err := time.Parse("2006-01-02 15:04:05", reflect.ValueOf(objVal).Interface().(string))
			if err == nil {
				fieldVal.Set(reflect.ValueOf(timeValue))
			}
		} else {
			convert := parse.ConvertValue(objVal, fieldType.Type)
			fieldVal.Set(convert)
		}
	}
}
func setListVal(objVal any, fieldVal reflect.Value) {
	if objVal != nil {
		objType := reflect.TypeOf(objVal)
		if fieldVal.Type().String() == objType.String() {
			fieldVal.Set(reflect.ValueOf(objVal))
		} else {
			if fieldVal.Type().Kind() == reflect.Struct {
				newObj := reflect.New(fieldVal.Type())
				val := reflect.ValueOf(objVal)
				itemType := val.Field(0) //.Type()
				toInfo := reflect.New(fieldVal.Type().Field(0).Type.Elem())
				setSliceVal(itemType, toInfo)
				if toInfo.Elem().Len() > 0 {
					newObj.Set(toInfo)
				}
				//setStruct(reflect.ValueOf(objVal), newObj.Elem())
				fieldVal.Set(newObj.Elem())
			} else {
				convert := parse.ConvertValue(objVal, fieldVal.Type())
				fieldVal.Set(convert)
			}
		}
	}
}

// 结构赋值
func setStructVal(anonymous bool, fieldType reflect.StructField, fieldVal reflect.Value, objMap map[string]any) {
	for j := 0; j < fieldVal.NumField(); j++ {
		itemType := fieldVal.Field(j).Type()
		name := fieldType.Name + fieldType.Type.Field(j).Name
		if anonymous {
			name = "anonymous_" + fieldType.Type.Field(j).Name
		}
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

// 结构体递归解析
func structAnalysis(anonymous bool, parentName string, fieldName string, fromStructVal reflect.Value, fromStructType reflect.Type, objMap map[string]any) {
	// 转换完成之后 执行初始化MapperInit方法
	for i := 0; i < fromStructVal.NumField(); i++ {
		fieldVal := fromStructVal.Field(i)
		itemType := fieldVal.Type()
		itemName := fieldName + fromStructType.Field(i).Name
		if anonymous {
			itemName = "anonymous_" + fromStructType.Field(i).Name
		}
		if types.IsDateTime(itemType) {
			objMap[itemName] = fieldVal.Interface()
		}
		if types.IsGoBasicType(itemType) {
			if fieldVal.CanInterface() {
				objMap[itemName] = fieldVal.Interface()
			}
			// map
		} else if itemType.Kind() == reflect.Map {
			mapAnalysis(parentName, fieldName, fieldVal, objMap)
			// struct
		} else if itemType.Kind() == reflect.Struct && itemType.String() != "dateTime.DateTime" && itemType.String() != "decimal.Decimal" {
			structAnalysis(anonymous, parentName, fromStructType.Field(i).Name, fieldVal, fromStructType.Field(i).Type, objMap)
		} else if itemType.Kind() == reflect.Slice {
			if fieldVal.CanInterface() {
				objMap[itemName] = fieldVal.Interface()
			}
		} else {
			if fieldVal.CanInterface() || itemType.String() == "decimal.Decimal" {
				objMap[itemName] = fieldVal.Interface()
			}
		}

		//else if itemType.Kind() == reflect.Pointer {
		//	itemName := fieldName
		//	itemValue := fieldVal
		//	objMap[itemName] = itemValue
		//}
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
	dicValue := reflect.ValueOf(dic)
	for i := 0; i < fsVal.NumField(); i++ {
		itemName := fsVal.Type().Field(i).Name
		itemValue := fsVal.Field(i)
		if fsVal.Type().Field(i).Type.Kind() != reflect.Interface {
			dicValue.SetMapIndex(reflect.ValueOf(itemName), itemValue)
		}
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

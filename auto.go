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
	sourceVal := reflect.Indirect(reflect.ValueOf(from))

	// 遍历来源对象
	sourceMap := analysis(sourceVal)
	//转换对象赋值操作
	//反射转换对象 to 指针使用Elem 获取具体值
	targetVal := reflect.ValueOf(to).Elem()
	//赋值操作
	assignment(targetVal, sourceMap)
	return nil
}

func analysis(sourceVal reflect.Value) map[string]any {
	// 定义存储map ,保存解析出来的字段和值
	sourceMap := make(map[string]any)

	switch sourceVal.Kind() {
	case reflect.Map:
		analysisMap(sourceVal, sourceMap)
	default:
		// 结构体
		for i := 0; i < sourceVal.NumField(); i++ {
			sourceNumFieldValue := sourceVal.Field(i)
			sourceNumFieldType := sourceVal.Type().Field(i)
			analysisField(sourceNumFieldValue, sourceNumFieldType, sourceMap)
		}
	}
	return sourceMap
}

func analysisMap(sourceVal reflect.Value, sourceMap map[string]any) {
	for _, key := range sourceVal.MapKeys() {
		sourceMapValue := sourceVal.MapIndex(key)
		field := reflect.StructField{
			Name:    key.String(),
			PkgPath: sourceMapValue.Type().PkgPath(),
		}
		analysisField(sourceMapValue, field, sourceMap)
	}
}

// 解析结构体
func analysisField(sourceFieldValue reflect.Value, sourceFieldType reflect.StructField, sourceMap map[string]any) {
	sourceFieldValueType := sourceFieldValue.Type()
	if sourceFieldValueType.Kind() == reflect.Interface && sourceFieldValue.CanInterface() && !types.IsNil(sourceFieldValue) {
		sourceFieldValueType = sourceFieldValue.Elem().Type()
	}
	// 结构体遍历
	if sourceFieldValueType.Kind() == reflect.Interface || !sourceFieldType.IsExported() {
		return
	}
	// 是否为集合
	if _, isList := types.IsList(sourceFieldValue); isList {
		array := types.ListToArray(sourceFieldValue)
		sourceMap[sourceFieldType.Name] = array
	} else if len(sourceFieldValueType.String()) > 8 && sourceFieldValueType.String()[len(sourceFieldValueType.String())-8:] == "ListType" {
		sourceMap[sourceFieldType.Name] = sourceFieldValue.Interface()
	} else if types.IsStruct(sourceFieldValueType) { // 结构体 sourceFieldValueType.Kind() == reflect.Struct && !types.IsGoBasicType(sourceFieldValueType) && sourceFieldValueType.String() != "dateTime.DateTime" && sourceFieldValueType.String() != "decimal.Decimal" {
		if types.IsNil(sourceFieldValue) {
			return
		}
		if sourceFieldValue.Kind() == reflect.Pointer {
			sourceFieldValue = sourceFieldValue.Elem()
		}
		analysisStruct(sourceFieldType.Anonymous, sourceFieldType.Name, sourceFieldType.Name, sourceFieldValue, sourceFieldType.Type, sourceMap)
	} else {
		// 非结构体遍历
		itemValue := sourceFieldValue.Interface()
		sourceMap[sourceFieldType.Name] = itemValue
	}
}

// 结构体递归解析
func analysisStruct(anonymous bool, parentName string, fieldName string, fromStructVal reflect.Value, fromStructType reflect.Type, sourceMap map[string]any) {
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
		if types.IsDateTime(itemType) {
			sourceMap[itemName] = fieldVal.Interface()
		} else if types.IsGoBasicType(itemType) {
			if fieldVal.CanInterface() {
				sourceMap[itemName] = fieldVal.Interface()
			}
			// map
		} else if itemType.Kind() == reflect.Map {
			mapAnalysis(parentName, fieldName, fieldVal, sourceMap)
			// struct
		} else if itemType.Kind() == reflect.Struct && itemType.String() != "dateTime.DateTime" && itemType.String() != "decimal.Decimal" {
			analysisStruct(anonymous, parentName, sourceFieldName, fieldVal, fromStructType.Field(i).Type, sourceMap)
		} else if itemType.Kind() == reflect.Slice {
			if fieldVal.CanInterface() {
				sourceMap[itemName] = fieldVal.Interface()
			}
		} else {
			if fieldVal.CanInterface() || itemType.String() == "decimal.Decimal" {
				sourceMap[itemName] = fieldVal.Interface()
			}
		}

		//else if itemType.Kind() == reflect.Pointer {
		//	itemName := fieldName
		//	itemValue := fieldVal
		//	objMap[itemName] = itemValue
		//}
	}
}

// 赋值操作
func assignment(targetVal reflect.Value, sourceMap map[string]any) {
	for i := 0; i < targetVal.NumField(); i++ {
		//获取单个字段类型
		targetNumFieldStructField := targetVal.Type().Field(i)
		targetNumFieldValue := targetVal.Field(i)
		targetNumFieldValueType := targetNumFieldValue.Type()
		sourceValue := sourceMap[targetNumFieldStructField.Name]

		// 忽略未导出的字段
		if !targetNumFieldStructField.IsExported() {
			continue
		}
		// 忽略字段
		tags := strings.Split(targetNumFieldStructField.Tag.Get("mapper"), ";")
		for _, tag := range tags {
			if tag == "ignore" {
				continue
			}
		}

		// 目标是指针类型，则先转换成非指针类型
		if targetNumFieldValueType.Kind() == reflect.Pointer {
			targetNumFieldValueType = targetNumFieldValueType.Elem()
		}

		//结构体赋值
		if _, isList := types.IsList(targetNumFieldValue); isList {
			if reflect.ValueOf(sourceValue).Kind() == reflect.Invalid {
				continue
			}
			sourceType := types.GetListItemArrayType(targetNumFieldValueType)
			newArr := reflect.New(sourceType)
			toList := types.ListNew(targetNumFieldValueType)
			setSliceVal(reflect.ValueOf(sourceValue), newArr)
			sliArray := reflect.Indirect(newArr)
			for i := 0; i < sliArray.Len(); i++ {
				//获取数组内的元素
				structObj := sliArray.Index(i)
				types.ListAdd(toList, structObj.Interface())
			}
			targetNumFieldValue.Set(toList.Elem())
		} else if len(targetNumFieldValueType.String()) > 8 && targetNumFieldValueType.String()[len(targetNumFieldValueType.String())-8:] == "ListType" {
			if reflect.ValueOf(sourceValue).Kind() == reflect.Invalid {
				continue
			}
			targetNumFieldValue.Set(reflect.ValueOf(sourceValue))
		} else if targetNumFieldValueType.Kind() == reflect.Slice {
			if reflect.ValueOf(sourceValue).Kind() == reflect.Invalid {
				continue
			}
			setSliceVal(reflect.ValueOf(sourceValue), targetNumFieldValue)
		} else if types.IsCollections(targetNumFieldValue.Type()) { // 集合，//list ,pagelist ,dic 转换 ，直接赋值
			setVal(sourceValue, targetNumFieldValue, targetNumFieldStructField)
		} else if types.IsStruct(targetNumFieldValueType) { // 结构体
			if types.IsGoBasicType(targetNumFieldValueType) {
				setVal(sourceValue, targetNumFieldValue, targetNumFieldStructField)
			} else {
				// 结构内字段转换 赋值
				if types.IsNil(targetNumFieldValue) {
					targetNumFieldValue.Set(reflect.New(targetNumFieldValueType))
					targetNumFieldValue = targetNumFieldValue.Elem()
				}
				setStructVal(targetNumFieldStructField.Anonymous, targetNumFieldStructField, targetNumFieldValue, sourceMap)
			}

		} else {
			//正常字段转换
			setVal(sourceValue, targetNumFieldValue, targetNumFieldStructField)
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
func setVal(sourceValue any, targetFieldValue reflect.Value, targetFieldType reflect.StructField) {
	if sourceValue == nil {
		return
	}
	sourceValueType := reflect.TypeOf(sourceValue)
	// 类型一样
	if targetFieldType.Type.String() == sourceValueType.String() {
		targetFieldValue.Set(reflect.ValueOf(sourceValue))
	} else if targetFieldType.Type.String() == "string" && sourceValueType.String() == "time.Time" { // time.Time转string
		stringValue := reflect.ValueOf(sourceValue).Interface().(time.Time).Format("2006-01-02 15:04:05")
		targetFieldValue.Set(reflect.ValueOf(stringValue))
	} else if targetFieldType.Type.String() == "time.Time" && sourceValueType.String() == "string" { // string转time.Time
		// 将字符串转换为时间类型并赋值给time.Time类型字段
		timeValue, err := time.Parse("2006-01-02 15:04:05", reflect.ValueOf(sourceValue).Interface().(string))
		if err == nil {
			targetFieldValue.Set(reflect.ValueOf(timeValue))
		}
	} else {
		convert := parse.ConvertValue(sourceValue, targetFieldType.Type)
		targetFieldValue.Set(convert)
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
func setStructVal(targetAnonymous bool, targetFieldType reflect.StructField, targetFieldValue reflect.Value, sourceMap map[string]any) {
	for j := 0; j < targetFieldValue.NumField(); j++ {
		itemType := targetFieldValue.Field(j).Type()
		// 目标字段的名称
		targetNumFieldName := targetFieldValue.Type().Field(j).Name
		name := targetFieldType.Name + targetNumFieldName
		objVal := sourceMap[name]
		if targetAnonymous && objVal == nil {
			name = "anonymous_" + targetNumFieldName
			objVal = sourceMap[name]
		}
		if objVal == nil {
			objVal = sourceMap[targetNumFieldName]
		}
		if objVal == nil {
			continue
		}
		objType := reflect.TypeOf(objVal)
		if types.IsTime(itemType) && types.IsDateTime(objType) {
			targetFieldValue.Field(j).Set(parse.ConvertValue(objVal, itemType))
		} else if types.IsDateTime(itemType) && types.IsTime(objType) {
			targetFieldValue.Field(j).Set(parse.ConvertValue(objVal, itemType))
		} else if itemType.Kind() == objType.Kind() {
			targetFieldValue.Field(j).Set(reflect.ValueOf(objVal))
		} else {
			targetFieldValue.Field(j).Set(parse.ConvertValue(objVal, itemType))
		}
	}
	// 转换完成之后 执行初始化MapperInit方法

	defer execInitFunc(targetFieldValue.Addr())
	//defer execInitFunc(reflect.ValueOf(tsVal.Field(i).Interface()))
}

// map 解析
func mapAnalysis(parentName string, fieldName string, fieldVal reflect.Value, sourceMap map[string]any) {
	newMaps := make(map[string]string)
	maps := fieldVal.MapRange()
	for maps.Next() {
		str := fmt.Sprintf("%v=%v", maps.Key(), maps.Value())
		array := strings.Split(str, "=")
		newMaps[array[0]] = array[1]
	}
	dic := collections.NewDictionaryFromMap(newMaps)
	sourceMap[parentName] = dic
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
func execInitFunc(targetFieldValue reflect.Value) {
	// 是否实现了IMapperInit
	var actionMapperInit = reflect.TypeOf((*core.IMapperInit)(nil)).Elem()
	if actionMapperInit != nil {
		isImplActionMapperInit := targetFieldValue.Type().Implements(actionMapperInit)
		if isImplActionMapperInit {
			//执行方法
			targetFieldValue.MethodByName("MapperInit").Call([]reflect.Value{})
			return
		}
	}
	actionMapperInit = reflect.TypeOf((core.IMapperInit)(nil))
	if actionMapperInit != nil {
		isImplActionMapperInit := targetFieldValue.Type().Implements(actionMapperInit)
		if isImplActionMapperInit {
			//执行方法
			targetFieldValue.MethodByName("MapperInit").Call([]reflect.Value{})
		}
	}
}

package mapper

import (
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
	"time"
)

type assignObj struct {
}

// 赋值操作
func (receiver *assignObj) assignment(targetVal reflect.Value, sourceMap map[string]any) {
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
			receiver.setSliceVal(reflect.ValueOf(sourceValue), newArr)
			sliArray := reflect.Indirect(newArr)
			for i := 0; i < sliArray.Len(); i++ {
				//获取数组内的元素
				structObj := sliArray.Index(i)
				types.ListAdd(toList, structObj.Interface())
			}
			targetNumFieldValue.Set(toList.Elem())
			continue
		}

		if len(targetNumFieldValueType.String()) > 8 && targetNumFieldValueType.String()[len(targetNumFieldValueType.String())-8:] == "ListType" {
			if reflect.ValueOf(sourceValue).Kind() == reflect.Invalid {
				continue
			}
			targetNumFieldValue.Set(reflect.ValueOf(sourceValue))
			continue
		}

		if targetNumFieldValueType.Kind() == reflect.Slice {
			if reflect.ValueOf(sourceValue).Kind() == reflect.Invalid {
				continue
			}
			receiver.setSliceVal(reflect.ValueOf(sourceValue), targetNumFieldValue)
			continue
		}

		// 集合，//list ,pagelist ,dic 转换 ，直接赋值
		if types.IsCollections(targetNumFieldValue.Type()) {
			receiver.setVal(sourceValue, targetNumFieldValue, targetNumFieldStructField)
			continue
		}

		// 结构体
		if types.IsStruct(targetNumFieldValueType) {
			if types.IsGoBasicType(targetNumFieldValueType) {
				receiver.setVal(sourceValue, targetNumFieldValue, targetNumFieldStructField)
				continue
			}

			// 目标是否为指针
			if types.IsNil(targetNumFieldValue) {
				// 判断源值是否为nil
				targetFieldName := targetNumFieldStructField.Name
				if targetNumFieldStructField.Anonymous {
					targetFieldName = "anonymous_" + targetFieldName
				}
				sourceHaveVal := false
				for k, _ := range sourceMap {
					if strings.HasPrefix(k, targetFieldName) {
						sourceHaveVal = true
						break
					}
				}
				if sourceHaveVal {
					// 结构内字段转换 赋值。（目标字段是指针结构体，需要先初始化）
					targetNumFieldValue.Set(reflect.New(targetNumFieldValueType))
					targetNumFieldValue = targetNumFieldValue.Elem()

					// 指针类型，只有在源值存在的情况下，才赋值。否则跳过
					receiver.setStructVal(targetNumFieldStructField.Anonymous, targetNumFieldStructField, targetNumFieldValue, sourceMap)
				}
				continue
			}

			// 非指针，正常走逻辑
			receiver.setStructVal(targetNumFieldStructField.Anonymous, targetNumFieldStructField, targetNumFieldValue, sourceMap)
			continue
		}

		//正常字段转换
		receiver.setVal(sourceValue, targetNumFieldValue, targetNumFieldStructField)

	}
}

// 数组设置值
func (receiver *assignObj) setSliceVal(objVal reflect.Value, fieldVal reflect.Value) {
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
					receiver.setStruct(itemSubValue, field)
				} else if itemSubValue.Kind() == reflect.Slice {
					receiver.setSliceVal(itemSubValue, field)
				} else if itemSubValue.Kind() == reflect.Map {
					receiver.setSliceValMap(itemSubValue.Interface(), field)
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
func (receiver *assignObj) setSliceValMap(objVal any, fieldVal reflect.Value) {
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
				receiver.setStruct(v.Elem(), newObj)
				fieldVal.SetMapIndex(k, newObj)
			} else if vKind == reflect.Slice {
				receiver.setSliceVal(v.Elem(), fieldVal)
			} else if vKind == reflect.Map {
				receiver.setSliceValMap(v.Elem().Interface(), fieldVal)
			} else {
				fieldVal.Set(v.Elem())
			}
		} else {
			//非指针类型的
			if item.Kind() == reflect.Struct {
				newObj := reflect.New(fieldVal.Type().Elem())
				receiver.setStruct(v, newObj)
				fieldVal.SetMapIndex(k, newObj.Elem())
			} else if item.Kind() == reflect.Slice {
				receiver.setSliceVal(v, fieldVal)
			} else if item.Kind() == reflect.Map {
				receiver.setSliceValMap(value, fieldVal)
			} else {
				fieldVal.Set(v)
			}
		}

	}
}

// 数组结构赋值
func (receiver *assignObj) setStruct(fieldVal reflect.Value, fields reflect.Value) {
	//list ,pagelist ,dic 转换 ，直接赋值
	if types.IsCollections(fieldVal.Type()) {
		receiver.setListVal(fieldVal.Interface(), fields)
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
					receiver.setStruct(itemSubValue, fields)
				} else if itemSubValue.Kind() == reflect.Slice {
					receiver.setSliceVal(itemSubValue.Elem(), fields)
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
					receiver.setStruct(itemSubValue, fields)
				} else if itemSubValue.Kind() == reflect.Slice {
					receiver.setSliceVal(itemSubValue.Elem(), fields)
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
func (receiver *assignObj) setVal(sourceValue any, targetFieldValue reflect.Value, targetFieldType reflect.StructField) {
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
	} else if sourceValueType.Kind() == reflect.Map {
		targetFieldValue.Set(reflect.ValueOf(sourceValue))
	} else {
		convert := parse.ConvertValue(sourceValue, targetFieldType.Type)
		targetFieldValue.Set(convert)
	}
}

func (receiver *assignObj) setListVal(objVal any, fieldVal reflect.Value) {
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
				receiver.setSliceVal(itemType, toInfo)
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
func (receiver *assignObj) setStructVal(targetAnonymous bool, targetFieldType reflect.StructField, targetFieldValue reflect.Value, sourceMap map[string]any) {
	for j := 0; j < targetFieldValue.NumField(); j++ {
		targetNumFieldValue := targetFieldValue.Field(j)
		targetNumFieldValueType := targetNumFieldValue.Type()
		// 目标字段的名称
		targetNumField := targetFieldValue.Type().Field(j)
		targetNumFieldName := targetNumField.Name
		name := targetFieldType.Name + targetNumFieldName

		if types.IsStruct(targetNumFieldValueType) {
			if targetNumFieldValue.Kind() == reflect.Pointer {
				if types.IsNil(targetNumFieldValue) {
					// 结构内字段转换 赋值。（目标字段是指针结构体，需要先初始化）
					targetNumFieldValue.Set(reflect.New(targetNumFieldValueType.Elem()))
				}
				targetNumFieldValue = targetNumFieldValue.Elem()
			}

			for i := 0; i < targetNumFieldValue.NumField(); i++ {
				itemSubType := targetNumFieldValue.Field(i).Type()
				itemSubNumFieldName := targetNumFieldValue.Type().Field(i).Name
				name = targetNumFieldName + itemSubNumFieldName

				receiver.setFieldValue(targetAnonymous, targetNumFieldValue.Field(i), itemSubType, targetNumField, name, sourceMap)

				//objVal := sourceMap[name]
				//
				//if targetAnonymous && objVal == nil {
				//	name = "anonymous_" + targetNumFieldName
				//	objVal = sourceMap[name]
				//}
				//if objVal == nil {
				//	objVal = sourceMap[targetNumFieldName]
				//}
				//if objVal == nil {
				//	continue
				//}
				//objType := reflect.TypeOf(objVal)
				//if types.IsTime(itemSubType) && types.IsDateTime(objType) {
				//	targetFieldValue.Field(j).Field(i).Set(parse.ConvertValue(objVal, itemSubType))
				//} else if types.IsDateTime(itemSubType) && types.IsTime(objType) {
				//	targetFieldValue.Field(j).Field(i).Set(parse.ConvertValue(objVal, itemSubType))
				//} else if itemSubType.Kind() == objType.Kind() {
				//	targetFieldValue.Field(j).Field(i).Set(reflect.ValueOf(objVal))
				//} else {
				//	targetFieldValue.Field(j).Field(i).Set(parse.ConvertValue(objVal, itemSubType))
				//}
			}
		} else {
			receiver.setFieldValue(targetAnonymous, targetFieldValue.Field(j), targetNumFieldValueType, targetNumField, name, sourceMap)
			//objVal := sourceMap[name]
			//if targetAnonymous && objVal == nil {
			//	name = "anonymous_" + targetNumFieldName
			//	objVal = sourceMap[name]
			//}
			//if objVal == nil {
			//	objVal = sourceMap[targetNumFieldName]
			//}
			//if objVal == nil {
			//	continue
			//}
			//objType := reflect.TypeOf(objVal)
			//if types.IsTime(itemType) && types.IsDateTime(objType) {
			//	targetFieldValue.Field(j).Set(parse.ConvertValue(objVal, itemType))
			//} else if types.IsDateTime(itemType) && types.IsTime(objType) {
			//	targetFieldValue.Field(j).Set(parse.ConvertValue(objVal, itemType))
			//} else if itemType.Kind() == objType.Kind() {
			//	targetFieldValue.Field(j).Set(reflect.ValueOf(objVal))
			//} else {
			//	targetFieldValue.Field(j).Set(parse.ConvertValue(objVal, itemType))
			//}
		}

	}
	// 转换完成之后 执行初始化MapperInit方法

	defer execInitFunc(targetFieldValue.Addr())
	//defer execInitFunc(reflect.ValueOf(tsVal.Field(i).Interface()))
}

func (receiver *assignObj) setFieldValue(targetAnonymous bool, targetFieldValue reflect.Value, targetFieldType reflect.Type, targetNumField reflect.StructField, name string, sourceMap map[string]any) {
	// 忽略未导出的字段
	if !targetNumField.IsExported() {
		return
	}
	// 忽略字段
	tags := strings.Split(targetNumField.Tag.Get("mapper"), ";")
	for _, tag := range tags {
		if tag == "ignore" {
			return
		}
	}
	objVal := sourceMap[name]
	if targetAnonymous && objVal == nil {
		name = "anonymous_" + targetNumField.Name
		objVal = sourceMap[name]
	}
	if objVal == nil {
		objVal = sourceMap[targetNumField.Name]
	}
	if objVal == nil {
		return
	}
	objType := reflect.TypeOf(objVal)
	if types.IsTime(targetFieldType) && types.IsDateTime(objType) {
		targetFieldValue.Set(parse.ConvertValue(objVal, targetFieldType))
	} else if types.IsDateTime(targetFieldType) && types.IsTime(objType) {
		targetFieldValue.Set(parse.ConvertValue(objVal, targetFieldType))
	} else if targetFieldType.Kind() == objType.Kind() {
		targetFieldValue.Set(reflect.ValueOf(objVal))
	} else {
		targetFieldValue.Set(parse.ConvertValue(objVal, targetFieldType))
	}
}

package mapper

import (
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
	"time"
)

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

	//defer execInitFunc(reflect.ReflectValue(tsVal.Field(i).Interface()))
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
				//setStruct(reflect.ReflectValue(objVal), newObj.Elem())
				fieldVal.Set(newObj.Elem())
			} else {
				convert := parse.ConvertValue(objVal, fieldVal.Type())
				fieldVal.Set(convert)
			}
		}
	}
}

// 结构赋值
func (receiver *assignObj) setStructVal(targetAnonymous bool, targetFieldType reflect.StructField, targetFieldValue reflect.Value) {
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
				if targetAnonymous {
					name = itemSubNumFieldName
				}
				receiver.setFieldValue(targetAnonymous, targetNumFieldValue.Field(i), itemSubType, targetNumField, name)

			}
		} else {
			if targetAnonymous {
				name = targetNumFieldName
			}
			receiver.setFieldValue(targetAnonymous, targetFieldValue.Field(j), targetNumFieldValueType, targetNumField, name)
		}

	}
	// 转换完成之后 执行初始化MapperInit方法

	defer execInitFunc(targetFieldValue.Addr())
	//defer execInitFunc(reflect.ReflectValue(tsVal.Field(i).Interface()))
}

func (receiver *assignObj) setFieldValue(targetAnonymous bool, targetFieldValue reflect.Value, targetFieldType reflect.Type, targetNumField reflect.StructField, name string) {
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
	// 实体转换时，多个层级
	objVal := receiver.sourceMap[name]
	// 当匿名字段时，检索map内是否有值
	if targetAnonymous && objVal == nil {
		objVal = receiver.sourceMap["anonymous_"+name]
	}
	// map转实体时，并且实体时匿名字段时走的逻辑
	if objVal == nil {
		objVal = receiver.sourceMap[targetNumField.Name+name]
	}
	if objVal == nil {
		objVal = receiver.sourceMap[targetNumField.Name]
	}
	if objVal == nil {
		return
	}
	switch objVal.Type {
	case GoBasicType:
		targetFieldValue.Set(parse.ConvertValue(objVal.ValueAny, targetFieldType))
	default:
		if targetFieldType.Kind() == objVal.ReflectType.Kind() {
			targetFieldValue.Set(reflect.ValueOf(objVal))
		} else {
			targetFieldValue.Set(parse.ConvertValue(objVal, targetFieldType))
		}
	}
}

package mapper

import (
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

// Array 数组转换
// fromSlice=数组切片
func Array[T any](fromSlice any) []T {
	var toSlice []T
	_ = mapper.MapperSlice(fromSlice, &toSlice)
	return toSlice
}

// 数组转换
func ArrayDOtoDTO[T any](fromDO any) []T {
	var toDTO []T

	//fs := reflect.TypeOf(fromDO)
	//for i := 0; i < fs.Len(); i++ {
	//	item := fs.i
	//}
	return toDTO
}
func AutoMapper(fromDO, toDTO any) error {
	fs := reflect.TypeOf(fromDO)
	if fs.Kind() != reflect.Ptr {
		return fmt.Errorf("fromDO must be a struct pointer")
	}
	ts := reflect.TypeOf(toDTO)
	if ts.Kind() != reflect.Ptr {
		return fmt.Errorf("toDTO must be a struct pointer")
	}
	objMap := structToMap(reflect.ValueOf(fromDO), "")

	//fmt.Println(objMap)
	tsVal := reflect.ValueOf(toDTO).Elem()
	setStructVal(objMap, &tsVal, "")
	return nil
}

// 设置值
func setStructVal(objMap map[string]interface{}, tsVal *reflect.Value, pre string) {
	for i := 0; i < tsVal.NumField(); i++ {
		f := tsVal.Type().Field(i)
		name := pre + f.Name
		cv := tsVal.Field(i)
		var objVal any
		objVal = objMap[name]
		objType := reflect.TypeOf(objVal)
		if objVal != nil && tsVal.Field(i).Kind() == reflect.Struct && objType.Kind() == reflect.Map {
			//dic := collections.NewDictionaryFromMap(objVal.(map[string]interface{}))
			//tsVal.Field(i).Set(reflect.ValueOf(dic))
			continue
		} else if tsVal.Field(i).Kind() == reflect.Struct {
			setStructVal(objMap, &cv, name)
			continue
		}
		if objVal == nil {
			continue
		}
		if f.Type.Kind() == objType.Kind() {
			tsVal.Field(i).Set(reflect.ValueOf(objVal))
		}
	}
}

// struct转map
func structToMap(obj reflect.Value, name string) map[string]any {
	objMap := make(map[string]interface{})
	switch obj.Kind() {
	case reflect.Ptr:
		if name != "" {
			objMap[name] = structToMap(obj.Elem(), name)
		} else {
			objMap = structToMap(obj.Elem(), name)
		}
	case reflect.Struct:
		for i := 0; i < obj.NumField(); i++ {
			f := obj.Field(i)
			name := name + obj.Type().Field(i).Name
			if f.Kind() == reflect.Struct || f.Kind() == reflect.Ptr {
				cMap := structToMap(f, name)
				objMap[name] = cMap
				for mk, mv := range cMap {
					objMap[mk] = mv
				}
			} else {
				if f.CanInterface() {
					objMap[name] = f.Interface()
				} else if f.Kind() == reflect.Map {
					m := make(map[any]any)
					iter := f.MapRange()
					for iter.Next() {
						k := iter.Key()
						v := iter.Value()
						m[k] = v
					}
					objMap[name] = m
				}
			}
		}
	}
	return objMap
}

// 单例实现相互转换
func MapDOtoDTO(from, to any) error {
	// 参数校验
	//fs := reflect.TypeOf(fromDO)
	//if fs.Kind() != reflect.Ptr {
	//	return fmt.Errorf("fromDO must be a struct pointer")
	//}
	ts := reflect.TypeOf(to)
	if ts.Kind() != reflect.Ptr {
		return fmt.Errorf("toDTO must be a struct pointer")
	}
	fsVal := reflect.ValueOf(from)
	objMap := make(map[string]interface{})
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
	//tsObjMap := reflect.ValueOf(objMap)
	tsVal := reflect.ValueOf(to).Elem()
	for i := 0; i < tsVal.NumField(); i++ {
		// 在源结构体中查询有数据结构体中相同属性和类型的字段，有则修改其值
		// name := sTypeOfT.Field(i).Name
		item := tsVal.Field(i).Type()
		if item.Kind() == reflect.Struct {
			f := tsVal.Type().Field(i)
			var structObj = tsVal.Field(i)
			for j := 0; j < structObj.NumField(); j++ {
				itemType := structObj.Field(j).Type()
				name := f.Name + f.Type.Field(j).Name
				if f.Type.Field(j).Name == "source" {
					name = f.Name
					objVal := objMap[name]
					if objVal == nil {
						continue
					}
					maps := reflect.ValueOf(objVal).Field(0)
					mapSliceVal := reflect.ValueOf(objVal).Elem()
					mapSliceType, _ := types.IsMap(mapSliceVal)

					// make....
					mapSliceVal.Set(reflect.MakeMap(mapSliceType))
					//// make....
					for index, key := range maps.MapKeys() {
						val := maps.MapIndex(key)
						//fmt.Println(key)
						//fmt.Println(val)
						// nil说明map不存在这个key
						if val == reflect.ValueOf(nil) {
							val = reflect.MakeSlice(mapSliceType.Elem(), 0, 0)
							mapSliceVal.SetMapIndex(key, val)
						}
						mapValue := reflect.Append(val, reflect.ValueOf(maps.MapKeys()[index]))
						mapSliceVal.SetMapIndex(key, mapValue)

					}
					tsVal.Field(i).Field(j).Set(mapSliceVal)
					//tsVal.Field(i).Elem().Set(mapSliceVal)
					//
					//for _, item := range *objVal.source {
					//	// 生成key
					//	key := reflect.ValueOf(getMapKeyFunc(item))
					//	// 根据key，找到map的value
					//	findMapValue := mapSliceVal.MapIndex(key)
					//	// nil说明map不存在这个key
					//	if findMapValue == reflect.ValueOf(nil) {
					//		findMapValue = reflect.MakeSlice(mapSliceType.Elem(), 0, 0)
					//		mapSliceVal.SetMapIndex(key, findMapValue)
					//	}
					//	mapValue := reflect.Append(findMapValue, reflect.ValueOf(item))
					//	mapSliceVal.SetMapIndex(key, mapValue)
					//}
					//
					//refmap := reflect.MakeMap(itemType)
					//fmt.Println(objVal)
					//
					////for k := 0; k < maps.NumField(); k++ {
					////	key := maps[k]
					////}
					//fmt.Println(maps)
					////tsVal.Field(i).Field(j).Set(maps)
					//tsVal.Field(i).Elem().Set(maps)
				} else {
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
			if f.Type.Kind() == objType.Kind() {
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

// Single 单个转换
// fromObjPtr=struct的指针
func Single[TEntity any](fromObjPtr any) TEntity {
	var toObj TEntity
	_ = mapper.AutoMapper(fromObjPtr, &toObj)
	return toObj
}

// ToMap 结构体转Map
// fromObjPtr=struct的指针
func ToMap[K comparable, V any](fromObjPtr any) map[K]V {
	dic := make(map[K]V)
	_ = mapper.Mapper(fromObjPtr, &dic)
	return dic
}

// ToPageList 转换成core.PageList
// fromSlice=数组切片
func ToPageList[TEntity any](sliceOrListOrListAny any, recordCount int64) collections.PageList[TEntity] {
	lst := ToList[TEntity](sliceOrListOrListAny)
	return collections.NewPageList(lst, recordCount)
}

// ToList 支持：ListAny、List[xx]、[]xx转List[yy]
func ToList[TEntity any](sliceOrListOrListAny any) collections.List[TEntity] {
	sliceOrListOrListAnyValue := reflect.ValueOf(sliceOrListOrListAny)
	if sliceOrListOrListAnyValue.Kind() == reflect.Ptr {
		sliceOrListOrListAnyValue = sliceOrListOrListAnyValue.Elem()
	}
	sliceOrListOrListAnyType := sliceOrListOrListAnyValue.Type()

	// 切片类型
	if sliceOrListOrListAnyValue.Kind() == reflect.Slice {
		var arr []TEntity
		_ = mapper.MapperSlice(sliceOrListOrListAny, &arr)
		return collections.NewList[TEntity](arr...)
	}

	// List类型、ListAny类型
	if strings.HasPrefix(sliceOrListOrListAnyType.String(), "collections.List[") || strings.HasPrefix(sliceOrListOrListAnyType.String(), "collections.ListAny") {
		var arr []TEntity
		items := collections.ReflectToArray(sliceOrListOrListAnyValue)
		_ = mapper.MapperSlice(items, &arr)
		return collections.NewList[TEntity](arr...)
	}
	panic("sliceOrListOrListAny入参必须为切片或collections.List、collections.ListAny集合")
}

// ToListAny 切片转ToListAny
func ToListAny(sliceOrList any) collections.ListAny {
	sliceOrListVal := reflect.ValueOf(sliceOrList)
	if sliceOrListVal.Kind() == reflect.Ptr {
		sliceOrListVal = sliceOrListVal.Elem()
	}
	sliceOrListType := sliceOrListVal.Type()

	// 切片类型
	if sliceOrListVal.Kind() == reflect.Slice || sliceOrListVal.Kind() == reflect.Array {
		lst := collections.NewListAny()
		for i := 0; i < sliceOrListVal.Len(); i++ {
			itemValue := sliceOrListVal.Index(i).Interface()
			lst.Add(itemValue)
		}
		return lst
	}
	if strings.HasPrefix(sliceOrListType.String(), "collections.List[") {
		arr := collections.ReflectToArray(sliceOrListVal)
		return collections.NewListAny(arr...)
	}
	panic("sliceOrList入参必须为切片或collections.List集合")
}

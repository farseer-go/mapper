package mapper

import (
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/farseer-go/collections"
	"reflect"
	"strings"
)

// Array 数组转换
// fromSlice=数组切片
func Array[T any](fromSlice any) []T {
	var toSlice []T
	//获取到具体的值信息
	sliArray := reflect.Indirect(reflect.ValueOf(fromSlice))
	for i := 0; i < sliArray.Len(); i++ {
		item := sliArray.Index(i)
		var tInfo T
		_ = Auto(item.Interface(), &tInfo)
		toSlice = append(toSlice, tInfo)
	}
	return toSlice
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
		//fmt.Println(name)
		if objVal == nil && tsVal.Field(i).Kind() == reflect.Struct {
			setStructVal(objMap, &cv, name)
			continue
		} else if objVal == nil {
			continue
		}
		objType := reflect.TypeOf(objVal)
		//fmt.Println(name, f.Type.Kind(), objType.Kind(), objVal)
		if f.Type.Kind() == objType.Kind() {
			tsVal.Field(i).Set(reflect.ValueOf(objVal))
		} else if tsVal.Field(i).Kind() == reflect.Struct {
			setStructVal(objMap, &cv, name)
			continue
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
			prename := name
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

					newMaps := make(map[string]string)
					maps := f.MapRange()
					for maps.Next() {
						str := fmt.Sprintf("%v=%v", maps.Key(), maps.Value())
						array := strings.Split(str, "=")
						newMaps[array[0]] = array[1]
					}
					dic := collections.NewDictionaryFromMap(newMaps)
					objMap[prename] = dic
					//m := make(map[any]any)
					//iter := f.MapRange()
					//for iter.Next() {
					//	k := iter.Key()
					//	v := iter.Value()
					//	m[k] = v
					//}
					//objMap[name] = m
				}
			}
		}
	}
	return objMap
}

// Single 单个转换
// fromObjPtr=struct的指针
func Single[TEntity any](fromObjPtr any) TEntity {
	var toObj TEntity
	_ = Auto(fromObjPtr, &toObj)
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
		//var arr []TEntity
		arr := Array[TEntity](sliceOrListOrListAny)
		//_ = mapper.MapperSlice(sliceOrListOrListAny, &arr)
		return collections.NewList[TEntity](arr...)
	}

	// List类型、ListAny类型
	if strings.HasPrefix(sliceOrListOrListAnyType.String(), "collections.List[") || strings.HasPrefix(sliceOrListOrListAnyType.String(), "collections.ListAny") {
		//var arr []TEntity
		items := collections.ReflectToArray(sliceOrListOrListAnyValue)
		arr := Array[TEntity](items)
		//_ = mapper.MapperSlice(items, &arr)
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

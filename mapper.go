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
	_ = mapper.MapperSlice(fromSlice, &toSlice)
	return toSlice
}

// 数组转换
func ArrayDOtoDTO[T any](fromDO any) []T {
	var toDTO []T

	fs := reflect.TypeOf(fromDO)
	for i := 0; i < fs.Len(); i++ {
		item := fs.i
	}
	return toDTO
}

// 单例实现相互转换
func MapDOtoDTO(fromDO, toDTO interface{}) error {
	// 参数校验
	fs := reflect.TypeOf(fromDO)
	if fs.Kind() != reflect.Ptr {
		return fmt.Errorf("fromDO must be a struct pointer")
	}
	ts := reflect.TypeOf(toDTO)
	if ts.Kind() != reflect.Ptr {
		return fmt.Errorf("toDTO must be a struct pointer")
	}
	fsVal := reflect.ValueOf(fromDO).Elem()
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
	tsVal := reflect.ValueOf(toDTO).Elem()
	for i := 0; i < tsVal.NumField(); i++ {
		// 在源结构体中查询有数据结构体中相同属性和类型的字段，有则修改其值
		// name := sTypeOfT.Field(i).Name
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

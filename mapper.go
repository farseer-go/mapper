package mapper

import (
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

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

// ToPageList 转换成core.PageList
// fromSlice=数组切片
func ToPageList[TEntity any](fromSlice any, recordCount int64) collections.PageList[TEntity] {
	arr := Array[TEntity](fromSlice)
	return collections.NewPageList(collections.NewList(arr...), recordCount)
}

// ToList ListAny转List泛型
func ToList[TEntity any](source collections.ListAny) collections.List[TEntity] {
	if source.Count() == 0 {
		return collections.NewList[TEntity]()
	}
	toSlice := Array[TEntity](source.ToArray())
	lst := collections.NewList[TEntity](toSlice...)
	return lst
}

// ToListAny 切片转ToListAny
func ToListAny(sliceOrList any) collections.ListAny {
	sliceOrListVal := reflect.ValueOf(sliceOrList)
	if sliceOrListVal.Kind() == reflect.Ptr {
		sliceOrListVal = sliceOrListVal.Elem()
	}
	sliceOrListType := sliceOrListVal.Type()

	lst := collections.NewListAny()
	// 切片类型
	if sliceOrListVal.Kind() == reflect.Slice || sliceOrListVal.Kind() == reflect.Array {
		for i := 0; i < sliceOrListVal.Len(); i++ {
			itemValue := sliceOrListVal.Index(i).Interface()
			lst.Add(itemValue)
		}
		return lst
	}
	if strings.HasPrefix(sliceOrListType.String(), "collections.List[") {
		arrValue := sliceOrListVal.MethodByName("ToArray").Call(nil)[0]
		for i := 0; i < arrValue.Len(); i++ {
			itemValue := arrValue.Index(i)
			lst.Add(itemValue.Interface())
		}
		return lst
	}
	panic("sliceOrList入参必须为切片或collections.List集合")
}

package mapper

import (
	"github.com/devfeel/mapper"
	"github.com/farseer-go/collections"
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

// PageList 转换成core.PageList
// fromSlice=数组切片
func PageList[TEntity any](fromSlice any, recordCount int64) collections.PageList[TEntity] {
	arr := Array[TEntity](fromSlice)
	return collections.NewPageList(collections.NewList(arr...), recordCount)
}

// ToList ListAny转List泛型
func ToList[TEntity any](source collections.ListAny) collections.List[TEntity] {
	toSlice := Array[TEntity](source.ToArray())
	lst := collections.NewList[TEntity](toSlice...)
	return lst
}

package mapper

import (
	"github.com/devfeel/mapper"
	"github.com/farseer-go/collections"
)

// Array 数组转换
func Array[T any](fromSlice any) []T {
	var toSlice []T
	_ = mapper.MapperSlice(fromSlice, &toSlice)
	return toSlice
}

// Single 单个转换
func Single[T any](fromObj any) T {
	var toObj T
	_ = mapper.MapperSlice(fromObj, &toObj)
	return toObj
}

// PageList 转换成core.PageList
func PageList[TData any](fromObj any, recordCount int64) collections.PageList[TData] {
	lst := Array[TData](fromObj)
	return collections.NewPageList(collections.NewList(lst...), recordCount)
}

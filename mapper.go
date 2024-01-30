package mapper

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/fastReflect"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

// Array 数组转换
// fromSlice=数组切片
func Array[T any](fromSlice any) []T {
	// 临时加入埋点
	if container.IsRegister[trace.IManager]() {
		traceHand := container.Resolve[trace.IManager]().TraceHand("mapper.Array")
		defer traceHand.End(nil)
	}

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

// Single 单个转换
func Single[TEntity any](object any) TEntity {
	var toObj TEntity
	_ = Auto(object, &toObj)
	return toObj
}

// ToMap 结构体转Map
func ToMap[K comparable, V any](entity any) map[K]V {
	dic := make(map[K]V)
	_ = StructToMap(entity, dic)
	return dic
}

// ToPageList 转换成core.PageList
// fromSlice=数组切片
func ToPageList[TEntity any](pageList any) collections.PageList[TEntity] {
	list, recordCount := types.GetPageList(pageList)
	lst := ToList[TEntity](list)
	return collections.NewPageList(lst, recordCount)
}

// ToList 支持：ListAny、List[xx]、[]xx转List[yy]
func ToList[TEntity any](sliceOrListOrListAny any) collections.List[TEntity] {
	// 临时加入埋点
	if container.IsRegister[trace.IManager]() {
		traceHand := container.Resolve[trace.IManager]().TraceHand("mapper.ToList")
		defer traceHand.End(nil)
	}

	pointerMeta := fastReflect.PointerOf(sliceOrListOrListAny)
	switch pointerMeta.Type {
	case fastReflect.Slice:
		//var arr []TEntity
		arr := Array[TEntity](sliceOrListOrListAny)
		return collections.NewList[TEntity](arr...)
	case fastReflect.List:
		sliceOrListOrListAnyValue := reflect.ValueOf(sliceOrListOrListAny)
		if sliceOrListOrListAnyValue.Kind() == reflect.Ptr {
			sliceOrListOrListAnyValue = sliceOrListOrListAnyValue.Elem()
		}

		//var arr []TEntity
		items := types.GetListToArray(sliceOrListOrListAnyValue)
		arr := Array[TEntity](items)
		return collections.NewList[TEntity](arr...)
	default:
		//sliceOrListOrListAnyType := sliceOrListOrListAnyValue.Type()
		//toArrayMethod := sliceOrListOrListAnyValue.MethodByName("ToArray")
		//if !toArrayMethod.IsNil() {
		//	//var arr []TEntity
		//	arrValue := toArrayMethod.Call(nil)[0]
		//	var items []TEntity
		//	for i := 0; i < arrValue.Len(); i++ {
		//		item := Single[TEntity](arrValue.Index(i).Interface())
		//		items = append(items, item)
		//	}
		//	return collections.NewList[TEntity](items...)
		//}
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
		arr := types.GetListToArray(sliceOrListVal)
		return collections.NewListAny(arr...)
	}
	panic("sliceOrList入参必须为切片或collections.List集合")
}

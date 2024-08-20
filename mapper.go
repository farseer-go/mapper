package mapper

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/types"
	"reflect"
	"strings"
)

var actionMapperInitAddr = reflect.TypeOf((*core.IMapperInit)(nil)).Elem()

// Single 单个转换
func Single[TEntity any](object any, set ...func(*TEntity)) TEntity {
	itemType := reflect.ValueOf(object)

	var toObj TEntity
	_ = auto(itemType, &toObj, itemType.Type().Implements(actionMapperInitAddr))
	if set != nil {
		set[0](&toObj)
	}
	return toObj
}

// ToList 支持：ListAny、List[xx]、[]xx转List[yy]
func ToList[TEntity any](sliceOrListOrListAny any, set ...func(*TEntity, any)) collections.List[TEntity] {
	sliceOrListOrListAnyValue := reflect.ValueOf(sliceOrListOrListAny)
	kind := sliceOrListOrListAnyValue.Kind()
	if kind == reflect.Ptr {
		sliceOrListOrListAnyValue = sliceOrListOrListAnyValue.Elem()
		kind = sliceOrListOrListAnyValue.Kind()
	}
	switch kind {
	case reflect.Slice:
		arr := Array[TEntity](sliceOrListOrListAny, set...)
		return collections.NewList[TEntity](arr...)
	case reflect.Struct:
		if _, isOk := types.IsList(sliceOrListOrListAnyValue); isOk {
			items := types.GetListToArrayValue(sliceOrListOrListAnyValue)
			arr := arrayByReflectValue[TEntity](items, set...)
			return collections.NewList[TEntity](arr...)
		}
	default:
	}

	panic("sliceOrListOrListAny入参必须为切片或collections.List、collections.ListAny集合")
}

// Array 数组转换
// fromSlice=数组切片
func Array[TEntity any](fromSlice any, set ...func(*TEntity, any)) []TEntity {
	//获取到具体的值信息
	sliArray := reflect.Indirect(reflect.ValueOf(fromSlice))
	return arrayByReflectValue(sliArray, set...)
}

// arrayByReflectValue 数组转换
// fromSlice=数组切片
func arrayByReflectValue[TEntity any](sliArray reflect.Value, set ...func(*TEntity, any)) []TEntity {
	var toSlice []TEntity
	// 元素是否为基础类型
	isGoBasicType := false
	// 元素是否实现了MapperInit
	isImplementsActionMapperInitAddr := false
	arrCount := sliArray.Len()
	if arrCount > 0 {
		itemType := sliArray.Index(0).Type()
		isGoBasicType = types.IsGoBasicType(itemType)
		isImplementsActionMapperInitAddr = itemType.Implements(actionMapperInitAddr)
	}

	// 基础类型
	if isGoBasicType {
		for i := 0; i < arrCount; i++ {
			item := sliArray.Index(i)
			var toObj = (item.Interface()).(TEntity)
			toSlice = append(toSlice, toObj)
		}
		return toSlice
	}

	// BenchmarkSample-12    	 1896634	       633 ns/op	     264 B/op	       7 allocs/op
	// BenchmarkSample-12    	  375018	      3164 ns/op	     264 B/op	       7 allocs/op
	// 复合类型
	for i := 0; i < arrCount; i++ {
		var toObj TEntity
		// BenchmarkSample-12    	   33852	     33642 ns/op	     264 B/op	       7 allocs/op
		item := sliArray.Index(i)
		// 基础类型
		_ = auto(item, &toObj, isImplementsActionMapperInitAddr)
		if set != nil {
			set[0](&toObj, item.Interface())
		}
		toSlice = append(toSlice, toObj)
	}
	return toSlice
}

// ToMap 结构体转Map
func ToMap[K comparable, V any](entity any) map[K]V {
	dic := make(map[K]V)
	_ = structToMap(entity, dic)
	return dic
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

// ToPageList 转换成core.PageList
// fromSlice=数组切片
func ToPageList[TEntity any](pageList any, set ...func(*TEntity, any)) collections.PageList[TEntity] {
	list, recordCount := types.GetPageList(pageList)
	lst := ToList[TEntity](list, set...)
	return collections.NewPageList(lst, recordCount)
}

// structToMap 结构转map
func structToMap(fromObjPtr any, dic any) error {
	fsVal := reflect.Indirect(reflect.ValueOf(fromObjPtr))
	dicValue := reflect.ValueOf(dic)
	for i := 0; i < fsVal.NumField(); i++ {
		itemName := fsVal.Type().Field(i).Name
		itemValue := fsVal.Field(i)
		if fsVal.Type().Field(i).Type.Kind() != reflect.Interface {
			dicValue.SetMapIndex(reflect.ValueOf(itemName), itemValue)
		}
	}
	return nil
}

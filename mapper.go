package mapper

import (
	"reflect"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/types"
)

// Single 单个转换
func Single[TEntity any](object any, set ...func(*TEntity)) TEntity {
	// 快速路径1：类型完全相同
	if obj, ok := object.(TEntity); ok {
		if len(set) > 0 {
			set[0](&obj)
		}
		return obj
	}

	itemType := reflect.ValueOf(object)
	var toObj TEntity
	toObjType := reflect.TypeOf(toObj)

	// 快速路径2：使用类型缓存
	if canUseFastPath(itemType.Type(), toObjType) {
		toObj = fastCopy[TEntity](object)
		if len(set) > 0 {
			set[0](&toObj)
		}
		return toObj
	}

	// 普通路径
	_ = auto(itemType, &toObj)
	if len(set) > 0 {
		set[0](&toObj)
	}
	return toObj
}

// ToList 支持：ListAny、List[xx]、[]xx转List[yy]
func ToList[TEntity any](sliceOrListOrListAny any, set ...func(*TEntity, any)) collections.List[TEntity] {
	sliceOrListOrListAnyValue := reflect.ValueOf(sliceOrListOrListAny)
	if sliceOrListOrListAnyValue.Kind() == reflect.Ptr {
		sliceOrListOrListAnyValue = sliceOrListOrListAnyValue.Elem()
	}
	kind := sliceOrListOrListAnyValue.Kind()

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
	arrCount := sliArray.Len()
	if arrCount == 0 {
		return []TEntity{}
	}

	// 预分配切片容量
	toSlice := make([]TEntity, 0, arrCount)

	// 元素是否为基础类型
	itemType := sliArray.Index(0).Type()
	isGoBasicType := types.IsGoBasicType(itemType)

	// 基础类型
	if isGoBasicType {
		for i := 0; i < arrCount; i++ {
			item := sliArray.Index(i)
			var toObj = (item.Interface()).(TEntity)
			toSlice = append(toSlice, toObj)
		}
		return toSlice
	}

	// 复合类型
	hasSet := len(set) > 0
	for i := 0; i < arrCount; i++ {
		var toObj TEntity
		item := sliArray.Index(i)

		_ = auto(item, &toObj)
		if hasSet {
			itemInterface := item.Interface()
			set[0](&toObj, itemInterface)
		}
		toSlice = append(toSlice, toObj)
	}
	return toSlice
}

// ToMap 结构体转Map
func ToMap[K comparable, V any](entity any) map[K]V {
	fsVal := reflect.Indirect(reflect.ValueOf(entity))
	fieldCount := fsVal.NumField()
	dic := make(map[K]V, fieldCount)
	_ = structToMap(entity, dic)
	return dic
}

// ToListAny 切片转ToListAny
func ToListAny(sliceOrList any) collections.ListAny {
	sliceOrListVal := reflect.ValueOf(sliceOrList)
	if sliceOrListVal.Kind() == reflect.Ptr {
		sliceOrListVal = sliceOrListVal.Elem()
	}

	// 切片类型
	if sliceOrListVal.Kind() == reflect.Slice || sliceOrListVal.Kind() == reflect.Array {
		length := sliceOrListVal.Len()
		lst := collections.NewListAny()
		for i := 0; i < length; i++ {
			itemValue := sliceOrListVal.Index(i).Interface()
			lst.Add(itemValue)
		}
		return lst
	}

	// List类型
	if _, isOk := types.IsList(sliceOrListVal); isOk {
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
		field := fsVal.Type().Field(i)
		itemName := field.Name
		itemValue := fsVal.Field(i)

		// 只处理导出的字段（大写字母开头）
		if field.IsExported() && field.Type.Kind() != reflect.Interface {
			dicValue.SetMapIndex(reflect.ValueOf(itemName), itemValue)
		}
	}
	return nil
}

package mapper

import (
	"fmt"
	"github.com/farseer-go/fs/core"
	"reflect"
)

// Auto 对象相互转换
func Auto(from, to any) error {
	ts := reflect.TypeOf(to)
	//判断是否指针
	if ts.Kind() != reflect.Ptr {
		return fmt.Errorf("toDTO must be a struct pointer")
	}

	// 转换完成之后 执行初始化MapperInit方法
	defer execInitFunc(reflect.ValueOf(to))

	// 遍历来源对象
	var ao analysisOjb
	ao.analysis(from)

	// 赋值
	var so assignObj
	targetVal := reflect.ValueOf(to).Elem()
	so.assignment(targetVal, ao.sourceMap)

	return nil
}

// StructToMap 结构转map
func StructToMap(fromObjPtr any, dic any) error {
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

// execInitFunc map转换完成之后执行 初始化方法
func execInitFunc(targetFieldValue reflect.Value) {
	// 是否实现了IMapperInit
	var actionMapperInit = reflect.TypeOf((*core.IMapperInit)(nil)).Elem()
	if actionMapperInit != nil {
		isImplActionMapperInit := targetFieldValue.Type().Implements(actionMapperInit)
		if isImplActionMapperInit {
			//执行方法
			targetFieldValue.MethodByName("MapperInit").Call([]reflect.Value{})
			return
		}
	}
	actionMapperInit = reflect.TypeOf((core.IMapperInit)(nil))
	if actionMapperInit != nil {
		isImplActionMapperInit := targetFieldValue.Type().Implements(actionMapperInit)
		if isImplActionMapperInit {
			//执行方法
			targetFieldValue.MethodByName("MapperInit").Call([]reflect.Value{})
		}
	}
}

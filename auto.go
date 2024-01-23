package mapper

import (
	"fmt"
	"github.com/farseer-go/fs/core"
	"reflect"
)

// Auto 对象相互转换
func Auto(from, to any) error {
	targetVal := reflect.ValueOf(to)
	//判断是否指针
	if targetVal.Kind() != reflect.Pointer {
		return fmt.Errorf("toDTO must be a struct pointer")
	}

	// 转换完成之后 执行初始化MapperInit方法
	defer execInitFunc(targetVal)

	// 遍历来源对象
	var fAnalysis AnalysisOjb
	fAnalysis.Analysis(from)

	// 赋值
	var tAssign assignObj
	return tAssign.assignment(targetVal, fAnalysis.sourceMap)
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

var actionMapperInitAddr = reflect.TypeOf((*core.IMapperInit)(nil)).Elem()
var actionMapperInit = reflect.TypeOf((core.IMapperInit)(nil))

// execInitFunc map转换完成之后执行 初始化方法
func execInitFunc(targetFieldValue reflect.Value) {
	// 是否实现了IMapperInit
	if actionMapperInitAddr != nil {
		isImplActionMapperInit := targetFieldValue.Type().Implements(actionMapperInitAddr)
		if isImplActionMapperInit {
			//执行方法
			targetFieldValue.MethodByName("MapperInit").Call([]reflect.Value{})
			return
		}
	}
}

package mapper

import (
	"reflect"
)

// Auto 对象相互转换
func Auto(from, to any) error {
	targetVal := reflect.ValueOf(from)
	return auto(targetVal, to, targetVal.Type().Implements(actionMapperInitAddr))
}

// 对象相互转换
func auto(from reflect.Value, target any, isImplementsActionMapperInitAddr bool) error {
	targetVal := reflect.ValueOf(target).Elem()

	// 判断是否指针 外部需保证为指针类型
	//if targetVal.Kind() != reflect.Pointer {
	//	return fmt.Errorf("target must be a struct pointer")
	//}

	// 遍历来源对象
	var fAnalysis analysisOjb
	// BenchmarkSample2-12    	      32	  36,675612 ns/op	39772403 B/op	  212752 allocs/op
	// BenchmarkSample2-12    	      50	  20,069469 ns/op	37280084 B/op	   80000 allocs/op
	sourceSlice := fAnalysis.entry(from)

	// 倒序
	//var sourceSliceDesc []valueMeta
	//for i := len(sourceSlice) - 1; i >= 0; i-- {
	//	sourceSliceDesc = append(sourceSliceDesc, sourceSlice[i])
	//}

	// 赋值
	var tAssign assignObj
	err := tAssign.entry(targetVal, from, sourceSlice)

	/*
		// 转换完成之后 执行初始化MapperInit方法
		if err == nil && isImplementsActionMapperInitAddr {
			// BenchmarkSample-12    	    2606	    460037 ns/op	  960265 B/op	   10007 allocs/op
			targetVal.MethodByName("MapperInit").Call(nil)
			// BenchmarkSample-12    	    2491	    459648 ns/op	  960268 B/op	   10007 allocs/op
			//types.ExecuteMapperInit(targetVal)
		}
	*/
	return err
}

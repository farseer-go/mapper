package mapper

import (
	"fmt"
	"reflect"
)

// Auto 对象相互转换
func Auto(from, target any) error {
	fromVal := reflect.ValueOf(from)
	// 判断是否指针 外部需保证为指针类型
	if reflect.ValueOf(target).Kind() != reflect.Pointer {
		return fmt.Errorf("target must be a struct pointer")
	}

	return auto(fromVal, target)
}

// 对象相互转换
func auto(from reflect.Value, target any) error {
	targetVal := reflect.ValueOf(target).Elem()

	// 判断是否指针 外部需保证为指针类型
	//if targetVal.Kind() != reflect.Pointer {
	//	return fmt.Errorf("target must be a struct pointer")
	//}
	// Benchmark2-12    	 1277652	       933.5 ns/op	    1176 B/op	      11 allocs/op
	// return nil
	// 遍历来源对象
	var fAnalysis analysisOjb
	// defer func() {
	// 	fAnalysis.source = nil
	// 	fAnalysis.fromMeta.Parent = nil
	// 	fAnalysis.fromMeta = valueMeta{}
	// }()
	sourceSlice := fAnalysis.entry(from)
	// Benchmark2-12    	   67772	     17156 ns/op	   21744 B/op	     153 allocs/op
	// return nil

	// 赋值
	var tAssign assignObj
	// defer func() {
	// 	sourceSlice = nil
	// 	tAssign.sourceSlice = nil
	// 	tAssign.valueMeta.Parent = nil
	// 	tAssign.valueMeta = valueMeta{}
	// }()
	err := tAssign.entry(targetVal, from, sourceSlice)

	// // Benchmark2-12    	   41016	     27463 ns/op	   23960 B/op	     207 allocs/op
	// return nil

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

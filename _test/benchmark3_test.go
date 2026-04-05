package test

import (
	"testing"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/data/decimal"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/mapper"
)

// 优化前：1244 ms
// 第1次优化：1115 ms 加入缓存
// 第2次优化：1072 ms 移除ValueAny
// 第3次优化：1038 ms 缓存NumField
// 第4次优化：1038 ms 1075858606 ns/op	688635696 B/op	12366560 allocs/op
// 第5次优化：986 ms 										 Benchmark3-12    	       1	1018,944959 ns/op	581959224 B/op	11226437 allocs/op
// 第6次优化：825 ms 										 Benchmark3-12    	       2	 827,262450 ns/op	544626832 B/op	 9735507 allocs/op
// 第7次优化：602 ms 										 Benchmark3-12    	       2	 607,876002 ns/op	512319544 B/op	 5675569 allocs/op
// 第8次优化：516 ms 										 Benchmark3-12    	       2	 713,063544 ns/op	577931260 B/op	 5045624 allocs/op
// 第9次优化：451 ms 										 Benchmark3-12    	       3	 450,659945 ns/op	340175312 B/op	 3908975 allocs/op
// 第10次优化：439 ms 										 Benchmark3-12    	   	   3	 439,052122 ns/op	327125178 B/op	 3798923 allocs/op
// 第11次优化：281 ms 										 Benchmark3-12    	       4	 281,863772 ns/op	247210794 B/op	 1980033 allocs/op
// 第12次优化：274 ms 										 Benchmark3-12    	       4	 274,476495 ns/op	246091152 B/op	 1990033 allocs/op
// 第13次优化：268 ms arena分配器替代sync.Pool per-object分配  Benchmark3-12    	     46	   268,071261 ns/op	  143713232 B/op   2550028 allocs/op
// 第14次优化：224 ms assignObj也使用arena                    Benchmark3-12    	        50	  223,734870 ns/op	 95542651 B/op	  2120052 allocs/op
// 第15次优化：217 ms sourceMap用pool复用                     Benchmark3-12    	        54	  217,022196 ns/op	 60198918 B/op	  2090062 allocs/op
// 第16次优化：210 ms source slice用pool复用                  Benchmark3-12    	        55	  209,637780 ns/op	 39723320 B/op	  2050050 allocs/op
// 第17次优化：												 Benchmark3-12             6     173,084254 ns/op   32761560 B/op    1690043 allocs/op
func Benchmark3(b *testing.B) {
	lst := collections.NewList[UserVO]()
	float3 := decimal.NewFromInt(3) // 需要时再打开
	for i := 0; i < 10000; i++ {
		lst.Add(UserVO{
			List: collections.NewList(CountVO{Count: 0}, CountVO{Count: 0}, CountVO{Count: 0}, CountVO{Count: 0}, CountVO{Count: 0}, CountVO{Count: 0}),
			Id:   555,
			Name: "aaaa",
			User3: UserVO3{
				Id:    0,
				Name:  "",
				Time2: dateTime.Now(),
				Time3: time.Time{},
				Time:  time.Time{},
				Date:  dateTime.Now(),
				Dec:   float3, // 需要时再打开
				Ts:    333333,
			},
			Array3: []string{"a", "a", "a", "a", "a", "a"},
			Count:  map[int]*CountVO{1: {Count: 0}, 2: {Count: 0}, 3: {Count: 0}, 4: {Count: 0}, 5: {Count: 0}, 6: {Count: 0}},
			Count2: map[int]CountVO{1: {Count: 0}, 2: {Count: 0}, 3: {Count: 0}, 4: {Count: 0}, 5: {Count: 0}, 6: {Count: 0}},
			Count3: map[int]CountVO2{1: {Count: 0}, 2: {Count: 0}, 3: {Count: 0}, 4: {Count: 0}, 5: {Count: 0}, 6: {Count: 0}},
		})
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapper.ToList[UserVO2](lst)
	}
}

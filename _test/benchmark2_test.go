package test

import (
	"testing"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/data/decimal"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/mapper"
)

type CountVO struct {
	Count int // 出现的次数
}

type CountVO2 struct {
	Count int // 出现的次数
}
type UserVO3 struct {
	Id    int64
	Name  string
	Time2 dateTime.DateTime
	Time3 time.Time
	Time  time.Time
	Date  dateTime.DateTime
	Dec   decimal.Decimal // 需要时再打开
	Ts    time.Duration
}
type UserVO struct {
	List   collections.List[CountVO]
	Id     int64
	Name   string
	User3  UserVO3
	Array3 []string
	Count  map[int]*CountVO
	Count2 map[int]CountVO
	Count3 map[int]CountVO2
}
type UserVO2 struct {
	List2  collections.List[CountVO]
	List   collections.List[CountVO]
	Id     int64
	Name   string
	User3  UserVO4
	Array3 []string
	Count  map[int]*CountVO2
	Count2 map[int]CountVO2
	Count3 map[int]CountVO2
}
type UserVO4 struct {
	Id    int64
	Name  string
	Time2 time.Time
	Time3 dateTime.DateTime
	Time  time.Time
	Date  dateTime.DateTime
	Dec   decimal.Decimal
	Ts    time.Duration
}

// 优化前：1244 ms
// 第1次优化：1115 ms 加入缓存
// 第2次优化：1072 ms 移除ValueAny
// 第3次优化：1038 ms 缓存NumField
// 第4次优化：1038 ms 1075858606 ns/op	688635696 B/op	12366560 allocs/op
// 第5次优化：986 ms Benchmark2-12    	       1	1018,944959 ns/op	581959224 B/op	11226437 allocs/op
// 第6次优化：825 ms Benchmark2-12    	       2	 827,262450 ns/op	544626832 B/op	 9735507 allocs/op
// 第7次优化：602 ms Benchmark2-12    	       2	 607,876002 ns/op	512319544 B/op	 5675569 allocs/op
// 第8次优化：516 ms Benchmark2-12    	       2	 713,063544 ns/op	577931260 B/op	 5045624 allocs/op
// 第9次优化：451 ms Benchmark2-12    	       3	 450,659945 ns/op	340175312 B/op	 3908975 allocs/op
// 第10次优化：439 ms Benchmark2-12    	   	   3	 439,052122 ns/op	327125178 B/op	 3798923 allocs/op
// 第11次优化：281 ms Benchmark2-12    	       4	 281,863772 ns/op	247210794 B/op	 1980033 allocs/op
// 第12次优化：274 ms Benchmark2-12    	       4	 274,476495 ns/op	246091152 B/op	 1990033 allocs/op
// Benchmark2-10           					  8    132,386391 ns/op        245971292 B/op   1985036 allocs/op
// Benchmark2-10                 			 10    109,808892 ns/op        174187560 B/op   2564035 allocs/op
// Benchmark2-12    	   41016	     27463 ns/op	   23960 B/op	     207 allocs/op
// Benchmark2-12    	   63397	     17751 ns/op	   21744 B/op	     153 allocs/op
// Benchmark2-12    	   41462	     27493 ns/op	   23848 B/op	     208 allocs/op
// Benchmark2-10           95859         12198 ns/op           23848 B/op        208 allocs/op
// Benchmark2-10          116956         10324 ns/op           16672 B/op        266 allocs/op
func Benchmark2(b *testing.B) {
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

package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/mapper"
	"github.com/shopspring/decimal"
	"testing"
	"time"
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
	Dec   decimal.Decimal
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
	//List2  collections.List[CountVO]
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
// 第5次优化：986 ms BenchmarkMapperToList-12    	       1	1018,944959 ns/op	581959224 B/op	11226437 allocs/op
// 第6次优化：825 ms BenchmarkMapperToList-12    	       2	 827,262450 ns/op	544626832 B/op	 9735507 allocs/op
// 第7次优化：602 ms BenchmarkMapperToList-12    	       2	 607,876002 ns/op	512319544 B/op	 5675569 allocs/op
// 第8次优化：516 ms BenchmarkMapperToList-12    	       2	 713,063544 ns/op	577931260 B/op	 5045624 allocs/op
// 第9次优化：451 ms BenchmarkMapperToList-12    	       3	 450,659945 ns/op	340175312 B/op	 3908975 allocs/op
// 第10次优化：439 ms BenchmarkMapperToList-12    	   3	 439,052122 ns/op	327125178 B/op	 3798923 allocs/op
// 第11次优化：338 ms BenchmarkMapperToList-12    	   3	 344,835250 ns/op	257127912 B/op	 2253367 allocs/op
func Benchmark2(b *testing.B) {
	lst := collections.NewList[UserVO]()
	for i := 0; i < 10000; i++ {
		lst.Add(UserVO{
			List: collections.NewList[CountVO](CountVO{Count: 0}, CountVO{Count: 0}, CountVO{Count: 0}, CountVO{Count: 0}, CountVO{Count: 0}, CountVO{Count: 0}),
			Id:   555,
			Name: "aaaa",
			User3: UserVO3{
				Id:    0,
				Name:  "",
				Time2: dateTime.Now(),
				Time3: time.Time{},
				Time:  time.Time{},
				Date:  dateTime.Now(),
				Dec:   decimal.NewFromFloat32(3),
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

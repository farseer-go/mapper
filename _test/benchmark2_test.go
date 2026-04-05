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
// 第1次优化：Benchmark2-12              54452             0,021936 ns/op            4327 B/op        214 allocs/op
// 第2次优化：Benchmark2-12              65491             0,017990 ns/op            3451 B/op        176 allocs/op
func Benchmark2(b *testing.B) {
	lst := collections.NewList[UserVO]()
	float3 := decimal.NewFromInt(3) // 需要时再打开
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

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapper.ToList[UserVO2](lst)
	}
}

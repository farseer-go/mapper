package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

type SamplePO struct {
	Dec       decimal.Decimal
	ArrayStr  []string
	UserName  string
	CreateAt  time.Time
	IsEnable  bool
	Id        int
	ArrayInt  []int
	Id8       int8
	Id16      int16
	Id32      int32
	Id64      int64
	IdUint    uint
	IdUint8   uint8
	IdUint16  uint16
	IdUint32  uint32
	IdUint64  uint64
	IdFloat32 float32
	IdFloat64 float64
}

// BenchmarkMapperToList-12    	      12	  92,737303 ns/op	99900376 B/op	  565703 allocs/op
func BenchmarkMapperToList(b *testing.B) {
	lst := collections.NewList[SamplePO]()
	for i := 0; i < 10000; i++ {
		lst.Add(SamplePO{
			Dec:       decimal.NewFromFloat(66.88),
			ArrayStr:  []string{"a", "b"},
			UserName:  "UserName",
			CreateAt:  time.Now(),
			IsEnable:  true,
			Id:        1,
			ArrayInt:  []int{3, 4},
			Id8:       8,
			Id16:      16,
			Id32:      32,
			Id64:      64,
			IdUint:    1,
			IdUint8:   8,
			IdUint16:  16,
			IdUint32:  32,
			IdUint64:  64,
			IdFloat32: 32.32,
			IdFloat64: 64.64,
		})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapper.ToList[SamplePO](&lst)
	}
}

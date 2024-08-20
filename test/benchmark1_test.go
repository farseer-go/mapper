package test

import (
	"github.com/farseer-go/mapper"
	"testing"
)

type SamplePO struct {
	UserName  string
	IsEnable  bool
	Id        int
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

// BenchmarkCopyStruct-12    	   14	  81,855216 ns/op	12800284 B/op	  610000 allocs/op （jinzhu）
// BenchmarkSample2-12    	       26	  38,841798 ns/op	39680215 B/op	  100001 allocs/op
func BenchmarkSample1(b *testing.B) {
	po := SamplePO{
		UserName:  "UserName",
		IsEnable:  true,
		Id:        1,
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
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 10000; i++ {
			mapper.Single[SamplePO](po)
		}
	}
}

func BenchmarkSample2(b *testing.B) {
	b.ReportAllocs()

	a := []int{1}
	for i := 0; i < b.N; i++ {
		var b []any
		for _, ai := range a {
			b = append(b, ai)
		}
	}
}

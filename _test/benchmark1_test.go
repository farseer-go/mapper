package test

import (
	"testing"

	"github.com/farseer-go/mapper"
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

// BenchmarkCopyStruct-12    	   	14	   			818.5 ns/op	 		128 B/op	    6 allocs/op （jinzhu）
// BenchmarkSample1-12     	 		25506892        45.63 ns/op           96 B/op       1 allocs/op
// BenchmarkSample1-12     			24262320        48.09 ns/op           96 B/op       1 allocs/op
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
		mapper.Single[SamplePO](po)
	}
}

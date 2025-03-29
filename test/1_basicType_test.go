package test

import (
	"testing"
	"time"

	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/mapper"
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
)

type SamplePO1 struct {
	UserName     string
	IsEnable     bool
	IsEnableStr  string
	IsEnableBool bool
	Id           int
	Id8          int8
	Id16         int16
	Id32         int32
	Id64         int64
	IdUint       uint
	IdUint8      uint8
	IdUint16     uint16
	IdUint32     uint32
	IdUint64     uint64
	IdFloat32    float32
	IdFloat64    float64

	Id8Str  string
	Id16Str string

	Dec      decimal.Decimal // 3ms
	ArrayStr []string
	ArrayInt []int

	// time
	CreateAt time.Time
	UpdateAt dateTime.DateTime
	LastAt   time.Time
	FirstAt  dateTime.DateTime
	Ts       time.Duration
}

type SamplePO2 struct {
	UserName     string
	IsEnable     bool
	IsEnableStr  bool
	IsEnableBool string
	Id           int
	Id8          int8
	Id16         int16
	Id32         int32
	Id64         int64
	IdUint       uint
	IdUint8      uint8
	IdUint16     uint16
	IdUint32     uint32
	IdUint64     uint64
	IdFloat32    float32
	IdFloat64    float64

	Id8Str  int8
	Id16Str int16

	Dec      decimal.Decimal // 3ms
	ArrayStr []string
	ArrayInt []int

	// time
	CreateAt dateTime.DateTime
	UpdateAt time.Time
	LastAt   time.Time
	FirstAt  dateTime.DateTime
	Ts       time.Duration
}

// 基础类型测试
func TestBasicType(t *testing.T) {
	float66_88, _ := decimal.NewFromFloat64(66.88)
	po := SamplePO1{
		UserName:     "UserName",
		IsEnable:     true,
		IsEnableStr:  "true",
		IsEnableBool: true,
		Id:           1,
		Id8:          8,
		Id16:         16,
		Id32:         32,
		Id64:         64,
		IdUint:       1,
		IdUint8:      8,
		IdUint16:     16,
		IdUint32:     32,
		IdUint64:     64,
		IdFloat32:    32.32,
		IdFloat64:    64.64,

		Id8Str:  "8",
		Id16Str: "16",

		Dec:      float66_88,
		ArrayStr: []string{"a", "b"},
		ArrayInt: []int{3, 4},
		CreateAt: time.Now(),
		UpdateAt: dateTime.Now(),
		LastAt:   time.Now(),
		FirstAt:  dateTime.Now(),
		Ts:       888,
	}

	do := mapper.Single[SamplePO2](po)
	assert.Equal(t, po.UserName, do.UserName)
	assert.Equal(t, po.IsEnable, do.IsEnable)
	assert.Equal(t, true, do.IsEnableStr)
	assert.Equal(t, "true", do.IsEnableBool)
	assert.Equal(t, po.Id, do.Id)
	assert.Equal(t, po.Id8, do.Id8)
	assert.Equal(t, po.Id16, do.Id16)
	assert.Equal(t, po.Id32, do.Id32)
	assert.Equal(t, po.Id64, do.Id64)
	assert.Equal(t, po.IdUint, do.IdUint)
	assert.Equal(t, po.IdUint8, do.IdUint8)
	assert.Equal(t, po.IdUint16, do.IdUint16)
	assert.Equal(t, po.IdUint32, do.IdUint32)
	assert.Equal(t, po.IdUint64, do.IdUint64)
	assert.Equal(t, po.IdFloat32, do.IdFloat32)
	assert.Equal(t, po.IdFloat64, do.IdFloat64)
	assert.Equal(t, int8(8), do.Id8Str)
	assert.Equal(t, int16(16), do.Id16Str)
	assert.Equal(t, po.Dec.String(), do.Dec.String())
	assert.Equal(t, po.ArrayStr[0], do.ArrayStr[0])
	assert.Equal(t, po.ArrayStr[1], do.ArrayStr[1])
	assert.Equal(t, po.ArrayInt[0], do.ArrayInt[0])
	assert.Equal(t, po.ArrayInt[1], do.ArrayInt[1])

	assert.Equal(t, po.CreateAt.String(), do.CreateAt.ToTime().String())
	assert.Equal(t, po.UpdateAt.ToTime().String(), do.UpdateAt.String())
	assert.Equal(t, po.LastAt.String(), do.LastAt.String())
	assert.Equal(t, po.FirstAt.ToTime().String(), do.FirstAt.ToTime().String())
}

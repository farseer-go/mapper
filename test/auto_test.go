package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/mapper"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestDtoToDoList(t *testing.T) {
	dto := TaskDTO{
		Array:   []UserVO{{List: collections.NewList[CountVO](CountVO{Count: 123})}},
		List:    collections.NewList[CountVO](CountVO{Count: 123}),
		List2:   collections.NewList[CountVO](CountVO{Count: 464}),
		LstType: ListType(collections.NewList[CountVO](CountVO{Count: 678})),
	}
	var do TaskDO
	_ = mapper.Auto(dto, &do)
	assert.Equal(t, dto.Array[0].List.Index(0).Count, do.Array[0].List.Index(0).Count)
	assert.Equal(t, dto.LstType.Index(0).Count, do.LstType.Index(0).Count)
	assert.Equal(t, dto.List.Index(0).Count, do.List.Index(0).Count)
	assert.Equal(t, dto.List2.Index(0).Count, do.List2.Index(0).Count)
}

func TestDtoToDoTimeDateTime(t *testing.T) {
	dto := TaskDTO{
		TimeInfo2:    time.Now(),
		TimeInfo:     dateTime.Now(),
		Time:         time.Now(),
		UpdateAt:     time.Now(),
		LastUpdateAt: dateTime.Now(),
		CreateAt:     time.Now(),
	}
	dto.UserVO3.Time = time.Now()
	dto.UserVO3.Time2 = dateTime.Now()
	dto.UserVO3.Time3 = time.Now()
	dto.UserVO3.Date = dateTime.Now()
	dto.UserVO3.Ts = time.Duration(90)
	var do TaskDO
	_ = mapper.Auto(dto, &do)
	assert.Equal(t, dto.TimeInfo2.Format("2006-01-02 15:04:05"), do.TimeInfo2.ToString("2006-01-02 15:04:05"))
	assert.Equal(t, dto.TimeInfo.ToString("2006-01-02 15:04:05"), do.TimeInfo.Format("2006-01-02 15:04:05"))
	assert.Equal(t, dto.Time.Format("2006-01-02 15:04:05"), do.Time)
	assert.Equal(t, dto.UpdateAt.Format("2006-01-02 15:04:05"), do.UpdateAt.ToString("yyyy-MM-dd HH:mm:ss"))
	assert.Equal(t, dto.LastUpdateAt.ToString("yyyy-MM-dd HH:mm:ss"), do.LastUpdateAt.Format("2006-01-02 15:04:05"))

	assert.Equal(t, dto.UserVO3.Time.Format("2006-01-02 15:04:05"), do.UserVO3.Time.Format("2006-01-02 15:04:05"))
	assert.Equal(t, dto.UserVO3.Time2.ToString("yyyy-MM-dd HH:mm:ss"), do.UserVO3.Time2.Format("2006-01-02 15:04:05"))
	assert.Equal(t, dto.UserVO3.Time3.Format("2006-01-02 15:04:05"), do.UserVO3.Time3.ToString("yyyy-MM-dd HH:mm:ss"))
	assert.Equal(t, dto.UserVO3.Date.ToString("yyyy-MM-dd HH:mm:ss"), do.UserVO3.Date.ToString("yyyy-MM-dd HH:mm:ss"))
	assert.Equal(t, dto.UserVO3.Ts, do.UserVO3.Ts)
	assert.Equal(t, dto.CreateAt.String(), do.CreateAt.String())

}

func TestDtoToDoDic(t *testing.T) {
	dto := TaskDTO{
		Data: collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
	}
	var do TaskDO
	_ = mapper.Auto(dto, &do)
	assert.Equal(t, dto.Data.Count(), do.Data.Count())
	assert.Equal(t, dto.Data.GetValue("age"), do.Data.GetValue("age"))
	assert.Equal(t, dto.Data.GetValue("price"), do.Data.GetValue("price"))
}

func TestDtoToDoMap(t *testing.T) {
	dto := TaskDTO{
		ClusterVer: map[int64]*SubMapEntity{
			2: {
				Age:     33,
				Caption: "测试map[64]*",
			},
		},
	}
	var do TaskDO
	_ = mapper.Auto(dto, &do)
	assert.Equal(t, dto.ClusterVer[2].Age, do.ClusterVer[2].Age)
	assert.Equal(t, dto.ClusterVer[2].Caption, do.ClusterVer[2].Caption)
}

func TestDtoToDoStruct(t *testing.T) {
	dto := TaskDTO{
		User: UserVO{
			Id:   1,
			Name: "steden",
		},
		ClientId:   1000,
		ClientIp:   "127.0.0.1",
		ClientName: "node",
		Status:     Pending,
	}
	dto.UserVO3.Name = "USER03"
	dto.UserVO3.Id = 123123
	dto.UserVO3.Stat = Pending
	var do TaskDO
	_ = mapper.Auto(dto, &do)
	assert.Equal(t, dto.User.Id, do.UserId)
	assert.Equal(t, dto.User.Name, do.UserName)
	assert.Equal(t, dto.ClientId, do.Client.Id)
	assert.Equal(t, dto.ClientIp, do.Client.Ip)
	assert.Equal(t, dto.ClientName, do.Client.Name)
	assert.Equal(t, dto.Status, do.Status)
	assert.Equal(t, dto.UserVO3.Name, do.UserVO3.Name)
	assert.Equal(t, dto.UserVO3.Id, do.UserVO3.Id)
	assert.Equal(t, dto.UserVO3.Stat, do.UserVO3.Stat)
}

func TestDtoToDoBaseType(t *testing.T) {
	dto := TaskDTO{
		IsEnable:     true,
		IsEnableStr:  "true",
		IsEnableBool: true,
		Id8:          8,
		Id8Str:       "8",
		Id16:         16,
		Id16Str:      "16",
		Id32:         32,
		Id64:         64,
		IdUint:       11,
		IdUint8:      88,
		IdUint16:     1616,
		IdUint32:     3232,
		IdUint64:     6464,
		IdFloat32:    32.32,
		IdFloat64:    64.64,
		Dec:          decimal.NewFromFloat32(12.22),
	}
	var do TaskDO
	_ = mapper.Auto(dto, &do)
	assert.Equal(t, dto.IsEnable, do.IsEnable)
	assert.Equal(t, dto.IsEnableStr, strconv.FormatBool(do.IsEnableStr))
	assert.Equal(t, strconv.FormatBool(dto.IsEnableBool), do.IsEnableBool)
	assert.Equal(t, dto.Id8, do.Id8)
	pi8, _ := strconv.ParseInt(dto.Id8Str, 0, 8)
	assert.Equal(t, int8(pi8), do.Id8Str)
	assert.Equal(t, dto.Id16, do.Id16)
	pi16, _ := strconv.ParseInt(dto.Id16Str, 0, 16)
	assert.Equal(t, int16(pi16), do.Id16Str)
	assert.Equal(t, dto.Id32, do.Id32)
	assert.Equal(t, dto.Id64, do.Id64)
	assert.Equal(t, dto.IdUint, do.IdUint)
	assert.Equal(t, dto.IdUint8, do.IdUint8)
	assert.Equal(t, dto.IdUint16, do.IdUint16)
	assert.Equal(t, dto.IdUint32, do.IdUint32)
	assert.Equal(t, dto.IdUint64, do.IdUint64)
	assert.Equal(t, dto.IdFloat32, do.IdFloat32)
	assert.Equal(t, dto.IdFloat64, do.IdFloat64)
	assert.Equal(t, dto.Dec, do.Dec)
}

func TestDoToDto(t *testing.T) {
	do := TaskDO{
		Id: 1,
		Client: ClientVO{
			Id:   2,
			Ip:   "192.168.1.1",
			Name: "node2",
		},
		Status:       Stopped,
		UserId:       666,
		UserName:     "harlen",
		Data:         collections.NewDictionaryFromMap(map[string]string{"age": "16", "price": "888.88"}),
		CreateAt:     time.Now(),
		IsEnable:     true,
		IsEnableStr:  true,
		IsEnableBool: "true",
		Id8:          8,
		Id8Str:       8,
		Id16:         16,
		Id16Str:      16,
		Id32:         32,
		Id64:         64,
		IdUint:       11,
		IdUint8:      88,
		IdUint16:     1616,
		IdUint32:     3232,
		IdUint64:     6464,
		IdFloat32:    32.32,
		IdFloat64:    64.64,
		UpdateAt:     dateTime.Now(),
		LastUpdateAt: time.Now(),
	}
	var dto TaskDTO
	_ = mapper.Auto(do, &dto)

	assert.Equal(t, do.Id, dto.Id)
	assert.Equal(t, do.Client.Id, dto.ClientId)
	assert.Equal(t, do.Client.Ip, dto.ClientIp)
	assert.Equal(t, do.Client.Name, dto.ClientName)
	assert.Equal(t, do.Status, dto.Status)
	assert.Equal(t, do.UserId, dto.User.Id)
	assert.Equal(t, do.UserName, dto.User.Name)
	assert.Equal(t, do.Data.Count(), dto.Data.Count())
	assert.Equal(t, do.Data.GetValue("age"), dto.Data.GetValue("age"))
	assert.Equal(t, do.Data.GetValue("price"), dto.Data.GetValue("price"))
	assert.Equal(t, do.CreateAt.String(), dto.CreateAt.String())
	assert.Equal(t, do.IsEnable, dto.IsEnable)
	assert.Equal(t, strconv.FormatBool(do.IsEnableStr), dto.IsEnableStr)
	assert.Equal(t, do.IsEnableBool, strconv.FormatBool(dto.IsEnableBool))
	assert.Equal(t, do.Id8, dto.Id8)
	pi8, _ := strconv.ParseInt(dto.Id8Str, 0, 8)
	assert.Equal(t, do.Id8Str, int8(pi8))
	assert.Equal(t, do.Id16, dto.Id16)
	pi16, _ := strconv.ParseInt(dto.Id16Str, 0, 16)
	assert.Equal(t, do.Id16Str, int16(pi16))
	assert.Equal(t, do.Id32, dto.Id32)
	assert.Equal(t, do.Id64, dto.Id64)
	assert.Equal(t, do.IdUint, dto.IdUint)
	assert.Equal(t, do.IdUint8, dto.IdUint8)
	assert.Equal(t, do.IdUint16, dto.IdUint16)
	assert.Equal(t, do.IdUint32, dto.IdUint32)
	assert.Equal(t, do.IdUint64, dto.IdUint64)
	assert.Equal(t, do.IdFloat32, dto.IdFloat32)
	assert.Equal(t, do.IdFloat64, dto.IdFloat64)
	assert.Equal(t, do.UpdateAt.ToString("yyyy-MM-dd HH:mm:ss"), dto.UpdateAt.Format("2006-01-02 15:04:05"))
	assert.Equal(t, do.LastUpdateAt.Format("2006-01-02 15:04:05"), dto.LastUpdateAt.ToString("yyyy-MM-dd HH:mm:ss"))
}

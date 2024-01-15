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

func TestDtoToDo(t *testing.T) {
	mapArray := map[int]*CountVO{0: {Count: 999}}
	mapArray2 := map[int]CountVO{0: {Count: 888}}
	mapArray3 := map[int]CountVO2{0: {Count: 777}}

	arrayStr := []string{0: "数组字符串测试"}
	lst3 := ListType(collections.NewList[CountVO](CountVO{Count: 678}))
	lst := collections.NewList[CountVO](CountVO{Count: 123})
	lst2 := collections.NewList[CountVO](CountVO{Count: 464})
	arrayUser := []UserVO{{List: lst, Id: 33, Name: "san", Array3: arrayStr, User3: UserVO3{Id: 55, Name: "user3"}, Count: mapArray, Count2: mapArray2, Count3: mapArray3}}

	dto := TaskDTO{
		ClusterVer: map[int64]*SubMapEntity{
			2: {
				Age:     33,
				Caption: "测试map[64]*",
			},
		},
		Data:       collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
		TimeInfo2:  time.Now(),
		TimeInfo:   dateTime.Now(),
		Time:       time.Now(),
		Dec:        decimal.NewFromFloat32(12.22),
		LstType:    lst3,
		ClientId:   1000,
		ClientIp:   "127.0.0.1",
		ClientName: "node",
		List:       lst,
		List2:      lst2,
		Array:      arrayUser,
		ArrayStr:   arrayStr,
		Id:         1,
		Status:     Pending,
		User: UserVO{
			Id:   1,
			Name: "steden",
		},

		CreateAt:     time.Now(),
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
		UpdateAt:     time.Now(),
		LastUpdateAt: dateTime.Now(),
	}
	dto.UserVO3.Name = "USER03"
	dto.UserVO3.Id = 123123
	dto.UserVO3.Time = time.Now()
	dto.UserVO3.Time2 = dateTime.Now()
	dto.UserVO3.Time3 = time.Now()
	dto.UserVO3.Date = dateTime.Now()
	dto.UserVO3.Stat = Pending
	dto.UserVO3.Dec = decimal.NewFromFloat32(12.22)
	dto.UserVO3.Ts = time.Duration(90)

	var do TaskDO
	_ = mapper.Auto(dto, &do)

	assert.Equal(t, dto.Array[0].User3.Id, do.Array[0].User3.Id)
	assert.Equal(t, dto.Array[0].User3.Name, do.Array[0].User3.Name)
	assert.Equal(t, dto.Array[0].Id, do.Array[0].Id)
	assert.Equal(t, dto.Array[0].Name, do.Array[0].Name)
	assert.Equal(t, dto.Array[0].Array3[0], do.Array[0].Array3[0])
	assert.Equal(t, dto.Array[0].Count[0].Count, do.Array[0].Count[0].Count)
	assert.Equal(t, dto.Array[0].Count2[0].Count, do.Array[0].Count2[0].Count)
	assert.Equal(t, dto.Array[0].Count3[0].Count, do.Array[0].Count3[0].Count)
	assert.Equal(t, dto.Array[0].List.Index(0).Count, do.Array[0].List.Index(0).Count)

	assert.Equal(t, dto.TimeInfo2.Format("2006-01-02 15:04:05"), do.TimeInfo2.ToString("2006-01-02 15:04:05"))
	assert.Equal(t, dto.TimeInfo.ToString("2006-01-02 15:04:05"), do.TimeInfo.Format("2006-01-02 15:04:05"))
	assert.Equal(t, dto.Time.Format("2006-01-02 15:04:05"), do.Time)
	assert.Equal(t, dto.Dec, do.Dec)

	assert.Equal(t, dto.LstType.Index(0).Count, do.LstType.Index(0).Count)
	assert.Equal(t, dto.List.Index(0).Count, do.List.Index(0).Count)
	assert.Equal(t, dto.List2.Index(0).Count, do.List2.Index(0).Count)
	assert.Equal(t, dto.ClientId+1, do.Client.Id)
	assert.Equal(t, dto.ClientIp, do.Client.Ip)
	assert.Equal(t, dto.ClientName, do.Client.Name)
	assert.Equal(t, dto.Status, do.Status)
	assert.Equal(t, dto.User.Id, do.UserId)
	assert.Equal(t, dto.User.Name, do.UserName)
	assert.Equal(t, dto.Data.Count(), do.Data.Count())
	assert.Equal(t, dto.Data.GetValue("age"), do.Data.GetValue("age"))
	assert.Equal(t, dto.Data.GetValue("price"), do.Data.GetValue("price"))
	assert.Equal(t, dto.CreateAt.String(), do.CreateAt.String())
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
	assert.Equal(t, dto.UpdateAt.Format("2006-01-02 15:04:05"), do.UpdateAt.ToString("yyyy-MM-dd HH:mm:ss"))
	assert.Equal(t, dto.LastUpdateAt.ToString("yyyy-MM-dd HH:mm:ss"), do.LastUpdateAt.Format("2006-01-02 15:04:05"))
	assert.Equal(t, dto.ClusterVer[2].Age, do.ClusterVer[2].Age)
	assert.Equal(t, dto.ClusterVer[2].Caption, do.ClusterVer[2].Caption)
}

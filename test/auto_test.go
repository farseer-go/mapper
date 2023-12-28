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

type State int

const (
	Running State = iota
	Pending
	Stopped
)

func (s State) String() string {
	switch s {
	case Running:
		return "Running"
	case Pending:
		return "Pending"
	case Stopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}

type ClientVO struct {
	Id   int64
	Ip   string
	Name string
}

func (vo *ClientVO) MapperInit() {
	vo.Id = vo.Id + 1
	println("已执行 ClientVO 初始化方法 MapperInit ")
}

type ListType collections.List[CountVO]

type UserVO struct {
	//List2  collections.List[CountVO2]
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
type UserVO3 struct {
	Id    int64
	Name  string
	Time2 dateTime.DateTime
	Time3 time.Time
	Time  time.Time
	Date  dateTime.DateTime
	Dec   decimal.Decimal
	Stat  State
	Ts    time.Duration
}
type UserVO4 struct {
	Id    int64
	Name  string
	Time2 time.Time
	Time3 dateTime.DateTime
	Time  time.Time
	Date  dateTime.DateTime
	Dec   decimal.Decimal
	Stat  State
	Ts    time.Duration
}

type CountVO struct {
	Count int // 出现的次数
}
type CountVO2 struct {
	Count int // 出现的次数
}
type TaskDO struct {
	UserVO3      UserVO4
	TimeInfo2    dateTime.DateTime
	TimeInfo     time.Time
	Time         string
	Dec          decimal.Decimal
	LstType      ListType
	Client       ClientVO
	List         collections.List[CountVO2]
	List2        collections.List[CountVO]
	Array        []UserVO2
	ArrayStr     []string
	Id           int
	Status       State
	UserId       int64
	UserName     string
	Data         collections.Dictionary[string, string]
	CreateAt     time.Time
	IsEnable     bool
	IsEnableStr  bool
	IsEnableBool string
	Id8          int8
	Id8Str       int8
	Id16         int16
	Id16Str      int16
	Id32         int32
	Id64         int64
	IdUint       uint
	IdUint8      uint8
	IdUint16     uint16
	IdUint32     uint32
	IdUint64     uint64
	IdFloat32    float32
	IdFloat64    float64
	UpdateAt     dateTime.DateTime
	LastUpdateAt time.Time
}

type TaskDTO struct {
	UserVO3      UserVO3
	TimeInfo2    time.Time
	TimeInfo     dateTime.DateTime
	Time         time.Time
	Dec          decimal.Decimal
	LstType      ListType
	List         collections.List[CountVO]
	List2        collections.List[CountVO]
	Array        []UserVO
	ArrayStr     []string
	Id           int
	ClientId     int64
	ClientIp     string
	ClientName   string
	Status       State
	User         UserVO
	Data         collections.Dictionary[string, string]
	CreateAt     time.Time
	UpdateAt     time.Time
	IsEnable     bool
	IsEnableStr  string
	IsEnableBool bool
	Id8          int8
	Id8Str       string
	Id16         int16
	Id16Str      string
	Id32         int32
	Id64         int64
	IdUint       uint
	IdUint8      uint8
	IdUint16     uint16
	IdUint32     uint32
	IdUint64     uint64
	IdFloat32    float32
	IdFloat64    float64
	LastUpdateAt dateTime.DateTime
	product      IProduct
}
type IProduct interface {
}

func (do *TaskDO) MapperInit() {
	do.Id = do.Id + 1
	println("已执行 TaskDO 初始化方法 MapperInit ")
}

func TestDtoToDo(t *testing.T) {
	mapArray := make(map[int]*CountVO, 1)
	mapArray[0] = &CountVO{Count: 999}
	mapArray2 := make(map[int]CountVO)
	mapArray2[0] = CountVO{Count: 888}
	mapArray3 := make(map[int]CountVO2)
	mapArray3[0] = CountVO2{Count: 777}
	arrayUser := make([]UserVO, 1)
	arrayStr := make([]string, 1)
	arrayStr[0] = "数组字符串测试"
	lst3 := ListType(collections.NewList[CountVO](CountVO{Count: 678}))
	lst := collections.NewList[CountVO](CountVO{Count: 123})
	lst2 := collections.NewList[CountVO](CountVO{Count: 464})
	arrayUser[0] = UserVO{List: lst, Id: 33, Name: "san", Array3: arrayStr, User3: UserVO3{Id: 55, Name: "user3"}, Count: mapArray, Count2: mapArray2, Count3: mapArray3}
	dto := TaskDTO{
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
		Data:         collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
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

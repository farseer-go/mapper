package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/shopspring/decimal"
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

type TaskDO struct {
	ClusterVer   map[int64]*SubMapEntity
	Data         collections.Dictionary[string, string]
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

func (do *TaskDO) MapperInit() {
	do.Id = do.Id + 1
	println("已执行 TaskDO 初始化方法 MapperInit ")
}

type IProduct interface {
}

type CountVO struct {
	Count int // 出现的次数
}

type CountVO2 struct {
	Count int // 出现的次数
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

type TaskDTO struct {
	ClusterVer   map[int64]*SubMapEntity
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

package mapper

import (
	"github.com/farseer-go/collections"
	"github.com/stretchr/testify/assert"
	"testing"
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

type UserVO struct {
	Id   int64
	Name string
}

type TaskDO struct {
	Id       int
	Client   ClientVO
	Status   State
	UserId   int64
	UserName string
	Data     collections.Dictionary[string, string]
}

type TaskDTO struct {
	Id         int
	ClientId   int64
	ClientIp   string
	ClientName string
	Status     State
	User       UserVO
	Data       collections.Dictionary[string, string]
}

func TestAutoMapper(t *testing.T) {
	t.Run("dto转do", func(t *testing.T) {
		dto := TaskDTO{
			Id:         1,
			ClientId:   1000,
			ClientIp:   "127.0.0.1",
			ClientName: "node",
			Status:     Pending,
			User: UserVO{
				Id:   88,
				Name: "steden",
			},
			Data: collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
		}

		var do TaskDO
		_ = MapDOtoDTO(dto, &do)
		assert.Equal(t, dto.Id, do.Id)
		assert.Equal(t, dto.ClientId, do.Client.Id)
		assert.Equal(t, dto.ClientIp, do.Client.Ip)
		assert.Equal(t, dto.ClientName, do.Client.Name)
		assert.Equal(t, dto.Status, do.Status)
		assert.Equal(t, dto.User.Id, do.UserId)
		assert.Equal(t, dto.User.Name, do.UserName)
		assert.Equal(t, dto.Data.Count(), do.Data.Count())
		assert.Equal(t, dto.Data.GetValue("age"), do.Data.GetValue("age"))
		assert.Equal(t, dto.Data.GetValue("price"), do.Data.GetValue("price"))
	})

	//t.Run("do转dto", func(t *testing.T) {
	//	do := TaskDO{
	//		Id: 1,
	//		Client: ClientVO{
	//			Id:   2,
	//			Ip:   "192.168.1.1",
	//			Name: "node2",
	//		},
	//		Status:   Stopped,
	//		UserId:   666,
	//		UserName: "harlen",
	//		Data:     collections.NewDictionaryFromMap(map[string]string{"age": "16", "price": "888.88"}),
	//	}
	//	var dto TaskDTO
	//	_ = MapDOtoDTO(do, &dto)
	//
	//	assert.Equal(t, do.Id, dto.Id)
	//	assert.Equal(t, do.Client.Id, dto.ClientId)
	//	assert.Equal(t, do.Client.Ip, dto.ClientIp)
	//	assert.Equal(t, do.Client.Name, dto.ClientName)
	//	assert.Equal(t, do.Status, dto.Status)
	//	assert.Equal(t, do.UserId, dto.User.Id)
	//	assert.Equal(t, do.UserName, dto.User.Name)
	//	assert.Equal(t, do.Data.Count(), dto.Data.Count())
	//	assert.Equal(t, do.Data.GetValue("age"), dto.Data.GetValue("age"))
	//	assert.Equal(t, do.Data.GetValue("price"), dto.Data.GetValue("price"))
	//})
}

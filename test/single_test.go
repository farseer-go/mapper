package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSingle(t *testing.T) {
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

	doSingle := mapper.Single[TaskDO](dto)
	doSingle2 := mapper.Single[TaskDTO](doSingle)
	assert.Equal(t, dto.Id+1, doSingle.Id)
	assert.Equal(t, dto.ClientId+1, doSingle.Client.Id)
	assert.Equal(t, doSingle2.Id, doSingle.Id)
	assert.Equal(t, doSingle2.ClientId, doSingle.Client.Id)
}

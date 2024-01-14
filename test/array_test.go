package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArray(t *testing.T) {
	arrDto := []TaskDTO{{
		Id:         1,
		ClientId:   1000,
		ClientIp:   "127.0.0.1",
		ClientName: "node",
		Status:     Pending,
		User:       UserVO{Id: 88, Name: "steden"},
		Data:       collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
	}, {
		Id:         2,
		ClientId:   1000,
		ClientIp:   "127.0.0.1",
		ClientName: "node",
		Status:     Pending,
		User:       UserVO{Id: 20, Name: "steden1"},
		Data:       collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
	}}

	arrDO := mapper.Array[TaskDO](arrDto)
	assert.Equal(t, len(arrDto), len(arrDO))

	for i := 0; i < len(arrDto); i++ {
		assert.Equal(t, arrDto[i].User.Name, arrDO[i].UserName)
		assert.Equal(t, arrDto[i].User.Id, arrDO[i].UserId)
	}
}

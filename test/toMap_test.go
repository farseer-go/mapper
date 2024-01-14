package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToMap(t *testing.T) {
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
	dic := mapper.ToMap[string, any](dto)
	assert.Equal(t, "steden", dic["User"].(UserVO).Name)
	assert.Equal(t, int64(88), dic["User"].(UserVO).Id)
}

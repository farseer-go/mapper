package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPageList(t *testing.T) {
	arrDto := collections.NewPageList[TaskDTO](collections.NewList(TaskDTO{
		Id:         1,
		ClientId:   1000,
		ClientIp:   "127.0.0.1",
		ClientName: "node",
		Status:     Pending,
		User:       UserVO{Id: 88, Name: "steden"},
		Data:       collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
	}, TaskDTO{
		Id:         2,
		ClientId:   1000,
		ClientIp:   "127.0.0.1",
		ClientName: "node",
		Status:     Pending,
		User:       UserVO{Id: 20, Name: "steden1"},
		Data:       collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
	}), 10)

	lst := mapper.ToPageList[TaskDO](arrDto)

	assert.Equal(t, arrDto.List.Count(), lst.List.Count())

	assert.Equal(t, lst.RecordCount, int64(10))
	for i := 0; i < arrDto.List.Count(); i++ {
		assert.Equal(t, arrDto.List.Index(i).User.Name, lst.List.Index(i).UserName)
		assert.Equal(t, arrDto.List.Index(i).User.Id, lst.List.Index(i).UserId)
	}
}

package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToListAny(t *testing.T) {
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

	listAny := mapper.ToListAny(arrDto)

	assert.Equal(t, listAny.Count(), len(arrDto))
	for i := 0; i < listAny.Count(); i++ {
		assert.Equal(t, listAny.Index(i).(TaskDTO).User.Name, arrDto[i].User.Name)
		assert.Equal(t, listAny.Index(i).(TaskDTO).User.Id, arrDto[i].User.Id)
	}

	lst := collections.NewList(arrDto...)
	listAny = mapper.ToListAny(lst)

	assert.Equal(t, listAny.Count(), len(arrDto))
	for i := 0; i < listAny.Count(); i++ {
		assert.Equal(t, listAny.Index(i).(TaskDTO).User.Name, arrDto[i].User.Name)
		assert.Equal(t, listAny.Index(i).(TaskDTO).User.Id, arrDto[i].User.Id)
	}
}

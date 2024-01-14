package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type po struct {
	Name string
	Age  int
}
type do struct {
	Name string
	Age  int
}

func TestToList(t *testing.T) {
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
	lst := collections.NewList(arrDto...)
	lstDO := mapper.ToList[TaskDO](lst)

	assert.Equal(t, lst.Count(), lstDO.Count())

	for i := 0; i < lst.Count(); i++ {
		assert.Equal(t, lst.Index(i).User.Name, lstDO.Index(i).UserName)
		assert.Equal(t, lst.Index(i).User.Id, lstDO.Index(i).UserId)
	}

	lstAny := lst.ToListAny()
	lstDO = mapper.ToList[TaskDO](lstAny)

	assert.Equal(t, lstAny.Count(), lstDO.Count())

	for i := 0; i < lstAny.Count(); i++ {
		assert.Equal(t, lstAny.Index(i).(TaskDTO).User.Name, lstDO.Index(i).UserName)
		assert.Equal(t, lstAny.Index(i).(TaskDTO).User.Id, lstDO.Index(i).UserId)
	}

	arr := lst.ToArray()
	lstDO = mapper.ToList[TaskDO](arr)

	assert.Equal(t, len(arr), lstDO.Count())

	for i := 0; i < lstAny.Count(); i++ {
		assert.Equal(t, arr[i].User.Name, lstDO.Index(i).UserName)
		assert.Equal(t, arr[i].User.Id, lstDO.Index(i).UserId)
	}
}

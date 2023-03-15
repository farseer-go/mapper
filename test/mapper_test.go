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
	assert.Equal(t, dto.Id, doSingle.Id)
	assert.Equal(t, dto.ClientId, doSingle.Client.Id)
	assert.Equal(t, doSingle2.Id, doSingle.Id)
	assert.Equal(t, doSingle2.ClientId, doSingle.Client.Id)
}

func TestPageList(t *testing.T) {
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

	lst := mapper.ToPageList[TaskDO](arrDto, 10)

	assert.Equal(t, len(arrDto), lst.List.Count())

	assert.Equal(t, lst.RecordCount, int64(10))
	for i := 0; i < len(arrDto); i++ {
		assert.Equal(t, arrDto[i].User.Name, lst.List.Index(i).UserName)
		assert.Equal(t, arrDto[i].User.Id, lst.List.Index(i).UserId)
	}
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

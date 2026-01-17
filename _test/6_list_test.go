package test

import (
	"testing"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
)

func TestToList(t *testing.T) {
	type userVO struct {
		Id   int64
		Name string
	}

	type s1 struct {
		Id   int
		User userVO
	}

	type s2 struct {
		Id       int
		UserId   int64
		UserName string
	}

	lst := collections.NewList([]s1{{
		Id:   1,
		User: userVO{Id: 88, Name: "steden"},
	}, {
		Id:   2,
		User: userVO{Id: 20, Name: "steden1"},
	}}...)

	lstDO := mapper.ToList[s2](lst)

	assert.Equal(t, lst.Count(), lstDO.Count())
	for i := 0; i < lst.Count(); i++ {
		assert.Equal(t, lst.Index(i).Id, lstDO.Index(i).Id)
		assert.Equal(t, lst.Index(i).User.Name, lstDO.Index(i).UserName)
		assert.Equal(t, lst.Index(i).User.Id, lstDO.Index(i).UserId)
	}

	lstAny := lst.ToListAny()
	lstDO = mapper.ToList[s2](lstAny)

	assert.Equal(t, lstAny.Count(), lstDO.Count())

	for i := 0; i < lstAny.Count(); i++ {
		assert.Equal(t, lstAny.Index(i).(s1).Id, lstDO.Index(i).Id)
		assert.Equal(t, lstAny.Index(i).(s1).User.Name, lstDO.Index(i).UserName)
		assert.Equal(t, lstAny.Index(i).(s1).User.Id, lstDO.Index(i).UserId)
	}

	arr := lst.ToArray()
	lstDO = mapper.ToList[s2](arr)
	assert.Equal(t, len(arr), lstDO.Count())
	for i := 0; i < lstAny.Count(); i++ {
		assert.Equal(t, arr[i].Id, lstDO.Index(i).Id)
		assert.Equal(t, arr[i].User.Name, lstDO.Index(i).UserName)
		assert.Equal(t, arr[i].User.Id, lstDO.Index(i).UserId)
	}

	// 结体体数组转成any数组
	lstAny2 := mapper.ToList[any](arr)
	assert.Equal(t, len(arr), lstAny2.Count())
	for i := 0; i < lstAny.Count(); i++ {
		assert.Equal(t, arr[i].Id, lstAny2.Index(i).(s1).Id)
		assert.Equal(t, arr[i].User.Name, lstAny2.Index(i).(s1).User.Name)
		assert.Equal(t, arr[i].User.Id, lstAny2.Index(i).(s1).User.Id)
	}
}

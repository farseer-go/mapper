package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToListAny(t *testing.T) {
	type userVO struct {
		Id   int64
		Name string
	}

	type s1 struct {
		Id   int
		User userVO
	}

	arrDto := []s1{{
		Id:   1,
		User: userVO{Id: 88, Name: "steden"},
	}, {
		Id:   2,
		User: userVO{Id: 20, Name: "steden1"},
	}}

	listAny := mapper.ToListAny(arrDto)

	assert.Equal(t, listAny.Count(), len(arrDto))
	for i := 0; i < listAny.Count(); i++ {
		assert.Equal(t, listAny.Index(i).(s1).Id, arrDto[i].Id)
		assert.Equal(t, listAny.Index(i).(s1).User.Name, arrDto[i].User.Name)
		assert.Equal(t, listAny.Index(i).(s1).User.Id, arrDto[i].User.Id)
	}

	lst := collections.NewList(arrDto...)
	listAny = mapper.ToListAny(lst)

	assert.Equal(t, listAny.Count(), len(arrDto))
	for i := 0; i < listAny.Count(); i++ {
		assert.Equal(t, listAny.Index(i).(s1).Id, arrDto[i].Id)
		assert.Equal(t, listAny.Index(i).(s1).User.Name, arrDto[i].User.Name)
		assert.Equal(t, listAny.Index(i).(s1).User.Id, arrDto[i].User.Id)
	}
}

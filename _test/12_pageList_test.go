package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPageList(t *testing.T) {
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

	arrDto := collections.NewPageList[s1](collections.NewList(s1{
		Id:   1,
		User: userVO{Id: 88, Name: "steden"},
	}, s1{
		Id:   2,
		User: userVO{Id: 20, Name: "steden1"},
	}), 10)

	lst := mapper.ToPageList[s2](arrDto)

	assert.Equal(t, arrDto.List.Count(), lst.List.Count())
	assert.Equal(t, lst.RecordCount, int64(10))

	for i := 0; i < arrDto.List.Count(); i++ {
		assert.Equal(t, arrDto.List.Index(i).Id, lst.List.Index(i).Id)
		assert.Equal(t, arrDto.List.Index(i).User.Name, lst.List.Index(i).UserName)
		assert.Equal(t, arrDto.List.Index(i).User.Id, lst.List.Index(i).UserId)
	}
}

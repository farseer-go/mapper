package test

import (
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试匿名 + 名称匹配
func TestAnonymousNameMatch(t *testing.T) {
	type user1VO struct {
		Id   int64
		Name string
	}
	type User2VO struct {
		UserId   int64
		UserName string
	}

	type arr1 struct {
		Id   int
		User user1VO
	}

	type arr2 struct {
		Id int
		User2VO
		UserName string
	}

	arrDto := []arr1{{
		Id:   1,
		User: user1VO{Id: 88, Name: "steden"},
	}, {
		Id:   2,
		User: user1VO{Id: 20, Name: "steden1"},
	}}

	arrDO := mapper.Array[arr2](arrDto)
	assert.Equal(t, len(arrDto), len(arrDO))

	for i := 0; i < len(arrDto); i++ {
		assert.Equal(t, arrDto[i].User.Name, arrDO[i].UserName)
		assert.Equal(t, arrDto[i].User.Name, arrDO[i].User2VO.UserName)
		assert.Equal(t, arrDto[i].User.Id, arrDO[i].User2VO.UserId)
	}
}

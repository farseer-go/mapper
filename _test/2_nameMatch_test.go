package test

import (
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试名称匹配
func TestNameMatch(t *testing.T) {
	type userMatchVO struct {
		Name string
		Id   int
	}

	type NameMatch1 struct {
		UserName string
		UserId   int
		Admin    userMatchVO
	}

	type NameMatch2 struct {
		User      userMatchVO
		AdminName string
		AdminId   int
	}

	po := NameMatch1{
		UserName: "aaa",
		UserId:   222,
		Admin: userMatchVO{
			Name: "bbb",
			Id:   333,
		},
	}

	do := mapper.Single[NameMatch2](po)
	assert.Equal(t, po.UserName, do.User.Name)
	assert.Equal(t, po.UserId, do.User.Id)
	assert.Equal(t, po.Admin.Id, do.AdminId)
	assert.Equal(t, po.Admin.Name, do.AdminName)
}

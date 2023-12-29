package test

import (
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type BaseEntity struct {
	AppId   int64  // 应用ID
	AppName string // 应用名称
}
type MapEntity struct {
	BaseEntity
	UserId   int
	UserName string
}

func TestMapToEntity(t *testing.T) {
	m := make(map[string]any)
	m["AppId"] = int64(1)
	m["AppName"] = "test"
	m["UserId"] = 888
	m["UserName"] = "steden"
	entity := mapper.Single[MapEntity](m)
	assert.Equal(t, m["AppId"], entity.AppId)
	assert.Equal(t, m["AppName"], entity.AppName)
	assert.Equal(t, m["UserId"], entity.UserId)
	assert.Equal(t, m["UserName"], entity.UserName)
}

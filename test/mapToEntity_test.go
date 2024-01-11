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
type SubMapEntity struct {
	Age     int64  // 应用ID
	Caption string // 应用名称
}
type MapEntity struct {
	BaseEntity
	UserId   int
	UserName string
	Sub      SubMapEntity
}

func TestMapToEntity(t *testing.T) {
	m := make(map[string]any)
	m["AppId"] = int64(1)
	m["AppName"] = "test"
	m["UserId"] = 888
	m["UserName"] = "steden"
	m["Sub"] = map[string]any{
		"Age":     18,
		"Caption": "有值吗",
	}

	entity := mapper.Single[MapEntity](m)
	assert.Equal(t, m["AppId"], entity.AppId)
	assert.Equal(t, m["AppName"], entity.AppName)
	assert.Equal(t, m["UserId"], entity.UserId)
	assert.Equal(t, m["UserName"], entity.UserName)
	assert.Equal(t, m["Sub"].(map[string]any)["Age"], entity.Sub.Age)
	assert.Equal(t, m["Sub"].(map[string]any)["Caption"], entity.Sub.Caption)
}

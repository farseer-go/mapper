package test

import (
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type BaseEntity struct {
	Exception SubMapEntity
	AppId     int64  // 应用ID
	AppName   string // 应用名称
	CreateAt  dateTime.DateTime
}
type SubMapEntity struct {
	Age     int    // 应用ID
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
	m["Sub"] = map[string]any{
		"Age":     18,
		"Caption": "有值吗",
	}
	m["Exception"] = map[string]any{
		"Age":     22,
		"Caption": "嵌入字段",
	}
	m["AppId"] = int64(1)
	m["AppName"] = "test"
	m["UserId"] = 888
	m["UserName"] = "steden"
	m["CreateAt"] = dateTime.Now()

	entity := mapper.Single[MapEntity](m)
	assert.Equal(t, m["AppId"], entity.AppId)
	assert.Equal(t, m["AppName"], entity.AppName)
	assert.Equal(t, m["UserId"], entity.UserId)
	assert.Equal(t, m["UserName"], entity.UserName)
	assert.Equal(t, m["Sub"].(map[string]any)["Age"], entity.Sub.Age)
	assert.Equal(t, m["Sub"].(map[string]any)["Caption"], entity.Sub.Caption)
	assert.Equal(t, m["Exception"].(map[string]any)["Age"], entity.Exception.Age)
	assert.Equal(t, m["Exception"].(map[string]any)["Caption"], entity.Exception.Caption)
	assert.Equal(t, (m["CreateAt"].(dateTime.DateTime)).ToString("yyyyMMdd"), entity.CreateAt.ToString("yyyyMMdd"))
}

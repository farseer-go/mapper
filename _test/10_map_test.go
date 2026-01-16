package test

import (
	"testing"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	type subMapEntity struct {
		Age     int    // 应用ID
		Caption string // 应用名称
	}

	type s1 struct {
		ClusterVer map[int64]*subMapEntity
		A          map[string]any
		B          collections.Dictionary[string, string]
	}
	type s2 struct {
		ClusterVer map[int64]subMapEntity
		A          collections.Dictionary[string, string]
		B          map[string]any
	}

	dto := s1{
		ClusterVer: map[int64]*subMapEntity{
			2: {
				Age:     33,
				Caption: "测试map[64]*",
			},
		},
		A: map[string]any{"1": "2"},
		B: collections.NewDictionaryFromMap(map[string]string{"1": "2"}),
	}
	do := mapper.Single[s2](dto)
	assert.Equal(t, dto.ClusterVer[2].Age, do.ClusterVer[2].Age)
	assert.Equal(t, dto.ClusterVer[2].Caption, do.ClusterVer[2].Caption)
	assert.Equal(t, dto.A["1"], do.A.GetValue("1"))
	assert.Equal(t, dto.B.GetValue("1"), do.B["1"])
}

func TestStructToMap(t *testing.T) {
	type state int
	const (
		Running state = iota
		Pending
		Stopped
	)

	type userVO struct {
		Id   int64
		Name string
	}
	type s1 struct {
		Id     int
		Status state
		User   userVO
		aa     string
	}

	dto := s1{
		Id:     1,
		Status: Pending,
		User: userVO{
			Id:   88,
			Name: "steden",
		},
		aa: "fff",
	}
	dic := mapper.ToMap[string, any](dto)
	assert.Equal(t, "steden", dic["User"].(userVO).Name)
	assert.Equal(t, int64(88), dic["User"].(userVO).Id)
	assert.Equal(t, 1, dic["Id"].(int))
}

func TestMapToStruct(t *testing.T) {
	type subMapEntity struct {
		Age     int    // 应用ID
		Caption string // 应用名称
	}

	m := make(map[string]any)
	m["CreateAt"] = dateTime.Now()
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
	m["ClusterVer"] = map[int64]*subMapEntity{
		2: {
			Age:     33,
			Caption: "测试map[64]*",
		},
	}
	m["Head"] = map[string]any{
		"Content-Type": "application/json",
	}
	m["CallType"] = 1

	type eumCallTypeEnum int
	type BaseEntity struct {
		Exception subMapEntity
		AppId     int64  // 应用ID
		AppName   string // 应用名称
		CreateAt  dateTime.DateTime
	}
	type mapEntity struct {
		BaseEntity
		UserId     int
		UserName   string
		Sub        subMapEntity
		ClusterVer map[int64]*subMapEntity
		Head       collections.Dictionary[string, string]
		CallType   eumCallTypeEnum
	}

	entity := mapper.Single[mapEntity](m)
	assert.Equal(t, m["AppId"], entity.AppId)
	assert.Equal(t, m["AppName"], entity.AppName)
	assert.Equal(t, m["UserId"], entity.UserId)
	assert.Equal(t, m["UserName"], entity.UserName)
	assert.Equal(t, m["Sub"].(map[string]any)["Age"], entity.Sub.Age)
	assert.Equal(t, m["Sub"].(map[string]any)["Caption"], entity.Sub.Caption)
	assert.Equal(t, m["Exception"].(map[string]any)["Age"], entity.Exception.Age)
	assert.Equal(t, m["Exception"].(map[string]any)["Caption"], entity.Exception.Caption)
	assert.Equal(t, (m["CreateAt"].(dateTime.DateTime)).ToString("yyyyMMdd"), entity.CreateAt.ToString("yyyyMMdd"))

	assert.Equal(t, m["ClusterVer"].(map[int64]*subMapEntity)[2].Age, entity.ClusterVer[2].Age)
	assert.Equal(t, m["ClusterVer"].(map[int64]*subMapEntity)[2].Caption, entity.ClusterVer[2].Caption)
	assert.Equal(t, m["Head"].(map[string]any)["Content-Type"], entity.Head.GetValue("Content-Type"))
	assert.Equal(t, eumCallTypeEnum(parse.ToInt(m["CallType"])), entity.CallType)
}

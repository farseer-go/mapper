package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDic(t *testing.T) {
	type s1 struct {
		Data1 collections.Dictionary[string, string]
		Data2 collections.Dictionary[string, any]
		Data3 collections.Dictionary[string, string]
	}
	type s2 struct {
		Data1 collections.Dictionary[string, string]
		Data2 collections.Dictionary[string, any]
		Data3 collections.Dictionary[string, any]
	}

	dto := s1{
		Data1: collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
		Data2: collections.NewDictionaryFromMap(map[string]any{"age": 18, "price": 88.88}),
		Data3: collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
	}

	do := mapper.Single[s2](dto)
	assert.Equal(t, dto.Data1.Count(), do.Data1.Count())
	assert.Equal(t, dto.Data1.GetValue("age"), do.Data1.GetValue("age"))
	assert.Equal(t, dto.Data1.GetValue("price"), do.Data1.GetValue("price"))

	assert.Equal(t, dto.Data2.Count(), do.Data2.Count())
	assert.Equal(t, dto.Data2.GetValue("age"), do.Data2.GetValue("age"))
	assert.Equal(t, dto.Data2.GetValue("price"), do.Data2.GetValue("price"))

	assert.Equal(t, dto.Data3.Count(), do.Data3.Count())
	assert.Equal(t, dto.Data3.GetValue("age"), do.Data3.GetValue("age"))
	assert.Equal(t, dto.Data3.GetValue("price"), do.Data3.GetValue("price"))
}

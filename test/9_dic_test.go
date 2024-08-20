package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDic(t *testing.T) {
	type s1 struct {
		Data collections.Dictionary[string, string]
	}

	type s2 struct {
		Data collections.Dictionary[string, string]
	}

	dto := s1{
		Data: collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
	}

	do := mapper.Single[s2](dto)
	assert.Equal(t, dto.Data.Count(), do.Data.Count())
	assert.Equal(t, dto.Data.GetValue("age"), do.Data.GetValue("age"))
	assert.Equal(t, dto.Data.GetValue("price"), do.Data.GetValue("price"))
}

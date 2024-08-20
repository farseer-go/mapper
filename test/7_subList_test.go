package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubListToSubList(t *testing.T) {
	type countVO struct {
		Count int // 出现的次数
	}

	type countVO2 struct {
		Count int // 出现的次数
	}

	type userVO struct {
		List   collections.List[countVO]
		Id     int64
		Name   string
		Array3 []string
		Count  map[int]*countVO
		Count2 map[int]countVO
		Count3 map[int]countVO2
	}
	type userVO2 struct {
		List   collections.List[countVO]
		Id     int64
		Name   string
		Array3 []string
		Count  map[int]*countVO2
		Count2 map[int]countVO2
		Count3 map[int]countVO2
	}
	type a struct {
		Array []userVO
		List2 collections.List[countVO]
	}
	type b struct {
		Array []userVO2
		List2 collections.List[countVO]
	}

	dto := a{
		Array: []userVO{{
			List:   collections.NewList[countVO](countVO{Count: 123}),
			Id:     33,
			Name:   "abc",
			Array3: []string{"a1", "b2"},
			Count:  map[int]*countVO{11: {Count: 123}},
			Count2: map[int]countVO{22: {Count: 1234}},
			Count3: map[int]countVO2{33: {Count: 12345}},
		}},
		List2: collections.NewList[countVO](countVO{Count: 464}),
	}

	do := mapper.Single[b](dto)
	assert.Equal(t, dto.Array[0].List.Index(0).Count, do.Array[0].List.Index(0).Count)
	assert.Equal(t, dto.Array[0].Id, do.Array[0].Id)
	assert.Equal(t, dto.Array[0].Name, do.Array[0].Name)
	assert.Equal(t, dto.Array[0].Array3[0], do.Array[0].Array3[0])
	assert.Equal(t, dto.Array[0].Array3[1], do.Array[0].Array3[1])
	assert.Equal(t, dto.Array[0].Count[11].Count, do.Array[0].Count[11].Count)
	assert.Equal(t, dto.Array[0].Count2[22].Count, do.Array[0].Count2[22].Count)
	assert.Equal(t, dto.Array[0].Count3[33].Count, do.Array[0].Count3[33].Count)

	assert.Equal(t, dto.List2.Index(0).Count, do.List2.Index(0).Count)
}

package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 数组 转 List
func TestArray(t *testing.T) {
	// 字符串
	var array = []string{"A", "B", "C", "D"}
	lstStr := mapper.ToList[string](array)
	assert.Equal(t, lstStr.Index(0), array[0])
	assert.Equal(t, lstStr.Index(1), array[1])
	assert.Equal(t, lstStr.Index(2), array[2])
	assert.Equal(t, lstStr.Index(3), array[3])

	// int类型
	var arrayInt = []uint8{1, 2, 3, 4}
	lstInt := mapper.ToList[uint8](arrayInt)
	assert.Equal(t, lstInt.Index(0), arrayInt[0])
	assert.Equal(t, lstInt.Index(1), arrayInt[1])
	assert.Equal(t, lstInt.Index(2), arrayInt[2])
	assert.Equal(t, lstInt.Index(3), arrayInt[3])

	type countVO struct {
		Count int // 出现的次数
	}

	type countVO2 struct {
		Count int // 出现的次数
	}
	type userVO3 struct {
		Id   int64
		Name string
	}
	type userVO struct {
		List   collections.List[countVO]
		Id     int64
		Name   string
		User3  userVO3
		Array3 []string
		Count  map[int]*countVO
		Count2 map[int]countVO
		Count3 map[int]countVO2
	}

	type s1 struct {
		List  collections.List[countVO]
		List2 collections.List[countVO]
		Array []userVO
		User  userVO
	}

	dto := s1{
		List:  collections.NewList[countVO](countVO{Count: 123}),
		List2: collections.NewList[countVO](countVO{Count: 464}),
		Array: []userVO{{
			List:   collections.NewList[countVO](countVO{Count: 123}),
			Id:     33,
			Name:   "san",
			Array3: []string{0: "数组字符串测试"},
			User3: userVO3{
				Id:   55,
				Name: "user3",
			},
			Count:  map[int]*countVO{0: {Count: 999}},
			Count2: map[int]countVO{0: {Count: 888}},
			Count3: map[int]countVO2{0: {Count: 777}},
		}},
		User: userVO{
			Id:   1,
			Name: "steden",
		},
	}

	do := mapper.Single[s1](dto)

	assert.Equal(t, dto.List.Index(0).Count, do.List.Index(0).Count)
	assert.Equal(t, dto.List2.Index(0).Count, do.List2.Index(0).Count)
	assert.Equal(t, dto.Array[0].List.Index(0).Count, do.Array[0].List.Index(0).Count)
	assert.Equal(t, dto.Array[0].Id, do.Array[0].Id)
	assert.Equal(t, dto.Array[0].Name, do.Array[0].Name)
	assert.Equal(t, dto.Array[0].Array3[0], do.Array[0].Array3[0])
	assert.Equal(t, dto.Array[0].User3.Id, do.Array[0].User3.Id)
	assert.Equal(t, dto.Array[0].User3.Name, do.Array[0].User3.Name)
	assert.Equal(t, dto.Array[0].Count[0].Count, do.Array[0].Count[0].Count)
	assert.Equal(t, dto.Array[0].Count2[0].Count, do.Array[0].Count2[0].Count)
	assert.Equal(t, dto.Array[0].Count3[0].Count, do.Array[0].Count3[0].Count)
	assert.Equal(t, dto.User.Id, do.User.Id)
	assert.Equal(t, dto.User.Name, do.User.Name)
}

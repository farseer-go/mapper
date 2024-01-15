package test

import (
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type DtoToDoArrayDTO struct {
	Array []UserVO
}

type ArrayDTO struct {
	Array []UserVO2
}

func TestDtoToDoArray(t *testing.T) {
	mapArray := map[int]*CountVO{0: {Count: 999}}
	mapArray2 := map[int]CountVO{0: {Count: 888}}
	mapArray3 := map[int]CountVO2{0: {Count: 777}}
	arrayStr := []string{0: "数组字符串测试"}

	dto := DtoToDoArrayDTO{
		Array: []UserVO{{Id: 33, Name: "san", Array3: arrayStr, User3: UserVO3{Id: 55, Name: "user3"}, Count: mapArray, Count2: mapArray2, Count3: mapArray3}},
	}

	do := mapper.Single[ArrayDTO](dto)
	assert.Equal(t, dto.Array[0].User3.Id, do.Array[0].User3.Id)
	assert.Equal(t, dto.Array[0].User3.Name, do.Array[0].User3.Name)
	assert.Equal(t, dto.Array[0].Id, do.Array[0].Id)
	assert.Equal(t, dto.Array[0].Name, do.Array[0].Name)
	assert.Equal(t, dto.Array[0].Array3[0], do.Array[0].Array3[0])
	assert.Equal(t, dto.Array[0].Count[0].Count, do.Array[0].Count[0].Count)
	assert.Equal(t, dto.Array[0].Count2[0].Count, do.Array[0].Count2[0].Count)
	assert.Equal(t, dto.Array[0].Count3[0].Count, do.Array[0].Count3[0].Count)
}

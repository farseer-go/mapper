package test

import (
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrayToList(t *testing.T) {
	// 字符串
	var array = []string{"A", "B", "C", "D"}
	lst := mapper.ToList[string](array)
	assert.Equal(t, lst.Index(0), array[0])
	assert.Equal(t, lst.Index(1), array[1])
	assert.Equal(t, lst.Index(2), array[2])
	assert.Equal(t, lst.Index(3), array[3])

	// int类型
	var arrayInt = []uint8{1, 2, 3, 4}
	lstInt := mapper.ToList[uint8](arrayInt)
	assert.Equal(t, lstInt.Index(0), arrayInt[0])
	assert.Equal(t, lstInt.Index(1), arrayInt[1])
	assert.Equal(t, lstInt.Index(2), arrayInt[2])
	assert.Equal(t, lstInt.Index(3), arrayInt[3])
}

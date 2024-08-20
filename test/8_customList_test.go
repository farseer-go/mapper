package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试自定义List
func TestCustomList(t *testing.T) {
	type CountVO struct {
		Count int // 出现的次数
	}

	type ListType collections.List[CountVO]

	type customList1 struct {
		LstType ListType
	}
	type customList2 struct {
		LstType ListType
	}

	dto := customList1{
		LstType: ListType(collections.NewList[CountVO](CountVO{Count: 678})),
	}
	do := mapper.Single[customList2](dto)
	assert.Equal(t, dto.LstType.Index(0).Count, do.LstType.Index(0).Count)
}

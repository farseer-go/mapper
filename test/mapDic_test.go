package test

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mtdDO struct {
	A map[string]any
}

type mtdPO struct {
	A collections.Dictionary[string, string]
}

func TestMapDic(t *testing.T) {
	do := mtdDO{A: map[string]any{"1": "2"}}
	po := mapper.Single[mtdPO](do)

	assert.Equal(t, do.A["1"], po.A.GetValue("1"))

	po = mtdPO{A: collections.NewDictionaryFromMap(map[string]string{"1": "2"})}
	do = mapper.Single[mtdDO](po)
	assert.Equal(t, do.A["1"], po.A.GetValue("1"))
}

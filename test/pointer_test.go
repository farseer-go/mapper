package test

import (
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestPointerDO struct {
	Exception1 *trace.ExceptionStack
	Exception2 *trace.ExceptionStack
	Exception3 trace.ExceptionStack
	Exception4 *trace.ExceptionStack
	Exception5 *ExceptionStack
}
type TestPointerPO struct {
	Exception1 *trace.ExceptionStack
	Exception2 trace.ExceptionStack
	Exception3 *trace.ExceptionStack
	Exception4 *ExceptionStack
	Exception5 *trace.ExceptionStack
}
type ExceptionStack struct {
	ExceptionCallFile     string // 调用者文件路径
	ExceptionCallLine     int    // 调用者行号
	ExceptionCallFuncName string // 调用者函数名称
	ExceptionIsException  bool   // 是否执行异常
	ExceptionMessage      string // 异常信息
}

func TestPointer(t *testing.T) {
	do := TestPointerDO{
		Exception1: &trace.ExceptionStack{
			ExceptionCallFile:     "1",
			ExceptionCallLine:     2,
			ExceptionCallFuncName: "3",
			ExceptionIsException:  true,
			ExceptionMessage:      "4",
		},
		Exception2: &trace.ExceptionStack{
			ExceptionCallFile:     "1",
			ExceptionCallLine:     2,
			ExceptionCallFuncName: "3",
			ExceptionIsException:  true,
			ExceptionMessage:      "4",
		},
		Exception3: trace.ExceptionStack{
			ExceptionCallFile:     "1",
			ExceptionCallLine:     2,
			ExceptionCallFuncName: "3",
			ExceptionIsException:  true,
			ExceptionMessage:      "4",
		},
		Exception4: nil,
	}
	po := mapper.Single[TestPointerPO](do)

	assert.Equal(t, do.Exception1.ExceptionCallFile, po.Exception1.ExceptionCallFile)
	assert.Equal(t, do.Exception1.ExceptionCallLine, po.Exception1.ExceptionCallLine)
	assert.Equal(t, do.Exception1.ExceptionCallFuncName, po.Exception1.ExceptionCallFuncName)
	assert.Equal(t, do.Exception1.ExceptionIsException, po.Exception1.ExceptionIsException)
	assert.Equal(t, do.Exception1.ExceptionMessage, po.Exception1.ExceptionMessage)

	assert.Equal(t, do.Exception2.ExceptionCallFile, po.Exception2.ExceptionCallFile)
	assert.Equal(t, do.Exception2.ExceptionCallLine, po.Exception2.ExceptionCallLine)
	assert.Equal(t, do.Exception2.ExceptionCallFuncName, po.Exception2.ExceptionCallFuncName)
	assert.Equal(t, do.Exception2.ExceptionIsException, po.Exception2.ExceptionIsException)
	assert.Equal(t, do.Exception2.ExceptionMessage, po.Exception2.ExceptionMessage)

	assert.Equal(t, do.Exception3.ExceptionCallFile, po.Exception3.ExceptionCallFile)
	assert.Equal(t, do.Exception3.ExceptionCallLine, po.Exception3.ExceptionCallLine)
	assert.Equal(t, do.Exception3.ExceptionCallFuncName, po.Exception3.ExceptionCallFuncName)
	assert.Equal(t, do.Exception3.ExceptionIsException, po.Exception3.ExceptionIsException)
	assert.Equal(t, do.Exception3.ExceptionMessage, po.Exception3.ExceptionMessage)

	assert.True(t, po.Exception4 == nil)
	assert.True(t, po.Exception5 == nil)
}

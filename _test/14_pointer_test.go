package test

import (
	"testing"

	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
)

func TestPointer(t *testing.T) {
	type ExceptionStack struct {
		ExceptionCallFile     string // 调用者文件路径
		ExceptionCallLine     int    // 调用者行号
		ExceptionCallFuncName string // 调用者函数名称
		ExceptionIsException  bool   // 是否执行异常
		ExceptionMessage      string // 异常信息
	}
	type s1 struct {
		Exception1 *trace.ExceptionStack
		Exception2 *trace.ExceptionStack
		Exception3 trace.ExceptionStack
		Exception4 *trace.ExceptionStack
		Exception5 *ExceptionStack
	}
	type s2 struct {
		Exception1 *trace.ExceptionStack
		Exception2 trace.ExceptionStack
		Exception3 *trace.ExceptionStack
		Exception4 *ExceptionStack
		Exception5 *trace.ExceptionStack
	}
	do := s1{
		Exception1: &trace.ExceptionStack{
			ExceptionDetails: []trace.ExceptionStackDetail{trace.ExceptionStackDetail{
				ExceptionCallFile:     "1",
				ExceptionCallLine:     2,
				ExceptionCallFuncName: "3",
			}},
			ExceptionIsException: true,
			ExceptionMessage:     "4",
		},
		Exception2: &trace.ExceptionStack{
			ExceptionDetails: []trace.ExceptionStackDetail{trace.ExceptionStackDetail{
				ExceptionCallFile:     "11",
				ExceptionCallLine:     22,
				ExceptionCallFuncName: "33",
			}},
			ExceptionIsException: true,
			ExceptionMessage:     "44",
		},
		Exception3: trace.ExceptionStack{
			ExceptionDetails: []trace.ExceptionStackDetail{trace.ExceptionStackDetail{
				ExceptionCallFile:     "111",
				ExceptionCallLine:     222,
				ExceptionCallFuncName: "333",
			}},
			ExceptionIsException: true,
			ExceptionMessage:     "444",
		},
		Exception4: nil,
	}
	po := mapper.Single[s2](do)

	assert.Equal(t, do.Exception1.ExceptionDetails[0].ExceptionCallFile, po.Exception1.ExceptionDetails[0].ExceptionCallFile)
	assert.Equal(t, do.Exception1.ExceptionDetails[0].ExceptionCallLine, po.Exception1.ExceptionDetails[0].ExceptionCallLine)
	assert.Equal(t, do.Exception1.ExceptionDetails[0].ExceptionCallFuncName, po.Exception1.ExceptionDetails[0].ExceptionCallFuncName)
	assert.Equal(t, do.Exception1.ExceptionIsException, po.Exception1.ExceptionIsException)
	assert.Equal(t, do.Exception1.ExceptionMessage, po.Exception1.ExceptionMessage)

	assert.Equal(t, do.Exception2.ExceptionDetails[0].ExceptionCallFile, po.Exception2.ExceptionDetails[0].ExceptionCallFile)
	assert.Equal(t, do.Exception2.ExceptionDetails[0].ExceptionCallLine, po.Exception2.ExceptionDetails[0].ExceptionCallLine)
	assert.Equal(t, do.Exception2.ExceptionDetails[0].ExceptionCallFuncName, po.Exception2.ExceptionDetails[0].ExceptionCallFuncName)
	assert.Equal(t, do.Exception2.ExceptionIsException, po.Exception2.ExceptionIsException)
	assert.Equal(t, do.Exception2.ExceptionMessage, po.Exception2.ExceptionMessage)

	assert.Equal(t, do.Exception3.ExceptionDetails[0].ExceptionCallFile, po.Exception3.ExceptionDetails[0].ExceptionCallFile)
	assert.Equal(t, do.Exception3.ExceptionDetails[0].ExceptionCallLine, po.Exception3.ExceptionDetails[0].ExceptionCallLine)
	assert.Equal(t, do.Exception3.ExceptionDetails[0].ExceptionCallFuncName, po.Exception3.ExceptionDetails[0].ExceptionCallFuncName)
	assert.Equal(t, do.Exception3.ExceptionIsException, po.Exception3.ExceptionIsException)
	assert.Equal(t, do.Exception3.ExceptionMessage, po.Exception3.ExceptionMessage)

	assert.True(t, po.Exception4 == nil)
	assert.True(t, po.Exception5 == nil)
}

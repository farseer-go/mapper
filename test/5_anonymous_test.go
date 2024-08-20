package test

import (
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumCallType"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TraceDetailDatabase struct {
	trace.BaseTraceDetail
	DbName       string // 数据库名
	RowsAffected int64  // 影响行数
}
type BaseTraceDetailPO struct {
	TraceId   string                `gorm:"not null;default:'';comment:上下文ID"`
	Level     int                   `gorm:"not null;comment:当前层级（入口为0层）"`
	CallType  eumCallType.Enum      `gorm:"not null;comment:调用类型"`
	Timeline  time.Duration         `gorm:"not null;default:0;comment:从入口开始统计（微秒）"`
	Exception *trace.ExceptionStack `gorm:"json;not null;comment:异常信息"`
}
type TraceDetailDatabasePO struct {
	BaseTraceDetailPO `gorm:"embedded"`
	DbName            string `gorm:"not null;default:'';comment:数据库名"`
	RowsAffected      int64  `gorm:"not null;default:0;comment:影响行数"`
}

// 测试匿名类型
func TestAnonymous(t *testing.T) {
	po := TraceDetailDatabasePO{
		BaseTraceDetailPO: BaseTraceDetailPO{
			TraceId:  "123456",
			Level:    1,
			CallType: 2,
			Timeline: 3,
			Exception: &trace.ExceptionStack{
				ExceptionCallFile:     "4",
				ExceptionCallLine:     5,
				ExceptionCallFuncName: "6",
				ExceptionIsException:  true,
				ExceptionMessage:      "7",
			},
		},
		DbName:       "8",
		RowsAffected: 9,
	}
	do := mapper.Single[TraceDetailDatabase](po)

	assert.Equal(t, do.TraceId, po.TraceId)
	assert.Equal(t, do.Level, po.Level)
	assert.Equal(t, do.CallType, po.CallType)
	assert.Equal(t, do.Timeline, po.Timeline)
	assert.Equal(t, do.Exception.ExceptionCallFile, po.Exception.ExceptionCallFile)
	assert.Equal(t, do.Exception.ExceptionCallLine, po.Exception.ExceptionCallLine)
	assert.Equal(t, do.Exception.ExceptionCallFuncName, po.Exception.ExceptionCallFuncName)
	assert.Equal(t, do.Exception.ExceptionIsException, po.Exception.ExceptionIsException)
	assert.Equal(t, do.Exception.ExceptionMessage, po.Exception.ExceptionMessage)
	assert.Equal(t, do.DbName, po.DbName)
	assert.Equal(t, do.RowsAffected, po.RowsAffected)
}

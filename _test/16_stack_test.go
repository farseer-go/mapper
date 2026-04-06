package test

import (
	"testing"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
)

type Service interface {
}
type state int

type VO struct {
	Top       *VO
	Parent    *VO
	UserName  string
	Nick      string
	StateType state
	CreateAt  time.Time
	Login     Service
}

type LoginEO struct {
	*VO
	IsMaster       bool
	LoginName      string
	UserName       string
	Nick           string
	StateType      state
	LoginIp        string
	LoginAt        time.Time
	Permission     string
	PermissionList collections.List[int]
	LimitNameList  collections.List[string]
}

type EO struct {
	UserName        string
	Nick            string
	LoginName       string
	LoginPwd        string
	BindLoginIp     string
	ChangePwdAt     time.Time
	LoginIp         string
	LoginAt         time.Time
	LoginSalt       string
	IsOld           bool
	Permission      string
	LoginErrorCount int
	CreateAt        time.Time
	LimitNameList   []string
}

func TestStack(t *testing.T) {
	eo := EO{
		UserName:        "UserName",
		Nick:            "Nick",
		LoginName:       "LoginName",
		LoginPwd:        "LoginPwd",
		BindLoginIp:     "BindLoginIp",
		ChangePwdAt:     time.Now(),
		LoginIp:         "LoginIp",
		LoginAt:         time.Now(),
		LoginSalt:       "LoginSalt",
		IsOld:           true,
		Permission:      "Permission",
		LoginErrorCount: 1,
		CreateAt:        time.Now(),
		LimitNameList:   []string{"LimitNameList"},
	}
	loginEO := mapper.Single[LoginEO](eo)
	assert.Equal(t, eo.UserName, loginEO.UserName)
	assert.Equal(t, eo.Nick, loginEO.Nick)
	assert.Equal(t, eo.LoginName, loginEO.LoginName)
	assert.Equal(t, eo.LoginIp, loginEO.LoginIp)
	assert.Equal(t, eo.LoginAt.Unix(), loginEO.LoginAt.Unix())
	assert.Equal(t, eo.Permission, loginEO.Permission)
	assert.Equal(t, eo.LimitNameList[0], loginEO.LimitNameList.Index(0))

	loginEO.VO = &VO{
		UserName:  "UserName",
		Nick:      "Nick",
		StateType: 1,
		CreateAt:  time.Now(),
		Parent: &VO{
			UserName:  "UserName2",
			Nick:      "Nick2",
			StateType: 1,
			CreateAt:  time.Now(),
		},
	}
	eO := mapper.Single[EO](loginEO)
	assert.Equal(t, eo.UserName, eO.UserName)
	assert.Equal(t, eo.Nick, eO.Nick)
	assert.Equal(t, eo.LoginName, eO.LoginName)
	assert.Equal(t, eo.LoginIp, eO.LoginIp)
	assert.Equal(t, eo.LoginAt.Unix(), eO.LoginAt.Unix())
	assert.Equal(t, eo.Permission, eO.Permission)
	assert.Equal(t, eo.LimitNameList[0], eO.LimitNameList[0])
}

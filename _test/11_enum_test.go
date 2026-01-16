package test

import (
	"github.com/farseer-go/mapper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnum(t *testing.T) {
	type state int
	const (
		Running state = iota
		Pending
		Stopped
	)

	type sub struct {
		Stat state
	}
	type s1 struct {
		Status state
		Sub    sub
	}

	dto := s1{
		Status: Pending,
		Sub: sub{
			Stat: Pending,
		},
	}

	do := mapper.Single[s1](dto)
	assert.Equal(t, dto.Status, do.Status)
	assert.Equal(t, dto.Sub.Stat, do.Sub.Stat)
}

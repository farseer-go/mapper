package test

import (
	"context"
	"github.com/farseer-go/mapper"
	"testing"
)

func TestInterface(t *testing.T) {
	type s struct {
		Ctx context.Context
	}

	var s1 s
	var cancel context.CancelFunc
	s1.Ctx, cancel = context.WithCancel(context.Background())
	mapper.Single[s](s1)
	cancel()
}

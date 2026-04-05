package mapper

import (
	"fmt"
	"reflect"
	"testing"
)

func TestArenaUsage(t *testing.T) {
	// 模拟一个 UserVO 的 analysis
	type inner struct{ Count int }
	type sub struct {
		Id   int64
		Name string
	}
	type vo struct {
		Id    int64
		Name  string
		Sub   sub
		Arr   []string
		Count map[int]*inner
	}

	from := vo{Id: 1, Name: "x", Sub: sub{1, "y"}, Arr: []string{"a", "b"}, Count: map[int]*inner{1: {0}, 2: {0}}}
	fromVal := reflect.ValueOf(from)
	a := getArena()
	var obj analysisOjb
	obj.arena = a
	obj.source = make([]*valueMeta, 0, 32)
	obj.fromMeta = a.alloc()
	obj.fromMeta.setReflectValue(fromVal)
	obj.fromMeta.PointerMeta = obj.fromMeta.PointerMeta // skip fastReflect
	fmt.Printf("arena size=%d, pos after root=%d\n", len(a.buf), a.pos)
	putArena(a)
}

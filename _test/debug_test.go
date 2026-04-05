package mapper

import (
	"fmt"
	"reflect"
	"testing"
)

type DebugVO struct {
	A int
	B string
}

type DebugVO2 struct {
	A int
	B string
}

func TestDebugSourceSliceSize(t *testing.T) {
	from := DebugVO{A: 1, B: "x"}
	fromVal := reflect.ValueOf(from)
	var fAnalysis analysisOjb
	sourceSlice := fAnalysis.entry(fromVal)
	fmt.Printf("sourceSlice len: %d\n", len(sourceSlice))
	for _, m := range sourceSlice {
		fmt.Printf("  FullName=%q Type=%v\n", m.FullName, m.Type)
	}
}

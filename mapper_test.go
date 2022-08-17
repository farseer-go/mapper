package mapper

import (
	"github.com/farseer-go/collections"
	"testing"
)

type po struct {
	Name string
	Age  int
}
type do struct {
	Name string
	Age  int
}

func TestArray(t *testing.T) {
	arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
	arrDO := Array[do](arrPO)
	if len(arrPO) != len(arrDO) {
		t.Error()
	}

	for i := 0; i < len(arrPO); i++ {
		if arrPO[i].Name != arrDO[i].Name || arrPO[i].Age != arrDO[i].Age {
			t.Error()
		}
	}
}

func TestSingle(t *testing.T) {
	poSingle := po{Name: "steden", Age: 18}
	doSingle := Single[do](&poSingle)
	if poSingle.Name != doSingle.Name || poSingle.Age != doSingle.Age {
		t.Error()
	}
}

func TestPageList(t *testing.T) {
	arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
	lst := PageList[do](arrPO, 10)
	if len(arrPO) != lst.List.Count() {
		t.Error()
	}

	if lst.RecordCount != 10 {
		t.Error()
	}
	for i := 0; i < len(arrPO); i++ {
		if arrPO[i].Name != lst.List.Index(i).Name || arrPO[i].Age != lst.List.Index(i).Age {
			t.Error()
		}
	}
}

func TestToList(t *testing.T) {
	lstAny := collections.NewListAny(po{Name: "steden", Age: 18}, po{Name: "steden1", Age: 20})
	lstDO := ToList[do](lstAny)

	if lstAny.Count() != lstDO.Count() {
		t.Error()
	}

	for i := 0; i < lstAny.Count(); i++ {
		po := lstAny.Index(i).(po)
		if po.Name != lstDO.Index(i).Name || po.Age != lstDO.Index(i).Age {
			t.Error()
		}
	}
}

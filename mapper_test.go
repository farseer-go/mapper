package mapper

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/stretchr/testify/assert"
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

func TestMapperSingle(t *testing.T) {
	maps := make(map[string]string)
	maps["name"] = "steden"
	maps["age"] = "18"
	dic := collections.NewDictionaryFromMap(maps)
	dic.Add("name2", "harlen")
	arrDO := TaskDO{Id: 1, Client: ClientVO{Id: 2, Ip: "127.0.0.1", Name: "电脑"}, Status: Pending, Data: dic}
	arrDTO := TaskDTO{}
	err := MapDOtoDTO(&arrDO, &arrDTO)
	fmt.Println(arrDTO)
	println(err)
}

func TestArray(t *testing.T) {
	arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
	arrDO := Array[do](arrPO)
	assert.Equal(t, len(arrPO), len(arrDO))

	for i := 0; i < len(arrPO); i++ {
		assert.Equal(t, arrPO[i].Name, arrDO[i].Name)
		assert.Equal(t, arrPO[i].Age, arrDO[i].Age)
	}
}

func TestSingle(t *testing.T) {
	poSingle := po{Name: "steden", Age: 18}
	doSingle := Single[do](&poSingle)

	assert.Equal(t, poSingle.Name, doSingle.Name)
	assert.Equal(t, poSingle.Age, doSingle.Age)
}

func TestPageList(t *testing.T) {
	arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
	lst := ToPageList[do](arrPO, 10)

	assert.Equal(t, len(arrPO), lst.List.Count())

	assert.Equal(t, lst.RecordCount, int64(10))
	for i := 0; i < len(arrPO); i++ {
		assert.Equal(t, arrPO[i].Name, lst.List.Index(i).Name)
		assert.Equal(t, arrPO[i].Age, lst.List.Index(i).Age)
	}
}

func TestToList(t *testing.T) {
	lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden1", Age: 20})
	lstDO := ToList[do](lst)

	assert.Equal(t, lst.Count(), lstDO.Count())

	for i := 0; i < lst.Count(); i++ {
		assert.Equal(t, lst.Index(i).Name, lstDO.Index(i).Name)
		assert.Equal(t, lst.Index(i).Age, lstDO.Index(i).Age)
	}

	lstAny := lst.ToListAny()
	lstDO = ToList[do](lstAny)

	assert.Equal(t, lstAny.Count(), lstDO.Count())

	for i := 0; i < lstAny.Count(); i++ {
		po := lstAny.Index(i).(po)
		assert.Equal(t, po.Name, lstDO.Index(i).Name)
		assert.Equal(t, po.Age, lstDO.Index(i).Age)
	}

	arr := lst.ToArray()
	lstDO = ToList[do](arr)

	assert.Equal(t, len(arr), lstDO.Count())

	for i := 0; i < lstAny.Count(); i++ {
		assert.Equal(t, arr[i].Name, lstDO.Index(i).Name)
		assert.Equal(t, arr[i].Age, lstDO.Index(i).Age)
	}
}

func TestToListAny(t *testing.T) {
	arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}

	listAny := ToListAny(arrPO)

	assert.Equal(t, listAny.Count(), len(arrPO))
	for i := 0; i < listAny.Count(); i++ {
		po := listAny.Index(i).(po)

		assert.Equal(t, po.Name, arrPO[i].Name)
		assert.Equal(t, po.Age, arrPO[i].Age)
	}

	lst := collections.NewList(arrPO...)
	listAny = ToListAny(lst)

	assert.Equal(t, listAny.Count(), len(arrPO))
	for i := 0; i < listAny.Count(); i++ {
		po := listAny.Index(i).(po)

		assert.Equal(t, po.Name, arrPO[i].Name)
		assert.Equal(t, po.Age, arrPO[i].Age)
	}
}

func TestToMap(t *testing.T) {
	arrPO := po{Name: "steden", Age: 18}
	dic := ToMap[string, any](&arrPO)
	assert.Equal(t, "steden", dic["Name"])
	assert.Equal(t, 18, dic["Age"])
}

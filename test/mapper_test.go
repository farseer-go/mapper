package test

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
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
	var arrDO []TaskDO
	var arrDTO []TaskDTO
	arrDO = append(arrDO, TaskDO{Id: 1, Client: ClientVO{Id: 2, Ip: "127.0.0.1", Name: "电脑"}, Status: Pending, Data: dic})
	//Single[arrDTO](arrDO)
	fmt.Println(arrDTO)
}

func TestArray(t *testing.T) {
	arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
	arrDO := mapper.Array[do](arrPO)
	assert.Equal(t, len(arrPO), len(arrDO))

	for i := 0; i < len(arrPO); i++ {
		assert.Equal(t, arrPO[i].Name, arrDO[i].Name)
		assert.Equal(t, arrPO[i].Age, arrDO[i].Age)
	}
}

func TestSingle(t *testing.T) {

	dto := TaskDTO{
		Id:         1,
		ClientId:   1000,
		ClientIp:   "127.0.0.1",
		ClientName: "node",
		Status:     Pending,
		User: UserVO{
			Id:   88,
			Name: "steden",
		},
		Data: collections.NewDictionaryFromMap(map[string]string{"age": "18", "price": "88.88"}),
	}

	doSingle := mapper.Single[TaskDO](dto)
	doSingle2 := mapper.Single[TaskDTO](doSingle)
	assert.Equal(t, dto.Id, doSingle.Id)
	assert.Equal(t, dto.ClientId, doSingle.Client.Id)
	assert.Equal(t, doSingle2.Id, doSingle.Id)
	assert.Equal(t, doSingle2.ClientId, doSingle.Client.Id)
}

func TestPageList(t *testing.T) {
	arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
	lst := mapper.ToPageList[do](arrPO, 10)

	assert.Equal(t, len(arrPO), lst.List.Count())

	assert.Equal(t, lst.RecordCount, int64(10))
	for i := 0; i < len(arrPO); i++ {
		assert.Equal(t, arrPO[i].Name, lst.List.Index(i).Name)
		assert.Equal(t, arrPO[i].Age, lst.List.Index(i).Age)
	}
}

func TestToList(t *testing.T) {
	lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden1", Age: 20})
	lstDO := mapper.ToList[do](lst)

	assert.Equal(t, lst.Count(), lstDO.Count())

	for i := 0; i < lst.Count(); i++ {
		assert.Equal(t, lst.Index(i).Name, lstDO.Index(i).Name)
		assert.Equal(t, lst.Index(i).Age, lstDO.Index(i).Age)
	}

	lstAny := lst.ToListAny()
	lstDO = mapper.ToList[do](lstAny)

	assert.Equal(t, lstAny.Count(), lstDO.Count())

	for i := 0; i < lstAny.Count(); i++ {
		po := lstAny.Index(i).(po)
		assert.Equal(t, po.Name, lstDO.Index(i).Name)
		assert.Equal(t, po.Age, lstDO.Index(i).Age)
	}

	arr := lst.ToArray()
	lstDO = mapper.ToList[do](arr)

	assert.Equal(t, len(arr), lstDO.Count())

	for i := 0; i < lstAny.Count(); i++ {
		assert.Equal(t, arr[i].Name, lstDO.Index(i).Name)
		assert.Equal(t, arr[i].Age, lstDO.Index(i).Age)
	}
}

func TestToListAny(t *testing.T) {
	arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}

	listAny := mapper.ToListAny(arrPO)

	assert.Equal(t, listAny.Count(), len(arrPO))
	for i := 0; i < listAny.Count(); i++ {
		po := listAny.Index(i).(po)

		assert.Equal(t, po.Name, arrPO[i].Name)
		assert.Equal(t, po.Age, arrPO[i].Age)
	}

	lst := collections.NewList(arrPO...)
	listAny = mapper.ToListAny(lst)

	assert.Equal(t, listAny.Count(), len(arrPO))
	for i := 0; i < listAny.Count(); i++ {
		po := listAny.Index(i).(po)

		assert.Equal(t, po.Name, arrPO[i].Name)
		assert.Equal(t, po.Age, arrPO[i].Age)
	}
}

func TestToMap(t *testing.T) {
	arrPO := po{Name: "steden", Age: 18}
	dic := mapper.ToMap[string, any](&arrPO)
	assert.Equal(t, "steden", dic["Name"])
	assert.Equal(t, 18, dic["Age"])
}

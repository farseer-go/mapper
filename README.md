# mapper对象转换
> 包：`"github.com/farseer-go/mapper"`

- `文档`
    - [English](https://farseer-go.gitee.io/#/en-us/)
    - [中文](https://farseer-go.gitee.io/)
    - [English](https://farseer-go.github.io/doc/#/en-us/)
- `源代码`
    - [github](https://github.com/farseer-go/fs)
  

![](https://img.shields.io/github/stars/farseer-go?style=social)
![](https://img.shields.io/github/license/farseer-go/mapper)
![](https://img.shields.io/github/go-mod/go-version/farseer-go/mapper)
![](https://img.shields.io/github/v/release/farseer-go/mapper)
![Codecov](https://img.shields.io/codecov/c/github/farseer-go/mapper)
![](https://img.shields.io/github/languages/code-size/farseer-go/mapper)
![](https://img.shields.io/github/directory-file-count/farseer-go/mapper)
![](https://goreportcard.com/badge/github.com/farseer-go/mapper)

## 概述
我们在开发当中经常会遇到实体与实体，实体集合与实体集合之间的相互转换，以及各种各样的相互转换问题。

该组件很好的解决了上述转换过程中的人为赋值问题。

组件支持 Map之间的转换，支持 数组之间的转换，支持 数组转List ，支持实体转Map ，支持List转List。



## 1、方法Single

_函数定义_
```go
// DO和DTO之间相互转换
func Single[TEntity any](fromObjPtr any) TEntity {
```
_用法：_
```go

type State int

const (
    Running State = iota
    Pending
    Stopped
)

func (s State) String() string {
    switch s {
        case Running:
        return "Running"
        case Pending:
        return "Pending"
        case Stopped:
        return "Stopped"
        default:
        return "Unknown"
    }
}

type ClientVO struct {
    Id   int64
    Ip   string
    Name string
}

type UserVO struct {
    Id   int64
    Name string
}

type TaskDO struct {
    Id       int
    Client   ClientVO
    Status   State
    UserId   int64
    UserName string
    Data     collections.Dictionary[string, string]
}

type TaskDTO struct {
    Id         int
    ClientId   int64
    ClientIp   string
    ClientName string
    Status     State
    User       UserVO
    Data       collections.Dictionary[string, string]
}

//DO 转 DTO
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

do := Single[TaskDO](dto)


//DTO 转 DO
dtoSingle := Single[TaskDTO](do)


```

## 2、方法ToMap

_函数定义_
```go
// Map转Map
func ToMap[K comparable, V any](fromObjPtr any) map[K]V
```
_用法：_
```go
type po struct {
    Name string
    Age  int
}
type do struct {
    Name string
    Age  int
}

//实体对象转map
arrPO := po{Name: "steden", Age: 18}
map := ToMap[string, any](&arrPO)
```

## 3、方法Array

_函数定义_
```go
// 数组转换
func Array[T any](fromSlice any) []T
```
_用法：_
```go
type po struct {
    Name string
    Age  int
}
type do struct {
    Name string
    Age  int
}

arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
arrDO := Array[do](arrPO)

assert.Equal(t, len(arrPO), len(arrDO))

for i := 0; i < len(arrPO); i++ {
    assert.Equal(t, arrPO[i].Name, arrDO[i].Name)
    assert.Equal(t, arrPO[i].Age, arrDO[i].Age)
}

```

## 4、方法ToPageList

_函数定义_
```go
// 转成分页集合
func ToPageList[TEntity any](sliceOrListOrListAny any, recordCount int64) collections.PageList[TEntity]
```
_用法：_
```go
type po struct {
    Name string
    Age  int
}
type do struct {
    Name string
    Age  int
}

arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
lst := ToPageList[do](arrPO, 10)

//返回数据数量
lst.RecordCount

```

## 5、方法ToList

_函数定义_
```go
// 转换成List
func ToList[TEntity any](sliceOrListOrListAny any) collections.List[TEntity] 
```
_用法：_
```go
type po struct {
    Name string
    Age  int
}
type do struct {
    Name string
    Age  int
}

lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden1", Age: 20})
lstDO := ToList[do](lst)

assert.Equal(t, lst.Count(), lstDO.Count())

for i := 0; i < lst.Count(); i++ {
    assert.Equal(t, lst.Index(i).Name, lstDO.Index(i).Name)
    assert.Equal(t, lst.Index(i).Age, lstDO.Index(i).Age)
}

//List转List
lstAny := lst.ToListAny()
lstDO = ToList[do](lstAny)

assert.Equal(t, lstAny.Count(), lstDO.Count())

for i := 0; i < lstAny.Count(); i++ {
    po := lstAny.Index(i).(po)
    assert.Equal(t, po.Name, lstDO.Index(i).Name)
    assert.Equal(t, po.Age, lstDO.Index(i).Age)
}
//数组转List
arr := lst.ToArray()
lstDO = ToList[do](arr)

assert.Equal(t, len(arr), lstDO.Count())

```
## 6、方法ToListAny

_函数定义_
```go
// 切片转ToListAny
func ToListAny(sliceOrList any) collections.ListAny
```
_用法：_
```go
type po struct {
    Name string
    Age  int
}
type do struct {
    Name string
    Age  int
}

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
```

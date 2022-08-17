## What are the functions?
* mapper（对象转换）
    * func
        * Array （数组转换）
        * Single （单个转换）
        * PageList （转换成core.PageList）
        * ToList （ListAny转List泛型）

## Getting Started
```go
type po struct {
    Name string
    Age  int
}
type do struct {
    Name string
    Age  int
}
```
### Array
slice struct mapTo slice struct
```go
arrPO := []PO{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
Array[DO](arrPO)       // return []DO{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
```

### Single
struct mapTo struct
```go
poSingle := po{Name: "steden", Age: 18}
Single[do](&poSingle)   // return do{Name: "steden", Age: 18}
```

### PageList
slice struct mapTo collections.PageList
```go
arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
PageList[po](arrPO, 10) // return collections.PageList[do]
```

### ToList
collections.ListAny mapTo collections.List
```go
lstAny := collections.NewListAny(po{Name: "steden", Age: 18}, po{Name: "steden1", Age: 20})
ToList[do](lstAny)      // return collections.List[do]
```
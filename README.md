# mapper
对象转换

## What are the functions?
* mapper
    * func
        * Array （数组转换）
        * Single （单个转换）
        * ToPageList （转换成core.PageList）
        * ToList （ListAny、List[xx]、[]xx转List[yy]）
        * ToListAny （切片、List转ToListAny）
        * ToMap （结构体转Map）

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
`slice` mapTo `slice`
```go
arrPO := []PO{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
mapper.Array[DO](arrPO)       // return []DO{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
```

### Single
`struct` mapTo `struct`
```go
poSingle := po{Name: "steden", Age: 18}
mapper.Single[do](&poSingle)   // return do{Name: "steden", Age: 18}
```

### ToMap
`struct` mapTo `map`
```go
arrPO := po{Name: "steden", Age: 18}
mapper.ToMap[string, any](&arrPO)  // return map["Name"] = "steden", map["Age"] = 18
```

### ToPageList
`slice` mapTo `collections.PageList`
```go
arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
mapper.ToPageList[po](arrPO, 10) // return collections.PageList[do]
```

### ToList
`ListAny、List[xx]、[]xx` mapTo `List[yy]`
```go
lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden1", Age: 20})
mapper.ToList[do](lst)         // return collections.List[do]

lstAny := lst.ToListAny()
mapper.ToList[do](lstAny)      // return collections.List[do]

arr := lst.ToArray()
mapper.ToList[do](arr)         // return collections.List[do]
```

### ToListAny
`slice` mapTo `collections.ListAny`
```go
arrPO := []po{{Name: "steden", Age: 18}, {Name: "steden1", Age: 20}}
mapper.ToListAny(arrPO)     // return collections.ListAny

lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden1", Age: 20})
mapper.ToListAny(lst)       // return collections.ListAny
```
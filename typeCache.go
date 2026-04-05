package mapper

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/farseer-go/fs/fastReflect"
)

// typeCache 缓存类型映射关系
var typeCache sync.Map // map[typePair]*typeMeta

// typePair 类型对
type typePair struct {
	from reflect.Type
	to   reflect.Type
}

// typeMeta 类型元数据缓存
type typeMeta struct {
	isSimpleStruct bool           // 是否为简单结构体（全基础类型）
	fieldMappings  []fieldMapping // 字段映射关系
	sourceFields   []*valueMeta   // 源字段元数据（预分析）
}

// fieldMapping 字段映射
type fieldMapping struct {
	srcIndex int         // 源字段索引
	dstIndex int         // 目标字段索引
	srcType  reflect.Type
	dstType  reflect.Type
	copyFunc func(src, dst reflect.Value) // 拷贝函数
}

// getOrCreateTypeMeta 获取或创建类型元数据
func getOrCreateTypeMeta(fromType, toType reflect.Type) *typeMeta {
	pair := typePair{from: fromType, to: toType}

	if cached, ok := typeCache.Load(pair); ok {
		return cached.(*typeMeta)
	}

	// 创建新的类型元数据
	meta := &typeMeta{}
	meta.analyze(fromType, toType)

	typeCache.Store(pair, meta)
	return meta
}

// analyze 分析类型映射关系
func (tm *typeMeta) analyze(fromType, toType reflect.Type) {
	// 去除指针
	if fromType.Kind() == reflect.Ptr {
		fromType = fromType.Elem()
	}
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		tm.isSimpleStruct = false
		return
	}

	// 构建目标字段名映射
	toFieldMap := make(map[string]int)
	for i := 0; i < toType.NumField(); i++ {
		field := toType.Field(i)
		if field.IsExported() {
			toFieldMap[field.Name] = i
		}
	}

	// 检查是否为简单结构体（所有字段名完全匹配且类型相同）
	tm.isSimpleStruct = true
	tm.fieldMappings = make([]fieldMapping, 0, fromType.NumField())

	for i := 0; i < fromType.NumField(); i++ {
		fromField := fromType.Field(i)
		if !fromField.IsExported() {
			continue
		}

		// 查找匹配的目标字段
		if toIndex, exists := toFieldMap[fromField.Name]; exists {
			toField := toType.Field(toIndex)

			mapping := fieldMapping{
				srcIndex: i,
				dstIndex: toIndex,
				srcType:  fromField.Type,
				dstType:  toField.Type,
			}

			// 判断字段类型
			srcPointerMeta := fastReflect.PointerOf(fromField.Type)
			if srcPointerMeta.Type != fastReflect.GoBasicType {
				tm.isSimpleStruct = false
			}

			// 设置拷贝函数
			if fromField.Type == toField.Type {
				// 类型相同，直接赋值
				mapping.copyFunc = func(src, dst reflect.Value) {
					dst.Set(src)
				}
			} else {
				// 类型不同，需要转换
				tm.isSimpleStruct = false
				mapping.copyFunc = nil
			}

			tm.fieldMappings = append(tm.fieldMappings, mapping)
		} else {
			// 字段名不匹配，不是简单结构体
			tm.isSimpleStruct = false
		}
	}
}

// fastCopyStruct 快速拷贝结构体（仅适用于简单结构体）
func (tm *typeMeta) fastCopyStruct(from reflect.Value, to reflect.Value) {
	if !tm.isSimpleStruct {
		return
	}

	fromVal := from
	toVal := to

	// 去除指针
	if fromVal.Kind() == reflect.Ptr {
		fromVal = fromVal.Elem()
	}
	if toVal.Kind() == reflect.Ptr {
		toVal = toVal.Elem()
	}

	// 逐字段拷贝
	for _, mapping := range tm.fieldMappings {
		srcField := fromVal.Field(mapping.srcIndex)
		dstField := toVal.Field(mapping.dstIndex)

		if mapping.copyFunc != nil {
			mapping.copyFunc(srcField, dstField)
		}
	}
}

// canUseFastPath 判断是否可以使用快速路径
func canUseFastPath(fromType, toType reflect.Type) bool {
	// 类型完全相同
	if fromType == toType {
		return true
	}

	// 去除指针后类型相同
	if fromType.Kind() == reflect.Ptr {
		fromType = fromType.Elem()
	}
	if toType.Kind() == reflect.Ptr {
		toType = toType.Elem()
	}

	if fromType == toType {
		return true
	}

	// 都是结构体，检查缓存
	if fromType.Kind() == reflect.Struct && toType.Kind() == reflect.Struct {
		meta := getOrCreateTypeMeta(fromType, toType)
		return meta.isSimpleStruct
	}

	return false
}

// fastCopy 快速拷贝（使用 unsafe 优化）
func fastCopy[TEntity any](object any) TEntity {
	var toObj TEntity

	fromVal := reflect.ValueOf(object)
	toVal := reflect.ValueOf(&toObj).Elem()

	fromType := fromVal.Type()
	toType := toVal.Type()

	// 类型完全相同，直接转换
	if fromType == toType {
		return object.(TEntity)
	}

	// 使用类型缓存
	meta := getOrCreateTypeMeta(fromType, toType)
	if meta.isSimpleStruct {
		meta.fastCopyStruct(fromVal, toVal)
		return toObj
	}

	// 降级到普通拷贝
	_ = auto(fromVal, &toObj)
	return toObj
}

// metaArena 是每次 auto() 调用使用的局部分配器
// 预分配一个大 slice，按需 bump-allocate，调用结束后整体归还，避免 per-object pool.Get 竞争
type metaArena struct {
	buf  []valueMeta // 预分配的内存块
	pos  int         // 当前分配位置
	ptrs []*valueMeta // 分配出去的指针（用于 source slice）
}

const arenaSize = 1024 // 每个 arena 预分配的 valueMeta 数量

// arenaPool 缓存 metaArena 对象
var arenaPool = sync.Pool{
	New: func() interface{} {
		return &metaArena{
			buf:  make([]valueMeta, arenaSize),
			ptrs: make([]*valueMeta, 0, arenaSize),
		}
	},
}

// getArena 从 pool 获取一个 arena
func getArena() *metaArena {
	a := arenaPool.Get().(*metaArena)
	a.pos = 0
	a.ptrs = a.ptrs[:0]
	return a
}

// putArena 归还 arena 到 pool
func putArena(a *metaArena) {
	if a != nil {
		arenaPool.Put(a)
	}
}

// alloc 从 arena 分配一个 valueMeta
func (a *metaArena) alloc() *valueMeta {
	if a.pos < len(a.buf) {
		m := &a.buf[a.pos]
		a.pos++
		// 重置字段
		m.Parent = nil
		m.FullName = ""
		m.FullHash = 0
		m.FieldName = ""
		m.IsNil = false
		m.IsAnonymous = false
		m.IsIgnore = false
		m.Level = 0
		return m
	}
	// arena 满了，直接 new（极少发生）
	return &valueMeta{}
}

// 全局当前调用的 arena，通过 analysisOjb/assignObj 传递
// （不用 goroutine-local 因为实现复杂，直接在结构体里携带）

// 保留以下函数供兼容，内部改为 arena 分配
// getMetaFromPool 兼容旧接口（实际由 arena 替代）
func getMetaFromPool() *valueMeta {
	return &valueMeta{}
}

// putMetaToPool 兼容旧接口（arena 模式下无需单个归还）
func putMetaToPool(meta *valueMeta) {
	// no-op in arena mode
}

// sourceMapPool 复用 map[uint64]*valueMeta，避免每次 make
var sourceMapPool = sync.Pool{
	New: func() interface{} {
		m := make(map[uint64]*valueMeta, 64)
		return &m
	},
}

func getSourceMap() map[uint64]*valueMeta {
	mp := sourceMapPool.Get().(*map[uint64]*valueMeta)
	m := *mp
	// 清空 map 保留容量
	for k := range m {
		delete(m, k)
	}
	return m
}

func putSourceMap(m map[uint64]*valueMeta) {
	if m != nil {
		sourceMapPool.Put(&m)
	}
}

// slicePool 用于缓存 source slice
var slicePool = sync.Pool{
	New: func() interface{} {
		s := make([]*valueMeta, 0, 128)
		return &s
	},
}

// getSliceFromPool 从对象池获取切片
func getSliceFromPool() *[]*valueMeta {
	slice := slicePool.Get().(*[]*valueMeta)
	*slice = (*slice)[:0] // 清空但保留容量
	return slice
}

// putSliceToPool 归还切片到对象池
func putSliceToPool(slice *[]*valueMeta) {
	if slice != nil {
		slicePool.Put(slice)
	}
}

// 用于 unsafe 操作的辅助函数
func unsafeCopyBasicField(src, dst unsafe.Pointer, size uintptr) {
	// 直接内存拷贝
	*(*[]byte)(unsafe.Pointer(&dst)) = *(*[]byte)(unsafe.Pointer(&src))
}

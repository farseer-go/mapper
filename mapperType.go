package mapper

type FieldType int

const (
	List        FieldType = iota // 集合类型
	PageList    FieldType = iota // 集合类型
	CustomList                   // 自定义List类型
	Slice                        // 结构体类型
	ArrayType                    // 结构体类型
	Map                          // map类型
	Dic                          // 字典类型
	GoBasicType                  // 基础类型
	Struct                       // 结构体类型
	Chan                         // 结构体类型
	Interface                    // 结构体类型
	Func                         // 结构体类型
	Invalid                      // 结构体类型
	Unknown                      // 结构体类型
)

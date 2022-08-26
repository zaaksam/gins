package orm

// FieldType 字段类型约束
type FieldType interface {
	~string | ~bool | FieldTypeInt | FieldTypeUint | FieldTypeFloat
}

// FieldTypeInt 字段类型约束 - 整型
type FieldTypeInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// FieldTypeUint 字段类型约束 - 无符号整型
type FieldTypeUint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// FieldTypeFloat 字段类型约束 - 浮点
type FieldTypeFloat interface {
	~float32 | ~float64
}

package orm

// 基础数据模型接口
type IModel interface {
	// 本服务的 db 接口要求
	DatabaseAlias() string
	FieldNames(isSetValOpt ...bool) []string
	FieldReset()

	// xorm 的接口要求
	TableName() string
}

// Model 基础数据模型
type Model struct {
	fields map[IField]struct{}
}

func (m *Model) setField(field IField) {
	if m.fields == nil {
		m.fields = make(map[IField]struct{})
	}

	m.fields[field] = struct{}{}
}

// FieldNames 获取字段列表，可指定是否赋值
func (m *Model) FieldNames(isSetValOpt ...bool) (names []string) {
	names = make([]string, 0, len(m.fields))

	for field := range m.fields {
		if len(isSetValOpt) == 1 && isSetValOpt[0] != field.IsSetVal() {
			continue
		}

		names = append(names, field.FieldName())
	}

	return
}

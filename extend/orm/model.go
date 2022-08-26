package orm

// 基础数据模型接口
type IModel interface {
	// 本服务的 db 接口要求
	DatabaseAlias() string
	FieldNames(isSetValOpt ...bool) []string

	// xorm 的接口要求
	TableName() string
}

// Model 基础数据模型
type Model struct {
	fields map[string]bool // name:isSetVal
}

// 仅初始化一次
func (m *Model) onceInit() {
	if m.fields != nil {
		return
	}

	m.fields = make(map[string]bool)
}

// 设置字段为赋值状态
func (m *Model) setVal(fieldName string) {
	m.onceInit()

	m.fields[fieldName] = true
}

// 设置字段为非赋值状态
func (m *Model) unsetVal(fieldName string) {
	m.onceInit()

	m.fields[fieldName] = false
}

// 判断字段是否赋值
func (m *Model) isSetVal(fieldName string) (isSetVal bool) {
	m.onceInit()

	isSetVal = m.fields[fieldName]
	return
}

// FieldNames 获取字段列表，可指定是否赋值
func (m *Model) FieldNames(isSetValOpt ...bool) (names []string) {
	m.onceInit()

	names = make([]string, 0, len(m.fields))

	for name, isSetVal := range m.fields {
		if len(isSetValOpt) == 1 && isSetValOpt[0] != isSetVal {
			continue
		}

		names = append(names, name)
	}

	return
}

package orm

type IModel interface {
	// xorm 的接口要求
	TableName() string

	// 本服务的接口要求
	DatabaseAlias() string

	// 以下接口要求为 github.com/amoeba-inc/easyjson 生成
	FieldMarks() []string
}

// ModelList 通用模型列表
type ModelList[T any] struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"pageSize"`
	PageCount int   `json:"pageCount"`
	Total     int64 `json:"total"`
	Items     []T   `json:"items"`
}

func init() {

}

// Model orm基础模型
// 与定制的 easyjson 搭配使用
type Model struct {
	fieldMark    map[string]bool
	propertyMark map[string]bool
}

// FieldMarks 列出所有已赋值的字段名称列表
func (m *Model) FieldMarks() []string {
	names := make([]string, 0, len(m.fieldMark))
	for k, v := range m.fieldMark {
		if v {
			names = append(names, k)
		}
	}

	return names
}

// HasFieldMark 指定的字段名称是否已赋值
func (m *Model) HasFieldMark(fieldName string) bool {
	if m.fieldMark == nil {
		return false
	}

	return m.fieldMark[fieldName]
}

// SetFieldMark 设置字段的赋值标识，isMark不传时，默认:true
func (m *Model) SetFieldMark(fieldName string, isMark ...bool) {
	if m.fieldMark == nil {
		return
	}

	if len(isMark) == 1 {
		m.fieldMark[fieldName] = isMark[0]
		return
	}

	m.fieldMark[fieldName] = true
}

// ResetFieldMark 重置所有字段的赋值标识为:false，字段内容并不会清空
func (m *Model) ResetFieldMark() {
	if m.fieldMark == nil {
		m.fieldMark = make(map[string]bool)
		return
	}

	for k := range m.fieldMark {
		m.fieldMark[k] = false
	}
}

// PropertyMarks 列出所有已赋值的属性名称列表
func (m *Model) PropertyMarks() []string {
	names := make([]string, 0, len(m.propertyMark))
	for k, v := range m.propertyMark {
		if v {
			names = append(names, k)
		}
	}

	return names
}

// HasPropertyMark 指定的字段名称是否已赋值
func (m *Model) HasPropertyMark(markKey string) bool {
	if m.propertyMark == nil {
		return false
	}

	return m.propertyMark[markKey]
}

// SetPropertyMark 设置字段的赋值标识，isMark不传时，默认:true
func (m *Model) SetPropertyMark(markKey string, isMark ...bool) {
	if m.propertyMark == nil {
		return
	}

	if len(isMark) == 1 {
		m.propertyMark[markKey] = isMark[0]
		return
	}

	m.propertyMark[markKey] = true
}

// ResetPropertyMark 重置所有字段的赋值标识为:false，字段内容并不会清空
func (m *Model) ResetPropertyMark() {
	if m.propertyMark == nil {
		m.propertyMark = make(map[string]bool)
		return
	}

	for k := range m.propertyMark {
		m.propertyMark[k] = false
	}
}

// ResetMark 重置所有字段、属于的赋值标识为:false，字段内容并不会清空
func (m *Model) ResetMark() {
	m.ResetFieldMark()
	m.ResetPropertyMark()
}

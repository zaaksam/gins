package orm

import "xorm.io/xorm/convert"

var _ convert.Conversion = &Field[uint8]{}

// ====== 实现 xorm convert.Conversion 接口 ======

// ToDB 实现
func (f *Field[T]) ToDB() ([]byte, error) {
	return f.marshal(false)
}

// FromDB 实现
func (f *Field[T]) FromDB(body []byte) error {
	return f.unmarshal(false, body)
}

package orm

import "encoding/json"

var (
	_ json.Marshaler   = &Field[uint8]{}
	_ json.Unmarshaler = &Field[uint8]{}
)

// MarshalJSON 实现 json.Marshaler 接口
func (f *Field[T]) MarshalJSON() ([]byte, error) {
	return f.marshal(true)
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (f *Field[T]) UnmarshalJSON(body []byte) error {
	return f.unmarshal(true, body)
}

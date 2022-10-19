package orm

import (
	"errors"
	"fmt"
)

const field_json_string_quote = `"`

// 基础字段接口
type IField interface {
	IsSetVal() bool
	FieldName() string
}

// Field 字段
type Field[T FieldType] struct {
	model           *Model
	isJSONTagString bool // json tag 是否有 string 标记
	jsonName        string
	fieldName       string
	isSetVal        bool
	val             T
}

// Reset 重设配置
func (f *Field[T]) Reset(model *Model, fieldName, jsonName string, isJSONTagString bool) {
	f.model = model
	f.fieldName = fieldName
	f.jsonName = jsonName
	f.isJSONTagString = isJSONTagString

	f.model.setField(f)
}

// IsSetVal 是否赋值
func (f *Field[T]) IsSetVal() bool {
	return f.isSetVal
}

// SetVal 设定赋值，不指定 valOpt 时仅为设置赋值状态
func (f *Field[T]) SetVal(valOpt ...T) {
	if len(valOpt) == 1 {
		f.val = valOpt[0]
	}

	f.isSetVal = true
}

// UnsetVal 取消赋值，不指定 valOpt 时仅为取消赋值状态
func (f *Field[T]) UnsetVal(valOpt ...T) {
	if len(valOpt) == 1 {
		f.val = valOpt[0]
	}

	f.isSetVal = false
}

// Val 值
func (f *Field[T]) Val() T {
	return f.val
}

// JSONName 成员名
func (f *Field[T]) JSONName() string {
	return f.jsonName
}

// FieldName 字段名
func (f *Field[T]) FieldName() string {
	return f.fieldName
}

// 值转化为 bytes
func (f *Field[T]) marshal(isJSONFormat bool) (body []byte, err error) {
	// logger.Debugf("marshal：%s=%s", f.fieldName, fmt.Sprintf("%v", f.val))

	var ti any = &f.val
	switch v := ti.(type) {
	case *string:
		if isJSONFormat {
			body = stringToJSONValue(*v)
		} else {
			body = []byte(*v)
		}
	case *bool:
		body = boolToBytes(*v)
	case *int:
		body = intToBytes(*v)
	case *int8:
		body = intToBytes(*v)
	case *int16:
		body = intToBytes(*v)
	case *int32:
		body = intToBytes(*v)
	case *int64:
		body = intToBytes(*v)
	case *uint:
		body = uintToBytes(*v)
	case *uint8:
		body = uintToBytes(*v)
	case *uint16:
		body = uintToBytes(*v)
	case *uint32:
		body = uintToBytes(*v)
	case *uint64:
		body = uintToBytes(*v)
	case *float32:
		body = floatToBytes(*v, 32)
	case *float64:
		body = floatToBytes(*v, 64)
	}

	if isJSONFormat && f.isJSONTagString {
		newBody := make([]byte, 0, len(body)+2)
		newBody = append(newBody, []byte(field_json_string_quote)...)
		newBody = append(newBody, body...)
		newBody = append(newBody, []byte(field_json_string_quote)...)

		body = newBody
	}

	return
}

// bytes 转化为值，相当于赋值
func (f *Field[T]) unmarshal(isJSONFormat bool, body []byte) (err error) {
	// logger.Debugf("unmarshal：%s=%s", f.fieldName, string(body))

	if isJSONFormat && f.isJSONTagString {
		body = body[1 : len(body)-1]
	}

	var ti any = &f.val
	switch v := ti.(type) {
	case *string:
		if isJSONFormat {
			str := string(body)
			l := len(str)
			if str[0:1] != field_json_string_quote && str[l-1:l] != field_json_string_quote {
				if f.jsonName != "" {
					err = fmt.Errorf("%s's value invalid, it must be string", f.jsonName)
				} else {
					err = errors.New("json string value invalid")
				}
				return
			}

			*v = str[1 : l-1]
		} else {
			*v = string(body)
		}
	case *bool:
		*v, err = bytesToBool(body)
	case *int:
		*v, err = bytesToInt[int](body)
	case *int8:
		*v, err = bytesToInt[int8](body)
	case *int16:
		*v, err = bytesToInt[int16](body)
	case *int32:
		*v, err = bytesToInt[int32](body)
	case *int64:
		*v, err = bytesToInt[int64](body)
	case *uint:
		*v, err = bytesToUint[uint](body)
	case *uint8:
		*v, err = bytesToUint[uint8](body)
	case *uint16:
		*v, err = bytesToUint[uint16](body)
	case *uint32:
		*v, err = bytesToUint[uint32](body)
	case *uint64:
		*v, err = bytesToUint[uint64](body)
	case *float32:
		*v, err = bytesToFloat[float32](body, 32)
	case *float64:
		*v, err = bytesToFloat[float64](body, 64)
	}

	// 上面已赋值，只需设置赋值状态
	f.SetVal()

	return err
}

package orm

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var (
	_ ValidatorValuer = Field[uint8]{}
)

// ValidatorValuer 自定义 validate 字段类型接口
type ValidatorValuer interface {
	ValidatorValue() (value any, isSetVal bool)
}

// ValidatorValue 实现 ValidatorValuer 接口
func (f Field[T]) ValidatorValue() (value any, isSetVal bool) {
	value, isSetVal = f.val, f.getModel().isSetVal(f.fieldName)
	return
}

// RegisterWithValidator 注册到 Validator
func RegisterWithValidator(validate *validator.Validate) {

	// 注册 FieldType 类型
	validate.RegisterCustomTypeFunc(
		validateValuer,
		Field[string]{},
		Field[bool]{},
		Field[int]{},
		Field[int8]{},
		Field[int16]{},
		Field[int32]{},
		Field[int64]{},
		Field[uint]{},
		Field[uint8]{},
		Field[uint16]{},
		Field[uint32]{},
		Field[uint64]{},
		Field[float32]{},
		Field[float64]{},
	)
}

// validateValuer 实现 validator.CustomTypeFunc 函数
func validateValuer(field reflect.Value) interface{} {
	valuer, ok := field.Interface().(ValidatorValuer)
	if !ok {
		return nil
	}

	val, isSetVal := valuer.ValidatorValue()
	if !isSetVal {
		return nil
	}

	return val
}

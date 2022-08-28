package errors

import "errors"

// 兼容 errors.New 原生函数
func New(text string) error {
	return errors.New(text)
}

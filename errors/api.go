package errors

import (
	"fmt"

	"github.com/zaaksam/gins/constant"
	"github.com/zaaksam/gins/extend/nanoid"
)

type APIError struct {
	Code constant.APICodeType
	Msg  string
	Data any

	err error
}

// Error 实现 errors 接口
func (ae *APIError) Error() string {
	return ae.Msg
}

// Unwrap 实现 errors 接口
func (ae *APIError) Unwrap() error {
	return ae.err
}

func newAPIError(err error, msg string, codeOpt ...constant.APICodeType) (ae *APIError) {
	ae = &APIError{
		Msg: msg,
		err: err,
	}

	if len(codeOpt) == 1 {
		ae.Code = codeOpt[0]
	} else {
		ae.Code = constant.API_ERROR
	}

	return ae
}

// NewAPIError 创建常规的 APIError，code 默认：API_ERROR
func NewAPIError(msg string, codeOpt ...constant.APICodeType) (ae *APIError) {
	ae = newAPIError(nil, msg, codeOpt...)
	return
}

// NewAPIErrorWrap 创建包装过原始错误信息的 APIError ，code 默认：API_CRASH
func NewAPIErrorWrap(err error) (ae *APIError) {
	msg := fmt.Sprintf("系统内部错误 [%s]", nanoid.New())
	ae = newAPIError(err, msg, constant.API_CRASH)
	return
}

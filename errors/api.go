package errors

import (
	"github.com/zaaksam/gins/constant"
	"github.com/zaaksam/gins/extend/nanoid"
)

type APIError struct {
	Code constant.APICodeType
	Msg  string
	Data any

	err   error
	stack string
}

// Error 实现 errors 接口
func (ae *APIError) Error() string {
	return ae.Msg
}

// Unwrap 实现 errors 接口
func (ae *APIError) Unwrap() error {
	return ae.err
}

func unwrapError(err error) (str string) {
	if err == nil {
		return
	}

	ae, ok := err.(*APIError)
	if !ok {
		str = err.Error()
		return
	}

	str = unwrapError(ae.err)
	if str == "" {
		str = ae.Msg
	} else {
		str = ae.Msg + "：\n" + str
	}
	if ae.stack != "" {
		str += ae.stack
	}
	return
}

// UnwrapError 输出所有包装层的错误信息
func (ae *APIError) UnwrapError() string {
	return unwrapError(ae)
}

func newAPIError(err error, msg, stack string, codeOpt ...constant.APICodeType) (ae *APIError) {
	ae = &APIError{
		Msg:   msg,
		err:   err,
		stack: stack,
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
	ae = newAPIError(nil, msg, "", codeOpt...)
	return
}

// NewAPIErrorWrap 创建包装过原始错误信息的 APIError，并提供 stack ，code 默认：API_CRASH
func NewAPIErrorWrap(err error) (ae *APIError) {
	ae, ok := err.(*APIError)
	if ok {
		if ae.Unwrap() != nil {
			return
		}
	}

	msg := "系统内部错误 [" + nanoid.New() + "]"
	stack := callers(2)
	ae = newAPIError(err, msg, stack, constant.API_CRASH)
	return
}

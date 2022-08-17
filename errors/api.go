package errors

import (
	"fmt"
	"time"

	"github.com/zaaksam/gins/constant"
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
func NewAPIError(msg string, codeOpt ...constant.APICodeType) (ae error) {
	ae = newAPIError(nil, msg, codeOpt...)
	return
}

// NewAPIErrorWrap 创建包装过原始错误信息的 APIError ，code 默认：API_CRASH
func NewAPIErrorWrap(err error) (ae error) {
	msg := fmt.Sprintf("系统内部错误 [%d]", time.Now().Unix())
	ae = newAPIError(err, msg, constant.API_CRASH)
	return
}

// WrapAPIError 包装原始错误为 APIError，code 默认：API_ERROR
// func WrapAPIError(err error, msg string, codeOpt ...constant.APICodeType) (ae error) {
// 	ae = newAPIError(err, msg, codeOpt...)
// 	return
// }

// 直接调用 errors.Unwrap 即可检出原始 err
// UnwrapAPIError 从 APIError 错误中检出原始 err
// func UnwrapAPIError(ae error) error {
// 	for {
// 		e := errors.Unwrap(ae)
// 		if e == nil {
// 			return nil
// 		}

// 		if ae, ok := e.(*APIError); ok {
// 			return ae
// 		}
// 	}
// }

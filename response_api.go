package gins

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zaaksam/gins/constant"
	"github.com/zaaksam/gins/errors"
	"github.com/zaaksam/gins/extend/logger"
)

type apiResponse struct {
	ctx *gin.Context

	data any
	kv   map[string]any

	result struct {
		Code constant.APICodeType `json:"code"`
		Msg  string               `json:"msg"`
		Data any                  `json:"data"`
	}
}

// GetAPIResponse 获取 API 响应对象
func GetAPIResponse(ctx *gin.Context) (res *apiResponse) {
	val, exists := ctx.Get(gins_context_response_api)
	if exists {
		res = val.(*apiResponse)
		return
	}

	res = &apiResponse{
		ctx: ctx,
	}

	ctx.Set(gins_context_response_api, res)

	return res
}

// SetData 设置响应数据，code 默认：API_OK
func (res *apiResponse) SetData(data any) {
	res.data = data

	res.SetOK()
}

// SetKV 设置响应键值对，code 默认：API_OK
func (res *apiResponse) SetKV(key string, val any) {
	if res.kv == nil {
		res.kv = make(map[string]any)
	}

	res.kv[key] = val

	res.SetOK()
}

// SetOK 设置无 data 数据的正常响应，code 默认：API_OK
func (res *apiResponse) SetOK() {
	res.SetCode(constant.API_OK)
}

// SetCode 设置 API Code
func (res *apiResponse) SetCode(code constant.APICodeType) {
	res.result.Code = code
}

// SetMsg 设置 msg 信息，code 默认：API_ERROR
func (res *apiResponse) SetMsg(msg string, codeOpt ...constant.APICodeType) {
	if len(codeOpt) == 1 {
		res.result.Code = codeOpt[0]
	} else {
		res.result.Code = constant.API_ERROR
	}

	res.result.Msg = msg
}

// SetError 设置错误，code 默认：API_ERROR
func (res *apiResponse) SetError(err error) {
	ae, ok := err.(*errors.APIError)
	if ok {
		res.result.Code = ae.Code
		res.result.Msg = ae.Msg

		wrapErr := ae.Unwrap()
		if wrapErr != nil {
			logger.Errorf("%s", ae.UnwrapError())
		}
		return
	}

	res.result.Code = constant.API_ERROR
	res.result.Msg = err.Error()
}

func (res *apiResponse) render() {
	var obj any
	if res.data != nil {
		obj = res.data
	} else if res.kv != nil {
		obj = res.kv
	}

	if obj == nil {
		obj = struct{}{}
	}
	res.result.Data = obj

	res.ctx.JSON(http.StatusOK, res.result)
}

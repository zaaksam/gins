package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/zaaksam/gins/errors"
)

const (
	request_body = "request_body"
)

func getRequestBody(ctx *gin.Context) (body []byte, err error) {
	val, exists := ctx.Get(request_body)
	if exists {
		body, exists = val.([]byte)
		if !exists {
			err = errors.NewAPIError("Request Body读取缓存错误")
			return
		}
	} else {
		body, err = ctx.GetRawData()
		if err != nil {
			err = errors.NewAPIErrorWrap(err)
			return
		}
	}

	ctx.Set(request_body, body)
	return
}

func parseModel(ctx *gin.Context, md any) (err error) {
	body, err := getRequestBody(ctx)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, md)
	return
}

func parseJSONResult(ctx *gin.Context) (jsonResult *gjson.Result, err error) {
	body, err := getRequestBody(ctx)
	if err != nil {
		return
	}

	jsonResult = new(gjson.Result)
	*jsonResult = gjson.ParseBytes(body)
	if !jsonResult.IsObject() {
		err = errors.NewAPIError("json参数格式错误")
	}
	return
}

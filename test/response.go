package test

import (
	"github.com/goccy/go-json"

	"github.com/zaaksam/gins/constant"
)

type Response struct {
	Code constant.APICodeType `json:"code"`
	Msg  string               `json:"msg"`
	Data any                  `json:"data"`

	content string
}

func NewResponse(body []byte) (res *Response, err error) {
	res = &Response{
		content: string(body),
	}
	err = json.Unmarshal(body, res)
	return
}

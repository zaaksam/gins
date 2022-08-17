package gins

import (
	"github.com/gin-gonic/gin"
)

const (
	gins_context_response_api = "gins_context_response_api"
	gins_context_response_web = "gins_context_response_web"
)

func handler(ctx *gin.Context) {
	ctx.Next()

	var res IResponse

	// API 响应
	val, exists := ctx.Get(gins_context_response_api)
	if exists {
		res = val.(*apiResponse)
	}

	if res == nil {
		// Web 响应
		val, exists = ctx.Get(gins_context_response_web)
		if exists {
			res = val.(*webResponse)
		}
	}

	if res != nil {
		res.render()
	}
}

package gins

import (
	"github.com/gin-gonic/gin"
)

func onHandler(gs *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// 提前中止
		if ctx.IsAborted() {
			return
		}

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
}

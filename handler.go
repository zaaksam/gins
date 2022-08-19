package gins

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func onHandler(gs *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		httpStatus := ctx.Writer.Status()

		if httpStatus == http.StatusNotFound {
			if gs.conf.On404 != nil {
				gs.conf.On404(ctx)
			}

			return
		}

		// TODO: 是否要处理 http.StatusInternalServerError 500 事件调用 On505

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

package gins

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zaaksam/gins/errors"
	"github.com/zaaksam/gins/extend/logger"
)

func onRecover(gs *Server) gin.HandlerFunc {
	return gin.CustomRecovery(func(ctx *gin.Context, err any) {
		if gs.conf.On500 != nil {
			gs.conf.On500(ctx, err)
			return
		}

		// API 响应
		val, exists := ctx.Get(gins_context_response_api)
		if exists {
			apiRes := val.(*apiResponse)

			apiErr := errors.NewAPIErrorWrap(fmt.Errorf("%s", err))
			apiRes.SetError(apiErr)
			return
		}

		// Web 响应
		// val, exists = ctx.Get(gins_context_response_web)
		// if exists {
		// 	// webRes := val.(*webResponse)
		// 	return
		// }

		ae := errors.NewAPIErrorWrap(fmt.Errorf("%s", err))
		logger.Warnf("%s：%s", ae.Msg, ae.Unwrap().Error())

		ctx.String(http.StatusInternalServerError, ae.Msg)
	})
}

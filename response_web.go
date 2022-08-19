package gins

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type webResponse struct {
	ctx *gin.Context
}

// TODO: web 处理包装

func (res *webResponse) render() {
	res.ctx.HTML(http.StatusOK, "template_name", nil)
}

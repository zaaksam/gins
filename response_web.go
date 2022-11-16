package gins

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type webResponse struct {
	ctx  *gin.Context
	name string
	data gin.H
}

// GetWebResponse 获取 Web 响应对象
func GetWebResponse(ctx *gin.Context) (res *webResponse) {
	val, exists := ctx.Get(gins_context_response_web)
	if exists {
		res = val.(*webResponse)
		return
	}

	res = &webResponse{
		ctx: ctx,
	}

	ctx.Set(gins_context_response_web, res)

	return res
}

// 设置 永久 跳转信息
func (res *webResponse) Set301Redirect(location string) {
	res.ctx.Redirect(http.StatusMovedPermanently, location)

	// 设置了跳转，中止后续处理
	res.ctx.Abort()
}

// 设置 临时 跳转信息
func (res *webResponse) Set302Redirect(location string) {
	res.ctx.Redirect(http.StatusFound, location)

	// 设置了跳转，中止后续处理
	res.ctx.Abort()
}

// 设置 header 信息
func (res *webResponse) SetHeader(key, val string) {
	res.ctx.Header(key, val)
}

// 设置 cookie 信息
func (res *webResponse) SetCookie(key, val string) {
	// maxAge: -1，在浏览器打开期间有效，关闭失效，单位：秒，0 表示马上过期
	// path: /
	// domain: 跟随请求 host
	// secure: true，只能在 https 下传输
	// httpOnly: true，无法被 js 读取

	domain := res.ctx.Request.URL.Hostname()
	res.ctx.SetCookie(key, val, -1, "/", domain, true, true)
}

// SetTemplate 设置渲染模板名称及数据
func (res *webResponse) SetTemplate(name string, data gin.H) {
	res.name = name
	res.data = data
}

func (res *webResponse) render() {
	res.ctx.HTML(http.StatusOK, res.name, res.data)
}

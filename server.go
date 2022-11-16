package gins

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zaaksam/gins/extend/logger"
)

type Server struct {
	conf       *Config
	engine     *gin.Engine
	httpServer *http.Server
	rootCtx    context.Context
	rootCancel context.CancelFunc

	templ *template.Template
}

// New 新服务
func New(conf *Config) (gs *Server, err error) {
	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	gs = &Server{
		conf:   conf,
		engine: gin.New(),
	}

	gs.httpServer = &http.Server{
		Addr:    gs.conf.Addr(),
		Handler: gs.engine,
	}

	gs.rootCtx, gs.rootCancel = context.WithCancel(context.Background())

	// 添加 panic 恢复处理
	gs.engine.Use(onRecover(gs))

	// 添加日志处理
	if conf.Debug {
		gs.engine.Use(gin.Logger())
	}

	// 使用自定义 404 处理
	if gs.conf.On404 != nil {
		gs.engine.NoRoute(gs.conf.On404)
	}

	// 使用自定义响应处理
	gs.engine.Use(onHandler(gs))

	return
}

// Engine gin engine
func (gs *Server) Engine() *gin.Engine {
	return gs.engine
}

// Run 启动
func (gs *Server) Run() {
	// 设置模板
	if gs.templ != nil {
		gs.Engine().SetHTMLTemplate(gs.templ)
	}

	logger.Infof("[%s %s]服务运行在：%s", gs.conf.Name, gs.conf.Version, gs.conf.BroadcastAddr())

	if !gs.conf.IsEnableKeepAlives {
		// FIXME: 不关闭的话，优雅退出时，会导致有连接挂起，总是需要超时退出
		// 目前关闭 keep-alive 状态并未成功
		gs.httpServer.SetKeepAlivesEnabled(false)
	}

	err := gs.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Errorf("服务运行时出错：%s", err)
	}
}

// Stop 停止
func (gs *Server) Stop() {
	logger.Infof("正在停止[%s %s]服务...", gs.conf.Name, gs.conf.Version)

	stopCtx, stopCancel := context.WithTimeout(gs.rootCtx, time.Duration(gs.conf.Timeout)*time.Second)

	err := gs.httpServer.Shutdown(stopCtx)
	if err != nil {
		logger.Warnf("[%s %s]服务停止出错：%s", gs.conf.Name, gs.conf.Version, err)
	} else {
		logger.Infof("[%s %s]服务已停止", gs.conf.Name, gs.conf.Version)
	}

	gs.rootCancel()
	stopCancel()

	// FIXME: 延时2秒退出，让超时任务 504 响应完成
	time.Sleep(2 * time.Second)
}

// AddTemplate 添加模板
func (gs *Server) AddTemplate(name string, body []byte) (err error) {
	if gs.templ == nil {
		gs.templ, err = template.New(name).Parse(string(body))
	} else {
		gs.templ, err = gs.templ.New(name).Parse(string(body))
	}
	return
}

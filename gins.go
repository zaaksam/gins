package gins

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/zaaksam/gins/extend/logger"
)

var (
	signalChan chan os.Signal
	instance   *Server

	initFuncs []func(gs *Server)
)

// Init 初始化实例，一般直接调用 Run 即可，内置了 Init 的调用
func Init(conf *Config) {
	if instance != nil {
		return
	}

	// 设置日志级别
	logger.SetLevel(conf.logrusLevel)

	var err error
	instance, err = New(conf)
	if err != nil {
		panic(err)
	}

	// 初始化事件
	for i, l := 0, len(initFuncs); i < l; i++ {
		initFuncs[i](instance)
	}

}

// Instance 实例
func Instance() *Server {
	return instance
}

// Run 启动
func Run(conf *Config) {
	signalChan = make(chan os.Signal, 1)

	// 初始化实例
	Init(conf)

	if !instance.conf.IsDisableSignal {
		//使用docker stop 命令去关闭Container时，该命令会发送SIGTERM 命令到Container主进程，让主进程处理该信号，关闭Container，如果在10s内，未关闭容器，Docker Damon会发送SIGKILL 信号将Container关闭。
		signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)
	}

	go instance.Run()

	<-signalChan

	if !instance.conf.IsDisableSignal {
		signal.Stop(signalChan)
	}

	instance.Stop()

	if instance.conf.IsDisableSignal {
		signalChan <- syscall.SIGINT
	}
}

// Stop 停止
func Stop() {
	if signalChan == nil {
		return
	}

	signalChan <- syscall.SIGINT

	if instance.conf.IsDisableSignal {
		<-signalChan
	}

	signalChan = nil
}

// AddInit 安全添加初始化方法
func AddInit(fn func(gs *Server)) {
	if len(initFuncs) == 0 {
		initFuncs = make([]func(*Server), 0, 50)
	}

	initFuncs = append(initFuncs, fn)
}

// GetValidate 获取 gin 内部的 Validate 对象
func GetValidate() *validator.Validate {
	return binding.Validator.Engine().(*validator.Validate)
}

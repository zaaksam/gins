package gins

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Config 服务器配置
type Config struct {
	Name          string // 服务名称，必填
	Version       string // 服务版本，必填
	Host          string // 域名主机
	IP            string // 运行地址，必填
	BroadcastIP   string // 广播的运行地址，默认为：IP
	Port          int    // 运行端口，必填
	BroadcastPort int    // 广播的运行端口，默认为：Port
	Timeout       int    // 优雅退出时的超时机制，默认：30秒
	Debug         bool   // 是否开启调试
	Pprof         bool   // 是否监控性能

	LogLevel        string       // logger 日志级别，支持：trace、debug、info、warn、error、fatal、panic，默认：info
	logrusLevel     logrus.Level // LogLevel 的 logrus 格式化
	IsDisableSignal bool         // 是否关闭 signal 信号监听退出，默认：false，设置为 true 时，需主动调用 gins.Stop() 来触发优雅退出

	On500 gin.RecoveryFunc // 自定义 500 处理，相当于 panic 的 recover 处理
	On404 gin.HandlerFunc  // 自定义 404 处理
}

// Addr 运行地址
func (conf *Config) Addr() string {
	return fmt.Sprintf("%s:%d", conf.IP, conf.Port)
}

// BroadcastAddr 广播的运行地址
func (conf *Config) BroadcastAddr() string {
	return fmt.Sprintf("%s:%d", conf.BroadcastIP, conf.BroadcastPort)
}

// Init 初始化配置
func (conf *Config) Init() (err error) {
	if conf.Name == "" {
		err = errors.New("配置 Name 不能为空")
	} else if conf.Version == "" {
		err = errors.New("配置 Version 不能为空")
	} else if conf.IP == "" {
		err = errors.New("配置 IP 不能为空")
	} else if conf.Port <= 0 {
		err = errors.New("配置 Port 不能为空")
	}
	if err != nil {
		return
	}

	if conf.BroadcastIP == "" {
		conf.BroadcastIP = conf.IP
	}

	if conf.BroadcastPort <= 0 {
		conf.BroadcastPort = conf.Port
	}

	if conf.Timeout <= 0 {
		conf.Timeout = 30
	}

	if conf.LogLevel == "" {
		conf.LogLevel = "info"
	}

	conf.logrusLevel, err = logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		err = fmt.Errorf("配置 LogLevel 错误：%s", err)
	}
	return
}

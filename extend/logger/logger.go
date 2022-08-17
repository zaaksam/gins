package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	isInit bool
	logger *logrus.Logger
)

// Init 初始化，默认不需要调用，除非需要指定配置
func Init(level logrus.Level, formatterOpt ...logrus.Formatter) {
	if isInit {
		return
	}
	isInit = true

	var formatter logrus.Formatter

	if len(formatterOpt) == 1 {
		formatter = formatterOpt[0]
	} else {
		textFormatter := &logrus.TextFormatter{}
		textFormatter.TimestampFormat = "2006-01-02 15:04:05"
		textFormatter.FullTimestamp = true

		formatter = textFormatter
	}

	logger = &logrus.Logger{
		Out:       os.Stderr,
		Formatter: formatter,
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}
}

// 默认初始化
func defaultInit() {
	if isInit {
		return
	}

	Init(logrus.InfoLevel)
}

// Instance Logrus实例
func Instance() (instance *logrus.Logger) {
	defaultInit()

	return logger
}

// SetLevel 设置日志级别
func SetLevel(level logrus.Level) {
	defaultInit()

	logger.SetLevel(level)
}

// Panicf 日志
func Panicf(format string, args ...any) {
	defaultInit()

	logger.Panicf(format, args...)
}

// Panic 日志
func Panic(args ...any) {
	defaultInit()

	logger.Panic(args...)
}

// Fatalf 日志
func Fatalf(format string, args ...any) {
	defaultInit()

	logger.Fatalf(format, args...)
}

// Fatal 日志
func Fatal(args ...any) {
	defaultInit()

	logger.Fatal(args...)
}

// Errorf 日志
func Errorf(format string, args ...any) {
	defaultInit()

	logger.Errorf(format, args...)
}

// Error 日志
func Error(args ...any) {
	defaultInit()

	logger.Error(args...)
}

// Warnf 日志
func Warnf(format string, args ...any) {
	defaultInit()

	logger.Warnf(format, args...)
}

// Warn 日志
func Warn(args ...any) {
	defaultInit()

	logger.Warn(args...)
}

// Infof 日志
func Infof(format string, args ...any) {
	defaultInit()

	logger.Infof(format, args...)
}

// Info 日志
func Info(args ...any) {
	defaultInit()

	logger.Info(args...)
}

// Debugf 日志
func Debugf(format string, args ...any) {
	defaultInit()

	logger.Debugf(format, args...)
}

// Debug 日志
func Debug(args ...any) {
	defaultInit()

	logger.Debug(args...)
}

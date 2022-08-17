package orm

import (
	"github.com/zaaksam/gins/extend/logger"
	xormLog "xorm.io/xorm/log"
)

// logger xorm Logger 接口实现
type myLogger struct {
	isShowSQL bool
	level     xormLog.LogLevel
}

func (l *myLogger) Level() xormLog.LogLevel {
	return l.level
}

func (l *myLogger) SetLevel(level xormLog.LogLevel) {
	l.level = level
}

func (l *myLogger) ShowSQL(show ...bool) {
	if len(show) == 1 {
		l.isShowSQL = show[0]
	}
}

func (l *myLogger) IsShowSQL() bool {
	return l.isShowSQL
}

func (l *myLogger) Debug(v ...interface{}) {
	logger.Debug(v...)
}

func (l *myLogger) Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func (l *myLogger) Info(v ...interface{}) {
	logger.Info(v...)
}

func (l *myLogger) Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func (l *myLogger) Warn(v ...interface{}) {
	logger.Warn(v...)
}

func (l *myLogger) Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func (l *myLogger) Error(v ...interface{}) {
	logger.Error(v...)
}

func (l *myLogger) Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

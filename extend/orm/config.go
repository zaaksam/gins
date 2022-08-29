package orm

import (
	"errors"
	"fmt"
	"strconv"

	"xorm.io/xorm/log"
)

// Config 数据库配置
type Config struct {
	Alias        string       `json:"alias"`
	Type         string       `json:"type"`
	Server       string       `json:"server"`
	Port         int          `json:"port"`
	Database     string       `json:"database"`
	User         string       `json:"user"`
	Password     string       `json:"password"`
	IsShowSQL    bool         `json:"isShowSQL"` // 是否打印 SQL 语句
	LogLevel     string       `json:"logLevel"`  // logger 日志级别，支持：debug、info、warn、error、off，默认：off
	xormLogLevel log.LogLevel // LogLevel 的 xorm log 格式化
}

// Init 初始化
func (conf *Config) Init() (err error) {
	if conf.LogLevel == "" {
		conf.LogLevel = "off"
	}

	switch conf.LogLevel {
	case "debug":
		conf.xormLogLevel = log.LOG_DEBUG
	case "info":
		conf.xormLogLevel = log.LOG_INFO
	case "warn":
		conf.xormLogLevel = log.LOG_WARNING
	case "error":
		conf.xormLogLevel = log.LOG_ERR
	case "off":
		conf.xormLogLevel = log.LOG_OFF
	default:
		err = errors.New("数据库 LogLevel 配置错误：" + conf.LogLevel)
	}

	return
}

// getConnect 获取连接信息
func (conf *Config) getConnect() (driverName, dataSourceName string, err error) {
	switch conf.Type {
	// case "mssql":
	// 	driverName = "mssql"
	// 	dataSourceName += "server=" + conf.Server
	// 	dataSourceName += ";port=" + strconv.Itoa(conf.Port)
	// 	dataSourceName += ";database=" + conf.Database
	// 	dataSourceName += ";user id=" + conf.User
	// 	dataSourceName += ";password=" + conf.Password
	//	case "odbc":
	//		//lunny mssql driver
	//		driverName = "odbc"
	//		dataSourceName = "driver={SQL Server}"
	//		dataSourceName += ";Server=" + db["server"] + "," + db["port"]
	//		dataSourceName += ";Database=" + db["database"]
	//		dataSourceName += ";uid=" + db["user"] + ";pwd=" + db["password"] + ";"
	case "mysql":
		driverName = "mysql"
		dataSourceName = conf.User + ":" + conf.Password
		dataSourceName += "@(" + conf.Server + ":" + strconv.Itoa(conf.Port) + ")/"
		dataSourceName += conf.Database + "?charset=utf8mb4&loc=Asia%2FShanghai&multiStatements=true"
	case "sqlite3":
		driverName = "sqlite3"
		dataSourceName = conf.Database
	default:
		err = fmt.Errorf("不支持数据库 %s 的类型：%s", conf.Alias, conf.Type)
	}
	return
}

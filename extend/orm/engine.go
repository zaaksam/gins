package orm

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/zaaksam/gins/extend/logger"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
	"xorm.io/xorm/schemas"

	_ "github.com/go-sql-driver/mysql"
	// _ "github.com/yingshengtech/go-mssqldb"
	//_ "github.com/mattn/go-sqlite3"
)

// Engine 数据库引擎管理对象
type Engine struct {
	confs       map[string]*Config
	xormEngines map[string]*xorm.Engine
}

// newEngine 创建
func newEngine(confs []*Config) (engine *Engine, err error) {
	l := len(confs)
	if l == 0 {
		err = errors.New("数据库 Config 列表不能为空")
		return
	}

	engine = &Engine{
		confs: make(map[string]*Config),
	}
	for i := 0; i < l; i++ {
		err = confs[i].Init()
		if err != nil {
			return
		}

		engine.confs[confs[i].Alias] = confs[i]
	}

	engine.xormEngines = make(map[string]*xorm.Engine)
	for alias := range dbRegistry {
		// 检查数据库连接是否禁用，不创建
		_, ok := dbDisable[alias]
		if ok {
			logger.Infof("禁用 %s 数据库引擎", alias)
			continue
		}

		logger.Infof("初始化 %s 数据库引擎", alias)

		// 检查注册的数据库存在配置信息
		conf, ok := engine.confs[alias]
		if !ok {
			err = fmt.Errorf("数据库 %s 不存在 Config 配置", alias)
			return
		}

		// 获取数据库连接配置
		var driverName, dataSourceName string
		driverName, dataSourceName, err = conf.getConnect()
		if err != nil {
			return
		}

		// 创建 xorm 引擎
		var xormEngine *xorm.Engine
		xormEngine, err = xorm.NewEngine(driverName, dataSourceName)
		if err != nil {
			err = fmt.Errorf("数据库 %s 创建引擎时出错：%s", alias, err)
			return
		}

		// 测试连通
		err = xormEngine.Ping()
		if err != nil {
			err = fmt.Errorf("数据库 %s 测试连接时出错：%s", alias, err)
			return
		}

		// 连接池的空闲数大小
		// xormEngine.SetMaxIdleConns(5)
		// 最大打开连接数
		// xormEngine.SetMaxConns(500)

		// 结构体命名与数据库一致
		xormEngine.SetMapper(names.NewCacheMapper(new(names.SameMapper)))

		// 设置日志接口
		myLog := &myLogger{
			isShowSQL: conf.IsShowSQL,
			level:     conf.xormLogLevel,
		}
		xormEngine.SetLogger(myLog)
		xormEngine.ShowSQL(myLog.isShowSQL)
		xormEngine.SetLogLevel(myLog.level)

		// 设置默认超时
		setTimeOut(xormEngine)

		engine.xormEngines[alias] = xormEngine
	}

	return
}

// func (engine *Engine) getDBType(alias string) string {
// 	conf, ok := engine.confs[alias]
// 	if ok {
// 		return conf.Type
// 	}

// 	return ""
// }

func (engine *Engine) getXormEngine(alias string) (xormEngine *xorm.Engine, err error) {
	xormEngine, ok := engine.xormEngines[alias]
	if ok {
		return
	}

	err = fmt.Errorf("数据库引擎：%s 不存在", alias)
	return
}

// setTimeOut 设置连接超时时间
func setTimeOut(xormEngine *xorm.Engine) {
	if xormEngine.Dialect().URI().DBType != schemas.MYSQL {
		logger.Warnf("不支持为除 %s 之外的其他类型的数据库设置连接超时时间，请自行设置", schemas.MYSQL)
		return
	}
	rs, err := xormEngine.Query("show  variables like '%wait_timeout%'")
	if err != nil {
		logger.Errorf("设置数据库连接超时失败，无法获取数据库变量，原因是【%v】", err)
		return
	}
	//单位秒
	var timeOut uint64
	for _, v := range rs {
		if string(v["Variable_name"]) == "wait_timeout" {
			timeOutStr := string(v["Value"])
			timeOut, err = strconv.ParseUint(timeOutStr, 10, 64)
			if err != nil {
				logger.Errorf("设置数据库连接超时失败，wait_timeout【%v】环境变量不是数字格式", timeOutStr)
				return
			}
			//取数据库超时时间的前2s作为客户端连接超时时间
			var offset uint64 = 2
			if timeOut-offset == 0 {
				timeOut = 0
				logger.Warn("数据库wait_timeout时间低于%v秒，设置客户端超时时间为0，即永不超时 \n", offset)
			} else {
				timeOut = timeOut - offset
			}
			xormEngine.SetConnMaxLifetime(time.Duration(timeOut) * time.Second)
			return
		}
	}
}

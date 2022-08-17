package db

import (
	"github.com/zaaksam/gins/example/config"
	"github.com/zaaksam/gins/extend/orm"
)

// Init 初始化
func Init() (err error) {
	l := len(config.Instance.DBs)
	confs := make([]*orm.Config, 0, l)

	for i := 0; i < l; i++ {
		conf := &orm.Config{
			Alias:     config.Instance.DBs[i].Alias,
			Type:      config.Instance.DBs[i].Type,
			Server:    config.Instance.DBs[i].Server,
			Port:      config.Instance.DBs[i].Port,
			Database:  config.Instance.DBs[i].Database,
			User:      config.Instance.DBs[i].User,
			Password:  config.Instance.DBs[i].Password,
			IsShowSQL: config.Instance.DBs[i].IsShowSQL,
			LogLevel:  config.Instance.DBs[i].LogLevel,
		}

		confs = append(confs, conf)
	}

	err = orm.Init(confs)
	return
}

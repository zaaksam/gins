package test

import (
	"github.com/zaaksam/gins"
	"github.com/zaaksam/gins/example/config"
	"github.com/zaaksam/gins/example/db"
	"github.com/zaaksam/gins/extend/orm"
	"github.com/zaaksam/gins/test"
)

func init() {
	config.Instance.Gins.Name = "gins_example"
	config.Instance.Gins.IP = "localhost"
	config.Instance.Gins.Port = 8080
	config.Instance.Gins.Debug = true
	config.Instance.LogLevel = "trace"

	err := config.Instance.Init("0.0.0")
	if err != nil {
		panic(err)
	}

	err = db.Init()
	if err != nil {
		panic(err)
	}

	gins.AddInit(func(*gins.Server) {
		// 注册 orm 到 validator
		validate := gins.GetValidate()
		orm.RegisterWithValidator(validate)
	})

	test.Init(config.Instance.Gins)
}

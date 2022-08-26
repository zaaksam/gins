package main

import (
	"github.com/zaaksam/gins"
	"github.com/zaaksam/gins/example/config"
	_ "github.com/zaaksam/gins/example/controller"
	"github.com/zaaksam/gins/example/db"
	"github.com/zaaksam/gins/extend/orm"
)

func main() {
	err := config.Instance.Init(VERSION)
	if err != nil {
		panic(err)
	}

	err = db.Init()
	if err != nil {
		panic(err)
	}

	gins.AddInit(func(*gins.Server) {
		// 注册 orm 自定义字段规则到 validator
		validate := gins.GetValidate()
		orm.RegisterWithValidator(validate)
	})

	// 注意：在 gins.Run 之前，日志默认打印级别是：InfoLevel
	gins.Run(config.Instance.Gins)
}

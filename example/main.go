package main

import (
	"github.com/zaaksam/gins"
	"github.com/zaaksam/gins/example/config"
	_ "github.com/zaaksam/gins/example/controller"
	"github.com/zaaksam/gins/example/db"
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

	// 注意：在 gins.Run 之前，日志默认打印级别是：InfoLevel
	gins.Run(config.Instance.Gins)
}

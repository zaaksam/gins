package test

import (
	"github.com/zaaksam/gins"
)

var (
	instance *gins.Server
	conf     *gins.Config
	isInit   bool
)

// Init 初始化
func Init(c *gins.Config) {
	conf = c
}

func onceInit() {
	if isInit {
		return
	}
	isInit = true

	gins.Init(conf)
	instance = gins.Instance()
}

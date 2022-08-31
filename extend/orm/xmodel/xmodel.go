package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/zaaksam/gins/extend/logger"
	"github.com/zaaksam/gins/extend/orm/xmodel/command"
)

func main() {
	logger.SetLevel(logrus.DebugLevel)

	genCmd := command.Gen()

	app := &cli.App{
		Name:  "xmodel",
		Usage: "orm model 扩展代码生成工具",
		Commands: []*cli.Command{
			genCmd,
			command.Version(),
		},
	}

	// 设置默认命令为：gen
	app.Flags = genCmd.Flags
	app.Action = genCmd.Action

	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err)
		cli.Exit(err, 1)
	}
}

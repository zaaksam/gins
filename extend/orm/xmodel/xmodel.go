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

	app := &cli.App{
		Name:  "xmodel",
		Usage: "orm model 扩展代码生成工具",
		Commands: []*cli.Command{
			command.Gen(),
			command.Version(),
		},
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Error(err)
		cli.Exit(err, 1)
	}
}

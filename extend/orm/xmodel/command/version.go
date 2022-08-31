package command

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Version() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Usage:   "版本号",
		Aliases: []string{"v", "ver"},
		Action: func(ctx *cli.Context) error {
			fmt.Println("v1.0.1")
			return nil
		},
	}
}

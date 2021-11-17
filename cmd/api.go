// Package cmd manage all commands required to launch cypress-parallel-cli
package cmd

import (
	"github.com/urfave/cli/v2"
)

// API command options
func API(c *cli.Context) (z *cli.Command) {
	return &cli.Command{
		Name:  "api",
		Usage: "Start api server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-file",
				Aliases:     []string{"c"},
				Usage:       "Config file name",
				Required:    true,
				Destination: &cmdValidate.ConfigDir,
			},
		},
		Action: func(c *cli.Context) error {
			cmdAPI.Run()
			return nil
		},
	}
}

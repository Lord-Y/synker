// Package cmd manage all commands required to launch cypress-parallel-cli
package cmd

import (
	"github.com/Lord-Y/synker/api"
	"github.com/urfave/cli/v2"
)

// API command options
func API(c *cli.Context) (z *cli.Command) {
	return &cli.Command{
		Name:  "api",
		Usage: "Start api server",
		Action: func(c *cli.Context) error {
			api.Run()
			return nil
		},
	}
}

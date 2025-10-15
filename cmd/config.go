package cmd

import (
	"context"

	"github.com/Lord-Y/synker/logger"
	"github.com/urfave/cli/v3"
)

// Validate commands
func Validate() (z *cli.Command) {
	return &cli.Command{
		Name:  "validate",
		Usage: "options related to configuration file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-dir",
				Aliases:     []string{"c"},
				Usage:       "Config dir name holding files",
				Required:    true,
				Destination: &cmdValidate.ConfigDir,
			},
		},
		Action: func(context.Context, *cli.Command) error {
			cmdValidate.Logger = logger.NewLogger()
			return cmdValidate.ParseAndValidateConfig()
		},
	}
}

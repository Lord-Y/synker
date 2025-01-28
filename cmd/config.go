package cmd

import (
	"github.com/Lord-Y/synker/logger"
	"github.com/urfave/cli/v2"
)

// Validate commands
func Validate(c *cli.Context) (z *cli.Command) {
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
		Action: func(c *cli.Context) error {
			cmdValidate.Logger = logger.NewLogger()
			cmdValidate.ParseAndValidateConfig()
			return nil
		},
	}
}

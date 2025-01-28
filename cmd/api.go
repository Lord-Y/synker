// Package cmd manage all commands required to launch cypress-parallel-cli
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Lord-Y/synker/logger"

	"github.com/urfave/cli/v2"
)

// API command options
func API(c *cli.Context) (z *cli.Command) {
	return &cli.Command{
		Name:  "api",
		Usage: "Start api server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-dir",
				Aliases:     []string{"c"},
				Usage:       "Config dir name holding files",
				Required:    true,
				Destination: &cmdValidate.ConfigDir,
			},
			&cli.BoolFlag{
				Name:        "init",
				Aliases:     []string{"i"},
				Usage:       "Perform prerequisites related to elasticsearch / kafka / cockroachdb",
				Required:    false,
				Destination: &cmdValidate.Init,
			},
		},
		Action: func(c *cli.Context) error {
			cmdValidate.Logger = logger.NewLogger()

			if strings.TrimSpace(os.Getenv("SYNKER_PG_URI")) == "" {
				msg := "SYNKER_PG_URI environment variable must be set"
				cmdValidate.Logger.Fatal().Err(fmt.Errorf("%s", msg)).Msg(msg)
			}
			if strings.TrimSpace(os.Getenv("SYNKER_ELASTICSEARCH_URI")) == "" {
				msg := "SYNKER_ELASTICSEARCH_URI environment variable must be set"
				cmdValidate.Logger.Fatal().Err(fmt.Errorf("%s", msg)).Msg(msg)
			}
			if strings.TrimSpace(os.Getenv("SYNKER_KAFKA_URI")) == "" {
				msg := "SYNKER_KAFKA_URI environment variable must be set"
				cmdValidate.Logger.Fatal().Err(fmt.Errorf("%s", msg)).Msg(msg)
			}

			cmdValidate.ParseAndValidateConfig()
			cmdValidate.RunAPI()
			return nil
		},
	}
}

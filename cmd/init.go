// Package cmd manage all commands required to launch cypress-parallel-cli
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Lord-Y/synker/logger"
	"github.com/urfave/cli/v2"
)

// Init command options
func Init(c *cli.Context) (z *cli.Command) {
	return &cli.Command{
		Name:  "init",
		Usage: "Perform prerequisites related to elasticsearch / kafka / cockroachdb",
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
			cmdValidate.RunPrerequisitesOnly()
			return nil
		},
	}
}

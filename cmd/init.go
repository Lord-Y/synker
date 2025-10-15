// Package cmd manage all commands required to launch cypress-parallel-cli
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Lord-Y/synker/logger"
	"github.com/urfave/cli/v3"
)

// Init command options
func Init() (z *cli.Command) {
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
		Action: func(context.Context, *cli.Command) error {
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

			if err := cmdValidate.ParseAndValidateConfig(); err != nil {
				return err
			}
			cmdValidate.RunPrerequisitesOnly()
			return nil
		},
	}
}

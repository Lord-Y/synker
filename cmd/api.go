// Package cmd manage all commands required to launch cypress-parallel-cli
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

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
		},
		Action: func(c *cli.Context) error {
			if strings.TrimSpace(os.Getenv("SYNKER_PG_URI")) == "" {
				msg := "SYNKER_PG_URI environment variable must be set"
				log.Fatal().Err(fmt.Errorf(msg)).Msg(msg)
			}
			if strings.TrimSpace(os.Getenv("SYNKER_ELASTICSEARCH_URI")) == "" {
				msg := "SYNKER_ELASTICSEARCH_URI environment variable must be set"
				log.Fatal().Err(fmt.Errorf(msg)).Msg(msg)
			}
			if strings.TrimSpace(os.Getenv("SYNKER_KAFKA_URI")) == "" {
				msg := "SYNKER_KAFKA_URI environment variable must be set"
				log.Fatal().Err(fmt.Errorf(msg)).Msg(msg)
			}

			cmdValidate.Run()
			cmdAPI.Run(&cmdValidate)
			return nil
		},
	}
}

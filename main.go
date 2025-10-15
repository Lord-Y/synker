package main

import (
	"context"
	"os"
	"strings"

	"github.com/Lord-Y/synker/cmd"
	"github.com/Lord-Y/synker/logger"
	"github.com/urfave/cli/v3"
)

// Make vars public for unit testing
var (
	VersionDetails *cli.Command
	CmdValidate    *cli.Command
	CmdAPI         *cli.Command
	CmdInit        *cli.Command
)

func init() {
	CmdValidate = cmd.Validate()
	VersionDetails = cmd.VersionDetails()
	CmdAPI = cmd.API()
	CmdInit = cmd.Init()
}

func main() {
	cmd := cli.Command{
		Name:                  "synker",
		Usage:                 "A new cli with an api to sync data between databases and elasticsearch",
		Description:           "synker permit you to validate config and sync data from databases to elasticsearch",
		Version:               cmd.Version,
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			VersionDetails,
			CmdValidate,
			CmdInit,
			CmdAPI,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		if !strings.Contains(err.Error(), "Required flag") {
			logger.NewLogger().Fatal().Err(err).Msg("Error occured while executing the program")
		}
	}
}

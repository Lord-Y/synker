package main

import (
	"os"
	"strings"

	"github.com/Lord-Y/synker/cmd"
	"github.com/Lord-Y/synker/logger"
	cli "github.com/urfave/cli/v2"
)

// Make vars public for unit testing
var (
	VersionDetails *cli.Command
	CmdValidate    *cli.Command
	CmdAPI         *cli.Command
	CmdInit        *cli.Command
)

func init() {
	CmdValidate = cmd.Validate(&cli.Context{})
	VersionDetails = cmd.VersionDetails(&cli.Context{})
	CmdAPI = cmd.API(&cli.Context{})
	CmdInit = cmd.Init(&cli.Context{})
}

func main() {
	app := cli.NewApp()
	app.Name = "synker"
	app.Usage = "A new cli with an api to sync data between databases and elasticsearch"
	app.Description = "synker permit you to validate config and sync data from databases to elasticsearch"
	app.Version = cmd.Version
	app.EnableBashCompletion = true

	app.Commands = []*cli.Command{
		VersionDetails,
		CmdValidate,
		CmdInit,
		CmdAPI,
	}

	if err := app.Run(os.Args); err != nil {
		if !strings.Contains(err.Error(), "Required flag") {
			logger.NewLogger().Error().Err(err).Msg("Error occured while executing the program")
		}
		os.Exit(1)
	}
}

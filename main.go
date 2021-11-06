package main

import (
	"os"

	"github.com/Lord-Y/synker/cmd"
	"github.com/Lord-Y/synker/logger"
	"github.com/rs/zerolog/log"
	cli "github.com/urfave/cli/v2"
)

// Make vars public for unit testing
var (
	VersionDetails *cli.Command
	CmdValidate    *cli.Command
	CmdAPI         *cli.Command
)

func init() {
	logger.SetCLILoggerLogLevel()

	CmdValidate = cmd.Validate(&cli.Context{})
	VersionDetails = cmd.VersionDetails(&cli.Context{})
	CmdAPI = cmd.API(&cli.Context{})
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
		CmdAPI,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Error occured while executing the program")
	}
}

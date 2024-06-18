package main

import (
	"os"
	"strings"

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
	os.Setenv("SYNKER_BATCH_LOG", "true")
	defer os.Unsetenv("APP_BATCH_LOG")
	logger.SetLoggerLogLevel()

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

	if err := app.Run(os.Args); err != nil {
		if !strings.Contains(err.Error(), "Required flag") {
			log.Error().Err(err).Msg("Error occured while executing the program")
		}
		os.Exit(1)
	}
}

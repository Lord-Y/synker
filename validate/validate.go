package validate

import (
	"github.com/Lord-Y/synker/logger"
	"github.com/Lord-Y/synker/models"
	"github.com/rs/zerolog/log"
)

func init() {
	logger.SetCLILoggerLogLevel()
}

type Validate models.Configuration

// Run will run the validate command
func (c *Validate) Run() {
	log.Info().Msgf("Config file is %s", c.ConfigFile)
}

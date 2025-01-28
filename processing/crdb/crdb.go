// Package etcd assemble all requirements to manipulate data received from CockroachDB
package crdb

import (
	"io"

	"github.com/Lord-Y/synker/logger"
	"github.com/gin-gonic/gin"
)

func Trigger(c *gin.Context) {
	var body []byte
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.NewLogger().Error().Err(err).Msg("blooooo")
	}
	logger.NewLogger().Info().Msgf("resp method %s content %+v", c.Request.Method, string(body))
}

// Package etcd assemble all requirements to manipulate data received from CockroachDB
package crdb

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Trigger(c *gin.Context) {
	var body []byte
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("blooooo")
	}
	log.Info().Msgf("resp method %s content %+v", c.Request.Method, string(body))
}

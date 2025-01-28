// Package routers assemble all routes for the api
package routers

import (
	"fmt"
	"os"
	"strings"

	cLogger "github.com/Lord-Y/synker/logger"
	"github.com/Lord-Y/synker/processing/crdb"
	"github.com/Lord-Y/synker/processing/health"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SetupRouter func handle all routes of the api
func SetupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(requestid.New())
	router.Use(gin.Recovery())

	router.Use(
		logger.SetLogger(
			logger.WithUTC(true),
			logger.WithLogger(
				func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
					log.Logger = *cLogger.NewLogger()
					return log.With().Str("requestId", requestid.Get(c)).Logger()
				},
			),
		),
	)

	// disable during unit testing
	if strings.TrimSpace(os.Getenv("SYNKER_PROMETHEUS")) != "" {
		prometheus_port := "9101"
		if strings.TrimSpace(os.Getenv("SYNKER_PROMETHEUS_PORT")) != "" {
			prometheus_port = strings.TrimSpace(os.Getenv("SYNKER_PROMETHEUS_PORT"))
		}
		p := ginprometheus.NewPrometheus("http")
		p.SetListenAddressWithRouter(fmt.Sprintf(":%s", prometheus_port), router)
		p.Use(router)
	}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", health.Health)
		v1.GET("/healthz", health.Healthz)

		v1.POST("/crdb/trigger", crdb.Trigger)
		v1.PUT("/crdb/trigger", crdb.Trigger)
	}
	return router
}

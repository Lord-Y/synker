// Package routers assemble all routes for the api
package routers

import (
	"fmt"
	"os"
	"strings"

	"github.com/Lord-Y/synker/api/crdb"
	"github.com/Lord-Y/synker/api/health"
	apiLogger "github.com/Lord-Y/synker/logger"
	"github.com/Lord-Y/synker/tools"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	os.Setenv("SYNKER_BATCH_LOG", "true")
	defer os.Unsetenv("APP_BATCH_LOG")
	apiLogger.SetLoggerLogLevel()
}

// SetupRouter func handle all routes of the api
func SetupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	requestID := tools.RandStringInt(32)

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	router := gin.New()
	router.Use(requestid.New())
	router.Use(gin.Recovery())

	router.Use(
		logger.SetLogger(
			logger.WithUTC(true),
			logger.WithLogger(
				func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
					return zerolog.New(os.Stdout).
						With().
						Timestamp().
						Str("requestId", requestid.Get(c)).
						Logger()
				},
			),
		),
	)
	headerHandler := func(c *gin.Context) {
		if c.GetHeader("X-Request-Id") == "" {
			c.Request.Header.Set("X-Request-Id", requestID)
			c.Next()
		}
	}
	router.Use(headerHandler)
	// disable during unit testing
	if strings.TrimSpace(os.Getenv("SYNKER_PROMETHEUS")) != "" {
		var SYNKER_prometheus_port string = "9101"
		if strings.TrimSpace(os.Getenv("SYNKER_PROMETHEUS_PORT")) != "" {
			SYNKER_prometheus_port = strings.TrimSpace(os.Getenv("SYNKER_PROMETHEUS_PORT"))
		}
		p := ginprometheus.NewPrometheus("http")
		p.SetListenAddress(fmt.Sprintf(":%s", SYNKER_prometheus_port))
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

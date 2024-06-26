// Package api assemble all requirements to start the api
package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Lord-Y/synker/api/routers"
	apiLogger "github.com/Lord-Y/synker/logger"
	"github.com/Lord-Y/synker/models"
	"github.com/Lord-Y/synker/processing"
	"github.com/rs/zerolog/log"
)

type API models.Configuration

// init func
func init() {
	apiLogger.SetAPILoggerLogLevel()

	if strings.TrimSpace(os.Getenv("SYNKER_PG_URI")) == "" {
		msg := "SYNKER_PG_URI environment variable must be set"
		log.Fatal().Err(fmt.Errorf(msg)).Msg(msg)
		return
	}
	if strings.TrimSpace(os.Getenv("SYNKER_ELASTICSEARCH_URI")) == "" {
		msg := "SYNKER_ELASTICSEARCH_URI environment variable must be set"
		log.Fatal().Err(fmt.Errorf(msg)).Msg(msg)
		return
	}
	if strings.TrimSpace(os.Getenv("SYNKER_KAFKA_URI")) == "" {
		msg := "SYNKER_KAFKA_URI environment variable must be set"
		log.Fatal().Err(fmt.Errorf(msg)).Msg(msg)
		return
	}
}

// Run will start the api server
func (c *API) Run(validated *processing.Validate) {
	var srv *http.Server
	router := routers.SetupRouter()

	os.Setenv("SYNKER_CONFIG_DIR", c.ConfigDir)
	defer os.Unsetenv("SYNKER_CONFIG_DIR")

	err := validated.ManageTopics()
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to manage topics")
		return
	}
	err = validated.ManageElasticsearchIndex()
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to manage elasticsearch indexes")
		return
	}
	err = validated.ManageChangeFeed()
	if err != nil {
		log.Fatal().Err(err).Msg("Fail to manage changefeed")
		return
	}

	appPort := strings.TrimSpace(os.Getenv("SYNKER_API_PORT"))
	if appPort != "" {
		srv = &http.Server{
			Addr:    fmt.Sprintf(":%s", appPort),
			Handler: router,
		}
		log.Info().Msgf("Starting api server on port %s", appPort)
	} else {
		srv = &http.Server{
			Addr:    ":8080",
			Handler: router,
		}
		log.Info().Msg("Starting api server on port 8080")
	}

	go func() {
		validated.Processing()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Startup api server failed")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down api server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("API server shutted down abruptly")
	}
	log.Info().Msg("API server exited successfully")
}

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
	os.Setenv("SYNKER_BATCH_LOG", "true")
	defer os.Unsetenv("APP_BATCH_LOG")
	apiLogger.SetLoggerLogLevel()
}

// Run will start the api server
func (c *API) Run(validated *processing.Validate) {
	var appPort string
	port := strings.TrimSpace(os.Getenv("SYNKER_API_PORT"))

	switch strings.HasPrefix(port, ":") {
	case true:
		appPort = port
	case false:
		if port == "" {
			appPort = ":8080"
		} else {
			appPort = fmt.Sprintf(":%s", port)
		}
	}

	router := routers.SetupRouter()
	srv := &http.Server{
		Addr:    appPort,
		Handler: router,
	}
	log.Info().Msgf("Starting api server on port %s", appPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Startup api server failed")
		}
	}()

	go c.runPrerequisitesAndStartProcessing(validated)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 minutes.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down api server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("API server shutted down abruptly")
	}
	log.Info().Msg("API server exited successfully")
}

// runPrerequisitesAndStartProcessing permit to run all functions
// related to Kafka, elasticsearch and cockroach feeds
// and then start processing all kafka messages
func (c *API) runPrerequisitesAndStartProcessing(validated *processing.Validate) {
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

	validated.Processing()
}

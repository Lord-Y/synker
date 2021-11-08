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
	"github.com/rs/zerolog/log"
)

type API models.Configuration

// init func
func init() {
	apiLogger.SetAPILoggerLogLevel()

	if strings.TrimSpace(os.Getenv("SKR_PG_URI")) == "" {
		msg := "SKR_PG_URI environment variable must be set"
		log.Fatal().Err(fmt.Errorf(msg)).Msg(msg)
		return
	}
	if strings.TrimSpace(os.Getenv("SKR_ELASTICSEARCH_URI")) == "" {
		msg := "SKR_ELASTICSEARCH_URI environment variable must be set"
		log.Fatal().Err(fmt.Errorf(msg)).Msg(msg)
		return
	}
}

// Run will start the api server
func (c *API) Run() {
	var srv *http.Server
	router := routers.SetupRouter()

	os.Setenv("SKR_CONFIG_FILE", c.ConfigFile)
	defer os.Unsetenv("SKR_CONFIG_FILE")
	appPort := strings.TrimSpace(os.Getenv("SKR_API_PORT"))
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

// Package processing provide all requirements to process change data capture
package processing

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Lord-Y/synker/processing/routers"
)

// RunAPI will start the api server
func (c *Validate) RunAPI() {
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
	c.Logger.Info().Msgf("Starting api server on port %s", appPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.Logger.Fatal().Err(err).Msg("Startup api server failed")
		}
	}()

	go c.runPrerequisitesAndStartProcessing()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 minutes.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.Logger.Info().Msg("Shutting down api server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		c.Logger.Fatal().Err(err).Msg("API server shutted down abruptly")
	}
	c.Logger.Info().Msg("API server exited successfully")
}

// RunPrerequisitesOnly permit to run all functions
// related to Kafka, elasticsearch and cockroach feeds
func (c *Validate) RunPrerequisitesOnly() {
	os.Setenv("SYNKER_CONFIG_DIR", c.ConfigDir)
	defer os.Unsetenv("SYNKER_CONFIG_DIR")

	err := c.manageTopics()
	if err != nil {
		c.Logger.Fatal().Err(err).Msg("Fail to manage topics")
		return
	}
	err = c.manageElasticsearchIndex()
	if err != nil {
		c.Logger.Fatal().Err(err).Msg("Fail to manage elasticsearch indexes")
		return
	}
	err = c.manageChangeFeed()
	if err != nil {
		c.Logger.Fatal().Err(err).Msg("Fail to manage changefeed")
		return
	}
}

// runPrerequisitesAndStartProcessing permit to run all functions
// related to Kafka, elasticsearch and cockroach feeds
// and then start processing all kafka messages
func (c *Validate) runPrerequisitesAndStartProcessing() {
	os.Setenv("SYNKER_CONFIG_DIR", c.ConfigDir)
	defer os.Unsetenv("SYNKER_CONFIG_DIR")

	if c.Init {
		c.RunPrerequisitesOnly()
	}
	c.processing()
}

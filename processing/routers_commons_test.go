// Package processing provide all requirements to process change data capture
package processing

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/Lord-Y/synker/logger"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	var c Validate
	c.Logger = logger.NewLogger()
	router := c.setupRouter()
	w, err := performRequest(router, headers, "GET", "/api/v1/health", "")
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to perform http GET request")
		assert.FailNow("Failed to perform http GET request")
		return
	}

	assert.Equal(200, w.Code, "Failed to perform http GET request")
	assert.Contains(w.Body.String(), `{"health":"OK"}`, "Failed to get right body content")
}

func TestHealth_no_xrequestid(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	var c Validate
	router := c.setupRouter()
	w, err := performRequest(router, headers, "GET", "/api/v1/health", "")
	if err != nil {
		assert.FailNow("Failed to perform http GET request")
		return
	}

	assert.Equal(200, w.Code, "Failed to perform http GET request")
	assert.Contains(w.Body.String(), `{"health":"OK"}`, "Failed to get right body content")
}

func TestHealthz(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	var c Validate
	router := c.setupRouter()
	w, err := performRequest(router, headers, "GET", "/api/v1/healthz", "")
	if err != nil {
		assert.FailNow("Failed to perform http GET request")
		return
	}

	assert.Equal(200, w.Code, "Failed to perform http GET request")
	assert.Contains(w.Body.String(), `{"health":"OK"}`, "Failed to get right body content")
}

func TestHealth_prometheus(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	os.Setenv("SYNKER_PROMETHEUS", "1")
	defer os.Unsetenv("SYNKER_PROMETHEUS")
	var c Validate
	c.Logger = logger.NewLogger()
	router := c.setupRouter()
	w, err := performRequest(router, headers, "GET", "/api/v1/health", "")
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to perform http GET request")
		assert.FailNow("Failed to perform http GET request")
		return
	}

	assert.Equal(200, w.Code, "Failed to perform http GET request")
	assert.Contains(w.Body.String(), `{"health":"OK"}`, "Failed to get right body content")
}

func TestHealth_prometheus_port(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	os.Setenv("SYNKER_PROMETHEUS", "1")
	os.Setenv("SYNKER_PROMETHEUS_PORT", "9111")
	defer os.Unsetenv("SYNKER_PROMETHEUS")
	defer os.Unsetenv("SYNKER_PROMETHEUS_PORT")
	var c Validate
	c.Logger = logger.NewLogger()
	router := c.setupRouter()
	w, err := performRequest(router, headers, "GET", "/api/v1/health", "")
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to perform http GET request")
		assert.FailNow("Failed to perform http GET request")
		return
	}

	assert.Equal(200, w.Code, "Failed to perform http GET request")
	assert.Contains(w.Body.String(), `{"health":"OK"}`, "Failed to get right body content")

	resp, err := http.Get("http://localhost:9111/metrics")
	if err != nil {
		c.Logger.Error().Err(err).Msg("Failed to perform http GET metrics endpoint")
		assert.FailNow("Failed to perform http metrics GET endpoint")
	}
	body, errb := io.ReadAll(resp.Body)
	if errb != nil {
		assert.FailNow("Failed get body")
	}
	defer resp.Body.Close()
	assert.Equal(200, resp.StatusCode)
	assert.NotNil(body)
}

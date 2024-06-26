package routers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	assert := assert.New(t)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	router := SetupRouter()
	w, err := performRequest(router, headers, "GET", "/api/v1/health", "")
	if err != nil {
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

	router := SetupRouter()
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

	router := SetupRouter()
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
	router := SetupRouter()
	w, err := performRequest(router, headers, "GET", "/api/v1/health", "")
	if err != nil {
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
	os.Setenv("SYNKER_PROMETHEUS_PORT", "9101")
	defer os.Unsetenv("SYNKER_PROMETHEUS")
	defer os.Unsetenv("SYNKER_PROMETHEUS_PORT")
	router := SetupRouter()
	w, err := performRequest(router, headers, "GET", "/api/v1/health", "")
	if err != nil {
		assert.FailNow("Failed to perform http GET request")
		return
	}

	assert.Equal(200, w.Code, "Failed to perform http GET request")
	assert.Contains(w.Body.String(), `{"health":"OK"}`, "Failed to get right body content")
}

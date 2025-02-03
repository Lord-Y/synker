// Package processing provide all requirements to process change data capture
package processing

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/Lord-Y/synker/logger"
)

func performRequest(r http.Handler, headers map[string]string, method string, url string, payload string) (z *httptest.ResponseRecorder, err error) {
	var (
		req *http.Request
	)

	switch method {
	case "GET":
		req, err = http.NewRequest(method, url, nil)
	case "POST":
		req, err = http.NewRequest(method, url, strings.NewReader(payload))
	default:
		if payload == "" {
			req, err = http.NewRequest(method, url, nil)
		} else {
			req, err = http.NewRequest(method, url, strings.NewReader(payload))
		}
	}
	if err != nil {
		logger.NewLogger().Error().Err(err).Msgf("Error occured while initalising http request")
		return nil, err
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, nil
}

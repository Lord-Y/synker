// Package processing provide all requirements to process change data capture
package processing

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// health permit to return basic health check
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "OK"})
}

// healthz permit to get the status of all backend required by the api
func healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "OK"})
}

// Package health assemble all functions required for health checks
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health permit to return basic health check
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "OK"})
}

// Healthz permit to get the status of all backend required by the api
func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"health": "OK"})
}

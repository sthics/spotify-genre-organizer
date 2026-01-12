package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "spotify-genre-organizer",
	})
}

func NotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}

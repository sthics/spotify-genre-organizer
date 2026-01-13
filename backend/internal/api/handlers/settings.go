package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/database"
)

func GetSettings(c *gin.Context) {
	// userID, _ := c.Cookie("user_id") // Should use authenticated user from middleware if available
	// For now, getting from cookie is fine as per existing pattern
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	settings, err := database.GetUserSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

type UpdateSettingsRequest struct {
	NameTemplate        string `json:"name_template"`
	DescriptionTemplate string `json:"description_template"`
}

func UpdateSettings(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate templates
	if !strings.Contains(req.NameTemplate, "{genre}") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name template must contain {genre}"})
		return
	}

	settings, err := database.GetUserSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch settings"})
		return
	}

	settings.NameTemplate = req.NameTemplate
	settings.DescriptionTemplate = req.DescriptionTemplate

	if err := database.SaveUserSettings(settings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

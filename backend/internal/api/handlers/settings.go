package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/models"
)

// In-memory storage for MVP (replace with database later)
var userSettingsStore = make(map[string]*models.UserSettings)

func GetSettings(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	settings, ok := userSettingsStore[userID]
	if !ok {
		// Return defaults if not found
		settings = models.DefaultSettings(userID)
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

	// Validate templates contain {genre}
	if !strings.Contains(req.NameTemplate, "{genre}") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name_template must contain {genre}"})
		return
	}

	settings := &models.UserSettings{
		UserID:              userID,
		NameTemplate:        req.NameTemplate,
		DescriptionTemplate: req.DescriptionTemplate,
	}

	// Check if user was previously premium
	if existing, ok := userSettingsStore[userID]; ok {
		settings.IsPremium = existing.IsPremium
	}

	userSettingsStore[userID] = settings

	c.JSON(http.StatusOK, settings)
}

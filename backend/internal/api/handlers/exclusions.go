package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/models"
)

type CreateExclusionRequest struct {
	Type      string  `json:"type" binding:"required,oneof=artist song"`
	SpotifyID string  `json:"spotify_id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Scope     *string `json:"scope" binding:"omitempty,oneof=all_appearances primary_only"`
}

// ListExclusions returns all exclusion rules for the user
func ListExclusions(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	rules, err := database.GetExclusionRules(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch exclusions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

// CreateExclusion creates a new exclusion rule and returns impact preview
func CreateExclusion(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req CreateExclusionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build the rule
	rule := &models.ExclusionRule{
		UserID:        userID,
		ExclusionType: models.ExclusionType(req.Type),
		SpotifyID:     req.SpotifyID,
		Name:          req.Name,
	}

	if req.Scope != nil {
		scope := models.ExclusionScope(*req.Scope)
		rule.Scope = &scope
	}

	// Save to database
	if err := database.CreateExclusionRule(rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create exclusion"})
		return
	}

	// Calculate impact preview (placeholder - will be implemented in Task 8)
	preview, err := calculateExclusionImpact(userID, accessToken, rule)
	if err != nil {
		// Rule created but preview failed - still return success
		c.JSON(http.StatusCreated, gin.H{
			"rule_id": rule.ID,
			"preview": nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"rule_id": rule.ID,
		"preview": preview,
	})
}

// Placeholder for Task 8
func calculateExclusionImpact(userID, accessToken string, rule *models.ExclusionRule) (interface{}, error) {
	return nil, nil
}

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type LibraryCountResponse struct {
	Count    int       `json:"count"`
	CachedAt time.Time `json:"cached_at"`
}

// Simple in-memory cache (per-user)
var countCache = make(map[string]*LibraryCountResponse)

func GetLibraryCount(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userID, _ := c.Cookie("user_id")

	// Check cache (5 minute TTL)
	if cached, ok := countCache[userID]; ok {
		if time.Since(cached.CachedAt) < 5*time.Minute {
			c.JSON(http.StatusOK, cached)
			return
		}
	}

	// Fetch count from Spotify
	count, err := spotify.GetLikedSongsCount(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch count"})
		return
	}

	response := &LibraryCountResponse{
		Count:    count,
		CachedAt: time.Now(),
	}
	countCache[userID] = response

	c.JSON(http.StatusOK, response)
}

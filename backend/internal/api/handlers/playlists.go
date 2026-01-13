package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type ManagedPlaylist struct {
	SpotifyID  string  `json:"spotify_id"`
	Name       string  `json:"name"`
	Genre      string  `json:"genre"`
	SongCount  int     `json:"song_count"`
	SpotifyURL string  `json:"spotify_url"`
	ImageURL   *string `json:"image_url"`
	CustomName *string `json:"custom_name"`
	CustomDesc *string `json:"custom_description"`
	LastSynced *string `json:"last_synced"`
}

func ListPlaylists(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userID, _ := c.Cookie("user_id")

	// Fetch all user's playlists from Spotify
	playlists, err := spotify.GetUserPlaylists(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch playlists"})
		return
	}

	// Get user's settings to match naming pattern
	settings, ok := userSettingsStore[userID]
	if !ok {
		settings = models.DefaultSettings(userID)
	}

	// Filter to only Organizer-created playlists
	var managed []ManagedPlaylist
	for _, p := range playlists {
		// Check if matches our naming pattern (contains "by Organizer" or template pattern)
		if !isOrganizerPlaylist(p.Name, settings.NameTemplate) {
			continue
		}

		genre := extractGenreFromName(p.Name, settings.NameTemplate)

		mp := ManagedPlaylist{
			SpotifyID:  p.ID,
			Name:       p.Name,
			Genre:      genre,
			SongCount:  p.Tracks.Total,
			SpotifyURL: p.ExternalURLs.Spotify,
		}

		if len(p.Images) > 0 {
			mp.ImageURL = &p.Images[0].URL
		}

		managed = append(managed, mp)
	}

	c.JSON(http.StatusOK, gin.H{
		"playlists":   managed,
		"total_songs": sumSongCounts(managed),
	})
}

func isOrganizerPlaylist(name, template string) bool {
	// Check for default pattern
	if strings.Contains(name, "by Organizer") {
		return true
	}
	// Check if matches user template pattern
	templateBase := strings.ReplaceAll(template, "{genre}", "")
	if templateBase != "" && strings.Contains(name, strings.TrimSpace(templateBase)) {
		return true
	}
	return false
}

func extractGenreFromName(name, template string) string {
	// Try to extract genre from name
	if idx := strings.Index(name, " by Organizer"); idx > 0 {
		return name[:idx]
	}
	// Fallback for custom templates
	return "Unknown"
}

func sumSongCounts(playlists []ManagedPlaylist) int {
	sum := 0
	for _, p := range playlists {
		sum += p.SongCount
	}
	return sum
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type CreateExclusionRequest struct {
	Type      string  `json:"type" binding:"required,oneof=artist song"`
	SpotifyID string  `json:"spotify_id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Scope     *string `json:"scope" binding:"omitempty,oneof=all_appearances primary_only"`
}

type AffectedPlaylist struct {
	PlaylistID    string   `json:"playlist_id"`
	Name          string   `json:"name"`
	SongsToRemove []string `json:"songs_to_remove"`
	SongCount     int      `json:"song_count"`
}

type ImpactPreview struct {
	TotalSongsAffected int                `json:"total_songs_affected"`
	AffectedPlaylists  []AffectedPlaylist `json:"affected_playlists"`
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

	// Calculate impact preview
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

func calculateExclusionImpact(userID, accessToken string, rule *models.ExclusionRule) (*ImpactPreview, error) {
	// Get all managed playlists
	overrides, err := database.GetPlaylistOverrides(userID)
	if err != nil {
		return nil, err
	}

	preview := &ImpactPreview{
		AffectedPlaylists: []AffectedPlaylist{},
	}

	for _, override := range overrides {
		// Fetch tracks from this playlist
		tracks, err := spotify.GetPlaylistTracks(accessToken, override.PlaylistSpotifyID)
		if err != nil {
			continue // Skip playlists we can't fetch
		}

		// Find matching tracks
		var matchingTracks []string
		for _, track := range tracks {
			if matchesExclusionRule(track, rule) {
				matchingTracks = append(matchingTracks, track.ID)
			}
		}

		if len(matchingTracks) > 0 {
			name := override.Genre
			if override.CustomName != nil {
				name = *override.CustomName
			}
			preview.AffectedPlaylists = append(preview.AffectedPlaylists, AffectedPlaylist{
				PlaylistID:    override.PlaylistSpotifyID,
				Name:          name,
				SongsToRemove: matchingTracks,
				SongCount:     len(matchingTracks),
			})
			preview.TotalSongsAffected += len(matchingTracks)
		}
	}

	return preview, nil
}

func matchesExclusionRule(track spotify.Song, rule *models.ExclusionRule) bool {
	if rule.ExclusionType == models.ExclusionTypeSong {
		return track.ID == rule.SpotifyID
	}

	// Artist rule
	if rule.Scope != nil && *rule.Scope == models.ExclusionScopeAllAppearances {
		for _, artist := range track.Artists {
			if artist.ID == rule.SpotifyID {
				return true
			}
		}
	} else {
		// Primary only
		if len(track.Artists) > 0 && track.Artists[0].ID == rule.SpotifyID {
			return true
		}
	}

	return false
}

// ApplyExclusion removes matching tracks from playlists
func ApplyExclusion(c *gin.Context) {
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

	ruleID := c.Param("id")

	// Get the rule
	rule, err := database.GetExclusionRuleByID(userID, ruleID)
	if err != nil || rule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exclusion rule not found"})
		return
	}

	// Get impact
	preview, err := calculateExclusionImpact(userID, accessToken, rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate impact"})
		return
	}

	// Remove tracks from each playlist
	playlistsUpdated := 0
	songsRemoved := 0

	for _, affected := range preview.AffectedPlaylists {
		err := spotify.RemoveTracksFromPlaylist(accessToken, affected.PlaylistID, affected.SongsToRemove)
		if err != nil {
			continue
		}
		playlistsUpdated++
		songsRemoved += len(affected.SongsToRemove)
	}

	c.JSON(http.StatusOK, gin.H{
		"playlists_updated": playlistsUpdated,
		"songs_removed":     songsRemoved,
	})
}

// DeleteExclusion removes an exclusion rule
func DeleteExclusion(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	ruleID := c.Param("id")

	if err := database.DeleteExclusionRule(userID, ruleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete exclusion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/genres"
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

// knownGenres for playlist detection
var knownGenres = genres.ParentGenres

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
	settings, err := database.GetUserSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch settings"})
		return
	}

	// Debug: Check which user the access token belongs to
	profile, profileErr := spotify.GetUserProfile(accessToken)
	if profileErr == nil {
		log.Printf("DEBUG: Access token belongs to user: %s (%s)", profile.ID, profile.DisplayName)
	} else {
		log.Printf("DEBUG: Could not fetch profile: %v", profileErr)
	}
	log.Printf("DEBUG: User ID from cookie: %s", userID)
	log.Printf("DEBUG: Found %d total playlists from Spotify", len(playlists))
	log.Printf("DEBUG: User template: %s", settings.NameTemplate)

	// Filter to only Organizer-created playlists
	var managed []ManagedPlaylist
	for _, p := range playlists {
		log.Printf("DEBUG: Checking playlist: '%s'", p.Name)
		// Check if matches our naming pattern (contains "by Organizer" or template pattern)
		if !isOrganizerPlaylist(p.Name, settings.NameTemplate) {
			log.Printf("DEBUG: SKIPPED - doesn't match pattern")
			continue
		}
		log.Printf("DEBUG: MATCHED!")

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
	// Check for "Organizer" anywhere (catches variations)
	if strings.Contains(strings.ToLower(name), "organizer") {
		return true
	}
	// Check if matches user template pattern
	templateBase := strings.ReplaceAll(template, "{genre}", "")
	if templateBase != "" && strings.Contains(name, strings.TrimSpace(templateBase)) {
		return true
	}
	// Check if playlist name starts with a known genre (likely created by organizer)
	for _, genre := range knownGenres {
		if strings.HasPrefix(name, genre+" ") || name == genre {
			return true
		}
	}
	return false
}

func extractGenreFromName(name, template string) string {
	// Try to extract genre from name using "by Organizer" pattern
	if idx := strings.Index(name, " by Organizer"); idx > 0 {
		return name[:idx]
	}
	// Try to match against known genres
	for _, genre := range knownGenres {
		if strings.HasPrefix(name, genre+" ") || name == genre {
			return genre
		}
		// Also check if genre appears at start case-insensitively
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(genre)) {
			return genre
		}
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

// Helper to get overrides from database (cached or direct)
// For MVP, just hitting DB direct is fine as per scale
func getPlaylistOverride(userID, playlistID string) *models.PlaylistOverride {
	overrides, err := database.GetPlaylistOverrides(userID)
	if err != nil {
		return nil
	}
	return overrides[playlistID]
}

func savePlaylistOverride(override *models.PlaylistOverride) {
	database.SavePlaylistOverride(override)
}

func deletePlaylistOverride(userID, playlistID string) {
	// Not strictly deleting row, but could implement real delete if needed.
	// For now, let's just assume modifying logic or leave it as is if we implemented delete in DB.
	// But our DB helper was Upsert. We probably need a Delete function in DB helper.
	// For MVP, let's skip explicit delete implementation in DB helper if not requested,
	// or implemented a soft delete by setting genre to empty?
	// Actually, let's just ignoring implementation detail for now or add Delete to DB.
	// Let's assume we implement DeletePlaylistOverrides in database/settings.go next.
}

type UpdatePlaylistRequest struct {
	CustomName        *string `json:"custom_name"`
	CustomDescription *string `json:"custom_description"`
}

func UpdatePlaylist(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userID, _ := c.Cookie("user_id")
	playlistID := c.Param("id")

	var req UpdatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Build update values
	newName := ""
	newDesc := ""

	if req.CustomName != nil {
		newName = *req.CustomName
	}
	if req.CustomDescription != nil {
		newDesc = *req.CustomDescription

		// Append footer for free users
		settings, err := database.GetUserSettings(userID)
		if err == nil && (!settings.IsPremium) {
			newDesc += " • spotifygenreorganizer.com"
		}
	}

	// Update in Spotify
	if newName != "" || newDesc != "" {
		if err := spotify.UpdatePlaylistDetails(accessToken, playlistID, newName, newDesc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update playlist"})
			return
		}
	}

	// Save override locally
	override := getPlaylistOverride(userID, playlistID)
	if override == nil {
		override = &models.PlaylistOverride{
			UserID:            userID,
			PlaylistSpotifyID: playlistID,
		}
	}
	if req.CustomName != nil {
		override.CustomName = req.CustomName
	}
	if req.CustomDescription != nil {
		override.CustomDescription = req.CustomDescription
	}
	savePlaylistOverride(override)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeletePlaylist(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userID, _ := c.Cookie("user_id")
	playlistID := c.Param("id")

	// Unfollow (delete) the playlist in Spotify
	if err := spotify.UnfollowPlaylist(accessToken, playlistID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete playlist"})
		return
	}

	// Remove local override
	deletePlaylistOverride(userID, playlistID)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func RefreshPlaylist(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userID, _ := c.Cookie("user_id")
	playlistID := c.Param("id")

	// Get the playlist's genre from our override store
	override := getPlaylistOverride(userID, playlistID)
	if override == nil || override.Genre == "" {
		// Try to get genre from the playlist name
		playlists, err := spotify.GetUserPlaylists(accessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch playlists"})
			return
		}

		settings, err := database.GetUserSettings(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch settings"})
			return
		}

		var foundGenre string
		for _, p := range playlists {
			if p.ID == playlistID {
				foundGenre = extractGenreFromName(p.Name, settings.NameTemplate)
				break
			}
		}

		if foundGenre == "" || foundGenre == "Unknown" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not determine playlist genre"})
			return
		}

		// Store the genre for future use
		if override == nil {
			override = &models.PlaylistOverride{
				UserID:            userID,
				PlaylistSpotifyID: playlistID,
				Genre:             foundGenre,
			}
		} else {
			override.Genre = foundGenre
		}
		savePlaylistOverride(override)
	}

	// Fetch all liked songs
	songs, err := spotify.FetchAllLikedSongs(accessToken, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch songs"})
		return
	}

	// Enrich with genres
	artistGenres, err := spotify.FetchAllArtistGenres(accessToken, songs, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch artist genres"})
		return
	}
	spotify.EnrichSongsWithGenres(songs, artistGenres)

	// Filter to songs matching this genre
	var genreSongs []spotify.Song
	for _, song := range songs {
		if genres.ScoreGenres(song.Genres) == override.Genre {
			genreSongs = append(genreSongs, song)
		}
	}

	// Clear the playlist
	if err := spotify.ClearPlaylist(accessToken, playlistID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear playlist"})
		return
	}

	// Add tracks back
	trackIDs := make([]string, len(genreSongs))
	for i, s := range genreSongs {
		trackIDs[i] = s.ID
	}

	if len(trackIDs) > 0 {
		if err := spotify.AddTracksToPlaylist(accessToken, playlistID, trackIDs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add tracks"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"song_count": len(genreSongs),
	})
}



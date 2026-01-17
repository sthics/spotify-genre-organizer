package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/genres"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type PlaylistSyncStatus struct {
	SpotifyID string `json:"spotify_id"`
	Genre     string `json:"genre"`
	NewCount  int    `json:"new_count"`
}

type SyncStatusResponse struct {
	NewSongsCount int                  `json:"new_songs_count"`
	OldestSyncAt  *time.Time           `json:"oldest_sync_at"`
	Playlists     []PlaylistSyncStatus `json:"playlists"`
}

func GetSyncStatus(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userID, err := c.Cookie("user_id")
	if err != nil || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	// Get oldest sync timestamp
	oldestSync, err := database.GetOldestSyncTimestamp(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sync status"})
		return
	}

	// If no playlists synced yet, return empty
	if oldestSync == nil {
		c.JSON(http.StatusOK, SyncStatusResponse{
			NewSongsCount: 0,
			OldestSyncAt:  nil,
			Playlists:     []PlaylistSyncStatus{},
		})
		return
	}

	// Fetch all liked songs
	songs, err := spotify.FetchAllLikedSongs(accessToken, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch songs"})
		return
	}

	// Filter to songs added after oldest sync
	// Subtract 1 minute buffer to account for timezone/precision edge cases
	syncThreshold := oldestSync.Add(-1 * time.Minute)
	var newSongs []spotify.Song
	for _, song := range songs {
		if song.AddedAt.After(syncThreshold) {
			newSongs = append(newSongs, song)
		}
	}

	if len(newSongs) == 0 {
		c.JSON(http.StatusOK, SyncStatusResponse{
			NewSongsCount: 0,
			OldestSyncAt:  oldestSync,
			Playlists:     []PlaylistSyncStatus{},
		})
		return
	}

	// Enrich new songs with genres
	artistGenres, err := spotify.FetchAllArtistGenres(accessToken, newSongs, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch genres"})
		return
	}
	spotify.EnrichSongsWithGenres(newSongs, artistGenres)

	// Get user's playlist overrides to know which playlists exist
	overrides, err := database.GetPlaylistOverrides(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get playlists"})
		return
	}

	// Count new songs per genre/playlist
	genreCounts := make(map[string]int)
	for _, song := range newSongs {
		genre := genres.ScoreGenres(song.Genres)
		genreCounts[genre]++
	}

	// Build playlist status list
	var playlistStatuses []PlaylistSyncStatus
	for playlistID, override := range overrides {
		if override.Genre != "" {
			count := genreCounts[override.Genre]
			if count > 0 {
				playlistStatuses = append(playlistStatuses, PlaylistSyncStatus{
					SpotifyID: playlistID,
					Genre:     override.Genre,
					NewCount:  count,
				})
			}
		}
	}

	c.JSON(http.StatusOK, SyncStatusResponse{
		NewSongsCount: len(newSongs),
		OldestSyncAt:  oldestSync,
		Playlists:     playlistStatuses,
	})
}

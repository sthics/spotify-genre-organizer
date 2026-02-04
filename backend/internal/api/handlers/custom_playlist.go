package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/organizer"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type CustomPlaylistRequest struct {
	SubGenres []string `json:"sub_genres" binding:"required,min=1"`
	Mode      string   `json:"mode" binding:"required,oneof=combined separate"`
	Name      string   `json:"name"` // Only used for combined mode
}

func StartCustomPlaylist(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req CustomPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get cached genre analysis (must exist for this endpoint)
	cache, err := database.GetGenreAnalysisCache(userID)
	if err != nil || cache == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please analyze your library first"})
		return
	}

	// Create job
	jobID := uuid.New().String()
	job := &JobStatus{
		ID:     jobID,
		Status: "pending",
		Stage:  "initializing",
	}

	jobsMu.Lock()
	jobs[jobID] = job
	jobsMu.Unlock()

	// Start async processing
	go processCustomPlaylistJob(job, accessToken, userID, req, cache)

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": jobID,
		"status": "pending",
	})
}

func processCustomPlaylistJob(job *JobStatus, accessToken, userID string, req CustomPlaylistRequest, cache *models.GenreAnalysisCache) {
	updateJob := func() {
		jobsMu.Lock()
		jobs[job.ID] = job
		jobsMu.Unlock()
	}

	job.Status = "processing"
	job.Stage = "filtering"
	updateJob()

	// Build set of selected sub-genres for quick lookup
	selectedGenres := make(map[string]bool)
	for _, sg := range req.SubGenres {
		selectedGenres[sg] = true
	}

	// Filter songs that match selected sub-genres
	var matchingTrackIDs []string
	for trackID, trackGenres := range cache.AnalysisData.SongGenreMap {
		for _, g := range trackGenres {
			if selectedGenres[g] {
				matchingTrackIDs = append(matchingTrackIDs, trackID)
				break
			}
		}
	}

	// Apply exclusion rules to track IDs
	exclusionRules, err := database.GetExclusionRules(userID)
	if err != nil {
		log.Printf("custom playlist job %s: warning - failed to fetch exclusion rules: %v", job.ID, err)
	}
	if len(exclusionRules) > 0 {
		matchingTrackIDs = filterTrackIDsByExclusions(matchingTrackIDs, exclusionRules)
	}

	if len(matchingTrackIDs) == 0 {
		job.Status = "failed"
		job.Error = "No songs found matching selected genres"
		updateJob()
		return
	}

	job.TotalSongs = len(matchingTrackIDs)
	job.Stage = "creating"
	updateJob()

	// Get user settings for name template
	settings, _ := database.GetUserSettings(userID)
	if settings == nil {
		settings = models.DefaultSettings(userID)
	}

	var playlists []organizer.PlaylistResult

	if req.Mode == "combined" {
		// Create single playlist
		name := req.Name
		if name == "" {
			name = "Custom Mix by Organizer"
		}
		desc := "Custom playlist created by Spotify Genre Organizer"

		playlist, err := spotify.CreatePlaylist(accessToken, userID, name, desc)
		if err != nil {
			log.Printf("custom playlist job %s: failed to create playlist: %v", job.ID, err)
			job.Status = "failed"
			job.Error = "Failed to create playlist"
			updateJob()
			return
		}

		// Add tracks in chunks
		if err := spotify.AddTracksToPlaylist(accessToken, playlist.ID, matchingTrackIDs); err != nil {
			log.Printf("custom playlist job %s: failed to add tracks: %v", job.ID, err)
			job.Status = "failed"
			job.Error = "Failed to add tracks to playlist"
			updateJob()
			return
		}

		playlists = append(playlists, organizer.PlaylistResult{
			SpotifyID:  playlist.ID,
			Name:       name,
			Genre:      "Custom",
			SongCount:  len(matchingTrackIDs),
			SpotifyURL: playlist.ExternalURL,
		})
	} else {
		// Create separate playlists per sub-genre
		for _, subGenre := range req.SubGenres {
			// Filter tracks for this sub-genre
			var genreTrackIDs []string
			for trackID, trackGenres := range cache.AnalysisData.SongGenreMap {
				for _, g := range trackGenres {
					if g == subGenre {
						genreTrackIDs = append(genreTrackIDs, trackID)
						break
					}
				}
			}

			if len(genreTrackIDs) == 0 {
				continue
			}

			name := settings.BuildPlaylistName(subGenre)
			desc := settings.BuildDescription(subGenre)

			playlist, err := spotify.CreatePlaylist(accessToken, userID, name, desc)
			if err != nil {
				log.Printf("custom playlist job %s: failed to create playlist for %s: %v", job.ID, subGenre, err)
				continue
			}

			if err := spotify.AddTracksToPlaylist(accessToken, playlist.ID, genreTrackIDs); err != nil {
				log.Printf("custom playlist job %s: failed to add tracks for %s: %v", job.ID, subGenre, err)
				continue
			}

			playlists = append(playlists, organizer.PlaylistResult{
				SpotifyID:  playlist.ID,
				Name:       name,
				Genre:      subGenre,
				SongCount:  len(genreTrackIDs),
				SpotifyURL: playlist.ExternalURL,
			})

			job.SongsProcessed += len(genreTrackIDs)
			updateJob()
		}
	}

	job.Status = "completed"
	job.Stage = "done"
	job.GenresDiscovered = req.SubGenres
	job.Result = &organizer.OrganizeResult{
		Playlists: playlists,
	}
	updateJob()
}

// filterTrackIDsByExclusions removes excluded tracks
// Note: This is a simplified version - full song data not available from cache
func filterTrackIDsByExclusions(trackIDs []string, rules []models.ExclusionRule) []string {
	blockedSongs := make(map[string]bool)
	for _, rule := range rules {
		if rule.ExclusionType == models.ExclusionTypeSong {
			blockedSongs[rule.SpotifyID] = true
		}
	}

	// For artist exclusions, we'd need track->artist mapping which isn't in cache
	// This will be handled when we have full song data in organize flow
	// Custom playlist from cache only filters song-level exclusions

	result := make([]string, 0, len(trackIDs))
	for _, id := range trackIDs {
		if !blockedSongs[id] {
			result = append(result, id)
		}
	}
	return result
}

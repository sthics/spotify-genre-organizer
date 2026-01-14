package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spotify-genre-organizer/backend/internal/organizer"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type OrganizeRequest struct {
	PlaylistCount   int  `json:"playlist_count" binding:"required,min=1,max=50"`
	ReplaceExisting bool `json:"replace_existing"`
}

type JobStatus struct {
	ID               string                    `json:"id"`
	Status           string                    `json:"status"`
	Stage            string                    `json:"stage"`
	SongsProcessed   int                       `json:"songs_processed"`
	TotalSongs       int                       `json:"total_songs"`
	GenresDiscovered []string                  `json:"genres_discovered"`
	Result           *organizer.OrganizeResult `json:"result,omitempty"`
	Error            string                    `json:"error,omitempty"`
}

var (
	jobs   = make(map[string]*JobStatus)
	jobsMu sync.RWMutex
)

func StartOrganize(c *gin.Context) {
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

	var req OrganizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	go processOrganizeJob(job, accessToken, userID, req)

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": jobID,
		"status": "pending",
	})
}

func processOrganizeJob(job *JobStatus, accessToken, userID string, req OrganizeRequest) {
	updateJob := func() {
		jobsMu.Lock()
		jobs[job.ID] = job
		jobsMu.Unlock()
	}

	job.Status = "processing"
	job.Stage = "fetching"
	updateJob()

	// Fetch liked songs
	songs, err := spotify.FetchAllLikedSongs(accessToken, func(processed, total int) {
		job.SongsProcessed = processed
		job.TotalSongs = total
		updateJob()
	})
	if err != nil {
		log.Printf("organize job %s: failed to fetch songs: %v", job.ID, err)
		job.Status = "failed"
		job.Error = "Failed to fetch your liked songs. Please try again."
		updateJob()
		return
	}

	job.Stage = "analyzing"
	updateJob()

	// Fetch artist genres
	artistGenres, err := spotify.FetchAllArtistGenres(accessToken, songs, nil)
	if err != nil {
		log.Printf("organize job %s: failed to fetch artist genres: %v", job.ID, err)
		job.Status = "failed"
		job.Error = "Failed to analyze song genres. Please try again."
		updateJob()
		return
	}

	// Enrich songs with genres
	spotify.EnrichSongsWithGenres(songs, artistGenres)

	// Collect discovered genres for UI
	genreSet := make(map[string]bool)
	for _, song := range songs {
		for _, g := range song.Genres {
			genreSet[g] = true
		}
	}
	for g := range genreSet {
		job.GenresDiscovered = append(job.GenresDiscovered, g)
	}
	updateJob()

	job.Stage = "creating"
	updateJob()

	// Organize into playlists
	result, err := organizer.OrganizeSongs(
		accessToken,
		userID,
		songs,
		req.PlaylistCount,
		req.ReplaceExisting,
		func(stage string, processed, total int) {
			job.SongsProcessed = processed
			job.TotalSongs = total
			updateJob()
		},
	)
	if err != nil {
		log.Printf("organize job %s: failed to create playlists: %v", job.ID, err)
		job.Status = "failed"
		job.Error = "Failed to create playlists. Please try again."
		updateJob()
		return
	}

	job.Status = "completed"
	job.Stage = "done"
	job.Result = result
	updateJob()
}

func GetOrganizeStatus(c *gin.Context) {
	jobID := c.Param("id")

	jobsMu.RLock()
	job, exists := jobs[jobID]
	jobsMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

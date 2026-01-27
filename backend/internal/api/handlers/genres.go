package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/genres"
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type LibraryGenresResponse struct {
	ParentGenres []models.ParentGenreCount `json:"parent_genres"`
	AnalyzedAt   time.Time                 `json:"analyzed_at"`
	TotalSongs   int                       `json:"total_songs"`
}

func GetLibraryGenres(c *gin.Context) {
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

	refresh := c.Query("refresh") == "true"

	// Check cache first (unless refresh requested)
	if !refresh {
		cached, err := database.GetGenreAnalysisCache(userID)
		if err == nil && cached != nil {
			// Check if cache is less than 24 hours old
			if time.Since(cached.AnalyzedAt) < 24*time.Hour {
				c.JSON(http.StatusOK, LibraryGenresResponse{
					ParentGenres: cached.AnalysisData.ParentGenres,
					AnalyzedAt:   cached.AnalyzedAt,
					TotalSongs:   cached.AnalysisData.TotalSongs,
				})
				return
			}
		}
	}

	// Fetch and analyze library
	songs, err := spotify.FetchAllLikedSongs(accessToken, nil)
	if err != nil {
		log.Printf("GetLibraryGenres: failed to fetch songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch library"})
		return
	}

	// Fetch artist genres
	artistGenres, err := spotify.FetchAllArtistGenres(accessToken, songs, nil)
	if err != nil {
		log.Printf("GetLibraryGenres: failed to fetch artist genres: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze genres"})
		return
	}

	// Enrich songs with genres
	spotify.EnrichSongsWithGenres(songs, artistGenres)

	// Analyze genres
	genreResults, songGenreMap := genres.AnalyzeLibraryGenres(songs)

	// Convert genres.ParentGenreCount to models.ParentGenreCount
	parentGenres := make([]models.ParentGenreCount, len(genreResults))
	for i, pg := range genreResults {
		subGenres := make([]models.SubGenreCount, len(pg.SubGenres))
		for j, sg := range pg.SubGenres {
			subGenres[j] = models.SubGenreCount{
				Name:  sg.Name,
				Count: sg.Count,
			}
		}
		parentGenres[i] = models.ParentGenreCount{
			Name:      pg.Name,
			Count:     pg.Count,
			SubGenres: subGenres,
		}
	}

	// Cache the result
	analysisData := models.GenreAnalysisData{
		ParentGenres: parentGenres,
		TotalSongs:   len(songs),
		SongGenreMap: songGenreMap,
	}

	cache := &models.GenreAnalysisCache{
		ID:           uuid.New().String(),
		UserID:       userID,
		AnalysisData: analysisData,
		AnalyzedAt:   time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := database.SaveGenreAnalysisCache(cache); err != nil {
		log.Printf("GetLibraryGenres: failed to cache analysis: %v", err)
	}

	c.JSON(http.StatusOK, LibraryGenresResponse{
		ParentGenres: parentGenres,
		AnalyzedAt:   cache.AnalyzedAt,
		TotalSongs:   len(songs),
	})
}

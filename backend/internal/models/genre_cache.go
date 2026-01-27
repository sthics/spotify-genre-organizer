package models

import (
	"time"
)

type SubGenreCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ParentGenreCount struct {
	Name      string          `json:"name"`
	Count     int             `json:"count"`
	SubGenres []SubGenreCount `json:"sub_genres"`
}

type GenreAnalysisData struct {
	ParentGenres []ParentGenreCount `json:"parent_genres"`
	TotalSongs   int                `json:"total_songs"`
	// Maps track URI to list of sub-genres for quick filtering
	SongGenreMap map[string][]string `json:"song_genre_map"`
}

type GenreAnalysisCache struct {
	ID           string            `json:"id" db:"id"`
	UserID       string            `json:"user_id" db:"user_id"`
	AnalysisData GenreAnalysisData `json:"analysis_data" db:"analysis_data"`
	AnalyzedAt   time.Time         `json:"analyzed_at" db:"analyzed_at"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
}

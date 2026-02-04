package models

import "time"

// ExclusionType represents the type of exclusion (artist or song)
type ExclusionType string

const (
	ExclusionTypeArtist ExclusionType = "artist"
	ExclusionTypeSong   ExclusionType = "song"
)

// ExclusionScope represents how artist exclusions apply
type ExclusionScope string

const (
	ExclusionScopeAllAppearances ExclusionScope = "all_appearances"
	ExclusionScopePrimaryOnly    ExclusionScope = "primary_only"
)

// ExclusionRule represents a user's blocklist entry
type ExclusionRule struct {
	ID            string          `json:"id" db:"id"`
	UserID        string          `json:"user_id" db:"user_id"`
	ExclusionType ExclusionType   `json:"exclusion_type" db:"exclusion_type"`
	SpotifyID     string          `json:"spotify_id" db:"spotify_id"`
	Name          string          `json:"name" db:"name"`
	Scope         *ExclusionScope `json:"scope,omitempty" db:"scope"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

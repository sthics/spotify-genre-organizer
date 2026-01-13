package models

import (
	"strings"
	"time"
)

type UserSettings struct {
	UserID              string    `json:"user_id" db:"user_id"`
	NameTemplate        string    `json:"name_template" db:"name_template"`
	DescriptionTemplate string    `json:"description_template" db:"description_template"`
	IsPremium           bool      `json:"is_premium" db:"is_premium"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type PlaylistOverride struct {
	ID                string    `json:"id" db:"id"`
	UserID            string    `json:"user_id" db:"user_id"`
	PlaylistSpotifyID string    `json:"playlist_spotify_id" db:"playlist_spotify_id"`
	CustomName        *string   `json:"custom_name" db:"custom_name"`
	CustomDescription *string   `json:"custom_description" db:"custom_description"`
	Genre             string    `json:"genre" db:"genre"`
	LastSyncedAt      *time.Time `json:"last_synced_at" db:"last_synced_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

func DefaultSettings(userID string) *UserSettings {
	return &UserSettings{
		UserID:              userID,
		NameTemplate:        "{genre} by Organizer",
		DescriptionTemplate: "Organized by Spotify Genre Organizer",
		IsPremium:           false,
	}
}

// BuildPlaylistName replaces tokens in the template with actual values
func (s *UserSettings) BuildPlaylistName(genre string) string {
	name := s.NameTemplate
	name = strings.ReplaceAll(name, "{genre}", genre)
	name = strings.ReplaceAll(name, "{year}", time.Now().Format("2006"))
	return name
}

// BuildDescription replaces tokens in the description template
func (s *UserSettings) BuildDescription(genre string) string {
	desc := s.DescriptionTemplate
	desc = strings.ReplaceAll(desc, "{genre}", genre)
	desc = strings.ReplaceAll(desc, "{year}", time.Now().Format("2006"))
	return desc
}

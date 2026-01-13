package models

import (
	"strings"
	"time"
)

type UserSettings struct {
	UserID              string    `json:"user_id"`
	NameTemplate        string    `json:"name_template"`
	DescriptionTemplate string    `json:"description_template"`
	IsPremium           bool      `json:"is_premium"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type PlaylistOverride struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	PlaylistSpotifyID string     `json:"playlist_spotify_id"`
	CustomName        *string    `json:"custom_name"`
	CustomDescription *string    `json:"custom_description"`
	Genre             string     `json:"genre"`
	LastSyncedAt      *time.Time `json:"last_synced_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// DefaultSettings returns default user settings
func DefaultSettings(userID string) *UserSettings {
	return &UserSettings{
		UserID:              userID,
		NameTemplate:        "{genre} by Organizer",
		DescriptionTemplate: "Organized by Spotify Genre Organizer",
		IsPremium:           false,
	}
}

// BuildPlaylistName applies template to generate playlist name
func (s *UserSettings) BuildPlaylistName(genre string) string {
	return replaceTokens(s.NameTemplate, genre)
}

// BuildDescription applies template and appends footer for free users
func (s *UserSettings) BuildDescription(genre string) string {
	desc := replaceTokens(s.DescriptionTemplate, genre)
	if !s.IsPremium {
		desc += " • spotifygenreorganizer.com"
	}
	return desc
}

func replaceTokens(template, genre string) string {
	result := template
	result = strings.ReplaceAll(result, "{genre}", genre)
	return result
}

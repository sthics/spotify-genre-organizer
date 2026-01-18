package database

import (
	"time"
)

// User represents a user in the database
type User struct {
	ID            string    `json:"id"`
	SpotifyID     string    `json:"spotify_id"`
	DisplayName   string    `json:"display_name"`
	Email         string    `json:"email"`
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UpsertUser creates or updates a user in the database
func UpsertUser(spotifyID, displayName, email, accessToken string, tokenExpiresAt time.Time) error {
	user := map[string]interface{}{
		"spotify_id":       spotifyID,
		"display_name":     displayName,
		"email":            email,
		"access_token":     accessToken,
		"token_expires_at": tokenExpiresAt,
		"updated_at":       time.Now(),
	}

	_, _, err := Client.From("users").
		Upsert(user, "spotify_id", "", "").
		Execute()

	return err
}

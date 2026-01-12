package models

import (
	"time"
)

type User struct {
	ID             string     `json:"id"`
	SpotifyID      string     `json:"spotify_id"`
	DisplayName    string     `json:"display_name"`
	Email          string     `json:"email"`
	AccessToken    string     `json:"-"`
	RefreshToken   string     `json:"-"`
	TokenExpiresAt *time.Time `json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

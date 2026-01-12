package spotify

import (
	"os"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func NewConfig() *Config {
	return &Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("SPOTIFY_REDIRECT_URI"),
	}
}

const (
	AuthURL  = "https://accounts.spotify.com/authorize"
	TokenURL = "https://accounts.spotify.com/api/token"
	APIURL   = "https://api.spotify.com/v1"
)

var Scopes = []string{
	"user-library-read",
	"playlist-modify-public",
	"playlist-modify-private",
	"user-read-email",
	"user-read-private",
}

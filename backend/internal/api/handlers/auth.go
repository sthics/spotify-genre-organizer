package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

var spotifyConfig *spotify.Config

func getSpotifyConfig() *spotify.Config {
	if spotifyConfig == nil {
		spotifyConfig = spotify.NewConfig()
	}
	return spotifyConfig
}

func generateState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func isProduction() bool {
	return os.Getenv("ENV") == "production"
}

func Login(c *gin.Context) {
	state := generateState()
	c.SetCookie("oauth_state", state, 600, "/", "", isProduction(), true)
	authURL := getSpotifyConfig().GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error="+errorParam)
		return
	}

	storedState, _ := c.Cookie("oauth_state")
	if state != storedState {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=state_mismatch")
		return
	}

	tokens, err := getSpotifyConfig().ExchangeCode(code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=token_exchange_failed")
		return
	}

	profile, err := spotify.GetUserProfile(tokens.AccessToken)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=profile_fetch_failed")
		return
	}

	expiresAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
	secure := isProduction()
	c.SetCookie("user_id", profile.ID, tokens.ExpiresIn, "/", "", secure, true)
	c.SetCookie("access_token", tokens.AccessToken, tokens.ExpiresIn, "/", "", secure, true)
	_ = expiresAt

	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"/dashboard")
}

func Me(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	profile, err := spotify.GetUserProfile(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           profile.ID,
		"display_name": profile.DisplayName,
		"email":        profile.Email,
	})
}

func Logout(c *gin.Context) {
	secure := isProduction()
	c.SetCookie("user_id", "", -1, "/", "", secure, true)
	c.SetCookie("access_token", "", -1, "/", "", secure, true)
	c.SetCookie("oauth_state", "", -1, "/", "", secure, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

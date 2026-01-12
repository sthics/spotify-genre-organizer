package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spotify-genre-organizer/backend/internal/api"
	"github.com/spotify-genre-organizer/backend/internal/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Debug: Check if Spotify credentials are loaded
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientID == "" {
		log.Println("WARNING: SPOTIFY_CLIENT_ID is empty!")
	} else {
		log.Printf("Spotify Client ID loaded: %s...", clientID[:8])
	}

	if err := database.Init(); err != nil {
		log.Printf("Warning: Could not connect to Supabase: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	api.SetupRoutes(r)

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

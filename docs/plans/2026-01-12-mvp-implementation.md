# Spotify Genre Organizer MVP - Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a web app that organizes Spotify liked songs into genre-based playlists with one click.

**Architecture:** Next.js frontend calls Go/Gin backend API. Backend handles Spotify OAuth, fetches liked songs, maps genres, and creates playlists. Supabase provides PostgreSQL database and auth token storage.

**Tech Stack:** Next.js 14 + TypeScript + Tailwind (frontend), Go 1.21 + Gin (backend), Supabase (database), Spotify Web API

---

## Phase 1: Foundation

### Task 1: Initialize Go Backend

**Files:**
- Create: `backend/go.mod`
- Create: `backend/cmd/api/main.go`
- Create: `backend/internal/api/routes.go`
- Create: `backend/internal/api/handlers/health.go`

**Step 1: Create Go module**

```bash
mkdir -p backend && cd backend
go mod init github.com/spotify-genre-organizer/backend
```

**Step 2: Install dependencies**

```bash
cd backend
go get github.com/gin-gonic/gin
go get github.com/joho/godotenv
```

**Step 3: Create main.go entry point**

Create `backend/cmd/api/main.go`:

```go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spotify-genre-organizer/backend/internal/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
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
```

**Step 4: Create routes.go**

Create `backend/internal/api/routes.go`:

```go
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/login", handlers.NotImplemented)
			auth.GET("/callback", handlers.NotImplemented)
			auth.GET("/me", handlers.NotImplemented)
			auth.POST("/logout", handlers.NotImplemented)
		}

		api.POST("/organize", handlers.NotImplemented)
		api.GET("/organize/:id", handlers.NotImplemented)
	}
}
```

**Step 5: Create health handler**

Create `backend/internal/api/handlers/health.go`:

```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "spotify-genre-organizer",
	})
}

func NotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "not implemented",
	})
}
```

**Step 6: Create .env.example**

Create `backend/.env.example`:

```env
PORT=8080
ENV=development

# Spotify OAuth
SPOTIFY_CLIENT_ID=your_client_id
SPOTIFY_CLIENT_SECRET=your_client_secret
SPOTIFY_REDIRECT_URI=http://localhost:8080/api/auth/callback

# Supabase
SUPABASE_URL=your_supabase_url
SUPABASE_KEY=your_supabase_anon_key

# Frontend
FRONTEND_URL=http://localhost:3000
```

**Step 7: Run and verify**

```bash
cd backend
cp .env.example .env
go run cmd/api/main.go
```

Test: `curl http://localhost:8080/health`
Expected: `{"service":"spotify-genre-organizer","status":"ok"}`

**Step 8: Commit**

```bash
git add backend/
git commit -m "feat: initialize Go backend with Gin and health endpoint"
```

---

### Task 2: Initialize Next.js Frontend

**Files:**
- Create: `frontend/` (entire Next.js project)
- Modify: `frontend/tailwind.config.ts`
- Create: `frontend/src/app/globals.css`

**Step 1: Create Next.js project**

```bash
npx create-next-app@14 frontend --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"
```

Select defaults when prompted.

**Step 2: Update tailwind.config.ts with design tokens**

Replace `frontend/tailwind.config.ts`:

```typescript
import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        'bg-dark': '#1a1a1a',
        'bg-card': '#252525',
        'text-cream': '#f5f0e6',
        'text-muted': '#8a8a8a',
        'accent-orange': '#e85d04',
        'accent-orange-hover': '#ff6b0a',
        'success-green': '#2d936c',
      },
      fontFamily: {
        display: ['Instrument Serif', 'Georgia', 'serif'],
        body: ['IBM Plex Sans', '-apple-system', 'sans-serif'],
      },
      animation: {
        'spin-slow': 'spin 8s linear infinite',
        'spin-vinyl': 'spin 3s linear infinite',
        'bounce-in': 'bounceIn 0.5s ease-out',
        'fade-in': 'fadeIn 0.5s ease-out',
        'drop-in': 'dropIn 0.4s ease-out',
      },
      keyframes: {
        bounceIn: {
          '0%': { transform: 'scale(0.3)', opacity: '0' },
          '50%': { transform: 'scale(1.05)' },
          '70%': { transform: 'scale(0.9)' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
        fadeIn: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        dropIn: {
          '0%': { transform: 'translateY(-50px)', opacity: '0' },
          '60%': { transform: 'translateY(5px)' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
  plugins: [],
};

export default config;
```

**Step 3: Update globals.css**

Replace `frontend/src/app/globals.css`:

```css
@import url('https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@400;500;600&display=swap');

@tailwind base;
@tailwind components;
@tailwind utilities;

/* Instrument Serif from Google Fonts */
@import url('https://fonts.googleapis.com/css2?family=Instrument+Serif&display=swap');

body {
  font-family: 'IBM Plex Sans', -apple-system, sans-serif;
  background-color: #1a1a1a;
  color: #f5f0e6;
}

/* Noise/grain overlay */
.grain-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  opacity: 0.03;
  z-index: 1000;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%' height='100%' filter='url(%23noise)'/%3E%3C/svg%3E");
}

/* Vinyl record animation */
@keyframes vinyl-spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.vinyl-spinning {
  animation: vinyl-spin 3s linear infinite;
}
```

**Step 4: Create layout.tsx**

Replace `frontend/src/app/layout.tsx`:

```tsx
import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Spotify Genre Organizer",
  description: "Organize your Spotify liked songs into genre playlists automatically",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-bg-dark text-text-cream antialiased">
        <div className="grain-overlay" />
        {children}
      </body>
    </html>
  );
}
```

**Step 5: Create placeholder page.tsx**

Replace `frontend/src/app/page.tsx`:

```tsx
export default function Home() {
  return (
    <main className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <h1 className="font-display text-5xl mb-4">Spotify Genre Organizer</h1>
        <p className="text-text-muted">Coming soon...</p>
      </div>
    </main>
  );
}
```

**Step 6: Run and verify**

```bash
cd frontend
npm run dev
```

Visit http://localhost:3000 - should see styled placeholder page with dark background and grain overlay.

**Step 7: Commit**

```bash
git add frontend/
git commit -m "feat: initialize Next.js frontend with design tokens and fonts"
```

---

### Task 3: Set Up Supabase Database

**Files:**
- Create: `supabase/migrations/001_initial_schema.sql`
- Create: `backend/internal/database/supabase.go`

**Step 1: Create migrations directory**

```bash
mkdir -p supabase/migrations
```

**Step 2: Create initial migration**

Create `supabase/migrations/001_initial_schema.sql`:

```sql
-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  spotify_id VARCHAR(255) UNIQUE NOT NULL,
  display_name VARCHAR(255),
  email VARCHAR(255),
  access_token TEXT,
  refresh_token TEXT,
  token_expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Organize jobs table
CREATE TABLE organize_jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  playlist_count INTEGER NOT NULL,
  replace_existing BOOLEAN DEFAULT TRUE,
  status VARCHAR(50) DEFAULT 'pending',
  songs_processed INTEGER DEFAULT 0,
  total_songs INTEGER,
  playlists_created JSONB DEFAULT '[]'::jsonb,
  error_message TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  completed_at TIMESTAMPTZ
);

-- Genre mappings table (seeded data)
CREATE TABLE genre_mappings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  micro_genre VARCHAR(255) NOT NULL,
  parent_genre VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_spotify_id ON users(spotify_id);
CREATE INDEX idx_organize_jobs_user_id ON organize_jobs(user_id);
CREATE INDEX idx_organize_jobs_status ON organize_jobs(status);
CREATE INDEX idx_genre_mappings_micro ON genre_mappings(micro_genre);

-- Enable Row Level Security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE organize_jobs ENABLE ROW LEVEL SECURITY;
ALTER TABLE genre_mappings ENABLE ROW LEVEL SECURITY;

-- RLS Policies (service role bypasses these)
CREATE POLICY "Users can view own data" ON users
  FOR SELECT USING (true);

CREATE POLICY "Users can view own jobs" ON organize_jobs
  FOR SELECT USING (true);

CREATE POLICY "Anyone can read genre mappings" ON genre_mappings
  FOR SELECT USING (true);
```

**Step 3: Create Supabase client**

Create `backend/internal/database/supabase.go`:

```go
package database

import (
	"os"

	"github.com/supabase-community/supabase-go"
)

var Client *supabase.Client

func Init() error {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	client, err := supabase.NewClient(url, key, nil)
	if err != nil {
		return err
	}

	Client = client
	return nil
}
```

**Step 4: Install Supabase Go client**

```bash
cd backend
go get github.com/supabase-community/supabase-go
```

**Step 5: Initialize database in main.go**

Update `backend/cmd/api/main.go`:

```go
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
```

**Step 6: Commit**

```bash
git add supabase/ backend/
git commit -m "feat: add Supabase database schema and Go client"
```

---

### Task 4: Implement Spotify OAuth Flow

**Files:**
- Create: `backend/internal/spotify/client.go`
- Create: `backend/internal/spotify/oauth.go`
- Modify: `backend/internal/api/handlers/auth.go`
- Create: `backend/internal/models/user.go`

**Step 1: Create Spotify client**

Create `backend/internal/spotify/client.go`:

```go
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

// Scopes needed for our app
var Scopes = []string{
	"user-library-read",
	"playlist-modify-public",
	"playlist-modify-private",
	"user-read-email",
	"user-read-private",
}
```

**Step 2: Create OAuth handler**

Create `backend/internal/spotify/oauth.go`:

```go
package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type UserProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *Config) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", c.ClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", c.RedirectURI)
	params.Set("scope", strings.Join(Scopes, " "))
	params.Set("state", state)

	return fmt.Sprintf("%s?%s", AuthURL, params.Encode())
}

func (c *Config) ExchangeCode(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectURI)

	req, err := http.NewRequest("POST", TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.ClientID + ":" + c.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %d", resp.StatusCode)
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func GetUserProfile(accessToken string) (*UserProfile, error) {
	req, err := http.NewRequest("GET", APIURL+"/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user profile: %d", resp.StatusCode)
	}

	var profile UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
```

**Step 3: Create user model**

Create `backend/internal/models/user.go`:

```go
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
```

**Step 4: Create auth handlers**

Create `backend/internal/api/handlers/auth.go`:

```go
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

var spotifyConfig = spotify.NewConfig()

func generateState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func Login(c *gin.Context) {
	state := generateState()

	// In production, store state in session/cookie for verification
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	authURL := spotifyConfig.GetAuthURL(state)
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

	// Verify state (simplified - in production check against stored state)
	storedState, _ := c.Cookie("oauth_state")
	if state != storedState {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=state_mismatch")
		return
	}

	// Exchange code for tokens
	tokens, err := spotifyConfig.ExchangeCode(code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=token_exchange_failed")
		return
	}

	// Get user profile
	profile, err := spotify.GetUserProfile(tokens.AccessToken)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL")+"?error=profile_fetch_failed")
		return
	}

	// TODO: Store user in database
	// For now, set a simple session cookie with user ID
	expiresAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)

	c.SetCookie("user_id", profile.ID, tokens.ExpiresIn, "/", "", false, true)
	c.SetCookie("access_token", tokens.AccessToken, tokens.ExpiresIn, "/", "", false, true)

	_ = expiresAt // Will use when storing in DB

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
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
```

**Step 5: Update routes.go**

Update `backend/internal/api/routes.go`:

```go
package api

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/login", handlers.Login)
			auth.GET("/callback", handlers.Callback)
			auth.GET("/me", handlers.Me)
			auth.POST("/logout", handlers.Logout)
		}

		api.POST("/organize", handlers.NotImplemented)
		api.GET("/organize/:id", handlers.NotImplemented)
	}
}
```

**Step 6: Install CORS middleware**

```bash
cd backend
go get github.com/gin-contrib/cors
```

**Step 7: Test OAuth flow**

1. Start backend: `cd backend && go run cmd/api/main.go`
2. Visit: http://localhost:8080/api/auth/login
3. Should redirect to Spotify login
4. After auth, redirects to frontend with session cookie

**Step 8: Commit**

```bash
git add backend/
git commit -m "feat: implement Spotify OAuth flow with login, callback, me, logout"
```

---

## Phase 2: Core Logic

### Task 5: Fetch Liked Songs from Spotify

**Files:**
- Create: `backend/internal/spotify/library.go`
- Create: `backend/internal/spotify/library_test.go`

**Step 1: Write the failing test**

Create `backend/internal/spotify/library_test.go`:

```go
package spotify

import (
	"testing"
)

func TestParseLikedSongsResponse(t *testing.T) {
	// Test parsing of Spotify API response
	jsonData := `{
		"items": [
			{
				"track": {
					"id": "track123",
					"name": "Test Song",
					"artists": [
						{
							"id": "artist123",
							"name": "Test Artist",
							"genres": ["indie rock", "alternative"]
						}
					]
				}
			}
		],
		"total": 1,
		"next": null
	}`

	songs, total, next, err := ParseLikedSongsResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}

	if next != "" {
		t.Errorf("expected empty next, got %s", next)
	}

	if len(songs) != 1 {
		t.Fatalf("expected 1 song, got %d", len(songs))
	}

	if songs[0].ID != "track123" {
		t.Errorf("expected track ID track123, got %s", songs[0].ID)
	}
}
```

**Step 2: Run test to verify it fails**

```bash
cd backend
go test ./internal/spotify/... -v
```

Expected: FAIL with "undefined: ParseLikedSongsResponse"

**Step 3: Write implementation**

Create `backend/internal/spotify/library.go`:

```go
package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Song struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Artists  []Artist `json:"artists"`
	Genres   []string `json:"genres"` // Aggregated from artists
}

type Artist struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

type likedSongsResponse struct {
	Items []struct {
		Track struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"track"`
	} `json:"items"`
	Total int     `json:"total"`
	Next  *string `json:"next"`
}

func ParseLikedSongsResponse(data []byte) ([]Song, int, string, error) {
	var resp likedSongsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, 0, "", err
	}

	songs := make([]Song, len(resp.Items))
	for i, item := range resp.Items {
		artists := make([]Artist, len(item.Track.Artists))
		for j, a := range item.Track.Artists {
			artists[j] = Artist{
				ID:   a.ID,
				Name: a.Name,
			}
		}
		songs[i] = Song{
			ID:      item.Track.ID,
			Name:    item.Track.Name,
			Artists: artists,
		}
	}

	next := ""
	if resp.Next != nil {
		next = *resp.Next
	}

	return songs, resp.Total, next, nil
}

func FetchLikedSongs(accessToken string, limit, offset int) ([]Song, int, string, error) {
	url := fmt.Sprintf("%s/me/tracks?limit=%d&offset=%d", APIURL, limit, offset)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, "", fmt.Errorf("failed to fetch liked songs: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", err
	}

	return ParseLikedSongsResponse(body)
}

func FetchAllLikedSongs(accessToken string, progressCallback func(processed, total int)) ([]Song, error) {
	var allSongs []Song
	limit := 50
	offset := 0
	total := 0

	for {
		songs, t, _, err := FetchLikedSongs(accessToken, limit, offset)
		if err != nil {
			return nil, err
		}

		if total == 0 {
			total = t
		}

		allSongs = append(allSongs, songs...)

		if progressCallback != nil {
			progressCallback(len(allSongs), total)
		}

		if len(songs) < limit || len(allSongs) >= total {
			break
		}

		offset += limit
	}

	return allSongs, nil
}
```

**Step 4: Run test to verify it passes**

```bash
cd backend
go test ./internal/spotify/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/spotify/
git commit -m "feat: add liked songs fetching from Spotify API"
```

---

### Task 6: Build Genre Consolidation Mapping

**Files:**
- Create: `backend/internal/genres/mapping.go`
- Create: `backend/internal/genres/mapping_test.go`

**Step 1: Write the failing test**

Create `backend/internal/genres/mapping_test.go`:

```go
package genres

import (
	"testing"
)

func TestConsolidateGenre(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"indie rock", "Rock"},
		{"alternative rock", "Rock"},
		{"classic rock", "Rock"},
		{"edm", "Electronic"},
		{"house", "Electronic"},
		{"hip hop", "Hip-Hop"},
		{"rap", "Hip-Hop"},
		{"jazz fusion", "Jazz"},
		{"unknown genre xyz", "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ConsolidateGenre(tt.input)
			if result != tt.expected {
				t.Errorf("ConsolidateGenre(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetParentGenres(t *testing.T) {
	genres := GetParentGenres()

	if len(genres) < 10 {
		t.Errorf("expected at least 10 parent genres, got %d", len(genres))
	}

	// Check some expected genres exist
	expected := []string{"Rock", "Pop", "Hip-Hop", "Electronic", "Jazz", "Classical"}
	for _, e := range expected {
		found := false
		for _, g := range genres {
			if g == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected parent genre %q not found", e)
		}
	}
}
```

**Step 2: Run test to verify it fails**

```bash
cd backend
go test ./internal/genres/... -v
```

Expected: FAIL with "undefined: ConsolidateGenre"

**Step 3: Write implementation**

Create `backend/internal/genres/mapping.go`:

```go
package genres

import "strings"

// Parent genres - the consolidated categories
var ParentGenres = []string{
	"Rock",
	"Pop",
	"Hip-Hop",
	"Electronic",
	"R&B",
	"Jazz",
	"Classical",
	"Country",
	"Metal",
	"Folk",
	"Latin",
	"Blues",
	"Reggae",
	"Punk",
	"Indie",
	"Soul",
	"Funk",
	"World",
	"Other",
}

// Mapping from micro-genres to parent genres
var genreMapping = map[string]string{
	// Rock
	"rock":             "Rock",
	"indie rock":       "Rock",
	"alternative rock": "Rock",
	"garage rock":      "Rock",
	"classic rock":     "Rock",
	"hard rock":        "Rock",
	"soft rock":        "Rock",
	"progressive rock": "Rock",
	"psychedelic rock": "Rock",
	"art rock":         "Rock",
	"glam rock":        "Rock",
	"grunge":           "Rock",
	"post-rock":        "Rock",
	"shoegaze":         "Rock",
	"britpop":          "Rock",

	// Pop
	"pop":          "Pop",
	"indie pop":    "Pop",
	"synth-pop":    "Pop",
	"electropop":   "Pop",
	"dance pop":    "Pop",
	"art pop":      "Pop",
	"dream pop":    "Pop",
	"chamber pop":  "Pop",
	"power pop":    "Pop",
	"teen pop":     "Pop",
	"k-pop":        "Pop",
	"j-pop":        "Pop",

	// Hip-Hop
	"hip hop":          "Hip-Hop",
	"rap":              "Hip-Hop",
	"trap":             "Hip-Hop",
	"conscious hip hop":"Hip-Hop",
	"gangsta rap":      "Hip-Hop",
	"underground hip hop": "Hip-Hop",
	"boom bap":         "Hip-Hop",
	"drill":            "Hip-Hop",
	"crunk":            "Hip-Hop",
	"grime":            "Hip-Hop",

	// Electronic
	"electronic":    "Electronic",
	"edm":           "Electronic",
	"house":         "Electronic",
	"techno":        "Electronic",
	"trance":        "Electronic",
	"dubstep":       "Electronic",
	"drum and bass": "Electronic",
	"ambient":       "Electronic",
	"idm":           "Electronic",
	"downtempo":     "Electronic",
	"trip hop":      "Electronic",
	"chillwave":     "Electronic",
	"synthwave":     "Electronic",
	"deep house":    "Electronic",
	"tech house":    "Electronic",
	"progressive house": "Electronic",

	// R&B
	"r&b":             "R&B",
	"rnb":             "R&B",
	"contemporary r&b": "R&B",
	"neo soul":        "R&B",
	"new jack swing":  "R&B",
	"quiet storm":     "R&B",

	// Jazz
	"jazz":         "Jazz",
	"jazz fusion":  "Jazz",
	"smooth jazz":  "Jazz",
	"bebop":        "Jazz",
	"cool jazz":    "Jazz",
	"free jazz":    "Jazz",
	"acid jazz":    "Jazz",
	"nu jazz":      "Jazz",
	"swing":        "Jazz",
	"big band":     "Jazz",

	// Classical
	"classical":          "Classical",
	"baroque":            "Classical",
	"romantic":           "Classical",
	"contemporary classical": "Classical",
	"opera":              "Classical",
	"orchestral":         "Classical",
	"chamber music":      "Classical",
	"symphony":           "Classical",

	// Country
	"country":          "Country",
	"country rock":     "Country",
	"alt-country":      "Country",
	"bluegrass":        "Country",
	"americana":        "Country",
	"outlaw country":   "Country",
	"country pop":      "Country",

	// Metal
	"metal":           "Metal",
	"heavy metal":     "Metal",
	"thrash metal":    "Metal",
	"death metal":     "Metal",
	"black metal":     "Metal",
	"doom metal":      "Metal",
	"power metal":     "Metal",
	"progressive metal": "Metal",
	"nu metal":        "Metal",
	"metalcore":       "Metal",

	// Folk
	"folk":         "Folk",
	"indie folk":   "Folk",
	"folk rock":    "Folk",
	"freak folk":   "Folk",
	"contemporary folk": "Folk",
	"traditional folk": "Folk",

	// Latin
	"latin":      "Latin",
	"reggaeton":  "Latin",
	"salsa":      "Latin",
	"bachata":    "Latin",
	"cumbia":     "Latin",
	"bossa nova": "Latin",
	"latin pop":  "Latin",
	"latin rock": "Latin",

	// Blues
	"blues":         "Blues",
	"electric blues": "Blues",
	"delta blues":   "Blues",
	"chicago blues": "Blues",
	"blues rock":    "Blues",

	// Reggae
	"reggae":    "Reggae",
	"dub":       "Reggae",
	"ska":       "Reggae",
	"dancehall": "Reggae",
	"roots reggae": "Reggae",

	// Punk
	"punk":         "Punk",
	"punk rock":    "Punk",
	"pop punk":     "Punk",
	"post-punk":    "Punk",
	"hardcore punk": "Punk",
	"emo":          "Punk",
	"skate punk":   "Punk",

	// Indie
	"indie":      "Indie",
	"lo-fi":      "Indie",
	"bedroom pop": "Indie",

	// Soul
	"soul":          "Soul",
	"motown":        "Soul",
	"northern soul": "Soul",
	"southern soul": "Soul",

	// Funk
	"funk":     "Funk",
	"p-funk":   "Funk",
	"funk rock": "Funk",
	"disco":    "Funk",

	// World
	"world":    "World",
	"afrobeat": "World",
	"afropop":  "World",
	"celtic":   "World",
	"flamenco": "World",
	"indian":   "World",
	"middle eastern": "World",
}

func ConsolidateGenre(microGenre string) string {
	normalized := strings.ToLower(strings.TrimSpace(microGenre))

	// Direct lookup
	if parent, ok := genreMapping[normalized]; ok {
		return parent
	}

	// Try partial matching for compound genres
	for micro, parent := range genreMapping {
		if strings.Contains(normalized, micro) || strings.Contains(micro, normalized) {
			return parent
		}
	}

	return "Other"
}

func GetParentGenres() []string {
	return ParentGenres
}

// ConsolidateGenres takes a list of micro-genres and returns unique parent genres
func ConsolidateGenres(microGenres []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, g := range microGenres {
		parent := ConsolidateGenre(g)
		if !seen[parent] {
			seen[parent] = true
			result = append(result, parent)
		}
	}

	return result
}
```

**Step 4: Run test to verify it passes**

```bash
cd backend
go test ./internal/genres/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/genres/
git commit -m "feat: add genre consolidation mapping (100+ micro-genres to 19 parent categories)"
```

---

### Task 7: Fetch Artist Genres from Spotify

**Files:**
- Create: `backend/internal/spotify/artists.go`
- Modify: `backend/internal/spotify/library.go`

**Step 1: Create artists.go**

Create `backend/internal/spotify/artists.go`:

```go
package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ArtistDetails struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

type artistsResponse struct {
	Artists []ArtistDetails `json:"artists"`
}

// FetchArtists fetches details for multiple artists (max 50 per request)
func FetchArtists(accessToken string, artistIDs []string) ([]ArtistDetails, error) {
	if len(artistIDs) == 0 {
		return nil, nil
	}

	if len(artistIDs) > 50 {
		artistIDs = artistIDs[:50]
	}

	url := fmt.Sprintf("%s/artists?ids=%s", APIURL, strings.Join(artistIDs, ","))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch artists: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result artistsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Artists, nil
}

// FetchAllArtistGenres fetches genres for all unique artists in songs
func FetchAllArtistGenres(accessToken string, songs []Song, progressCallback func(processed, total int)) (map[string][]string, error) {
	// Collect unique artist IDs
	artistSet := make(map[string]bool)
	for _, song := range songs {
		for _, artist := range song.Artists {
			artistSet[artist.ID] = true
		}
	}

	artistIDs := make([]string, 0, len(artistSet))
	for id := range artistSet {
		artistIDs = append(artistIDs, id)
	}

	// Fetch in batches of 50
	genreMap := make(map[string][]string)
	batchSize := 50
	total := len(artistIDs)

	for i := 0; i < len(artistIDs); i += batchSize {
		end := i + batchSize
		if end > len(artistIDs) {
			end = len(artistIDs)
		}

		batch := artistIDs[i:end]
		artists, err := FetchArtists(accessToken, batch)
		if err != nil {
			return nil, err
		}

		for _, artist := range artists {
			genreMap[artist.ID] = artist.Genres
		}

		if progressCallback != nil {
			progressCallback(end, total)
		}

		// Small delay to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	return genreMap, nil
}
```

**Step 2: Add helper to enrich songs with genres**

Add to `backend/internal/spotify/library.go`:

```go
// EnrichSongsWithGenres adds genre information to songs from artist data
func EnrichSongsWithGenres(songs []Song, artistGenres map[string][]string) {
	for i := range songs {
		genreSet := make(map[string]bool)
		for _, artist := range songs[i].Artists {
			if genres, ok := artistGenres[artist.ID]; ok {
				for _, g := range genres {
					genreSet[g] = true
				}
			}
		}

		songs[i].Genres = make([]string, 0, len(genreSet))
		for g := range genreSet {
			songs[i].Genres = append(songs[i].Genres, g)
		}
	}
}
```

**Step 3: Commit**

```bash
git add backend/internal/spotify/
git commit -m "feat: add artist genre fetching and song enrichment"
```

---

### Task 8: Implement Playlist Creation

**Files:**
- Create: `backend/internal/spotify/playlists.go`
- Create: `backend/internal/spotify/playlists_test.go`

**Step 1: Write the failing test**

Create `backend/internal/spotify/playlists_test.go`:

```go
package spotify

import (
	"testing"
)

func TestBuildPlaylistName(t *testing.T) {
	tests := []struct {
		genre    string
		expected string
	}{
		{"Rock", "Rock by Organizer"},
		{"Hip-Hop", "Hip-Hop by Organizer"},
		{"R&B", "R&B by Organizer"},
	}

	for _, tt := range tests {
		result := BuildPlaylistName(tt.genre)
		if result != tt.expected {
			t.Errorf("BuildPlaylistName(%q) = %q, want %q", tt.genre, result, tt.expected)
		}
	}
}

func TestChunkTrackIDs(t *testing.T) {
	ids := make([]string, 150)
	for i := range ids {
		ids[i] = "track" + string(rune('0'+i%10))
	}

	chunks := ChunkTrackIDs(ids, 100)

	if len(chunks) != 2 {
		t.Errorf("expected 2 chunks, got %d", len(chunks))
	}

	if len(chunks[0]) != 100 {
		t.Errorf("expected first chunk to have 100 items, got %d", len(chunks[0]))
	}

	if len(chunks[1]) != 50 {
		t.Errorf("expected second chunk to have 50 items, got %d", len(chunks[1]))
	}
}
```

**Step 2: Run test to verify it fails**

```bash
cd backend
go test ./internal/spotify/... -v -run TestBuildPlaylistName
```

Expected: FAIL

**Step 3: Write implementation**

Create `backend/internal/spotify/playlists.go`:

```go
package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Playlist struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ExternalURL string `json:"external_urls"`
	TracksTotal int    `json:"tracks_total"`
}

type createPlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type addTracksRequest struct {
	URIs []string `json:"uris"`
}

func BuildPlaylistName(genre string) string {
	return genre + " by Organizer"
}

func ChunkTrackIDs(ids []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}
		chunks = append(chunks, ids[i:end])
	}
	return chunks
}

func CreatePlaylist(accessToken, userID, name, description string) (*Playlist, error) {
	url := fmt.Sprintf("%s/users/%s/playlists", APIURL, userID)

	body := createPlaylistRequest{
		Name:        name,
		Description: description,
		Public:      false,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create playlist: %d - %s", resp.StatusCode, string(respBody))
	}

	var playlist struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return nil, err
	}

	return &Playlist{
		ID:          playlist.ID,
		Name:        playlist.Name,
		ExternalURL: playlist.ExternalURLs.Spotify,
	}, nil
}

func AddTracksToPlaylist(accessToken, playlistID string, trackIDs []string) error {
	// Spotify requires track URIs in format "spotify:track:ID"
	uris := make([]string, len(trackIDs))
	for i, id := range trackIDs {
		uris[i] = "spotify:track:" + id
	}

	// Spotify limits to 100 tracks per request
	chunks := ChunkTrackIDs(uris, 100)

	for _, chunk := range chunks {
		url := fmt.Sprintf("%s/playlists/%s/tracks", APIURL, playlistID)

		body := addTracksRequest{URIs: chunk}
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to add tracks: %d", resp.StatusCode)
		}

		time.Sleep(100 * time.Millisecond) // Rate limiting
	}

	return nil
}

func ClearPlaylist(accessToken, playlistID string) error {
	url := fmt.Sprintf("%s/playlists/%s/tracks", APIURL, playlistID)

	// First get all tracks
	req, err := http.NewRequest("GET", url+"?limit=100", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tracksResp struct {
		Items []struct {
			Track struct {
				URI string `json:"uri"`
			} `json:"track"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tracksResp); err != nil {
		return err
	}

	if len(tracksResp.Items) == 0 {
		return nil
	}

	// Remove all tracks
	tracks := make([]map[string]string, len(tracksResp.Items))
	for i, item := range tracksResp.Items {
		tracks[i] = map[string]string{"uri": item.Track.URI}
	}

	deleteBody, _ := json.Marshal(map[string]interface{}{"tracks": tracks})

	req, err = http.NewRequest("DELETE", url, bytes.NewReader(deleteBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func FindExistingPlaylist(accessToken, playlistName string) (*Playlist, error) {
	url := fmt.Sprintf("%s/me/playlists?limit=50", APIURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var playlistsResp struct {
		Items []struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&playlistsResp); err != nil {
		return nil, err
	}

	for _, p := range playlistsResp.Items {
		if p.Name == playlistName {
			return &Playlist{
				ID:          p.ID,
				Name:        p.Name,
				ExternalURL: p.ExternalURLs.Spotify,
			}, nil
		}
	}

	return nil, nil
}
```

**Step 4: Run tests**

```bash
cd backend
go test ./internal/spotify/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/spotify/
git commit -m "feat: add playlist creation, track adding, and existing playlist lookup"
```

---

### Task 9: Implement Organize Job Handler

**Files:**
- Create: `backend/internal/api/handlers/organize.go`
- Create: `backend/internal/organizer/organizer.go`
- Modify: `backend/internal/api/routes.go`

**Step 1: Create organizer service**

Create `backend/internal/organizer/organizer.go`:

```go
package organizer

import (
	"sort"

	"github.com/spotify-genre-organizer/backend/internal/genres"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type OrganizeResult struct {
	Playlists []PlaylistResult `json:"playlists"`
}

type PlaylistResult struct {
	Name       string `json:"name"`
	Genre      string `json:"genre"`
	SpotifyID  string `json:"spotify_id"`
	SpotifyURL string `json:"spotify_url"`
	SongCount  int    `json:"song_count"`
}

type ProgressCallback func(stage string, processed, total int)

func OrganizeSongs(
	accessToken string,
	userID string,
	songs []spotify.Song,
	playlistCount int,
	replaceExisting bool,
	progress ProgressCallback,
) (*OrganizeResult, error) {
	// Group songs by parent genre
	genreGroups := make(map[string][]spotify.Song)
	for _, song := range songs {
		if len(song.Genres) == 0 {
			genreGroups["Other"] = append(genreGroups["Other"], song)
			continue
		}

		// Use the first consolidated genre
		parentGenre := genres.ConsolidateGenre(song.Genres[0])
		genreGroups[parentGenre] = append(genreGroups[parentGenre], song)
	}

	// Sort genres by song count (descending)
	type genreCount struct {
		genre string
		count int
	}
	var sortedGenres []genreCount
	for genre, songs := range genreGroups {
		sortedGenres = append(sortedGenres, genreCount{genre, len(songs)})
	}
	sort.Slice(sortedGenres, func(i, j int) bool {
		return sortedGenres[i].count > sortedGenres[j].count
	})

	// Limit to requested playlist count
	if len(sortedGenres) > playlistCount {
		// Merge smaller genres into "Other"
		for i := playlistCount; i < len(sortedGenres); i++ {
			genreGroups["Other"] = append(genreGroups["Other"], genreGroups[sortedGenres[i].genre]...)
			delete(genreGroups, sortedGenres[i].genre)
		}
		sortedGenres = sortedGenres[:playlistCount]
	}

	// Create playlists
	var results []PlaylistResult
	total := len(sortedGenres)

	for i, gc := range sortedGenres {
		if progress != nil {
			progress("creating", i+1, total)
		}

		playlistName := spotify.BuildPlaylistName(gc.genre)
		songs := genreGroups[gc.genre]

		var playlist *spotify.Playlist
		var err error

		if replaceExisting {
			// Check for existing playlist
			playlist, err = spotify.FindExistingPlaylist(accessToken, playlistName)
			if err != nil {
				return nil, err
			}

			if playlist != nil {
				// Clear existing tracks
				if err := spotify.ClearPlaylist(accessToken, playlist.ID); err != nil {
					return nil, err
				}
			}
		}

		if playlist == nil {
			// Create new playlist
			playlist, err = spotify.CreatePlaylist(
				accessToken,
				userID,
				playlistName,
				"Organized by Spotify Genre Organizer",
			)
			if err != nil {
				return nil, err
			}
		}

		// Add tracks
		trackIDs := make([]string, len(songs))
		for i, s := range songs {
			trackIDs[i] = s.ID
		}

		if err := spotify.AddTracksToPlaylist(accessToken, playlist.ID, trackIDs); err != nil {
			return nil, err
		}

		results = append(results, PlaylistResult{
			Name:       playlistName,
			Genre:      gc.genre,
			SpotifyID:  playlist.ID,
			SpotifyURL: playlist.ExternalURL,
			SongCount:  len(songs),
		})
	}

	return &OrganizeResult{Playlists: results}, nil
}
```

**Step 2: Create organize handler**

Create `backend/internal/api/handlers/organize.go`:

```go
package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spotify-genre-organizer/backend/internal/organizer"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type OrganizeRequest struct {
	PlaylistCount   int  `json:"playlist_count" binding:"required,min=1,max=50"`
	ReplaceExisting bool `json:"replace_existing"`
}

type JobStatus struct {
	ID              string                      `json:"id"`
	Status          string                      `json:"status"`
	Stage           string                      `json:"stage"`
	SongsProcessed  int                         `json:"songs_processed"`
	TotalSongs      int                         `json:"total_songs"`
	GenresDiscovered []string                   `json:"genres_discovered"`
	Result          *organizer.OrganizeResult   `json:"result,omitempty"`
	Error           string                      `json:"error,omitempty"`
}

var (
	jobs   = make(map[string]*JobStatus)
	jobsMu sync.RWMutex
)

func StartOrganize(c *gin.Context) {
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

	var req OrganizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create job
	jobID := uuid.New().String()
	job := &JobStatus{
		ID:     jobID,
		Status: "pending",
		Stage:  "initializing",
	}

	jobsMu.Lock()
	jobs[jobID] = job
	jobsMu.Unlock()

	// Start async processing
	go processOrganizeJob(job, accessToken, userID, req)

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": jobID,
		"status": "pending",
	})
}

func processOrganizeJob(job *JobStatus, accessToken, userID string, req OrganizeRequest) {
	updateJob := func() {
		jobsMu.Lock()
		jobs[job.ID] = job
		jobsMu.Unlock()
	}

	job.Status = "processing"
	job.Stage = "fetching"
	updateJob()

	// Fetch liked songs
	songs, err := spotify.FetchAllLikedSongs(accessToken, func(processed, total int) {
		job.SongsProcessed = processed
		job.TotalSongs = total
		updateJob()
	})
	if err != nil {
		job.Status = "failed"
		job.Error = err.Error()
		updateJob()
		return
	}

	job.Stage = "analyzing"
	updateJob()

	// Fetch artist genres
	artistGenres, err := spotify.FetchAllArtistGenres(accessToken, songs, nil)
	if err != nil {
		job.Status = "failed"
		job.Error = err.Error()
		updateJob()
		return
	}

	// Enrich songs with genres
	spotify.EnrichSongsWithGenres(songs, artistGenres)

	// Collect discovered genres for UI
	genreSet := make(map[string]bool)
	for _, song := range songs {
		for _, g := range song.Genres {
			genreSet[g] = true
		}
	}
	for g := range genreSet {
		job.GenresDiscovered = append(job.GenresDiscovered, g)
	}
	updateJob()

	job.Stage = "creating"
	updateJob()

	// Organize into playlists
	result, err := organizer.OrganizeSongs(
		accessToken,
		userID,
		songs,
		req.PlaylistCount,
		req.ReplaceExisting,
		func(stage string, processed, total int) {
			job.SongsProcessed = processed
			job.TotalSongs = total
			updateJob()
		},
	)
	if err != nil {
		job.Status = "failed"
		job.Error = err.Error()
		updateJob()
		return
	}

	job.Status = "completed"
	job.Stage = "done"
	job.Result = result
	updateJob()
}

func GetOrganizeStatus(c *gin.Context) {
	jobID := c.Param("id")

	jobsMu.RLock()
	job, exists := jobs[jobID]
	jobsMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}
```

**Step 3: Install uuid package**

```bash
cd backend
go get github.com/google/uuid
```

**Step 4: Update routes**

Update `backend/internal/api/routes.go` to use new handlers:

```go
package api

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/login", handlers.Login)
			auth.GET("/callback", handlers.Callback)
			auth.GET("/me", handlers.Me)
			auth.POST("/logout", handlers.Logout)
		}

		api.POST("/organize", handlers.StartOrganize)
		api.GET("/organize/:id", handlers.GetOrganizeStatus)
	}
}
```

**Step 5: Commit**

```bash
git add backend/
git commit -m "feat: implement organize job handler with async processing"
```

---

## Phase 3: Frontend

### Task 10: Create Landing Page

**Files:**
- Create: `frontend/src/components/VinylIcon.tsx`
- Create: `frontend/src/components/Button.tsx`
- Modify: `frontend/src/app/page.tsx`

**Step 1: Create VinylIcon component**

Create `frontend/src/components/VinylIcon.tsx`:

```tsx
interface VinylIconProps {
  className?: string;
  spinning?: boolean;
  size?: number;
}

export function VinylIcon({ className = "", spinning = false, size = 64 }: VinylIconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 100 100"
      className={`${spinning ? 'animate-spin-slow' : ''} ${className}`}
    >
      {/* Outer ring */}
      <circle cx="50" cy="50" r="48" fill="#1a1a1a" stroke="#333" strokeWidth="2" />

      {/* Grooves */}
      <circle cx="50" cy="50" r="40" fill="none" stroke="#2a2a2a" strokeWidth="1" />
      <circle cx="50" cy="50" r="35" fill="none" stroke="#252525" strokeWidth="1" />
      <circle cx="50" cy="50" r="30" fill="none" stroke="#2a2a2a" strokeWidth="1" />
      <circle cx="50" cy="50" r="25" fill="none" stroke="#252525" strokeWidth="1" />
      <circle cx="50" cy="50" r="20" fill="none" stroke="#2a2a2a" strokeWidth="1" />

      {/* Label */}
      <circle cx="50" cy="50" r="15" fill="#e85d04" />
      <circle cx="50" cy="50" r="12" fill="#ff6b0a" />

      {/* Center hole */}
      <circle cx="50" cy="50" r="3" fill="#1a1a1a" />

      {/* Shine effect */}
      <ellipse cx="35" cy="35" rx="8" ry="4" fill="rgba(255,255,255,0.05)" transform="rotate(-45 35 35)" />
    </svg>
  );
}
```

**Step 2: Create Button component**

Create `frontend/src/components/Button.tsx`:

```tsx
import { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode;
  variant?: 'primary' | 'secondary' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
}

export function Button({
  children,
  variant = 'primary',
  size = 'md',
  className = '',
  ...props
}: ButtonProps) {
  const baseStyles = 'font-body font-medium rounded-lg transition-all duration-200 transform active:scale-95';

  const variants = {
    primary: 'bg-accent-orange hover:bg-accent-orange-hover text-white shadow-lg hover:shadow-xl',
    secondary: 'bg-bg-card hover:bg-bg-card/80 text-text-cream border border-text-muted/20',
    ghost: 'bg-transparent hover:bg-bg-card text-text-muted hover:text-text-cream',
  };

  const sizes = {
    sm: 'px-4 py-2 text-sm',
    md: 'px-6 py-3 text-base',
    lg: 'px-8 py-4 text-lg',
  };

  return (
    <button
      className={`${baseStyles} ${variants[variant]} ${sizes[size]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}
```

**Step 3: Create Landing Page**

Replace `frontend/src/app/page.tsx`:

```tsx
'use client';

import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { useEffect, useState } from 'react';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function Home() {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    setIsVisible(true);
  }, []);

  const handleConnect = () => {
    window.location.href = `${API_URL}/api/auth/login`;
  };

  return (
    <main className="min-h-screen flex flex-col items-center justify-center px-4">
      {/* Hero Section */}
      <div className={`text-center max-w-2xl transition-all duration-1000 ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}`}>
        {/* Logo */}
        <div className="flex items-center justify-center gap-4 mb-8">
          <h1 className="font-display text-4xl md:text-5xl lg:text-6xl text-text-cream">
            Spotify Genre
            <br />
            Organizer
          </h1>
          <VinylIcon spinning size={80} />
        </div>

        {/* Tagline */}
        <div className="mb-12 space-y-2">
          <p className="font-display text-2xl md:text-3xl text-text-cream italic">
            "2,000 liked songs.
          </p>
          <p className="font-display text-2xl md:text-3xl text-text-cream italic">
            Zero organization.
          </p>
          <p className="font-display text-2xl md:text-3xl text-text-cream italic">
            Sound familiar?"
          </p>
        </div>

        {/* CTA Button */}
        <Button
          size="lg"
          onClick={handleConnect}
          className="flex items-center gap-2 mx-auto"
        >
          <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z"/>
          </svg>
          Connect with Spotify
        </Button>
      </div>

      {/* Value Props */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-20 max-w-4xl w-full px-4">
        {[
          { title: 'Analyze', desc: 'your library' },
          { title: 'Organize', desc: 'into playlists' },
          { title: 'Enjoy', desc: 'your music' },
        ].map((item, index) => (
          <div
            key={item.title}
            className={`bg-bg-card p-6 rounded-xl text-center transition-all duration-700 delay-${index * 200} ${isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}`}
            style={{ transitionDelay: `${800 + index * 200}ms` }}
          >
            <h3 className="font-display text-xl text-accent-orange mb-2">{item.title}</h3>
            <p className="text-text-muted">{item.desc}</p>
          </div>
        ))}
      </div>
    </main>
  );
}
```

**Step 4: Create env file**

Create `frontend/.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

**Step 5: Test**

```bash
cd frontend
npm run dev
```

Visit http://localhost:3000 - should see styled landing page.

**Step 6: Commit**

```bash
git add frontend/
git commit -m "feat: create landing page with vinyl animation and connect button"
```

---

### Task 11: Create Dashboard Page

**Files:**
- Create: `frontend/src/app/dashboard/page.tsx`
- Create: `frontend/src/components/Slider.tsx`
- Create: `frontend/src/hooks/useUser.ts`
- Create: `frontend/src/lib/api.ts`

**Step 1: Create API client**

Create `frontend/src/lib/api.ts`:

```typescript
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function fetchUser() {
  const res = await fetch(`${API_URL}/api/auth/me`, {
    credentials: 'include',
  });

  if (!res.ok) {
    throw new Error('Not authenticated');
  }

  return res.json();
}

export async function startOrganize(playlistCount: number, replaceExisting: boolean) {
  const res = await fetch(`${API_URL}/api/organize`, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      playlist_count: playlistCount,
      replace_existing: replaceExisting,
    }),
  });

  if (!res.ok) {
    throw new Error('Failed to start organize');
  }

  return res.json();
}

export async function getOrganizeStatus(jobId: string) {
  const res = await fetch(`${API_URL}/api/organize/${jobId}`, {
    credentials: 'include',
  });

  if (!res.ok) {
    throw new Error('Failed to get status');
  }

  return res.json();
}

export async function logout() {
  await fetch(`${API_URL}/api/auth/logout`, {
    method: 'POST',
    credentials: 'include',
  });
}
```

**Step 2: Create useUser hook**

Create `frontend/src/hooks/useUser.ts`:

```typescript
'use client';

import { useEffect, useState } from 'react';
import { fetchUser } from '@/lib/api';
import { useRouter } from 'next/navigation';

interface User {
  id: string;
  display_name: string;
  email: string;
}

export function useUser() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    fetchUser()
      .then(setUser)
      .catch(() => {
        router.push('/');
      })
      .finally(() => setLoading(false));
  }, [router]);

  return { user, loading };
}
```

**Step 3: Create Slider component**

Create `frontend/src/components/Slider.tsx`:

```tsx
interface SliderProps {
  value: number;
  onChange: (value: number) => void;
  min: number;
  max: number;
  label?: string;
}

export function Slider({ value, onChange, min, max, label }: SliderProps) {
  const percentage = ((value - min) / (max - min)) * 100;

  return (
    <div className="w-full">
      {label && (
        <label className="block font-body text-text-muted mb-2">{label}</label>
      )}
      <div className="relative">
        <input
          type="range"
          min={min}
          max={max}
          value={value}
          onChange={(e) => onChange(Number(e.target.value))}
          className="w-full h-2 bg-bg-card rounded-lg appearance-none cursor-pointer accent-accent-orange"
          style={{
            background: `linear-gradient(to right, #e85d04 0%, #e85d04 ${percentage}%, #252525 ${percentage}%, #252525 100%)`,
          }}
        />
        <div className="flex justify-between text-sm text-text-muted mt-1">
          <span>{min}</span>
          <span className="text-accent-orange font-medium text-lg">{value}</span>
          <span>{max}</span>
        </div>
      </div>
    </div>
  );
}
```

**Step 4: Create Dashboard page**

Create `frontend/src/app/dashboard/page.tsx`:

```tsx
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { Slider } from '@/components/Slider';
import { useUser } from '@/hooks/useUser';
import { startOrganize, logout } from '@/lib/api';

export default function Dashboard() {
  const { user, loading } = useUser();
  const router = useRouter();
  const [playlistCount, setPlaylistCount] = useState(12);
  const [replaceExisting, setReplaceExisting] = useState(true);
  const [isOrganizing, setIsOrganizing] = useState(false);

  // Mock liked songs count - in real app, fetch this from API
  const likedSongsCount = 1247;
  const songsPerPlaylist = Math.round(likedSongsCount / playlistCount);

  const handleOrganize = async () => {
    setIsOrganizing(true);
    try {
      const { job_id } = await startOrganize(playlistCount, replaceExisting);
      router.push(`/processing?job=${job_id}`);
    } catch (error) {
      console.error('Failed to start organize:', error);
      setIsOrganizing(false);
    }
  };

  const handleLogout = async () => {
    await logout();
    router.push('/');
  };

  if (loading) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <VinylIcon spinning size={64} />
      </main>
    );
  }

  return (
    <main className="min-h-screen flex flex-col items-center justify-center px-4 py-12">
      {/* Header */}
      <div className="w-full max-w-xl flex items-center justify-between mb-12">
        <div>
          <h1 className="font-display text-2xl text-text-cream">
            Hey, {user?.display_name?.split(' ')[0] || 'there'}.
          </h1>
          <p className="text-text-muted">
            You've got <span className="text-text-cream font-medium">{likedSongsCount.toLocaleString()}</span> liked songs.
          </p>
        </div>
        <button
          onClick={handleLogout}
          className="text-text-muted hover:text-text-cream transition-colors"
          title="Logout"
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </button>
      </div>

      {/* Main Card */}
      <div className="w-full max-w-xl bg-bg-card rounded-2xl p-8 shadow-xl">
        {/* Playlist Count Slider */}
        <div className="mb-8">
          <h2 className="font-display text-xl text-text-cream mb-4">
            How many playlists?
          </h2>
          <Slider
            value={playlistCount}
            onChange={setPlaylistCount}
            min={1}
            max={50}
          />
          <div className="flex items-center gap-2 mt-4 text-text-muted">
            <VinylIcon size={20} />
            <span>~{songsPerPlaylist} songs per playlist</span>
          </div>
        </div>

        {/* Replace Toggle */}
        <div className="mb-8">
          <div className="space-y-3">
            <label className="flex items-start gap-3 cursor-pointer group">
              <input
                type="radio"
                name="replace"
                checked={replaceExisting}
                onChange={() => setReplaceExisting(true)}
                className="mt-1 accent-accent-orange"
              />
              <div>
                <span className="text-text-cream group-hover:text-accent-orange transition-colors">
                  Update existing playlists
                </span>
                <p className="text-sm text-text-muted">
                  Replaces songs in "Rock by Organizer", etc.
                </p>
              </div>
            </label>
            <label className="flex items-start gap-3 cursor-pointer group">
              <input
                type="radio"
                name="replace"
                checked={!replaceExisting}
                onChange={() => setReplaceExisting(false)}
                className="mt-1 accent-accent-orange"
              />
              <div>
                <span className="text-text-cream group-hover:text-accent-orange transition-colors">
                  Create fresh playlists
                </span>
                <p className="text-sm text-text-muted">
                  Keeps your old ones, makes new
                </p>
              </div>
            </label>
          </div>
        </div>

        {/* Organize Button */}
        <Button
          size="lg"
          className="w-full flex items-center justify-center gap-2"
          onClick={handleOrganize}
          disabled={isOrganizing}
        >
          {isOrganizing ? (
            <>
              <VinylIcon spinning size={24} />
              Starting...
            </>
          ) : (
            <>
              Organize My Library
              <span className="text-xl">●</span>
            </>
          )}
        </Button>
      </div>
    </main>
  );
}
```

**Step 5: Commit**

```bash
git add frontend/
git commit -m "feat: create dashboard page with slider and organize options"
```

---

### Task 12: Create Processing Page

**Files:**
- Create: `frontend/src/app/processing/page.tsx`
- Create: `frontend/src/components/ProgressBar.tsx`
- Create: `frontend/src/components/GenreTag.tsx`

**Step 1: Create ProgressBar component**

Create `frontend/src/components/ProgressBar.tsx`:

```tsx
interface ProgressBarProps {
  progress: number; // 0-100
  label?: string;
}

export function ProgressBar({ progress, label }: ProgressBarProps) {
  return (
    <div className="w-full">
      <div className="h-2 bg-bg-card rounded-full overflow-hidden">
        <div
          className="h-full bg-text-cream transition-all duration-300 ease-out"
          style={{ width: `${Math.min(100, Math.max(0, progress))}%` }}
        />
      </div>
      {label && (
        <p className="text-center text-text-muted mt-2">{label}</p>
      )}
    </div>
  );
}
```

**Step 2: Create GenreTag component**

Create `frontend/src/components/GenreTag.tsx`:

```tsx
interface GenreTagProps {
  genre: string;
  index: number;
}

export function GenreTag({ genre, index }: GenreTagProps) {
  return (
    <span
      className="inline-block px-3 py-1 bg-bg-dark text-text-cream rounded-full text-sm font-body animate-bounce-in"
      style={{ animationDelay: `${index * 100}ms` }}
    >
      {genre}
    </span>
  );
}
```

**Step 3: Create Processing page**

Create `frontend/src/app/processing/page.tsx`:

```tsx
'use client';

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { ProgressBar } from '@/components/ProgressBar';
import { GenreTag } from '@/components/GenreTag';
import { getOrganizeStatus } from '@/lib/api';

interface JobStatus {
  id: string;
  status: string;
  stage: string;
  songs_processed: number;
  total_songs: number;
  genres_discovered: string[];
  result?: {
    playlists: Array<{
      name: string;
      genre: string;
      spotify_id: string;
      spotify_url: string;
      song_count: number;
    }>;
  };
  error?: string;
}

function ProcessingContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const jobId = searchParams.get('job');

  const [status, setStatus] = useState<JobStatus | null>(null);
  const [tonearmAngle, setTonearmAngle] = useState(0);

  useEffect(() => {
    if (!jobId) {
      router.push('/dashboard');
      return;
    }

    const pollStatus = async () => {
      try {
        const data = await getOrganizeStatus(jobId);
        setStatus(data);

        // Update tonearm based on progress
        if (data.total_songs > 0) {
          const progress = data.songs_processed / data.total_songs;
          setTonearmAngle(progress * 30); // 0-30 degree sweep
        }

        if (data.status === 'completed') {
          // Redirect to success page
          router.push(`/success?job=${jobId}`);
        } else if (data.status === 'failed') {
          // Handle error
          console.error('Job failed:', data.error);
        } else {
          // Continue polling
          setTimeout(pollStatus, 1000);
        }
      } catch (error) {
        console.error('Failed to get status:', error);
        setTimeout(pollStatus, 2000);
      }
    };

    pollStatus();
  }, [jobId, router]);

  const getStageText = (stage: string) => {
    switch (stage) {
      case 'fetching':
        return 'Analyzing your library...';
      case 'analyzing':
        return 'Detecting genres...';
      case 'creating':
        return 'Creating playlists...';
      default:
        return 'Processing...';
    }
  };

  const progress = status?.total_songs
    ? (status.songs_processed / status.total_songs) * 100
    : 0;

  return (
    <main className="min-h-screen flex flex-col items-center justify-center px-4">
      {/* Vinyl with Tonearm */}
      <div className="relative mb-8">
        <VinylIcon spinning size={120} />
        {/* Tonearm */}
        <div
          className="absolute top-0 right-0 w-16 h-1 bg-text-muted origin-right transition-transform duration-500"
          style={{
            transform: `rotate(${-45 + tonearmAngle}deg)`,
            transformOrigin: 'right center',
          }}
        >
          <div className="absolute right-0 top-1/2 -translate-y-1/2 w-3 h-3 bg-text-cream rounded-full" />
        </div>
      </div>

      {/* Status Text */}
      <h2 className="font-display text-2xl text-text-cream mb-4">
        {getStageText(status?.stage || '')}
      </h2>

      {/* Progress Bar */}
      <div className="w-full max-w-md mb-8">
        <ProgressBar
          progress={progress}
          label={status?.total_songs
            ? `${status.songs_processed.toLocaleString()} / ${status.total_songs.toLocaleString()} songs`
            : undefined
          }
        />
      </div>

      {/* Discovered Genres */}
      {status?.genres_discovered && status.genres_discovered.length > 0 && (
        <div className="w-full max-w-lg">
          <p className="text-text-muted text-center mb-4">Genres discovered:</p>
          <div className="flex flex-wrap gap-2 justify-center">
            {status.genres_discovered.slice(0, 12).map((genre, index) => (
              <GenreTag key={genre} genre={genre} index={index} />
            ))}
            {status.genres_discovered.length > 12 && (
              <span className="text-text-muted">
                +{status.genres_discovered.length - 12} more
              </span>
            )}
          </div>
        </div>
      )}
    </main>
  );
}

export default function Processing() {
  return (
    <Suspense fallback={
      <main className="min-h-screen flex items-center justify-center">
        <VinylIcon spinning size={64} />
      </main>
    }>
      <ProcessingContent />
    </Suspense>
  );
}
```

**Step 4: Commit**

```bash
git add frontend/
git commit -m "feat: create processing page with vinyl animation and genre discovery"
```

---

### Task 13: Create Success Page

**Files:**
- Create: `frontend/src/app/success/page.tsx`
- Create: `frontend/src/components/PlaylistCard.tsx`

**Step 1: Create PlaylistCard component**

Create `frontend/src/components/PlaylistCard.tsx`:

```tsx
interface PlaylistCardProps {
  name: string;
  genre: string;
  songCount: number;
  spotifyUrl: string;
  index: number;
}

const genreColors: Record<string, string> = {
  'Rock': '#dc2626',
  'Pop': '#ec4899',
  'Hip-Hop': '#8b5cf6',
  'Electronic': '#06b6d4',
  'R&B': '#f59e0b',
  'Jazz': '#10b981',
  'Classical': '#6366f1',
  'Country': '#ca8a04',
  'Metal': '#374151',
  'Folk': '#84cc16',
  'Latin': '#f97316',
  'Blues': '#3b82f6',
  'Reggae': '#22c55e',
  'Punk': '#ef4444',
  'Indie': '#a855f7',
  'Soul': '#eab308',
  'Funk': '#d946ef',
  'World': '#14b8a6',
  'Other': '#6b7280',
};

export function PlaylistCard({ name, genre, songCount, spotifyUrl, index }: PlaylistCardProps) {
  const accentColor = genreColors[genre] || genreColors['Other'];

  return (
    <div
      className="bg-bg-card rounded-xl p-4 transform hover:-translate-y-1 hover:shadow-xl transition-all duration-200 animate-drop-in cursor-pointer group"
      style={{ animationDelay: `${index * 100}ms` }}
      onClick={() => window.open(spotifyUrl, '_blank')}
    >
      {/* Accent bar */}
      <div
        className="h-1 rounded-full mb-3"
        style={{ backgroundColor: accentColor }}
      />

      {/* Content */}
      <div className="flex items-center gap-2 mb-2">
        <span className="text-lg">♫</span>
        <h3 className="font-display text-lg text-text-cream group-hover:text-accent-orange transition-colors">
          {genre}
        </h3>
      </div>

      <p className="text-text-muted text-sm mb-3">
        {songCount} songs
      </p>

      <button className="text-sm text-accent-orange hover:text-accent-orange-hover transition-colors">
        Open in Spotify →
      </button>
    </div>
  );
}
```

**Step 2: Create Success page**

Create `frontend/src/app/success/page.tsx`:

```tsx
'use client';

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { PlaylistCard } from '@/components/PlaylistCard';
import { getOrganizeStatus } from '@/lib/api';

interface Playlist {
  name: string;
  genre: string;
  spotify_id: string;
  spotify_url: string;
  song_count: number;
}

function SuccessContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const jobId = searchParams.get('job');

  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [showConfetti, setShowConfetti] = useState(true);

  useEffect(() => {
    if (!jobId) {
      router.push('/dashboard');
      return;
    }

    getOrganizeStatus(jobId)
      .then((data) => {
        if (data.status === 'completed' && data.result) {
          setPlaylists(data.result.playlists);
        }
      })
      .catch(console.error);

    // Hide confetti after 3 seconds
    const timer = setTimeout(() => setShowConfetti(false), 3000);
    return () => clearTimeout(timer);
  }, [jobId, router]);

  const handleOpenAll = () => {
    // Open first playlist - Spotify will show related playlists
    if (playlists.length > 0) {
      window.open(playlists[0].spotify_url, '_blank');
    }
  };

  return (
    <main className="min-h-screen flex flex-col items-center px-4 py-12 relative overflow-hidden">
      {/* Confetti particles */}
      {showConfetti && (
        <div className="fixed inset-0 pointer-events-none">
          {[...Array(20)].map((_, i) => (
            <div
              key={i}
              className="absolute w-2 h-2 rounded-full animate-fade-in"
              style={{
                left: `${Math.random() * 100}%`,
                top: `-10px`,
                backgroundColor: ['#e85d04', '#f5f0e6', '#2d936c'][i % 3],
                animation: `fall ${3 + Math.random() * 2}s linear forwards`,
                animationDelay: `${Math.random() * 2}s`,
              }}
            />
          ))}
        </div>
      )}

      {/* Success Header */}
      <div className="text-center mb-12">
        <div className="text-4xl mb-4">✓</div>
        <h1 className="font-display text-4xl text-text-cream mb-2">Done!</h1>
        <p className="text-text-muted text-xl">
          {playlists.length} playlists ready to play
        </p>
      </div>

      {/* Playlist Grid */}
      <div className="w-full max-w-4xl grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-12">
        {playlists.map((playlist, index) => (
          <PlaylistCard
            key={playlist.spotify_id}
            name={playlist.name}
            genre={playlist.genre}
            songCount={playlist.song_count}
            spotifyUrl={playlist.spotify_url}
            index={index}
          />
        ))}
      </div>

      {/* Actions */}
      <div className="flex flex-col sm:flex-row gap-4 items-center">
        <Button size="lg" onClick={handleOpenAll}>
          Open All in Spotify
          <span className="ml-2">●</span>
        </Button>
        <Button
          variant="ghost"
          onClick={() => router.push('/dashboard')}
        >
          Organize Again
        </Button>
      </div>

      {/* Add falling animation */}
      <style jsx>{`
        @keyframes fall {
          to {
            transform: translateY(100vh) rotate(720deg);
            opacity: 0;
          }
        }
      `}</style>
    </main>
  );
}

export default function Success() {
  return (
    <Suspense fallback={
      <main className="min-h-screen flex items-center justify-center">
        <VinylIcon spinning size={64} />
      </main>
    }>
      <SuccessContent />
    </Suspense>
  );
}
```

**Step 3: Commit**

```bash
git add frontend/
git commit -m "feat: create success page with playlist cards and confetti animation"
```

---

## Phase 4: Polish

### Task 14: Add Error Handling

**Files:**
- Create: `frontend/src/components/ErrorBoundary.tsx`
- Modify: `frontend/src/app/dashboard/page.tsx`
- Modify: `backend/internal/spotify/oauth.go` (add token refresh)

**Step 1: Create ErrorBoundary component**

Create `frontend/src/components/ErrorBoundary.tsx`:

```tsx
'use client';

import { Component, ReactNode } from 'react';
import { Button } from './Button';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen flex flex-col items-center justify-center p-4">
          <h1 className="font-display text-2xl text-text-cream mb-4">
            Something went wrong
          </h1>
          <p className="text-text-muted mb-8">
            {this.state.error?.message || 'An unexpected error occurred'}
          </p>
          <Button onClick={() => window.location.href = '/'}>
            Go Home
          </Button>
        </div>
      );
    }

    return this.props.children;
  }
}
```

**Step 2: Add token refresh to backend**

Add to `backend/internal/spotify/oauth.go`:

```go
func (c *Config) RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.ClientID + ":" + c.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed: %d", resp.StatusCode)
	}

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}
```

**Step 3: Commit**

```bash
git add frontend/ backend/
git commit -m "feat: add error boundary and token refresh support"
```

---

### Task 15: Add Mobile Responsiveness

**Files:**
- Modify: `frontend/src/app/page.tsx`
- Modify: `frontend/src/app/dashboard/page.tsx`
- Modify: `frontend/src/app/success/page.tsx`

**Step 1: Audit and update responsive classes**

The pages already use responsive classes (`md:`, `lg:`, `sm:`). Verify they work on mobile by testing in browser dev tools.

Key responsive patterns used:
- `text-4xl md:text-5xl lg:text-6xl` - Scaling text
- `grid-cols-1 sm:grid-cols-2 lg:grid-cols-3` - Responsive grids
- `px-4` - Mobile padding
- `max-w-xl`, `max-w-4xl` - Content width constraints

**Step 2: Test on mobile viewport**

```bash
cd frontend
npm run dev
```

Open Chrome DevTools → Toggle device toolbar → Test on iPhone, iPad, etc.

**Step 3: Commit if changes needed**

```bash
git add frontend/
git commit -m "chore: verify mobile responsiveness"
```

---

### Task 16: Final Integration Test

**Step 1: Start backend**

```bash
cd backend
cp .env.example .env
# Edit .env with real Spotify credentials
go run cmd/api/main.go
```

**Step 2: Start frontend**

```bash
cd frontend
npm run dev
```

**Step 3: Test full flow**

1. Visit http://localhost:3000
2. Click "Connect with Spotify"
3. Authorize the app
4. Adjust slider, click "Organize My Library"
5. Watch progress
6. View created playlists
7. Click to open in Spotify

**Step 4: Final commit**

```bash
git add -A
git commit -m "chore: complete MVP implementation"
```

---

## Summary

**Total Tasks:** 16

**Phase 1 - Foundation (Tasks 1-4):**
- Go backend with Gin
- Next.js frontend with Tailwind
- Supabase database
- Spotify OAuth

**Phase 2 - Core Logic (Tasks 5-9):**
- Fetch liked songs
- Genre consolidation (100+ genres → 19 categories)
- Artist genre fetching
- Playlist creation
- Organize job handler

**Phase 3 - Frontend (Tasks 10-13):**
- Landing page with vinyl animation
- Dashboard with slider
- Processing page with progress
- Success page with playlist cards

**Phase 4 - Polish (Tasks 14-16):**
- Error handling
- Mobile responsiveness
- Integration testing

---

*Plan created: 2026-01-12*

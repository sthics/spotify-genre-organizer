# Exclusion Rules Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Allow users to block artists and songs from appearing in any playlist, with retroactive removal from existing playlists.

**Architecture:** New `exclusion_rules` table stores user blocklists. A filtering function runs after fetching songs but before genre scoring. API endpoints handle CRUD + impact preview. Frontend adds contextual blocking UI and settings management.

**Tech Stack:** Go/Gin (backend), PostgreSQL/Supabase (database), Next.js/TypeScript/Tailwind (frontend)

---

## Task 1: Database Migration

**Files:**
- Create: `supabase/migrations/003_exclusion_rules.sql`

**Step 1: Write the migration file**

```sql
-- Exclusion rules for blocking artists/songs from playlists
CREATE TABLE IF NOT EXISTS exclusion_rules (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id TEXT NOT NULL REFERENCES users(spotify_id) ON DELETE CASCADE,
  exclusion_type TEXT NOT NULL CHECK (exclusion_type IN ('artist', 'song')),
  spotify_id TEXT NOT NULL,
  name TEXT NOT NULL,
  scope TEXT CHECK (scope IN ('all_appearances', 'primary_only')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, exclusion_type, spotify_id)
);

-- Index for FK (PostgreSQL doesn't auto-create)
CREATE INDEX idx_exclusion_rules_user ON exclusion_rules(user_id);

-- RLS
ALTER TABLE exclusion_rules ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can manage own exclusion rules"
  ON exclusion_rules FOR ALL
  USING (user_id = current_setting('app.user_id', true));
```

**Step 2: Commit**

```bash
git add supabase/migrations/003_exclusion_rules.sql
git commit -m "feat(db): add exclusion_rules table migration"
```

---

## Task 2: Exclusion Rule Model

**Files:**
- Create: `backend/internal/models/exclusion_rule.go`

**Step 1: Write the model**

```go
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
```

**Step 2: Commit**

```bash
git add backend/internal/models/exclusion_rule.go
git commit -m "feat(models): add ExclusionRule model"
```

---

## Task 3: Database CRUD Functions

**Files:**
- Create: `backend/internal/database/exclusion_rules.go`

**Step 1: Write database functions**

```go
package database

import (
	"encoding/json"
	"time"

	"github.com/spotify-genre-organizer/backend/internal/models"
)

// GetExclusionRules returns all exclusion rules for a user
func GetExclusionRules(userID string) ([]models.ExclusionRule, error) {
	res, _, err := Client.From("exclusion_rules").
		Select("*", "", false).
		Eq("user_id", userID).
		Order("created_at", nil).
		Execute()

	if err != nil {
		return nil, err
	}

	var rules []models.ExclusionRule
	if err := json.Unmarshal(res, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

// CreateExclusionRule creates a new exclusion rule
func CreateExclusionRule(rule *models.ExclusionRule) error {
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	res, _, err := Client.From("exclusion_rules").
		Insert(rule, false, "", "", "").
		Execute()

	if err != nil {
		return err
	}

	// Parse response to get the generated ID
	var created []models.ExclusionRule
	if err := json.Unmarshal(res, &created); err != nil {
		return err
	}
	if len(created) > 0 {
		rule.ID = created[0].ID
	}

	return nil
}

// DeleteExclusionRule removes an exclusion rule by ID
func DeleteExclusionRule(userID, ruleID string) error {
	_, _, err := Client.From("exclusion_rules").
		Delete("", "").
		Eq("id", ruleID).
		Eq("user_id", userID).
		Execute()

	return err
}

// GetExclusionRuleByID fetches a single rule
func GetExclusionRuleByID(userID, ruleID string) (*models.ExclusionRule, error) {
	res, _, err := Client.From("exclusion_rules").
		Select("*", "", false).
		Eq("id", ruleID).
		Eq("user_id", userID).
		Single().
		Execute()

	if err != nil {
		return nil, err
	}

	var rule models.ExclusionRule
	if err := json.Unmarshal(res, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}
```

**Step 2: Commit**

```bash
git add backend/internal/database/exclusion_rules.go
git commit -m "feat(db): add exclusion rules CRUD functions"
```

---

## Task 4: Filtering Logic

**Files:**
- Create: `backend/internal/filters/exclusions.go`
- Create: `backend/internal/filters/exclusions_test.go`

**Step 1: Write the failing test**

```go
package filters

import (
	"testing"

	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

func TestApplyExclusions_BlocksArtistAllAppearances(t *testing.T) {
	songs := []spotify.Song{
		{ID: "1", Name: "Song A", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
		{ID: "2", Name: "Song B", Artists: []spotify.Artist{{ID: "artist2", Name: "Artist Two"}}},
		{ID: "3", Name: "Collab", Artists: []spotify.Artist{{ID: "artist2", Name: "Artist Two"}, {ID: "artist3", Name: "Artist Three"}}},
	}

	allAppearances := models.ExclusionScopeAllAppearances
	rules := []models.ExclusionRule{
		{ExclusionType: models.ExclusionTypeArtist, SpotifyID: "artist2", Scope: &allAppearances},
	}

	result := ApplyExclusions(songs, rules)

	if len(result) != 1 {
		t.Errorf("expected 1 song, got %d", len(result))
	}
	if result[0].ID != "1" {
		t.Errorf("expected song 1, got %s", result[0].ID)
	}
}

func TestApplyExclusions_BlocksArtistPrimaryOnly(t *testing.T) {
	songs := []spotify.Song{
		{ID: "1", Name: "Song A", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
		{ID: "2", Name: "Song B", Artists: []spotify.Artist{{ID: "artist2", Name: "Artist Two"}}},
		{ID: "3", Name: "Collab", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}, {ID: "artist2", Name: "Artist Two"}}},
	}

	primaryOnly := models.ExclusionScopePrimaryOnly
	rules := []models.ExclusionRule{
		{ExclusionType: models.ExclusionTypeArtist, SpotifyID: "artist2", Scope: &primaryOnly},
	}

	result := ApplyExclusions(songs, rules)

	// Song 2 blocked (artist2 is primary), Song 3 kept (artist2 is feature)
	if len(result) != 2 {
		t.Errorf("expected 2 songs, got %d", len(result))
	}
}

func TestApplyExclusions_BlocksSong(t *testing.T) {
	songs := []spotify.Song{
		{ID: "1", Name: "Song A", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
		{ID: "2", Name: "Song B", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
	}

	rules := []models.ExclusionRule{
		{ExclusionType: models.ExclusionTypeSong, SpotifyID: "1"},
	}

	result := ApplyExclusions(songs, rules)

	if len(result) != 1 {
		t.Errorf("expected 1 song, got %d", len(result))
	}
	if result[0].ID != "2" {
		t.Errorf("expected song 2, got %s", result[0].ID)
	}
}

func TestApplyExclusions_NoRules(t *testing.T) {
	songs := []spotify.Song{
		{ID: "1", Name: "Song A", Artists: []spotify.Artist{{ID: "artist1", Name: "Artist One"}}},
	}

	result := ApplyExclusions(songs, nil)

	if len(result) != 1 {
		t.Errorf("expected 1 song, got %d", len(result))
	}
}
```

**Step 2: Run test to verify it fails**

```bash
cd backend && go test -v ./internal/filters/...
```

Expected: FAIL - package doesn't exist

**Step 3: Write the implementation**

```go
package filters

import (
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

// ApplyExclusions filters out songs matching exclusion rules
func ApplyExclusions(songs []spotify.Song, rules []models.ExclusionRule) []spotify.Song {
	if len(rules) == 0 {
		return songs
	}

	// Build lookup maps for O(1) checks
	blockedSongs := make(map[string]bool)
	blockedArtistsAll := make(map[string]bool)
	blockedArtistsPrimary := make(map[string]bool)

	for _, rule := range rules {
		switch rule.ExclusionType {
		case models.ExclusionTypeSong:
			blockedSongs[rule.SpotifyID] = true
		case models.ExclusionTypeArtist:
			if rule.Scope != nil && *rule.Scope == models.ExclusionScopeAllAppearances {
				blockedArtistsAll[rule.SpotifyID] = true
			} else {
				blockedArtistsPrimary[rule.SpotifyID] = true
			}
		}
	}

	// Filter songs
	result := make([]spotify.Song, 0, len(songs))
	for _, song := range songs {
		if blockedSongs[song.ID] {
			continue
		}

		// Check primary artist (first in list)
		if len(song.Artists) > 0 {
			primaryID := song.Artists[0].ID
			if blockedArtistsPrimary[primaryID] || blockedArtistsAll[primaryID] {
				continue
			}
		}

		// Check all artists for "all_appearances" scope
		blocked := false
		for _, artist := range song.Artists {
			if blockedArtistsAll[artist.ID] {
				blocked = true
				break
			}
		}
		if blocked {
			continue
		}

		result = append(result, song)
	}

	return result
}
```

**Step 4: Run tests to verify they pass**

```bash
cd backend && go test -v ./internal/filters/...
```

Expected: PASS

**Step 5: Commit**

```bash
git add backend/internal/filters/
git commit -m "feat(filters): add ApplyExclusions function with tests"
```

---

## Task 5: Integrate Filtering into Organize Flow

**Files:**
- Modify: `backend/internal/api/handlers/organize.go`

**Step 1: Add filtering to processOrganizeJob**

In `organize.go`, add import:
```go
import (
	// ... existing imports
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/filters"
)
```

In `processOrganizeJob`, after `EnrichSongsWithGenres` and before collecting genres, add:
```go
	// Apply exclusion rules
	exclusionRules, err := database.GetExclusionRules(userID)
	if err != nil {
		log.Printf("organize job %s: warning - failed to fetch exclusion rules: %v", job.ID, err)
		// Continue without filtering - don't fail the job
	}
	songs = filters.ApplyExclusions(songs, exclusionRules)
```

**Step 2: Commit**

```bash
git add backend/internal/api/handlers/organize.go
git commit -m "feat(organize): integrate exclusion filtering"
```

---

## Task 6: Integrate Filtering into Custom Playlist Flow

**Files:**
- Modify: `backend/internal/api/handlers/custom_playlist.go`

**Step 1: Add filtering to processCustomPlaylistJob**

Add import for filters package, then in `processCustomPlaylistJob`, after building `matchingTrackIDs` and before checking if empty, add:

```go
	// Apply exclusion rules to track IDs
	exclusionRules, err := database.GetExclusionRules(userID)
	if err != nil {
		log.Printf("custom playlist job %s: warning - failed to fetch exclusion rules: %v", job.ID, err)
	}
	if len(exclusionRules) > 0 {
		matchingTrackIDs = filterTrackIDsByExclusions(matchingTrackIDs, exclusionRules, cache)
	}
```

Add helper function at bottom of file:

```go
// filterTrackIDsByExclusions removes excluded tracks
// Note: This is a simplified version - full song data not available from cache
func filterTrackIDsByExclusions(trackIDs []string, rules []models.ExclusionRule, cache *models.GenreAnalysisCache) []string {
	blockedSongs := make(map[string]bool)
	for _, rule := range rules {
		if rule.ExclusionType == models.ExclusionTypeSong {
			blockedSongs[rule.SpotifyID] = true
		}
	}

	// For artist exclusions, we'd need track->artist mapping which isn't in cache
	// This will be handled when we have full song data in organize flow
	// Custom playlist from cache only filters song-level exclusions

	result := make([]string, 0, len(trackIDs))
	for _, id := range trackIDs {
		if !blockedSongs[id] {
			result = append(result, id)
		}
	}
	return result
}
```

**Step 2: Commit**

```bash
git add backend/internal/api/handlers/custom_playlist.go
git commit -m "feat(custom-playlist): integrate song exclusion filtering"
```

---

## Task 7: Exclusions API Handler - List & Create

**Files:**
- Create: `backend/internal/api/handlers/exclusions.go`

**Step 1: Write handlers**

```go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/models"
)

type CreateExclusionRequest struct {
	Type      string  `json:"type" binding:"required,oneof=artist song"`
	SpotifyID string  `json:"spotify_id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	Scope     *string `json:"scope" binding:"omitempty,oneof=all_appearances primary_only"`
}

// ListExclusions returns all exclusion rules for the user
func ListExclusions(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	rules, err := database.GetExclusionRules(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch exclusions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

// CreateExclusion creates a new exclusion rule and returns impact preview
func CreateExclusion(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	var req CreateExclusionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build the rule
	rule := &models.ExclusionRule{
		UserID:        userID,
		ExclusionType: models.ExclusionType(req.Type),
		SpotifyID:     req.SpotifyID,
		Name:          req.Name,
	}

	if req.Scope != nil {
		scope := models.ExclusionScope(*req.Scope)
		rule.Scope = &scope
	}

	// Save to database
	if err := database.CreateExclusionRule(rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create exclusion"})
		return
	}

	// Calculate impact preview
	preview, err := calculateExclusionImpact(userID, accessToken, rule)
	if err != nil {
		// Rule created but preview failed - still return success
		c.JSON(http.StatusCreated, gin.H{
			"rule_id": rule.ID,
			"preview": nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"rule_id": rule.ID,
		"preview": preview,
	})
}
```

**Step 2: Commit**

```bash
git add backend/internal/api/handlers/exclusions.go
git commit -m "feat(api): add ListExclusions and CreateExclusion handlers"
```

---

## Task 8: Impact Preview & Apply Functions

**Files:**
- Modify: `backend/internal/api/handlers/exclusions.go`

**Step 1: Add impact calculation and apply functions**

```go
// Add to exclusions.go

type AffectedPlaylist struct {
	PlaylistID     string   `json:"playlist_id"`
	Name           string   `json:"name"`
	SongsToRemove  []string `json:"songs_to_remove"`
	SongCount      int      `json:"song_count"`
}

type ImpactPreview struct {
	TotalSongsAffected int                `json:"total_songs_affected"`
	AffectedPlaylists  []AffectedPlaylist `json:"affected_playlists"`
}

func calculateExclusionImpact(userID, accessToken string, rule *models.ExclusionRule) (*ImpactPreview, error) {
	// Get all managed playlists
	overrides, err := database.GetPlaylistOverrides(userID)
	if err != nil {
		return nil, err
	}

	preview := &ImpactPreview{
		AffectedPlaylists: []AffectedPlaylist{},
	}

	for _, override := range overrides {
		// Fetch tracks from this playlist
		tracks, err := spotify.GetPlaylistTracks(accessToken, override.PlaylistSpotifyID)
		if err != nil {
			continue // Skip playlists we can't fetch
		}

		// Find matching tracks
		var matchingTracks []string
		for _, track := range tracks {
			if matchesExclusionRule(track, rule) {
				matchingTracks = append(matchingTracks, track.ID)
			}
		}

		if len(matchingTracks) > 0 {
			name := override.Genre
			if override.CustomName != nil {
				name = *override.CustomName
			}
			preview.AffectedPlaylists = append(preview.AffectedPlaylists, AffectedPlaylist{
				PlaylistID:    override.PlaylistSpotifyID,
				Name:          name,
				SongsToRemove: matchingTracks,
				SongCount:     len(matchingTracks),
			})
			preview.TotalSongsAffected += len(matchingTracks)
		}
	}

	return preview, nil
}

func matchesExclusionRule(track spotify.Song, rule *models.ExclusionRule) bool {
	if rule.ExclusionType == models.ExclusionTypeSong {
		return track.ID == rule.SpotifyID
	}

	// Artist rule
	if rule.Scope != nil && *rule.Scope == models.ExclusionScopeAllAppearances {
		for _, artist := range track.Artists {
			if artist.ID == rule.SpotifyID {
				return true
			}
		}
	} else {
		// Primary only
		if len(track.Artists) > 0 && track.Artists[0].ID == rule.SpotifyID {
			return true
		}
	}

	return false
}

// ApplyExclusion removes matching tracks from playlists
func ApplyExclusion(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	ruleID := c.Param("id")

	// Get the rule
	rule, err := database.GetExclusionRuleByID(userID, ruleID)
	if err != nil || rule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exclusion rule not found"})
		return
	}

	// Get impact
	preview, err := calculateExclusionImpact(userID, accessToken, rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate impact"})
		return
	}

	// Remove tracks from each playlist
	playlistsUpdated := 0
	songsRemoved := 0

	for _, affected := range preview.AffectedPlaylists {
		err := spotify.RemoveTracksFromPlaylist(accessToken, affected.PlaylistID, affected.SongsToRemove)
		if err != nil {
			continue
		}
		playlistsUpdated++
		songsRemoved += len(affected.SongsToRemove)
	}

	c.JSON(http.StatusOK, gin.H{
		"playlists_updated": playlistsUpdated,
		"songs_removed":     songsRemoved,
	})
}

// DeleteExclusion removes an exclusion rule
func DeleteExclusion(c *gin.Context) {
	userID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	ruleID := c.Param("id")

	if err := database.DeleteExclusionRule(userID, ruleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete exclusion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}
```

**Step 2: Commit**

```bash
git add backend/internal/api/handlers/exclusions.go
git commit -m "feat(api): add impact preview and apply/delete exclusion handlers"
```

---

## Task 9: Add Spotify Playlist Tracks Functions

**Files:**
- Modify: `backend/internal/spotify/playlists.go`

**Step 1: Add GetPlaylistTracks and RemoveTracksFromPlaylist**

```go
// Add to playlists.go

// GetPlaylistTracks fetches all tracks from a playlist
func GetPlaylistTracks(accessToken, playlistID string) ([]Song, error) {
	var allTracks []Song
	limit := 100
	offset := 0

	for {
		url := fmt.Sprintf("%s/playlists/%s/tracks?limit=%d&offset=%d&fields=items(track(id,name,artists(id,name))),total",
			APIURL, playlistID, limit, offset)

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
			return nil, fmt.Errorf("failed to get playlist tracks: %d", resp.StatusCode)
		}

		var result struct {
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
			Total int `json:"total"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		for _, item := range result.Items {
			if item.Track.ID == "" {
				continue // Skip local/unavailable tracks
			}
			artists := make([]Artist, len(item.Track.Artists))
			for i, a := range item.Track.Artists {
				artists[i] = Artist{ID: a.ID, Name: a.Name}
			}
			allTracks = append(allTracks, Song{
				ID:      item.Track.ID,
				Name:    item.Track.Name,
				Artists: artists,
			})
		}

		if len(result.Items) < limit || len(allTracks) >= result.Total {
			break
		}
		offset += limit
	}

	return allTracks, nil
}

// RemoveTracksFromPlaylist removes tracks from a playlist
func RemoveTracksFromPlaylist(accessToken, playlistID string, trackIDs []string) error {
	// Spotify allows max 100 tracks per request
	for i := 0; i < len(trackIDs); i += 100 {
		end := i + 100
		if end > len(trackIDs) {
			end = len(trackIDs)
		}
		chunk := trackIDs[i:end]

		tracks := make([]map[string]string, len(chunk))
		for j, id := range chunk {
			tracks[j] = map[string]string{"uri": "spotify:track:" + id}
		}

		body, _ := json.Marshal(map[string]any{"tracks": tracks})

		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/playlists/%s/tracks", APIURL, playlistID),
			bytes.NewReader(body))
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

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to remove tracks: %d", resp.StatusCode)
		}
	}

	return nil
}
```

**Step 2: Add bytes import to playlists.go**

```go
import (
	"bytes"
	// ... other imports
)
```

**Step 3: Commit**

```bash
git add backend/internal/spotify/playlists.go
git commit -m "feat(spotify): add GetPlaylistTracks and RemoveTracksFromPlaylist"
```

---

## Task 10: Register API Routes

**Files:**
- Modify: `backend/internal/api/routes.go`

**Step 1: Add exclusion routes**

In `SetupRoutes`, add after the playlists group:

```go
		exclusions := api.Group("/exclusions")
		{
			exclusions.GET("", handlers.ListExclusions)
			exclusions.POST("", handlers.CreateExclusion)
			exclusions.POST("/:id/apply", handlers.ApplyExclusion)
			exclusions.DELETE("/:id", handlers.DeleteExclusion)
		}
```

**Step 2: Commit**

```bash
git add backend/internal/api/routes.go
git commit -m "feat(api): register exclusion routes"
```

---

## Task 11: Frontend API Client

**Files:**
- Modify: `frontend/src/lib/api.ts`

**Step 1: Add exclusion types and functions**

```typescript
// Add to api.ts

// Exclusion Types
export interface ExclusionRule {
  id: string;
  user_id: string;
  exclusion_type: 'artist' | 'song';
  spotify_id: string;
  name: string;
  scope?: 'all_appearances' | 'primary_only';
  created_at: string;
}

export interface AffectedPlaylist {
  playlist_id: string;
  name: string;
  songs_to_remove: string[];
  song_count: number;
}

export interface ImpactPreview {
  total_songs_affected: number;
  affected_playlists: AffectedPlaylist[];
}

export interface CreateExclusionResponse {
  rule_id: string;
  preview: ImpactPreview | null;
}

export interface ApplyExclusionResponse {
  playlists_updated: number;
  songs_removed: number;
}

// Exclusion API Functions
export async function getExclusions(): Promise<{ rules: ExclusionRule[] }> {
  const response = await fetch(`${API_URL}/api/exclusions`, {
    credentials: 'include',
  });
  handleApiResponse(response);
  return response.json();
}

export async function createExclusion(
  type: 'artist' | 'song',
  spotifyId: string,
  name: string,
  scope?: 'all_appearances' | 'primary_only'
): Promise<CreateExclusionResponse> {
  const response = await fetch(`${API_URL}/api/exclusions`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      type,
      spotify_id: spotifyId,
      name,
      scope,
    }),
  });
  handleApiResponse(response);
  return response.json();
}

export async function applyExclusion(ruleId: string): Promise<ApplyExclusionResponse> {
  const response = await fetch(`${API_URL}/api/exclusions/${ruleId}/apply`, {
    method: 'POST',
    credentials: 'include',
  });
  handleApiResponse(response);
  return response.json();
}

export async function deleteExclusion(ruleId: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/exclusions/${ruleId}`, {
    method: 'DELETE',
    credentials: 'include',
  });
  handleApiResponse(response);
}
```

**Step 2: Commit**

```bash
git add frontend/src/lib/api.ts
git commit -m "feat(frontend): add exclusion API client functions"
```

---

## Task 12: BlockButton Component

**Files:**
- Create: `frontend/src/components/blocking/BlockButton.tsx`

**Step 1: Create the component**

```tsx
'use client';

import { useState } from 'react';

interface BlockButtonProps {
  type: 'artist' | 'song';
  spotifyId: string;
  name: string;
  imageUrl?: string;
  onBlock: () => void;
}

export function BlockButton({ onBlock }: BlockButtonProps) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <button
      onClick={(e) => {
        e.stopPropagation();
        onBlock();
      }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      className={`p-1.5 rounded-full transition-all duration-200 ${
        isHovered
          ? 'bg-red-500/20 text-red-400'
          : 'text-text-muted hover:text-red-400'
      }`}
      title="Block"
    >
      <svg
        className="w-4 h-4"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
        />
      </svg>
    </button>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/src/components/blocking/BlockButton.tsx
git commit -m "feat(ui): add BlockButton component"
```

---

## Task 13: BlockArtistModal Component

**Files:**
- Create: `frontend/src/components/blocking/BlockArtistModal.tsx`

**Step 1: Create the modal**

```tsx
'use client';

import { useState } from 'react';
import { Button } from '@/components/Button';
import { createExclusion, applyExclusion, CreateExclusionResponse } from '@/lib/api';

interface BlockArtistModalProps {
  artist: { id: string; name: string; imageUrl?: string };
  songCount?: number;
  isOpen: boolean;
  onClose: () => void;
  onBlocked: () => void;
}

type Step = 'scope' | 'preview' | 'loading';

export function BlockArtistModal({
  artist,
  songCount,
  isOpen,
  onClose,
  onBlocked,
}: BlockArtistModalProps) {
  const [step, setStep] = useState<Step>('scope');
  const [scope, setScope] = useState<'all_appearances' | 'primary_only'>('all_appearances');
  const [preview, setPreview] = useState<CreateExclusionResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleBlock = async () => {
    setStep('loading');
    setError(null);

    try {
      const result = await createExclusion('artist', artist.id, artist.name, scope);
      setPreview(result);

      if (result.preview && result.preview.total_songs_affected > 0) {
        setStep('preview');
      } else {
        // No songs affected, just close
        onBlocked();
        onClose();
      }
    } catch (err) {
      setError('Failed to create block rule. Please try again.');
      setStep('scope');
    }
  };

  const handleApply = async () => {
    if (!preview) return;

    setStep('loading');
    try {
      await applyExclusion(preview.rule_id);
      onBlocked();
      onClose();
    } catch (err) {
      setError('Failed to remove songs. Please try again.');
      setStep('preview');
    }
  };

  const handleSkip = () => {
    onBlocked();
    onClose();
  };

  const resetAndClose = () => {
    setStep('scope');
    setScope('all_appearances');
    setPreview(null);
    setError(null);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4">
      <div className="bg-bg-card rounded-2xl max-w-md w-full p-6 shadow-2xl animate-scale-in">
        {/* Header */}
        <div className="flex justify-between items-start mb-6">
          <div className="flex items-center gap-4">
            {artist.imageUrl ? (
              <img
                src={artist.imageUrl}
                alt={artist.name}
                className="w-12 h-12 rounded-full object-cover"
              />
            ) : (
              <div className="w-12 h-12 rounded-full bg-bg-dark flex items-center justify-center">
                <svg className="w-6 h-6 text-text-muted" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 14.016q2.906 0 4.945 2.039T19.031 21H4.969q0-2.906 2.039-4.945T12 14.016zm0-1.032q-1.641 0-2.813-1.172T8.015 9t1.172-2.813T12 5.015t2.813 1.172T15.985 9t-1.172 2.813T12 12.984z"/>
                </svg>
              </div>
            )}
            <div>
              <h2 className="font-display text-xl text-text-cream">
                {step === 'preview' ? `Blocking "${artist.name}"` : `Block "${artist.name}"?`}
              </h2>
              {songCount && step === 'scope' && (
                <p className="text-sm text-text-muted">{songCount} songs in your library</p>
              )}
            </div>
          </div>
          <button onClick={resetAndClose} className="text-text-muted hover:text-text-cream">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400 text-sm">
            {error}
          </div>
        )}

        {step === 'scope' && (
          <>
            <div className="space-y-3 mb-6">
              <button
                onClick={() => setScope('all_appearances')}
                className={`w-full p-4 rounded-xl border-2 text-left transition-all ${
                  scope === 'all_appearances'
                    ? 'border-accent-orange bg-accent-orange/10'
                    : 'border-bg-dark hover:border-text-muted'
                }`}
              >
                <div className="font-medium text-text-cream">All songs with this artist</div>
                <div className="text-sm text-text-muted">Including features & collabs</div>
              </button>
              <button
                onClick={() => setScope('primary_only')}
                className={`w-full p-4 rounded-xl border-2 text-left transition-all ${
                  scope === 'primary_only'
                    ? 'border-accent-orange bg-accent-orange/10'
                    : 'border-bg-dark hover:border-text-muted'
                }`}
              >
                <div className="font-medium text-text-cream">Only as main artist</div>
                <div className="text-sm text-text-muted">Keep their features</div>
              </button>
            </div>

            <div className="flex gap-3 justify-end">
              <Button variant="ghost" onClick={resetAndClose}>Cancel</Button>
              <Button
                onClick={handleBlock}
                className="bg-red-600 hover:bg-red-500"
              >
                Block Artist
              </Button>
            </div>
          </>
        )}

        {step === 'preview' && preview?.preview && (
          <>
            <p className="text-text-cream mb-4">
              This will remove {preview.preview.total_songs_affected} songs from{' '}
              {preview.preview.affected_playlists.length} playlists:
            </p>
            <div className="max-h-48 overflow-y-auto mb-4 space-y-2">
              {preview.preview.affected_playlists.map((p) => (
                <div key={p.playlist_id} className="flex justify-between items-center p-3 bg-bg-dark rounded-lg">
                  <span className="text-text-cream">{p.name}</span>
                  <span className="text-text-muted text-sm">{p.song_count} songs</span>
                </div>
              ))}
            </div>
            <p className="text-sm text-text-muted mb-6">
              Rule saved. Future playlists will also exclude this artist.
            </p>
            <div className="flex gap-3 justify-end">
              <Button variant="secondary" onClick={handleSkip}>Keep in playlists</Button>
              <Button onClick={handleApply}>Remove now</Button>
            </div>
          </>
        )}

        {step === 'loading' && (
          <div className="flex items-center justify-center py-8">
            <div className="w-8 h-8 border-2 border-accent-orange border-t-transparent rounded-full animate-spin" />
          </div>
        )}
      </div>
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/src/components/blocking/BlockArtistModal.tsx
git commit -m "feat(ui): add BlockArtistModal component"
```

---

## Task 14: BlockSongModal Component

**Files:**
- Create: `frontend/src/components/blocking/BlockSongModal.tsx`

**Step 1: Create the modal**

```tsx
'use client';

import { useState } from 'react';
import { Button } from '@/components/Button';
import { createExclusion, applyExclusion, CreateExclusionResponse } from '@/lib/api';

interface BlockSongModalProps {
  song: { id: string; name: string; artistName: string };
  isOpen: boolean;
  onClose: () => void;
  onBlocked: () => void;
}

type Step = 'confirm' | 'preview' | 'loading';

export function BlockSongModal({
  song,
  isOpen,
  onClose,
  onBlocked,
}: BlockSongModalProps) {
  const [step, setStep] = useState<Step>('confirm');
  const [preview, setPreview] = useState<CreateExclusionResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleBlock = async () => {
    setStep('loading');
    setError(null);

    try {
      const result = await createExclusion('song', song.id, song.name);
      setPreview(result);

      if (result.preview && result.preview.total_songs_affected > 0) {
        setStep('preview');
      } else {
        onBlocked();
        onClose();
      }
    } catch (err) {
      setError('Failed to block song. Please try again.');
      setStep('confirm');
    }
  };

  const handleApply = async () => {
    if (!preview) return;

    setStep('loading');
    try {
      await applyExclusion(preview.rule_id);
      onBlocked();
      onClose();
    } catch (err) {
      setError('Failed to remove song. Please try again.');
      setStep('preview');
    }
  };

  const handleSkip = () => {
    onBlocked();
    onClose();
  };

  const resetAndClose = () => {
    setStep('confirm');
    setPreview(null);
    setError(null);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4">
      <div className="bg-bg-card rounded-2xl max-w-md w-full p-6 shadow-2xl animate-scale-in">
        <div className="flex justify-between items-start mb-6">
          <div>
            <h2 className="font-display text-xl text-text-cream">Block "{song.name}"?</h2>
            <p className="text-sm text-text-muted">{song.artistName}</p>
          </div>
          <button onClick={resetAndClose} className="text-text-muted hover:text-text-cream">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400 text-sm">
            {error}
          </div>
        )}

        {step === 'confirm' && (
          <>
            <p className="text-text-muted mb-6">
              This song will be removed from all playlists and never added again.
            </p>
            <div className="flex gap-3 justify-end">
              <Button variant="ghost" onClick={resetAndClose}>Cancel</Button>
              <Button onClick={handleBlock} className="bg-red-600 hover:bg-red-500">
                Block Song
              </Button>
            </div>
          </>
        )}

        {step === 'preview' && preview?.preview && (
          <>
            <p className="text-text-cream mb-4">
              This will remove the song from {preview.preview.affected_playlists.length} playlists.
            </p>
            <div className="flex gap-3 justify-end">
              <Button variant="secondary" onClick={handleSkip}>Keep in playlists</Button>
              <Button onClick={handleApply}>Remove now</Button>
            </div>
          </>
        )}

        {step === 'loading' && (
          <div className="flex items-center justify-center py-8">
            <div className="w-8 h-8 border-2 border-accent-orange border-t-transparent rounded-full animate-spin" />
          </div>
        )}
      </div>
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/src/components/blocking/BlockSongModal.tsx
git commit -m "feat(ui): add BlockSongModal component"
```

---

## Task 15: Exclusions Settings Page

**Files:**
- Create: `frontend/src/app/settings/exclusions/page.tsx`

**Step 1: Create the page**

```tsx
'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/Button';
import { getExclusions, deleteExclusion, ExclusionRule } from '@/lib/api';
import { useToast } from '@/contexts/ToastContext';

export default function ExclusionsSettings() {
  const router = useRouter();
  const { showToast } = useToast();
  const [loading, setLoading] = useState(true);
  const [rules, setRules] = useState<ExclusionRule[]>([]);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  useEffect(() => {
    loadRules();
  }, []);

  const loadRules = async () => {
    try {
      const data = await getExclusions();
      setRules(data.rules || []);
    } catch (err) {
      showToast('Failed to load exclusions', 'error');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (rule: ExclusionRule) => {
    if (!confirm(`Unblock "${rule.name}"? Their songs can appear in your playlists again.`)) {
      return;
    }

    setDeletingId(rule.id);
    try {
      await deleteExclusion(rule.id);
      setRules(rules.filter((r) => r.id !== rule.id));
      showToast(`${rule.name} unblocked`, 'success');
    } catch (err) {
      showToast('Failed to unblock', 'error');
    } finally {
      setDeletingId(null);
    }
  };

  const artistRules = rules.filter((r) => r.exclusion_type === 'artist');
  const songRules = rules.filter((r) => r.exclusion_type === 'song');

  if (loading) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <div className="text-text-cream animate-pulse">Loading exclusions...</div>
      </main>
    );
  }

  return (
    <main className="min-h-screen flex flex-col items-center px-4 py-12">
      <div className="w-full max-w-2xl">
        <button
          onClick={() => router.push('/settings')}
          className="text-text-muted hover:text-text-cream transition-colors mb-8"
        >
          &larr; Back to Settings
        </button>

        <h1 className="font-display text-3xl text-text-cream mb-8">Blocked Artists & Songs</h1>

        {/* Artists Section */}
        <section className="mb-8">
          <h2 className="font-display text-xl text-text-cream mb-4">Blocked Artists</h2>
          {artistRules.length === 0 ? (
            <div className="bg-bg-card rounded-xl p-6 text-center text-text-muted">
              No blocked artists yet. Block artists while browsing your library.
            </div>
          ) : (
            <div className="bg-bg-card rounded-xl overflow-hidden divide-y divide-bg-dark">
              {artistRules.map((rule) => (
                <div key={rule.id} className="flex items-center justify-between p-4">
                  <div>
                    <div className="font-medium text-text-cream">{rule.name}</div>
                    <div className="text-sm text-text-muted">
                      {rule.scope === 'all_appearances' ? 'All appearances' : 'Main artist only'}
                    </div>
                  </div>
                  <button
                    onClick={() => handleDelete(rule)}
                    disabled={deletingId === rule.id}
                    className="p-2 text-text-muted hover:text-red-400 transition-colors disabled:opacity-50"
                  >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Songs Section */}
        <section>
          <h2 className="font-display text-xl text-text-cream mb-4">Blocked Songs</h2>
          {songRules.length === 0 ? (
            <div className="bg-bg-card rounded-xl p-6 text-center text-text-muted">
              No blocked songs yet. Block songs while browsing your library or playlists.
            </div>
          ) : (
            <div className="bg-bg-card rounded-xl overflow-hidden divide-y divide-bg-dark">
              {songRules.map((rule) => (
                <div key={rule.id} className="flex items-center justify-between p-4">
                  <div className="font-medium text-text-cream">{rule.name}</div>
                  <button
                    onClick={() => handleDelete(rule)}
                    disabled={deletingId === rule.id}
                    className="p-2 text-text-muted hover:text-red-400 transition-colors disabled:opacity-50"
                  >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          )}
        </section>

        {/* Tip */}
        <div className="mt-8 p-4 bg-bg-card/50 rounded-xl text-sm text-text-muted">
          <span className="mr-2">💡</span>
          Tip: You can block artists and songs directly while browsing your library or viewing playlists.
        </div>
      </div>
    </main>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/src/app/settings/exclusions/page.tsx
git commit -m "feat(ui): add exclusions settings page"
```

---

## Task 16: Add Link to Exclusions from Settings

**Files:**
- Modify: `frontend/src/app/settings/page.tsx`

**Step 1: Add link to exclusions page**

Add after the existing settings card, before the closing `</div>`:

```tsx
        {/* Exclusions Link */}
        <div className="mt-6">
          <button
            onClick={() => router.push('/settings/exclusions')}
            className="w-full bg-bg-card rounded-xl p-6 text-left hover:bg-bg-card/80 transition-colors group"
          >
            <div className="flex items-center justify-between">
              <div>
                <h2 className="font-display text-lg text-text-cream">Blocked Artists & Songs</h2>
                <p className="text-sm text-text-muted">Manage your blocklist</p>
              </div>
              <svg
                className="w-5 h-5 text-text-muted group-hover:text-text-cream transition-colors"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </button>
        </div>
```

**Step 2: Commit**

```bash
git add frontend/src/app/settings/page.tsx
git commit -m "feat(ui): add link to exclusions from settings page"
```

---

## Task 17: Export Blocking Components

**Files:**
- Create: `frontend/src/components/blocking/index.ts`

**Step 1: Create index file**

```typescript
export { BlockButton } from './BlockButton';
export { BlockArtistModal } from './BlockArtistModal';
export { BlockSongModal } from './BlockSongModal';
```

**Step 2: Commit**

```bash
git add frontend/src/components/blocking/index.ts
git commit -m "feat(ui): export blocking components"
```

---

## Task 18: Final Integration Test

**Step 1: Run backend tests**

```bash
cd backend && go test -v ./...
```

Expected: All tests pass

**Step 2: Run frontend lint**

```bash
cd frontend && npm run lint
```

Expected: Pass (or existing issues only)

**Step 3: Commit any fixes**

If any issues found, fix and commit.

---

## Task 19: Final Commit & Summary

**Step 1: Review all changes**

```bash
git log --oneline feature/exclusion-rules ^main
```

**Step 2: Ensure working tree is clean**

```bash
git status
```

If any uncommitted changes, commit them.

---

## Summary

This plan implements:

1. **Database:** New `exclusion_rules` table with proper indexes and RLS
2. **Backend:**
   - Model for exclusion rules
   - CRUD database functions
   - `ApplyExclusions` filter function with tests
   - Integration into organize and custom playlist flows
   - API handlers for list, create, apply, delete
   - Spotify helper functions for playlist tracks
3. **Frontend:**
   - API client functions
   - BlockButton, BlockArtistModal, BlockSongModal components
   - Exclusions settings page
   - Link from main settings

Total: 19 tasks with frequent commits.

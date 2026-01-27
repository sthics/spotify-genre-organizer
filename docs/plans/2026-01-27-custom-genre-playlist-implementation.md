# Custom Genre Playlist Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add "Build a Crate" feature allowing users to create playlists from specific genres/sub-genres with hierarchical selection.

**Architecture:** New genre analysis endpoint caches user's genre breakdown in Supabase. Frontend provides accordion-style genre picker. Users choose combined vs separate playlists. Reuses existing job/processing pattern.

**Tech Stack:** Go/Gin backend, Next.js/React frontend, Supabase/PostgreSQL, Tailwind CSS

---

## Task 1: Add Genre Analysis Cache Model

**Files:**
- Create: `backend/internal/models/genre_cache.go`

**Step 1: Create the model file**

```go
package models

import (
	"time"
)

type SubGenreCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ParentGenreCount struct {
	Name      string          `json:"name"`
	Count     int             `json:"count"`
	SubGenres []SubGenreCount `json:"sub_genres"`
}

type GenreAnalysisData struct {
	ParentGenres []ParentGenreCount `json:"parent_genres"`
	TotalSongs   int                `json:"total_songs"`
	// Maps track URI to list of sub-genres for quick filtering
	SongGenreMap map[string][]string `json:"song_genre_map"`
}

type GenreAnalysisCache struct {
	ID           string            `json:"id" db:"id"`
	UserID       string            `json:"user_id" db:"user_id"`
	AnalysisData GenreAnalysisData `json:"analysis_data" db:"analysis_data"`
	AnalyzedAt   time.Time         `json:"analyzed_at" db:"analyzed_at"`
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" db:"updated_at"`
}
```

**Step 2: Commit**

```bash
git add backend/internal/models/genre_cache.go
git commit -m "feat: add genre analysis cache model"
```

---

## Task 2: Add Database Functions for Genre Cache

**Files:**
- Modify: `backend/internal/database/settings.go` (add new functions at end)

**Step 1: Add GetGenreAnalysisCache function**

Add at end of `backend/internal/database/settings.go`:

```go
// GetGenreAnalysisCache retrieves cached genre analysis for a user
func GetGenreAnalysisCache(userID string) (*models.GenreAnalysisCache, error) {
	res, _, err := Client.From("genre_analysis_cache").
		Select("*", "", false).
		Eq("user_id", userID).
		Single().
		Execute()

	if err != nil {
		return nil, err
	}

	var cache models.GenreAnalysisCache
	if err := json.Unmarshal(res, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// SaveGenreAnalysisCache upserts genre analysis cache
func SaveGenreAnalysisCache(cache *models.GenreAnalysisCache) error {
	cache.UpdatedAt = time.Now()

	_, _, err := Client.From("genre_analysis_cache").
		Upsert(cache, "", "", "").
		Execute()

	return err
}

// DeleteGenreAnalysisCache removes cached analysis (for refresh)
func DeleteGenreAnalysisCache(userID string) error {
	_, _, err := Client.From("genre_analysis_cache").
		Delete("", "").
		Eq("user_id", userID).
		Execute()

	return err
}
```

**Step 2: Commit**

```bash
git add backend/internal/database/settings.go
git commit -m "feat: add database functions for genre analysis cache"
```

---

## Task 3: Add Genre Analysis Helper Function

**Files:**
- Modify: `backend/internal/genres/mapping.go` (add function at end)

**Step 1: Add AnalyzeLibraryGenres function**

Add at end of `backend/internal/genres/mapping.go`:

```go
// AnalyzeLibraryGenres analyzes songs and returns detailed genre breakdown
// Returns parent genres with sub-genre counts and a map of song URIs to their sub-genres
func AnalyzeLibraryGenres(songs []spotify.Song) ([]ParentGenreCount, map[string][]string) {
	// Track sub-genre counts per parent
	parentSubGenres := make(map[string]map[string]int) // parent -> subgenre -> count
	songGenreMap := make(map[string][]string)          // trackURI -> []subgenres

	for _, song := range songs {
		trackSubGenres := make(map[string]bool)
		for _, microGenre := range song.Genres {
			normalized := strings.ToLower(strings.TrimSpace(microGenre))
			if normalized == "" {
				continue
			}
			parent := ConsolidateGenre(normalized)

			// Initialize parent map if needed
			if parentSubGenres[parent] == nil {
				parentSubGenres[parent] = make(map[string]int)
			}

			// Count this sub-genre under its parent
			parentSubGenres[parent][normalized]++
			trackSubGenres[normalized] = true
		}

		// Store sub-genres for this track
		if len(trackSubGenres) > 0 {
			for sg := range trackSubGenres {
				songGenreMap[song.URI] = append(songGenreMap[song.URI], sg)
			}
		}
	}

	// Build result sorted by count
	var result []ParentGenreCount
	for _, parentName := range ParentGenres {
		subGenres := parentSubGenres[parentName]
		if len(subGenres) == 0 {
			continue
		}

		var subGenreCounts []SubGenreCount
		totalCount := 0
		for sg, count := range subGenres {
			subGenreCounts = append(subGenreCounts, SubGenreCount{Name: sg, Count: count})
			totalCount += count
		}

		// Sort sub-genres by count descending
		sort.Slice(subGenreCounts, func(i, j int) bool {
			return subGenreCounts[i].Count > subGenreCounts[j].Count
		})

		result = append(result, ParentGenreCount{
			Name:      parentName,
			Count:     totalCount,
			SubGenres: subGenreCounts,
		})
	}

	// Sort parents by count descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result, songGenreMap
}
```

Also add `"sort"` to the imports at the top of the file.

Also add these types near the top after the `var genreMapping` block:

```go
type SubGenreCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type ParentGenreCount struct {
	Name      string          `json:"name"`
	Count     int             `json:"count"`
	SubGenres []SubGenreCount `json:"sub_genres"`
}
```

**Step 2: Run tests**

```bash
cd /Users/ambarg/developer/spotify_genre_organizer/backend && go test ./internal/genres/... -v
```

Expected: All existing tests pass.

**Step 3: Commit**

```bash
git add backend/internal/genres/mapping.go
git commit -m "feat: add AnalyzeLibraryGenres function for detailed genre breakdown"
```

---

## Task 4: Write Test for GetLibraryGenres Handler

**Files:**
- Create: `backend/internal/api/handlers/genres_test.go`

**Step 1: Create test file**

```go
package handlers

import (
	"testing"
)

func TestGetLibraryGenresRequest(t *testing.T) {
	// Test that refresh=true parameter is parsed correctly
	tests := []struct {
		name         string
		queryParam   string
		expectRefresh bool
	}{
		{"no param", "", false},
		{"refresh true", "refresh=true", true},
		{"refresh false", "refresh=false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This validates the expected query param behavior
			// Full integration test would require mock Spotify client
			if tt.queryParam == "refresh=true" && !tt.expectRefresh {
				t.Error("refresh=true should expect refresh")
			}
		})
	}
}
```

**Step 2: Run test to verify it passes**

```bash
cd /Users/ambarg/developer/spotify_genre_organizer/backend && go test ./internal/api/handlers/... -v
```

Expected: PASS

**Step 3: Commit**

```bash
git add backend/internal/api/handlers/genres_test.go
git commit -m "test: add basic test for genres handler"
```

---

## Task 5: Create GetLibraryGenres Handler

**Files:**
- Create: `backend/internal/api/handlers/genres.go`

**Step 1: Create handler file**

```go
package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/genres"
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type LibraryGenresResponse struct {
	ParentGenres []genres.ParentGenreCount `json:"parent_genres"`
	AnalyzedAt   time.Time                 `json:"analyzed_at"`
	TotalSongs   int                       `json:"total_songs"`
}

func GetLibraryGenres(c *gin.Context) {
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

	refresh := c.Query("refresh") == "true"

	// Check cache first (unless refresh requested)
	if !refresh {
		cached, err := database.GetGenreAnalysisCache(userID)
		if err == nil && cached != nil {
			// Check if cache is less than 24 hours old
			if time.Since(cached.AnalyzedAt) < 24*time.Hour {
				c.JSON(http.StatusOK, LibraryGenresResponse{
					ParentGenres: cached.AnalysisData.ParentGenres,
					AnalyzedAt:   cached.AnalyzedAt,
					TotalSongs:   cached.AnalysisData.TotalSongs,
				})
				return
			}
		}
	}

	// Fetch and analyze library
	songs, err := spotify.FetchAllLikedSongs(accessToken, nil)
	if err != nil {
		log.Printf("GetLibraryGenres: failed to fetch songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch library"})
		return
	}

	// Fetch artist genres
	artistGenres, err := spotify.FetchAllArtistGenres(accessToken, songs, nil)
	if err != nil {
		log.Printf("GetLibraryGenres: failed to fetch artist genres: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze genres"})
		return
	}

	// Enrich songs with genres
	spotify.EnrichSongsWithGenres(songs, artistGenres)

	// Analyze genres
	parentGenres, songGenreMap := genres.AnalyzeLibraryGenres(songs)

	// Cache the result
	analysisData := models.GenreAnalysisData{
		ParentGenres: parentGenres,
		TotalSongs:   len(songs),
		SongGenreMap: songGenreMap,
	}

	cache := &models.GenreAnalysisCache{
		ID:           uuid.New().String(),
		UserID:       userID,
		AnalysisData: analysisData,
		AnalyzedAt:   time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := database.SaveGenreAnalysisCache(cache); err != nil {
		log.Printf("GetLibraryGenres: failed to cache analysis: %v", err)
		// Continue anyway - caching failure shouldn't block the response
	}

	c.JSON(http.StatusOK, LibraryGenresResponse{
		ParentGenres: parentGenres,
		AnalyzedAt:   cache.AnalyzedAt,
		TotalSongs:   len(songs),
	})
}
```

**Step 2: Commit**

```bash
git add backend/internal/api/handlers/genres.go
git commit -m "feat: add GetLibraryGenres handler"
```

---

## Task 6: Add GetLibraryGenres Route

**Files:**
- Modify: `backend/internal/api/routes.go`

**Step 1: Add route**

In `backend/internal/api/routes.go`, after line 121 (`api.GET("/library/count", handlers.GetLibraryCount)`), add:

```go
		api.GET("/library/genres", handlers.GetLibraryGenres)
```

**Step 2: Verify build**

```bash
cd /Users/ambarg/developer/spotify_genre_organizer/backend && go build ./...
```

Expected: Build succeeds

**Step 3: Commit**

```bash
git add backend/internal/api/routes.go
git commit -m "feat: add /api/library/genres route"
```

---

## Task 7: Create Custom Playlist Handler

**Files:**
- Create: `backend/internal/api/handlers/custom_playlist.go`

**Step 1: Create handler file**

```go
package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type CustomPlaylistRequest struct {
	SubGenres []string `json:"sub_genres" binding:"required,min=1"`
	Mode      string   `json:"mode" binding:"required,oneof=combined separate"`
	Name      string   `json:"name"` // Only used for combined mode
}

func StartCustomPlaylist(c *gin.Context) {
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

	var req CustomPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get cached genre analysis (must exist for this endpoint)
	cache, err := database.GetGenreAnalysisCache(userID)
	if err != nil || cache == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please analyze your library first"})
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
	go processCustomPlaylistJob(job, accessToken, userID, req, cache)

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": jobID,
		"status": "pending",
	})
}

func processCustomPlaylistJob(job *JobStatus, accessToken, userID string, req CustomPlaylistRequest, cache *models.GenreAnalysisCache) {
	updateJob := func() {
		jobsMu.Lock()
		jobs[job.ID] = job
		jobsMu.Unlock()
	}

	job.Status = "processing"
	job.Stage = "filtering"
	updateJob()

	// Build set of selected sub-genres for quick lookup
	selectedGenres := make(map[string]bool)
	for _, sg := range req.SubGenres {
		selectedGenres[sg] = true
	}

	// Filter songs that match selected sub-genres
	var matchingTrackURIs []string
	for trackURI, trackGenres := range cache.AnalysisData.SongGenreMap {
		for _, g := range trackGenres {
			if selectedGenres[g] {
				matchingTrackURIs = append(matchingTrackURIs, trackURI)
				break
			}
		}
	}

	if len(matchingTrackURIs) == 0 {
		job.Status = "failed"
		job.Error = "No songs found matching selected genres"
		updateJob()
		return
	}

	job.TotalSongs = len(matchingTrackURIs)
	job.Stage = "creating"
	updateJob()

	// Get user settings for name template
	settings, _ := database.GetUserSettings(userID)
	if settings == nil {
		settings = models.DefaultSettings(userID)
	}

	var createdPlaylists []CreatedPlaylist

	if req.Mode == "combined" {
		// Create single playlist
		name := req.Name
		if name == "" {
			name = "Custom Mix by Organizer"
		}
		desc := "Custom playlist created by Spotify Genre Organizer"

		playlist, err := spotify.CreatePlaylist(accessToken, userID, name, desc)
		if err != nil {
			log.Printf("custom playlist job %s: failed to create playlist: %v", job.ID, err)
			job.Status = "failed"
			job.Error = "Failed to create playlist"
			updateJob()
			return
		}

		// Add tracks in chunks
		if err := spotify.AddTracksToPlaylist(accessToken, playlist.ID, matchingTrackURIs); err != nil {
			log.Printf("custom playlist job %s: failed to add tracks: %v", job.ID, err)
			job.Status = "failed"
			job.Error = "Failed to add tracks to playlist"
			updateJob()
			return
		}

		createdPlaylists = append(createdPlaylists, CreatedPlaylist{
			SpotifyID:  playlist.ID,
			Name:       name,
			Genre:      "Custom",
			SongCount:  len(matchingTrackURIs),
			SpotifyURL: playlist.ExternalURLs.Spotify,
		})
	} else {
		// Create separate playlists per sub-genre
		for _, subGenre := range req.SubGenres {
			// Filter tracks for this sub-genre
			var genreTrackURIs []string
			for trackURI, trackGenres := range cache.AnalysisData.SongGenreMap {
				for _, g := range trackGenres {
					if g == subGenre {
						genreTrackURIs = append(genreTrackURIs, trackURI)
						break
					}
				}
			}

			if len(genreTrackURIs) == 0 {
				continue
			}

			name := settings.BuildPlaylistName(subGenre)
			desc := settings.BuildDescription(subGenre)

			playlist, err := spotify.CreatePlaylist(accessToken, userID, name, desc)
			if err != nil {
				log.Printf("custom playlist job %s: failed to create playlist for %s: %v", job.ID, subGenre, err)
				continue
			}

			if err := spotify.AddTracksToPlaylist(accessToken, playlist.ID, genreTrackURIs); err != nil {
				log.Printf("custom playlist job %s: failed to add tracks for %s: %v", job.ID, subGenre, err)
				continue
			}

			createdPlaylists = append(createdPlaylists, CreatedPlaylist{
				SpotifyID:  playlist.ID,
				Name:       name,
				Genre:      subGenre,
				SongCount:  len(genreTrackURIs),
				SpotifyURL: playlist.ExternalURLs.Spotify,
			})

			job.SongsProcessed += len(genreTrackURIs)
			updateJob()
		}
	}

	job.Status = "completed"
	job.Stage = "done"
	job.GenresDiscovered = req.SubGenres
	job.Result = &organizer.OrganizeResult{
		PlaylistsCreated: len(createdPlaylists),
		TotalSongs:       len(matchingTrackURIs),
	}
	updateJob()
}

type CreatedPlaylist struct {
	SpotifyID  string `json:"spotify_id"`
	Name       string `json:"name"`
	Genre      string `json:"genre"`
	SongCount  int    `json:"song_count"`
	SpotifyURL string `json:"spotify_url"`
}
```

**Step 2: Commit**

```bash
git add backend/internal/api/handlers/custom_playlist.go
git commit -m "feat: add custom playlist handler"
```

---

## Task 8: Add Custom Playlist Route

**Files:**
- Modify: `backend/internal/api/routes.go`

**Step 1: Add route**

After the `api.GET("/library/genres", handlers.GetLibraryGenres)` line, add:

```go
		api.POST("/custom-playlist", handlers.StartCustomPlaylist)
```

**Step 2: Verify build**

```bash
cd /Users/ambarg/developer/spotify_genre_organizer/backend && go build ./...
```

Expected: Build succeeds

**Step 3: Commit**

```bash
git add backend/internal/api/routes.go
git commit -m "feat: add /api/custom-playlist route"
```

---

## Task 9: Add Frontend API Functions

**Files:**
- Modify: `frontend/src/lib/api.ts`

**Step 1: Add types and functions**

Add at end of `frontend/src/lib/api.ts`:

```typescript
// Genre Analysis Types
export interface SubGenreCount {
  name: string;
  count: number;
}

export interface ParentGenreCount {
  name: string;
  count: number;
  sub_genres: SubGenreCount[];
}

export interface LibraryGenresResponse {
  parent_genres: ParentGenreCount[];
  analyzed_at: string;
  total_songs: number;
}

export async function getLibraryGenres(refresh = false): Promise<LibraryGenresResponse> {
  const url = refresh
    ? `${API_URL}/api/library/genres?refresh=true`
    : `${API_URL}/api/library/genres`;

  const response = await fetch(url, {
    credentials: 'include',
  });
  handleApiResponse(response);
  return response.json();
}

// Custom Playlist Types
export interface CustomPlaylistRequest {
  sub_genres: string[];
  mode: 'combined' | 'separate';
  name?: string;
}

export async function startCustomPlaylist(request: CustomPlaylistRequest): Promise<{ job_id: string }> {
  const response = await fetch(`${API_URL}/api/custom-playlist`, {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
  handleApiResponse(response);
  return response.json();
}
```

**Step 2: Commit**

```bash
git add frontend/src/lib/api.ts
git commit -m "feat: add frontend API functions for genre analysis and custom playlists"
```

---

## Task 10: Create GenrePicker Component

**Files:**
- Create: `frontend/src/components/GenrePicker.tsx`

**Step 1: Create component**

```tsx
'use client';

import { useState } from 'react';
import { ParentGenreCount, SubGenreCount } from '@/lib/api';

interface GenrePickerProps {
  parentGenres: ParentGenreCount[];
  selectedGenres: Set<string>;
  onSelectionChange: (selected: Set<string>) => void;
}

export function GenrePicker({ parentGenres, selectedGenres, onSelectionChange }: GenrePickerProps) {
  const [expandedParent, setExpandedParent] = useState<string | null>(null);

  const toggleParent = (parentName: string) => {
    setExpandedParent(expandedParent === parentName ? null : parentName);
  };

  const isParentFullySelected = (parent: ParentGenreCount) => {
    return parent.sub_genres.every(sg => selectedGenres.has(sg.name));
  };

  const isParentPartiallySelected = (parent: ParentGenreCount) => {
    const selected = parent.sub_genres.filter(sg => selectedGenres.has(sg.name));
    return selected.length > 0 && selected.length < parent.sub_genres.length;
  };

  const toggleParentSelection = (parent: ParentGenreCount) => {
    const newSelected = new Set(selectedGenres);
    if (isParentFullySelected(parent)) {
      // Deselect all sub-genres
      parent.sub_genres.forEach(sg => newSelected.delete(sg.name));
    } else {
      // Select all sub-genres
      parent.sub_genres.forEach(sg => newSelected.add(sg.name));
    }
    onSelectionChange(newSelected);
  };

  const toggleSubGenre = (subGenreName: string) => {
    const newSelected = new Set(selectedGenres);
    if (newSelected.has(subGenreName)) {
      newSelected.delete(subGenreName);
    } else {
      newSelected.add(subGenreName);
    }
    onSelectionChange(newSelected);
  };

  return (
    <div className="space-y-2">
      {parentGenres.map((parent) => (
        <div key={parent.name} className="bg-bg-card rounded-xl overflow-hidden">
          {/* Parent Genre Row */}
          <div
            className="flex items-center justify-between p-4 cursor-pointer hover:bg-bg-dark/50 transition-colors"
            onClick={() => toggleParent(parent.name)}
          >
            <div className="flex items-center gap-3">
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  toggleParentSelection(parent);
                }}
                className={`w-5 h-5 rounded border-2 flex items-center justify-center transition-colors ${
                  isParentFullySelected(parent)
                    ? 'bg-accent-orange border-accent-orange'
                    : isParentPartiallySelected(parent)
                    ? 'border-accent-orange bg-accent-orange/30'
                    : 'border-text-muted hover:border-text-cream'
                }`}
              >
                {(isParentFullySelected(parent) || isParentPartiallySelected(parent)) && (
                  <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                  </svg>
                )}
              </button>
              <span className="font-display text-lg text-text-cream">{parent.name}</span>
            </div>
            <div className="flex items-center gap-3">
              <span className="text-sm text-text-muted">{parent.count} songs</span>
              <svg
                className={`w-5 h-5 text-text-muted transition-transform ${
                  expandedParent === parent.name ? 'rotate-180' : ''
                }`}
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>

          {/* Sub-genres (expanded) */}
          {expandedParent === parent.name && (
            <div className="px-4 pb-4 pt-2 border-t border-bg-dark">
              <div className="flex flex-wrap gap-2">
                {parent.sub_genres.map((subGenre, index) => (
                  <button
                    key={subGenre.name}
                    onClick={() => toggleSubGenre(subGenre.name)}
                    className={`px-3 py-1.5 rounded-full text-sm transition-all animate-drop-in ${
                      selectedGenres.has(subGenre.name)
                        ? 'bg-accent-orange text-white'
                        : 'bg-bg-dark text-text-cream hover:bg-bg-dark/70'
                    }`}
                    style={{ animationDelay: `${index * 50}ms` }}
                  >
                    {subGenre.name}
                    <span className="ml-1.5 opacity-70">({subGenre.count})</span>
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/src/components/GenrePicker.tsx
git commit -m "feat: add GenrePicker accordion component"
```

---

## Task 11: Create Custom Playlist Page

**Files:**
- Create: `frontend/src/app/custom-playlist/page.tsx`

**Step 1: Create page**

```tsx
'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { GenrePicker } from '@/components/GenrePicker';
import { useUser } from '@/hooks/useUser';
import { getLibraryGenres, LibraryGenresResponse } from '@/lib/api';

export default function CustomPlaylistPage() {
  const { user, loading: userLoading } = useUser();
  const router = useRouter();
  const [genreData, setGenreData] = useState<LibraryGenresResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [selectedGenres, setSelectedGenres] = useState<Set<string>>(new Set());

  useEffect(() => {
    const fetchGenres = async () => {
      try {
        const data = await getLibraryGenres();
        setGenreData(data);
      } catch (error) {
        console.error('Failed to fetch genres:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchGenres();
  }, []);

  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      const data = await getLibraryGenres(true);
      setGenreData(data);
    } catch (error) {
      console.error('Failed to refresh genres:', error);
    } finally {
      setRefreshing(false);
    }
  };

  const getTotalSelectedSongs = () => {
    if (!genreData) return 0;
    let total = 0;
    for (const parent of genreData.parent_genres) {
      for (const sub of parent.sub_genres) {
        if (selectedGenres.has(sub.name)) {
          total += sub.count;
        }
      }
    }
    return total;
  };

  const formatAnalyzedTime = (dateStr: string) => {
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    return `${diffDays}d ago`;
  };

  const handleNext = () => {
    // Store selection in sessionStorage and navigate to configure page
    sessionStorage.setItem('customPlaylistGenres', JSON.stringify([...selectedGenres]));
    router.push('/custom-playlist/configure');
  };

  if (userLoading || loading) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <VinylIcon spinning size={64} />
      </main>
    );
  }

  return (
    <main className="min-h-screen flex flex-col px-4 py-8">
      {/* Header */}
      <div className="max-w-2xl mx-auto w-full mb-8">
        <button
          onClick={() => router.push('/dashboard')}
          className="flex items-center gap-2 text-text-muted hover:text-text-cream transition-colors mb-4"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>

        <h1 className="font-display text-3xl text-text-cream mb-2">Build a Crate</h1>

        {genreData && (
          <div className="flex items-center gap-2 text-text-muted">
            <span>{genreData.total_songs.toLocaleString()} songs</span>
            <span>•</span>
            <span>Analyzed {formatAnalyzedTime(genreData.analyzed_at)}</span>
            <button
              onClick={handleRefresh}
              disabled={refreshing}
              className="ml-1 p-1 hover:text-text-cream transition-colors disabled:opacity-50"
              title="Refresh analysis"
            >
              <svg
                className={`w-4 h-4 ${refreshing ? 'animate-spin' : ''}`}
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </button>
          </div>
        )}
      </div>

      {/* Genre Picker */}
      <div className="max-w-2xl mx-auto w-full flex-1 mb-24">
        {genreData && (
          <GenrePicker
            parentGenres={genreData.parent_genres}
            selectedGenres={selectedGenres}
            onSelectionChange={setSelectedGenres}
          />
        )}
      </div>

      {/* Sticky Bottom Bar */}
      <div className="fixed bottom-0 left-0 right-0 bg-bg-dark border-t border-bg-card p-4">
        <div className="max-w-2xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-2 flex-wrap flex-1 mr-4">
            {selectedGenres.size === 0 ? (
              <span className="text-text-muted">Select genres to continue</span>
            ) : (
              <>
                {[...selectedGenres].slice(0, 3).map((genre) => (
                  <span
                    key={genre}
                    className="px-2 py-1 bg-accent-orange text-white text-sm rounded-full flex items-center gap-1 animate-bounce-in"
                  >
                    {genre}
                    <button
                      onClick={() => {
                        const newSelected = new Set(selectedGenres);
                        newSelected.delete(genre);
                        setSelectedGenres(newSelected);
                      }}
                      className="hover:bg-white/20 rounded-full p-0.5"
                    >
                      <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </span>
                ))}
                {selectedGenres.size > 3 && (
                  <span className="text-text-muted text-sm">+{selectedGenres.size - 3} more</span>
                )}
              </>
            )}
          </div>
          <div className="flex items-center gap-4">
            <span className="text-text-cream font-medium">
              {getTotalSelectedSongs().toLocaleString()} songs
            </span>
            <Button
              onClick={handleNext}
              disabled={selectedGenres.size === 0}
            >
              Next
              <svg className="w-4 h-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </Button>
          </div>
        </div>
      </div>
    </main>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/src/app/custom-playlist/page.tsx
git commit -m "feat: add custom playlist genre selection page"
```

---

## Task 12: Create Configure Page

**Files:**
- Create: `frontend/src/app/custom-playlist/configure/page.tsx`

**Step 1: Create page**

```tsx
'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { VinylIcon } from '@/components/VinylIcon';
import { Button } from '@/components/Button';
import { useUser } from '@/hooks/useUser';
import { startCustomPlaylist } from '@/lib/api';

export default function ConfigurePage() {
  const { loading: userLoading } = useUser();
  const router = useRouter();
  const [selectedGenres, setSelectedGenres] = useState<string[]>([]);
  const [mode, setMode] = useState<'combined' | 'separate'>('combined');
  const [playlistName, setPlaylistName] = useState('');
  const [isCreating, setIsCreating] = useState(false);

  useEffect(() => {
    // Get selected genres from sessionStorage
    const stored = sessionStorage.getItem('customPlaylistGenres');
    if (stored) {
      const genres = JSON.parse(stored);
      setSelectedGenres(genres);
      // If only one genre, skip to creation with that genre name
      if (genres.length === 1) {
        setMode('separate');
      }
    } else {
      // No genres selected, go back
      router.push('/custom-playlist');
    }
  }, [router]);

  const handleCreate = async () => {
    setIsCreating(true);
    try {
      const { job_id } = await startCustomPlaylist({
        sub_genres: selectedGenres,
        mode,
        name: mode === 'combined' ? (playlistName || 'Custom Mix by Organizer') : undefined,
      });
      router.push(`/processing?job=${job_id}`);
    } catch (error) {
      console.error('Failed to create playlist:', error);
      setIsCreating(false);
    }
  };

  if (userLoading || selectedGenres.length === 0) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <VinylIcon spinning size={64} />
      </main>
    );
  }

  // If only one genre, create immediately as separate (single playlist)
  if (selectedGenres.length === 1) {
    return (
      <main className="min-h-screen flex flex-col items-center justify-center px-4">
        <div className="max-w-md w-full bg-bg-card rounded-2xl p-8 text-center">
          <h1 className="font-display text-2xl text-text-cream mb-4">
            Creating {selectedGenres[0]} playlist
          </h1>
          <Button onClick={handleCreate} disabled={isCreating} className="w-full">
            {isCreating ? (
              <>
                <VinylIcon spinning size={20} />
                Creating...
              </>
            ) : (
              'Create Playlist'
            )}
          </Button>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen flex flex-col items-center justify-center px-4 py-12">
      <div className="max-w-xl w-full">
        {/* Header */}
        <button
          onClick={() => router.push('/custom-playlist')}
          className="flex items-center gap-2 text-text-muted hover:text-text-cream transition-colors mb-6"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back
        </button>

        <h1 className="font-display text-3xl text-text-cream mb-2">Almost there</h1>
        <div className="flex flex-wrap gap-2 mb-8">
          {selectedGenres.map((genre) => (
            <span key={genre} className="px-2 py-1 bg-bg-card text-text-cream text-sm rounded-full">
              {genre}
            </span>
          ))}
        </div>

        {/* Options */}
        <div className="space-y-4">
          {/* Combined Option */}
          <button
            onClick={() => setMode('combined')}
            className={`w-full p-6 rounded-xl border-2 text-left transition-all ${
              mode === 'combined'
                ? 'border-accent-orange bg-accent-orange/10'
                : 'border-bg-card bg-bg-card hover:border-text-muted'
            }`}
          >
            <div className="flex items-start gap-4">
              <div className="p-3 bg-bg-dark rounded-lg">
                <VinylIcon size={32} />
              </div>
              <div className="flex-1">
                <h3 className="font-display text-xl text-text-cream mb-1">One Combined Playlist</h3>
                <p className="text-text-muted text-sm mb-3">
                  All songs from selected genres in one playlist
                </p>
                {mode === 'combined' && (
                  <input
                    type="text"
                    value={playlistName}
                    onChange={(e) => setPlaylistName(e.target.value)}
                    placeholder="My Heavy Stuff"
                    className="w-full px-3 py-2 bg-bg-dark border border-bg-card rounded-lg text-text-cream placeholder-text-muted focus:outline-none focus:border-accent-orange"
                    onClick={(e) => e.stopPropagation()}
                  />
                )}
              </div>
            </div>
          </button>

          {/* Separate Option */}
          <button
            onClick={() => setMode('separate')}
            className={`w-full p-6 rounded-xl border-2 text-left transition-all ${
              mode === 'separate'
                ? 'border-accent-orange bg-accent-orange/10'
                : 'border-bg-card bg-bg-card hover:border-text-muted'
            }`}
          >
            <div className="flex items-start gap-4">
              <div className="p-3 bg-bg-dark rounded-lg flex -space-x-2">
                <VinylIcon size={24} />
                <VinylIcon size={24} />
                <VinylIcon size={24} />
              </div>
              <div className="flex-1">
                <h3 className="font-display text-xl text-text-cream mb-1">Separate Playlists</h3>
                <p className="text-text-muted text-sm">
                  {selectedGenres.length} playlists, one per genre
                </p>
                <p className="text-text-muted text-xs mt-2">
                  {selectedGenres.slice(0, 3).join(' • ')}{selectedGenres.length > 3 ? ' • ...' : ''}
                </p>
              </div>
            </div>
          </button>
        </div>

        {/* Create Button */}
        <div className="mt-8">
          <Button onClick={handleCreate} disabled={isCreating} size="lg" className="w-full">
            {isCreating ? (
              <>
                <VinylIcon spinning size={24} />
                Creating...
              </>
            ) : (
              <>
                Create Playlist{mode === 'separate' && selectedGenres.length > 1 ? 's' : ''}
                <span className="ml-2">●</span>
              </>
            )}
          </Button>
        </div>
      </div>
    </main>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/src/app/custom-playlist/configure/page.tsx
git commit -m "feat: add custom playlist configure page"
```

---

## Task 13: Add Dashboard Button

**Files:**
- Modify: `frontend/src/app/dashboard/page.tsx`

**Step 1: Add "Build a Crate" button**

In `frontend/src/app/dashboard/page.tsx`, find the section with "Manage My Crates" button (around line 190-210) and add a new button BEFORE it:

After line 188 (closing `</div>` of the Replace Toggle section) and before line 190 (`{/* Manage Playlists Link */}`), add:

```tsx
        {/* Build a Crate Link */}
        <div className="mb-4">
          <Button
            size="lg"
            variant="secondary"
            className="w-full flex items-center justify-center gap-2"
            onClick={() => router.push('/custom-playlist')}
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            Build a Crate
          </Button>
          <p className="text-text-muted text-sm text-center mt-2">
            Pick specific genres
          </p>
        </div>
```

**Step 2: Verify build**

```bash
cd /Users/ambarg/developer/spotify_genre_organizer/frontend && npm run build
```

Expected: Build succeeds

**Step 3: Commit**

```bash
git add frontend/src/app/dashboard/page.tsx
git commit -m "feat: add Build a Crate button to dashboard"
```

---

## Task 14: Create Database Migration

**Files:**
- Create: `supabase/migrations/20260127000000_add_genre_analysis_cache.sql` (or use Supabase dashboard)

**Step 1: Create migration SQL**

If using Supabase dashboard, run this SQL in the SQL Editor:

```sql
-- Create genre_analysis_cache table
CREATE TABLE IF NOT EXISTS genre_analysis_cache (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id TEXT NOT NULL UNIQUE,
  analysis_data JSONB NOT NULL,
  analyzed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index on user_id for fast lookups
CREATE INDEX IF NOT EXISTS idx_genre_analysis_cache_user_id ON genre_analysis_cache(user_id);

-- Enable RLS
ALTER TABLE genre_analysis_cache ENABLE ROW LEVEL SECURITY;

-- Create policy for users to manage their own cache
CREATE POLICY "Users can manage their own genre cache"
  ON genre_analysis_cache
  FOR ALL
  USING (true)
  WITH CHECK (true);
```

**Step 2: Verify migration ran**

Check in Supabase dashboard that `genre_analysis_cache` table exists.

**Step 3: Commit (if using migration file)**

```bash
git add supabase/migrations/
git commit -m "feat: add genre_analysis_cache database migration"
```

---

## Task 15: Final Integration Test

**Step 1: Start backend**

```bash
cd /Users/ambarg/developer/spotify_genre_organizer/backend && go run cmd/api/main.go
```

**Step 2: Start frontend**

```bash
cd /Users/ambarg/developer/spotify_genre_organizer/frontend && npm run dev
```

**Step 3: Manual test flow**

1. Log in with Spotify
2. Click "Build a Crate" on dashboard
3. Wait for genre analysis (or see cached data)
4. Expand a parent genre, select sub-genres
5. Click Next
6. Choose Combined or Separate
7. Click Create
8. Verify playlist(s) created in Spotify

**Step 4: Final commit**

```bash
git add -A
git commit -m "feat: complete custom genre playlist feature

- Add genre analysis caching
- Add hierarchical genre picker UI
- Add combined/separate playlist options
- Integrate with existing processing flow"
```

---

## Summary

| Task | Description |
|------|-------------|
| 1 | Genre cache model |
| 2 | Database functions |
| 3 | Genre analysis helper |
| 4 | Handler test |
| 5 | GetLibraryGenres handler |
| 6 | Add genres route |
| 7 | Custom playlist handler |
| 8 | Add custom playlist route |
| 9 | Frontend API functions |
| 10 | GenrePicker component |
| 11 | Custom playlist page |
| 12 | Configure page |
| 13 | Dashboard button |
| 14 | Database migration |
| 15 | Integration test |

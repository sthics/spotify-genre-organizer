# Auto-Sync Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add smart sync detection that shows users when new songs are waiting to be organized, with one-click sync for all playlists.

**Architecture:** Backend gets new endpoint to calculate sync status by comparing liked songs' `added_at` timestamps against playlist `last_synced_at`. Frontend shows badge on dashboard and Sync All button on playlists page. Toasts provide feedback.

**Tech Stack:** Go/Gin backend, Next.js/React frontend, Supabase PostgreSQL

---

## Task 1: Add `added_at` Field to Song Struct

**Files:**
- Modify: `backend/internal/spotify/library.go:11-16`

**Step 1: Update Song struct and parsing**

Add `AddedAt` field to track when song was liked:

```go
type Song struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Artists []Artist  `json:"artists"`
	Genres  []string  `json:"genres"`
	AddedAt time.Time `json:"added_at"`
}
```

**Step 2: Update likedSongsResponse struct**

```go
type likedSongsResponse struct {
	Items []struct {
		AddedAt string `json:"added_at"`
		Track   struct {
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
```

**Step 3: Update ParseLikedSongsResponse to parse AddedAt**

```go
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

		addedAt, _ := time.Parse(time.RFC3339, item.AddedAt)

		songs[i] = Song{
			ID:      item.Track.ID,
			Name:    item.Track.Name,
			Artists: artists,
			AddedAt: addedAt,
		}
	}

	next := ""
	if resp.Next != nil {
		next = *resp.Next
	}

	return songs, resp.Total, next, nil
}
```

**Step 4: Commit**

```bash
git add backend/internal/spotify/library.go
git commit -m "feat(backend): add added_at field to Song struct"
```

---

## Task 2: Add Database Helper for Oldest Sync Timestamp

**Files:**
- Modify: `backend/internal/database/settings.go`

**Step 1: Add GetOldestSyncTimestamp function**

```go
// GetOldestSyncTimestamp returns the oldest last_synced_at from user's playlist overrides
func GetOldestSyncTimestamp(userID string) (*time.Time, error) {
	res, _, err := Client.From("playlist_overrides").
		Select("last_synced_at", "", false).
		Eq("user_id", userID).
		Not("last_synced_at", "is", "null").
		Order("last_synced_at", &postgrest.OrderOpts{Ascending: true}).
		Limit(1, "").
		Execute()

	if err != nil {
		return nil, err
	}

	var results []struct {
		LastSyncedAt *time.Time `json:"last_synced_at"`
	}
	if err := json.Unmarshal(res, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 || results[0].LastSyncedAt == nil {
		return nil, nil
	}

	return results[0].LastSyncedAt, nil
}
```

**Step 2: Add import for postgrest at top of file**

```go
import (
	"encoding/json"
	"time"

	"github.com/spotify-genre-organizer/backend/internal/models"
	"github.com/supabase-community/postgrest-go"
)
```

**Step 3: Commit**

```bash
git add backend/internal/database/settings.go
git commit -m "feat(backend): add GetOldestSyncTimestamp helper"
```

---

## Task 3: Create Sync Status Handler

**Files:**
- Create: `backend/internal/api/handlers/sync.go`

**Step 1: Create the sync status handler file**

```go
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spotify-genre-organizer/backend/internal/database"
	"github.com/spotify-genre-organizer/backend/internal/genres"
	"github.com/spotify-genre-organizer/backend/internal/spotify"
)

type PlaylistSyncStatus struct {
	SpotifyID string `json:"spotify_id"`
	Genre     string `json:"genre"`
	NewCount  int    `json:"new_count"`
}

type SyncStatusResponse struct {
	NewSongsCount int                  `json:"new_songs_count"`
	OldestSyncAt  *time.Time           `json:"oldest_sync_at"`
	Playlists     []PlaylistSyncStatus `json:"playlists"`
}

func GetSyncStatus(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userID, _ := c.Cookie("user_id")

	// Get oldest sync timestamp
	oldestSync, err := database.GetOldestSyncTimestamp(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sync status"})
		return
	}

	// If no playlists synced yet, return empty
	if oldestSync == nil {
		c.JSON(http.StatusOK, SyncStatusResponse{
			NewSongsCount: 0,
			OldestSyncAt:  nil,
			Playlists:     []PlaylistSyncStatus{},
		})
		return
	}

	// Fetch all liked songs
	songs, err := spotify.FetchAllLikedSongs(accessToken, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch songs"})
		return
	}

	// Filter to songs added after oldest sync
	var newSongs []spotify.Song
	for _, song := range songs {
		if song.AddedAt.After(*oldestSync) {
			newSongs = append(newSongs, song)
		}
	}

	if len(newSongs) == 0 {
		c.JSON(http.StatusOK, SyncStatusResponse{
			NewSongsCount: 0,
			OldestSyncAt:  oldestSync,
			Playlists:     []PlaylistSyncStatus{},
		})
		return
	}

	// Enrich new songs with genres
	artistGenres, err := spotify.FetchAllArtistGenres(accessToken, newSongs, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch genres"})
		return
	}
	spotify.EnrichSongsWithGenres(newSongs, artistGenres)

	// Get user's playlist overrides to know which playlists exist
	overrides, err := database.GetPlaylistOverrides(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get playlists"})
		return
	}

	// Count new songs per genre/playlist
	genreCounts := make(map[string]int)
	for _, song := range newSongs {
		genre := genres.ScoreGenres(song.Genres)
		genreCounts[genre]++
	}

	// Build playlist status list
	var playlistStatuses []PlaylistSyncStatus
	for playlistID, override := range overrides {
		if override.Genre != "" {
			count := genreCounts[override.Genre]
			if count > 0 {
				playlistStatuses = append(playlistStatuses, PlaylistSyncStatus{
					SpotifyID: playlistID,
					Genre:     override.Genre,
					NewCount:  count,
				})
			}
		}
	}

	c.JSON(http.StatusOK, SyncStatusResponse{
		NewSongsCount: len(newSongs),
		OldestSyncAt:  oldestSync,
		Playlists:     playlistStatuses,
	})
}
```

**Step 2: Commit**

```bash
git add backend/internal/api/handlers/sync.go
git commit -m "feat(backend): add sync status handler"
```

---

## Task 4: Create Sync All Handler

**Files:**
- Modify: `backend/internal/api/handlers/sync.go`

**Step 1: Add SyncAllPlaylists handler to sync.go**

```go
type SyncAllResponse struct {
	PlaylistsUpdated int `json:"playlists_updated"`
	SongsAdded       int `json:"songs_added"`
}

func SyncAllPlaylists(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	userID, _ := c.Cookie("user_id")

	// Fetch all liked songs
	songs, err := spotify.FetchAllLikedSongs(accessToken, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch songs"})
		return
	}

	// Enrich with genres
	artistGenres, err := spotify.FetchAllArtistGenres(accessToken, songs, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch genres"})
		return
	}
	spotify.EnrichSongsWithGenres(songs, artistGenres)

	// Group songs by genre
	songsByGenre := make(map[string][]spotify.Song)
	for _, song := range songs {
		genre := genres.ScoreGenres(song.Genres)
		songsByGenre[genre] = append(songsByGenre[genre], song)
	}

	// Get user's playlist overrides
	overrides, err := database.GetPlaylistOverrides(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get playlists"})
		return
	}

	playlistsUpdated := 0
	totalSongsAdded := 0
	now := time.Now()

	for playlistID, override := range overrides {
		if override.Genre == "" {
			continue
		}

		genreSongs := songsByGenre[override.Genre]
		if len(genreSongs) == 0 {
			continue
		}

		// Clear and repopulate playlist
		if err := spotify.ClearPlaylist(accessToken, playlistID); err != nil {
			continue
		}

		trackIDs := make([]string, len(genreSongs))
		for i, s := range genreSongs {
			trackIDs[i] = s.ID
		}

		if err := spotify.AddTracksToPlaylist(accessToken, playlistID, trackIDs); err != nil {
			continue
		}

		// Update last_synced_at
		override.LastSyncedAt = &now
		database.SavePlaylistOverride(override)

		playlistsUpdated++
		totalSongsAdded += len(genreSongs)
	}

	c.JSON(http.StatusOK, SyncAllResponse{
		PlaylistsUpdated: playlistsUpdated,
		SongsAdded:       totalSongsAdded,
	})
}
```

**Step 2: Commit**

```bash
git add backend/internal/api/handlers/sync.go
git commit -m "feat(backend): add sync all playlists handler"
```

---

## Task 5: Update RefreshPlaylist to Set last_synced_at

**Files:**
- Modify: `backend/internal/api/handlers/playlists.go:259-357`

**Step 1: Add time import if not present**

At top of file, ensure `time` is imported.

**Step 2: Update RefreshPlaylist to set last_synced_at**

At the end of `RefreshPlaylist`, before the success response, add:

```go
	// Update last_synced_at
	now := time.Now()
	override.LastSyncedAt = &now
	savePlaylistOverride(override)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"song_count": len(genreSongs),
	})
```

**Step 3: Commit**

```bash
git add backend/internal/api/handlers/playlists.go
git commit -m "feat(backend): update last_synced_at on playlist refresh"
```

---

## Task 6: Register New Routes

**Files:**
- Modify: `backend/internal/api/routes.go:121-129`

**Step 1: Add sync routes**

After the existing playlist routes, add:

```go
		api.GET("/library/sync-status", handlers.GetSyncStatus)
		api.POST("/playlists/sync-all", handlers.SyncAllPlaylists)
```

**Step 2: Commit**

```bash
git add backend/internal/api/routes.go
git commit -m "feat(backend): add sync status and sync all routes"
```

---

## Task 7: Create Toast Component

**Files:**
- Create: `frontend/src/components/Toast.tsx`

**Step 1: Create the Toast component**

```tsx
'use client';

import { useEffect, useState } from 'react';

export interface ToastData {
  id: string;
  message: string;
  type: 'success' | 'error';
  action?: {
    label: string;
    onClick: () => void;
  };
}

interface ToastProps {
  toast: ToastData;
  onDismiss: (id: string) => void;
}

export function Toast({ toast, onDismiss }: ToastProps) {
  const [isExiting, setIsExiting] = useState(false);

  useEffect(() => {
    if (toast.type === 'success') {
      const timer = setTimeout(() => {
        setIsExiting(true);
        setTimeout(() => onDismiss(toast.id), 300);
      }, 4000);
      return () => clearTimeout(timer);
    }
  }, [toast.id, toast.type, onDismiss]);

  const handleDismiss = () => {
    setIsExiting(true);
    setTimeout(() => onDismiss(toast.id), 300);
  };

  const borderColor = toast.type === 'success' ? 'border-l-success-green' : 'border-l-red-500';

  return (
    <div
      className={`
        bg-bg-card rounded-lg shadow-xl border-l-4 ${borderColor}
        p-4 min-w-[280px] max-w-[380px]
        transform transition-all duration-300
        ${isExiting ? 'opacity-0 translate-x-4' : 'opacity-100 translate-x-0'}
      `}
    >
      <div className="flex items-start justify-between gap-3">
        <p className="text-text-cream text-sm">{toast.message}</p>
        <button
          onClick={handleDismiss}
          className="text-text-muted hover:text-text-cream transition-colors"
        >
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      {toast.action && (
        <button
          onClick={toast.action.onClick}
          className="mt-2 text-accent-orange hover:text-accent-orange-hover text-sm font-medium transition-colors"
        >
          {toast.action.label}
        </button>
      )}
    </div>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/src/components/Toast.tsx
git commit -m "feat(frontend): add Toast component"
```

---

## Task 8: Create Toast Context Provider

**Files:**
- Create: `frontend/src/contexts/ToastContext.tsx`

**Step 1: Create the Toast context and provider**

```tsx
'use client';

import { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { Toast, ToastData } from '@/components/Toast';

interface ToastContextType {
  showToast: (message: string, type: 'success' | 'error', action?: { label: string; onClick: () => void }) => void;
}

const ToastContext = createContext<ToastContextType | null>(null);

export function useToast() {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within ToastProvider');
  }
  return context;
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastData[]>([]);

  const showToast = useCallback((
    message: string,
    type: 'success' | 'error',
    action?: { label: string; onClick: () => void }
  ) => {
    const id = crypto.randomUUID();
    setToasts(prev => [...prev, { id, message, type, action }]);
  }, []);

  const dismissToast = useCallback((id: string) => {
    setToasts(prev => prev.filter(t => t.id !== id));
  }, []);

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
        {toasts.map(toast => (
          <Toast key={toast.id} toast={toast} onDismiss={dismissToast} />
        ))}
      </div>
    </ToastContext.Provider>
  );
}
```

**Step 2: Commit**

```bash
git add frontend/src/contexts/ToastContext.tsx
git commit -m "feat(frontend): add ToastContext provider"
```

---

## Task 9: Add Toast Provider to Layout

**Files:**
- Modify: `frontend/src/app/layout.tsx`

**Step 1: Import and wrap with ToastProvider**

Add import at top:

```tsx
import { ToastProvider } from '@/contexts/ToastContext';
```

Wrap the body content with ToastProvider:

```tsx
<body>
  <ToastProvider>
    <div className="grain-overlay" />
    {children}
  </ToastProvider>
</body>
```

**Step 2: Commit**

```bash
git add frontend/src/app/layout.tsx
git commit -m "feat(frontend): add ToastProvider to app layout"
```

---

## Task 10: Add Sync API Functions

**Files:**
- Modify: `frontend/src/lib/api.ts`

**Step 1: Add sync status types and function**

```typescript
export interface PlaylistSyncStatus {
  spotify_id: string;
  genre: string;
  new_count: number;
}

export interface SyncStatus {
  new_songs_count: number;
  oldest_sync_at: string | null;
  playlists: PlaylistSyncStatus[];
}

export async function getSyncStatus(): Promise<SyncStatus> {
  const response = await fetch(`${API_URL}/api/library/sync-status`, {
    credentials: 'include',
  });
  if (!response.ok) throw new Error('Failed to get sync status');
  return response.json();
}

export async function syncAllPlaylists(): Promise<{ playlists_updated: number; songs_added: number }> {
  const response = await fetch(`${API_URL}/api/playlists/sync-all`, {
    method: 'POST',
    credentials: 'include',
  });
  if (!response.ok) throw new Error('Failed to sync playlists');
  return response.json();
}
```

**Step 2: Commit**

```bash
git add frontend/src/lib/api.ts
git commit -m "feat(frontend): add sync status and sync all API functions"
```

---

## Task 11: Update Dashboard with Sync Badge

**Files:**
- Modify: `frontend/src/app/dashboard/page.tsx`

**Step 1: Import getSyncStatus and add state**

Add to imports:

```typescript
import { startOrganize, logout, getLibraryCount, getSyncStatus } from '@/lib/api';
```

Add state:

```typescript
const [newSongsCount, setNewSongsCount] = useState(0);
```

**Step 2: Fetch sync status in useEffect**

Add after the library count fetch:

```typescript
// Fetch sync status
const fetchSyncStatus = async () => {
  try {
    const status = await getSyncStatus();
    setNewSongsCount(status.new_songs_count);
  } catch (error) {
    // Silently fail - no playlists synced yet is fine
  }
};

fetchSyncStatus();
```

**Step 3: Update Manage My Crates button**

Replace the existing button with:

```tsx
{/* Manage Playlists Link */}
<div className="mb-4">
  <Button
    size="lg"
    variant="secondary"
    className="w-full flex items-center justify-center gap-2 relative"
    onClick={() => router.push('/playlists')}
  >
    <VinylIcon size={20} />
    Manage My Crates
    {newSongsCount > 0 && (
      <span className="bg-accent-orange text-white text-xs font-bold px-2 py-0.5 rounded-full ml-2 animate-fade-in">
        {newSongsCount}
      </span>
    )}
  </Button>
  {newSongsCount > 0 && (
    <p className="text-text-muted text-sm text-center mt-2 animate-fade-in">
      {newSongsCount} song{newSongsCount !== 1 ? 's' : ''} waiting to be organized
    </p>
  )}
</div>
```

**Step 4: Commit**

```bash
git add frontend/src/app/dashboard/page.tsx
git commit -m "feat(frontend): add sync badge to dashboard"
```

---

## Task 12: Update Playlists Page with Sync All

**Files:**
- Modify: `frontend/src/app/playlists/page.tsx`

**Step 1: Add imports**

```typescript
import { useToast } from '@/contexts/ToastContext';
import {
    getPlaylists,
    refreshPlaylist,
    updatePlaylist,
    deletePlaylist,
    getSyncStatus,
    syncAllPlaylists,
    ManagedPlaylist,
    SyncStatus,
} from '@/lib/api';
```

**Step 2: Add state for sync**

```typescript
const { showToast } = useToast();
const [syncStatus, setSyncStatus] = useState<SyncStatus | null>(null);
const [isSyncingAll, setIsSyncingAll] = useState(false);
```

**Step 3: Add sync status fetch**

Add after loadPlaylists:

```typescript
const loadSyncStatus = async () => {
    try {
        const status = await getSyncStatus();
        setSyncStatus(status);
    } catch (err) {
        // No synced playlists yet
    }
};

useEffect(() => {
    loadPlaylists();
    loadSyncStatus();
}, []);
```

**Step 4: Add handleSyncAll function**

```typescript
const handleSyncAll = async () => {
    if (!syncStatus || syncStatus.new_songs_count === 0) return;

    setIsSyncingAll(true);

    // Show starting toast with random message
    const messages = [
        'Sorting your new vinyl...',
        `Filing ${syncStatus.new_songs_count} tracks into your crates...`,
        'Spinning up your new additions...',
        'Dropping new records into your crates...',
        `${syncStatus.new_songs_count} new tracks heading to their crates...`,
    ];
    const startMessage = messages[Math.floor(Math.random() * messages.length)];
    showToast(startMessage, 'success');

    try {
        const result = await syncAllPlaylists();
        showToast(
            `${result.playlists_updated} crates updated • ${result.songs_added} songs added`,
            'success'
        );
        setSyncStatus(null);
        loadPlaylists();
    } catch (err) {
        showToast('Sync failed — Retry', 'error', {
            label: 'Retry',
            onClick: handleSyncAll,
        });
    } finally {
        setIsSyncingAll(false);
    }
};
```

**Step 5: Update header to show Sync All button**

Replace the header Button:

```tsx
{syncStatus && syncStatus.new_songs_count > 0 ? (
    <Button onClick={handleSyncAll} disabled={isSyncingAll}>
        {isSyncingAll ? (
            <>
                <VinylIcon spinning size={20} />
                Syncing...
            </>
        ) : (
            `Sync All (${syncStatus.new_songs_count} new)`
        )}
    </Button>
) : (
    <Button onClick={() => router.push('/dashboard')}>
        + New Organize
    </Button>
)}
```

**Step 6: Add per-playlist new count indicator**

In the playlist card, after the genre badge, add:

```tsx
{syncStatus?.playlists.find(p => p.spotify_id === playlist.spotify_id)?.new_count && (
    <span
        className="text-sm font-medium ml-auto"
        style={{ color }}
    >
        +{syncStatus.playlists.find(p => p.spotify_id === playlist.spotify_id)?.new_count} new
    </span>
)}
```

**Step 7: Update refresh handler to show toast**

```typescript
const handleRefresh = async (id: string) => {
    const playlist = playlists.find(p => p.spotify_id === id);
    setRefreshingId(id);
    try {
        const result = await refreshPlaylist(id);
        setPlaylists((prev) =>
            prev.map((p) =>
                p.spotify_id === id ? { ...p, song_count: result.song_count } : p
            )
        );
        showToast(`Added ${result.song_count} songs to ${playlist?.genre || 'playlist'}`, 'success');
        loadSyncStatus(); // Refresh sync status
    } catch (err) {
        showToast(`Couldn't sync ${playlist?.genre || 'playlist'} — Retry`, 'error', {
            label: 'Retry',
            onClick: () => handleRefresh(id),
        });
    }
    setRefreshingId(null);
};
```

**Step 8: Commit**

```bash
git add frontend/src/app/playlists/page.tsx
git commit -m "feat(frontend): add Sync All and per-playlist sync indicators"
```

---

## Task 13: Update Organize Handler to Set last_synced_at

**Files:**
- Modify: `backend/internal/api/handlers/organize.go`

**Step 1: Find where playlists are created in processOrganizeJob**

After each playlist is created successfully, save an override with last_synced_at:

```go
now := time.Now()
override := &models.PlaylistOverride{
    UserID:            userID,
    PlaylistSpotifyID: playlistID,
    Genre:             genre,
    LastSyncedAt:      &now,
}
database.SavePlaylistOverride(override)
```

**Step 2: Add necessary imports**

Ensure `time` and the database/models packages are imported.

**Step 3: Commit**

```bash
git add backend/internal/api/handlers/organize.go
git commit -m "feat(backend): set last_synced_at when creating playlists"
```

---

## Task 14: Manual Testing Checklist

**Test 1: New user flow**
- [ ] Fresh user sees no badge on dashboard
- [ ] After first organize, badge appears if user adds new songs

**Test 2: Sync status detection**
- [ ] Like a new song in Spotify
- [ ] Refresh dashboard - badge should appear
- [ ] Badge shows correct count

**Test 3: Sync All**
- [ ] Click Sync All button
- [ ] Toast shows progress message
- [ ] On complete, summary toast appears
- [ ] Badge clears

**Test 4: Individual sync**
- [ ] Expand a playlist
- [ ] Click Refresh
- [ ] Toast shows "Added X songs to Genre"

**Test 5: Error handling**
- [ ] Disconnect network during sync
- [ ] Error toast with Retry appears
- [ ] Retry works

---

## Summary

This implementation adds:

1. **Backend:**
   - `GET /api/library/sync-status` - Returns new song count and per-playlist breakdown
   - `POST /api/playlists/sync-all` - Syncs all playlists at once
   - Updated `RefreshPlaylist` and organize to set `last_synced_at`

2. **Frontend:**
   - Toast component and context for notifications
   - Dashboard badge showing new songs waiting
   - Playlists page with Sync All button
   - Per-playlist new song indicators
   - Success/error toasts for all sync operations

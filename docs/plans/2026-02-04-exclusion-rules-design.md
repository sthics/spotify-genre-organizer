# Exclusion Rules Feature Design

## Overview

Allow users to block artists and songs from appearing in any playlist created by the app. Rules are persistent and apply retroactively to existing managed playlists.

## Requirements

- Block artists with scope choice: "all appearances" (including features) or "primary artist only"
- Block individual songs (from contextual UI only, not search)
- Retroactive enforcement: show preview of affected playlists, then remove on confirmation
- Filters apply to both auto-organize and Build a Crate flows
- Management UI in settings to view/add/remove exclusion rules

## Data Model

```sql
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

CREATE INDEX idx_exclusion_rules_user ON exclusion_rules(user_id);

ALTER TABLE exclusion_rules ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can manage own exclusion rules"
  ON exclusion_rules FOR ALL
  USING (user_id = current_setting('app.user_id', true));
```

## API Endpoints

### `GET /api/exclusions`
Returns all exclusion rules for the authenticated user.

**Response:**
```json
{
  "rules": [
    {
      "id": "uuid",
      "type": "artist",
      "spotify_id": "abc123",
      "name": "The Beatles",
      "scope": "all_appearances",
      "created_at": "2026-02-04T..."
    }
  ]
}
```

### `POST /api/exclusions`
Creates a new exclusion rule and returns impact preview.

**Request:**
```json
{
  "type": "artist",
  "spotify_id": "abc123",
  "name": "The Beatles",
  "scope": "all_appearances"
}
```

**Response:**
```json
{
  "rule_id": "uuid",
  "preview": {
    "total_songs_affected": 12,
    "affected_playlists": [
      {
        "playlist_id": "xyz",
        "name": "Classic Rock",
        "songs_to_remove": [
          { "id": "track1", "name": "Hey Jude" },
          { "id": "track2", "name": "Let It Be" }
        ]
      }
    ]
  }
}
```

### `POST /api/exclusions/{rule_id}/apply`
Confirms retroactive removal from playlists.

**Response:**
```json
{
  "playlists_updated": 3,
  "songs_removed": 12
}
```

### `DELETE /api/exclusions/{rule_id}`
Removes an exclusion rule. Does not restore previously removed songs.

## Backend Logic

### Filtering Integration

Exclusion filtering happens after fetching/enriching songs, before genre scoring:

```
FetchAllLikedSongs()
  -> EnrichSongsWithGenres()
  -> ApplyExclusionRules()  <- NEW
  -> ScoreGenres() / AnalyzeLibraryGenres()
  -> Create playlists
```

### ApplyExclusions Function

```go
func ApplyExclusions(songs []models.Song, rules []models.ExclusionRule) []models.Song {
    // Build lookup maps for O(1) checks
    blockedSongs := map[string]bool{}
    blockedArtistsAll := map[string]bool{}
    blockedArtistsPrimary := map[string]bool{}

    for _, rule := range rules {
        switch rule.Type {
        case "song":
            blockedSongs[rule.SpotifyID] = true
        case "artist":
            if rule.Scope == "all_appearances" {
                blockedArtistsAll[rule.SpotifyID] = true
            } else {
                blockedArtistsPrimary[rule.SpotifyID] = true
            }
        }
    }

    // Filter songs
    var result []models.Song
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

### Impact Preview Function

```go
func PreviewExclusionImpact(userID string, rule ExclusionRule) ([]AffectedPlaylist, error) {
    // 1. Get all playlist_overrides for user
    overrides, err := db.GetPlaylistOverrides(userID)

    // 2. For each playlist, fetch tracks from Spotify
    var affected []AffectedPlaylist
    for _, override := range overrides {
        tracks, err := spotify.GetPlaylistTracks(override.PlaylistSpotifyID)

        // 3. Check each track against the rule
        var matches []Track
        for _, track := range tracks {
            if matchesRule(track, rule) {
                matches = append(matches, track)
            }
        }

        if len(matches) > 0 {
            affected = append(affected, AffectedPlaylist{
                PlaylistID: override.PlaylistSpotifyID,
                Name: override.CustomName,
                SongsToRemove: matches,
            })
        }
    }

    return affected, nil
}
```

## Frontend Components

### BlockButton

Small contextual button that appears on hover over artist/song rows:

```tsx
interface BlockButtonProps {
  type: 'artist' | 'song';
  spotifyId: string;
  name: string;
  imageUrl?: string;
  onBlock: () => void;
}
```

### BlockArtistModal

Two-step modal:
1. Scope selection (all appearances vs primary only)
2. Impact preview with confirm/cancel

Props:
```tsx
interface BlockArtistModalProps {
  artist: { id: string; name: string; imageUrl?: string };
  songCount?: number;
  isOpen: boolean;
  onClose: () => void;
  onBlocked: () => void;
}
```

### BlockSongModal

Simple confirmation modal (no scope choice for songs).

### ImpactPreviewModal

Shows affected playlists with song counts, confirm/skip buttons.

### ExclusionsSettings

New tab in settings page with:
- List of blocked artists (with unblock action)
- List of blocked songs (with unblock action)
- "+ Add" button that opens artist search modal
- Empty states for each section

### ArtistSearchModal

Search input that queries Spotify API, shows results, clicking result opens BlockArtistModal.

## UI/UX Details

### Contextual Blocking Locations

1. **Build a Crate** (`/custom-playlist`) - When browsing genres, show artist breakdown with block option
2. **Playlist review** - After creating playlists, show track list with block option per song/artist
3. **Managed playlists view** - When viewing existing playlists, allow blocking from there

### Toast Messages

| Action | Message |
|--------|---------|
| Rule created, no playlists affected | "{name} blocked" |
| Rule created + removal confirmed | "{name} blocked - {n} songs removed from {m} playlists" |
| Rule created, removal skipped | "{name} blocked - Songs kept in existing playlists" |
| Rule deleted | "{name} unblocked" |
| Error | "Couldn't block {type}. Try again." |

### Loading States

- Search results: Skeleton rows
- Impact preview: Spinner with "Checking your playlists..."
- Removal in progress: "Removing songs..." with optional progress

### Edge Cases

| Scenario | Behavior |
|----------|----------|
| Block artist with 0 songs in library | Allow - blocks future additions |
| Block song not in managed playlists | Save rule, skip preview, show simple toast |
| Spotify API rate limited | Queue remaining, show partial success message |
| Block then immediately unblock | Rule deleted, removed songs stay removed |

## File Structure

```
Backend:
  internal/models/exclusion_rule.go      # New model
  internal/database/exclusion_rules.go   # CRUD operations
  internal/filters/exclusions.go         # ApplyExclusions function
  internal/api/handlers/exclusions.go    # HTTP handlers
  supabase/migrations/003_exclusion_rules.sql

Frontend:
  src/components/blocking/BlockButton.tsx
  src/components/blocking/BlockArtistModal.tsx
  src/components/blocking/BlockSongModal.tsx
  src/components/blocking/ImpactPreviewModal.tsx
  src/components/blocking/ArtistSearchModal.tsx
  src/app/settings/exclusions/page.tsx   # Or tab in existing settings
  src/lib/api.ts                         # Add exclusion API functions
```

## Future Considerations

This is Feature #1 of the Smart Filters & Rules system. The remaining features to design:

2. **Priority/routing rules** - "Always put The Beatles in Classic Rock, not Pop"
3. **Date-based filters** - "Only songs from 2020+"
4. **Popularity filters** - "Only deep cuts (popularity < 50)"
5. **Duration filters** - "No songs under 2 minutes"
6. **Sort by most liked artists** - Order playlist by artist frequency

These will share the filtering infrastructure built for exclusions.

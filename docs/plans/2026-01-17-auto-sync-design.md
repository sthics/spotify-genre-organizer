# Auto-Sync Feature Design

**Date:** 2026-01-17
**Status:** Approved

## Overview

On-demand sync with smart detection. When users open the app, we check for new liked songs and show them what's waiting to be organized. Users control when to sync.

## User Flow

```
User opens dashboard
    ↓
API fetches: current library count + oldest sync timestamp
    ↓
If new songs exist:
    • Badge appears on "Manage My Crates" button
    • Text below: "X songs waiting to be organized"
    ↓
User clicks through to Playlists page
    ↓
Header shows: "Sync All Crates" button (with total count)
Each playlist card shows: "+X new" indicator if that genre has matches
    ↓
User syncs (all or individual)
    ↓
Inline spinner → Toast notification → Updated counts
```

## Dashboard UI

**When new songs exist:**

- Badge on "Manage My Crates" button (orange pill with count, inside button on right)
- Helper text below button: "X songs waiting to be organized"
- Badge uses `accent-orange` (#e85d04) with subtle fade+scale entrance animation
- Helper text uses `text-muted`, `text-sm`

**When no new songs:**

- No badge, no helper text
- Button appears as normal

## Playlists Page UI

**Header:**

- "Sync All (X new)" primary button replaces "+ New Organize" when new songs exist
- When no new songs: reverts to "+ New Organize"

**Playlist cards:**

- "+X new" indicator in top-right, using genre's accent color
- Only appears on playlists with matching new songs
- Expanded state: "↻ Refresh" becomes "↻ Sync X new songs"

**After sync:**

- Indicators disappear from synced playlists
- Header reverts when all synced

## Sync Behavior

### Small Syncs (<50 songs)

- Runs inline
- Button shows spinning vinyl icon
- Toast on completion

### Large Syncs (50+ songs)

- Runs as background job (same architecture as organize)
- Toast on start (rotating messages):
  1. "Sorting your new vinyl..."
  2. "Filing X tracks into your crates..."
  3. "Spinning up your new additions..."
  4. "Dropping new records into your crates..."
  5. "X new tracks heading to their crates..."
- User can navigate away
- Toast on completion: "Sync complete — X songs sorted"

### Progress Display (optional)

- If user stays on playlists page during large sync
- Small text under header: "Syncing... 45/127 songs processed"

## Toast Notifications

### Success

| Scenario | Message |
|----------|---------|
| Single playlist | "Added 4 songs to Rock" |
| Multi-playlist | "5 crates updated • 23 songs added" |
| Large sync complete | "Sync complete — 127 songs sorted" |
| No new matches | "All crates up to date" |

### Errors

| Scenario | Message |
|----------|---------|
| Single playlist fails | "Couldn't sync Rock — Retry" |
| Partial batch fail | "Synced 4 crates, 1 failed — Retry" |
| Auth expired | "Session expired — Log in again" |

**Toast styling:**

- Success: dark background, subtle green left border
- Error: dark background, red left border
- Position: bottom-right
- Auto-dismiss: 4 seconds (errors persist until dismissed)
- "Retry" is clickable link within toast

## Backend

### New Endpoint: Sync Status

```
GET /api/library/sync-status
```

Response:
```json
{
  "new_songs_count": 12,
  "oldest_sync_at": "2026-01-15T10:30:00Z",
  "playlists": [
    { "spotify_id": "abc123", "genre": "Rock", "new_count": 4 },
    { "spotify_id": "def456", "genre": "Electronic", "new_count": 3 }
  ]
}
```

**Logic:**

1. Find user's oldest `last_synced_at` from `playlist_overrides`
2. Fetch liked songs added after that timestamp
3. Determine genre for each new song
4. Count matches per playlist

### New Endpoint: Sync All

```
POST /api/playlists/sync-all
```

Response:
```json
{
  "playlists_updated": 5,
  "songs_added": 23
}
```

For large syncs (50+ songs), returns job ID instead for polling.

### Database

Uses existing `playlist_overrides.last_synced_at` column. Ensure it updates on:

- Manual refresh
- Individual sync
- Sync all
- Full organize

## New Song Tracking

**Baseline:** Songs liked after oldest playlist sync timestamp

**Why oldest:** Shows maximum possible new songs — "these haven't been organized yet"

**Reset triggers:**

- Full "Organize My Library" resets all timestamps
- Individual sync resets that playlist's timestamp
- "Sync All" resets all timestamps

## Edge Cases

**First-time user (no playlists):**

- No badge on dashboard
- Feature activates after first organize

**User runs full organize:**

- All `last_synced_at` timestamps update
- New song count resets to 0

**User deletes a playlist:**

- Oldest sync recalculates from remaining playlists

**All playlists up-to-date:**

- No badge, normal state
- Manual refresh still available

## Implementation Notes

- Job persistence: in-memory for MVP (matches current organize)
- Future: migrate to database-backed jobs
- Sync status check happens on dashboard mount (same as library count)

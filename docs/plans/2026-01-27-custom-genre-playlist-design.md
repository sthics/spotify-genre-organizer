# Custom Genre Playlist Feature Design

## Overview

Add a "Build a Crate" feature allowing users to create playlists from specific genres/sub-genres instead of organizing their entire library at once.

## User Flow

```
Dashboard
    ↓ Click "Build a Crate"
/custom-playlist
    ↓ Select genres (hierarchical picker)
/custom-playlist/configure
    ↓ Choose combined vs separate, name playlist(s)
/processing
    ↓ Watch creation progress
/success
    ↓ View results, build another
```

## Frontend Changes

### 1. Dashboard Addition

Add new button below "Manage My Crates":

```
┌─────────────────────────────────────┐
│  🎵  Build a Crate                  │
│      Pick specific genres           │
└─────────────────────────────────────┘
```

Secondary style button, links to `/custom-playlist`.

### 2. Genre Selection Page (`/custom-playlist`)

**Header:**
- Title: "Build a Crate" (Instrument Serif)
- Subtitle: "{count} liked songs • Analyzed {time} ago" with refresh button
- Back arrow to dashboard

**Genre Picker (accordion-style):**
- Parent genres as expandable rows with song count badges
- Chevron icon to expand/collapse
- When expanded, sub-genres appear indented as selectable chips
- Checkbox on parent = select all sub-genres under it
- Individual sub-genre checkboxes for granular control
- Genres with 0 songs shown muted/disabled

**Sticky Bottom Bar:**
- Left: Selected genres as removable orange chips
- Center: Total song count
- Right: "Next →" button (disabled until selection made)

**Behavior:**
- On mount: fetch `/api/library/genres`
- If no cached data or stale (>24h): show loading state, trigger analysis
- Refresh button invalidates cache and re-analyzes

### 3. Playlist Configuration Page (`/custom-playlist/configure`)

**Header:**
- Title: "Almost there" (Instrument Serif)
- Selected genres shown as small chips

**Two selectable cards:**

**Card 1: "One Combined Playlist"**
- Single vinyl icon
- "All {count} songs in one playlist"
- Text input for custom name (placeholder: "My Heavy Stuff")
- Orange border when selected

**Card 2: "Separate Playlists"**
- Stacked records icon
- "{n} playlists, one per genre"
- Preview list: "metalcore (52) • death metal (61) • hardcore (34)"
- Names use user's template from settings

**Bottom:**
- "Back" link
- "Create Playlist(s)" primary button

**Edge case:** If only one sub-genre selected, skip this page and go directly to creation.

### 4. Processing & Success

Reuse existing pages:
- `/processing?job={job_id}` — same spinning vinyl, stages update for custom flow
- `/success?job={job_id}` — shows created playlist(s), adds "Build Another Crate" button

## Backend Changes

### New Endpoints

#### `GET /api/library/genres`

Returns user's genre breakdown from cache or fresh analysis.

**Response:**
```json
{
  "parent_genres": [
    {
      "name": "Metal",
      "count": 89,
      "sub_genres": [
        { "name": "metalcore", "count": 52 },
        { "name": "death metal", "count": 24 },
        { "name": "nu metal", "count": 13 }
      ]
    },
    {
      "name": "Rock",
      "count": 156,
      "sub_genres": [
        { "name": "alternative rock", "count": 78 },
        { "name": "indie rock", "count": 45 },
        { "name": "hard rock", "count": 33 }
      ]
    }
  ],
  "analyzed_at": "2026-01-27T10:30:00Z",
  "total_songs": 412
}
```

**Logic:**
1. Check `genre_analysis_cache` for user
2. If exists and < 24h old, return cached data
3. Otherwise, fetch liked songs, analyze genres, cache results, return

**Query param:** `?refresh=true` to force re-analysis

#### `POST /api/custom-playlist`

Creates custom playlist(s) from selected genres.

**Request:**
```json
{
  "sub_genres": ["metalcore", "death metal"],
  "mode": "combined",
  "name": "My Heavy Stuff"
}
```

Or for separate mode:
```json
{
  "sub_genres": ["metalcore", "death metal", "hardcore"],
  "mode": "separate"
}
```

**Response:**
```json
{
  "job_id": "abc123"
}
```

**Logic:**
1. Fetch liked songs (or use recent cache)
2. Filter to songs matching selected sub-genres
3. If mode = "combined": create single playlist with custom name
4. If mode = "separate": create one playlist per sub-genre using name template
5. Return job_id for polling

### Database Schema

**New table: `genre_analysis_cache`**

```sql
CREATE TABLE genre_analysis_cache (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id TEXT NOT NULL UNIQUE REFERENCES users(spotify_id),
  analysis_data JSONB NOT NULL,
  analyzed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_genre_analysis_cache_user_id ON genre_analysis_cache(user_id);
```

**analysis_data structure:**
```json
{
  "parent_genres": [...],
  "total_songs": 412,
  "song_genre_map": {
    "spotify:track:abc": ["metalcore", "hardcore"],
    "spotify:track:def": ["death metal"]
  }
}
```

The `song_genre_map` allows quick filtering when creating custom playlists without re-fetching artist genres.

## Caching Strategy for Returning Users

1. **First visit:** Full analysis, cache results
2. **Return visit (<24h):** Instant genre list from cache
3. **Return visit (>24h):** Show stale data immediately, offer refresh
4. **Manual refresh:** Invalidate cache, re-analyze
5. **After "Organize My Library":** Update cache as side effect

**UI indicators:**
- "Analyzed 2 hours ago" timestamp
- Refresh icon button
- "12 new songs since last analysis" badge (compare liked songs count)

## File Structure

```
frontend/src/app/
  custom-playlist/
    page.tsx           # Genre selection
    configure/
      page.tsx         # Combine vs separate choice

frontend/src/components/
  GenrePicker.tsx      # Accordion genre selector
  GenreChip.tsx        # Selectable/removable genre chip

backend/internal/api/handlers/
  genres.go            # GET /api/library/genres
  custom_playlist.go   # POST /api/custom-playlist
```

## Open Questions

None — design approved through brainstorming session.

## Implementation Order

1. Database migration for `genre_analysis_cache`
2. Backend: `GET /api/library/genres` endpoint
3. Backend: `POST /api/custom-playlist` endpoint
4. Frontend: GenrePicker component
5. Frontend: `/custom-playlist` page
6. Frontend: `/custom-playlist/configure` page
7. Frontend: Dashboard button addition
8. Update processing/success pages to handle custom playlist jobs

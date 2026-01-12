# Spotify Genre Organizer - MVP Design

**Date:** 2026-01-12
**Status:** Approved
**Target:** 1-month functional MVP

---

## Overview

A web app that automatically organizes Spotify liked songs into genre-based playlists. Users connect their Spotify account, choose how many playlists they want, and get organized playlists created in their library.

### The Elevator Pitch

> "2,000 liked songs. Zero organization. Sound familiar? Connect Spotify, pick your playlist count, and we'll sort your music in seconds."

---

## MVP Scope

### What's IN

| Feature | Description |
|---------|-------------|
| Spotify OAuth | Login with Spotify, request library-read and playlist-modify scopes |
| Fetch liked songs | Pull all user's liked songs via Spotify API |
| Genre detection | Use Spotify's artist/track genre tags |
| Smart consolidation | Map 100s of micro-genres → ~15-20 parent categories |
| Playlist count slider | User chooses 1-50 playlists |
| Create/update playlists | Write playlists directly to user's Spotify account |
| Replace toggle | Option to update existing vs. create fresh (default: update) |

### What's OUT (v1.1+)

- Auto-sync (automatic updates when new songs liked)
- Last.fm integration (enhanced genre detection)
- Recommendations (discover similar songs)
- Custom genre mappings (user-defined categories)
- Analytics (listening trends, stats)
- Premium tier / payments
- Mobile app

---

## Technical Architecture

### Stack (Simplified for MVP)

| Layer | Technology |
|-------|------------|
| Frontend | Next.js 14, TypeScript, Tailwind CSS |
| Backend | Go 1.21+ with Gin |
| Database | PostgreSQL via Supabase |
| Auth | Supabase Auth (handles Spotify OAuth tokens) |
| Hosting | Vercel (frontend) + Railway (backend) |

### What We're Cutting from README Plans

| Original Plan | MVP Decision |
|---------------|--------------|
| Redis caching | Cut - premature optimization |
| Full Clean Architecture (DDD) | Simplify - fewer layers |
| JWT + encrypted tokens | Simplify - Supabase handles this |
| Docker + CI/CD | Defer - deploy direct |
| Integration & E2E tests | Defer - unit tests only |

### Database Schema (MVP)

```sql
-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  spotify_id VARCHAR(255) UNIQUE NOT NULL,
  display_name VARCHAR(255),
  email VARCHAR(255),
  access_token TEXT,
  refresh_token TEXT,
  token_expires_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Organize jobs (track requests)
CREATE TABLE organize_jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  playlist_count INTEGER NOT NULL,
  replace_existing BOOLEAN DEFAULT TRUE,
  status VARCHAR(50) DEFAULT 'pending', -- pending, processing, completed, failed
  songs_processed INTEGER DEFAULT 0,
  total_songs INTEGER,
  playlists_created JSONB, -- [{name, spotify_id, song_count}]
  error_message TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  completed_at TIMESTAMP
);

-- Genre mappings (seeded, not user-editable in MVP)
CREATE TABLE genre_mappings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  micro_genre VARCHAR(255) NOT NULL, -- e.g., "indie rock", "garage rock"
  parent_genre VARCHAR(100) NOT NULL, -- e.g., "Rock"
  created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_spotify_id ON users(spotify_id);
CREATE INDEX idx_organize_jobs_user_id ON organize_jobs(user_id);
CREATE INDEX idx_genre_mappings_micro ON genre_mappings(micro_genre);
```

### API Endpoints (MVP)

```
# Auth
GET  /api/auth/login          → Redirect to Spotify OAuth
GET  /api/auth/callback       → Handle OAuth callback
GET  /api/auth/me             → Get current user
POST /api/auth/logout         → Clear session

# Organize
POST /api/organize            → Start organize job
GET  /api/organize/:id        → Get job status (polling)
```

### Genre Consolidation Logic

Map Spotify's micro-genres to parent categories:

```go
var genreMap = map[string]string{
    // Rock
    "indie rock":       "Rock",
    "alternative rock": "Rock",
    "garage rock":      "Rock",
    "classic rock":     "Rock",
    "hard rock":        "Rock",
    "punk rock":        "Rock",

    // Electronic
    "edm":              "Electronic",
    "house":            "Electronic",
    "techno":           "Electronic",
    "dubstep":          "Electronic",
    "drum and bass":    "Electronic",
    "ambient":          "Electronic",

    // Hip-Hop
    "hip hop":          "Hip-Hop",
    "rap":              "Hip-Hop",
    "trap":             "Hip-Hop",
    "conscious hip hop":"Hip-Hop",

    // ... etc for ~15-20 parent categories
}
```

---

## UI/UX Design

### Aesthetic Direction: "Vinyl Record Store Meets Digital"

| Element | Choice |
|---------|--------|
| **Tone** | Warm, tactile, slightly retro-modern |
| **Typography** | Display: Instrument Serif / Body: IBM Plex Sans |
| **Colors** | Dark base (#1a1a1a), warm cream (#f5f0e6), accent orange (#e85d04) |
| **Texture** | Subtle noise/grain overlay, soft shadows |
| **Motion** | Vinyl spin animations, staggered reveals, bounce entrances |

### The Memorable Moment

When organizing completes, playlist cards "drop" onto the screen like vinyl records landing in crates - satisfying, physical, delightful.

---

## Screen Designs

### Screen 1: Landing Page

```
┌─────────────────────────────────────────────────────────────┐
│  [dark background with subtle grain texture]                │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  SPOTIFY GENRE                                       │   │
│  │  ORGANIZER          [spinning vinyl icon]            │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                                                       │   │
│  │   "2,000 liked songs.                                │   │
│  │    Zero organization.                                │   │
│  │    Sound familiar?"                                  │   │
│  │                                                       │   │
│  │   ┌─────────────────────────────────┐                │   │
│  │   │  Connect with Spotify  ●────────│                │   │
│  │   └─────────────────────────────────┘                │   │
│  │                                                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐              │
│  │ Analyze   │  │ Organize  │  │ Enjoy     │              │
│  │ your      │  │ into      │  │ your      │              │
│  │ library   │  │ playlists │  │ music     │              │
│  └───────────┘  └───────────┘  └───────────┘              │
└─────────────────────────────────────────────────────────────┘
```

**Details:**
- Hero text: Instrument Serif, large, cream-colored
- Button: Orange (#e85d04), subtle hover glow
- Vinyl icon: CSS animation, 8s rotation
- Value prop cards: Staggered fade-in on scroll

---

### Screen 2: Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│  [dark bg + grain]                                          │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Hey, Marcus.                        [avatar] ⚙️     │   │
│  │  You've got 1,247 liked songs.                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                                                       │   │
│  │  ┌─────────────────────────────────────────────┐     │   │
│  │  │  HOW MANY PLAYLISTS?                         │     │   │
│  │  │                                              │     │   │
│  │  │     1  ●━━━━━━━━━━━━━○━━━━━━━━━━━━━━  50    │     │   │
│  │  │                      ↑                       │     │   │
│  │  │                     12                       │     │   │
│  │  │                                              │     │   │
│  │  │  [vinyl icon] ~104 songs per playlist        │     │   │
│  │  └─────────────────────────────────────────────┘     │   │
│  │                                                       │   │
│  │  ┌─────────────────────────────────────────────┐     │   │
│  │  │  ◉ Update existing playlists                 │     │   │
│  │  │    (replaces songs in "Rock by Organizer")   │     │   │
│  │  │                                              │     │   │
│  │  │  ○ Create fresh playlists                    │     │   │
│  │  │    (keeps your old ones, makes new)          │     │   │
│  │  └─────────────────────────────────────────────┘     │   │
│  │                                                       │   │
│  │       ┌─────────────────────────────────┐            │   │
│  │       │    ORGANIZE MY LIBRARY    ◉     │            │   │
│  │       └─────────────────────────────────┘            │   │
│  │                                                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Details:**
- Slider: Snaps with tactile feel, live-updates song count estimate
- Radio buttons: Styled as vinyl grooves
- Main button: 3D press-down animation on click
- Settings cog: Opens drawer for logout only (MVP)

---

### Screen 3: Processing

```
┌─────────────────────────────────────────────────────────────┐
│  [dark bg + grain]                                          │
│                                                             │
│                    ┌─────────────┐                          │
│                    │             │                          │
│                    │   ◉━━━━━    │   [vinyl record          │
│                    │     ╲       │    spinning,             │
│                    │      ╲      │    tonearm moving]       │
│                    │             │                          │
│                    └─────────────┘                          │
│                                                             │
│              "Analyzing your library..."                    │
│                                                             │
│                    ━━━━━━━━━━━━━━━━━                        │
│                    [progress bar - cream fill]              │
│                    412 / 1,247 songs                        │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Genres discovered:                                  │   │
│  │                                                       │   │
│  │  ┌────────┐ ┌──────┐ ┌─────────┐ ┌──────┐           │   │
│  │  │  Rock  │ │ Jazz │ │Electronic│ │ Pop  │           │   │
│  │  └────────┘ └──────┘ └─────────┘ └──────┘           │   │
│  │  ┌───────┐ ┌────────┐ ┌───────┐                      │   │
│  │  │Hip-Hop│ │  Indie │ │ R&B   │  ...                 │   │
│  │  └───────┘ └────────┘ └───────┘                      │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

**Details:**
- Vinyl: Spins continuously (CSS, 3s rotation)
- Tonearm: Moves from outer edge to center as progress increases
- Genre tags: Pop in one-by-one with bounce animation
- Progress states: "Analyzing..." → "Sorting..." → "Creating..."

---

### Screen 4: Success

```
┌─────────────────────────────────────────────────────────────┐
│  [dark bg + grain + subtle confetti particles]              │
│                                                             │
│                      ✓ Done!                                │
│            "12 playlists ready to play"                     │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  [playlist cards drop in like vinyl into crates]     │   │
│  │                                                       │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐             │   │
│  │  │ ♫ Rock   │ │ ♫ Jazz   │ │♫Electronic│             │   │
│  │  │ 156 songs│ │  42 songs│ │ 203 songs│             │   │
│  │  │  [Open]  │ │  [Open]  │ │  [Open]  │             │   │
│  │  └──────────┘ └──────────┘ └──────────┘             │   │
│  │                                                       │   │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐             │   │
│  │  │ ♫ Pop    │ │ ♫ Hip-Hop│ │ ♫ Indie  │             │   │
│  │  │ 312 songs│ │  89 songs│ │ 127 songs│             │   │
│  │  │  [Open]  │ │  [Open]  │ │  [Open]  │             │   │
│  │  └──────────┘ └──────────┘ └──────────┘             │   │
│  │                                                       │   │
│  │            [+ 6 more - scroll to see]                │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│       ┌─────────────────────────────────────────────────┐  │
│       │         OPEN ALL IN SPOTIFY         ●           │  │
│       └─────────────────────────────────────────────────┘  │
│                                                             │
│                    [Organize Again]                         │
└─────────────────────────────────────────────────────────────┘
```

**Details:**
- Cards: Drop from top, staggered 100ms, bounce on land
- Card hover: Lifts 4px, shadow deepens
- Each card has genre-colored accent bar at top
- Confetti: Subtle particles drift for 3 seconds
- "Open" links directly to Spotify playlist

---

## User Flow Summary

```
Landing Page
     │
     ▼
[Connect with Spotify] ──→ Spotify OAuth
     │
     ▼
Dashboard
  • See liked song count
  • Choose playlist count (1-50)
  • Toggle: update vs. create fresh
     │
     ▼
[Organize My Library]
     │
     ▼
Processing Screen
  • Vinyl animation
  • Progress bar
  • Genre tags appear live
     │
     ▼
Success Screen
  • Playlist cards drop in
  • Open in Spotify buttons
  • Organize again option
```

---

## Implementation Priorities

### Phase 1: Foundation
- [ ] Set up Next.js + Tailwind project
- [ ] Set up Go + Gin backend
- [ ] Configure Supabase (database + auth)
- [ ] Implement Spotify OAuth flow

### Phase 2: Core Logic
- [ ] Fetch liked songs from Spotify API
- [ ] Build genre consolidation mapping
- [ ] Create playlist generation logic
- [ ] Implement job status tracking

### Phase 3: Frontend
- [ ] Landing page with animations
- [ ] Dashboard with slider + options
- [ ] Processing screen with vinyl animation
- [ ] Success screen with card drop animation

### Phase 4: Polish
- [ ] Error handling (token refresh, API limits)
- [ ] Loading states
- [ ] Mobile responsiveness
- [ ] Basic testing

---

## Success Metrics (MVP)

| Metric | Target |
|--------|--------|
| Time to organize 1,000 songs | < 30 seconds |
| Successful organize rate | > 95% |
| User returns within 7 days | > 20% |

---

## Open Questions (Post-MVP)

1. How to handle songs with no genre data?
2. What's the right number of default parent genres?
3. Should we name playlists "[Genre] by Organizer" or let user choose?

---

## Appendix: Color Palette

```css
:root {
  --bg-dark: #1a1a1a;
  --bg-card: #252525;
  --text-cream: #f5f0e6;
  --text-muted: #8a8a8a;
  --accent-orange: #e85d04;
  --accent-orange-hover: #ff6b0a;
  --success-green: #2d936c;
}
```

## Appendix: Font Stack

```css
:root {
  --font-display: 'Instrument Serif', Georgia, serif;
  --font-body: 'IBM Plex Sans', -apple-system, sans-serif;
}
```

---

*Design approved: 2026-01-12*

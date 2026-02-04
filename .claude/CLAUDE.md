# Spotify Genre Organizer

A web app that organizes your Spotify liked songs into genre-based playlists.

## Project Structure

```
backend/     # Go 1.23 + Gin framework
frontend/    # Next.js 14 + TypeScript + Tailwind
supabase/    # Database migrations
```

## Build & Run

### Backend
```bash
cd backend
go mod download
go run cmd/server/main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

### Tests
```bash
# Backend
cd backend && go test -v ./...

# Frontend
cd frontend && npm run lint
```

## Code Style

### Go (Backend)
- Use `goimports` for formatting
- Follow standard Go conventions
- Group imports: stdlib, third-party, internal
- Use `errors.Is()` / `errors.As()` for error checking
- Add comments for exported functions

### TypeScript (Frontend)
- Use TypeScript strict mode
- Follow existing component patterns
- Tailwind for styling (use existing design tokens)

## Design Tokens

```
Colors:
  bg-dark: #1a1a1a
  bg-card: (dark card background)
  text-cream: #f5f0e6
  text-muted: (muted text)
  accent-orange: (primary action color)

Typography:
  font-display: Instrument Serif (headings)
  font-body: IBM Plex Sans (body text)
```

## Architecture

### Backend Packages
- `internal/api/handlers/` - HTTP handlers
- `internal/spotify/` - Spotify API integration
- `internal/genres/` - Genre consolidation logic
- `internal/organizer/` - Playlist creation orchestration
- `internal/database/` - Supabase data layer
- `internal/models/` - Domain entities

### Frontend Structure
- `src/app/` - Next.js App Router pages
- `src/components/` - Reusable UI components
- `src/lib/api.ts` - Backend API client
- `src/hooks/` - Custom React hooks
- `src/contexts/` - React contexts

## Database

Uses Supabase (PostgreSQL). Migrations in `supabase/migrations/`.

Key tables: `users`, `user_settings`, `playlist_overrides`, `exclusion_rules`

## Worktrees

Use `.worktrees/` for feature branches.

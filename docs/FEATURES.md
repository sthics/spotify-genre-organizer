# Spotify Genre Organizer - Feature List

## ✅ Implemented Features

### 🔐 Authentication
- **Spotify OAuth 2.0 Login** - Secure login via Spotify account
- **Session Management** - Cookie-based session with access token storage
- **Logout** - Clear session and redirect to home

---

### 🎵 Core Organization
- **Automatic Genre Detection** - Analyzes artist genres from Spotify metadata
- **Smart Genre Grouping** - Groups 4,000+ Spotify sub-genres into ~20 parent categories (Rock, Electronic, Hip-Hop, Jazz, Pop, Metal, Folk, etc.)
- **Weighted Genre Scoring** - Uses algorithmic scoring across all artist genres for accurate classification
- **Liked Songs Import** - Fetches user's entire liked songs library from Spotify
- **Configurable Playlist Count** - Slider to choose 1-50 genre playlists
- **Merge Small Genres** - Automatically merges low-count genres into "Other"

---

### 📋 Playlist Creation
- **Create Genre Playlists** - Creates playlists on Spotify with organized songs
- **Replace or Create New** - Option to update existing playlists or create fresh ones
- **Custom Naming Templates** - Configurable patterns using `{genre}` and `{year}` tokens
  - Example: `{genre} by Organizer` → "Rock by Organizer"
- **Custom Description Templates** - Same token system for playlist descriptions
- **Real-time Progress Tracking** - Processing page with stage updates and progress bar

---

### 🎛️ Playlist Management (Crates)
- **View All Organized Playlists** - Grid view of all created genre playlists
- **Color-coded by Genre** - Visual genre tags (Rock=orange, Electronic=cyan, etc.)
- **Song Count Display** - Shows number of songs per playlist
- **Expandable Details** - Accordion-style cards with actions
- **Edit Playlist Details** - Rename or update description per-playlist
- **Refresh/Sync Playlist** - Re-sync songs from liked library to playlist
- **Delete Playlist** - Remove from Spotify
- **Open in Spotify** - Direct link to playlist on Spotify

---

### ⚙️ Settings
- **Playlist Name Pattern** - Global template for new playlist names
- **Description Pattern** - Global template for descriptions
- **Live Preview** - Real-time preview of how playlists will appear
- **Database-backed Settings** - Persisted per-user in Supabase

---

### 🛡️ Security & Performance
- **Rate Limiting** - 100 requests/minute per IP
- **CORS Protection** - Whitelisted frontend origins only
- **Row-Level Security (RLS)** - Supabase database-level access control
- **HttpOnly Cookies** - Prevents XSS token theft

---

### 🎨 Frontend UI/UX
- **Dark Theme** - Spotify-inspired dark aesthetic
- **Vinyl/Record Theming** - Spinning vinyl icons, "Crates" terminology
- **Responsive Design** - Works on desktop and mobile
- **Loading States** - Skeleton loaders and spinners
- **Error Handling** - Error boundary component

---

## 📍 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/auth/login` | Initiate Spotify OAuth |
| GET | `/api/auth/callback` | OAuth callback handler |
| GET | `/api/auth/me` | Get current user profile |
| POST | `/api/auth/logout` | End session |
| POST | `/api/organize` | Start organization job |
| GET | `/api/organize/:id` | Get job status |
| GET | `/api/library/count` | Get liked songs count |
| GET | `/api/settings` | Get user settings |
| PUT | `/api/settings` | Update settings |
| GET | `/api/playlists` | List managed playlists |
| PATCH | `/api/playlists/:id` | Update playlist details |
| DELETE | `/api/playlists/:id` | Delete playlist |
| POST | `/api/playlists/:id/refresh` | Refresh playlist songs |

---

## 🗺️ Not Yet Implemented (Roadmap from README)

- [ ] Last.fm integration for enhanced genre detection
- [ ] Auto-sync (automatic updates when new songs liked)
- [ ] AI-powered recommendations
- [ ] Analytics dashboard
- [ ] Mobile app (React Native)
- [ ] Collaborative playlists
- [ ] Mood-based organization
- [ ] Apple Music support

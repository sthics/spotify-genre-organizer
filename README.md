# Spotify Genre Organizer

A scalable, production-ready application that automatically organizes your Spotify liked songs into genre-based playlists with smart recommendations.

## 🎯 Features

- **Automatic Genre Organization**: Sorts your liked songs into custom genre playlists
- **Smart Recommendations**: AI-powered song discovery based on your taste
- **Auto-Sync**: Automatically updates playlists when you like new songs
- **Custom Genre Mappings**: Create your own genre categories
- **Last.fm Integration**: Enhanced genre detection (coming soon)
- **Analytics**: Visualize your music taste and trends

## 🏗️ Architecture

### Tech Stack

**Backend:**
- Go 1.21+ with Gin web framework
- Clean Architecture (Domain-driven design)
- PostgreSQL (via Supabase)
- Redis (optional caching)

**Frontend:**
- Next.js 14
- TypeScript
- Tailwind CSS
- Supabase client

**Infrastructure:**
- Docker & Docker Compose
- GitHub Actions CI/CD
- Railway/Fly.io deployment

### Project Structure

```
spotify-genre-organizer/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Business logic (no external dependencies)
│   │   ├── models/                 # Domain entities
│   │   ├── repositories/           # Repository interfaces
│   │   └── services/               # Service interfaces
│   ├── application/                # Use cases & orchestration
│   │   ├── genre_service.go
│   │   ├── sync_service.go
│   │   └── recommendation_service.go
│   ├── infrastructure/             # External dependencies
│   │   ├── spotify/                # Spotify API client
│   │   ├── supabase/               # Database repositories
│   │   └── lastfm/                 # Last.fm client
│   └── api/                        # HTTP layer
│       ├── handlers/               # Request handlers
│       ├── middleware/             # Auth, logging, rate limiting
│       └── routes.go               # Route definitions
├── pkg/                            # Shared utilities
│   ├── logger/
│   └── errors/
├── migrations/                     # Database migrations
├── tests/                          # Integration & E2E tests
└── docker-compose.yml
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Spotify Developer Account
- Supabase Account

### 1. Spotify API Setup

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Add redirect URI: `http://localhost:8080/api/v1/auth/callback`
4. Note your Client ID and Client Secret

### 2. Supabase Setup

1. Create a new project at [Supabase](https://supabase.com)
2. Get your project URL and anon key from Settings > API
3. Run migrations (see below)

### 3. Environment Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/spotify-genre-organizer.git
cd spotify-genre-organizer

# Copy environment template
cp .env.example .env

# Edit .env with your credentials
nano .env
```

### 4. Run with Docker

```bash
# Start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### 5. Run Locally (Development)

```bash
# Install development tools
make install-tools

# Download dependencies
make deps

# Run database migrations
make migrate-up

# Start with hot reload
make run
```

The API will be available at `http://localhost:8080`

## 📊 Database Migrations

```bash
# Create a new migration
make migrate-create name=add_user_preferences

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint
```

## 🔐 Security Features

- **JWT Authentication**: Secure token-based auth
- **Encrypted Tokens**: Spotify tokens encrypted at rest in database
- **Rate Limiting**: 100 requests/minute per IP
- **CORS Protection**: Whitelisted frontend origins only
- **Row-Level Security**: Database-level access control
- **Input Validation**: All user inputs sanitized
- **HTTPS Only**: Forced in production

## 📡 API Endpoints

### Authentication
- `POST /api/v1/auth/login` - Initiate OAuth flow
- `GET /api/v1/auth/callback` - OAuth callback
- `GET /api/v1/me` - Get current user (protected)

### Playlists
- `GET /api/v1/playlists` - List user's playlists (protected)
- `POST /api/v1/playlists/organize` - Organize liked songs (protected)
- `GET /api/v1/playlists/:id` - Get playlist details (protected)
- `DELETE /api/v1/playlists/:id` - Delete playlist (protected)

### Genre Mappings
- `GET /api/v1/genres/mappings` - List genre mappings (protected)
- `POST /api/v1/genres/mappings` - Create mapping (protected)
- `PUT /api/v1/genres/mappings/:id` - Update mapping (protected)
- `DELETE /api/v1/genres/mappings/:id` - Delete mapping (protected)
- `GET /api/v1/genres/distribution` - Get genre distribution (protected)

### Recommendations
- `POST /api/v1/recommendations` - Generate recommendations (protected)

### Sync
- `POST /api/v1/sync/trigger` - Trigger manual sync (protected)
- `GET /api/v1/sync/status/:job_id` - Get sync status (protected)
- `PUT /api/v1/sync/auto-sync` - Update auto-sync setting (protected)

### Health
- `GET /health` - Health check (public)

## 🎨 Frontend Setup (Next.js)

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build
```

## 🚢 Deployment

### Railway

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Deploy
railway up
```

### Fly.io

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Login
fly auth login

# Deploy
fly deploy
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `ENV` | Environment (development/production) | Yes |
| `PORT` | Server port | Yes |
| `SPOTIFY_CLIENT_ID` | Spotify OAuth client ID | Yes |
| `SPOTIFY_CLIENT_SECRET` | Spotify OAuth client secret | Yes |
| `SPOTIFY_REDIRECT_URI` | OAuth redirect URI | Yes |
| `SUPABASE_URL` | Supabase project URL | Yes |
| `SUPABASE_KEY` | Supabase anon key | Yes |
| `JWT_SECRET` | JWT signing secret | Yes |
| `FRONTEND_URL` | Frontend URL for CORS | Yes |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | No |
| `REDIS_URL` | Redis connection URL | No |

## 📈 Monitoring

The application includes:
- Structured logging (zerolog)
- Request/response logging
- Error tracking
- Health check endpoint
- Optional metrics (Prometheus-compatible)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go best practices
- Run `make lint` before committing
- Maintain test coverage above 70%
- Write meaningful commit messages

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Spotify Web API
- Supabase
- Go community
- All contributors

## 📞 Support

- 📧 Email: support@example.com
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/spotify-genre-organizer/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/yourusername/spotify-genre-organizer/discussions)

## 🗺️ Roadmap

- [ ] Last.fm integration for better genre detection
- [ ] Mobile app (React Native)
- [ ] Collaborative playlists
- [ ] Advanced analytics dashboard
- [ ] Mood-based organization
- [ ] Apple Music support
- [ ] Multi-language support
- [ ] Public API

---

Made with ❤️ by [Your Name]
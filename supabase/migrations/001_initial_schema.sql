-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  spotify_id VARCHAR(255) UNIQUE NOT NULL,
  display_name VARCHAR(255),
  email VARCHAR(255),
  access_token TEXT,
  refresh_token TEXT,
  token_expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Organize jobs table
CREATE TABLE organize_jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  playlist_count INTEGER NOT NULL,
  replace_existing BOOLEAN DEFAULT TRUE,
  status VARCHAR(50) DEFAULT 'pending',
  songs_processed INTEGER DEFAULT 0,
  total_songs INTEGER,
  playlists_created JSONB DEFAULT '[]'::jsonb,
  error_message TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  completed_at TIMESTAMPTZ
);

-- Genre mappings table (seeded data)
CREATE TABLE genre_mappings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  micro_genre VARCHAR(255) NOT NULL,
  parent_genre VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_spotify_id ON users(spotify_id);
CREATE INDEX idx_organize_jobs_user_id ON organize_jobs(user_id);
CREATE INDEX idx_organize_jobs_status ON organize_jobs(status);
CREATE INDEX idx_genre_mappings_micro ON genre_mappings(micro_genre);

-- Enable Row Level Security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE organize_jobs ENABLE ROW LEVEL SECURITY;
ALTER TABLE genre_mappings ENABLE ROW LEVEL SECURITY;

-- RLS Policies (service role bypasses these)
CREATE POLICY "Users can view own data" ON users
  FOR SELECT USING (true);

CREATE POLICY "Users can view own jobs" ON organize_jobs
  FOR SELECT USING (true);

CREATE POLICY "Anyone can read genre mappings" ON genre_mappings
  FOR SELECT USING (true);

-- User settings for playlist templates
CREATE TABLE IF NOT EXISTS user_settings (
  user_id TEXT PRIMARY KEY REFERENCES users(spotify_id) ON DELETE CASCADE,
  name_template TEXT NOT NULL DEFAULT '{genre} by Organizer',
  description_template TEXT NOT NULL DEFAULT 'Organized by Spotify Genre Organizer',
  is_premium BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Per-playlist custom overrides
CREATE TABLE IF NOT EXISTS playlist_overrides (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id TEXT NOT NULL REFERENCES users(spotify_id) ON DELETE CASCADE,
  playlist_spotify_id TEXT NOT NULL,
  custom_name TEXT,
  custom_description TEXT,
  genre TEXT NOT NULL,
  last_synced_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, playlist_spotify_id)
);

-- Index for faster lookups
CREATE INDEX idx_playlist_overrides_user ON playlist_overrides(user_id);

-- RLS policies
ALTER TABLE user_settings ENABLE ROW LEVEL SECURITY;
ALTER TABLE playlist_overrides ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can manage own settings"
  ON user_settings FOR ALL
  USING (user_id = current_setting('app.user_id', true));

CREATE POLICY "Users can manage own playlist overrides"
  ON playlist_overrides FOR ALL
  USING (user_id = current_setting('app.user_id', true));

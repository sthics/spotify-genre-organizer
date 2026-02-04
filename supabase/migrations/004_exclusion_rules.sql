-- Exclusion rules for blocking artists/songs from playlists
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

-- Index for FK (PostgreSQL doesn't auto-create)
CREATE INDEX idx_exclusion_rules_user ON exclusion_rules(user_id);

-- RLS
ALTER TABLE exclusion_rules ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can manage own exclusion rules"
  ON exclusion_rules FOR ALL
  USING (user_id = current_setting('app.user_id', true));

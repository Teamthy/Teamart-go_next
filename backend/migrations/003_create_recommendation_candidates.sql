-- Create table for recommendation candidates
CREATE TABLE IF NOT EXISTS recommendation_candidates (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  category TEXT,
  signals JSONB,
  meta JSONB,
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_recommendation_candidates_category ON recommendation_candidates (category);

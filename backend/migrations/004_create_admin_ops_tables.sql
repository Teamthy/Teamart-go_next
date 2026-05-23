-- Create tables for admin operations: disputes, payouts, fraud alerts

CREATE TABLE IF NOT EXISTS disputes (
  id TEXT PRIMARY KEY,
  order_id TEXT,
  user_id BIGINT,
  amount NUMERIC,
  reason TEXT,
  status TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_disputes_status ON disputes (status);

CREATE TABLE IF NOT EXISTS payouts (
  id TEXT PRIMARY KEY,
  creator_id TEXT,
  amount NUMERIC,
  status TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_payouts_status ON payouts (status);

CREATE TABLE IF NOT EXISTS fraud_alerts (
  id TEXT PRIMARY KEY,
  subject_id TEXT,
  score NUMERIC,
  data JSONB,
  created_at TIMESTAMPTZ DEFAULT now()
);

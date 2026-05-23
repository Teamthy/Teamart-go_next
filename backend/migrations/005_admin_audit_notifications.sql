-- Admin audit logs, notifications, and payout approvals

CREATE TABLE IF NOT EXISTS audit_logs (
  id TEXT PRIMARY KEY,
  actor_id BIGINT,
  action TEXT,
  resource_type TEXT,
  resource_id TEXT,
  details JSONB,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_actor ON audit_logs (actor_id);

CREATE TABLE IF NOT EXISTS notifications (
  id TEXT PRIMARY KEY,
  recipient_id TEXT,
  channel TEXT,
  payload JSONB,
  sent_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS payout_approvals (
  id TEXT PRIMARY KEY,
  payout_id TEXT,
  requested_by BIGINT,
  status TEXT,
  notes TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  approved_by BIGINT,
  approved_at TIMESTAMPTZ
);

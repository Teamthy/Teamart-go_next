package admin

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/teamart/commerce-api/config"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

func TestIntegration_AdminFlow(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "1" {
		t.Skip("set INTEGRATION_TEST=1 to run integration tests")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		user := os.Getenv("POSTGRES_USER")
		if user == "" {
			user = "postgres"
		}
		pass := os.Getenv("POSTGRES_PASSWORD")
		if pass == "" {
			pass = "postgres"
		}
		db := os.Getenv("POSTGRES_DB")
		if db == "" {
			db = "teamart"
		}
		host := os.Getenv("POSTGRES_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("POSTGRES_PORT")
		if port == "" {
			port = "5432"
		}
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
	}

	cfg := &config.DatabaseConfig{URL: dbURL}
	log := logger.NewNoop()
	pool, err := database.NewPool(context.Background(), cfg, log)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	// Ensure admin tables exist
	_, err = pool.Exec(context.Background(), `
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
CREATE TABLE IF NOT EXISTS payouts (
  id TEXT PRIMARY KEY,
  creator_id TEXT,
  amount NUMERIC,
  status TEXT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);
CREATE TABLE IF NOT EXISTS fraud_alerts (
  id TEXT PRIMARY KEY,
  subject_id TEXT,
  score NUMERIC,
  data JSONB,
  created_at TIMESTAMPTZ DEFAULT now()
);
`)
	if err != nil {
		t.Fatalf("failed to ensure admin tables: %v", err)
	}

	// Cleanup
	pool.Exec(context.Background(), "TRUNCATE disputes, payouts, fraud_alerts")
	defer pool.Exec(context.Background(), "TRUNCATE disputes, payouts, fraud_alerts")

	repo := NewPostgresAdminRepository(log)

	// Create dispute
	d := Dispute{ID: "d-1", OrderID: "o-1", UserID: int64(123), Amount: 9.99, Reason: "test", Status: DisputeOpen}
	if err := repo.SaveDispute(context.Background(), pool, d); err != nil {
		t.Fatalf("save dispute failed: %v", err)
	}

	// Refund (update dispute status)
	if err := repo.UpdateDisputeStatus(context.Background(), pool, d.ID, string(DisputeClosed)); err != nil {
		t.Fatalf("update dispute failed: %v", err)
	}

	// Create payout
	p := Payout{ID: "p-1", CreatorID: "c-1", Amount: 5.00, Status: PayoutPending}
	if err := repo.SavePayout(context.Background(), pool, p); err != nil {
		t.Fatalf("save payout failed: %v", err)
	}

	// Approve payout
	if err := repo.UpdatePayoutStatus(context.Background(), pool, p.ID, string(PayoutPaid)); err != nil {
		t.Fatalf("approve payout failed: %v", err)
	}

	// Verify state
	disputes, err := repo.ListDisputes(context.Background(), pool)
	if err != nil {
		t.Fatalf("list disputes failed: %v", err)
	}
	found := false
	for _, dd := range disputes {
		if dd.ID == d.ID && dd.Status == DisputeClosed {
			found = true
		}
	}
	if !found {
		t.Fatalf("dispute not closed as expected")
	}

	payouts, err := repo.ListPayouts(context.Background(), pool)
	if err != nil {
		t.Fatalf("list payouts failed: %v", err)
	}
	foundP := false
	for _, pp := range payouts {
		if pp.ID == p.ID && pp.Status == PayoutPaid {
			foundP = true
		}
	}
	if !foundP {
		t.Fatalf("payout not paid as expected")
	}
}

package recommendation

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/teamart/commerce-api/config"
	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

func TestIntegration_IngestAndRecommend(t *testing.T) {
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

	// Ensure table exists for test (migration may already have run in docker-compose)
	_, err = pool.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS recommendation_candidates (
  id TEXT PRIMARY KEY,
  type TEXT NOT NULL,
  category TEXT,
  signals JSONB,
  meta JSONB,
  updated_at TIMESTAMPTZ DEFAULT now()
);
`)
	if err != nil {
		t.Fatalf("failed to ensure table exists: %v", err)
	}

	// Clean up table before and after
	pool.Exec(context.Background(), "TRUNCATE recommendation_candidates")
	defer pool.Exec(context.Background(), "TRUNCATE recommendation_candidates")

	repo := NewPostgresRepository(log)

	cand := RecommendationCandidate{
		ID:       "cand-integ-1",
		Type:     ItemTypeProduct,
		Category: "sports",
		Signals: RankingSignals{
			WatchTime:        2.0,
			Purchases:        1.0,
			Reactions:        3.0,
			Follows:          0.5,
			CategoryAffinity: map[string]float64{"sports": 1.2},
		},
		Meta: map[string]interface{}{"title": "Integration Test Product"},
	}

	if err := repo.SaveCandidate(context.Background(), pool, cand); err != nil {
		t.Fatalf("save candidate failed: %v", err)
	}

	cands, err := repo.ListCandidates(context.Background(), pool)
	if err != nil {
		t.Fatalf("list candidates failed: %v", err)
	}

	found := false
	for _, c := range cands {
		if c.ID == cand.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("candidate not found after save")
	}

	weights := Weights{WatchTime: 1, Purchases: 3, Reactions: 0.5, Follows: 0.2, CategoryAffinity: 2}
	svc := NewInMemoryRecommendationService(cands, weights)
	recs, err := svc.RecommendForUser("user-1", 10)
	if err != nil {
		t.Fatalf("recommend failed: %v", err)
	}
	if len(recs) == 0 {
		t.Fatalf("no recommendations returned")
	}
	if recs[0].ID != cand.ID {
		t.Fatalf("expected top recommendation to be %s got %s", cand.ID, recs[0].ID)
	}
}

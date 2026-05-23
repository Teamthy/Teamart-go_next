package recommendation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

// Repository provides persistence for recommendation candidates.
type Repository interface {
	SaveCandidate(ctx context.Context, pool *database.Pool, c RecommendationCandidate) error
	ListCandidates(ctx context.Context, pool *database.Pool) ([]RecommendationCandidate, error)
}

// PostgresRepository implements Repository using Postgres JSONB columns.
type PostgresRepository struct {
	log *logger.Logger
}

func NewPostgresRepository(log *logger.Logger) *PostgresRepository {
	return &PostgresRepository{log: log}
}

// SaveCandidate inserts candidate data into a simple table. The table is expected to
// exist; this function focuses on demonstrating an ingestion path.
func (r *PostgresRepository) SaveCandidate(ctx context.Context, pool *database.Pool, c RecommendationCandidate) error {
	signalsB, err := json.Marshal(c.Signals)
	if err != nil {
		return fmt.Errorf("marshal signals: %w", err)
	}
	metaB, err := json.Marshal(c.Meta)
	if err != nil {
		return fmt.Errorf("marshal meta: %w", err)
	}

	sql := `INSERT INTO recommendation_candidates (id, type, category, signals, meta)
VALUES ($1, $2, $3, $4::jsonb, $5::jsonb)
ON CONFLICT (id) DO UPDATE SET type = EXCLUDED.type, category = EXCLUDED.category, signals = EXCLUDED.signals, meta = EXCLUDED.meta`

	_, err = pool.Exec(ctx, sql, c.ID, string(c.Type), c.Category, signalsB, metaB)
	if err != nil {
		return fmt.Errorf("save candidate: %w", err)
	}
	return nil
}

func (r *PostgresRepository) ListCandidates(ctx context.Context, pool *database.Pool) ([]RecommendationCandidate, error) {
	rows, err := pool.Query(ctx, `SELECT id, type, category, signals, meta FROM recommendation_candidates`)
	if err != nil {
		return nil, fmt.Errorf("query candidates: %w", err)
	}
	defer rows.Close()

	var out []RecommendationCandidate
	for rows.Next() {
		var id, t, category string
		var signalsB, metaB []byte
		if err := rows.Scan(&id, &t, &category, &signalsB, &metaB); err != nil {
			return nil, fmt.Errorf("scan candidate: %w", err)
		}
		var signals RankingSignals
		var meta map[string]interface{}
		_ = json.Unmarshal(signalsB, &signals)
		_ = json.Unmarshal(metaB, &meta)

		out = append(out, RecommendationCandidate{
			ID:       id,
			Type:     ItemType(t),
			Category: category,
			Signals:  signals,
			Meta:     meta,
		})
	}
	return out, nil
}

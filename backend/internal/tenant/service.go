package tenant

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

type Service struct {
	db  *database.Pool
	log *logger.Logger
}

func NewService(db *database.Pool, log *logger.Logger) *Service {
	return &Service{db: db, log: log}
}

type TenantSetting struct {
	ID        int64           `json:"id"`
	TenantID  int64           `json:"tenant_id"`
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func (s *Service) GetSetting(ctx context.Context, tenantID int64, key string) (*TenantSetting, error) {
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant ID is required")
	}
	if key == "" {
		return nil, fmt.Errorf("setting key is required")
	}

	row := s.db.QueryRow(ctx, `
		SELECT id, tenant_id, key, value, created_at, updated_at
		FROM tenant_settings
		WHERE tenant_id = $1 AND key = $2
	`, tenantID, key)

	var setting TenantSetting
	if err := row.Scan(&setting.ID, &setting.TenantID, &setting.Key, &setting.Value, &setting.CreatedAt, &setting.UpdatedAt); err != nil {
		return nil, err
	}

	return &setting, nil
}

func (s *Service) UpsertSetting(ctx context.Context, tenantID int64, key string, value json.RawMessage) (*TenantSetting, error) {
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant ID is required")
	}
	if key == "" {
		return nil, fmt.Errorf("setting key is required")
	}

	row := s.db.QueryRow(ctx, `
		INSERT INTO tenant_settings (tenant_id, key, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (tenant_id, key)
		DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP
		RETURNING id, tenant_id, key, value, created_at, updated_at
	`, tenantID, key, value)

	var setting TenantSetting
	if err := row.Scan(&setting.ID, &setting.TenantID, &setting.Key, &setting.Value, &setting.CreatedAt, &setting.UpdatedAt); err != nil {
		return nil, err
	}

	return &setting, nil
}

func (s *Service) ListSettings(ctx context.Context, tenantID int64) ([]TenantSetting, error) {
	if tenantID == 0 {
		return nil, fmt.Errorf("tenant ID is required")
	}

	rows, err := s.db.Query(ctx, `
		SELECT id, tenant_id, key, value, created_at, updated_at
		FROM tenant_settings
		WHERE tenant_id = $1
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var settings []TenantSetting
	for rows.Next() {
		var setting TenantSetting
		if err := rows.Scan(&setting.ID, &setting.TenantID, &setting.Key, &setting.Value, &setting.CreatedAt, &setting.UpdatedAt); err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

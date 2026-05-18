package migrations

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

// Runner implements MigrationRunner for PostgreSQL
// This provides explicit migration control with full observability
type Runner struct {
	pool       *database.Pool
	logger     *logger.Logger
	migrations []Migration
}

// NewRunner creates a new migration runner
func NewRunner(pool *database.Pool, log *logger.Logger, migrations []Migration) *Runner {
	// Sort migrations by version to ensure consistent order
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return &Runner{
		pool:       pool,
		logger:     log,
		migrations: migrations,
	}
}

// InitializeMigrationTable creates the migration tracking table
// This must be called before any migrations are run
func (r *Runner) InitializeMigrationTable(ctx context.Context) error {
	r.logger.Debugf("initializing migration table")

	// Create migrations table if it doesn't exist
	// This table tracks which migrations have been applied
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		id SERIAL PRIMARY KEY,
		version VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL,
		applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		executed_in BIGINT NOT NULL,
		success BOOLEAN NOT NULL DEFAULT true,
		error TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_schema_migrations_version ON schema_migrations(version);
	CREATE INDEX IF NOT EXISTS idx_schema_migrations_applied_at ON schema_migrations(applied_at);
	`

	_, err := r.pool.Exec(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize migration table: %w", err)
	}

	r.logger.Infof("migration table ready")
	return nil
}

// Migrate applies all pending migrations in order
// This is explicit schema evolution - we own the process
func (r *Runner) Migrate(ctx context.Context) error {
	r.logger.Infof("starting database migrations")

	// Initialize migration table
	if err := r.InitializeMigrationTable(ctx); err != nil {
		return err
	}

	// Get already applied migrations
	applied, err := r.getAppliedVersions(ctx)
	if err != nil {
		return err
	}

	appliedMap := make(map[string]bool)
	for _, v := range applied {
		appliedMap[v] = true
	}

	// Find pending migrations
	var pending []Migration
	for _, m := range r.migrations {
		if !appliedMap[m.Version] {
			pending = append(pending, m)
		}
	}

	if len(pending) == 0 {
		r.logger.Infof("no pending migrations")
		return nil
	}

	r.logger.Infof("found %d pending migration(s)", len(pending))

	// Apply each pending migration in a transaction
	for _, m := range pending {
		if err := r.applyMigration(ctx, m); err != nil {
			r.logger.Errorf("migration failed: %v", err)
			return fmt.Errorf("migration %s failed: %w", m.Version, err)
		}
	}

	r.logger.Infof("all migrations applied successfully")
	return nil
}

// applyMigration applies a single migration in a transaction
// Transaction ensures atomicity - either fully applied or fully rolled back
func (r *Runner) applyMigration(ctx context.Context, m Migration) error {
	r.logger.Infof("applying migration: %s (%s)", m.Version, m.Name)

	startTime := time.Now()

	// Start transaction
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Execute migration SQL
	_, err = tx.Exec(ctx, m.UpSQL)
	if err != nil {
		recordErr := r.recordMigration(ctx, m, time.Since(startTime), false, err.Error())
		if recordErr != nil {
			r.logger.Errorf("failed to record migration error: %v", recordErr)
		}
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record successful migration
	err = r.recordMigration(ctx, m, time.Since(startTime), true, "")
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	r.logger.Infof("migration applied: %s (took %v)", m.Version, time.Since(startTime))
	return nil
}

// recordMigration records migration execution in the database
func (r *Runner) recordMigration(ctx context.Context, m Migration, duration time.Duration, success bool, errMsg string) error {
	insertSQL := `
	INSERT INTO schema_migrations (version, name, executed_in, success, error)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, insertSQL, m.Version, m.Name, duration.Nanoseconds(), success, errMsg)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return nil
}

// Rollback reverts the last n migrations
// This provides explicit rollback safety
func (r *Runner) Rollback(ctx context.Context, steps int) error {
	r.logger.Warnf("rolling back %d migration(s)", steps)

	// Get applied migrations in reverse order
	applied, err := r.GetApplied(ctx)
	if err != nil {
		return err
	}

	// Sort by applied_at descending (most recent first)
	sort.Slice(applied, func(i, j int) bool {
		return applied[i].AppliedAt.After(applied[j].AppliedAt)
	})

	// Limit to the number of steps requested
	if steps > len(applied) {
		return fmt.Errorf("cannot rollback %d steps, only %d migrations have been applied", steps, len(applied))
	}

	toRollback := applied[:steps]

	// Find migration definitions for each applied migration
	migrationMap := make(map[string]Migration)
	for _, m := range r.migrations {
		migrationMap[m.Version] = m
	}

	// Rollback each migration
	for _, record := range toRollback {
		m, ok := migrationMap[record.Version]
		if !ok {
			return fmt.Errorf("migration definition not found for version %s", record.Version)
		}

		if err := r.rollbackMigration(ctx, m); err != nil {
			r.logger.Errorf("rollback failed: %v", err)
			return fmt.Errorf("rollback of %s failed: %w", m.Version, err)
		}
	}

	r.logger.Infof("rolled back %d migration(s)", steps)
	return nil
}

// rollbackMigration reverts a single migration
func (r *Runner) rollbackMigration(ctx context.Context, m Migration) error {
	r.logger.Infof("rolling back migration: %s (%s)", m.Version, m.Name)

	startTime := time.Now()

	// Start transaction
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Execute rollback SQL
	_, err = tx.Exec(ctx, m.DownSQL)
	if err != nil {
		return fmt.Errorf("failed to execute rollback SQL: %w", err)
	}

	// Remove migration record
	deleteSQL := `DELETE FROM schema_migrations WHERE version = $1`
	_, err = tx.Exec(ctx, deleteSQL, m.Version)
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	r.logger.Infof("migration rolled back: %s (took %v)", m.Version, time.Since(startTime))
	return nil
}

// Status returns current migration status
func (r *Runner) Status(ctx context.Context) (*MigrationStatus, error) {
	applied, err := r.GetApplied(ctx)
	if err != nil {
		return nil, err
	}

	pending, err := r.GetPending(ctx)
	if err != nil {
		return nil, err
	}

	var currentVersion string
	var lastAppliedAt time.Time
	if len(applied) > 0 {
		currentVersion = applied[len(applied)-1].Version
		lastAppliedAt = applied[len(applied)-1].AppliedAt
	}

	return &MigrationStatus{
		CurrentVersion:    currentVersion,
		PendingMigrations: len(pending),
		TotalMigrations:   len(r.migrations),
		AppliedMigrations: applied,
		LastAppliedAt:     lastAppliedAt,
	}, nil
}

// GetApplied returns all applied migrations
func (r *Runner) GetApplied(ctx context.Context) ([]MigrationRecord, error) {
	// Initialize table if needed
	_ = r.InitializeMigrationTable(ctx)

	querySQL := `
	SELECT id, version, name, applied_at, executed_in, success, error
	FROM schema_migrations
	WHERE success = true
	ORDER BY applied_at ASC
	`

	rows, err := r.pool.Query(ctx, querySQL)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	var records []MigrationRecord
	for rows.Next() {
		var record MigrationRecord
		if err := rows.Scan(
			&record.ID,
			&record.Version,
			&record.Name,
			&record.AppliedAt,
			&record.ExecutedIn,
			&record.Success,
			&record.Error,
		); err != nil {
			return nil, fmt.Errorf("failed to scan migration record: %w", err)
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

// GetPending returns all pending migrations
func (r *Runner) GetPending(ctx context.Context) ([]Migration, error) {
	applied, err := r.getAppliedVersions(ctx)
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]bool)
	for _, v := range applied {
		appliedMap[v] = true
	}

	var pending []Migration
	for _, m := range r.migrations {
		if !appliedMap[m.Version] {
			pending = append(pending, m)
		}
	}

	return pending, nil
}

// getAppliedVersions returns just the versions of applied migrations
func (r *Runner) getAppliedVersions(ctx context.Context) ([]string, error) {
	applied, err := r.GetApplied(ctx)
	if err != nil {
		// If table doesn't exist yet, return empty list
		return []string{}, nil
	}

	versions := make([]string, len(applied))
	for i, a := range applied {
		versions[i] = a.Version
	}

	return versions, nil
}

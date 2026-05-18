package migrations

import (
	"context"
	"time"
)

// Migration represents a database migration
// Migrations are versioned, timestamped, and reversible
type Migration struct {
	// Version is the unique identifier for this migration
	// Format: YYYYMMDDHHMMSS_description
	// Example: 20260101120000_create_users_table
	Version string

	// Name is the human-readable name
	Name string

	// Description explains what this migration does
	Description string

	// UpSQL is the SQL to apply the migration
	UpSQL string

	// DownSQL is the SQL to rollback the migration
	DownSQL string

	// CreatedAt is when this migration was created
	CreatedAt time.Time
}

// MigrationRecord represents a migration execution record in the database
// This tracks which migrations have been applied
type MigrationRecord struct {
	ID         int64
	Version    string
	Name       string
	AppliedAt  time.Time
	ExecutedIn int64 // nanoseconds
	Success    bool
	Error      string
}

// MigrationRunner executes migrations with explicit control
type MigrationRunner interface {
	// Migrate applies all pending migrations in order
	Migrate(ctx context.Context) error

	// Rollback reverts the last n migrations
	Rollback(ctx context.Context, steps int) error

	// Status returns the current migration status
	Status(ctx context.Context) (*MigrationStatus, error)

	// GetApplied returns all applied migrations
	GetApplied(ctx context.Context) ([]MigrationRecord, error)

	// GetPending returns all pending migrations
	GetPending(ctx context.Context) ([]Migration, error)
}

// MigrationStatus represents the current state of migrations
type MigrationStatus struct {
	// CurrentVersion is the latest applied migration
	CurrentVersion string

	// PendingMigrations is the count of unapplied migrations
	PendingMigrations int

	// TotalMigrations is the total count of migrations
	TotalMigrations int

	// AppliedMigrations is the list of applied migrations
	AppliedMigrations []MigrationRecord

	// LastError is the last error encountered during migration
	LastError string

	// LastAppliedAt is when the last migration was applied
	LastAppliedAt time.Time
}

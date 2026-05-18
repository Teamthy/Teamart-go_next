package queries

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DBTX is the interface satisfied by pgx connections and transactions
// This is the core interface for SQLC query execution
type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

// Queries contains all generated query methods
// This is populated by SQLC with compiled, type-safe query methods
type Queries struct {
	db DBTX
}

// New creates a new Queries instance
// This wraps the pgx connection with SQLC-generated query methods
func New(db DBTX) *Queries {
	return &Queries{db: db}
}

// WithTx creates a new Queries with a transaction
// Allows multiple queries in a single transaction
func (q *Queries) WithTx(tx pgx.Tx) *Queries {
	return &Queries{db: tx}
}

// Db returns the underlying database connection
func (q *Queries) Db() DBTX {
	return q.db
}

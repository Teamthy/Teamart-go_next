package database

import (
	"github.com/jackc/pgx/v5"
)

// QueryScanner represents a single row query result
// This allows explicit control over row scanning
type QueryScanner interface {
	Scan(dest ...interface{}) error
}

// Rows represents a result set from a query
// This is an alias to pgx.Rows for explicit type declaration
type Rows = pgx.Rows

// Tx represents a database transaction
// Explicit transaction ownership - caller must commit or rollback
type Tx = pgx.Tx

// TxOptions represents transaction options
type TxOptions = pgx.TxOptions

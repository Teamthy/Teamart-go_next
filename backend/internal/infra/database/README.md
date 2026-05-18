# Database Foundation - PostgreSQL Connection Layer

## Overview

This document explains the **PostgreSQL Connection Layer** - the first step in building a SQL-first backend architecture using SQLC.

## Architecture Principles

### Why SQL-First?

Senior backend engineering teams prefer explicit SQL over ORMs because:

- **Explicit queries**: You see exactly what SQL is executed
- **Explicit transactions**: You control when queries are committed/rolled back
- **Explicit indexes**: You understand query performance implications
- **Debugging clarity**: Errors point to specific SQL, not ORM magic
- **Performance ownership**: You write queries that match your data patterns

This backend uses **SQLC** for compile-safe SQL generation, ensuring:
- Type-safe queries at compile time
- Zero runtime magic
- Clear performance characteristics

## Connection Pooling Architecture

### The Pool

The `database.Pool` manages a PostgreSQL connection pool with explicit lifecycle management:

```
Application Startup
    ↓
database.NewPool() - Create connection pool
    ↓
pgxpool.Config - Configure pool settings:
    - MaxConns: 25 (default)
    - MinConns: 5 (default)
    - MaxConnLifetime: 15 minutes
    - MaxConnIdleTime: 5 minutes
    ↓
pool.Ping() - Validate connectivity
    ↓
App.RegisterCleanup() - Register pool.Close()
    ↓
Application Running
    ↓
Shutdown Signal
    ↓
pool.Close() - Graceful connection cleanup
```

### Configuration

Set these in `.env`:

```env
# Connection string (required)
DATABASE_URL=postgres://user:password@localhost:5432/database?sslmode=disable

# Pool settings (optional, defaults shown)
DATABASE_MAX_CONNECTIONS=25    # Maximum connections in pool
DATABASE_MIN_CONNECTIONS=5     # Minimum idle connections maintained
DATABASE_CONN_MAX_LIFETIME=15m # Maximum connection age
DATABASE_CONN_MAX_IDLE_TIME=5m # Maximum time before idle connection closed
```

### Pool Statistics

Monitor pool health in real-time:

```go
stats := db.Stats()
// stats.AcquiredConns - Connections currently in use
// stats.IdleConns - Connections available
// stats.TotalConns - Total connections
// stats.ConstructingConns - Connections being created
```

## Health Validation

### Explicit Health Checks

The pool provides explicit, observable health validation:

```go
health, err := db.Health(ctx)
// health.Status: "healthy", "degraded", or "unhealthy"
// health.ResponseTime: Ping latency in milliseconds
// health.PoolStats: Current connection pool status
```

### HTTP Health Endpoints

The application provides health check endpoints:

- `GET /health` - Basic application health
- `GET /ready` - Readiness check (includes DB validation)
- `GET /api/v1/diagnostics/db` - Database diagnostics (includes pool stats)

Example readiness check:

```bash
curl http://localhost:8000/ready
# {"status":"ready","database":"healthy"}
```

Example diagnostics:

```bash
curl http://localhost:8000/api/v1/diagnostics/db
# {
#   "status": "healthy",
#   "response_time_ms": 2,
#   "pool_stats": {
#     "acquired_conns": 0,
#     "idle_conns": 5,
#     "total_conns": 5,
#     "constructing_conns": 0
#   }
# }
```

## Infrastructure Ownership

### Explicit Lifecycle Management

All database resources are explicitly managed:

1. **Initialization**: `database.NewPool()` creates the pool with validated configuration
2. **Query Execution**: Explicit `db.Query()`, `db.QueryRow()`, `db.Exec()` methods
3. **Transactions**: Explicit `db.BeginTx()` returns a transaction handle
4. **Cleanup**: Pool cleanup is registered in application shutdown chain

### Error Handling

All database errors are explicit and traceable:

```go
// Query error example
rows, err := db.Query(ctx, "SELECT * FROM users WHERE id = $1", userID)
if err != nil {
    // Error includes context: "query error: <underlying error>"
    return fmt.Errorf("failed to fetch user: %w", err)
}
defer rows.Close()
```

## Integration with Application Lifecycle

The database pool is integrated into the application's graceful shutdown:

```go
// In main.go
db, err := database.NewPool(context.Background(), &cfg.Database, log)
if err != nil {
    log.Errorf("failed to initialize database: %v", err)
    os.Exit(1)
}

// Register cleanup - this runs when application shuts down
application.RegisterCleanup(func(shutdownCtx context.Context) error {
    log.Info("closing database connection pool...")
    db.Close()
    log.Info("database connection pool closed")
    return nil
})
```

When the application receives a shutdown signal (SIGTERM/SIGINT):

1. HTTP server stops accepting new requests
2. Cleanup handlers run in reverse order
3. Database pool closes gracefully
4. All connections are released

## Usage Patterns

### Simple Query

```go
var name string
err := db.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", userID).Scan(&name)
if err != nil {
    return "", fmt.Errorf("failed to fetch user: %w", err)
}
```

### Rows Query

```go
rows, err := db.Query(ctx, "SELECT id, name, email FROM users WHERE active = true")
if err != nil {
    return nil, fmt.Errorf("failed to fetch users: %w", err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var name, email string
    if err := rows.Scan(&id, &name, &email); err != nil {
        return nil, fmt.Errorf("failed to scan row: %w", err)
    }
    // Process user
}
```

### Transaction

```go
tx, err := db.BeginTx(ctx, pgx.TxOptions{})
if err != nil {
    return fmt.Errorf("failed to begin transaction: %w", err)
}
defer tx.Rollback(ctx)

// Execute queries on tx
_, err = tx.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", name)
if err != nil {
    return fmt.Errorf("failed to insert user: %w", err)
}

// Explicit commit
if err := tx.Commit(ctx); err != nil {
    return fmt.Errorf("failed to commit transaction: %w", err)
}
```

### Exec (INSERT/UPDATE/DELETE)

```go
rowsAffected, err := db.Exec(
    ctx,
    "UPDATE users SET updated_at = NOW() WHERE id = $1",
    userID,
)
if err != nil {
    return fmt.Errorf("failed to update user: %w", err)
}
log.Infof("updated %d users", rowsAffected)
```

## Monitoring

### Connection Pool Metrics

Monitor these metrics in production:

- **Acquired Connections**: How many connections are currently in use
- **Idle Connections**: How many connections are available
- **Total Connections**: Should stabilize near MaxConns
- **Connection Errors**: Spike indicates database or network issues

### Performance Characteristics

The pool automatically manages:

- Connection creation: Created on-demand up to MaxConns
- Connection reuse: Idle connections are reused (lowest latency)
- Connection closure: Idle connections > MaxConnIdleTime are closed
- Connection timeout: Connections > MaxConnLifetime are recycled

## Docker & Local Development

### Running PostgreSQL Locally

```bash
# Using Docker Compose
docker-compose up -d postgres

# Verify connection
psql -U postgres -d teamart -h localhost
```

### Connection String Examples

**Local development**:
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/teamart?sslmode=disable
```

**Docker Compose**:
```
DATABASE_URL=postgres://postgres:postgres@db:5432/teamart?sslmode=disable
```

**Production (CloudSQL, RDS, etc)**:
```
DATABASE_URL=postgres://user:password@host:5432/teamart?sslmode=require
```

## Next Steps

The database foundation is now complete. Next phase:

1. ✅ **PostgreSQL Connection Layer** - Connection pooling ✓
2. **Migration System** - Schema versioning and rollback
3. **SQLC Integration** - Type-safe SQL code generation

The next step will add schema versioning capabilities, allowing you to evolve your database safely.

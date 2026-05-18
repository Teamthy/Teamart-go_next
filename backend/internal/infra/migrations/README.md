# Migration System

## Overview

The migration system provides explicit, versioned schema evolution with full rollback safety. This is the core of deployment consistency and safe database evolution.

## Architecture

### Concepts

**Migration**: A versioned, timestamped change to the database schema.

```
20260101000000_init_schema
├── Version: 20260101000000 (YYYYMMDDHHMMSS format)
├── Name: init_schema (human-readable)
├── UpSQL: Create tables, indexes, etc.
└── DownSQL: Drop tables, reverse changes
```

**Migration Record**: Tracks execution in `schema_migrations` table.

```
schema_migrations
├── id: Auto-incrementing
├── version: Unique migration identifier
├── name: Human-readable name
├── applied_at: Timestamp when applied
├── executed_in: Duration in nanoseconds
├── success: Whether migration succeeded
└── error: Error message if failed
```

**MigrationRunner**: Executes migrations with explicit control.

```
Runner
├── Migrate() - Apply all pending migrations
├── Rollback() - Revert last n migrations
├── Status() - Get current migration status
├── GetApplied() - List applied migrations
└── GetPending() - List unapplied migrations
```

## Usage

### Automatic Migrations on Startup

Migrations run automatically when the application starts:

```go
// In cmd/api/main.go
runner := migrations.NewRunner(db, log, migrations.Migrations)
if err := runner.Migrate(context.Background()); err != nil {
    log.Errorf("migration failed: %v", err)
    os.Exit(1)
}
```

This ensures:
- Schema is always up-to-date on startup
- Failed migrations prevent application from running
- Migration status is logged for observability

### Check Migration Status

```bash
curl http://localhost:8000/api/v1/diagnostics/migrations

# Response:
{
  "current_version": "20260101000000",
  "applied": 1,
  "pending": 0,
  "total": 1
}
```

### Create New Migrations

Add to `internal/infra/migrations/migrations.go`:

```go
var migrationCreateUsersTable = Migration{
    Version:     "20260102120000",
    Name:        "create_users_table",
    Description: "Create users table with authentication",

    UpSQL: `
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `,

    DownSQL: `
    DROP TABLE IF EXISTS users;
    `,
}

// Then add to Migrations slice:
var Migrations = []Migration{
    migrationInitSchema,
    migrationCreateUsersTable,  // New migration
}
```

## Version Format

Versions use timestamp format: `YYYYMMDDHHMMSS`

```
20260101000000 = January 1, 2026, 00:00:00 UTC
20260102150930 = January 2, 2026, 15:09:30 UTC
```

This ensures:
- Versions are unique
- Versions sort chronologically
- Easy to identify when migration was created

## Migration Lifecycle

### 1. Pending State
```
Application starts
    ↓
Migration runner looks for pending migrations
    ↓
Compares against schema_migrations table
    ↓
Finds new migrations not in table
    ↓
Migrations are "pending"
```

### 2. Application State
```
Pending Migration
    ↓
BEGIN TRANSACTION
    ↓
Execute UP SQL
    ↓
Record in schema_migrations
    ↓
COMMIT
    ↓
Applied Migration
```

### 3. Tracking State
```
Applied = recorded in schema_migrations with success=true
Pending = in migrations.Migrations but not in schema_migrations
```

## Rollback Safety

### Explicit Rollback

```go
// Rollback last 2 migrations
err := runner.Rollback(ctx, 2)
if err != nil {
    log.Errorf("rollback failed: %v", err)
}
```

### Safety Mechanisms

1. **Transaction Safety**: Each migration in its own transaction
   - Success: Fully applied
   - Failure: Fully rolled back

2. **Record Tracking**: Migration record deleted on rollback
   - Ensures rollback can be re-applied
   - No duplicate key errors

3. **Explicit DownSQL**: Every migration has rollback SQL
   - NOT reverse-engineered
   - YOU write the safe rollback
   - Full control over what's reversed

### Rollback Limitations

- Can only rollback applied migrations
- Cannot rollback in production (use migration to fix instead)
- Requires DownSQL for each migration

## Production Best Practices

### 1. Write Both Up and Down
```go
UpSQL: `CREATE TABLE users (...);`,      // What you add
DownSQL: `DROP TABLE users;`,            // How to remove it
```

### 2. Test Migrations Locally First
```bash
# Run application locally
go run ./cmd/api

# Verify migrations applied
curl http://localhost:8000/api/v1/diagnostics/migrations

# Check schema
psql -U postgres -d teamart -c "\dt"
```

### 3. Never Modify Applied Migrations
```
❌ WRONG: Change existing migration SQL
✅ RIGHT: Create new migration to fix issues
```

### 4. Deploy Migrations First
```
1. Deploy new code with migrations included
2. Migrations run automatically on startup
3. Application starts after migrations complete
4. All servers have consistent schema
```

### 5. Large Migrations
For large tables, consider:
```go
UpSQL: `
    ALTER TABLE users ADD COLUMN new_field VARCHAR(255);
    UPDATE users SET new_field = 'default' WHERE new_field IS NULL;
    ALTER TABLE users ALTER COLUMN new_field SET NOT NULL;
`,
```

## Monitoring

### Health Endpoint
```bash
# Check if migrations are up-to-date
curl http://localhost:8000/ready

# Response if migrations pending:
{"status":"not_ready","database":"unavailable"}
```

### Migration Status Endpoint
```bash
curl http://localhost:8000/api/v1/diagnostics/migrations

# Returns:
{
  "current_version": "20260101000000",
  "applied": 5,
  "pending": 0,
  "total": 5
}
```

### Database Logs
```
[INFO]  running database migrations...
[INFO]  applying migration: 20260101000000 (init_schema)
[INFO]  migration applied: 20260101000000 (took 123ms)
[INFO]  migrations status: 1 applied, 0 pending
```

## Example: Create New Migration

### 1. Add Migration Definition

File: `internal/infra/migrations/migrations.go`

```go
var migrationAddInventoryTable = Migration{
    Version:     "20260103140000",  // New timestamp
    Name:        "add_inventory_table",
    Description: "Create inventory tracking table",

    UpSQL: `
    CREATE TABLE inventory (
        id SERIAL PRIMARY KEY,
        product_id INTEGER NOT NULL REFERENCES products(id),
        quantity INTEGER NOT NULL DEFAULT 0,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX idx_inventory_product_id ON inventory(product_id);
    `,

    DownSQL: `
    DROP TABLE IF EXISTS inventory CASCADE;
    `,
}

// Add to Migrations slice:
var Migrations = []Migration{
    migrationInitSchema,
    migrationAddInventoryTable,  // NEW
}
```

### 2. Test Locally

```bash
# Start app
go run ./cmd/api

# Check migrations were applied
curl http://localhost:8000/api/v1/diagnostics/migrations

# Verify table exists
psql -U postgres -d teamart -c "\d inventory"
```

### 3. Deploy

```bash
# Build new executable with migration included
go build ./cmd/api

# Deploy to production
# Migration runs automatically on startup
```

## Files

- `types.go` - Migration type definitions
- `runner.go` - Migration execution engine
- `migrations.go` - Migration collection
- `README.md` - This documentation

## Integration Points

### Database Layer
```go
runner := migrations.NewRunner(db, log, migrations.Migrations)
```

### Application Startup
```go
// Runs in cmd/api/main.go
if err := runner.Migrate(context.Background()); err != nil {
    os.Exit(1)
}
```

### HTTP Endpoints
- `GET /api/v1/diagnostics/migrations` - Status endpoint

## Comparison: Before vs After

### Before (No Migration System)
```
❌ Manual schema changes
❌ Track changes in comments or README
❌ No rollback capability
❌ Inconsistent schemas across environments
❌ No versioning
```

### After (Migration System)
```
✅ Versioned, timestamped changes
✅ Explicit up/down SQL
✅ Automatic tracking in schema_migrations
✅ Consistent deployment across all servers
✅ Full rollback capability
✅ Observable status
```

## Next Steps

The migration system is now complete. Next phase:

1. ✅ **PostgreSQL Connection Layer** - Connection pooling ✓
2. ✅ **Migration System** - Schema versioning ✓
3. **SQLC Integration** - Type-safe SQL code generation

SQLC will allow compile-safe queries with automatic typing.

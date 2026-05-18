# SQLC Integration - Type-Safe SQL Code Generation

## Overview

SQLC is a tool that generates type-safe Go code from explicit SQL queries. This is the final piece of the SQL-first architecture.

Instead of:
- ❌ ORMs with hidden queries
- ❌ String concatenation for SQL
- ❌ Runtime type errors

We have:
- ✅ Explicit, version-controlled SQL
- ✅ Compile-time type validation
- ✅ No runtime magic
- ✅ Full query ownership

## Architecture

```
Your SQL Queries
    ↓
sqlc reads queries + schema
    ↓
Generates type-safe Go code
    ↓
You use generated code in your handlers
    ↓
100% type-safe, compile-checked queries
```

## Configuration

### sqlc.yaml

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./internal/infra/queries"      # Where SQL files live
    schema: "./internal/infra/migrations"    # Where schema is defined
    gen:
      go:
        out: "./internal/infra/queries/generated"  # Where code is generated
        package: "queries"
```

This tells SQLC:
1. Read SQL from `internal/infra/queries/*.sql`
2. Validate against schema in `internal/infra/migrations/`
3. Generate Go code in `internal/infra/queries/generated/`

## Query Files

### Directory Structure

```
internal/infra/queries/
├── users.sql           # User-related queries
├── products.sql        # Product queries
├── orders.sql          # Order queries
├── order_items.sql     # Order item queries
└── generated/          # SQLC-generated code
    ├── models.go       # Data types (User, Product, Order, etc.)
    ├── querier.go      # Querier interface
    ├── users.sql.go    # User query implementations
    ├── products.sql.go # Product query implementations
    ├── orders.sql.go   # Order query implementations
    └── order_items.sql.go  # Order item implementations
```

## Query Syntax

### Query Comments

Every query needs a comment with:
- `-- name: FunctionName :operation`
- `operation` is one of: `:one`, `:many`, `:exec`, `:execrows`

```sql
-- name: GetUserByID :one
SELECT id, email, name, password_hash, created_at, updated_at
FROM users
WHERE id = $1;
```

This generates:
```go
func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error)
```

### Operations

| Operation | Returns | Use Case |
|-----------|---------|----------|
| `:one` | Single row | `SELECT` returning one record |
| `:many` | Multiple rows | `SELECT` returning many records |
| `:exec` | No rows | `INSERT`, `UPDATE`, `DELETE` |
| `:execrows` | Row count | `INSERT`, `UPDATE`, `DELETE` with count |

### Parameter Types

SQL parameters use `$1, $2, etc.` (PostgreSQL style):

```sql
-- name: CreateUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING id, email, name, password_hash, created_at, updated_at;
```

SQLC infers parameter types from the schema:
- `$1` → `email` column (VARCHAR) → `string` parameter
- `$2` → `name` column (VARCHAR) → `string` parameter
- `$3` → `password_hash` column (VARCHAR) → `string` parameter

Generated function signature:
```go
func (q *Queries) CreateUser(ctx context.Context, email string, name string, passwordHash string) (User, error)
```

## Provided Queries

### Users

**CreateUser** - Insert new user
```go
user, err := queries.CreateUser(ctx, "user@example.com", "John Doe", "hashed_password")
```

**GetUserByID** - Fetch by ID
```go
user, err := queries.GetUserByID(ctx, 1)
```

**GetUserByEmail** - Fetch by email
```go
user, err := queries.GetUserByEmail(ctx, "user@example.com")
```

**ListUsers** - Paginated list
```go
users, err := queries.ListUsers(ctx, limit, offset)
```

**UpdateUser** - Update user
```go
user, err := queries.UpdateUser(ctx, userID, newName)
```

**DeleteUser** - Delete user
```go
err := queries.DeleteUser(ctx, userID)
```

**CountUsers** - Total count
```go
count, err := queries.CountUsers(ctx)
```

### Products

**CreateProduct** - Insert product
**GetProductByID** - Fetch by ID
**GetProductBySKU** - Fetch by SKU
**ListProducts** - Paginated list
**SearchProducts** - Search by name/description
**UpdateProduct** - Update product
**DeleteProduct** - Delete product
**CountProducts** - Total count

### Orders

**CreateOrder** - Insert order
**GetOrderByID** - Fetch by ID
**ListOrdersByUserID** - User's orders
**ListOrdersByStatus** - Orders by status
**ListAllOrders** - All orders
**UpdateOrderStatus** - Change status
**UpdateOrderTotalAmount** - Change amount
**DeleteOrder** - Delete order
**CountOrdersByUserID** - User's order count
**CountOrdersByStatus** - Orders by status count

### Order Items

**CreateOrderItem** - Insert item
**GetOrderItemByID** - Fetch by ID
**ListOrderItemsByOrderID** - Items in order
**GetOrderItemsByProductID** - Product purchases
**UpdateOrderItemQuantity** - Change quantity
**UpdateOrderItemPrice** - Change price
**DeleteOrderItem** - Delete item
**DeleteOrderItems** - Delete all order items
**CountOrderItems** - Item count
**GetOrderItemsWithProduct** - Items with product details

## Generated Code

### Models

SQLC generates type-safe structs:

```go
type User struct {
    ID           int32     `db:"id" json:"id"`
    Email        string    `db:"email" json:"email"`
    Name         string    `db:"name" json:"name"`
    PasswordHash string    `db:"password_hash" json:"password_hash"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
```

- `db` tags show column names
- `json` tags for serialization
- Type-safe fields (not `interface{}`)

### Query Methods

Each query becomes a type-safe method:

```go
// Generated from: -- name: GetUserByID :one
func (q *Queries) GetUserByID(ctx context.Context, id int32) (User, error) {
    row := q.db.QueryRow(ctx,
        `SELECT id, email, name, password_hash, created_at, updated_at
         FROM users WHERE id = $1`, id)
    
    var i User
    err := row.Scan(&i.ID, &i.Email, &i.Name, &i.PasswordHash, &i.CreatedAt, &i.UpdatedAt)
    return i, err
}
```

### Benefits

1. **Type Safety**: Compiler checks field names and types
2. **Zero Runtime Magic**: Plain SQL with Go code
3. **Query Ownership**: You write and control the SQL
4. **Refactoring Safety**: Rename columns → compiler error guides fixes
5. **Performance**: Hand-written SQL, optimized by you

## Usage in Application

### Initialize Queries

```go
import "github.com/teamart/commerce-api/internal/infra/queries"

// Create queries instance
db := &database.Pool{...}
q := queries.New(db.pool)
```

### Basic Query

```go
user, err := q.GetUserByID(ctx, 1)
if err != nil {
    log.Errorf("failed to fetch user: %v", err)
    return
}
log.Infof("user: %s (%s)", user.Name, user.Email)
```

### List with Pagination

```go
limit := int32(20)
offset := int32(0)

users, err := q.ListUsers(ctx, limit, offset)
if err != nil {
    return err
}

for _, user := range users {
    log.Infof("user: %s", user.Email)
}
```

### Insert

```go
hashedPwd, _ := hashPassword(password)

user, err := q.CreateUser(ctx,
    "newuser@example.com",
    "New User",
    hashedPwd)
if err != nil {
    return err
}

log.Infof("created user: %d", user.ID)
```

### Transaction

```go
// Begin transaction
tx, err := db.BeginTx(ctx, pgx.TxOptions{})
if err != nil {
    return err
}
defer tx.Rollback(ctx)

// Create new querier for transaction
qtx := q.WithTx(tx)

// Execute multiple queries in transaction
order, err := qtx.CreateOrder(ctx, userID, totalAmount, "pending")
if err != nil {
    return err
}

item1, err := qtx.CreateOrderItem(ctx, order.ID, product1ID, qty1, price1)
if err != nil {
    return err
}

item2, err := qtx.CreateOrderItem(ctx, order.ID, product2ID, qty2, price2)
if err != nil {
    return err
}

// Commit if all succeeded
if err = tx.Commit(ctx); err != nil {
    return err
}
```

## SQL Best Practices

### 1. Explicit Column Selection

```sql
❌ WRONG:
SELECT * FROM users;

✅ RIGHT:
SELECT id, email, name, password_hash, created_at, updated_at FROM users;
```

The second is more explicit and changes are visible.

### 2. Parameterize Everything

```sql
❌ WRONG:
SELECT * FROM users WHERE id = 1;  -- hardcoded value

✅ RIGHT:
SELECT * FROM users WHERE id = $1;  -- parameterized
```

### 3. Use Meaningful Column Names

```sql
❌ WRONG:
SELECT u.id, u.name, p.id, p.name FROM users u JOIN products p ON ...

✅ RIGHT:
SELECT u.id, u.name, p.id AS product_id, p.name AS product_name FROM users u JOIN products p ON ...
```

SQLC uses column names to name struct fields.

### 4. Write Explicit JOINs

```sql
❌ WRONG:
SELECT * FROM users, orders WHERE users.id = orders.user_id;

✅ RIGHT:
SELECT u.id, u.email, o.id, o.status FROM users u
JOIN orders o ON u.id = o.user_id;
```

## Workflow

### 1. Write SQL

```sql
-- File: internal/infra/queries/users.sql
-- name: CreateUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING *;
```

### 2. Run SQLC

```bash
sqlc generate
```

This generates code in `internal/infra/queries/generated/`

### 3. Use Generated Code

```go
user, err := q.CreateUser(ctx, email, name, hash)
```

### 4. Modify Schema?

```bash
# Update migration
# Run migration
# Re-run sqlc generate
sqlc generate
```

SQLC validates SQL against new schema.

## Comparison: Before vs After

### Before (No SQLC)

```go
// Manual, error-prone, type-unsafe
row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
var id, name, email interface{}
err := row.Scan(&id, &name, &email)  // Wrong types!

// Later... runtime panic
name := email.(string)  // Type assertion fails
```

### After (SQLC)

```go
// Type-safe, compile-checked
user, err := q.GetUserByID(ctx, id)
if err != nil {
    return err
}

// Type already correct
name := user.Name  // string, guaranteed
```

## Testing Queries

Since queries are explicit SQL, you can:

1. **Test the SQL directly** in psql
2. **Test in Go** with real database
3. **Mock queries** by implementing the Querier interface

```go
// Mock for testing
type MockQueries struct{}

func (m *MockQueries) GetUserByID(ctx context.Context, id int32) (queries.User, error) {
    return queries.User{
        ID: id,
        Email: "test@example.com",
        Name: "Test User",
    }, nil
}
```

## Performance

SQLC generates identical queries to hand-written SQL:

```go
// SQLC generates this:
row := q.db.QueryRow(ctx,
    `SELECT id, email, name FROM users WHERE id = $1`,
    id)

// Same as if you wrote this manually:
row := db.QueryRow(`SELECT id, email, name FROM users WHERE id = $1`, id)
```

**Zero performance overhead** - SQLC is a code generator, not a runtime framework.

## Files in This Implementation

- `sqlc.yaml` - Configuration
- `internal/infra/queries/*.sql` - SQL query definitions
- `internal/infra/queries/generated/*.go` - Generated code
- `internal/infra/queries/querier.go` - Query wrapper

## Next Steps

With SQLC integration complete:

✅ Explicit connection pooling (STEP 1)
✅ Versioned migrations (STEP 2)
✅ Type-safe SQL generation (STEP 3)

The **SQL-first backend foundation** is now complete.

All infrastructure is explicit, observable, and production-ready:
- No hidden queries
- No runtime magic
- Full debugging visibility
- Type safety at compile time
- Performance ownership

This is senior backend engineering.

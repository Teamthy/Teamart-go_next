# Query Service Layer Architecture

This document describes the service layer architecture and how it integrates with SQLC-generated queries and HTTP handlers.

## Overview

The service layer provides a bridge between HTTP handlers and the database layer. It encapsulates business logic, validation, and error handling while using SQLC-generated queries for database access.

## Architecture Pattern

```
┌─────────────────────────────────────────────────────────────────┐
│                       HTTP Layer                                │
│  (Request parsing, response encoding, status codes)             │
│  - user_handler.go                                              │
│  - product_handler.go                                           │
│  - order_handler.go                                             │
└──────────────────────┬──────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                     Service Layer                               │
│  (Business logic, validation, error handling, logging)          │
│  - internal/users/service.go                                    │
│  - internal/products/service.go                                 │
│  - internal/orders/service.go                                   │
└──────────────────────┬──────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                  SQLC Query Layer                               │
│  (Type-safe database queries, transaction support)              │
│  - internal/infra/queries/                                      │
│    - queries.go         (Query interface)                       │
│    - users_gen.sql.go   (Generated user queries)                │
│    - products_gen.sql.go (Generated product queries)            │
│    - orders_gen.sql.go   (Generated order queries)              │
│    - models.go          (Generated models)                      │
└──────────────────────┬──────────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Database Layer                                │
│  (PostgreSQL connection pool, health checks)                    │
│  - internal/infra/database/database.go                          │
│  - pgx connection pool                                          │
└─────────────────────────────────────────────────────────────────┘
```

## Detailed Architecture

### HTTP Handler Layer

Responsibilities:
- Parse HTTP requests
- Validate HTTP-level requirements (path parameters, query strings)
- Call service methods
- Encode responses as JSON
- Set HTTP status codes

Example:
```go
func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var req CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // Call service
    output, err := h.service.CreateUser(r.Context(), &users.CreateUserInput{...})
    
    // Encode response
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(output)
}
```

### Service Layer

Responsibilities:
- Define input/output contracts via Input/Output types
- Validate business logic
- Transform data
- Log operations
- Call SQLC queries
- Handle errors consistently

Example:
```go
type Service struct {
    queries *queries.Queries
    logger  *logger.Logger
}

type CreateUserInput struct {
    Email        string
    Name         string
    PasswordHash string
}

func (s *Service) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
    // Validate
    if input.Email == "" {
        return nil, fmt.Errorf("email is required")
    }
    
    // Log
    s.logger.Debugf("creating user with email: %s", input.Email)
    
    // Call SQLC query
    user, err := s.queries.CreateUser(ctx, queries.CreateUserParams{
        Email:        input.Email,
        Name:         input.Name,
        PasswordHash: input.PasswordHash,
    })
    if err != nil {
        s.logger.Errorf("failed to create user: %v", err)
        return nil, err
    }
    
    // Return output
    return &CreateUserOutput{ID: user.ID, ...}, nil
}
```

### SQLC Query Layer

Responsibilities:
- Type-safe SQL execution
- Parameter binding
- Result mapping to structs
- Transaction support

Example generated code:
```go
// User model (generated)
type User struct {
    ID           int64     `db:"id" json:"id"`
    Email        string    `db:"email" json:"email"`
    Name         string    `db:"name" json:"name"`
    // ...
}

// Query parameters (generated)
type CreateUserParams struct {
    Email        string
    Name         string
    PasswordHash string
}

// Query method (generated)
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
    // Implementation: executes SQL with proper parameter binding
}
```

### Database Layer

Responsibilities:
- Manage connection pool
- Health checks
- Connection lifecycle

## Data Flow Examples

### Create User Request

```
1. HTTP Handler receives POST /users request
2. Handler parses JSON body into CreateUserRequest
3. Handler calls service.CreateUser(ctx, input)
4. Service validates input (email not empty, etc.)
5. Service logs "creating user with email: ..."
6. Service calls queries.CreateUser(ctx, params)
7. SQLC executes: INSERT INTO users ... RETURNING ...
8. Database returns inserted user row
9. SQLC maps row to User struct
10. Service wraps User in CreateUserOutput
11. Service logs "user created with ID: ..."
12. Handler receives CreateUserOutput
13. Handler encodes to JSON response
14. Handler sends HTTP 201 Created response
15. Client receives JSON response
```

### Get User by ID Request

```
1. HTTP Handler receives GET /users/1 request
2. Handler extracts ID from path parameter
3. Handler calls service.GetUserByID(ctx, input)
4. Service validates ID (not zero)
5. Service logs "fetching user with ID: 1"
6. Service calls queries.GetUserByID(ctx, id)
7. SQLC executes: SELECT ... FROM users WHERE id = $1
8. Database returns user row or no rows
9. SQLC maps row to User struct
10. Service wraps User in GetUserByIDOutput
11. Service logs "user found: ID=1, Email=..."
12. Handler receives GetUserByIDOutput
13. Handler encodes to JSON response
14. Handler sends HTTP 200 OK response
15. Client receives JSON response
```

### List Users with Pagination

```
1. HTTP Handler receives GET /users?limit=10&offset=0
2. Handler parses query parameters
3. Handler calls service.ListUsers(ctx, input)
4. Service validates limit (not zero, max 100)
5. Service logs "listing users with limit: 10, offset: 0"
6. Service calls queries.ListUsers(ctx, limit, offset)
7. SQLC executes: SELECT ... FROM users ... LIMIT $1 OFFSET $2
8. Database returns multiple user rows
9. SQLC maps rows to []User slice
10. Service wraps []User in ListUsersOutput
11. Service logs "fetched 10 users"
12. Handler receives ListUsersOutput
13. Handler encodes to JSON array response
14. Handler sends HTTP 200 OK response
15. Client receives JSON array response
```

## Service Layer Benefits

### 1. Clear Separation of Concerns
- HTTP layer: Protocol handling (JSON, HTTP status codes)
- Service layer: Business logic and validation
- Query layer: Database access
- Each layer has a single responsibility

### 2. Reusability
Services can be called from different sources:
```go
// From HTTP handler
output, err := service.CreateUser(ctx, input)

// From background job
output, err := service.CreateUser(ctx, input)

// From CLI command
output, err := service.CreateUser(ctx, input)
```

### 3. Testability
Easy to test each layer independently:
```go
// Test service in isolation
func TestCreateUser(t *testing.T) {
    mockQueries := &MockQueries{}
    service := users.NewService(mockQueries, logger)
    output, err := service.CreateUser(ctx, input)
    // Assert output and error
}

// Test handler in isolation
func TestHandleCreateUser(t *testing.T) {
    mockService := &MockService{}
    handler := NewUserHandler(mockService, logger)
    // Call handler with mock request/response
    // Assert response status and body
}
```

### 4. Consistent Error Handling
All errors flow through the service layer:
```go
// Service ensures consistent error messages
if input.Email == "" {
    return nil, fmt.Errorf("email is required")
}

// Handler can handle errors consistently
if err != nil {
    h.logger.Errorf("service error: %v", err)
    http.Error(w, err.Error(), http.StatusInternalServerError)
}
```

### 5. Centralized Logging
All operations log through the service:
```go
s.logger.Debugf("creating user with email: %s", input.Email)
user, err := s.queries.CreateUser(ctx, params)
if err != nil {
    s.logger.Errorf("failed to create user: %v", err)
    return nil, err
}
s.logger.Infof("user created with ID: %d", user.ID)
```

## Input/Output Type Pattern

Services use explicit Input/Output types for API contracts:

### Advantages
```go
// ✓ Clear, explicit inputs
type CreateUserInput struct {
    Email        string
    Name         string
    PasswordHash string
}

// vs

// ✗ Unclear what fields are used
func (s *Service) CreateUser(ctx context.Context, 
    email, name, password string) (...)

// ✓ Clear outputs and structure
type CreateUserOutput struct {
    ID        int64
    Email     string
    Name      string
    CreatedAt string
    UpdatedAt string
}

// vs

// ✗ Unclear what's being returned
func (s *Service) CreateUser(ctx context.Context, ...) (interface{}, error)
```

## Transaction Example

For operations requiring multiple queries:

```go
// Start transaction
tx, err := s.queries.db.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

// Create queries with transaction
txQueries := s.queries.WithTx(tx)

// Execute multiple queries atomically
user, err := txQueries.CreateUser(ctx, userParams)
if err != nil {
    return err
}

order, err := txQueries.CreateOrder(ctx, orderParams)
if err != nil {
    return err
}

// Commit if all succeeded
return tx.Commit(ctx)
```

## Adding New Service Methods

### Step 1: Create SQL Query
File: `internal/infra/queries/users.sql`
```sql
-- name: GetActiveUsers :many
SELECT id, email, name, password_hash, created_at, updated_at
FROM users
WHERE status = 'active'
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
```

### Step 2: Generate SQLC Code
```bash
sqlc generate
```

This creates the `GetActiveUsers` method in the generated code.

### Step 3: Create Service Method
File: `internal/users/service.go`
```go
type GetActiveUsersInput struct {
    Limit  int32
    Offset int32
}

type GetActiveUsersOutput struct {
    Users []UserData
    Limit int32
    Offset int32
}

func (s *Service) GetActiveUsers(ctx context.Context, input *GetActiveUsersInput) (*GetActiveUsersOutput, error) {
    if input.Limit == 0 {
        input.Limit = 10
    }
    
    s.logger.Debugf("fetching active users")
    
    users, err := s.queries.GetActiveUsers(ctx, input.Limit, input.Offset)
    if err != nil {
        s.logger.Errorf("failed to fetch active users: %v", err)
        return nil, fmt.Errorf("failed to fetch active users: %w", err)
    }
    
    output := &GetActiveUsersOutput{
        Users:  make([]UserData, len(users)),
        Limit:  input.Limit,
        Offset: input.Offset,
    }
    
    for i, user := range users {
        output.Users[i] = UserData{
            ID:    user.ID,
            Email: user.Email,
            // ...
        }
    }
    
    return output, nil
}
```

### Step 4: Create HTTP Handler
File: `internal/handlers/user_handler.go`
```go
func (h *UserHandler) HandleGetActiveUsers(w http.ResponseWriter, r *http.Request) {
    limit := int32(10)
    // ... parse limit from query params
    
    input := &users.GetActiveUsersInput{Limit: limit}
    output, err := h.service.GetActiveUsers(r.Context(), input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(output)
}
```

### Step 5: Register Route
File: `internal/handlers/setup.go`
```go
func RegisterUserRoutes(mux *http.ServeMux, handler *UserHandler) {
    // ... existing routes
    mux.HandleFunc("GET /users/active", handler.HandleGetActiveUsers)
}
```

## Best Practices

1. **Always use Input/Output types** - Clear contracts and easier to evolve
2. **Validate early** - Check inputs at service layer entry
3. **Log operations** - Debug and monitor with consistent logging
4. **Use context** - Pass context through all layers
5. **Handle errors explicitly** - Don't ignore or mask errors
6. **Use transactions** - For multi-query operations that must be atomic
7. **Keep services small** - Each service focuses on one entity
8. **Return structs** - Not pointers, unless there's a specific reason

## File Structure

```
internal/
├── users/
│   └── service.go           # User service layer
├── products/
│   └── service.go           # Product service layer
├── orders/
│   └── service.go           # Order service layer
├── handlers/
│   ├── user_handler.go      # User HTTP handlers
│   ├── product_handler.go   # Product HTTP handlers
│   ├── order_handler.go     # Order HTTP handlers
│   └── setup.go             # Handler setup and routing
└── infra/
    └── queries/
        ├── querier.go       # SQLC interface
        ├── models.go        # SQLC generated models
        ├── users_gen.sql.go # SQLC generated users
        └── users.sql        # User query definitions
```

## Related Documentation

- [SQLC Integration Guide](./README.md)
- [User Service](../users/service.go)
- [Product Service](../products/service.go)
- [Order Service](../orders/service.go)
- [User Handler](../handlers/user_handler.go)

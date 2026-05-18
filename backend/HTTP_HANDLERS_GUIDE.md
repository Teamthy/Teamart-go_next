# HTTP Handlers Quick Start

This guide shows how to use the HTTP handlers and service layer to handle API requests.

## Quick Start

### Setup Handlers in Main

```go
package main

import (
    "net/http"
    
    "github.com/teamart/commerce-api/internal/handlers"
    "github.com/teamart/commerce-api/internal/infra/database"
    "github.com/teamart/commerce-api/internal/infra/queries"
)

func main() {
    // Initialize database pool
    db, err := database.NewPool(context.Background(), &cfg.Database, log)
    if err != nil {
        log.Fatalf("failed to initialize database: %v", err)
    }
    
    // Create SQLC queries instance
    q := queries.New(db.Conn())
    
    // Setup HTTP handlers
    mux := http.NewServeMux()
    handlers.RegisterHealthRoutes(mux)      // Health check endpoint
    handlers.SetupHandlers(mux, q, log)     // All service handlers
    
    // Start server
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    log.Infof("server listening on %s", server.Addr)
    server.ListenAndServe()
}
```

## API Endpoints

All endpoints are organized by resource type.

### User Endpoints

#### Create User
```bash
POST /users
Content-Type: application/json

{
  "email": "user@example.com",
  "name": "John Doe",
  "password_hash": "hashed_password_here"
}

Response (201 Created):
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "password_hash": "hashed_password_here",
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:30:00Z"
}
```

#### Get User by ID
```bash
GET /users/1

Response (200 OK):
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "password_hash": "hashed_password_here",
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:30:00Z"
}
```

#### Get User by Email
```bash
GET /users/email/user@example.com

Response (200 OK):
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "password_hash": "hashed_password_here",
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:30:00Z"
}
```

#### List Users (Paginated)
```bash
GET /users?limit=10&offset=0

Response (200 OK):
{
  "users": [
    {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "password_hash": "hashed_password_here",
      "created_at": "2024-05-18T10:30:00Z",
      "updated_at": "2024-05-18T10:30:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

#### Update User
```bash
PUT /users/1
Content-Type: application/json

{
  "name": "Jane Doe"
}

Response (200 OK):
{
  "id": 1,
  "email": "user@example.com",
  "name": "Jane Doe",
  "password_hash": "hashed_password_here",
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:35:00Z"
}
```

#### Delete User
```bash
DELETE /users/1

Response (204 No Content)
```

### Product Endpoints

#### Create Product
```bash
POST /products
Content-Type: application/json

{
  "sku": "PROD001",
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99
}

Response (201 Created):
{
  "id": 1,
  "sku": "PROD001",
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:30:00Z"
}
```

#### Get Product by ID
```bash
GET /products/1

Response (200 OK):
{
  "id": 1,
  "sku": "PROD001",
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:30:00Z"
}
```

#### Get Product by SKU
```bash
GET /products/sku/PROD001

Response (200 OK):
{
  "id": 1,
  "sku": "PROD001",
  "name": "Laptop",
  "description": "High-performance laptop",
  "price": 999.99,
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:30:00Z"
}
```

#### List Products (Paginated)
```bash
GET /products?limit=10&offset=0

Response (200 OK):
{
  "products": [
    {
      "id": 1,
      "sku": "PROD001",
      "name": "Laptop",
      "description": "High-performance laptop",
      "price": 999.99,
      "created_at": "2024-05-18T10:30:00Z",
      "updated_at": "2024-05-18T10:30:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

#### Search Products
```bash
GET /products/search?q=laptop&limit=10&offset=0

Response (200 OK):
{
  "products": [
    {
      "id": 1,
      "sku": "PROD001",
      "name": "Laptop",
      "description": "High-performance laptop",
      "price": 999.99,
      "created_at": "2024-05-18T10:30:00Z",
      "updated_at": "2024-05-18T10:30:00Z"
    }
  ],
  "query": "laptop",
  "limit": 10,
  "offset": 0
}
```

### Order Endpoints

#### Create Order
```bash
POST /orders
Content-Type: application/json

{
  "user_id": 1,
  "total_amount": 1999.98,
  "status": "pending"
}

Response (201 Created):
{
  "id": 1,
  "user_id": 1,
  "total_amount": 1999.98,
  "status": "pending",
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:30:00Z"
}
```

#### Get Order by ID
```bash
GET /orders/1

Response (200 OK):
{
  "id": 1,
  "user_id": 1,
  "total_amount": 1999.98,
  "status": "pending",
  "created_at": "2024-05-18T10:30:00Z",
  "updated_at": "2024-05-18T10:30:00Z"
}
```

#### List All Orders (Paginated)
```bash
GET /orders?limit=10&offset=0

Response (200 OK):
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "total_amount": 1999.98,
      "status": "pending",
      "created_at": "2024-05-18T10:30:00Z",
      "updated_at": "2024-05-18T10:30:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

#### List Orders by User
```bash
GET /users/1/orders?limit=10&offset=0

Response (200 OK):
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "total_amount": 1999.98,
      "status": "pending",
      "created_at": "2024-05-18T10:30:00Z",
      "updated_at": "2024-05-18T10:30:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

#### List Orders by Status
```bash
GET /orders/status/pending?limit=10&offset=0

Response (200 OK):
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "total_amount": 1999.98,
      "status": "pending",
      "created_at": "2024-05-18T10:30:00Z",
      "updated_at": "2024-05-18T10:30:00Z"
    }
  ],
  "limit": 10,
  "offset": 0
}
```

### Health Check Endpoint

#### Check API Health
```bash
GET /health

Response (200 OK):
{
  "status": "healthy"
}
```

## Testing with cURL

### Create Test User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password_hash": "hashed_password"
  }'
```

### Get Test User
```bash
curl http://localhost:8080/users/1
```

### List Users
```bash
curl "http://localhost:8080/users?limit=5&offset=0"
```

### Search Products
```bash
curl "http://localhost:8080/products/search?q=laptop&limit=10"
```

### Create Order
```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "total_amount": 99.99,
    "status": "pending"
  }'
```

## Testing with Postman

Import the API endpoints as a Postman collection:

1. Create a new collection "Teamart Commerce API"
2. Create folders: Users, Products, Orders, Health
3. Add requests for each endpoint

Example request:
```
POST http://localhost:8080/users
Headers:
  Content-Type: application/json
Body (JSON):
{
  "email": "user@example.com",
  "name": "John Doe",
  "password_hash": "hashed_pwd"
}
```

## Testing with Go

### Example Test
```go
package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/teamart/commerce-api/internal/users"
)

func TestHandleCreateUser(t *testing.T) {
    // Mock service
    mockService := &MockUserService{
        CreateUserFn: func(ctx context.Context, input *users.CreateUserInput) (*users.CreateUserOutput, error) {
            return &users.CreateUserOutput{
                ID:    1,
                Email: input.Email,
            }, nil
        },
    }
    
    handler := NewUserHandler(mockService, logger)
    
    // Create request
    body := CreateUserRequest{
        Email:        "test@example.com",
        Name:         "Test User",
        PasswordHash: "hashed",
    }
    bodyBytes, _ := json.Marshal(body)
    req := httptest.NewRequest("POST", "/users", bytes.NewReader(bodyBytes))
    
    // Call handler
    w := httptest.NewRecorder()
    handler.HandleCreateUser(w, req)
    
    // Assert response
    if w.Code != http.StatusCreated {
        t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
    }
}
```

## Error Handling

All endpoints follow consistent error responses:

### Bad Request (400)
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{}'  # Missing required fields

Response (400 Bad Request):
email is required
```

### Not Found (500)
```bash
curl http://localhost:8080/users/999

Response (500 Internal Server Error):
failed to fetch user: sql: no rows in result set
```

## Rate Limiting

No rate limiting is currently implemented. For production, consider:
- Middleware-based rate limiting
- Database-driven rate limit tracking
- Per-user rate limits

## Pagination

All list endpoints support pagination with:
- `limit`: Number of items per page (default: 10, max: 100)
- `offset`: Number of items to skip (default: 0)

```bash
GET /products?limit=20&offset=40  # Get items 40-60
```

## Sorting

Currently, results are sorted by creation date (most recent first) and cannot be customized via API parameters.

## Filtering

Limited filtering is available:
- Search products by name/description: `GET /products/search?q=laptop`
- Filter orders by status: `GET /orders/status/pending`

For more complex filtering, consider adding filtering parameters to list endpoints.

## CORS

CORS is not currently enabled. To enable, add CORS middleware:

```go
import "github.com/rs/cors"

corsMiddleware := cors.Default()
handler := corsMiddleware.Handler(mux)
server.Handler = handler
```

## See Also

- [Service Layer Architecture](./SERVICE_LAYER_ARCHITECTURE.md)
- [SQLC Integration Guide](./internal/infra/queries/README.md)
- [User Service](./internal/users/service.go)
- [Product Service](./internal/products/service.go)
- [Order Service](./internal/orders/service.go)

package queries

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// User represents a user record from the database
type User struct {
	ID           int32     `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	Name         string    `db:"name" json:"name"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// Product represents a product record from the database
type Product struct {
	ID          int32          `db:"id" json:"id"`
	Sku         string         `db:"sku" json:"sku"`
	Name        string         `db:"name" json:"name"`
	Description pgtype.Text    `db:"description" json:"description"`
	Price       pgtype.Numeric `db:"price" json:"price"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

// Order represents an order record from the database
type Order struct {
	ID          int32          `db:"id" json:"id"`
	UserID      int32          `db:"user_id" json:"user_id"`
	TotalAmount pgtype.Numeric `db:"total_amount" json:"total_amount"`
	Status      string         `db:"status" json:"status"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

// OrderItem represents an order item record from the database
type OrderItem struct {
	ID        int32          `db:"id" json:"id"`
	OrderID   int32          `db:"order_id" json:"order_id"`
	ProductID int32          `db:"product_id" json:"product_id"`
	Quantity  int32          `db:"quantity" json:"quantity"`
	Price     pgtype.Numeric `db:"price" json:"price"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
}

package queries

import (
	"context"
	"fmt"
)

// CreateOrder creates a new order
func (q *Queries) CreateOrder(ctx context.Context, userID int32, totalAmount string, status string) (Order, error) {
	row := q.db.QueryRow(ctx,
		`INSERT INTO orders (user_id, total_amount, status)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, total_amount, status, created_at, updated_at`,
		userID, totalAmount, status)
	var i Order
	err := row.Scan(&i.ID, &i.UserID, &i.TotalAmount, &i.Status, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

// GetOrderByID retrieves an order by ID
func (q *Queries) GetOrderByID(ctx context.Context, id int32) (Order, error) {
	row := q.db.QueryRow(ctx,
		`SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders WHERE id = $1`,
		id)
	var i Order
	err := row.Scan(&i.ID, &i.UserID, &i.TotalAmount, &i.Status, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

// ListOrdersByUserID lists orders for a user
func (q *Queries) ListOrdersByUserID(ctx context.Context, userID int32, limit int32, offset int32) ([]Order, error) {
	rows, err := q.db.Query(ctx,
		`SELECT id, user_id, total_amount, status, created_at, updated_at
		 FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(&i.ID, &i.UserID, &i.TotalAmount, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ListOrdersByStatus lists orders by status
func (q *Queries) ListOrdersByStatus(ctx context.Context, status string, limit int32, offset int32) ([]Order, error) {
	rows, err := q.db.Query(ctx,
		`SELECT id, user_id, total_amount, status, created_at, updated_at
		 FROM orders WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(&i.ID, &i.UserID, &i.TotalAmount, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ListAllOrders lists all orders
func (q *Queries) ListAllOrders(ctx context.Context, limit int32, offset int32) ([]Order, error) {
	rows, err := q.db.Query(ctx,
		`SELECT id, user_id, total_amount, status, created_at, updated_at
		 FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Order
	for rows.Next() {
		var i Order
		if err := rows.Scan(&i.ID, &i.UserID, &i.TotalAmount, &i.Status, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// UpdateOrderStatus updates order status
func (q *Queries) UpdateOrderStatus(ctx context.Context, id int32, status string) (Order, error) {
	row := q.db.QueryRow(ctx,
		`UPDATE orders SET status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1
		 RETURNING id, user_id, total_amount, status, created_at, updated_at`,
		id, status)
	var i Order
	err := row.Scan(&i.ID, &i.UserID, &i.TotalAmount, &i.Status, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

// UpdateOrderTotalAmount updates order total amount
func (q *Queries) UpdateOrderTotalAmount(ctx context.Context, id int32, totalAmount string) (Order, error) {
	row := q.db.QueryRow(ctx,
		`UPDATE orders SET total_amount = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1
		 RETURNING id, user_id, total_amount, status, created_at, updated_at`,
		id, totalAmount)
	var i Order
	err := row.Scan(&i.ID, &i.UserID, &i.TotalAmount, &i.Status, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

// DeleteOrder deletes an order
func (q *Queries) DeleteOrder(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, `DELETE FROM orders WHERE id = $1`, id)
	return err
}

// CountOrdersByUserID counts orders for a user
func (q *Queries) CountOrdersByUserID(ctx context.Context, userID int32) (int64, error) {
	row := q.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE user_id = $1`, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

// CountOrdersByStatus counts orders by status
func (q *Queries) CountOrdersByStatus(ctx context.Context, status string) (int64, error) {
	row := q.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE status = $1`, status)
	var count int64
	err := row.Scan(&count)
	return count, err
}

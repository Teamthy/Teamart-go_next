package queries

import (
	"context"
	"fmt"
)

// CreateOrderItem creates a new order item
func (q *Queries) CreateOrderItem(ctx context.Context, orderID int32, productID int32, quantity int32, price string) (OrderItem, error) {
	row := q.db.QueryRow(ctx,
		`INSERT INTO order_items (order_id, product_id, quantity, price)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, order_id, product_id, quantity, price, created_at`,
		orderID, productID, quantity, price)
	var i OrderItem
	err := row.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.Price, &i.CreatedAt)
	return i, err
}

// GetOrderItemByID retrieves an order item by ID
func (q *Queries) GetOrderItemByID(ctx context.Context, id int32) (OrderItem, error) {
	row := q.db.QueryRow(ctx,
		`SELECT id, order_id, product_id, quantity, price, created_at FROM order_items WHERE id = $1`,
		id)
	var i OrderItem
	err := row.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.Price, &i.CreatedAt)
	return i, err
}

// ListOrderItemsByOrderID lists items in an order
func (q *Queries) ListOrderItemsByOrderID(ctx context.Context, orderID int32) ([]OrderItem, error) {
	rows, err := q.db.Query(ctx,
		`SELECT id, order_id, product_id, quantity, price, created_at
		 FROM order_items WHERE order_id = $1 ORDER BY created_at DESC`,
		orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.Price, &i.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// GetOrderItemsByProductID gets all order items for a product
func (q *Queries) GetOrderItemsByProductID(ctx context.Context, productID int32, limit int32, offset int32) ([]OrderItem, error) {
	rows, err := q.db.Query(ctx,
		`SELECT id, order_id, product_id, quantity, price, created_at
		 FROM order_items WHERE product_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.Price, &i.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// UpdateOrderItemQuantity updates quantity for an order item
func (q *Queries) UpdateOrderItemQuantity(ctx context.Context, id int32, quantity int32) (OrderItem, error) {
	row := q.db.QueryRow(ctx,
		`UPDATE order_items SET quantity = $2 WHERE id = $1
		 RETURNING id, order_id, product_id, quantity, price, created_at`,
		id, quantity)
	var i OrderItem
	err := row.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.Price, &i.CreatedAt)
	return i, err
}

// UpdateOrderItemPrice updates price for an order item
func (q *Queries) UpdateOrderItemPrice(ctx context.Context, id int32, price string) (OrderItem, error) {
	row := q.db.QueryRow(ctx,
		`UPDATE order_items SET price = $2 WHERE id = $1
		 RETURNING id, order_id, product_id, quantity, price, created_at`,
		id, price)
	var i OrderItem
	err := row.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.Price, &i.CreatedAt)
	return i, err
}

// DeleteOrderItem deletes an order item
func (q *Queries) DeleteOrderItem(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, `DELETE FROM order_items WHERE id = $1`, id)
	return err
}

// DeleteOrderItems deletes all items in an order
func (q *Queries) DeleteOrderItems(ctx context.Context, orderID int32) error {
	_, err := q.db.Exec(ctx, `DELETE FROM order_items WHERE order_id = $1`, orderID)
	return err
}

// CountOrderItems counts items in an order
func (q *Queries) CountOrderItems(ctx context.Context, orderID int32) (int64, error) {
	row := q.db.QueryRow(ctx, `SELECT COUNT(*) FROM order_items WHERE order_id = $1`, orderID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

// GetOrderItemsWithProductRow represents order items with product details
type GetOrderItemsWithProductRow struct {
	ID          int32
	OrderID     int32
	ProductID   int32
	Quantity    int32
	Price       string
	CreatedAt   int64
	Sku         string
	Name        string
	Description string
}

// GetOrderItemsWithProduct gets order items with product information
func (q *Queries) GetOrderItemsWithProduct(ctx context.Context, orderID int32) ([]GetOrderItemsWithProductRow, error) {
	rows, err := q.db.Query(ctx,
		`SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.created_at,
		        p.sku, p.name, p.description
		 FROM order_items oi
		 JOIN products p ON oi.product_id = p.id
		 WHERE oi.order_id = $1
		 ORDER BY oi.created_at DESC`,
		orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []GetOrderItemsWithProductRow
	for rows.Next() {
		var i GetOrderItemsWithProductRow
		if err := rows.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.Price, &i.CreatedAt,
			&i.Sku, &i.Name, &i.Description); err != nil {
			return nil, fmt.Errorf("failed to scan order item with product: %w", err)
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

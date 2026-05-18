package migrations

// Migrations is the collection of all database migrations
// Each migration is versioned and can be applied or rolled back
// Migrations are executed in version order
var Migrations = []Migration{
	migrationInitSchema,
}

// migrationInitSchema creates the initial database schema
// This is the first migration - sets up the basic tables
var migrationInitSchema = Migration{
	Version:     "20260101000000",
	Name:        "init_schema",
	Description: "Create initial database schema with core tables",

	// Up migration: Create tables
	UpSQL: `
	-- Users table
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX idx_users_email ON users(email);

	-- Products table
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		sku VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price NUMERIC(10, 2) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX idx_products_sku ON products(sku);

	-- Orders table
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
		total_amount NUMERIC(10, 2) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX idx_orders_user_id ON orders(user_id);
	CREATE INDEX idx_orders_status ON orders(status);

	-- Order items table
	CREATE TABLE IF NOT EXISTS order_items (
		id SERIAL PRIMARY KEY,
		order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
		quantity INTEGER NOT NULL,
		price NUMERIC(10, 2) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX idx_order_items_order_id ON order_items(order_id);
	CREATE INDEX idx_order_items_product_id ON order_items(product_id);
	`,

	// Down migration: Drop tables
	DownSQL: `
	DROP TABLE IF EXISTS order_items CASCADE;
	DROP TABLE IF EXISTS orders CASCADE;
	DROP TABLE IF EXISTS products CASCADE;
	DROP TABLE IF EXISTS users CASCADE;
	`,
}

// AddMigration adds a new migration to the collection
// Use this to extend migrations programmatically
func AddMigration(m Migration) {
	Migrations = append(Migrations, m)
}

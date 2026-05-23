package migrations

// Migrations is the collection of all database migrations
// Each migration is versioned and can be applied or rolled back
// Migrations are executed in version order
var Migrations = []Migration{
	migrationInitSchema,
	migrationTenantSchema,
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

// migrationTenantSchema creates merchant, tenant, and staff tables
var migrationTenantSchema = Migration{
	Version:     "20260523000000",
	Name:        "merchant_tenant_schema",
	Description: "Add merchant, store, staff, tenant settings, KYC, and payout account tables",

	UpSQL: `
	-- Merchants table
	CREATE TABLE IF NOT EXISTS merchants (
		id SERIAL PRIMARY KEY,
		owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
		name VARCHAR(255) NOT NULL,
		slug VARCHAR(255) NOT NULL UNIQUE,
		description TEXT,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		billing_plan VARCHAR(100) NOT NULL DEFAULT 'starter',
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX idx_merchants_owner_id ON merchants(owner_id);

	-- Stores table
	CREATE TABLE IF NOT EXISTS stores (
		id SERIAL PRIMARY KEY,
		merchant_id INTEGER NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
		owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
		creator_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
		name VARCHAR(255) NOT NULL,
		slug VARCHAR(255) NOT NULL,
		description TEXT,
		category VARCHAR(100),
		settings JSONB DEFAULT '{}'::JSONB,
		status VARCHAR(50) NOT NULL DEFAULT 'active',
		storefront_url VARCHAR(255),
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (merchant_id, slug)
	);
	CREATE INDEX idx_stores_merchant_id ON stores(merchant_id);
	CREATE INDEX idx_stores_owner_id ON stores(owner_id);

	-- Staff accounts table
	CREATE TABLE IF NOT EXISTS staff_accounts (
		id SERIAL PRIMARY KEY,
		merchant_id INTEGER NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		role VARCHAR(100) NOT NULL DEFAULT 'staff',
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (merchant_id, user_id)
	);
	CREATE INDEX idx_staff_accounts_merchant_id ON staff_accounts(merchant_id);
	CREATE INDEX idx_staff_accounts_user_id ON staff_accounts(user_id);

	-- Store members table
	CREATE TABLE IF NOT EXISTS store_members (
		id SERIAL PRIMARY KEY,
		store_id INTEGER NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
		staff_account_id INTEGER NOT NULL REFERENCES staff_accounts(id) ON DELETE CASCADE,
		role VARCHAR(100) NOT NULL DEFAULT 'staff',
		permissions TEXT[] DEFAULT ARRAY[]::TEXT[],
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		left_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (store_id, staff_account_id)
	);
	CREATE INDEX idx_store_members_store_id ON store_members(store_id);
	CREATE INDEX idx_store_members_staff_account_id ON store_members(staff_account_id);

	-- Tenant settings table
	CREATE TABLE IF NOT EXISTS tenant_settings (
		id SERIAL PRIMARY KEY,
		tenant_id INTEGER NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
		key VARCHAR(255) NOT NULL,
		value JSONB NOT NULL DEFAULT '{}'::JSONB,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (tenant_id, key)
	);
	CREATE INDEX idx_tenant_settings_tenant_id ON tenant_settings(tenant_id);

	-- Merchant KYC table
	CREATE TABLE IF NOT EXISTS merchant_kyc (
		id SERIAL PRIMARY KEY,
		merchant_id INTEGER NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
		legal_name VARCHAR(255) NOT NULL,
		tax_id VARCHAR(255),
		business_type VARCHAR(100),
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		submitted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		reviewed_at TIMESTAMP WITH TIME ZONE,
		approved_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX idx_merchant_kyc_merchant_id ON merchant_kyc(merchant_id);

	-- Merchant payout accounts table
	CREATE TABLE IF NOT EXISTS merchant_payout_accounts (
		id SERIAL PRIMARY KEY,
		merchant_id INTEGER NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
		provider VARCHAR(100) NOT NULL,
		account_holder_name VARCHAR(255) NOT NULL,
		account_type VARCHAR(100) NOT NULL,
		external_account_id VARCHAR(255) NOT NULL,
		currency VARCHAR(10) NOT NULL DEFAULT 'USD',
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		metadata JSONB DEFAULT '{}'::JSONB,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX idx_merchant_payout_accounts_merchant_id ON merchant_payout_accounts(merchant_id);
	`,

	DownSQL: `
	DROP TABLE IF EXISTS merchant_payout_accounts CASCADE;
	DROP TABLE IF EXISTS merchant_kyc CASCADE;
	DROP TABLE IF EXISTS tenant_settings CASCADE;
	DROP TABLE IF EXISTS store_members CASCADE;
	DROP TABLE IF EXISTS staff_accounts CASCADE;
	DROP TABLE IF EXISTS stores CASCADE;
	DROP TABLE IF EXISTS merchants CASCADE;
	`,
}

// AddMigration adds a new migration to the collection
// Use this to extend migrations programmatically
func AddMigration(m Migration) {
	Migrations = append(Migrations, m)
}

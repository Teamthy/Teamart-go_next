#!/bin/bash
# Migration: 004_order_management_system.sql
# Description: Create comprehensive order management system with cart, checkout, fulfillment
# Created: 2026-05-21

set -e

psql "$DATABASE_URL" <<'EOF'

-- ===== SHOPPING CART TABLES =====

CREATE TABLE IF NOT EXISTS carts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    store_id BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL, -- For direct purchases, NULL for multi-store
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, abandoned, converted
    subtotal DECIMAL(12, 2) NOT NULL DEFAULT 0,
    shipping_amount DECIMAL(12, 2) DEFAULT 0,
    tax_amount DECIMAL(12, 2) DEFAULT 0,
    discount_amount DECIMAL(12, 2) DEFAULT 0,
    total DECIMAL(12, 2) NOT NULL DEFAULT 0,
    coupon_code VARCHAR(100),
    estimated_shipping_method VARCHAR(100),
    estimated_shipping_cost DECIMAL(12, 2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    abandoned_at TIMESTAMP,
    converted_at TIMESTAMP
);

CREATE INDEX idx_carts_user_id ON carts(user_id);
CREATE INDEX idx_carts_status ON carts(status);
CREATE INDEX idx_carts_abandoned_at ON carts(abandoned_at);

-- ===== CART ITEMS =====

CREATE TABLE IF NOT EXISTS cart_items (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    variant_id BIGINT NOT NULL REFERENCES product_variants(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12, 2) NOT NULL,
    line_total DECIMAL(12, 2) NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_items_product_id ON cart_items(product_id);

-- ===== ORDERS =====

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(50) NOT NULL UNIQUE, -- e.g., ORD-2026-0001
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Order source
    source VARCHAR(50) NOT NULL DEFAULT 'web', -- web, mobile, livestream, admin
    creator_id BIGINT REFERENCES users(id) ON DELETE SET NULL, -- For livestream orders
    
    -- Customer info (denormalized for orders)
    customer_email VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(20),
    
    -- Status machine
    status VARCHAR(50) NOT NULL DEFAULT 'pending', 
    -- pending → paid → processing → shipped → delivered → completed
    -- Can also → cancelled, refunded, disputed at any point
    
    -- Financial
    subtotal DECIMAL(12, 2) NOT NULL,
    shipping_cost DECIMAL(12, 2) DEFAULT 0,
    tax_amount DECIMAL(12, 2) DEFAULT 0,
    discount_amount DECIMAL(12, 2) DEFAULT 0,
    total DECIMAL(12, 2) NOT NULL,
    paid_amount DECIMAL(12, 2) DEFAULT 0,
    
    -- Payment
    payment_method VARCHAR(50), -- card, wallet, bank_transfer, upi, etc.
    payment_id VARCHAR(255), -- External payment gateway ID
    payment_status VARCHAR(50) DEFAULT 'pending', -- pending, completed, failed
    
    -- Shipping
    shipping_method VARCHAR(100),
    tracking_number VARCHAR(255),
    estimated_delivery_date DATE,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    
    -- Billing & Shipping addresses
    billing_address JSONB,
    shipping_address JSONB,
    
    -- Notes & metadata
    notes TEXT,
    metadata JSONB,
    
    -- Cancellation & refund tracking
    cancelled_at TIMESTAMP,
    cancellation_reason VARCHAR(255),
    refunded_at TIMESTAMP,
    refund_reason VARCHAR(255),
    refund_amount DECIMAL(12, 2),
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_creator_id ON orders(creator_id);

-- ===== ORDER ITEMS =====

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    variant_id BIGINT NOT NULL REFERENCES product_variants(id),
    sku VARCHAR(255) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    variant_title VARCHAR(255),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(12, 2) NOT NULL,
    line_total DECIMAL(12, 2) NOT NULL,
    
    -- Fulfillment per item
    fulfillment_status VARCHAR(50) DEFAULT 'pending', -- pending, processing, shipped, delivered
    fulfillment_id VARCHAR(255), -- Link to fulfillment
    
    -- Returns
    return_status VARCHAR(50), -- none, pending, approved, rejected, completed
    return_quantity INTEGER DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_fulfillment_status ON order_items(fulfillment_status);

-- ===== ORDER TIMELINE / STATUS HISTORY =====

CREATE TABLE IF NOT EXISTS order_events (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- payment_received, shipped, delivered, cancelled, refunded, etc.
    status VARCHAR(50) NOT NULL, -- Old status
    new_status VARCHAR(50) NOT NULL, -- New status
    description TEXT,
    actor_id BIGINT REFERENCES users(id) ON DELETE SET NULL, -- Who made the change
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_events_order_id ON order_events(order_id);
CREATE INDEX idx_order_events_event_type ON order_events(event_type);

-- ===== FULFILLMENT / SHIPPING =====

CREATE TABLE IF NOT EXISTS fulfillments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, shipped, delivered, failed
    shipping_carrier VARCHAR(100),
    tracking_number VARCHAR(255),
    tracking_url VARCHAR(500),
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    estimated_delivery_date DATE,
    
    -- Address
    shipping_address JSONB,
    
    -- Items in this fulfillment (can have multiple fulfillments per order)
    items_count INTEGER,
    
    -- Notes
    notes TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fulfillments_order_id ON fulfillments(order_id);
CREATE INDEX idx_fulfillments_status ON fulfillments(status);
CREATE INDEX idx_fulfillments_tracking_number ON fulfillments(tracking_number);

-- ===== RETURNS & REFUNDS =====

CREATE TABLE IF NOT EXISTS order_returns (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', 
    -- pending → approved → received → refunded → completed
    
    reason VARCHAR(255) NOT NULL,
    description TEXT,
    
    initiator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT, -- Customer or support
    initiated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Authorization
    authorized_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    authorized_at TIMESTAMP,
    authorization_notes TEXT,
    
    -- Receipt & inspection
    received_at TIMESTAMP,
    received_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    inspection_notes TEXT,
    
    -- Refund
    refund_amount DECIMAL(12, 2),
    refund_status VARCHAR(50), -- pending, processed, failed
    refund_id VARCHAR(255), -- Payment gateway refund ID
    refunded_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_returns_order_id ON order_returns(order_id);
CREATE INDEX idx_order_returns_status ON order_returns(status);

-- ===== RETURN ITEMS =====

CREATE TABLE IF NOT EXISTS return_items (
    id BIGSERIAL PRIMARY KEY,
    return_id BIGINT NOT NULL REFERENCES order_returns(id) ON DELETE CASCADE,
    order_item_id BIGINT NOT NULL REFERENCES order_items(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    reason VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_return_items_return_id ON return_items(return_id);

-- ===== DISPUTES / CHARGEBACKS =====

CREATE TABLE IF NOT EXISTS order_disputes (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    -- open → under_review → resolved → escalated → closed
    
    dispute_type VARCHAR(100) NOT NULL, -- chargeback, not_received, quality_issue, fraud, etc.
    reason TEXT NOT NULL,
    
    -- Parties
    initiator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT, -- Customer or support
    initiated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Resolution
    assigned_to BIGINT REFERENCES users(id) ON DELETE SET NULL,
    resolution TEXT,
    resolved_at TIMESTAMP,
    resolved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    
    -- Amount
    disputed_amount DECIMAL(12, 2),
    resolved_amount DECIMAL(12, 2),
    
    -- Timeline
    deadline DATE,
    evidence TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_disputes_order_id ON order_disputes(order_id);
CREATE INDEX idx_order_disputes_status ON order_disputes(status);
CREATE INDEX idx_order_disputes_assigned_to ON order_disputes(assigned_to);

-- ===== DISCOUNTS & COUPONS =====

CREATE TABLE IF NOT EXISTS coupons (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    
    -- Discount type
    discount_type VARCHAR(50) NOT NULL, -- percentage, fixed_amount, free_shipping
    discount_value DECIMAL(12, 2) NOT NULL,
    
    -- Limits
    min_order_value DECIMAL(12, 2),
    max_discount_amount DECIMAL(12, 2),
    max_uses INTEGER,
    max_uses_per_customer INTEGER DEFAULT 1,
    
    -- Validity
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    
    -- Usage tracking
    times_used INTEGER DEFAULT 0,
    
    -- Status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coupons_code ON coupons(code);
CREATE INDEX idx_coupons_is_active ON coupons(is_active);

-- ===== COUPON USAGE TRACKING =====

CREATE TABLE IF NOT EXISTS coupon_usages (
    id BIGSERIAL PRIMARY KEY,
    coupon_id BIGINT NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    discount_amount DECIMAL(12, 2) NOT NULL,
    used_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coupon_usages_coupon_id ON coupon_usages(coupon_id);
CREATE INDEX idx_coupon_usages_user_id ON coupon_usages(user_id);

-- ===== ORDER NOTES / COMMENTS =====

CREATE TABLE IF NOT EXISTS order_notes (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    note_type VARCHAR(50) NOT NULL DEFAULT 'comment', -- comment, internal, system
    content TEXT NOT NULL,
    is_visible_to_customer BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_notes_order_id ON order_notes(order_id);

-- ===== GIFT CARDS =====

CREATE TABLE IF NOT EXISTS gift_cards (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    balance DECIMAL(12, 2) NOT NULL,
    original_balance DECIMAL(12, 2) NOT NULL,
    
    -- Ownership
    owner_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    recipient_email VARCHAR(255),
    
    -- Validity
    valid_from TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP,
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, redeemed, expired, cancelled
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_gift_cards_code ON gift_cards(code);
CREATE INDEX idx_gift_cards_owner_id ON gift_cards(owner_id);

-- ===== GIFT CARD USAGE =====

CREATE TABLE IF NOT EXISTS gift_card_usages (
    id BIGSERIAL PRIMARY KEY,
    gift_card_id BIGINT NOT NULL REFERENCES gift_cards(id) ON DELETE CASCADE,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    amount_used DECIMAL(12, 2) NOT NULL,
    used_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_gift_card_usages_gift_card_id ON gift_card_usages(gift_card_id);

-- ===== SHIPPING RATES =====

CREATE TABLE IF NOT EXISTS shipping_rates (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Conditions
    min_weight DECIMAL(8, 3),
    max_weight DECIMAL(8, 3),
    min_order_value DECIMAL(12, 2),
    max_order_value DECIMAL(12, 2),
    
    -- Destinations
    countries TEXT, -- JSON array of country codes
    
    -- Pricing
    base_rate DECIMAL(12, 2) NOT NULL,
    per_unit_rate DECIMAL(12, 2),
    
    -- Status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_shipping_rates_store_id ON shipping_rates(store_id);

EOF

echo "Migration completed: created order management system"

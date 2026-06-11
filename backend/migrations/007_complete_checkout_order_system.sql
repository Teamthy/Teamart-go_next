-- Migration: 007_complete_checkout_order_system.sql
-- Description: Complete checkout, payment, shipping, fulfillment, returns, support, and analytics system
-- This migration extends and normalizes the order management system with proper payment processing,
-- multi-carrier shipping, fulfillment workflows, returns/refunds, support tickets, and settlement tracking
-- Created: 2026-06-10

-- ===== PAYMENT SYSTEM =====

-- Payment intents for tracking payment lifecycle
CREATE TABLE IF NOT EXISTS payment_intents (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    external_id VARCHAR(255) UNIQUE, -- Payment gateway ID (Stripe, Paystack, etc.)
    provider VARCHAR(50) NOT NULL, -- stripe, paystack, flutterwave, apple_pay, google_pay, wallet, bank_transfer
    method VARCHAR(50) NOT NULL, -- card, bank_transfer, wallet, mobile_money, etc.
    amount_cents BIGINT NOT NULL, -- Store amounts in cents for precision
    currency VARCHAR(3) NOT NULL DEFAULT 'NGN',
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, authorized, captured, failed, cancelled
    metadata JSONB, -- Flexible storage for provider-specific data
    authorized_at TIMESTAMP,
    captured_at TIMESTAMP,
    failed_at TIMESTAMP,
    failure_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_intents_order_id ON payment_intents(order_id);
CREATE INDEX idx_payment_intents_external_id ON payment_intents(external_id);
CREATE INDEX idx_payment_intents_provider ON payment_intents(provider);
CREATE INDEX idx_payment_intents_status ON payment_intents(status);
CREATE INDEX idx_payment_intents_created_at ON payment_intents(created_at);

-- Payment transactions record each step (auth, capture, refund, etc.)
CREATE TABLE IF NOT EXISTS payment_transactions (
    id BIGSERIAL PRIMARY KEY,
    payment_intent_id BIGINT NOT NULL REFERENCES payment_intents(id) ON DELETE CASCADE,
    transaction_type VARCHAR(50) NOT NULL, -- auth, capture, refund, void, dispute
    provider_transaction_id VARCHAR(255), -- External transaction ID
    amount_cents BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, completed, failed, cancelled
    provider_response JSONB, -- Full provider response
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_transactions_payment_intent_id ON payment_transactions(payment_intent_id);
CREATE INDEX idx_payment_transactions_transaction_type ON payment_transactions(transaction_type);

-- Comprehensive audit trail for all payment operations
CREATE TABLE IF NOT EXISTS payment_audits (
    id BIGSERIAL PRIMARY KEY,
    payment_intent_id BIGINT NOT NULL REFERENCES payment_intents(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL, -- create_intent, authorize, capture, refund, fail, etc.
    actor_id BIGINT REFERENCES users(id) ON DELETE SET NULL, -- User who triggered action (admin, system, etc.)
    actor_type VARCHAR(50) NOT NULL DEFAULT 'system', -- system, user, admin, automated
    message TEXT,
    payload JSONB, -- Request/response data
    ip_address VARCHAR(45),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_audits_payment_intent_id ON payment_audits(payment_intent_id);
CREATE INDEX idx_payment_audits_action ON payment_audits(action);
CREATE INDEX idx_payment_audits_created_at ON payment_audits(created_at);

-- ===== SHIPPING CONFIGURATION & PROFILES =====

-- Merchant shipping profiles define zones, rates, and carrier rules
CREATE TABLE IF NOT EXISTS shipping_profiles (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Shipping zones (stored as JSON for flexibility)
    zones JSONB NOT NULL DEFAULT '[]', -- [{name: "Local", states: [...]}, ...]
    
    -- Rate calculation rules
    rates JSONB NOT NULL DEFAULT '[]', -- [{zone: "Local", baseRate: 500, weightRate: 10, minDays: 1, maxDays: 3}, ...]
    
    -- Carrier selection rules
    carriers JSONB NOT NULL DEFAULT '[]', -- [{carrier: "DHL", enabled: true, priority: 1}, ...]
    
    -- Fallback/default carrier
    default_carrier VARCHAR(100),
    
    -- Settings
    free_shipping_threshold_cents BIGINT, -- Minimum order amount for free shipping
    exclude_services JSONB DEFAULT '[]', -- Exclude certain services/zones
    
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_shipping_profiles_merchant_id ON shipping_profiles(merchant_id);
CREATE INDEX idx_shipping_profiles_is_default ON shipping_profiles(is_default);
CREATE INDEX idx_shipping_profiles_is_active ON shipping_profiles(is_active);

-- ===== SHIPMENTS & TRACKING =====

-- Detailed shipment information for each order (can be multiple per order)
CREATE TABLE IF NOT EXISTS shipments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    shipment_number VARCHAR(100) UNIQUE, -- e.g., SHP-2026-0001
    
    -- Carrier and routing
    carrier VARCHAR(100) NOT NULL, -- DHL, FedEx, UPS, GIG, Jumia Logistics, Custom
    carrier_service VARCHAR(255), -- Specific service level (Standard, Express, Overnight)
    tracking_number VARCHAR(255) UNIQUE,
    
    -- Shipping profile reference
    shipping_profile_id BIGINT REFERENCES shipping_profiles(id) ON DELETE SET NULL,
    
    -- Package details
    weight_grams INT,
    length_cm INT,
    width_cm INT,
    height_cm INT,
    items_count INT NOT NULL DEFAULT 0,
    
    -- Address
    shipping_address JSONB NOT NULL, -- {name, phone, street, city, state, postalCode, country}
    
    -- Cost
    shipping_cost_cents BIGINT NOT NULL,
    insurance_cost_cents BIGINT DEFAULT 0,
    total_cost_cents BIGINT NOT NULL,
    
    -- Labels & documentation
    label_url VARCHAR(500),
    label_generated_at TIMESTAMP,
    packing_slip_url VARCHAR(500),
    invoice_url VARCHAR(500),
    
    -- Timeline
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, label_generated, picked, packed, shipped, in_transit, delivered, cancelled
    expected_delivery_date DATE,
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancellation_reason TEXT,
    
    -- Notes
    notes TEXT,
    internal_notes TEXT,
    
    -- Metadata for carrier-specific data
    metadata JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_tracking_number ON shipments(tracking_number);
CREATE INDEX idx_shipments_carrier ON shipments(carrier);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipments_shipped_at ON shipments(shipped_at);
CREATE INDEX idx_shipments_created_at ON shipments(created_at);

-- Tracking events for real-time shipment status updates
CREATE TABLE IF NOT EXISTS tracking_events (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL REFERENCES shipments(id) ON DELETE CASCADE,
    
    -- Event details
    status VARCHAR(50) NOT NULL, -- picked_up, in_transit, out_for_delivery, delivered, failed_attempt, etc.
    location VARCHAR(255), -- Last known location
    location_details JSONB, -- {city, state, country, coordinates}
    
    -- Timeline
    event_time TIMESTAMP NOT NULL,
    
    -- Description and metadata
    description TEXT,
    carrier_code VARCHAR(100), -- Carrier-specific event code
    metadata JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tracking_events_shipment_id ON tracking_events(shipment_id);
CREATE INDEX idx_tracking_events_status ON tracking_events(status);
CREATE INDEX idx_tracking_events_event_time ON tracking_events(event_time);

-- ===== RETURNS SYSTEM =====

-- Return requests from customers
CREATE TABLE IF NOT EXISTS returns (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    return_number VARCHAR(100) UNIQUE,
    
    -- Request details
    return_type VARCHAR(50) NOT NULL, -- return, refund, replacement, exchange
    reason VARCHAR(100) NOT NULL, -- damaged, wrong_item, not_received, quality_issue, missing_parts, size_issue, color_issue, other
    description TEXT,
    
    -- Requestor
    requested_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    requested_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Workflow
    status VARCHAR(50) NOT NULL DEFAULT 'requested', -- requested, approved, rejected, escalated, in_transit, received, inspected, resolved, cancelled
    
    -- Approval
    approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP,
    approval_notes TEXT,
    
    rejected_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    rejected_at TIMESTAMP,
    rejection_reason TEXT,
    
    -- Receipt and inspection
    return_tracking_number VARCHAR(255),
    received_at TIMESTAMP,
    received_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    inspection_notes TEXT,
    
    -- Resolution
    resolution_type VARCHAR(50), -- refund, replacement, store_credit, compensation
    resolved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMP,
    resolution_notes TEXT,
    
    -- Refund details (if applicable)
    refund_amount_cents BIGINT,
    refund_id VARCHAR(255),
    
    -- Items being returned (reference count, details in return_items)
    items_count INT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_returns_order_id ON returns(order_id);
CREATE INDEX idx_returns_status ON returns(status);
CREATE INDEX idx_returns_requested_by ON returns(requested_by);
CREATE INDEX idx_returns_created_at ON returns(created_at);

-- ===== REFUNDS SYSTEM =====

-- Track all refunds separately for accounting and reconciliation
CREATE TABLE IF NOT EXISTS refunds (
    id BIGSERIAL PRIMARY KEY,
    payment_intent_id BIGINT NOT NULL REFERENCES payment_intents(id) ON DELETE CASCADE,
    return_id BIGINT REFERENCES returns(id) ON DELETE SET NULL,
    
    -- Refund details
    amount_cents BIGINT NOT NULL,
    reason VARCHAR(100) NOT NULL, -- return, cancellation, fraud, dispute, correction, other
    description TEXT,
    
    -- Processing
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed, cancelled
    provider_refund_id VARCHAR(255), -- Payment gateway refund ID
    
    -- Timeline
    requested_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    completed_at TIMESTAMP,
    
    -- Error handling
    failure_reason TEXT,
    retry_count INT DEFAULT 0,
    last_retry_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_refunds_payment_intent_id ON refunds(payment_intent_id);
CREATE INDEX idx_refunds_return_id ON refunds(return_id);
CREATE INDEX idx_refunds_status ON refunds(status);
CREATE INDEX idx_refunds_created_at ON refunds(created_at);

-- ===== SUPPORT SYSTEM =====

-- Support tickets linked to orders
CREATE TABLE IF NOT EXISTS support_tickets (
    id BIGSERIAL PRIMARY KEY,
    ticket_number VARCHAR(100) UNIQUE,
    
    -- Parties involved
    merchant_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    customer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    order_id BIGINT REFERENCES orders(id) ON DELETE SET NULL,
    
    -- Ticket details
    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(100), -- shipping, payment, product_quality, missing_items, damaged, returns, billing, account, other
    priority VARCHAR(50) NOT NULL DEFAULT 'normal', -- low, normal, high, urgent
    status VARCHAR(50) NOT NULL DEFAULT 'open', -- open, pending, in_progress, waiting_customer, resolved, closed, escalated
    
    -- Assignment
    assigned_to BIGINT REFERENCES users(id) ON DELETE SET NULL,
    assigned_at TIMESTAMP,
    
    -- SLA tracking
    sla_due_at TIMESTAMP,
    sla_met BOOLEAN,
    
    -- Resolution
    resolved_at TIMESTAMP,
    resolution_summary TEXT,
    resolved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    
    -- Escalation
    escalated_to BIGINT REFERENCES users(id) ON DELETE SET NULL,
    escalation_reason TEXT,
    escalated_at TIMESTAMP,
    
    -- Attachments count
    attachments_count INT DEFAULT 0,
    
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP
);

CREATE INDEX idx_support_tickets_ticket_number ON support_tickets(ticket_number);
CREATE INDEX idx_support_tickets_merchant_id ON support_tickets(merchant_id);
CREATE INDEX idx_support_tickets_customer_id ON support_tickets(customer_id);
CREATE INDEX idx_support_tickets_order_id ON support_tickets(order_id);
CREATE INDEX idx_support_tickets_status ON support_tickets(status);
CREATE INDEX idx_support_tickets_priority ON support_tickets(priority);
CREATE INDEX idx_support_tickets_assigned_to ON support_tickets(assigned_to);
CREATE INDEX idx_support_tickets_created_at ON support_tickets(created_at);

-- Support messages in ticket conversations
CREATE TABLE IF NOT EXISTS support_messages (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    
    -- Sender
    sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    sender_role VARCHAR(50) NOT NULL, -- customer, seller, support_agent, admin
    
    -- Content
    message TEXT NOT NULL,
    attachments_count INT DEFAULT 0,
    
    -- Status
    is_internal BOOLEAN NOT NULL DEFAULT FALSE, -- Internal note visible only to support team
    marked_as_solution BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_support_messages_ticket_id ON support_messages(ticket_id);
CREATE INDEX idx_support_messages_sender_id ON support_messages(sender_id);
CREATE INDEX idx_support_messages_created_at ON support_messages(created_at);

-- ===== FULFILLMENT SYSTEM =====

-- Batch operations for fulfillment (packing multiple orders together)
CREATE TABLE IF NOT EXISTS fulfillment_batches (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    batch_number VARCHAR(100) UNIQUE,
    
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, picking, packing, ready_to_ship, shipped, completed, cancelled
    
    -- Items in batch
    orders_count INT NOT NULL DEFAULT 0,
    items_count INT NOT NULL DEFAULT 0,
    
    -- Timeline
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    
    -- Notes
    notes TEXT
);

CREATE INDEX idx_fulfillment_batches_merchant_id ON fulfillment_batches(merchant_id);
CREATE INDEX idx_fulfillment_batches_status ON fulfillment_batches(status);

-- Inventory reservations for orders
CREATE TABLE IF NOT EXISTS inventory_reservations (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id BIGINT NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    
    quantity_reserved INT NOT NULL,
    quantity_released INT DEFAULT 0,
    
    status VARCHAR(50) NOT NULL DEFAULT 'reserved', -- reserved, partially_released, fully_released, expired
    
    reserved_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    release_requested_at TIMESTAMP,
    released_at TIMESTAMP,
    expires_at TIMESTAMP,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_inventory_reservations_order_id ON inventory_reservations(order_id);
CREATE INDEX idx_inventory_reservations_product_id ON inventory_reservations(product_id);
CREATE INDEX idx_inventory_reservations_status ON inventory_reservations(status);

-- ===== SETTLEMENT & PAYOUTS =====

-- Seller settlement periods and summaries
CREATE TABLE IF NOT EXISTS seller_settlements (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    settlement_number VARCHAR(100) UNIQUE,
    
    -- Period
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    -- Financial summary
    gross_sales_cents BIGINT NOT NULL DEFAULT 0,
    order_count INT NOT NULL DEFAULT 0,
    
    -- Deductions
    platform_fees_cents BIGINT NOT NULL DEFAULT 0,
    return_refunds_cents BIGINT NOT NULL DEFAULT 0,
    chargebacks_cents BIGINT NOT NULL DEFAULT 0,
    reversals_cents BIGINT NOT NULL DEFAULT 0,
    
    -- Net calculation
    net_payout_cents BIGINT NOT NULL DEFAULT 0,
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, reviewed, approved, processing, paid, failed, cancelled
    
    -- Processing
    payment_method VARCHAR(50), -- bank_transfer, wallet, check, etc.
    scheduled_payout_date DATE,
    actual_payout_date DATE,
    payout_reference VARCHAR(255),
    
    notes TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_seller_settlements_merchant_id ON seller_settlements(merchant_id);
CREATE INDEX idx_seller_settlements_period_start ON seller_settlements(period_start);
CREATE INDEX idx_seller_settlements_status ON seller_settlements(status);

-- Individual payout transactions
CREATE TABLE IF NOT EXISTS payouts (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    settlement_id BIGINT REFERENCES seller_settlements(id) ON DELETE SET NULL,
    
    amount_cents BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'NGN',
    
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, scheduled, in_transit, completed, failed, reversed
    
    payment_method VARCHAR(50) NOT NULL, -- bank_transfer, wallet, check, etc.
    payment_details JSONB, -- Bank account, wallet address, etc.
    
    provider VARCHAR(100), -- Payment gateway/provider name
    provider_reference VARCHAR(255),
    provider_response JSONB,
    
    scheduled_at TIMESTAMP,
    sent_at TIMESTAMP,
    completed_at TIMESTAMP,
    failed_at TIMESTAMP,
    
    failure_reason TEXT,
    retry_count INT DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payouts_merchant_id ON payouts(merchant_id);
CREATE INDEX idx_payouts_settlement_id ON payouts(settlement_id);
CREATE INDEX idx_payouts_status ON payouts(status);
CREATE INDEX idx_payouts_created_at ON payouts(created_at);

-- ===== ANALYTICS & REPORTING =====

-- Denormalized analytics table for fast dashboard queries (updated via triggers/jobs)
CREATE TABLE IF NOT EXISTS order_analytics (
    id BIGSERIAL PRIMARY KEY,
    merchant_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date_bucket DATE NOT NULL, -- Daily aggregation
    
    -- Order metrics
    orders_created INT DEFAULT 0,
    orders_completed INT DEFAULT 0,
    orders_cancelled INT DEFAULT 0,
    
    -- Financial metrics
    gross_revenue_cents BIGINT DEFAULT 0,
    refund_amount_cents BIGINT DEFAULT 0,
    net_revenue_cents BIGINT DEFAULT 0,
    
    -- Item metrics
    items_sold INT DEFAULT 0,
    items_returned INT DEFAULT 0,
    
    -- Average metrics
    avg_order_value_cents BIGINT,
    avg_discount_cents BIGINT,
    
    -- Return metrics
    return_rate DECIMAL(5, 2),
    refund_rate DECIMAL(5, 2),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_analytics_merchant_id ON order_analytics(merchant_id);
CREATE INDEX idx_order_analytics_date_bucket ON order_analytics(date_bucket);

-- ===== FRAUD & COMPLIANCE =====

-- Fraud flags and risk assessment
CREATE TABLE IF NOT EXISTS fraud_flags (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    
    -- Risk assessment
    risk_score INT, -- 0-100 score
    risk_level VARCHAR(50), -- low, medium, high, critical
    
    -- Flags (stored as JSONB array for flexibility)
    flags JSONB DEFAULT '[]', -- [{flag: "high_value", weight: 10}, {flag: "multiple_orders_same_day", weight: 5}]
    
    -- Action
    status VARCHAR(50) NOT NULL DEFAULT 'flagged', -- flagged, reviewed, approved, rejected, investigating
    reviewed_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMP,
    review_notes TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_fraud_flags_order_id ON fraud_flags(order_id);
CREATE INDEX idx_fraud_flags_risk_level ON fraud_flags(risk_level);
CREATE INDEX idx_fraud_flags_status ON fraud_flags(status);

-- ===== ADDRESSES =====

-- Saved customer addresses
CREATE TABLE IF NOT EXISTS customer_addresses (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    label VARCHAR(50), -- Home, Work, Apartment, etc.
    recipient_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    
    -- Address components
    street_address VARCHAR(255) NOT NULL,
    apartment_suite VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state_province VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20),
    country_code VARCHAR(2) NOT NULL,
    
    -- Geolocation (optional)
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    
    -- Preferences
    is_default BOOLEAN DEFAULT FALSE,
    is_billing_address BOOLEAN DEFAULT FALSE,
    is_shipping_address BOOLEAN DEFAULT TRUE,
    
    -- Validation
    is_validated BOOLEAN DEFAULT FALSE,
    validation_source VARCHAR(50), -- google_maps, internal, etc.
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_customer_addresses_customer_id ON customer_addresses(customer_id);
CREATE INDEX idx_customer_addresses_is_default ON customer_addresses(is_default);

-- ===== DISCOUNT CODES & PROMOS =====

-- Discount and promo code campaigns
CREATE TABLE IF NOT EXISTS discount_codes (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    campaign_name VARCHAR(255),
    
    -- Discount details
    discount_type VARCHAR(50) NOT NULL, -- percentage, fixed_amount, free_shipping, buy_x_get_y
    discount_value DECIMAL(10, 2) NOT NULL,
    discount_value_cents BIGINT, -- For fixed amount discounts
    
    -- Limits
    min_order_value_cents BIGINT,
    max_discount_cents BIGINT,
    max_uses_total INT,
    max_uses_per_customer INT,
    max_uses_per_day INT,
    
    -- Applicability
    applicable_to VARCHAR(50), -- all_products, specific_categories, specific_products, specific_merchants
    applicable_items JSONB, -- Product/category IDs if restricted
    
    -- Validity
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Usage tracking
    times_used INT DEFAULT 0,
    total_discount_given_cents BIGINT DEFAULT 0,
    
    -- Metadata
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    notes TEXT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_discount_codes_code ON discount_codes(code);
CREATE INDEX idx_discount_codes_is_active ON discount_codes(is_active);
CREATE INDEX idx_discount_codes_valid_from ON discount_codes(valid_from);

-- ===== CREATOR PROMO & FLASH SALES =====

-- Creator promotional pricing and flash sales
CREATE TABLE IF NOT EXISTS creator_promos (
    id BIGSERIAL PRIMARY KEY,
    creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Promo details
    name VARCHAR(255) NOT NULL,
    description TEXT,
    promo_type VARCHAR(50) NOT NULL, -- percentage_off, flat_discount, flash_sale, bundle_deal
    
    -- Discount configuration
    discount_percentage INT,
    discount_amount_cents BIGINT,
    
    -- Applicable products
    product_ids JSONB NOT NULL DEFAULT '[]',
    category_ids JSONB DEFAULT '[]',
    
    -- Timeline
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Livestream specific
    livestream_id BIGINT REFERENCES livestreams(id) ON DELETE SET NULL,
    livestream_only BOOLEAN DEFAULT FALSE,
    
    -- Usage
    min_order_value_cents BIGINT,
    max_uses INT,
    times_used INT DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_creator_promos_creator_id ON creator_promos(creator_id);
CREATE INDEX idx_creator_promos_is_active ON creator_promos(is_active);
CREATE INDEX idx_creator_promos_start_at ON creator_promos(start_at);

-- ===== NOTIFICATIONS & EVENTS =====

-- Notification queue for order events (checkout, payment, shipping, return, support)
CREATE TABLE IF NOT EXISTS order_notifications (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    
    recipient_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_type VARCHAR(50) NOT NULL, -- customer, merchant, admin
    
    -- Notification details
    event_type VARCHAR(100) NOT NULL, -- order_created, payment_received, shipped, delivered, return_approved, etc.
    title VARCHAR(255) NOT NULL,
    message TEXT,
    
    -- Channel
    channel VARCHAR(50) NOT NULL, -- email, sms, push, in_app
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, sent, failed, read, archived
    sent_at TIMESTAMP,
    read_at TIMESTAMP,
    
    -- Retry info
    retry_count INT DEFAULT 0,
    last_retry_at TIMESTAMP,
    error_message TEXT,
    
    -- Link to related entity
    related_entity_type VARCHAR(100), -- shipment, return, payment_intent, etc.
    related_entity_id BIGINT,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_notifications_order_id ON order_notifications(order_id);
CREATE INDEX idx_order_notifications_recipient_id ON order_notifications(recipient_id);
CREATE INDEX idx_order_notifications_status ON order_notifications(status);
CREATE INDEX idx_order_notifications_created_at ON order_notifications(created_at);

-- ===== PERFORMANCE INDEXES =====

-- Composite indexes for common queries
CREATE INDEX idx_orders_merchant_status_date ON orders(merchant_id, status, created_at DESC)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_shipments_merchant_carrier_status ON shipments(merchant_id, carrier, status)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_payment_intents_order_status ON payment_intents(order_id, status)
    WHERE created_at > CURRENT_TIMESTAMP - INTERVAL '90 days';

-- Partitioning ready (can be applied later for massive scaling)
-- PARTITION BY RANGE (created_at) for orders, tracking_events, order_notifications, order_analytics

COMMIT;

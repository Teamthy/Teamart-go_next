#!/bin/bash
# Migration: 005_payment_infrastructure.sql
# Description: Create payment processing infrastructure with multi-gateway support
# Created: 2026-05-21

set -e

psql "$DATABASE_URL" <<'EOF'

-- ===== PAYMENT METHODS =====

CREATE TABLE IF NOT EXISTS payment_methods (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- card, wallet, bank_account, upi, phone_money
    provider VARCHAR(50) NOT NULL, -- stripe, paypal, razorpay, paystack, flutterwave
    
    -- Card data (tokenized, never store full card)
    card_last_four VARCHAR(4),
    card_brand VARCHAR(50), -- visa, mastercard, amex, discover
    card_expiry_month INTEGER,
    card_expiry_year INTEGER,
    cardholder_name VARCHAR(255),
    
    -- Wallet/Account data
    account_email VARCHAR(255),
    account_phone VARCHAR(20),
    account_identifier VARCHAR(255), -- Account number, UPI ID, etc.
    
    -- Gateway references
    provider_id VARCHAR(255), -- Stripe: pm_*, Razorpay: card_*, etc.
    
    -- Status
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_payment_methods_user_id ON payment_methods(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payment_methods_provider ON payment_methods(provider);
CREATE INDEX idx_payment_methods_is_default ON payment_methods(is_default);

-- ===== PAYMENT INTENTS =====

CREATE TABLE IF NOT EXISTS payment_intents (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Amount
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD', -- USD, INR, NGN, etc.
    
    -- Payment flow
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- pending → authorized → processing → succeeded
    -- Or: cancelled, failed, expired
    
    -- Gateway info
    provider VARCHAR(50) NOT NULL, -- stripe, razorpay, paystack, etc.
    provider_intent_id VARCHAR(255), -- External ID
    provider_client_secret VARCHAR(500), -- For client-side handling
    
    -- Payment method
    payment_method_id BIGINT REFERENCES payment_methods(id) ON DELETE SET NULL,
    
    -- Risk & security
    risk_level VARCHAR(50), -- low, medium, high
    requires_3d_secure BOOLEAN DEFAULT FALSE,
    three_d_secure_status VARCHAR(50),
    
    -- Metadata
    description TEXT,
    metadata JSONB,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    succeeded_at TIMESTAMP,
    failed_at TIMESTAMP
);

CREATE INDEX idx_payment_intents_order_id ON payment_intents(order_id);
CREATE INDEX idx_payment_intents_user_id ON payment_intents(user_id);
CREATE INDEX idx_payment_intents_status ON payment_intents(status);
CREATE INDEX idx_payment_intents_provider ON payment_intents(provider);

-- ===== PAYMENT TRANSACTIONS =====

CREATE TABLE IF NOT EXISTS payment_transactions (
    id BIGSERIAL PRIMARY KEY,
    payment_intent_id BIGINT NOT NULL REFERENCES payment_intents(id) ON DELETE RESTRICT,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    
    -- Transaction details
    type VARCHAR(50) NOT NULL, -- charge, refund, dispute, payout
    status VARCHAR(50) NOT NULL, -- pending, succeeded, failed
    
    -- Amount
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Gateway info
    provider VARCHAR(50) NOT NULL,
    provider_transaction_id VARCHAR(255),
    provider_reference VARCHAR(255),
    
    -- Error handling
    error_code VARCHAR(100),
    error_message TEXT,
    
    -- Metadata
    metadata JSONB,
    raw_response JSONB, -- Full gateway response for debugging
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    settled_at TIMESTAMP
);

CREATE INDEX idx_payment_transactions_payment_intent_id ON payment_transactions(payment_intent_id);
CREATE INDEX idx_payment_transactions_order_id ON payment_transactions(order_id);
CREATE INDEX idx_payment_transactions_type ON payment_transactions(type);
CREATE INDEX idx_payment_transactions_status ON payment_transactions(status);

-- ===== WEBHOOKS =====

CREATE TABLE IF NOT EXISTS webhook_logs (
    id BIGSERIAL PRIMARY KEY,
    provider VARCHAR(50) NOT NULL, -- stripe, razorpay, paystack, flutterwave
    event_type VARCHAR(100) NOT NULL,
    webhook_id VARCHAR(255),
    payload JSONB NOT NULL,
    
    -- Processing
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMP,
    
    -- Results
    success BOOLEAN,
    error_message TEXT,
    
    -- Verification
    signature_valid BOOLEAN,
    signature_algorithm VARCHAR(50),
    
    -- Timestamps
    received_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_webhook_logs_provider ON webhook_logs(provider);
CREATE INDEX idx_webhook_logs_event_type ON webhook_logs(event_type);
CREATE INDEX idx_webhook_logs_processed ON webhook_logs(processed);

-- ===== PAYOUTS =====

CREATE TABLE IF NOT EXISTS payouts (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Payout details
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- pending → approved → processing → completed
    -- Or: rejected, failed, cancelled
    
    -- Bank details
    bank_account_id BIGINT NOT NULL REFERENCES payment_methods(id) ON DELETE RESTRICT,
    
    -- Period
    payout_period_start DATE,
    payout_period_end DATE,
    
    -- Gateway info
    provider VARCHAR(50),
    provider_payout_id VARCHAR(255),
    
    -- Processing
    reviewed_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMP,
    review_notes TEXT,
    
    -- Metadata
    metadata JSONB,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    scheduled_at TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX idx_payouts_store_id ON payouts(store_id);
CREATE INDEX idx_payouts_status ON payouts(status);
CREATE INDEX idx_payouts_created_at ON payouts(created_at);

-- ===== REFUNDS =====

CREATE TABLE IF NOT EXISTS refunds (
    id BIGSERIAL PRIMARY KEY,
    payment_intent_id BIGINT NOT NULL REFERENCES payment_intents(id) ON DELETE RESTRICT,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    order_return_id BIGINT REFERENCES order_returns(id) ON DELETE SET NULL,
    
    -- Refund details
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    reason VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- pending → processing → completed
    -- Or: failed, cancelled
    
    -- Gateway info
    provider VARCHAR(50),
    provider_refund_id VARCHAR(255),
    
    -- Approval workflow
    requested_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    requested_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    approved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMP,
    
    -- Processing
    processed_at TIMESTAMP,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_refunds_payment_intent_id ON refunds(payment_intent_id);
CREATE INDEX idx_refunds_order_id ON refunds(order_id);
CREATE INDEX idx_refunds_status ON refunds(status);

-- ===== WALLET / ACCOUNT BALANCE =====

CREATE TABLE IF NOT EXISTS user_wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    
    -- Balance
    balance DECIMAL(12, 2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Holds (for pending orders)
    held_amount DECIMAL(12, 2) NOT NULL DEFAULT 0,
    available_balance GENERATED ALWAYS AS (balance - held_amount) STORED,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_wallets_user_id ON user_wallets(user_id);

-- ===== WALLET TRANSACTIONS =====

CREATE TABLE IF NOT EXISTS wallet_transactions (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES user_wallets(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Transaction
    type VARCHAR(50) NOT NULL, -- deposit, withdrawal, payment, refund, transfer
    amount DECIMAL(12, 2) NOT NULL,
    previous_balance DECIMAL(12, 2) NOT NULL,
    new_balance DECIMAL(12, 2) NOT NULL,
    
    -- Reference
    reference_type VARCHAR(100), -- order_id, refund_id, payout_id
    reference_id VARCHAR(255),
    
    -- Description
    description TEXT,
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'completed', -- pending, completed, failed, reversed
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wallet_transactions_wallet_id ON wallet_transactions(wallet_id);
CREATE INDEX idx_wallet_transactions_type ON wallet_transactions(type);
CREATE INDEX idx_wallet_transactions_created_at ON wallet_transactions(created_at);

-- ===== ESCROW =====

CREATE TABLE IF NOT EXISTS escrow_accounts (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE REFERENCES orders(id) ON DELETE RESTRICT,
    
    -- Parties
    buyer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    seller_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    
    -- Funds
    amount DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    
    -- Status
    status VARCHAR(50) NOT NULL DEFAULT 'held',
    -- held → released (to seller), refunded (to buyer), disputed
    
    -- Timeline
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    release_date DATE,
    released_at TIMESTAMP,
    refunded_at TIMESTAMP,
    
    -- Metadata
    metadata JSONB
);

CREATE INDEX idx_escrow_accounts_order_id ON escrow_accounts(order_id);
CREATE INDEX idx_escrow_accounts_status ON escrow_accounts(status);

-- ===== ESCROW DISPUTES =====

CREATE TABLE IF NOT EXISTS escrow_disputes (
    id BIGSERIAL PRIMARY KEY,
    escrow_account_id BIGINT NOT NULL REFERENCES escrow_accounts(id) ON DELETE RESTRICT,
    
    -- Initiator
    initiated_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    initiated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Dispute details
    reason TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    
    -- Resolution
    resolved_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMP,
    resolution TEXT,
    
    -- Outcome (who gets the money)
    outcome VARCHAR(50), -- buyer, seller, split
    buyer_amount DECIMAL(12, 2),
    seller_amount DECIMAL(12, 2),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_escrow_disputes_escrow_account_id ON escrow_disputes(escrow_account_id);

-- ===== TAX & CURRENCY =====

CREATE TABLE IF NOT EXISTS transaction_taxes (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    
    -- Tax info
    tax_type VARCHAR(50) NOT NULL, -- vat, gst, sales_tax, etc.
    jurisdiction VARCHAR(100), -- Country/state code
    rate DECIMAL(5, 3) NOT NULL, -- Tax rate
    
    -- Amounts
    taxable_amount DECIMAL(12, 2) NOT NULL,
    tax_amount DECIMAL(12, 2) NOT NULL,
    
    -- Metadata
    metadata JSONB,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transaction_taxes_order_id ON transaction_taxes(order_id);

-- ===== PAYMENT GATEWAY CREDENTIALS (ENCRYPTED) =====

CREATE TABLE IF NOT EXISTS payment_gateway_config (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    
    -- Gateway
    provider VARCHAR(50) NOT NULL, -- stripe, razorpay, paystack, flutterwave
    mode VARCHAR(50) NOT NULL, -- test, live
    
    -- Credentials (encrypted at rest)
    public_key VARCHAR(500),
    secret_key VARCHAR(500),
    webhook_secret VARCHAR(500),
    
    -- Status
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Settings
    settings JSONB, -- Provider-specific settings
    
    -- Audit
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_payment_gateway_config_provider ON payment_gateway_config(provider);
CREATE INDEX idx_payment_gateway_config_store_id ON payment_gateway_config(store_id);

EOF

echo "Migration completed: created payment infrastructure"

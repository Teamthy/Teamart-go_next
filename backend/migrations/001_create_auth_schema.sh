#!/bin/bash
# Migration: 001_create_auth_schema.sql
# Description: Create users and sessions tables for authentication system
# Created: 2026-05-18

set -e

psql "$DATABASE_URL" <<'EOF'

-- Users table for identity management
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    
    -- Onboarding state machine
    onboarding_state VARCHAR(50) NOT NULL DEFAULT 'new',
    account_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Login attempt tracking
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    failed_login_last_attempt TIMESTAMP,
    locked_until TIMESTAMP,
    
    -- Password management
    password_changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    
    -- Contact information
    recovery_email VARCHAR(255),
    phone_number VARCHAR(20),
    
    -- MFA configuration
    requires_mfa BOOLEAN NOT NULL DEFAULT FALSE,
    mfa_method VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes for common queries
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_account_status ON users(account_status) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_is_active ON users(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_locked_until ON users(locked_until) WHERE locked_until > CURRENT_TIMESTAMP;
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Sessions table for session management
CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255) NOT NULL,
    device_fingerprint VARCHAR(255) NOT NULL,
    user_agent TEXT,
    ip_address VARCHAR(45) NOT NULL,
    
    -- Trust and security
    trust_level VARCHAR(50) NOT NULL DEFAULT 'untrusted',
    requires_mfa_step BOOLEAN NOT NULL DEFAULT FALSE,
    requires_password_verification BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Geolocation tracking
    geo_country VARCHAR(2),
    geo_city VARCHAR(100),
    geo_latitude NUMERIC(10, 8),
    geo_longitude NUMERIC(11, 8),
    geo_timezone VARCHAR(50),
    
    -- Verification tracking
    mfa_verified_at TIMESTAMP,
    password_verified_at TIMESTAMP,
    
    -- Revocation tracking
    revoked_at TIMESTAMP,
    revoke_reason VARCHAR(255),
    
    -- Activity tracking
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_activity_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

-- Indexes for common queries
CREATE INDEX idx_sessions_user_id ON sessions(user_id) WHERE revoked_at IS NULL;
CREATE INDEX idx_sessions_device_id ON sessions(device_id, user_id);
CREATE INDEX idx_sessions_ip_address ON sessions(ip_address);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at) WHERE revoked_at IS NULL;
CREATE INDEX idx_sessions_created_at ON sessions(created_at);
CREATE INDEX idx_sessions_last_activity ON sessions(last_activity_at) WHERE revoked_at IS NULL;
CREATE INDEX idx_sessions_revoked_at ON sessions(revoked_at) WHERE revoked_at IS NOT NULL;

-- Add trigger to update users.updated_at automatically
CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_updated_at_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_users_updated_at();

EOF

echo "Migration completed: created users and sessions tables"

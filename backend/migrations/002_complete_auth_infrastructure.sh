#!/bin/bash
# Migration: 002_complete_auth_infrastructure.sql
# Description: Create complete auth infrastructure tables
# Created: 2026-05-20

set -e

psql "$DATABASE_URL" <<'EOF'

-- ===== OTP (One-Time Password) Tables =====

CREATE TABLE IF NOT EXISTS otp_codes (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    code_hash VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 5,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_otp_codes_user_id ON otp_codes(user_id, type);
CREATE INDEX idx_otp_codes_destination ON otp_codes(destination);
CREATE INDEX idx_otp_codes_expires_at ON otp_codes(expires_at);

-- ===== Role-Based Access Control (RBAC) Tables =====

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_is_active ON roles(is_active);

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(50),
    resource VARCHAR(50),
    action VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_action ON permissions(action);
CREATE INDEX idx_permissions_category ON permissions(category);

-- Role to permissions mapping (many-to-many)
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id VARCHAR(100) NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (role_id, permission_id)
);

-- User to roles mapping (many-to-many)
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    granted_by BIGINT NOT NULL REFERENCES users(id),
    granted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_user_roles_expires_at ON user_roles(expires_at);

-- ===== Device Trust Tables =====

CREATE TABLE IF NOT EXISTS device_trusts (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255) NOT NULL,
    device_fingerprint VARCHAR(255) NOT NULL,
    device_name VARCHAR(255),
    device_type VARCHAR(50),
    ip_address VARCHAR(45),
    user_agent TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMP,
    last_seen_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_device_trusts_user_id ON device_trusts(user_id);
CREATE INDEX idx_device_trusts_device_id ON device_trusts(device_id, user_id);
CREATE INDEX idx_device_trusts_is_verified ON device_trusts(is_verified);
CREATE INDEX idx_device_trusts_expires_at ON device_trusts(expires_at);

-- ===== Password Reset & Account Recovery Tables =====

CREATE TABLE IF NOT EXISTS password_resets (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT FALSE,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_password_resets_user_id ON password_resets(user_id);
CREATE INDEX idx_password_resets_email ON password_resets(email);
CREATE INDEX idx_password_resets_expires_at ON password_resets(expires_at);
CREATE INDEX idx_password_resets_is_used ON password_resets(is_used);

CREATE TABLE IF NOT EXISTS account_recoveries (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recovery_code_hash VARCHAR(255) NOT NULL,
    backup_codes TEXT,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_account_recoveries_user_id ON account_recoveries(user_id);
CREATE INDEX idx_account_recoveries_expires_at ON account_recoveries(expires_at);

-- ===== Audit & Security Event Tables =====

CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    description TEXT,
    metadata JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_event_type ON audit_logs(event_type);
CREATE INDEX idx_audit_logs_severity ON audit_logs(severity);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

CREATE TABLE IF NOT EXISTS security_events (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    description TEXT,
    metadata JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    resolved BOOLEAN NOT NULL DEFAULT FALSE,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_security_events_user_id ON security_events(user_id);
CREATE INDEX idx_security_events_event_type ON security_events(event_type);
CREATE INDEX idx_security_events_severity ON security_events(severity);
CREATE INDEX idx_security_events_resolved ON security_events(resolved);
CREATE INDEX idx_security_events_created_at ON security_events(created_at);

CREATE TABLE IF NOT EXISTS login_attempts (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    email VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    success BOOLEAN NOT NULL,
    failure_reason VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_login_attempts_user_id ON login_attempts(user_id);
CREATE INDEX idx_login_attempts_email ON login_attempts(email);
CREATE INDEX idx_login_attempts_ip_address ON login_attempts(ip_address);
CREATE INDEX idx_login_attempts_success ON login_attempts(success);
CREATE INDEX idx_login_attempts_created_at ON login_attempts(created_at);

-- ===== KYC & Compliance Tables =====

CREATE TABLE IF NOT EXISTS kyc_submissions (
    id VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    submission_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    document_paths TEXT,
    submitted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMP,
    reviewed_by BIGINT REFERENCES users(id),
    rejection_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_kyc_submissions_user_id ON kyc_submissions(user_id);
CREATE INDEX idx_kyc_submissions_status ON kyc_submissions(status);
CREATE INDEX idx_kyc_submissions_submitted_at ON kyc_submissions(submitted_at);

-- ===== Onboarding State Tracking =====

CREATE TABLE IF NOT EXISTS onboarding_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    current_step VARCHAR(100) NOT NULL,
    step_data JSONB,
    completed_steps TEXT,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_onboarding_progress_user_id ON onboarding_progress(user_id);
CREATE INDEX idx_onboarding_progress_current_step ON onboarding_progress(current_step);

-- ===== Create Default Roles =====

INSERT INTO roles (id, name, description) VALUES 
    (1, 'admin', 'System administrator with full access'),
    (2, 'merchant', 'Merchant user'),
    (3, 'creator', 'Creator/influencer user'),
    (4, 'customer', 'Regular customer user'),
    (5, 'support', 'Support agent'),
    (6, 'moderator', 'Stream moderator')
ON CONFLICT DO NOTHING;

-- ===== Create Default Permissions =====

INSERT INTO permissions (id, name, description, resource, action) VALUES 
    ('users:read', 'Read users', 'Read user information', 'users', 'read'),
    ('users:create', 'Create users', 'Create new users', 'users', 'create'),
    ('users:update', 'Update users', 'Update user information', 'users', 'update'),
    ('users:delete', 'Delete users', 'Delete users', 'users', 'delete'),
    ('products:read', 'Read products', 'Read product information', 'products', 'read'),
    ('products:create', 'Create products', 'Create new products', 'products', 'create'),
    ('products:update', 'Update products', 'Update product information', 'products', 'update'),
    ('products:delete', 'Delete products', 'Delete products', 'products', 'delete'),
    ('streams:start', 'Start streams', 'Start live streams', 'streams', 'start'),
    ('streams:moderate', 'Moderate streams', 'Moderate live streams', 'streams', 'moderate'),
    ('orders:read', 'Read orders', 'Read order information', 'orders', 'read'),
    ('orders:create', 'Create orders', 'Create new orders', 'orders', 'create'),
    ('orders:refund', 'Refund orders', 'Issue order refunds', 'orders', 'refund'),
    ('admin:moderate', 'Admin moderation', 'Moderate content', 'admin', 'moderate')
ON CONFLICT DO NOTHING;

-- ===== Assign Permissions to Roles =====

-- Admin has all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions
ON CONFLICT DO NOTHING;

-- Customer permissions
INSERT INTO role_permissions (role_id, permission_id) VALUES 
    (4, 'users:read'),
    (4, 'products:read'),
    (4, 'orders:read'),
    (4, 'orders:create')
ON CONFLICT DO NOTHING;

-- Merchant permissions
INSERT INTO role_permissions (role_id, permission_id) VALUES 
    (2, 'users:read'),
    (2, 'products:create'),
    (2, 'products:update'),
    (2, 'products:read'),
    (2, 'streams:start'),
    (2, 'orders:read'),
    (2, 'orders:refund')
ON CONFLICT DO NOTHING;

-- Creator permissions
INSERT INTO role_permissions (role_id, permission_id) VALUES 
    (3, 'users:read'),
    (3, 'streams:start'),
    (3, 'products:read'),
    (3, 'orders:read')
ON CONFLICT DO NOTHING;

-- Support permissions
INSERT INTO role_permissions (role_id, permission_id) VALUES 
    (5, 'users:read'),
    (5, 'orders:read'),
    (5, 'admin:moderate')
ON CONFLICT DO NOTHING;

-- Moderator permissions
INSERT INTO role_permissions (role_id, permission_id) VALUES 
    (6, 'streams:moderate'),
    (6, 'admin:moderate')
ON CONFLICT DO NOTHING;

EOF

echo "Migration completed: created complete auth infrastructure"

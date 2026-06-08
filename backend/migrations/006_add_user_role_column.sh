#!/bin/bash
# Migration: 006_add_user_role_column.sql
# Description: Add role column to users table for existing databases
# Created: 2026-06-04

set -e

psql "$DATABASE_URL" <<'EOF'

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'customer';

EOF

echo "Migration completed: added role column to users table"

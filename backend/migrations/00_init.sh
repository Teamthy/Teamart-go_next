#!/bin/bash

# Database initialization script
# This script sets up the initial database schema

set -e

echo "Initializing Teamart database..."

# Create extensions
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Enable required extensions
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "vector";
    
    GRANT ALL PRIVILEGES ON DATABASE "$POSTGRES_DB" TO "$POSTGRES_USER";
    
    GRANT USAGE ON SCHEMA public TO "$POSTGRES_USER";
    GRANT CREATE ON SCHEMA public TO "$POSTGRES_USER";
    
    -- Create initial indexes
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_created_at ON users(created_at);
    
    COMMENT ON EXTENSION vector IS 'pgvector extension for semantic search';
    COMMENT ON EXTENSION "uuid-ossp" IS 'UUID generation extension';
    
EOSQL

echo "Database initialization complete!"

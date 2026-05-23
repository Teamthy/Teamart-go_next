#!/bin/bash
# Migration: 003_product_catalog_system.sql
# Description: Create comprehensive product catalog system with variants, inventory, media, categories
# Created: 2026-05-21

set -e

psql "$DATABASE_URL" <<'EOF'

-- ===== CATEGORIES =====

CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    parent_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    image_url VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_is_active ON categories(is_active);
CREATE INDEX idx_categories_display_order ON categories(display_order);

-- ===== PRODUCT ATTRIBUTES (SCHEMA) =====

CREATE TABLE IF NOT EXISTS product_attributes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL, -- text, select, number, color, size
    is_filterable BOOLEAN NOT NULL DEFAULT TRUE,
    is_searchable BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_attributes_name ON product_attributes(name);
CREATE INDEX idx_product_attributes_type ON product_attributes(type);

-- ===== ATTRIBUTE VALUES (FOR SELECT/ENUM ATTRIBUTES) =====

CREATE TABLE IF NOT EXISTS attribute_values (
    id BIGSERIAL PRIMARY KEY,
    attribute_id BIGINT NOT NULL REFERENCES product_attributes(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(attribute_id, value)
);

CREATE INDEX idx_attribute_values_attribute_id ON attribute_values(attribute_id);

-- ===== PRODUCTS =====

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Store owner (user_id)
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    sku VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    
    -- Pricing (base price - variants can override)
    price DECIMAL(12, 2) NOT NULL,
    compare_at_price DECIMAL(12, 2),
    cost DECIMAL(12, 2),
    
    -- Product type
    product_type VARCHAR(50) NOT NULL DEFAULT 'physical', -- physical, digital, service
    
    -- Status & visibility
    status VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, active, archived
    is_visible BOOLEAN NOT NULL DEFAULT FALSE,
    featured BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Inventory management
    track_inventory BOOLEAN NOT NULL DEFAULT TRUE,
    requires_shipping BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Media
    featured_image_url VARCHAR(500),
    
    -- SEO
    meta_title VARCHAR(255),
    meta_description VARCHAR(500),
    meta_keywords VARCHAR(500),
    
    -- Search vector for full-text search
    search_vector tsvector,
    
    -- Ratings & reviews
    rating DECIMAL(3, 2) DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    
    -- Analytics
    view_count INTEGER DEFAULT 0,
    sale_count INTEGER DEFAULT 0,
    
    -- Livestream integration
    livestream_pinnable BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_products_store_id ON products(store_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_sku ON products(sku) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_slug ON products(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_status ON products(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_is_visible ON products(is_visible) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_created_at ON products(created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_search_vector ON products USING gin(search_vector);

-- ===== PRODUCT VARIANTS =====

CREATE TABLE IF NOT EXISTS product_variants (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL, -- e.g., "Red / Large"
    
    -- Pricing (override product pricing if set)
    price DECIMAL(12, 2),
    compare_at_price DECIMAL(12, 2),
    cost DECIMAL(12, 2),
    
    -- Media
    image_url VARCHAR(500),
    
    -- Status
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Attributes for this variant (JSON for flexibility)
    -- e.g., {"color": "red", "size": "large"}
    attributes JSONB,
    
    -- Inventory tracking
    quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    available_quantity INTEGER GENERATED ALWAYS AS (quantity - reserved_quantity) STORED,
    
    -- Low stock alert threshold
    low_stock_threshold INTEGER,
    
    -- Barcode for physical products
    barcode VARCHAR(255),
    
    -- Dimensions and weight for shipping
    weight_kg DECIMAL(8, 3),
    weight_unit VARCHAR(10) DEFAULT 'kg', -- kg, lbs
    dimensions_length DECIMAL(8, 2),
    dimensions_width DECIMAL(8, 2),
    dimensions_height DECIMAL(8, 2),
    dimensions_unit VARCHAR(10) DEFAULT 'cm', -- cm, in
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_variants_sku ON product_variants(sku) WHERE deleted_at IS NULL;
CREATE INDEX idx_product_variants_is_available ON product_variants(is_available);

-- ===== PRODUCT VARIANT ATTRIBUTES (EXPLICIT MAPPING) =====

CREATE TABLE IF NOT EXISTS product_variant_attributes (
    id BIGSERIAL PRIMARY KEY,
    variant_id BIGINT NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    attribute_id BIGINT NOT NULL REFERENCES product_attributes(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(variant_id, attribute_id)
);

CREATE INDEX idx_product_variant_attributes_variant_id ON product_variant_attributes(variant_id);

-- ===== PRODUCT INVENTORY =====

CREATE TABLE IF NOT EXISTS product_inventory (
    id BIGSERIAL PRIMARY KEY,
    variant_id BIGINT NOT NULL UNIQUE REFERENCES product_variants(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    warehouse_id VARCHAR(100), -- Optional: for multi-warehouse tracking
    
    -- History tracking
    last_recount_at TIMESTAMP,
    last_recount_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_inventory_variant_id ON product_inventory(variant_id);
CREATE INDEX idx_product_inventory_warehouse_id ON product_inventory(warehouse_id);

-- ===== INVENTORY TRANSACTIONS (AUDIT TRAIL) =====

CREATE TABLE IF NOT EXISTS inventory_transactions (
    id BIGSERIAL PRIMARY KEY,
    variant_id BIGINT NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- purchase, sale, adjustment, return, damage
    quantity_change INTEGER NOT NULL,
    reference_type VARCHAR(100), -- order_id, return_id, adjustment_id
    reference_id VARCHAR(100),
    reason TEXT,
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_inventory_transactions_variant_id ON inventory_transactions(variant_id);
CREATE INDEX idx_inventory_transactions_type ON inventory_transactions(type);
CREATE INDEX idx_inventory_transactions_created_at ON inventory_transactions(created_at);

-- ===== PRODUCT MEDIA =====

CREATE TABLE IF NOT EXISTS product_media (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES product_variants(id) ON DELETE CASCADE,
    media_type VARCHAR(50) NOT NULL, -- image, video, document
    url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    display_order INTEGER NOT NULL DEFAULT 0,
    size_kb INTEGER,
    width INTEGER,
    height INTEGER,
    metadata JSONB, -- For storing additional info like video duration, etc.
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_media_product_id ON product_media(product_id);
CREATE INDEX idx_product_media_variant_id ON product_media(variant_id);
CREATE INDEX idx_product_media_media_type ON product_media(media_type);

-- ===== PRODUCT ATTRIBUTES (ACTUAL VALUES ON PRODUCTS) =====

CREATE TABLE IF NOT EXISTS product_attribute_values (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    attribute_id BIGINT NOT NULL REFERENCES product_attributes(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, attribute_id)
);

CREATE INDEX idx_product_attribute_values_product_id ON product_attribute_values(product_id);

-- ===== PRODUCT PRICING RULES =====

CREATE TABLE IF NOT EXISTS product_pricing_rules (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    rule_type VARCHAR(50) NOT NULL, -- volume_discount, buyer_segment, time_limited
    
    -- Conditions
    min_quantity INTEGER,
    max_quantity INTEGER,
    customer_segment VARCHAR(100), -- wholesale, vip, new_customer
    
    -- Discount
    discount_type VARCHAR(50), -- percentage, fixed
    discount_value DECIMAL(12, 2) NOT NULL,
    
    -- Validity
    valid_from TIMESTAMP,
    valid_until TIMESTAMP,
    
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_pricing_rules_product_id ON product_pricing_rules(product_id);
CREATE INDEX idx_product_pricing_rules_is_active ON product_pricing_rules(is_active);

-- ===== PRODUCT BUNDLES =====

CREATE TABLE IF NOT EXISTS product_bundles (
    id BIGSERIAL PRIMARY KEY,
    store_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Pricing
    price DECIMAL(12, 2) NOT NULL,
    original_price DECIMAL(12, 2),
    discount_percentage DECIMAL(5, 2),
    
    -- Status
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_bundles_store_id ON product_bundles(store_id);
CREATE INDEX idx_product_bundles_is_active ON product_bundles(is_active);

-- ===== BUNDLE ITEMS =====

CREATE TABLE IF NOT EXISTS bundle_items (
    id BIGSERIAL PRIMARY KEY,
    bundle_id BIGINT NOT NULL REFERENCES product_bundles(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES product_variants(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bundle_items_bundle_id ON bundle_items(bundle_id);
CREATE INDEX idx_bundle_items_product_id ON bundle_items(product_id);

-- ===== PRODUCT REVIEWS & RATINGS =====

CREATE TABLE IF NOT EXISTS product_reviews (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    order_id BIGINT, -- Optional: link to order for verification
    reviewer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    comment TEXT,
    
    -- Media
    media_urls TEXT[], -- JSON array of image URLs
    
    -- Moderation
    is_verified_purchase BOOLEAN NOT NULL DEFAULT FALSE,
    is_approved BOOLEAN NOT NULL DEFAULT FALSE,
    is_flagged BOOLEAN NOT NULL DEFAULT FALSE,
    flag_reason VARCHAR(255),
    
    -- Helpful votes
    helpful_count INTEGER DEFAULT 0,
    unhelpful_count INTEGER DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_reviews_product_id ON product_reviews(product_id);
CREATE INDEX idx_product_reviews_reviewer_id ON product_reviews(reviewer_id);
CREATE INDEX idx_product_reviews_is_approved ON product_reviews(is_approved);
CREATE INDEX idx_product_reviews_created_at ON product_reviews(created_at);

-- ===== PRODUCT SEARCH KEYWORDS =====

CREATE TABLE IF NOT EXISTS product_search_keywords (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    keyword VARCHAR(255) NOT NULL,
    search_count INTEGER DEFAULT 0,
    last_searched_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, keyword)
);

CREATE INDEX idx_product_search_keywords_product_id ON product_search_keywords(product_id);
CREATE INDEX idx_product_search_keywords_keyword ON product_search_keywords(keyword);

-- ===== INSERT DEFAULT ATTRIBUTES =====

INSERT INTO product_attributes (name, type, is_filterable, is_searchable) VALUES
    ('Color', 'color', true, true),
    ('Size', 'select', true, true),
    ('Material', 'select', true, true),
    ('Brand', 'select', true, true),
    ('Gender', 'select', true, true),
    ('Age Group', 'select', true, false),
    ('Style', 'select', true, false),
    ('Pattern', 'select', true, false)
ON CONFLICT DO NOTHING;

-- ===== INSERT DEFAULT ATTRIBUTE VALUES =====

-- Sizes
INSERT INTO attribute_values (attribute_id, value, display_order) 
SELECT id, value, ord FROM (
    VALUES 
        ('Size', 'XS', 1),
        ('Size', 'S', 2),
        ('Size', 'M', 3),
        ('Size', 'L', 4),
        ('Size', 'XL', 5),
        ('Size', 'XXL', 6)
) AS t(attr_name, value, ord)
JOIN product_attributes pa ON pa.name = t.attr_name
ON CONFLICT DO NOTHING;

-- Colors
INSERT INTO attribute_values (attribute_id, value, display_order)
SELECT id, value, ord FROM (
    VALUES
        ('Color', 'Black', 1),
        ('Color', 'White', 2),
        ('Color', 'Red', 3),
        ('Color', 'Blue', 4),
        ('Color', 'Green', 5),
        ('Color', 'Yellow', 6),
        ('Color', 'Purple', 7),
        ('Color', 'Pink', 8),
        ('Color', 'Orange', 9),
        ('Color', 'Brown', 10),
        ('Color', 'Gray', 11),
        ('Color', 'Navy', 12)
) AS t(attr_name, value, ord)
JOIN product_attributes pa ON pa.name = t.attr_name
ON CONFLICT DO NOTHING;

-- Materials
INSERT INTO attribute_values (attribute_id, value, display_order)
SELECT id, value, ord FROM (
    VALUES
        ('Material', 'Cotton', 1),
        ('Material', 'Polyester', 2),
        ('Material', 'Wool', 3),
        ('Material', 'Silk', 4),
        ('Material', 'Leather', 5),
        ('Material', 'Denim', 6),
        ('Material', 'Linen', 7)
) AS t(attr_name, value, ord)
JOIN product_attributes pa ON pa.name = t.attr_name
ON CONFLICT DO NOTHING;

EOF

echo "Migration completed: created product catalog system"

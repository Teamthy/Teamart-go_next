# Products Domain

Product catalog and management system.

## Features
- Product creation and management
- Product variants
- SKU management
- Category management
- AI-generated descriptions and tags
- Image uploads
- Semantic search vectors

## Entities
- Product
- ProductVariant
- Category
- ProductImage

## Services
- ProductService - Product CRUD operations
- CategoryService - Category management
- AIProductService - AI description/tag generation

## API
- POST /products
- GET /products
- GET /products/:id
- PATCH /products/:id
- DELETE /products/:id
- POST /products/:id/images

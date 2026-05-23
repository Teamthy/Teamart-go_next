package products

import (
	"time"
)

// ===== PRODUCT TYPES =====

// Product represents a product in the catalog
type Product struct {
	ID                 int64
	StoreID            int64
	CategoryID         int64
	SKU                string
	Name               string
	Slug               string
	Description        string
	Price              float64
	CompareAtPrice     *float64
	Cost               *float64
	ProductType        string // physical, digital, service
	Status             string // draft, active, archived
	IsVisible          bool
	Featured           bool
	TrackInventory     bool
	RequiresShipping   bool
	FeaturedImageURL   *string
	MetaTitle          *string
	MetaDescription    *string
	MetaKeywords       *string
	Rating             float64
	ReviewCount        int64
	ViewCount          int64
	SaleCount          int64
	LivestreamPinnable bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	PublishedAt        *time.Time
	DeletedAt          *time.Time
}

// ProductVariant represents a variant of a product (size, color, etc.)
type ProductVariant struct {
	ID                int64
	ProductID         int64
	SKU               string
	Title             string
	Price             *float64
	CompareAtPrice    *float64
	Cost              *float64
	ImageURL          *string
	IsAvailable       bool
	Attributes        map[string]string // e.g., {"color": "red", "size": "large"}
	Quantity          int64
	ReservedQuantity  int64
	AvailableQuantity int64 // Generated: quantity - reserved_quantity
	LowStockThreshold *int64
	Barcode           *string
	WeightKg          *float64
	WeightUnit        string // kg, lbs
	DimensionsLength  *float64
	DimensionsWidth   *float64
	DimensionsHeight  *float64
	DimensionsUnit    string // cm, in
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

// ProductMedia represents media (images, videos) for a product
type ProductMedia struct {
	ID           int64
	ProductID    int64
	VariantID    *int64
	MediaType    string // image, video, document
	URL          string
	AltText      *string
	DisplayOrder int64
	SizeKB       *int64
	Width        *int64
	Height       *int64
	Metadata     map[string]interface{} // Additional info
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ProductReview represents a customer review
type ProductReview struct {
	ID                 int64
	ProductID          int64
	OrderID            *int64
	ReviewerID         int64
	Rating             int64 // 1-5
	Title              *string
	Comment            *string
	MediaURLs          []string
	IsVerifiedPurchase bool
	IsApproved         bool
	IsFlagged          bool
	FlagReason         *string
	HelpfulCount       int64
	UnhelpfulCount     int64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// Category represents a product category
type Category struct {
	ID           int64
	Name         string
	Slug         string
	Description  *string
	ParentID     *int64
	ImageURL     *string
	IsActive     bool
	DisplayOrder int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ProductAttribute represents an attribute (e.g., color, size)
type ProductAttribute struct {
	ID           int64
	Name         string
	Type         string // text, select, number, color, size
	IsFilterable bool
	IsSearchable bool
	CreatedAt    time.Time
}

// AttributeValue represents a specific value for an attribute (e.g., "Red" for color)
type AttributeValue struct {
	ID           int64
	AttributeID  int64
	Value        string
	DisplayOrder int64
	CreatedAt    time.Time
}

// ProductBundle represents a bundle of products sold together
type ProductBundle struct {
	ID                 int64
	StoreID            int64
	Name               string
	Description        *string
	Price              float64
	OriginalPrice      *float64
	DiscountPercentage *float64
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// BundleItem represents an item in a bundle
type BundleItem struct {
	ID        int64
	BundleID  int64
	ProductID int64
	VariantID *int64
	Quantity  int64
	Order     int64
	CreatedAt time.Time
}

// InventoryTransaction tracks inventory changes
type InventoryTransaction struct {
	ID             int64
	VariantID      int64
	Type           string // purchase, sale, adjustment, return, damage
	QuantityChange int64
	ReferenceType  *string
	ReferenceID    *string
	Reason         *string
	CreatedBy      int64
	CreatedAt      time.Time
}

// ===== INPUT/OUTPUT TYPES =====

// CreateProductInput represents the input for creating a product
type CreateProductInput struct {
	CategoryID         int64
	SKU                string
	Name               string
	Slug               string
	Description        string
	Price              float64
	CompareAtPrice     *float64
	Cost               *float64
	ProductType        string
	TrackInventory     bool
	RequiresShipping   bool
	FeaturedImageURL   *string
	MetaTitle          *string
	MetaDescription    *string
	MetaKeywords       *string
	LivestreamPinnable bool
}

// UpdateProductInput represents the input for updating a product
type UpdateProductInput struct {
	Name             string
	Slug             string
	Description      string
	Price            float64
	CompareAtPrice   *float64
	Cost             *float64
	Status           string
	IsVisible        bool
	Featured         bool
	FeaturedImageURL *string
	MetaTitle        *string
	MetaDescription  *string
	MetaKeywords     *string
}

// CreateVariantInput represents the input for creating a product variant
type CreateVariantInput struct {
	ProductID         int64
	SKU               string
	Title             string
	Price             *float64
	CompareAtPrice    *float64
	Cost              *float64
	ImageURL          *string
	Attributes        map[string]string
	Quantity          int64
	LowStockThreshold *int64
	Barcode           *string
	WeightKg          *float64
	DimensionsLength  *float64
	DimensionsWidth   *float64
	DimensionsHeight  *float64
}

// UpdateVariantInput represents the input for updating a product variant
type UpdateVariantInput struct {
	Title             string
	Price             *float64
	CompareAtPrice    *float64
	Cost              *float64
	ImageURL          *string
	IsAvailable       bool
	Attributes        map[string]string
	Quantity          int64
	LowStockThreshold *int64
	WeightKg          *float64
	DimensionsLength  *float64
	DimensionsWidth   *float64
	DimensionsHeight  *float64
}

// CreateMediaInput represents the input for creating product media
type CreateMediaInput struct {
	ProductID    int64
	VariantID    *int64
	MediaType    string // image, video, document
	URL          string
	AltText      *string
	DisplayOrder int64
	SizeKB       *int64
	Width        *int64
	Height       *int64
}

// CreateReviewInput represents the input for creating a product review
type CreateReviewInput struct {
	ProductID          int64
	OrderID            *int64
	ReviewerID         int64
	Rating             int64 // 1-5
	Title              *string
	Comment            *string
	MediaURLs          []string
	IsVerifiedPurchase bool
}

// CreateCategoryInput represents the input for creating a category
type CreateCategoryInput struct {
	Name         string
	Slug         string
	Description  *string
	ParentID     *int64
	ImageURL     *string
	IsActive     bool
	DisplayOrder int64
}

// UpdateCategoryInput represents the input for updating a category
type UpdateCategoryInput struct {
	Name         string
	Slug         string
	Description  *string
	ImageURL     *string
	IsActive     bool
	DisplayOrder int64
}

// ===== SEARCH & FILTER TYPES =====

// ProductSearchCriteria represents search criteria for products
type ProductSearchCriteria struct {
	Query       string
	CategoryID  *int64
	MinPrice    *float64
	MaxPrice    *float64
	Rating      *float64
	Status      *string
	ProductType *string
	SortBy      string // rating, sales, newest, price_asc, price_desc
	Limit       int64
	Offset      int64
}

// ListProductsOptions represents options for listing products
type ListProductsOptions struct {
	Limit  int64
	Offset int64
	SortBy string // created_at, updated_at, rating, sales
}

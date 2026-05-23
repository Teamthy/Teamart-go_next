package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/teamart/commerce-api/internal/products"
	"github.com/teamart/commerce-api/pkg/logger"
)

// ProductHandler handles HTTP requests related to products
type ProductHandler struct {
	service *products.Service
	logger  *logger.Logger
}

// NewProductHandler creates a new product HTTP handler
func NewProductHandler(service *products.Service, logger *logger.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

// ===== REQUEST TYPES =====

// CreateProductRequest represents the HTTP request to create a product
type CreateProductRequest struct {
	CategoryID         int64    `json:"category_id" binding:"required"`
	SKU                string   `json:"sku" binding:"required"`
	Name               string   `json:"name" binding:"required"`
	Slug               string   `json:"slug"`
	Description        string   `json:"description"`
	Price              float64  `json:"price" binding:"required"`
	CompareAtPrice     *float64 `json:"compare_at_price"`
	Cost               *float64 `json:"cost"`
	ProductType        string   `json:"product_type"`
	TrackInventory     bool     `json:"track_inventory"`
	RequiresShipping   bool     `json:"requires_shipping"`
	FeaturedImageURL   *string  `json:"featured_image_url"`
	MetaTitle          *string  `json:"meta_title"`
	MetaDescription    *string  `json:"meta_description"`
	MetaKeywords       *string  `json:"meta_keywords"`
	LivestreamPinnable bool     `json:"livestream_pinnable"`
}

// UpdateProductRequest represents the HTTP request to update a product
type UpdateProductRequest struct {
	Name             string   `json:"name"`
	Slug             string   `json:"slug"`
	Description      string   `json:"description"`
	Price            float64  `json:"price"`
	CompareAtPrice   *float64 `json:"compare_at_price"`
	Cost             *float64 `json:"cost"`
	Status           string   `json:"status"`
	IsVisible        bool     `json:"is_visible"`
	Featured         bool     `json:"featured"`
	FeaturedImageURL *string  `json:"featured_image_url"`
	MetaTitle        *string  `json:"meta_title"`
	MetaDescription  *string  `json:"meta_description"`
	MetaKeywords     *string  `json:"meta_keywords"`
}

// CreateVariantRequest represents the HTTP request to create a variant
type CreateVariantRequest struct {
	SKU               string            `json:"sku" binding:"required"`
	Title             string            `json:"title" binding:"required"`
	Price             *float64          `json:"price"`
	CompareAtPrice    *float64          `json:"compare_at_price"`
	Cost              *float64          `json:"cost"`
	ImageURL          *string           `json:"image_url"`
	Attributes        map[string]string `json:"attributes"`
	Quantity          int64             `json:"quantity"`
	LowStockThreshold *int64            `json:"low_stock_threshold"`
	Barcode           *string           `json:"barcode"`
	WeightKg          *float64          `json:"weight_kg"`
	DimensionsLength  *float64          `json:"dimensions_length"`
	DimensionsWidth   *float64          `json:"dimensions_width"`
	DimensionsHeight  *float64          `json:"dimensions_height"`
}

// CreateMediaRequest represents the HTTP request to create product media
type CreateMediaRequest struct {
	MediaType    string  `json:"media_type" binding:"required"`
	URL          string  `json:"url" binding:"required"`
	AltText      *string `json:"alt_text"`
	DisplayOrder int64   `json:"display_order"`
	Width        *int64  `json:"width"`
	Height       *int64  `json:"height"`
}

// CreateReviewRequest represents the HTTP request to create a product review
type CreateReviewRequest struct {
	Rating             int64    `json:"rating" binding:"required,min=1,max=5"`
	Title              *string  `json:"title"`
	Comment            *string  `json:"comment"`
	MediaURLs          []string `json:"media_urls"`
	IsVerifiedPurchase bool     `json:"is_verified_purchase"`
}

// ===== RESPONSE TYPES =====

// ProductResponse represents the HTTP response for a product
type ProductResponse struct {
	ID                 int64                    `json:"id"`
	StoreID            int64                    `json:"store_id"`
	CategoryID         int64                    `json:"category_id"`
	SKU                string                   `json:"sku"`
	Name               string                   `json:"name"`
	Slug               string                   `json:"slug"`
	Description        string                   `json:"description"`
	Price              float64                  `json:"price"`
	CompareAtPrice     *float64                 `json:"compare_at_price"`
	Cost               *float64                 `json:"cost"`
	ProductType        string                   `json:"product_type"`
	Status             string                   `json:"status"`
	IsVisible          bool                     `json:"is_visible"`
	Featured           bool                     `json:"featured"`
	TrackInventory     bool                     `json:"track_inventory"`
	RequiresShipping   bool                     `json:"requires_shipping"`
	FeaturedImageURL   *string                  `json:"featured_image_url"`
	MetaTitle          *string                  `json:"meta_title"`
	MetaDescription    *string                  `json:"meta_description"`
	MetaKeywords       *string                  `json:"meta_keywords"`
	Rating             float64                  `json:"rating"`
	ReviewCount        int64                    `json:"review_count"`
	ViewCount          int64                    `json:"view_count"`
	SaleCount          int64                    `json:"sale_count"`
	LivestreamPinnable bool                     `json:"livestream_pinnable"`
	CreatedAt          string                   `json:"created_at"`
	UpdatedAt          string                   `json:"updated_at"`
	PublishedAt        *string                  `json:"published_at"`
	Variants           []ProductVariantResponse `json:"variants,omitempty"`
	Media              []ProductMediaResponse   `json:"media,omitempty"`
}

// ProductVariantResponse represents the HTTP response for a product variant
type ProductVariantResponse struct {
	ID                int64             `json:"id"`
	ProductID         int64             `json:"product_id"`
	SKU               string            `json:"sku"`
	Title             string            `json:"title"`
	Price             *float64          `json:"price"`
	CompareAtPrice    *float64          `json:"compare_at_price"`
	Cost              *float64          `json:"cost"`
	ImageURL          *string           `json:"image_url"`
	IsAvailable       bool              `json:"is_available"`
	Attributes        map[string]string `json:"attributes"`
	Quantity          int64             `json:"quantity"`
	ReservedQuantity  int64             `json:"reserved_quantity"`
	AvailableQuantity int64             `json:"available_quantity"`
	LowStockThreshold *int64            `json:"low_stock_threshold"`
	Barcode           *string           `json:"barcode"`
	WeightKg          *float64          `json:"weight_kg"`
	DimensionsLength  *float64          `json:"dimensions_length"`
	DimensionsWidth   *float64          `json:"dimensions_width"`
	DimensionsHeight  *float64          `json:"dimensions_height"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
}

// ProductMediaResponse represents the HTTP response for product media
type ProductMediaResponse struct {
	ID           int64   `json:"id"`
	ProductID    int64   `json:"product_id"`
	VariantID    *int64  `json:"variant_id"`
	MediaType    string  `json:"media_type"`
	URL          string  `json:"url"`
	AltText      *string `json:"alt_text"`
	DisplayOrder int64   `json:"display_order"`
	SizeKB       *int64  `json:"size_kb"`
	Width        *int64  `json:"width"`
	Height       *int64  `json:"height"`
	CreatedAt    string  `json:"created_at"`
}

// ProductReviewResponse represents the HTTP response for a product review
type ProductReviewResponse struct {
	ID                 int64    `json:"id"`
	ProductID          int64    `json:"product_id"`
	ReviewerID         int64    `json:"reviewer_id"`
	Rating             int64    `json:"rating"`
	Title              *string  `json:"title"`
	Comment            *string  `json:"comment"`
	MediaURLs          []string `json:"media_urls"`
	IsVerifiedPurchase bool     `json:"is_verified_purchase"`
	IsApproved         bool     `json:"is_approved"`
	HelpfulCount       int64    `json:"helpful_count"`
	UnhelpfulCount     int64    `json:"unhelpful_count"`
	CreatedAt          string   `json:"created_at"`
}

// ===== PRODUCT HANDLERS =====

// HandleCreateProduct handles POST /products
func (h *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get storeID from context (set by auth middleware)
	storeID := r.Context().Value("user_id").(int64)

	input := &products.CreateProductInput{
		CategoryID:         req.CategoryID,
		SKU:                req.SKU,
		Name:               req.Name,
		Slug:               req.Slug,
		Description:        req.Description,
		Price:              req.Price,
		CompareAtPrice:     req.CompareAtPrice,
		Cost:               req.Cost,
		ProductType:        req.ProductType,
		TrackInventory:     req.TrackInventory,
		RequiresShipping:   req.RequiresShipping,
		FeaturedImageURL:   req.FeaturedImageURL,
		MetaTitle:          req.MetaTitle,
		MetaDescription:    req.MetaDescription,
		MetaKeywords:       req.MetaKeywords,
		LivestreamPinnable: req.LivestreamPinnable,
	}

	product, err := h.service.CreateProduct(r.Context(), storeID, input)
	if err != nil {
		h.logger.Errorf("failed to create product: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.productToResponse(product))
}

// HandleGetProduct handles GET /products/{id}
func (h *ProductHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProductByID(r.Context(), productID)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.productToResponse(product))
}

// HandleUpdateProduct handles PUT /products/{id}
func (h *ProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := &products.UpdateProductInput{
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		Price:            req.Price,
		CompareAtPrice:   req.CompareAtPrice,
		Cost:             req.Cost,
		Status:           req.Status,
		IsVisible:        req.IsVisible,
		Featured:         req.Featured,
		FeaturedImageURL: req.FeaturedImageURL,
		MetaTitle:        req.MetaTitle,
		MetaDescription:  req.MetaDescription,
		MetaKeywords:     req.MetaKeywords,
	}

	product, err := h.service.UpdateProduct(r.Context(), productID, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.productToResponse(product))
}

// HandleDeleteProduct handles DELETE /products/{id}
func (h *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProduct(r.Context(), productID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandlePublishProduct handles POST /products/{id}/publish
func (h *ProductHandler) HandlePublishProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.PublishProduct(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.productToResponse(product))
}

// ===== PRODUCT VARIANT HANDLERS =====

// HandleCreateVariant handles POST /products/{id}/variants
func (h *ProductHandler) HandleCreateVariant(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req CreateVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := &products.CreateVariantInput{
		ProductID:         productID,
		SKU:               req.SKU,
		Title:             req.Title,
		Price:             req.Price,
		CompareAtPrice:    req.CompareAtPrice,
		Cost:              req.Cost,
		ImageURL:          req.ImageURL,
		Attributes:        req.Attributes,
		Quantity:          req.Quantity,
		LowStockThreshold: req.LowStockThreshold,
		Barcode:           req.Barcode,
		WeightKg:          req.WeightKg,
		DimensionsLength:  req.DimensionsLength,
		DimensionsWidth:   req.DimensionsWidth,
		DimensionsHeight:  req.DimensionsHeight,
	}

	variant, err := h.service.CreateProductVariant(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.variantToResponse(variant))
}

// HandleGetVariants handles GET /products/{id}/variants
func (h *ProductHandler) HandleGetVariants(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	variants, err := h.service.ListProductVariants(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]ProductVariantResponse, len(variants))
	for i, v := range variants {
		responses[i] = h.variantToResponse(v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// ===== PRODUCT MEDIA HANDLERS =====

// HandleCreateMedia handles POST /products/{id}/media
func (h *ProductHandler) HandleCreateMedia(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req CreateMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := &products.CreateMediaInput{
		ProductID:    productID,
		MediaType:    req.MediaType,
		URL:          req.URL,
		AltText:      req.AltText,
		DisplayOrder: req.DisplayOrder,
		Width:        req.Width,
		Height:       req.Height,
	}

	media, err := h.service.CreateProductMedia(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.mediaToResponse(media))
}

// HandleGetMedia handles GET /products/{id}/media
func (h *ProductHandler) HandleGetMedia(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	mediaList, err := h.service.GetProductMedia(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]ProductMediaResponse, len(mediaList))
	for i, m := range mediaList {
		responses[i] = h.mediaToResponse(m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// ===== PRODUCT REVIEW HANDLERS =====

// HandleCreateReview handles POST /products/{id}/reviews
func (h *ProductHandler) HandleCreateReview(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Get reviewerID from context (set by auth middleware)
	reviewerID := r.Context().Value("user_id").(int64)

	var req CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := &products.CreateReviewInput{
		ProductID:          productID,
		ReviewerID:         reviewerID,
		Rating:             req.Rating,
		Title:              req.Title,
		Comment:            req.Comment,
		MediaURLs:          req.MediaURLs,
		IsVerifiedPurchase: req.IsVerifiedPurchase,
	}

	review, err := h.service.CreateProductReview(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(h.reviewToResponse(review))
}

// HandleGetReviews handles GET /products/{id}/reviews
func (h *ProductHandler) HandleGetReviews(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	limit := int64(20)
	offset := int64(0)

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 64); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.ParseInt(o, 10, 64); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	reviews, err := h.service.GetProductReviews(r.Context(), productID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responses := make([]ProductReviewResponse, len(reviews))
	for i, rv := range reviews {
		responses[i] = h.reviewToResponse(rv)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responses)
}

// ===== RESPONSE CONVERTERS =====

func (h *ProductHandler) productToResponse(p *products.Product) ProductResponse {
	publishedAt := ""
	if p.PublishedAt != nil {
		publishedAt = p.PublishedAt.String()
	}

	return ProductResponse{
		ID:                 p.ID,
		StoreID:            p.StoreID,
		CategoryID:         p.CategoryID,
		SKU:                p.SKU,
		Name:               p.Name,
		Slug:               p.Slug,
		Description:        p.Description,
		Price:              p.Price,
		CompareAtPrice:     p.CompareAtPrice,
		Cost:               p.Cost,
		ProductType:        p.ProductType,
		Status:             p.Status,
		IsVisible:          p.IsVisible,
		Featured:           p.Featured,
		TrackInventory:     p.TrackInventory,
		RequiresShipping:   p.RequiresShipping,
		FeaturedImageURL:   p.FeaturedImageURL,
		MetaTitle:          p.MetaTitle,
		MetaDescription:    p.MetaDescription,
		MetaKeywords:       p.MetaKeywords,
		Rating:             p.Rating,
		ReviewCount:        p.ReviewCount,
		ViewCount:          p.ViewCount,
		SaleCount:          p.SaleCount,
		LivestreamPinnable: p.LivestreamPinnable,
		CreatedAt:          p.CreatedAt.String(),
		UpdatedAt:          p.UpdatedAt.String(),
		PublishedAt: func() *string {
			if publishedAt != "" {
				return &publishedAt
			}
			return nil
		}(),
	}
}

func (h *ProductHandler) variantToResponse(v *products.ProductVariant) ProductVariantResponse {
	return ProductVariantResponse{
		ID:                v.ID,
		ProductID:         v.ProductID,
		SKU:               v.SKU,
		Title:             v.Title,
		Price:             v.Price,
		CompareAtPrice:    v.CompareAtPrice,
		Cost:              v.Cost,
		ImageURL:          v.ImageURL,
		IsAvailable:       v.IsAvailable,
		Attributes:        v.Attributes,
		Quantity:          v.Quantity,
		ReservedQuantity:  v.ReservedQuantity,
		AvailableQuantity: v.AvailableQuantity,
		LowStockThreshold: v.LowStockThreshold,
		Barcode:           v.Barcode,
		WeightKg:          v.WeightKg,
		DimensionsLength:  v.DimensionsLength,
		DimensionsWidth:   v.DimensionsWidth,
		DimensionsHeight:  v.DimensionsHeight,
		CreatedAt:         v.CreatedAt.String(),
		UpdatedAt:         v.UpdatedAt.String(),
	}
}

func (h *ProductHandler) mediaToResponse(m *products.ProductMedia) ProductMediaResponse {
	return ProductMediaResponse{
		ID:           m.ID,
		ProductID:    m.ProductID,
		VariantID:    m.VariantID,
		MediaType:    m.MediaType,
		URL:          m.URL,
		AltText:      m.AltText,
		DisplayOrder: m.DisplayOrder,
		SizeKB:       m.SizeKB,
		Width:        m.Width,
		Height:       m.Height,
		CreatedAt:    m.CreatedAt.String(),
	}
}

func (h *ProductHandler) reviewToResponse(r *products.ProductReview) ProductReviewResponse {
	return ProductReviewResponse{
		ID:                 r.ID,
		ProductID:          r.ProductID,
		ReviewerID:         r.ReviewerID,
		Rating:             r.Rating,
		Title:              r.Title,
		Comment:            r.Comment,
		MediaURLs:          r.MediaURLs,
		IsVerifiedPurchase: r.IsVerifiedPurchase,
		IsApproved:         r.IsApproved,
		HelpfulCount:       r.HelpfulCount,
		UnhelpfulCount:     r.UnhelpfulCount,
		CreatedAt:          r.CreatedAt.String(),
	}
}

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

// CreateProductRequest represents the HTTP request body
type CreateProductRequest struct {
	SKU         string  `json:"sku" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
}

// ProductResponse represents the HTTP response body
type ProductResponse struct {
	ID          int64   `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// HandleCreateProduct handles POST /products requests
// Example: curl -X POST http://localhost:8080/products \
//   -H "Content-Type: application/json" \
//   -d '{"sku":"PROD001","name":"Product 1","description":"A great product","price":99.99}'
func (h *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("handling CreateProduct request")

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := &products.CreateProductInput{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	output, err := h.service.CreateProduct(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ProductResponse{
		ID:          output.ID,
		SKU:         output.SKU,
		Name:        output.Name,
		Description: output.Description,
		Price:       output.Price,
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	})
}

// HandleGetProduct handles GET /products/:id requests
// Example: curl http://localhost:8080/products/1
func (h *ProductHandler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.PathValue("id")
	if productIDStr == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling GetProduct request for product: %d", productID)

	input := &products.GetProductByIDInput{ProductID: productID}
	output, err := h.service.GetProductByID(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProductResponse{
		ID:          output.ID,
		SKU:         output.SKU,
		Name:        output.Name,
		Description: output.Description,
		Price:       output.Price,
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	})
}

// HandleGetProductBySKU handles GET /products/sku/:sku requests
// Example: curl http://localhost:8080/products/sku/PROD001
func (h *ProductHandler) HandleGetProductBySKU(w http.ResponseWriter, r *http.Request) {
	sku := r.PathValue("sku")
	if sku == "" {
		http.Error(w, "SKU is required", http.StatusBadRequest)
		return
	}

	h.logger.Debugf("handling GetProductBySKU request for SKU: %s", sku)

	input := &products.GetProductBySKUInput{SKU: sku}
	output, err := h.service.GetProductBySKU(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ProductResponse{
		ID:          output.ID,
		SKU:         output.SKU,
		Name:        output.Name,
		Description: output.Description,
		Price:       output.Price,
		CreatedAt:   output.CreatedAt,
		UpdatedAt:   output.UpdatedAt,
	})
}

// ListProductsResponse represents the HTTP response body
type ListProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Limit    int32             `json:"limit"`
	Offset   int32             `json:"offset"`
}

// HandleListProducts handles GET /products requests with pagination
// Example: curl "http://localhost:8080/products?limit=10&offset=0"
func (h *ProductHandler) HandleListProducts(w http.ResponseWriter, r *http.Request) {
	limit := int32(10)
	offset := int32(0)

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	h.logger.Debugf("handling ListProducts request with limit: %d, offset: %d", limit, offset)

	input := &products.ListProductsInput{Limit: limit, Offset: offset}
	output, err := h.service.ListProducts(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productResponses := make([]ProductResponse, len(output.Products))
	for i, product := range output.Products {
		productResponses[i] = ProductResponse{
			ID:          product.ID,
			SKU:         product.SKU,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListProductsResponse{
		Products: productResponses,
		Limit:    output.Limit,
		Offset:   output.Offset,
	})
}

// SearchProductsResponse represents the HTTP response body
type SearchProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Query    string            `json:"query"`
	Limit    int32             `json:"limit"`
	Offset   int32             `json:"offset"`
}

// HandleSearchProducts handles GET /products/search requests
// Example: curl "http://localhost:8080/products/search?q=laptop&limit=10"
func (h *ProductHandler) HandleSearchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limit := int32(10)
	offset := int32(0)

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil {
			offset = int32(o)
		}
	}

	h.logger.Debugf("handling SearchProducts request with query: %s", query)

	input := &products.SearchProductsInput{Query: query, Limit: limit, Offset: offset}
	output, err := h.service.SearchProducts(r.Context(), input)
	if err != nil {
		h.logger.Errorf("service error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productResponses := make([]ProductResponse, len(output.Products))
	for i, product := range output.Products {
		productResponses[i] = ProductResponse{
			ID:          product.ID,
			SKU:         product.SKU,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SearchProductsResponse{
		Products: productResponses,
		Query:    output.Query,
		Limit:    output.Limit,
		Offset:   output.Offset,
	})
}

// RegisterProductRoutes registers all product-related routes
func RegisterProductRoutes(mux *http.ServeMux, handler *ProductHandler) {
	// Product endpoints
	mux.HandleFunc("POST /products", handler.HandleCreateProduct)
	mux.HandleFunc("GET /products", handler.HandleListProducts)
	mux.HandleFunc("GET /products/{id}", handler.HandleGetProduct)
	mux.HandleFunc("GET /products/sku/{sku}", handler.HandleGetProductBySKU)
	mux.HandleFunc("GET /products/search", handler.HandleSearchProducts)
}

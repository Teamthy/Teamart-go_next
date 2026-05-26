package products

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/pkg/logger"
)

// Service provides business logic for product operations
type Service struct {
	queries *queries.Queries
	logger  *logger.Logger
}

// NewService creates a new product service
func NewService(queries *queries.Queries, logger *logger.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger,
	}
}

// CreateProductInput represents the input for creating a product
type CreateProductInput struct {
	SKU         string
	Name        string
	Description string
	Price       float64
}

// CreateProductOutput represents the output after creating a product
type CreateProductOutput struct {
	ID          int64
	SKU         string
	Name        string
	Description string
	Price       float64
	CreatedAt   string
	UpdatedAt   string
}

// CreateProduct creates a new product with validation
func (s *Service) CreateProduct(ctx context.Context, input *CreateProductInput) (*CreateProductOutput, error) {
	if input.SKU == "" {
		return nil, fmt.Errorf("SKU is required")
	}
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than zero")
	}

	s.logger.Debugf("creating product with SKU: %s", input.SKU)

	price := strconv.FormatFloat(input.Price, 'f', -1, 64)
	product, err := s.queries.CreateProduct(ctx, input.SKU, input.Name, nullableString(input.Description), price)
	if err != nil {
		s.logger.Errorf("failed to create product: %v", err)
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	priceValue, err := numericToFloat64(product.Price)
	if err != nil {
		s.logger.Errorf("failed to parse product price: %v", err)
		return nil, fmt.Errorf("failed to parse product price: %w", err)
	}

	s.logger.Infof("product created successfully with ID: %d", product.ID)

	return &CreateProductOutput{
		ID:          int64(product.ID),
		SKU:         product.Sku,
		Name:        product.Name,
		Description: product.Description.String,
		Price:       priceValue,
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}, nil
}

// GetProductByIDInput represents the input for getting a product by ID
type GetProductByIDInput struct {
	ProductID int64
}

// GetProductByIDOutput represents the output
type GetProductByIDOutput struct {
	ID          int64
	SKU         string
	Name        string
	Description string
	Price       float64
	CreatedAt   string
	UpdatedAt   string
}

// GetProductByID retrieves a product by its ID
func (s *Service) GetProductByID(ctx context.Context, input *GetProductByIDInput) (*GetProductByIDOutput, error) {
	if input.ProductID == 0 {
		return nil, fmt.Errorf("product ID is required")
	}

	s.logger.Debugf("fetching product with ID: %d", input.ProductID)

	product, err := s.queries.GetProductByID(ctx, int32(input.ProductID))
	if err != nil {
		s.logger.Errorf("failed to fetch product: %v", err)
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	priceValue, err := numericToFloat64(product.Price)
	if err != nil {
		s.logger.Errorf("failed to parse product price: %v", err)
		return nil, fmt.Errorf("failed to parse product price: %w", err)
	}

	return &GetProductByIDOutput{
		ID:          int64(product.ID),
		SKU:         product.Sku,
		Name:        product.Name,
		Description: product.Description.String,
		Price:       priceValue,
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}, nil
}

// GetProductBySKUInput represents the input
type GetProductBySKUInput struct {
	SKU string
}

// GetProductBySKUOutput represents the output
type GetProductBySKUOutput struct {
	ID          int64
	SKU         string
	Name        string
	Description string
	Price       float64
	CreatedAt   string
	UpdatedAt   string
}

// GetProductBySKU retrieves a product by its SKU
func (s *Service) GetProductBySKU(ctx context.Context, input *GetProductBySKUInput) (*GetProductBySKUOutput, error) {
	if input.SKU == "" {
		return nil, fmt.Errorf("SKU is required")
	}

	s.logger.Debugf("fetching product with SKU: %s", input.SKU)

	product, err := s.queries.GetProductBySKU(ctx, input.SKU)
	if err != nil {
		s.logger.Errorf("failed to fetch product by SKU: %v", err)
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	priceValue, err := numericToFloat64(product.Price)
	if err != nil {
		s.logger.Errorf("failed to parse product price: %v", err)
		return nil, fmt.Errorf("failed to parse product price: %w", err)
	}

	return &GetProductBySKUOutput{
		ID:          int64(product.ID),
		SKU:         product.Sku,
		Name:        product.Name,
		Description: product.Description.String,
		Price:       priceValue,
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}, nil
}

// ListProductsInput represents the input
type ListProductsInput struct {
	Limit  int32
	Offset int32
}

// ListProductsOutput represents the output
type ListProductsOutput struct {
	Products []ProductData
	Limit    int32
	Offset   int32
}

type ProductData struct {
	ID          int64
	SKU         string
	Name        string
	Description string
	Price       float64
	CreatedAt   string
	UpdatedAt   string
}

// ListProducts retrieves a list of products with pagination
func (s *Service) ListProducts(ctx context.Context, input *ListProductsInput) (*ListProductsOutput, error) {
	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	s.logger.Debugf("listing products with limit: %d, offset: %d", input.Limit, input.Offset)

	products, err := s.queries.ListProducts(ctx, input.Limit, input.Offset)
	if err != nil {
		s.logger.Errorf("failed to list products: %v", err)
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	output := &ListProductsOutput{
		Products: make([]ProductData, len(products)),
		Limit:    input.Limit,
		Offset:   input.Offset,
	}

	for i, product := range products {
		price, err := numericToFloat64(product.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to parse product price: %w", err)
		}
		output.Products[i] = ProductData{
			ID:          int64(product.ID),
			SKU:         product.Sku,
			Name:        product.Name,
			Description: product.Description.String,
			Price:       price,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		}
	}

	s.logger.Infof("fetched %d products", len(products))

	return output, nil
}

// SearchProductsInput represents the input for searching products
type SearchProductsInput struct {
	Query  string
	Limit  int32
	Offset int32
}

// SearchProductsOutput represents the output
type SearchProductsOutput struct {
	Products []ProductData
	Query    string
	Limit    int32
	Offset   int32
}

// SearchProducts searches for products by name or description
func (s *Service) SearchProducts(ctx context.Context, input *SearchProductsInput) (*SearchProductsOutput, error) {
	if input.Query == "" {
		return nil, fmt.Errorf("search query is required")
	}
	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	s.logger.Debugf("searching products with query: %s", input.Query)

	products, err := s.queries.SearchProducts(ctx, input.Query, input.Limit, input.Offset)
	if err != nil {
		s.logger.Errorf("failed to search products: %v", err)
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	output := &SearchProductsOutput{
		Products: make([]ProductData, len(products)),
		Query:    input.Query,
		Limit:    input.Limit,
		Offset:   input.Offset,
	}

	for i, product := range products {
		price, err := numericToFloat64(product.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to parse product price: %w", err)
		}
		output.Products[i] = ProductData{
			ID:          int64(product.ID),
			SKU:         product.Sku,
			Name:        product.Name,
			Description: product.Description.String,
			Price:       price,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		}
	}

	s.logger.Infof("found %d products matching query: %s", len(products), input.Query)

	return output, nil
}

func numericToFloat64(value pgtype.Numeric) (float64, error) {
	floatValue, err := value.Float64Value()
	if err != nil {
		return 0, err
	}
	return floatValue.Float64, nil
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

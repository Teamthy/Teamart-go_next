package livestream

import "context"

// ProductPinningService manages pinned product placement during livestreams.
type ProductPinningService struct{}

// NewProductPinningService creates a new product pinning service.
func NewProductPinningService() *ProductPinningService {
	return &ProductPinningService{}
}

// PinProduct pins a product into a livestream overlay.
func (s *ProductPinningService) PinProduct(ctx context.Context, streamID string, productID int64) error {
	return nil
}

// UnpinProduct removes a pinned product.
func (s *ProductPinningService) UnpinProduct(ctx context.Context, streamID string, productID int64) error {
	return nil
}

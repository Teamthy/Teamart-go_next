package trending

import (
	"context"
	"fmt"
)

// TrendingService computes trending products, creators, and streams.
type TrendingService struct {
	WindowMinutes int
}

// NewTrendingService creates a trending service.
func NewTrendingService(windowMinutes int) *TrendingService {
	return &TrendingService{WindowMinutes: windowMinutes}
}

// GetTrendingProducts returns a list of trending product IDs.
func (s *TrendingService) GetTrendingProducts(ctx context.Context, limit int) ([]int64, error) {
	fmt.Printf("fetching top %d trending products\n", limit)
	return []int64{}, nil
}

// GetTrendingCreators returns a list of trending creator IDs.
func (s *TrendingService) GetTrendingCreators(ctx context.Context, limit int) ([]int64, error) {
	fmt.Printf("fetching top %d trending creators\n", limit)
	return []int64{}, nil
}

// GetTrendingStreams returns a list of trending livestream IDs.
func (s *TrendingService) GetTrendingStreams(ctx context.Context, limit int) ([]string, error) {
	fmt.Printf("fetching top %d trending streams\n", limit)
	return []string{}, nil
}

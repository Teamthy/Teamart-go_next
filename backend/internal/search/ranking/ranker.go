package ranking

import (
	"context"
	"fmt"
)

// Ranker computes relevance scores for search results.
type Ranker struct {
	BoostFactors map[string]float64
}

// NewRanker creates a new ranking engine.
func NewRanker(boostFactors map[string]float64) *Ranker {
	return &Ranker{BoostFactors: boostFactors}
}

// ScoreProduct returns a relevance score for a product result.
func (r *Ranker) ScoreProduct(ctx context.Context, productID int64, metadata map[string]interface{}) float64 {
	score := 1.0
	fmt.Printf("scoring product %d with metadata %v\n", productID, metadata)
	return score
}

// ScoreCreator returns a relevance score for a creator result.
func (r *Ranker) ScoreCreator(ctx context.Context, creatorID int64, metadata map[string]interface{}) float64 {
	score := 1.0
	fmt.Printf("scoring creator %d with metadata %v\n", creatorID, metadata)
	return score
}

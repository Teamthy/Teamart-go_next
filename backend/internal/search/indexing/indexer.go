package indexing

import (
	"context"
	"fmt"
	"time"
)

// SearchIndexer defines the interface for indexing search data.
type SearchIndexer interface {
	IndexProduct(ctx context.Context, productID int64, payload map[string]interface{}) error
	IndexCreator(ctx context.Context, creatorID int64, payload map[string]interface{}) error
	DeleteIndex(ctx context.Context, indexName string, id string) error
}

// MeilisearchIndexer provides Meilisearch-backed indexing.
type MeilisearchIndexer struct {
	Endpoint string
	APIKey   string
}

// NewMeilisearchIndexer creates a new Meilisearch indexer.
func NewMeilisearchIndexer(endpoint, apiKey string) *MeilisearchIndexer {
	return &MeilisearchIndexer{Endpoint: endpoint, APIKey: apiKey}
}

// IndexProduct indexes a product for search.
func (m *MeilisearchIndexer) IndexProduct(ctx context.Context, productID int64, payload map[string]interface{}) error {
	// TODO: implement Meilisearch document upsert
	fmt.Printf("indexing product %d to Meilisearch\n", productID)
	return nil
}

// IndexCreator indexes a creator profile for discovery.
func (m *MeilisearchIndexer) IndexCreator(ctx context.Context, creatorID int64, payload map[string]interface{}) error {
	fmt.Printf("indexing creator %d to Meilisearch\n", creatorID)
	return nil
}

// DeleteIndex removes a document from the search index.
func (m *MeilisearchIndexer) DeleteIndex(ctx context.Context, indexName string, id string) error {
	fmt.Printf("deleting index document %s from %s\n", id, indexName)
	return nil
}

// RefreshIndex triggers a refresh cycle when new data arrives.
func (m *MeilisearchIndexer) RefreshIndex(ctx context.Context, indexName string) error {
	fmt.Printf("refreshing index %s at %s\n", indexName, time.Now().UTC())
	return nil
}

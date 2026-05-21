package analytics

import (
	"context"
	"fmt"
	"time"
)

// SearchAnalytics tracks queries, clicks, and conversion signals.
type SearchAnalytics struct {
	QueryLog map[string]int
}

// NewSearchAnalytics creates a new search analytics service.
func NewSearchAnalytics() *SearchAnalytics {
	return &SearchAnalytics{QueryLog: make(map[string]int)}
}

// TrackQuery logs a search query event.
func (s *SearchAnalytics) TrackQuery(ctx context.Context, userID int64, query string) {
	s.QueryLog[query]++
	fmt.Printf("tracked search query '%s' for user %d at %s\n", query, userID, time.Now().UTC())
}

// TrackClick logs a click on a search result.
func (s *SearchAnalytics) TrackClick(ctx context.Context, userID int64, resultID string) {
	fmt.Printf("tracked click on %s for user %d\n", resultID, userID)
}

// GetTopQueries returns the most frequent queries.
func (s *SearchAnalytics) GetTopQueries(limit int) []string {
	top := make([]string, 0, limit)
	for query := range s.QueryLog {
		top = append(top, query)
		if len(top) >= limit {
			break
		}
	}
	return top
}

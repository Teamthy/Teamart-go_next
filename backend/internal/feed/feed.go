package feed

import rec "github.com/teamart/commerce-api/internal/recommendation"

// FeedItem is an alias for recommendation.Recommendation to keep feed-focused semantics.
type FeedItem = rec.Recommendation

// Service composes a recommendation service to expose a feed API.
type Service struct {
	rec rec.RecommendationService
}

func NewService(r rec.RecommendationService) *Service {
	return &Service{rec: r}
}

// GetFeedForUser returns a ranked feed for a user. This is a thin wrapper
// around the RecommendationService to keep feed-specific concerns separate.
func (s *Service) GetFeedForUser(userID string, limit int) ([]FeedItem, error) {
	return s.rec.RecommendForUser(userID, limit)
}

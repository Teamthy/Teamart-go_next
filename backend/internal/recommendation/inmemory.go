package recommendation

import (
	"sort"
)

// RecommendationCandidate represents a candidate item which can be scored.
type RecommendationCandidate struct {
	ID       string
	Type     ItemType
	Category string
	Signals  RankingSignals
	Meta     map[string]interface{}
}

// InMemoryRecommendationService is a simple, test-friendly implementation.
type InMemoryRecommendationService struct {
	items   []RecommendationCandidate
	weights Weights
}

func NewInMemoryRecommendationService(cands []RecommendationCandidate, weights Weights) *InMemoryRecommendationService {
	return &InMemoryRecommendationService{items: cands, weights: weights}
}

func (s *InMemoryRecommendationService) RecommendForUser(userID string, limit int) ([]Recommendation, error) {
	results := make([]Recommendation, 0, len(s.items))
	for _, c := range s.items {
		score := Score(c.Signals, s.weights, c.Category)
		results = append(results, Recommendation{
			ID:    c.ID,
			Type:  c.Type,
			Score: score,
			Meta:  c.Meta,
		})
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	if limit <= 0 || limit > len(results) {
		limit = len(results)
	}
	return results[:limit], nil
}

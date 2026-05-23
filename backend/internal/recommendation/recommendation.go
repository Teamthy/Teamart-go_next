package recommendation

type ItemType string

const (
	ItemTypeLivestream ItemType = "livestream"
	ItemTypeCreator    ItemType = "creator"
	ItemTypeProduct    ItemType = "product"
)

// RankingSignals captures user and item interaction signals used for ranking.
type RankingSignals struct {
	WatchTime        float64
	Purchases        float64
	Reactions        float64
	Follows          float64
	CategoryAffinity map[string]float64
}

// Recommendation is a ranked item returned by the recommendation engine.
type Recommendation struct {
	ID    string
	Type  ItemType
	Score float64
	Meta  map[string]interface{}
}

// RecommendationService returns ranked recommendations for a user.
type RecommendationService interface {
	RecommendForUser(userID string, limit int) ([]Recommendation, error)
}

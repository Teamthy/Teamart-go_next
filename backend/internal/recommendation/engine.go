package recommendation

// Weights defines contribution of each signal to the final score.
type Weights struct {
	WatchTime        float64
	Purchases        float64
	Reactions        float64
	Follows          float64
	CategoryAffinity float64
}

// Score computes a simple linear combination of signals using provided weights.
func Score(signals RankingSignals, weights Weights, category string) float64 {
	score := weights.WatchTime*signals.WatchTime +
		weights.Purchases*signals.Purchases +
		weights.Reactions*signals.Reactions +
		weights.Follows*signals.Follows

	if affinity, ok := signals.CategoryAffinity[category]; ok {
		score += weights.CategoryAffinity * affinity
	}

	return score
}

package recommendation

import "testing"

func TestScore(t *testing.T) {
	signals := RankingSignals{
		WatchTime:        2.0,
		Purchases:        1.0,
		Reactions:        3.0,
		Follows:          0.5,
		CategoryAffinity: map[string]float64{"sports": 1.2},
	}

	weights := Weights{WatchTime: 1, Purchases: 1, Reactions: 1, Follows: 1, CategoryAffinity: 2}

	got := Score(signals, weights, "sports")
	want := 8.9 // 2 + 1 + 3 + 0.5 + 2*1.2 == 8.9

	if diff := got - want; diff > 1e-9 || diff < -1e-9 {
		t.Fatalf("Score() = %v, want %v", got, want)
	}
}

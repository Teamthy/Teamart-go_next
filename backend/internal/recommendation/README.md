Recommendation package

This package contains a minimal recommendation engine prototype used to
demonstrate ranking signals and a testable in-memory implementation.

Features:
- `RankingSignals` captures watch time, purchases, reactions, follows, and category affinity.
- `Score` implements a linear scoring function with configurable weights.
- `InMemoryRecommendationService` is a simple candidate scorer useful for local dev and tests.

Next steps:
- Add persistent candidate storage (Postgres/Redis)
- Add user-signal aggregation pipelines (events -> aggregated signals)
- Add ML model or more advanced ranking (GAMs, tree ensembles, neural networks)

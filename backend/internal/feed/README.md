Feed package

The `feed` package provides a thin abstraction over the recommendation
service to present ranked feed items to callers.

It is intentionally small — business logic (ranking, signals) lives in
the `internal/recommendation` package so it can be evolved independently.

Next steps:
- Add paging cursors (cursor-based pagination)
- Add personalization context (device, session, recent activity)
- Integrate with API handlers under `internal/handlers`

Admin & Operations package

This package contains lightweight admin domain types and a simple in-memory
service useful for local development and tests. It provides operations for:

- Dashboard summary
- Dispute management
- Payout approval
- Moderation actions
- Creator verification
- Fraud alerts
- Support actions (refunds, suspension)

This is intentionally minimal; production systems should wire these operations
to persistent stores and durable workflows.

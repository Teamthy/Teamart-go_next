-- queries/support.sql
-- Support tickets and messages queries

-- ===== SUPPORT TICKETS =====

-- name: CreateSupportTicket :one
INSERT INTO support_tickets (
    ticket_number, merchant_id, customer_id, order_id, subject, description,
    category, priority, status, created_by, sla_due_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: GetSupportTicket :one
SELECT * FROM support_tickets WHERE id = $1;

-- name: GetTicketByNumber :one
SELECT * FROM support_tickets WHERE ticket_number = $1;

-- name: ListMerchantTickets :many
SELECT * FROM support_tickets 
WHERE merchant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCustomerTickets :many
SELECT * FROM support_tickets 
WHERE customer_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTicketsByStatus :many
SELECT * FROM support_tickets 
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTicketsByPriority :many
SELECT * FROM support_tickets 
WHERE priority = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAssignedTickets :many
SELECT * FROM support_tickets 
WHERE assigned_to = $1
ORDER BY sla_due_at ASC;

-- name: ListOpenTickets :many
SELECT * FROM support_tickets 
WHERE status IN ('open', 'pending', 'in_progress')
ORDER BY priority DESC, created_at ASC;

-- name: UpdateTicketStatus :one
UPDATE support_tickets 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: AssignTicket :one
UPDATE support_tickets 
SET assigned_to = $2, assigned_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ResolveTicket :one
UPDATE support_tickets 
SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP,
    resolved_by = $2, resolution_summary = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CloseTicket :one
UPDATE support_tickets 
SET status = 'closed', closed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: EscalateTicket :one
UPDATE support_tickets 
SET status = 'escalated', escalated_to = $2, escalation_reason = $3,
    escalated_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== SUPPORT MESSAGES =====

-- name: CreateSupportMessage :one
INSERT INTO support_messages (
    ticket_id, sender_id, sender_role, message, is_internal
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetSupportMessage :one
SELECT * FROM support_messages WHERE id = $1;

-- name: ListTicketMessages :many
SELECT * FROM support_messages 
WHERE ticket_id = $1 AND is_internal = FALSE
ORDER BY created_at ASC;

-- name: ListTicketAllMessages :many
SELECT * FROM support_messages 
WHERE ticket_id = $1
ORDER BY created_at ASC;

-- name: ListUserMessages :many
SELECT * FROM support_messages 
WHERE sender_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountTicketMessages :one
SELECT COUNT(*) FROM support_messages WHERE ticket_id = $1 AND is_internal = FALSE;

-- name: MarkMessageAsSolution :one
UPDATE support_messages 
SET marked_as_solution = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- ===== ANALYTICS =====

-- name: GetSupportStats :one
SELECT 
    COUNT(*) as total_tickets,
    COUNT(CASE WHEN status = 'open' THEN 1 END) as open_tickets,
    COUNT(CASE WHEN status = 'resolved' THEN 1 END) as resolved_tickets,
    COUNT(CASE WHEN status = 'closed' THEN 1 END) as closed_tickets,
    COUNT(CASE WHEN sla_met = FALSE THEN 1 END) as sla_breaches
FROM support_tickets 
WHERE created_at BETWEEN $1 AND $2;

-- name: GetAverageResolutionTime :one
SELECT 
    AVG(EXTRACT(EPOCH FROM (resolved_at - created_at))) as avg_resolution_time_seconds,
    MIN(EXTRACT(EPOCH FROM (resolved_at - created_at))) as min_resolution_time_seconds,
    MAX(EXTRACT(EPOCH FROM (resolved_at - created_at))) as max_resolution_time_seconds
FROM support_tickets 
WHERE resolved_at IS NOT NULL AND created_at BETWEEN $1 AND $2;

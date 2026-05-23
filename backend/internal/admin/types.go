package admin

import "time"

type DisputeStatus string
type PayoutStatus string

const (
	DisputeOpen   DisputeStatus = "open"
	DisputeClosed DisputeStatus = "closed"

	PayoutPending PayoutStatus = "pending"
	PayoutPaid    PayoutStatus = "paid"
)

type Dispute struct {
	ID        string
	OrderID   string
	UserID    int64
	Amount    float64
	Reason    string
	Status    DisputeStatus
	CreatedAt time.Time
}

type Payout struct {
	ID        string
	CreatorID string
	Amount    float64
	Status    PayoutStatus
	CreatedAt time.Time
}

type ModerationAction struct {
	ID        string
	TargetID  string
	Action    string
	Reason    string
	Performed time.Time
}

type CreatorVerification struct {
	CreatorID string
	Verified  bool
	Notes     string
}

type FraudAlert struct {
	ID        string
	SubjectID string
	Score     float64
	Data      map[string]interface{}
	CreatedAt time.Time
}

type AuditLog struct {
	ID           string
	ActorID      int64
	Action       string
	ResourceType string
	ResourceID   string
	Details      map[string]interface{}
	CreatedAt    time.Time
}

type Notification struct {
	ID          string
	RecipientID string
	Channel     string
	Payload     map[string]interface{}
	SentAt      *time.Time
}

type PayoutApproval struct {
	ID          string
	PayoutID    string
	RequestedBy int64
	Status      string
	Notes       string
	CreatedAt   time.Time
	ApprovedBy  *int64
	ApprovedAt  *time.Time
}

type DashboardSummary struct {
	OpenDisputes   int
	PendingPayouts int
	FraudAlerts    int
}

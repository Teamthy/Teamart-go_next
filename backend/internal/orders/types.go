package orders

import (
	"time"
)

// ===== ORDER TYPES =====

// Order represents a marketplace order
type Order struct {
	ID                    int64
	OrderNumber           string
	UserID                int64
	Source                string // web, mobile, livestream, admin
	CreatorID             *int64 // For livestream orders
	CustomerEmail         string
	CustomerPhone         *string
	Status                string // pending, paid, processing, shipped, delivered, completed, cancelled, refunded, disputed
	Subtotal              float64
	ShippingCost          float64
	TaxAmount             float64
	DiscountAmount        float64
	Total                 float64
	PaidAmount            float64
	PaymentMethod         *string
	PaymentID             *string
	PaymentStatus         string // pending, completed, failed
	ShippingMethod        *string
	TrackingNumber        *string
	EstimatedDeliveryDate *time.Time
	ShippedAt             *time.Time
	DeliveredAt           *time.Time
	BillingAddress        map[string]interface{}
	ShippingAddress       map[string]interface{}
	Notes                 *string
	Metadata              map[string]interface{}
	CancelledAt           *time.Time
	CancellationReason    *string
	RefundedAt            *time.Time
	RefundReason          *string
	RefundAmount          *float64
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Items                 []*OrderItem `db:"-"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID                int64
	OrderID           int64
	ProductID         int64
	VariantID         int64
	SKU               string
	ProductName       string
	VariantTitle      *string
	Quantity          int64
	UnitPrice         float64
	LineTotal         float64
	FulfillmentStatus string // pending, processing, shipped, delivered
	FulfillmentID     *string
	ReturnStatus      *string // none, pending, approved, rejected, completed
	ReturnQuantity    int64
	CreatedAt         time.Time
}

// Cart represents a shopping cart
type Cart struct {
	ID                      int64
	UserID                  int64
	StoreID                 *int64
	Status                  string // active, abandoned, converted
	Subtotal                float64
	ShippingAmount          float64
	TaxAmount               float64
	DiscountAmount          float64
	Total                   float64
	CouponCode              *string
	EstimatedShippingMethod *string
	EstimatedShippingCost   *float64
	CreatedAt               time.Time
	UpdatedAt               time.Time
	AbandonedAt             *time.Time
	ConvertedAt             *time.Time
	Items                   []*CartItem `db:"-"`
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        int64
	CartID    int64
	ProductID int64
	VariantID int64
	Quantity  int64
	UnitPrice float64
	LineTotal float64
	AddedAt   time.Time
	UpdatedAt time.Time
}

// OrderEvent represents a status change in order timeline
type OrderEvent struct {
	ID          int64
	OrderID     int64
	EventType   string
	Status      string
	NewStatus   string
	Description *string
	ActorID     *int64
	Metadata    map[string]interface{}
	CreatedAt   time.Time
}

// Fulfillment represents shipment fulfillment
type Fulfillment struct {
	ID                    int64
	OrderID               int64
	Status                string // pending, processing, shipped, delivered, failed
	ShippingCarrier       *string
	TrackingNumber        *string
	TrackingURL           *string
	ShippedAt             *time.Time
	DeliveredAt           *time.Time
	EstimatedDeliveryDate *time.Time
	ShippingAddress       map[string]interface{}
	ItemsCount            int64
	Notes                 *string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// OrderReturn represents a product return
type OrderReturn struct {
	ID                 int64
	OrderID            int64
	Status             string // pending, approved, received, refunded, completed
	Reason             string
	Description        *string
	InitiatorID        int64
	InitiatedAt        time.Time
	AuthorizedBy       *int64
	AuthorizedAt       *time.Time
	AuthorizationNotes *string
	ReceivedAt         *time.Time
	ReceivedBy         *int64
	InspectionNotes    *string
	RefundAmount       *float64
	RefundStatus       *string
	RefundID           *string
	RefundedAt         *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Items              []*ReturnItem `db:"-"`
}

// ReturnItem represents an item being returned
type ReturnItem struct {
	ID          int64
	ReturnID    int64
	OrderItemID int64
	Quantity    int64
	Reason      *string
	CreatedAt   time.Time
}

// OrderDispute represents a chargeback or dispute
type OrderDispute struct {
	ID             int64
	OrderID        int64
	Status         string // open, under_review, resolved, escalated, closed
	DisputeType    string // chargeback, not_received, quality_issue, fraud, etc.
	Reason         string
	InitiatorID    int64
	InitiatedAt    time.Time
	AssignedTo     *int64
	Resolution     *string
	ResolvedAt     *time.Time
	ResolvedBy     *int64
	DisputedAmount float64
	ResolvedAmount *float64
	Deadline       *time.Time
	Evidence       *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Coupon represents a discount coupon
type Coupon struct {
	ID                 int64
	Code               string
	Description        *string
	DiscountType       string // percentage, fixed_amount, free_shipping
	DiscountValue      float64
	MinOrderValue      *float64
	MaxDiscountAmount  *float64
	MaxUses            *int64
	MaxUsesPerCustomer int64
	ValidFrom          time.Time
	ValidUntil         time.Time
	TimesUsed          int64
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// GiftCard represents a gift card
type GiftCard struct {
	ID              int64
	Code            string
	Balance         float64
	OriginalBalance float64
	OwnerID         *int64
	RecipientEmail  *string
	ValidFrom       time.Time
	ValidUntil      *time.Time
	Status          string // active, redeemed, expired, cancelled
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ShippingRate represents a shipping rate rule
type ShippingRate struct {
	ID            int64
	StoreID       int64
	Name          string
	Description   *string
	MinWeight     *float64
	MaxWeight     *float64
	MinOrderValue *float64
	MaxOrderValue *float64
	Countries     []string // JSON array
	BaseRate      float64
	PerUnitRate   *float64
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ===== INPUT/OUTPUT TYPES =====

// AddToCartInput represents the input for adding an item to cart
type AddToCartInput struct {
	VariantID int64
	Quantity  int64
}

// UpdateCartItemInput represents the input for updating a cart item
type UpdateCartItemInput struct {
	CartItemID int64
	Quantity   int64
}

// CreateOrderInput represents the input for creating an order
type CreateOrderInput struct {
	UserID          int64
	Source          string
	CreatorID       *int64
	CustomerEmail   string
	CustomerPhone   *string
	ShippingAddress map[string]interface{}
	BillingAddress  map[string]interface{}
	ShippingMethod  string
	CouponCode      *string
}

// UpdateOrderStatusInput represents the input for updating order status
type UpdateOrderStatusInput struct {
	Status string
	Reason *string
}

// CancelOrderInput represents the input for cancelling an order
type CancelOrderInput struct {
	Reason string
}

// CreateReturnInput represents the input for creating a return
type CreateReturnInput struct {
	Reason      string
	Description *string
	ItemIDs     []int64         // Order item IDs to return
	Quantities  map[int64]int64 // Item ID -> Quantity mapping
}

// ApproveReturnInput represents the input for approving a return
type ApproveReturnInput struct {
	Notes string
}

// ProcessRefundInput represents the input for processing a refund
type ProcessRefundInput struct {
	ReturnID int64
	Amount   float64
}

// CreateDisputeInput represents the input for creating a dispute
type CreateDisputeInput struct {
	DisputeType string
	Reason      string
	Evidence    *string
	Amount      float64
}

// ResolveDisputeInput represents the input for resolving a dispute
type ResolveDisputeInput struct {
	Resolution     string
	ResolvedAmount float64
}

// CreateCouponInput represents the input for creating a coupon
type CreateCouponInput struct {
	Code               string
	Description        *string
	DiscountType       string
	DiscountValue      float64
	MinOrderValue      *float64
	MaxDiscountAmount  *float64
	MaxUses            *int64
	MaxUsesPerCustomer int64
	ValidFrom          time.Time
	ValidUntil         time.Time
}

// ===== SEARCH & FILTER TYPES =====

// OrderSearchCriteria represents search criteria for orders
type OrderSearchCriteria struct {
	UserID    *int64
	Status    *string
	StartDate *time.Time
	EndDate   *time.Time
	MinAmount *float64
	MaxAmount *float64
	SortBy    string // created_at, total, status
	Limit     int64
	Offset    int64
}

// ===== RESPONSE/VO TYPES =====

// OrderSummary is a lightweight order representation
type OrderSummary struct {
	ID          int64
	OrderNumber string
	Status      string
	Total       float64
	ItemCount   int64
	CreatedAt   time.Time
}

// CartSummary is a lightweight cart representation
type CartSummary struct {
	ID        int64
	ItemCount int64
	Total     float64
	UpdatedAt time.Time
}

package commerce

import "time"

type ProductDetails struct {
	ID          int64   `json:"id"`
	SKU         string  `json:"sku,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock,omitempty"`
}

type ProductPin struct {
	StreamID     string    `json:"stream_id"`
	ProductID    int64     `json:"product_id"`
	DisplayTitle string    `json:"display_title,omitempty"`
	OverlayText  string    `json:"overlay_text,omitempty"`
	PromotionID  string    `json:"promotion_id,omitempty"`
	DisplayOrder int       `json:"display_order"`
	PinnedAt     time.Time `json:"pinned_at"`
	UnpinnedAt   time.Time `json:"unpinned_at,omitempty"`
}

type LivePromotion struct {
	ID              string    `json:"id"`
	Label           string    `json:"label"`
	DiscountPercent float64   `json:"discount_percent"`
	FlashDiscount   bool      `json:"flash_discount"`
	StartAt         time.Time `json:"start_at"`
	EndAt           time.Time `json:"end_at"`
	Active          bool      `json:"active"`
}

type LiveCartItem struct {
	ProductID  int64   `json:"product_id"`
	SKU        string  `json:"sku,omitempty"`
	Name       string  `json:"name,omitempty"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
}

type LiveCart struct {
	StreamID    string          `json:"stream_id"`
	UserID      int64           `json:"user_id"`
	Items       []LiveCartItem  `json:"items"`
	TotalAmount float64         `json:"total_amount"`
	Discount    float64         `json:"discount"`
	FinalAmount float64         `json:"final_amount"`
	Promotions  []LivePromotion `json:"promotions,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type CommissionSplit struct {
	Platform  float64 `json:"platform_cut"`
	Creator   float64 `json:"creator_cut"`
	Affiliate float64 `json:"affiliate_cut"`
	Seller    float64 `json:"seller_payout"`
}

type PurchaseAttribution struct {
	CreatorID   int64           `json:"creator_id"`
	AffiliateID *int64          `json:"affiliate_id,omitempty"`
	Split       CommissionSplit `json:"split"`
}

type PurchaseReceipt struct {
	PurchaseID       string              `json:"purchase_id"`
	StreamID         string              `json:"stream_id"`
	UserID           int64               `json:"user_id"`
	OrderReference   string              `json:"order_reference"`
	PaymentReference string              `json:"payment_reference"`
	Amount           float64             `json:"amount"`
	Attribution      PurchaseAttribution `json:"attribution"`
	CreatedAt        time.Time           `json:"created_at"`
}

type OrderInput struct {
	UserID      int64
	StreamID    string
	TotalAmount float64
	Status      string
}

type OrderResult struct {
	OrderID   int64
	Reference string
}

type PaymentResult struct {
	PaymentID string
	Method    string
}

type CommissionCalculator interface {
	Calculate(amount float64, creatorID int64, affiliateID *int64) CommissionSplit
}

type ProductFetcher interface {
	FetchProduct(productID int64) (*ProductDetails, error)
}

type OrderCreator interface {
	CreateOrder(input *OrderInput) (*OrderResult, error)
}

type PaymentProcessor interface {
	Charge(userID int64, amount float64, method string) (*PaymentResult, error)
}

type WalletManager interface {
	RecordPayout(userID int64, amount float64, reason string) error
}

type EventDispatcher interface {
	PublishEvent(eventType string, payload map[string]interface{}) error
}

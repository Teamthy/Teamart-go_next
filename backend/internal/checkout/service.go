package checkout

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/pkg/logger"
)

// CheckoutService handles cart validation, address validation, shipping rate selection,
// and order creation as part of the multi-step checkout flow
type CheckoutService struct {
	queries *queries.Queries
	logger  *logger.Logger
}

// NewCheckoutService creates a new checkout service instance
func NewCheckoutService(q *queries.Queries, logger *logger.Logger) *CheckoutService {
	return &CheckoutService{
		queries: q,
		logger:  logger,
	}
}

// ===== INPUT/OUTPUT TYPES =====

// CalculateCheckoutRequest validates cart and calculates totals
type CalculateCheckoutRequest struct {
	CartItems       []CartItem           `json:"cart_items"`
	AddressID       *int64               `json:"address_id"`
	ShippingAddress *AddressInput        `json:"shipping_address"`
	DiscountCode    *string              `json:"discount_code"`
	PromoCodeID     *int64               `json:"promo_code_id"`
	ShippingZone    string               `json:"shipping_zone"`
	MerchantShipping map[int64]ShippingConfig `json:"merchant_shipping"`
}

// CartItem represents an item in the checkout cart
type CartItem struct {
	ProductID int64  `json:"product_id"`
	VariantID int64  `json:"variant_id"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price_cents"` // Price in cents
	MerchantID int64 `json:"merchant_id"`
}

// AddressInput represents address information for checkout
type AddressInput struct {
	RecipientName  string `json:"recipient_name"`
	PhoneNumber    string `json:"phone_number"`
	StreetAddress  string `json:"street_address"`
	ApartmentSuite string `json:"apartment_suite"`
	City           string `json:"city"`
	StateProvince  string `json:"state_province"`
	PostalCode     string `json:"postal_code"`
	CountryCode    string `json:"country_code"`
}

// ShippingConfig represents merchant shipping configuration
type ShippingConfig struct {
	ProfileID int64  `json:"profile_id"`
	Carrier   string `json:"carrier"`
	Service   string `json:"service"`
}

// CalculateCheckoutResponse returns calculated totals
type CalculateCheckoutResponse struct {
	Subtotal               int64                      `json:"subtotal_cents"`
	DiscountAmount         int64                      `json:"discount_cents"`
	TaxAmount              int64                      `json:"tax_cents"`
	ShippingCost           int64                      `json:"shipping_cents"`
	Total                  int64                      `json:"total_cents"`
	DeliveryEstimate       DeliveryEstimate           `json:"delivery_estimate"`
	ShippingOptions        []ShippingOption           `json:"shipping_options"`
	AppliedDiscount        *AppliedDiscount           `json:"applied_discount,omitempty"`
	ItemsBreakdown         []ItemBreakdown            `json:"items_breakdown"`
	WarningsAndSuggestions []string                   `json:"warnings_and_suggestions"`
}

// DeliveryEstimate provides expected delivery timeline
type DeliveryEstimate struct {
	MinDays             int       `json:"min_days"`
	MaxDays             int       `json:"max_days"`
	EstimatedDeliveryDate string `json:"estimated_delivery_date"`
}

// ShippingOption represents available shipping methods
type ShippingOption struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Carrier     string `json:"carrier"`
	Service     string `json:"service"`
	Cost        int64  `json:"cost_cents"`
	MinDays     int    `json:"min_days"`
	MaxDays     int    `json:"max_days"`
	IsDefault   bool   `json:"is_default"`
	Restrictions string `json:"restrictions,omitempty"`
}

// AppliedDiscount represents an applied discount or promo
type AppliedDiscount struct {
	Code           string `json:"code"`
	DiscountType   string `json:"discount_type"` // percentage, fixed_amount, free_shipping
	DiscountValue  int64  `json:"discount_value_cents"`
	DiscountAmount int64  `json:"discount_amount_cents"`
	Description    string `json:"description"`
}

// ItemBreakdown shows per-item cost details
type ItemBreakdown struct {
	ProductID      int64  `json:"product_id"`
	VariantID      int64  `json:"variant_id"`
	Quantity       int    `json:"quantity"`
	UnitPrice      int64  `json:"unit_price_cents"`
	Subtotal       int64  `json:"subtotal_cents"`
	Discount       int64  `json:"discount_cents"`
	Tax            int64  `json:"tax_cents"`
	Total          int64  `json:"total_cents"`
	MerchantID     int64  `json:"merchant_id"`
	ShippingCost   int64  `json:"shipping_cents"`
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	CustomerID        int64           `json:"customer_id"`
	CustomerEmail     string          `json:"customer_email"`
	CustomerPhone     string          `json:"customer_phone"`
	CartItems         []CartItem      `json:"cart_items"`
	BillingAddress    *AddressInput   `json:"billing_address"`
	ShippingAddress   *AddressInput   `json:"shipping_address"`
	ShippingMethod    ShippingConfig  `json:"shipping_method"`
	PaymentMethod     string          `json:"payment_method"` // card, wallet, bank_transfer, etc.
	PaymentMethodID   *string         `json:"payment_method_id"`
	DiscountCode      *string         `json:"discount_code"`
	PromoCodeID       *int64          `json:"promo_code_id"`
	OrderNotes        *string         `json:"order_notes"`
	Source            string          `json:"source"` // web, mobile, livestream, admin
	CreatorID         *int64          `json:"creator_id"`
	Metadata          interface{}     `json:"metadata"`
}

// CreateOrderResponse returns created order details
type CreateOrderResponse struct {
	OrderID         int64  `json:"order_id"`
	OrderNumber     string `json:"order_number"`
	PaymentIntentID string `json:"payment_intent_id"`
	ClientSecret    string `json:"client_secret"`
	Status          string `json:"status"`
	Total           int64  `json:"total_cents"`
	NextAction      string `json:"next_action"` // redirect_to_payment, confirm_payment, etc.
}

// ===== SERVICE METHODS =====

// CalculateCheckout calculates totals and validates checkout data
func (s *CheckoutService) CalculateCheckout(ctx context.Context, req *CalculateCheckoutRequest) (*CalculateCheckoutResponse, error) {
	s.logger.Debugf("calculating checkout: items_count=%d", len(req.CartItems))

	// Validate cart items
	if len(req.CartItems) == 0 {
		return nil, errors.New("cart cannot be empty")
	}

	// Calculate subtotal
	subtotal := int64(0)
	for _, item := range req.CartItems {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}
		itemTotal := item.UnitPrice * int64(item.Quantity)
		subtotal += itemTotal
	}

	// Calculate discount
	discountAmount := int64(0)
	var appliedDiscount *AppliedDiscount
	if req.DiscountCode != nil {
		// TODO: Call discount service to validate and calculate discount
		// For now, placeholder
		discountAmount = subtotal / 10 // Example: 10% discount
		appliedDiscount = &AppliedDiscount{
			Code:           *req.DiscountCode,
			DiscountType:   "percentage",
			DiscountValue:  10,
			DiscountAmount: discountAmount,
		}
	}

	// Calculate tax (simplified - assumes 7.5% VAT for Nigeria)
	taxAmount := (subtotal - discountAmount) * 75 / 1000

	// Calculate shipping (TODO: integrate with shipping service)
	shippingCost := int64(500) // Example: 500 naira base shipping

	// Calculate total
	total := subtotal - discountAmount + taxAmount + shippingCost

	// Create items breakdown
	itemsBreakdown := make([]ItemBreakdown, len(req.CartItems))
	for i, item := range req.CartItems {
		itemSubtotal := item.UnitPrice * int64(item.Quantity)
		itemDiscount := (itemSubtotal * discountAmount) / subtotal // Proportional discount
		itemTax := (itemSubtotal - itemDiscount) * 75 / 1000
		itemShipping := shippingCost / int64(len(req.CartItems)) // Proportional shipping

		itemsBreakdown[i] = ItemBreakdown{
			ProductID:    item.ProductID,
			VariantID:    item.VariantID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			Subtotal:     itemSubtotal,
			Discount:     itemDiscount,
			Tax:          itemTax,
			Total:        itemSubtotal - itemDiscount + itemTax + itemShipping,
			MerchantID:   item.MerchantID,
			ShippingCost: itemShipping,
		}
	}

	// Estimate delivery
	deliveryEstimate := DeliveryEstimate{
		MinDays:             1,
		MaxDays:             3,
		EstimatedDeliveryDate: time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
	}

	// Shipping options
	shippingOptions := []ShippingOption{
		{
			ID:        "standard",
			Name:      "Standard Delivery",
			Carrier:   "GIG",
			Service:   "standard",
			Cost:      500,
			MinDays:   1,
			MaxDays:   3,
			IsDefault: true,
		},
		{
			ID:      "express",
			Name:    "Express Delivery",
			Carrier: "DHL",
			Service: "express",
			Cost:    1500,
			MinDays: 0,
			MaxDays: 1,
		},
	}

	return &CalculateCheckoutResponse{
		Subtotal:               subtotal,
		DiscountAmount:         discountAmount,
		TaxAmount:              taxAmount,
		ShippingCost:           shippingCost,
		Total:                  total,
		DeliveryEstimate:       deliveryEstimate,
		ShippingOptions:        shippingOptions,
		AppliedDiscount:        appliedDiscount,
		ItemsBreakdown:         itemsBreakdown,
		WarningsAndSuggestions: []string{},
	}, nil
}

// CreateOrder creates a new order from checkout data
func (s *CheckoutService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error) {
	s.logger.Debugf("creating order from checkout: customer_id=%d items_count=%d", req.CustomerID, len(req.CartItems))

	// Validate request
	if req.CustomerEmail == "" {
		return nil, errors.New("customer email is required")
	}
	if len(req.CartItems) == 0 {
		return nil, errors.New("cart items cannot be empty")
	}

	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d-%d", time.Now().Unix(), req.CustomerID)

	// TODO: Create order in database using queries
	// For now, placeholder implementation

	return &CreateOrderResponse{
		OrderID:         1,
		OrderNumber:     orderNumber,
		PaymentIntentID: "pi_test123",
		ClientSecret:    "secret_test123",
		Status:          "pending",
		Total:           0,
		NextAction:      "redirect_to_payment",
	}, nil
}

// ValidateShippingAddress validates address format and completeness
func (s *CheckoutService) ValidateShippingAddress(ctx context.Context, addr *AddressInput) error {
	if addr.RecipientName == "" {
		return errors.New("recipient name is required")
	}
	if addr.StreetAddress == "" {
		return errors.New("street address is required")
	}
	if addr.City == "" {
		return errors.New("city is required")
	}
	if addr.StateProvince == "" {
		return errors.New("state/province is required")
	}
	if addr.CountryCode == "" {
		return errors.New("country code is required")
	}
	if addr.PhoneNumber == "" {
		return errors.New("phone number is required")
	}

	s.logger.Debugf("address validation passed: city=%s", addr.City)
	return nil
}

// ApplyDiscountCode validates and applies a discount code
func (s *CheckoutService) ApplyDiscountCode(ctx context.Context, code string, subtotal int64) (*AppliedDiscount, error) {
	if code == "" {
		return nil, errors.New("discount code is required")
	}

	// TODO: Query discount codes table via SQLC
	s.logger.Debugf("applying discount code: code=%s", code)

	// Placeholder implementation
	return &AppliedDiscount{
		Code:           code,
		DiscountType:   "percentage",
		DiscountValue:  10,
		DiscountAmount: subtotal / 10,
	}, nil
}

package commerce

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

const (
	EventTypeLivestreamCommercePinUpdated = "livestream.commerce.pin.updated"
	EventTypeLivestreamCartUpdated        = "livestream.commerce.cart.updated"
	EventTypeOrderCreated                 = "order.created"
	EventTypeLivestreamPurchaseCompleted  = "livestream.purchase.completed"
	EventTypeLivestreamCommissionSettled  = "livestream.commission.settled"
)

type Service struct {
	mu                   sync.RWMutex
	pinnedProducts       map[string]map[int64]*ProductPin
	promotions           map[string]map[string]*LivePromotion
	carts                map[string]map[int64]*LiveCart
	stock                map[int64]int
	reservedStock        map[int64]int
	purchaseHistory      []*PurchaseReceipt
	productFetcher       ProductFetcher
	orderCreator         OrderCreator
	paymentProcessor     PaymentProcessor
	walletManager        WalletManager
	eventDispatcher      EventDispatcher
	commissionCalculator CommissionCalculator
	logger               *logger.Logger
}

func NewService(options ...Option) *Service {
	service := &Service{
		pinnedProducts:       make(map[string]map[int64]*ProductPin),
		promotions:           make(map[string]map[string]*LivePromotion),
		carts:                make(map[string]map[int64]*LiveCart),
		stock:                make(map[int64]int),
		reservedStock:        make(map[int64]int),
		purchaseHistory:      make([]*PurchaseReceipt, 0),
		logger:               logger.NewNoop(),
		commissionCalculator: DefaultCommissionCalculator{},
		paymentProcessor:     NoopPaymentProcessor{},
		walletManager:        NoopWalletManager{},
		eventDispatcher:      NoopEventDispatcher{},
	}
	for _, option := range options {
		option(service)
	}
	return service
}

type Option func(*Service)

func WithLogger(log *logger.Logger) Option {
	return func(s *Service) {
		s.logger = log
	}
}

func WithProductFetcher(fetcher ProductFetcher) Option {
	return func(s *Service) {
		s.productFetcher = fetcher
	}
}

func WithOrderCreator(creator OrderCreator) Option {
	return func(s *Service) {
		s.orderCreator = creator
	}
}

func WithPaymentProcessor(processor PaymentProcessor) Option {
	return func(s *Service) {
		s.paymentProcessor = processor
	}
}

func WithWalletManager(wallet WalletManager) Option {
	return func(s *Service) {
		s.walletManager = wallet
	}
}

func WithEventDispatcher(dispatcher EventDispatcher) Option {
	return func(s *Service) {
		s.eventDispatcher = dispatcher
	}
}

func WithCommissionCalculator(calculator CommissionCalculator) Option {
	return func(s *Service) {
		s.commissionCalculator = calculator
	}
}

func (s *Service) PinProduct(streamID string, productID int64, displayTitle, overlayText, promotionID string, displayOrder int) (*ProductPin, error) {
	if streamID == "" {
		return nil, fmt.Errorf("stream ID is required")
	}
	if productID == 0 {
		return nil, fmt.Errorf("product ID is required")
	}

	pin := &ProductPin{
		StreamID:     streamID,
		ProductID:    productID,
		DisplayTitle: displayTitle,
		OverlayText:  overlayText,
		PromotionID:  promotionID,
		DisplayOrder: displayOrder,
		PinnedAt:     time.Now(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.pinnedProducts[streamID] == nil {
		s.pinnedProducts[streamID] = make(map[int64]*ProductPin)
	}
	s.pinnedProducts[streamID][productID] = pin

	s.publishEvent(EventTypeLivestreamCommercePinUpdated, map[string]interface{}{
		"stream_id":     streamID,
		"product_id":    productID,
		"display_title": displayTitle,
		"promotion_id":  promotionID,
	})

	return pin, nil
}

func (s *Service) UnpinProduct(streamID string, productID int64) error {
	if streamID == "" {
		return fmt.Errorf("stream ID is required")
	}
	if productID == 0 {
		return fmt.Errorf("product ID is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	pins, ok := s.pinnedProducts[streamID]
	if !ok {
		return fmt.Errorf("no pinned products for stream %q", streamID)
	}
	pin, ok := pins[productID]
	if !ok {
		return fmt.Errorf("product %d is not pinned", productID)
	}
	pin.UnpinnedAt = time.Now()
	delete(pins, productID)
	return nil
}

func (s *Service) ListPins(streamID string) []ProductPin {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pins := make([]ProductPin, 0)
	for _, pin := range s.pinnedProducts[streamID] {
		pins = append(pins, *pin)
	}
	sort.Slice(pins, func(i, j int) bool {
		return pins[i].DisplayOrder < pins[j].DisplayOrder
	})
	return pins
}

func (s *Service) CreatePromotion(streamID string, promo LivePromotion) (*LivePromotion, error) {
	if streamID == "" {
		return nil, fmt.Errorf("stream ID is required")
	}
	if promo.ID == "" {
		promo.ID = fmt.Sprintf("promo-%d", time.Now().UnixNano())
	}
	if promo.Label == "" {
		return nil, fmt.Errorf("promotion label is required")
	}
	if promo.DiscountPercent <= 0 || promo.DiscountPercent > 100 {
		return nil, fmt.Errorf("discount percent must be greater than zero and at most 100")
	}
	if promo.StartAt.IsZero() {
		promo.StartAt = time.Now()
	}
	if promo.EndAt.IsZero() {
		promo.EndAt = promo.StartAt.Add(15 * time.Minute)
	}
	promo.Active = time.Now().After(promo.StartAt) && time.Now().Before(promo.EndAt)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.promotions[streamID] == nil {
		s.promotions[streamID] = make(map[string]*LivePromotion)
	}
	s.promotions[streamID][promo.ID] = &promo

	return &promo, nil
}

func (s *Service) ListPromotions(streamID string) []LivePromotion {
	s.mu.RLock()
	defer s.mu.RUnlock()

	promos := make([]LivePromotion, 0)
	for _, promo := range s.promotions[streamID] {
		promos = append(promos, *promo)
	}
	return promos
}

func (s *Service) SyncStock(productID int64, available int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stock[productID] = available
}

func (s *Service) GetStock(productID int64) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.stock[productID]
}

func (s *Service) AddToCart(streamID string, userID int64, productID int64, quantity int) (*LiveCart, error) {
	if streamID == "" {
		return nil, fmt.Errorf("stream ID is required")
	}
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if productID == 0 {
		return nil, fmt.Errorf("product ID is required")
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than zero")
	}
	if s.productFetcher == nil {
		return nil, fmt.Errorf("product fetcher is not configured")
	}

	product, err := s.productFetcher.FetchProduct(productID)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	available := s.stock[productID]
	reserved := s.reservedStock[productID]
	if available > 0 && reserved+quantity > available {
		return nil, fmt.Errorf("not enough stock for product %d", productID)
	}

	if s.carts[streamID] == nil {
		s.carts[streamID] = make(map[int64]*LiveCart)
	}

	cart := s.carts[streamID][userID]
	if cart == nil {
		cart = &LiveCart{StreamID: streamID, UserID: userID, UpdatedAt: time.Now()}
		s.carts[streamID][userID] = cart
	}

	found := false
	for i := range cart.Items {
		if cart.Items[i].ProductID == productID {
			cart.Items[i].Quantity += quantity
			cart.Items[i].TotalPrice = float64(cart.Items[i].Quantity) * cart.Items[i].UnitPrice
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, LiveCartItem{
			ProductID:  product.ID,
			SKU:        product.SKU,
			Name:       product.Name,
			Quantity:   quantity,
			UnitPrice:  product.Price,
			TotalPrice: float64(quantity) * product.Price,
		})
	}

	s.reservedStock[productID] += quantity
	cart.UpdatedAt = time.Now()
	s.calculateCartTotals(cart)

	s.publishEvent(EventTypeLivestreamCartUpdated, map[string]interface{}{
		"stream_id":  streamID,
		"user_id":    userID,
		"product_id": productID,
		"quantity":   quantity,
	})

	return cart, nil
}

func (s *Service) GetCart(streamID string, userID int64) (*LiveCart, error) {
	if streamID == "" {
		return nil, fmt.Errorf("stream ID is required")
	}
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	cart := s.carts[streamID][userID]
	if cart == nil {
		return &LiveCart{StreamID: streamID, UserID: userID, UpdatedAt: time.Now()}, nil
	}
	return cart, nil
}

func (s *Service) ApplyPromotion(streamID string, userID int64, promotionID string) (*LiveCart, error) {
	if streamID == "" {
		return nil, fmt.Errorf("stream ID is required")
	}
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if promotionID == "" {
		return nil, fmt.Errorf("promotion ID is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	cart := s.carts[streamID][userID]
	if cart == nil || len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	promo, ok := s.promotions[streamID][promotionID]
	if !ok {
		return nil, fmt.Errorf("promotion %q not found", promotionID)
	}
	if !promo.Active || time.Now().Before(promo.StartAt) || time.Now().After(promo.EndAt) {
		return nil, fmt.Errorf("promotion %q is not active", promotionID)
	}

	cart.Promotions = []LivePromotion{*promo}
	cart.UpdatedAt = time.Now()
	s.calculateCartTotals(cart)

	return cart, nil
}

func (s *Service) PurchaseCart(streamID string, userID int64, affiliateID *int64, paymentMethod string) (*PurchaseReceipt, error) {
	if streamID == "" {
		return nil, fmt.Errorf("stream ID is required")
	}
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if paymentMethod == "" {
		paymentMethod = "auto"
	}

	s.mu.Lock()
	cart := s.carts[streamID][userID]
	if cart == nil || len(cart.Items) == 0 {
		s.mu.Unlock()
		return nil, fmt.Errorf("cart is empty")
	}
	amount := cart.FinalAmount
	s.mu.Unlock()

	payment, err := s.paymentProcessor.Charge(userID, amount, paymentMethod)
	if err != nil {
		return nil, fmt.Errorf("payment failed: %w", err)
	}

	order, err := s.orderCreator.CreateOrder(&OrderInput{UserID: userID, StreamID: streamID, TotalAmount: amount, Status: "completed"})
	if err != nil {
		return nil, fmt.Errorf("order creation failed: %w", err)
	}

	split := s.commissionCalculator.Calculate(amount, userID, affiliateID)
	if err := s.walletManager.RecordPayout(userID, split.Creator, "creator_cut"); err != nil {
		s.logger.Errorf("failed to record creator payout: %v", err)
	}
	if affiliateID != nil {
		if err := s.walletManager.RecordPayout(*affiliateID, split.Affiliate, "affiliate_cut"); err != nil {
			s.logger.Errorf("failed to record affiliate payout: %v", err)
		}
	}
	if err := s.walletManager.RecordPayout(0, split.Platform, "platform_cut"); err != nil {
		s.logger.Errorf("failed to record platform payout: %v", err)
	}

	receipt := &PurchaseReceipt{
		PurchaseID:       fmt.Sprintf("purchase-%d", time.Now().UnixNano()),
		StreamID:         streamID,
		UserID:           userID,
		OrderReference:   strconv.FormatInt(order.OrderID, 10),
		PaymentReference: payment.PaymentID,
		Amount:           amount,
		Attribution: PurchaseAttribution{
			CreatorID:   userID,
			AffiliateID: affiliateID,
			Split:       split,
		},
		CreatedAt: time.Now(),
	}

	s.mu.Lock()
	s.purchaseHistory = append(s.purchaseHistory, receipt)
	for _, item := range cart.Items {
		s.reservedStock[item.ProductID] -= item.Quantity
		if s.stock[item.ProductID] > 0 {
			s.stock[item.ProductID] -= item.Quantity
		}
	}
	cart.Items = nil
	cart.Promotions = nil
	cart.TotalAmount = 0
	cart.Discount = 0
	cart.FinalAmount = 0
	cart.UpdatedAt = time.Now()
	s.mu.Unlock()

	s.publishEvent(EventTypeOrderCreated, map[string]interface{}{
		"stream_id":         streamID,
		"user_id":           userID,
		"order_reference":   receipt.OrderReference,
		"payment_reference": receipt.PaymentReference,
		"amount":            amount,
		"affiliate_id":      affiliateID,
		"creator_cut":       split.Creator,
		"affiliate_cut":     split.Affiliate,
		"platform_cut":      split.Platform,
	})

	s.publishEvent(EventTypeLivestreamPurchaseCompleted, map[string]interface{}{
		"stream_id":       streamID,
		"user_id":         userID,
		"purchase_id":     receipt.PurchaseID,
		"amount":          amount,
		"payment_method":  paymentMethod,
		"order_reference": receipt.OrderReference,
	})

	s.publishEvent(EventTypeLivestreamCommissionSettled, map[string]interface{}{
		"stream_id":     streamID,
		"purchase_id":   receipt.PurchaseID,
		"creator_cut":   split.Creator,
		"affiliate_cut": split.Affiliate,
		"platform_cut":  split.Platform,
		"seller_payout": split.Seller,
	})

	return receipt, nil
}

func (s *Service) calculateCartTotals(cart *LiveCart) {
	total := 0.0
	for _, item := range cart.Items {
		total += item.TotalPrice
	}

	discount := 0.0
	for _, promo := range cart.Promotions {
		if promo.Active && time.Now().After(promo.StartAt) && time.Now().Before(promo.EndAt) {
			discount += total * (promo.DiscountPercent / 100)
		}
	}

	cart.TotalAmount = total
	cart.Discount = discount
	cart.FinalAmount = total - discount
	if cart.FinalAmount < 0 {
		cart.FinalAmount = 0
	}
}

func (s *Service) ListPurchases(streamID string) []PurchaseReceipt {
	s.mu.RLock()
	defer s.mu.RUnlock()

	receipts := make([]PurchaseReceipt, 0)
	for _, receipt := range s.purchaseHistory {
		if receipt.StreamID == streamID {
			receipts = append(receipts, *receipt)
		}
	}
	return receipts
}

func (s *Service) publishEvent(eventType string, payload map[string]interface{}) {
	if s.eventDispatcher == nil {
		return
	}
	_ = s.eventDispatcher.PublishEvent(eventType, payload)
}

type DefaultCommissionCalculator struct{}

func (DefaultCommissionCalculator) Calculate(amount float64, creatorID int64, affiliateID *int64) CommissionSplit {
	platform := round(amount * 0.15)
	creator := round(amount * 0.70)
	affiliate := 0.0
	if affiliateID != nil {
		affiliate = round(amount * 0.10)
	}
	seller := amount - platform - creator - affiliate
	return CommissionSplit{
		Platform:  platform,
		Creator:   creator,
		Affiliate: affiliate,
		Seller:    seller,
	}
}

func round(value float64) float64 {
	rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return rounded
}

type NoopPaymentProcessor struct{}

func (NoopPaymentProcessor) Charge(userID int64, amount float64, method string) (*PaymentResult, error) {
	return &PaymentResult{PaymentID: fmt.Sprintf("payment-%d", time.Now().UnixNano()), Method: method}, nil
}

type NoopWalletManager struct{}

func (NoopWalletManager) RecordPayout(userID int64, amount float64, reason string) error {
	return nil
}

type NoopEventDispatcher struct{}

func (NoopEventDispatcher) PublishEvent(eventType string, payload map[string]interface{}) error {
	return nil
}

package merchant

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/teamart/commerce-api/pkg/logger"
)

type Service struct {
	repo Repository
	log  *logger.Logger
}

func NewService(repo Repository, log *logger.Logger) *Service {
	return &Service{repo: repo, log: log}
}

type CreateMerchantInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	BillingPlan string `json:"billing_plan,omitempty"`
	Currency    string `json:"currency,omitempty"`
}

type CreateStoreInput struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Category    string          `json:"category,omitempty"`
	BannerURL   string          `json:"banner_url,omitempty"`
	Settings    json.RawMessage `json:"settings,omitempty"`
	CreatorID   *int64          `json:"creator_id,omitempty"`
}

type StoreUpdateInput struct {
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Category    string          `json:"category,omitempty"`
	BannerURL   string          `json:"banner_url,omitempty"`
	Settings    json.RawMessage `json:"settings,omitempty"`
}

func (s *Service) CreateMerchant(ctx context.Context, ownerID int64, input *CreateMerchantInput) (*Merchant, error) {
	if ownerID == 0 {
		return nil, fmt.Errorf("owner ID is required")
	}
	if input.Name == "" {
		return nil, fmt.Errorf("merchant name is required")
	}

	merchant := Merchant{
		OwnerID:     ownerID,
		Name:        input.Name,
		Slug:        generateSlug(input.Name),
		Description: input.Description,
		Status:      MerchantStatusPending,
		BillingPlan: defaultString(input.BillingPlan, "starter"),
		Currency:    defaultString(input.Currency, "USD"),
	}

	created, err := s.repo.CreateMerchant(ctx, merchant)
	if err != nil {
		s.log.Errorf("failed to create merchant: %v", err)
		return nil, err
	}

	s.log.Infof("merchant created: %d for owner %d", created.ID, ownerID)
	return created, nil
}

func (s *Service) GetMerchant(ctx context.Context, merchantID int64) (*Merchant, error) {
	if merchantID == 0 {
		return nil, fmt.Errorf("merchant ID is required")
	}
	merchant, err := s.repo.GetMerchantByID(ctx, merchantID)
	if err != nil {
		s.log.Errorf("failed to fetch merchant %d: %v", merchantID, err)
		return nil, err
	}
	return merchant, nil
}

func (s *Service) GetMerchantByOwner(ctx context.Context, ownerID int64) (*Merchant, error) {
	if ownerID == 0 {
		return nil, fmt.Errorf("owner ID is required")
	}
	merchant, err := s.repo.GetMerchantByOwnerID(ctx, ownerID)
	if err != nil {
		s.log.Errorf("failed to resolve merchant for owner %d: %v", ownerID, err)
		return nil, err
	}
	return merchant, nil
}

func (s *Service) CreateStore(ctx context.Context, merchantID, ownerID int64, input *CreateStoreInput) (*Store, error) {
	if merchantID == 0 {
		return nil, fmt.Errorf("merchant ID is required")
	}
	if ownerID == 0 {
		return nil, fmt.Errorf("owner ID is required")
	}
	if input.Name == "" {
		return nil, fmt.Errorf("store name is required")
	}

	store := Store{
		MerchantID:    merchantID,
		OwnerID:       ownerID,
		CreatorID:     input.CreatorID,
		Name:          input.Name,
		Slug:          generateSlug(input.Name),
		Description:   input.Description,
		Category:      input.Category,
		Settings:      input.Settings,
		Status:        StoreStatusActive,
		StorefrontURL: "",
	}

	created, err := s.repo.CreateStore(ctx, store)
	if err != nil {
		s.log.Errorf("failed to create store for merchant %d: %v", merchantID, err)
		return nil, err
	}

	s.log.Infof("store created: %d for merchant %d", created.ID, merchantID)
	return created, nil
}

func (s *Service) GetStore(ctx context.Context, storeID int64) (*Store, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("store ID is required")
	}
	store, err := s.repo.GetStoreByID(ctx, storeID)
	if err != nil {
		s.log.Errorf("failed to fetch store %d: %v", storeID, err)
		return nil, err
	}
	return store, nil
}

func (s *Service) ListStoresForMerchant(ctx context.Context, merchantID int64) ([]Store, error) {
	if merchantID == 0 {
		return nil, fmt.Errorf("merchant ID is required")
	}
	stores, err := s.repo.ListStoresByMerchantID(ctx, merchantID)
	if err != nil {
		s.log.Errorf("failed to list stores for merchant %d: %v", merchantID, err)
		return nil, err
	}
	return stores, nil
}

func (s *Service) UpdateStore(ctx context.Context, storeID int64, input *StoreUpdateInput) (*Store, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("store ID is required")
	}

	store, err := s.repo.GetStoreByID(ctx, storeID)
	if err != nil {
		s.log.Errorf("failed to resolve store %d before update: %v", storeID, err)
		return nil, err
	}

	if input.Name != "" {
		store.Name = input.Name
		store.Slug = generateSlug(input.Name)
	}
	if input.Description != "" {
		store.Description = input.Description
	}
	if input.Category != "" {
		store.Category = input.Category
	}
	if len(input.Settings) > 0 {
		store.Settings = input.Settings
	}

	if err := s.repo.UpdateStore(ctx, *store); err != nil {
		s.log.Errorf("failed to update store %d: %v", storeID, err)
		return nil, err
	}

	return store, nil
}

func (s *Service) CreateMerchantKYC(ctx context.Context, merchantID int64, legalName, taxID, businessType string) (*MerchantKYC, error) {
	if merchantID == 0 {
		return nil, fmt.Errorf("merchant ID is required")
	}
	if legalName == "" {
		return nil, fmt.Errorf("legal name is required")
	}

	kyc := MerchantKYC{
		MerchantID:   merchantID,
		LegalName:    legalName,
		TaxID:        taxID,
		BusinessType: businessType,
		Status:       MerchantKYCStatusPending,
	}

	return s.repo.CreateMerchantKYC(ctx, kyc)
}

func (s *Service) CreateMerchantPayoutAccount(ctx context.Context, merchantID int64, provider, accountHolderName, accountType, externalAccountID, currency string, metadata json.RawMessage) (*MerchantPayoutAccount, error) {
	if merchantID == 0 {
		return nil, fmt.Errorf("merchant ID is required")
	}
	if provider == "" || accountHolderName == "" || accountType == "" || externalAccountID == "" {
		return nil, fmt.Errorf("payout provider, holder name, account type, and external account ID are required")
	}

	account := MerchantPayoutAccount{
		MerchantID:        merchantID,
		Provider:          provider,
		AccountHolderName: accountHolderName,
		AccountType:       accountType,
		ExternalAccountID: externalAccountID,
		Currency:          defaultString(currency, "USD"),
		Status:            "pending",
		Metadata:          metadata,
	}

	return s.repo.CreateMerchantPayoutAccount(ctx, account)
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func generateSlug(value string) string {
	slug := strings.TrimSpace(strings.ToLower(value))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, "--", "-")
	return slug
}

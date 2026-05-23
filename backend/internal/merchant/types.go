package merchant

import (
	"encoding/json"
	"time"
)

type MerchantStatus string

type StoreStatus string

type MerchantKYCStatus string

const (
	MerchantStatusPending   MerchantStatus = "pending"
	MerchantStatusActive    MerchantStatus = "active"
	MerchantStatusSuspended MerchantStatus = "suspended"

	StoreStatusActive    StoreStatus = "active"
	StoreStatusInactive  StoreStatus = "inactive"
	StoreStatusSuspended StoreStatus = "suspended"

	MerchantKYCStatusPending  MerchantKYCStatus = "pending"
	MerchantKYCStatusApproved MerchantKYCStatus = "approved"
	MerchantKYCStatusRejected MerchantKYCStatus = "rejected"
)

type Merchant struct {
	ID          int64          `json:"id"`
	OwnerID     int64          `json:"owner_id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description,omitempty"`
	Status      MerchantStatus `json:"status"`
	BillingPlan string         `json:"billing_plan"`
	Currency    string         `json:"currency"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Store struct {
	ID            int64           `json:"id"`
	MerchantID    int64           `json:"merchant_id"`
	OwnerID       int64           `json:"owner_id"`
	CreatorID     *int64          `json:"creator_id,omitempty"`
	Name          string          `json:"name"`
	Slug          string          `json:"slug"`
	Description   string          `json:"description,omitempty"`
	Category      string          `json:"category,omitempty"`
	Settings      json.RawMessage `json:"settings,omitempty"`
	Status        StoreStatus     `json:"status"`
	StorefrontURL string          `json:"storefront_url,omitempty"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type MerchantKYC struct {
	ID           int64             `json:"id"`
	MerchantID   int64             `json:"merchant_id"`
	LegalName    string            `json:"legal_name"`
	TaxID        string            `json:"tax_id,omitempty"`
	BusinessType string            `json:"business_type,omitempty"`
	Status       MerchantKYCStatus `json:"status"`
	SubmittedAt  time.Time         `json:"submitted_at"`
	ReviewedAt   *time.Time        `json:"reviewed_at,omitempty"`
	ApprovedAt   *time.Time        `json:"approved_at,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type MerchantPayoutAccount struct {
	ID                int64           `json:"id"`
	MerchantID        int64           `json:"merchant_id"`
	Provider          string          `json:"provider"`
	AccountHolderName string          `json:"account_holder_name"`
	AccountType       string          `json:"account_type"`
	ExternalAccountID string          `json:"external_account_id"`
	Currency          string          `json:"currency"`
	Status            string          `json:"status"`
	Metadata          json.RawMessage `json:"metadata,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

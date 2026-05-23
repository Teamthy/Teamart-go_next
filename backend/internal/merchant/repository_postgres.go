package merchant

import (
	"context"

	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

type PostgresRepository struct {
	db  *database.Pool
	log *logger.Logger
}

func NewPostgresRepository(db *database.Pool, log *logger.Logger) *PostgresRepository {
	return &PostgresRepository{db: db, log: log}
}

func (r *PostgresRepository) CreateMerchant(ctx context.Context, merchant Merchant) (*Merchant, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO merchants (owner_id, name, slug, description, status, billing_plan, currency)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, owner_id, name, slug, description, status, billing_plan, currency, created_at, updated_at
	`, merchant.OwnerID, merchant.Name, merchant.Slug, merchant.Description, merchant.Status, merchant.BillingPlan, merchant.Currency)

	var created Merchant
	if err := row.Scan(&created.ID, &created.OwnerID, &created.Name, &created.Slug, &created.Description, &created.Status, &created.BillingPlan, &created.Currency, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresRepository) GetMerchantByID(ctx context.Context, merchantID int64) (*Merchant, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, owner_id, name, slug, description, status, billing_plan, currency, created_at, updated_at
		FROM merchants
		WHERE id = $1
	`, merchantID)

	var merchant Merchant
	if err := row.Scan(&merchant.ID, &merchant.OwnerID, &merchant.Name, &merchant.Slug, &merchant.Description, &merchant.Status, &merchant.BillingPlan, &merchant.Currency, &merchant.CreatedAt, &merchant.UpdatedAt); err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (r *PostgresRepository) GetMerchantByOwnerID(ctx context.Context, ownerID int64) (*Merchant, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, owner_id, name, slug, description, status, billing_plan, currency, created_at, updated_at
		FROM merchants
		WHERE owner_id = $1
		LIMIT 1
	`, ownerID)

	var merchant Merchant
	if err := row.Scan(&merchant.ID, &merchant.OwnerID, &merchant.Name, &merchant.Slug, &merchant.Description, &merchant.Status, &merchant.BillingPlan, &merchant.Currency, &merchant.CreatedAt, &merchant.UpdatedAt); err != nil {
		return nil, err
	}
	return &merchant, nil
}

func (r *PostgresRepository) CreateStore(ctx context.Context, store Store) (*Store, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO stores (merchant_id, owner_id, creator_id, name, slug, description, category, settings, status, storefront_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, merchant_id, owner_id, creator_id, name, slug, description, category, settings, status, storefront_url, created_at, updated_at
	`, store.MerchantID, store.OwnerID, store.CreatorID, store.Name, store.Slug, store.Description, store.Category, store.Settings, store.Status, store.StorefrontURL)

	var created Store
	if err := row.Scan(&created.ID, &created.MerchantID, &created.OwnerID, &created.CreatorID, &created.Name, &created.Slug, &created.Description, &created.Category, &created.Settings, &created.Status, &created.StorefrontURL, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresRepository) GetStoreByID(ctx context.Context, storeID int64) (*Store, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, merchant_id, owner_id, creator_id, name, slug, description, category, settings, status, storefront_url, created_at, updated_at
		FROM stores
		WHERE id = $1
	`, storeID)

	var store Store
	if err := row.Scan(&store.ID, &store.MerchantID, &store.OwnerID, &store.CreatorID, &store.Name, &store.Slug, &store.Description, &store.Category, &store.Settings, &store.Status, &store.StorefrontURL, &store.CreatedAt, &store.UpdatedAt); err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *PostgresRepository) ListStoresByMerchantID(ctx context.Context, merchantID int64) ([]Store, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, merchant_id, owner_id, creator_id, name, slug, description, category, settings, status, storefront_url, created_at, updated_at
		FROM stores
		WHERE merchant_id = $1
		ORDER BY created_at DESC
	`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []Store
	for rows.Next() {
		var store Store
		if err := rows.Scan(&store.ID, &store.MerchantID, &store.OwnerID, &store.CreatorID, &store.Name, &store.Slug, &store.Description, &store.Category, &store.Settings, &store.Status, &store.StorefrontURL, &store.CreatedAt, &store.UpdatedAt); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func (r *PostgresRepository) UpdateStore(ctx context.Context, store Store) error {
	_, err := r.db.Exec(ctx, `
		UPDATE stores
		SET name = $1,
		    slug = $2,
		    description = $3,
		    category = $4,
		    settings = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`, store.Name, store.Slug, store.Description, store.Category, store.Settings, store.ID)
	return err
}

func (r *PostgresRepository) CreateMerchantKYC(ctx context.Context, kyc MerchantKYC) (*MerchantKYC, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO merchant_kyc (merchant_id, legal_name, tax_id, business_type, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, merchant_id, legal_name, tax_id, business_type, status, submitted_at, reviewed_at, approved_at, created_at, updated_at
	`, kyc.MerchantID, kyc.LegalName, kyc.TaxID, kyc.BusinessType, kyc.Status)

	var created MerchantKYC
	if err := row.Scan(&created.ID, &created.MerchantID, &created.LegalName, &created.TaxID, &created.BusinessType, &created.Status, &created.SubmittedAt, &created.ReviewedAt, &created.ApprovedAt, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresRepository) CreateMerchantPayoutAccount(ctx context.Context, payout MerchantPayoutAccount) (*MerchantPayoutAccount, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO merchant_payout_accounts (merchant_id, provider, account_holder_name, account_type, external_account_id, currency, status, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, merchant_id, provider, account_holder_name, account_type, external_account_id, currency, status, metadata, created_at, updated_at
	`, payout.MerchantID, payout.Provider, payout.AccountHolderName, payout.AccountType, payout.ExternalAccountID, payout.Currency, payout.Status, payout.Metadata)

	var created MerchantPayoutAccount
	if err := row.Scan(&created.ID, &created.MerchantID, &created.Provider, &created.AccountHolderName, &created.AccountType, &created.ExternalAccountID, &created.Currency, &created.Status, &created.Metadata, &created.CreatedAt, &created.UpdatedAt); err != nil {
		return nil, err
	}
	return &created, nil
}

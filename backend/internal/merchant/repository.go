package merchant

import (
	"context"
)

type Repository interface {
	CreateMerchant(ctx context.Context, merchant Merchant) (*Merchant, error)
	GetMerchantByID(ctx context.Context, merchantID int64) (*Merchant, error)
	GetMerchantByOwnerID(ctx context.Context, ownerID int64) (*Merchant, error)
	CreateStore(ctx context.Context, store Store) (*Store, error)
	GetStoreByID(ctx context.Context, storeID int64) (*Store, error)
	ListStoresByMerchantID(ctx context.Context, merchantID int64) ([]Store, error)
	UpdateStore(ctx context.Context, store Store) error
	CreateMerchantKYC(ctx context.Context, kyc MerchantKYC) (*MerchantKYC, error)
	CreateMerchantPayoutAccount(ctx context.Context, payout MerchantPayoutAccount) (*MerchantPayoutAccount, error)
}

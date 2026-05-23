package staff

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/teamart/commerce-api/internal/infra/database"
	"github.com/teamart/commerce-api/pkg/logger"
)

type Service struct {
	db  *database.Pool
	log *logger.Logger
}

func NewService(db *database.Pool, log *logger.Logger) *Service {
	return &Service{db: db, log: log}
}

type StaffAccount struct {
	ID         int64     `json:"id"`
	MerchantID int64     `json:"merchant_id"`
	UserID     int64     `json:"user_id"`
	Role       string    `json:"role"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type StoreMember struct {
	ID             int64      `json:"id"`
	StoreID        int64      `json:"store_id"`
	StaffAccountID int64      `json:"staff_account_id"`
	Role           string     `json:"role"`
	Permissions    []string   `json:"permissions"`
	IsActive       bool       `json:"is_active"`
	JoinedAt       time.Time  `json:"joined_at"`
	LeftAt         *time.Time `json:"left_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (s *Service) CreateStaffAccount(ctx context.Context, merchantID, userID int64, role string) (*StaffAccount, error) {
	if merchantID == 0 {
		return nil, fmt.Errorf("merchant ID is required")
	}
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if role == "" {
		role = "staff"
	}

	row := s.db.QueryRow(ctx, `
		INSERT INTO staff_accounts (merchant_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (merchant_id, user_id)
		DO UPDATE SET role = EXCLUDED.role, is_active = TRUE, updated_at = CURRENT_TIMESTAMP
		RETURNING id, merchant_id, user_id, role, is_active, created_at, updated_at
	`, merchantID, userID, role)

	var account StaffAccount
	if err := row.Scan(&account.ID, &account.MerchantID, &account.UserID, &account.Role, &account.IsActive, &account.CreatedAt, &account.UpdatedAt); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Service) ListStaffForMerchant(ctx context.Context, merchantID int64) ([]StaffAccount, error) {
	if merchantID == 0 {
		return nil, fmt.Errorf("merchant ID is required")
	}

	rows, err := s.db.Query(ctx, `
		SELECT id, merchant_id, user_id, role, is_active, created_at, updated_at
		FROM staff_accounts
		WHERE merchant_id = $1
		ORDER BY created_at DESC
	`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var staff []StaffAccount
	for rows.Next() {
		var account StaffAccount
		if err := rows.Scan(&account.ID, &account.MerchantID, &account.UserID, &account.Role, &account.IsActive, &account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, err
		}
		staff = append(staff, account)
	}

	return staff, nil
}

func (s *Service) AddStoreMember(ctx context.Context, storeID, staffAccountID int64, role string, permissions []string) (*StoreMember, error) {
	if storeID == 0 || staffAccountID == 0 {
		return nil, fmt.Errorf("store ID and staff account ID are required")
	}
	if role == "" {
		role = "staff"
	}

	row := s.db.QueryRow(ctx, `
		INSERT INTO store_members (store_id, staff_account_id, role, permissions)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (store_id, staff_account_id)
		DO UPDATE SET role = EXCLUDED.role, permissions = EXCLUDED.permissions, is_active = TRUE, updated_at = CURRENT_TIMESTAMP
		RETURNING id, store_id, staff_account_id, role, permissions, is_active, joined_at, left_at, created_at, updated_at
	`, storeID, staffAccountID, role, permissions)

	var member StoreMember
	if err := row.Scan(&member.ID, &member.StoreID, &member.StaffAccountID, &member.Role, &member.Permissions, &member.IsActive, &member.JoinedAt, &member.LeftAt, &member.CreatedAt, &member.UpdatedAt); err != nil {
		return nil, err
	}

	return &member, nil
}

func (s *Service) ListStoreMembers(ctx context.Context, storeID int64) ([]StoreMember, error) {
	if storeID == 0 {
		return nil, fmt.Errorf("store ID is required")
	}

	rows, err := s.db.Query(ctx, `
		SELECT id, store_id, staff_account_id, role, permissions, is_active, joined_at, left_at, created_at, updated_at
		FROM store_members
		WHERE store_id = $1
	`, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []StoreMember
	for rows.Next() {
		var member StoreMember
		if err := rows.Scan(&member.ID, &member.StoreID, &member.StaffAccountID, &member.Role, &member.Permissions, &member.IsActive, &member.JoinedAt, &member.LeftAt, &member.CreatedAt, &member.UpdatedAt); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (s *Service) GetActiveStaffAccountByUserID(ctx context.Context, userID int64) (*StaffAccount, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	row := s.db.QueryRow(ctx, `
		SELECT id, merchant_id, user_id, role, is_active, created_at, updated_at
		FROM staff_accounts
		WHERE user_id = $1 AND is_active = TRUE
		LIMIT 1
	`, userID)

	var account StaffAccount
	if err := row.Scan(&account.ID, &account.MerchantID, &account.UserID, &account.Role, &account.IsActive, &account.CreatedAt, &account.UpdatedAt); err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *Service) ResolveTenantForUser(ctx context.Context, userID int64) (int64, error) {
	account, err := s.GetActiveStaffAccountByUserID(ctx, userID)
	if err == nil && account != nil {
		return account.MerchantID, nil
	}
	return 0, fmt.Errorf("tenant resolution failed for user %d", userID)
}

func normalizeRole(role string) string {
	role = strings.TrimSpace(strings.ToLower(role))
	if role == "" {
		return "staff"
	}
	return role
}

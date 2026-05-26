package users

import (
	"context"
	"fmt"

	"github.com/teamart/commerce-api/internal/infra/queries"
	"github.com/teamart/commerce-api/pkg/logger"
)

// Service provides business logic for user operations
// It wraps the SQLC-generated queries with service-level operations
type Service struct {
	queries *queries.Queries
	logger  *logger.Logger
}

// NewService creates a new user service
func NewService(queries *queries.Queries, logger *logger.Logger) *Service {
	return &Service{
		queries: queries,
		logger:  logger,
	}
}

// CreateUserInput represents the input for creating a user
type CreateUserInput struct {
	Email        string
	Name         string
	PasswordHash string
}

// CreateUserOutput represents the output after creating a user
type CreateUserOutput struct {
	ID           int64
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    string
	UpdatedAt    string
}

// CreateUser creates a new user with validation
func (s *Service) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	// Validate input
	if input.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if input.PasswordHash == "" {
		return nil, fmt.Errorf("password hash is required")
	}

	s.logger.Debugf("creating user with email: %s", input.Email)

	// Execute the query using SQLC
	user, err := s.queries.CreateUser(ctx, input.Email, input.Name, input.PasswordHash)
	if err != nil {
		s.logger.Errorf("failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Infof("user created successfully with ID: %d", user.ID)

	return &CreateUserOutput{
		ID:           int64(user.ID),
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}

// GetUserByIDInput represents the input for getting a user by ID
type GetUserByIDInput struct {
	UserID int64
}

// GetUserByIDOutput represents the output after getting a user
type GetUserByIDOutput struct {
	ID           int64
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    string
	UpdatedAt    string
}

// GetUserByID retrieves a user by their ID
func (s *Service) GetUserByID(ctx context.Context, input *GetUserByIDInput) (*GetUserByIDOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	s.logger.Debugf("fetching user with ID: %d", input.UserID)

	user, err := s.queries.GetUserByID(ctx, int32(input.UserID))
	if err != nil {
		s.logger.Errorf("failed to fetch user: %v", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &GetUserByIDOutput{
		ID:           int64(user.ID),
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}

// GetUserByEmailInput represents the input for getting a user by email
type GetUserByEmailInput struct {
	Email string
}

// GetUserByEmailOutput represents the output after getting a user
type GetUserByEmailOutput struct {
	ID           int64
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    string
	UpdatedAt    string
}

// GetUserByEmail retrieves a user by their email
func (s *Service) GetUserByEmail(ctx context.Context, input *GetUserByEmailInput) (*GetUserByEmailOutput, error) {
	if input.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	s.logger.Debugf("fetching user with email: %s", input.Email)

	user, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err != nil {
		s.logger.Errorf("failed to fetch user by email: %v", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &GetUserByEmailOutput{
		ID:           int64(user.ID),
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}

// ListUsersInput represents the input for listing users
type ListUsersInput struct {
	Limit  int32
	Offset int32
}

// ListUsersOutput represents the output after listing users
type ListUsersOutput struct {
	Users  []UserData
	Total  int32
	Limit  int32
	Offset int32
}

type UserData struct {
	ID           int64
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    string
	UpdatedAt    string
}

// ListUsers retrieves a list of users with pagination
func (s *Service) ListUsers(ctx context.Context, input *ListUsersInput) (*ListUsersOutput, error) {
	if input.Limit == 0 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	s.logger.Debugf("listing users with limit: %d, offset: %d", input.Limit, input.Offset)

	users, err := s.queries.ListUsers(ctx, input.Limit, input.Offset)
	if err != nil {
		s.logger.Errorf("failed to list users: %v", err)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	output := &ListUsersOutput{
		Users:  make([]UserData, len(users)),
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	for i, user := range users {
		output.Users[i] = UserData{
			ID:           int64(user.ID),
			Email:        user.Email,
			Name:         user.Name,
			PasswordHash: user.PasswordHash,
			CreatedAt:    user.CreatedAt.String(),
			UpdatedAt:    user.UpdatedAt.String(),
		}
	}

	s.logger.Infof("fetched %d users", len(users))

	return output, nil
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	UserID int64
	Name   string
}

// UpdateUserOutput represents the output after updating a user
type UpdateUserOutput struct {
	ID           int64
	Email        string
	Name         string
	PasswordHash string
	CreatedAt    string
	UpdatedAt    string
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(ctx context.Context, input *UpdateUserInput) (*UpdateUserOutput, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	s.logger.Debugf("updating user with ID: %d", input.UserID)

	user, err := s.queries.UpdateUser(ctx, int32(input.UserID), input.Name)
	if err != nil {
		s.logger.Errorf("failed to update user: %v", err)
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	s.logger.Infof("user updated successfully with ID: %d", user.ID)

	return &UpdateUserOutput{
		ID:           int64(user.ID),
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
	}, nil
}

// DeleteUserInput represents the input for deleting a user
type DeleteUserInput struct {
	UserID int64
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, input *DeleteUserInput) error {
	if input.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}

	s.logger.Debugf("deleting user with ID: %d", input.UserID)

	err := s.queries.DeleteUser(ctx, int32(input.UserID))
	if err != nil {
		s.logger.Errorf("failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	s.logger.Infof("user deleted successfully with ID: %d", input.UserID)

	return nil
}

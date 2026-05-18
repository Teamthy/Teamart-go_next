package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/teamart/commerce-api/pkg/logger"
)

// OnboardingStateMachine manages user onboarding lifecycle
type OnboardingStateMachine struct {
	logger *logger.Logger
}

// NewOnboardingStateMachine creates a new onboarding state machine
func NewOnboardingStateMachine(logger *logger.Logger) *OnboardingStateMachine {
	return &OnboardingStateMachine{
		logger: logger,
	}
}

// TransitionInput represents input for state transition
type TransitionInput struct {
	UserID      int64
	CurrentState OnboardingState
	TargetState  OnboardingState
	Reason       string
}

// TransitionOutput represents the result of state transition
type TransitionOutput struct {
	UserID    int64
	PreviousState OnboardingState
	NewState  OnboardingState
	TransitionedAt time.Time
}

// Transition performs a state transition with validation
func (sm *OnboardingStateMachine) Transition(ctx context.Context, input *TransitionInput) (*TransitionOutput, error) {
	// Validate state transition is allowed
	if !ValidStateTransition(input.CurrentState, input.TargetState) {
		sm.logger.Warnf("invalid state transition for user %d: %s -> %s (reason: %s)",
			input.UserID, input.CurrentState, input.TargetState, input.Reason)
		return nil, fmt.Errorf("invalid state transition from %s to %s", input.CurrentState, input.TargetState)
	}

	sm.logger.Infof("transitioning user %d from %s to %s (reason: %s)",
		input.UserID, input.CurrentState, input.TargetState, input.Reason)

	return &TransitionOutput{
		UserID:    input.UserID,
		PreviousState: input.CurrentState,
		NewState:  input.TargetState,
		TransitionedAt: time.Now(),
	}, nil
}

// ===== Onboarding Phase Checklist =====

// OnboardingPhase represents a phase in the onboarding process
type OnboardingPhase string

const (
	// PhaseEmailVerification: Verify email address
	PhaseEmailVerification OnboardingPhase = "email_verification"
	
	// PhaseProfileCompletion: Complete user profile
	PhaseProfileCompletion OnboardingPhase = "profile_completion"
	
	// Phase2FA: Set up two-factor authentication (optional)
	Phase2FA OnboardingPhase = "two_factor_auth"
	
	// PhaseTermsAcceptance: Accept terms and conditions
	PhaseTermsAcceptance OnboardingPhase = "terms_acceptance"
)

// OnboardingProgress represents user's onboarding progress
type OnboardingProgress struct {
	UserID           int64
	CompletedPhases  []OnboardingPhase
	CurrentPhase     OnboardingPhase
	AllPhasesComplete bool
	ProgressPercentage int32 // 0-100
	LastUpdatedAt    time.Time
}

// GetProgress calculates onboarding progress
func (sm *OnboardingStateMachine) GetProgress(ctx context.Context, userID int64, currentState OnboardingState) *OnboardingProgress {
	progress := &OnboardingProgress{
		UserID:        userID,
		CompletedPhases: make([]OnboardingPhase, 0),
		LastUpdatedAt: time.Now(),
	}

	switch currentState {
	case StateNew:
		progress.CurrentPhase = PhaseEmailVerification
		progress.ProgressPercentage = 0
	case StateEmailVerified:
		progress.CompletedPhases = append(progress.CompletedPhases, PhaseEmailVerification)
		progress.CurrentPhase = PhaseProfileCompletion
		progress.ProgressPercentage = 25
	case StateProfileComplete:
		progress.CompletedPhases = append(progress.CompletedPhases, 
			PhaseEmailVerification, PhaseProfileCompletion)
		progress.CurrentPhase = Phase2FA
		progress.ProgressPercentage = 50
	case StateOnboarded:
		progress.CompletedPhases = append(progress.CompletedPhases,
			PhaseEmailVerification, PhaseProfileCompletion, 
			Phase2FA, PhaseTermsAcceptance)
		progress.AllPhasesComplete = true
		progress.ProgressPercentage = 100
	}

	return progress
}

// ===== Onboarding Events =====

// OnboardingEvent represents an event in the onboarding process
type OnboardingEvent struct {
	EventType  string // "email_verified", "profile_completed", etc.
	UserID     int64
	Timestamp  time.Time
	Metadata   map[string]interface{}
}

// ValidateEmailVerification validates email verification
func (sm *OnboardingStateMachine) ValidateEmailVerification(ctx context.Context, userID int64, email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	sm.logger.Infof("validating email verification for user %d: %s", userID, email)
	return nil
}

// ValidateProfileCompletion validates profile completion
func (sm *OnboardingStateMachine) ValidateProfileCompletion(ctx context.Context, userID int64, name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	sm.logger.Infof("validating profile completion for user %d: %s", userID, name)
	return nil
}

// IsOnboardingComplete checks if user has completed onboarding
func (sm *OnboardingStateMachine) IsOnboardingComplete(state OnboardingState) bool {
	return state == StateOnboarded
}

// GetRequiredSteps returns the required steps for onboarding
func (sm *OnboardingStateMachine) GetRequiredSteps() []OnboardingPhase {
	return []OnboardingPhase{
		PhaseEmailVerification,
		PhaseProfileCompletion,
		Phase2FA,
		PhaseTermsAcceptance,
	}
}

// CanSkipPhase determines if a phase can be skipped
func (sm *OnboardingStateMachine) CanSkipPhase(phase OnboardingPhase) bool {
	// Only 2FA is optional
	return phase == Phase2FA
}

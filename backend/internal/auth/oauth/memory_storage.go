package oauth

import (
	"context"
	"sync"
	"time"
)

// MemoryStateStorage stores OAuth state tokens in memory.
type MemoryStateStorage struct {
	mu     sync.RWMutex
	states map[string]*OAuthState
}

// NewMemoryStateStorage creates an in-memory OAuth state store.
func NewMemoryStateStorage() *MemoryStateStorage {
	return &MemoryStateStorage{states: make(map[string]*OAuthState)}
}

// SaveState stores the given OAuth state token.
func (s *MemoryStateStorage) SaveState(ctx context.Context, state *OAuthState) error {
	if state == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state.StateToken] = state
	return nil
}

// GetState returns the OAuth state for the given token.
func (s *MemoryStateStorage) GetState(ctx context.Context, stateToken string) (*OAuthState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.states[stateToken]
	if !ok {
		return nil, nil
	}
	return state, nil
}

// DeleteState removes the OAuth state token.
func (s *MemoryStateStorage) DeleteState(ctx context.Context, stateToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, stateToken)
	return nil
}

// IsStateExpired reports whether the state token has expired.
func (s *MemoryStateStorage) IsStateExpired(ctx context.Context, stateToken string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.states[stateToken]
	if !ok || state == nil {
		return true
	}
	return time.Now().After(state.ExpiresAt)
}

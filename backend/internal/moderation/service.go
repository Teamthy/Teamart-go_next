package moderation

import (
	"strings"
	"sync"
	"time"
)

type ModerationService struct {
	mu             sync.RWMutex
	moderators     map[int64]struct{}
	blockedUsers   map[int64]struct{}
	shadowBanned   map[int64]struct{}
	mutedUsers     map[int64]time.Time
	permissions    map[int64]map[string]struct{}
	messageHistory map[int64][]messageRecord
	bannedWords    []string
	floodWindow    time.Duration
	maxMessages    int
	repeatWindow   time.Duration
	repeatLimit    int
}

type messageRecord struct {
	text      string
	timestamp time.Time
}

type ModerationDecision struct {
	Allowed      bool       `json:"allowed"`
	ShadowBan    bool       `json:"shadow_ban"`
	RejectReason string     `json:"reject_reason,omitempty"`
	MuteExpires  *time.Time `json:"mute_expires,omitempty"`
}

type UserModerationStatus struct {
	UserID       int64    `json:"user_id"`
	IsModerator  bool     `json:"is_moderator"`
	Blocked      bool     `json:"blocked"`
	ShadowBanned bool     `json:"shadow_banned"`
	Muted        bool     `json:"muted"`
	MuteExpires  string   `json:"mute_expires,omitempty"`
	Permissions  []string `json:"permissions"`
}

type ModerationConfig struct {
	BannedWords  []string
	FloodWindow  time.Duration
	RepeatWindow time.Duration
	MaxMessages  int
	RepeatLimit  int
}

func NewService(cfg *ModerationConfig) *ModerationService {
	service := &ModerationService{
		moderators:     make(map[int64]struct{}),
		blockedUsers:   make(map[int64]struct{}),
		shadowBanned:   make(map[int64]struct{}),
		mutedUsers:     make(map[int64]time.Time),
		permissions:    make(map[int64]map[string]struct{}),
		messageHistory: make(map[int64][]messageRecord),
		bannedWords:    defaultBannedWords(),
		floodWindow:    time.Minute,
		maxMessages:    8,
		repeatWindow:   30 * time.Second,
		repeatLimit:    2,
	}

	if cfg != nil {
		if len(cfg.BannedWords) > 0 {
			service.bannedWords = cfg.BannedWords
		}
		if cfg.FloodWindow > 0 {
			service.floodWindow = cfg.FloodWindow
		}
		if cfg.MaxMessages > 0 {
			service.maxMessages = cfg.MaxMessages
		}
		if cfg.RepeatWindow > 0 {
			service.repeatWindow = cfg.RepeatWindow
		}
		if cfg.RepeatLimit > 0 {
			service.repeatLimit = cfg.RepeatLimit
		}
	}

	return service
}

func defaultBannedWords() []string {
	return []string{"fuck", "shit", "bitch", "asshole", "damn", "cunt", "whore", "nigger"}
}

func (s *ModerationService) AddModerator(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.moderators[userID] = struct{}{}
}

func (s *ModerationService) RemoveModerator(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.moderators, userID)
}

func (s *ModerationService) BlockUser(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.blockedUsers[userID] = struct{}{}
}

func (s *ModerationService) UnblockUser(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.blockedUsers, userID)
}

func (s *ModerationService) ShadowBanUser(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shadowBanned[userID] = struct{}{}
}

func (s *ModerationService) UnshadowBanUser(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.shadowBanned, userID)
}

func (s *ModerationService) MuteUser(userID int64, duration time.Duration) {
	if userID == 0 || duration <= 0 {
		return
	}
	expiresAt := time.Now().Add(duration)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mutedUsers[userID] = expiresAt
}

func (s *ModerationService) UnmuteUser(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.mutedUsers, userID)
}

func (s *ModerationService) HasPermission(userID int64, permission string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.permissions[userID] == nil {
		return false
	}
	_, ok := s.permissions[userID][permission]
	return ok
}

func (s *ModerationService) GrantPermission(userID int64, permission string) {
	if userID == 0 || permission == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.permissions[userID] == nil {
		s.permissions[userID] = make(map[string]struct{})
	}
	s.permissions[userID][permission] = struct{}{}
}

func (s *ModerationService) RevokePermission(userID int64, permission string) {
	if userID == 0 || permission == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.permissions[userID], permission)
}

func (s *ModerationService) IsBlocked(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.blockedUsers[userID]
	return ok
}

func (s *ModerationService) IsShadowBanned(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.shadowBanned[userID]
	return ok
}

func (s *ModerationService) IsMuted(userID int64) (bool, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	expire, ok := s.mutedUsers[userID]
	if !ok {
		return false, time.Time{}
	}
	if time.Now().After(expire) {
		return false, time.Time{}
	}
	return true, expire
}

func (s *ModerationService) IsModerator(userID int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.moderators[userID]
	return ok
}

func (s *ModerationService) EvaluateMessage(userID int64, text string, _ string) ModerationDecision {
	if userID == 0 {
		return ModerationDecision{Allowed: false, RejectReason: "invalid user"}
	}

	if s.IsBlocked(userID) {
		return ModerationDecision{Allowed: false, RejectReason: "user blocked"}
	}

	if muted, expiresAt := s.IsMuted(userID); muted {
		return ModerationDecision{Allowed: false, RejectReason: "user muted", MuteExpires: &expiresAt}
	}

	normalized := strings.ToLower(strings.TrimSpace(text))
	if normalized == "" {
		return ModerationDecision{Allowed: false, RejectReason: "message empty"}
	}

	if s.containsProfanity(normalized) {
		return ModerationDecision{Allowed: false, RejectReason: "profanity detected"}
	}

	if s.isFlooding(userID) {
		return ModerationDecision{Allowed: false, RejectReason: "flood protection triggered"}
	}

	if s.isSpam(userID, normalized) {
		return ModerationDecision{Allowed: false, RejectReason: "spam detected"}
	}

	if s.IsShadowBanned(userID) {
		return ModerationDecision{Allowed: true, ShadowBan: true}
	}

	s.recordMessage(userID, normalized)
	return ModerationDecision{Allowed: true}
}

func (s *ModerationService) containsProfanity(text string) bool {
	for _, word := range s.bannedWords {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

func (s *ModerationService) isFlooding(userID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	records := s.messageHistory[userID]
	if len(records) == 0 {
		return false
	}

	cutoff := time.Now().Add(-s.floodWindow)
	valid := make([]messageRecord, 0, len(records))
	for _, record := range records {
		if record.timestamp.After(cutoff) {
			valid = append(valid, record)
		}
	}
	s.messageHistory[userID] = valid
	return len(valid) >= s.maxMessages
}

func (s *ModerationService) isSpam(userID int64, text string) bool {
	s.mu.RLock()
	records := s.messageHistory[userID]
	s.mu.RUnlock()

	if len(records) == 0 {
		return false
	}

	cutoff := time.Now().Add(-s.repeatWindow)
	sameCount := 0
	for _, record := range records {
		if record.timestamp.Before(cutoff) {
			continue
		}
		if record.text == text {
			sameCount++
		}
	}
	return sameCount >= s.repeatLimit
}

func (s *ModerationService) recordMessage(userID int64, text string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	records := s.messageHistory[userID]
	records = append(records, messageRecord{text: text, timestamp: now})
	cutoff := now.Add(-s.floodWindow)
	keep := make([]messageRecord, 0, len(records))
	for _, record := range records {
		if record.timestamp.After(cutoff) {
			keep = append(keep, record)
		}
	}
	s.messageHistory[userID] = keep
}

func (s *ModerationService) GetUserStatus(userID int64) UserModerationStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	permissions := make([]string, 0, len(s.permissions[userID]))
	for p := range s.permissions[userID] {
		permissions = append(permissions, p)
	}

	muted, expiresAt := s.IsMuted(userID)
	status := UserModerationStatus{
		UserID:       userID,
		IsModerator:  s.IsModerator(userID),
		Blocked:      s.IsBlocked(userID),
		ShadowBanned: s.IsShadowBanned(userID),
		Muted:        muted,
		Permissions:  permissions,
	}
	if muted {
		status.MuteExpires = expiresAt.Format(time.RFC3339)
	}
	return status
}

func (s *ModerationService) ConfigureUserPermissions(userID int64, permissions []string) {
	if userID == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.permissions[userID] == nil {
		s.permissions[userID] = make(map[string]struct{})
	}
	for _, p := range permissions {
		s.permissions[userID][p] = struct{}{}
	}
}

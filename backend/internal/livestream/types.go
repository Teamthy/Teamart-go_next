package livestream

import "time"

type StreamState string

const (
	StreamStateScheduled StreamState = "scheduled"
	StreamStatePreparing StreamState = "preparing"
	StreamStateLive      StreamState = "live"
	StreamStatePaused    StreamState = "paused"
	StreamStateEnded     StreamState = "ended"
	StreamStateArchived  StreamState = "archived"
)

type EngagementType string

const (
	EngagementTypeLike     EngagementType = "like"
	EngagementTypeComment  EngagementType = "comment"
	EngagementTypeReaction EngagementType = "reaction"
	EngagementTypeShare    EngagementType = "share"
)

type StreamMetadata struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Category     string    `json:"category,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	CreatorID    int64     `json:"creator_id"`
	CreatorName  string    `json:"creator_name,omitempty"`
	CoHosts      []string  `json:"co_hosts,omitempty"`
	ScheduledAt  time.Time `json:"scheduled_at,omitempty"`
}

type StreamAnalytics struct {
	ViewerCount       int                    `json:"viewer_count"`
	UniqueViewerCount int                    `json:"unique_viewer_count"`
	TotalJoinCount    int                    `json:"total_join_count"`
	TotalLeaveCount   int                    `json:"total_leave_count"`
	LiveDuration      time.Duration          `json:"live_duration"`
	EngagementCounts  map[EngagementType]int `json:"engagement_counts"`
	StartedAt         time.Time              `json:"started_at,omitempty"`
	EndedAt           time.Time              `json:"ended_at,omitempty"`
}

type StreamInfo struct {
	Metadata  StreamMetadata  `json:"metadata"`
	State     StreamState     `json:"state"`
	Analytics StreamAnalytics `json:"analytics"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

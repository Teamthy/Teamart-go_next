package streaming

import "time"

type SessionState string

const (
	SessionStateIdle      SessionState = "idle"
	SessionStateIngesting SessionState = "ingesting"
	SessionStateLive      SessionState = "live"
	SessionStateEnded     SessionState = "ended"
)

type StreamingSession struct {
	ID           string       `json:"id"`
	StreamKey    string       `json:"stream_key"`
	Title        string       `json:"title"`
	State        SessionState `json:"state"`
	IngestURL    string       `json:"ingest_url,omitempty"`
	PlaybackURL  string       `json:"playback_url,omitempty"`
	HLSDirectory string       `json:"hls_directory,omitempty"`
	Profiles     []string     `json:"profiles,omitempty"`
	CDNProvider  string       `json:"cdn_provider,omitempty"`
	LastError    string       `json:"last_error,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

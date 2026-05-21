package livestream

import (
	"context"
	"fmt"
	"time"
)

// RTMPIngest represents a stubbed RTMP ingestion endpoint.
type RTMPIngest struct{}

// StartIngest starts consuming an RTMP stream and associates it with a session.
func (r *RTMPIngest) StartIngest(ctx context.Context, sessionID string, sourceURL string) error {
	// In production this would coordinate an ingest process (ffmpeg, media server).
	fmt.Printf("[rtmp] start ingest session=%s source=%s\n", sessionID, sourceURL)
	// Simulate async setup
	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("[rtmp] ingest started for %s\n", sessionID)
	}()
	return nil
}

// StopIngest stops the RTMP ingest for a session.
func (r *RTMPIngest) StopIngest(ctx context.Context, sessionID string) error {
	fmt.Printf("[rtmp] stop ingest session=%s\n", sessionID)
	return nil
}

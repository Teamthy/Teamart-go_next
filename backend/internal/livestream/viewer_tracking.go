package livestream

import "context"

// ViewerTracker collects analytics for livestream viewership.
type ViewerTracker struct{}

// NewViewerTracker creates a new viewer tracker.
func NewViewerTracker() *ViewerTracker {
	return &ViewerTracker{}
}

// TrackViewer records a viewer joining or leaving a stream.
func (t *ViewerTracker) TrackViewer(ctx context.Context, streamID string, userID int64, joined bool) error {
	return nil
}

// GetViewerAnalytics returns stream viewer analytics.
func (t *ViewerTracker) GetViewerAnalytics(ctx context.Context, streamID string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

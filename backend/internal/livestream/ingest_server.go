package livestream

import "context"

// IngestServer represents an RTMP ingest endpoint for livestream video.
type IngestServer struct {
	Address string
}

// NewIngestServer creates a new ingest server configuration.
func NewIngestServer(address string) *IngestServer {
	return &IngestServer{Address: address}
}

// Start begins processing ingested stream data.
func (s *IngestServer) Start(ctx context.Context) error {
	// Placeholder for RTMP ingest and stream handshake.
	return nil
}

// Stop terminates the ingest server.
func (s *IngestServer) Stop(ctx context.Context) error {
	return nil
}

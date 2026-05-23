package livestream

import "testing"

func TestStreamLifecycleAndAnalytics(t *testing.T) {
	service := NewService()

	metadata := StreamMetadata{
		ID:          "stream-1",
		Title:       "Launch Event",
		CreatorID:   10,
		CreatorName: "creator@example.com",
		Category:    "music",
		Tags:        []string{"concert", "live"},
		CoHosts:     []string{"host1", "host2"},
	}

	info, err := service.CreateStream(metadata, StreamStateScheduled)
	if err != nil {
		t.Fatalf("create stream failed: %v", err)
	}
	if info.State != StreamStateScheduled {
		t.Fatalf("expected scheduled state, got %s", info.State)
	}

	info, err = service.TransitionState(metadata.ID, StreamStatePreparing)
	if err != nil {
		t.Fatalf("transition failed: %v", err)
	}
	if info.State != StreamStatePreparing {
		t.Fatalf("expected preparing state, got %s", info.State)
	}

	info, err = service.TransitionState(metadata.ID, StreamStateLive)
	if err != nil {
		t.Fatalf("transition to live failed: %v", err)
	}
	if info.State != StreamStateLive {
		t.Fatalf("expected live state, got %s", info.State)
	}

	_, err = service.AddViewer(metadata.ID, 100)
	if err != nil {
		t.Fatalf("add viewer failed: %v", err)
	}
	_, err = service.AddViewer(metadata.ID, 101)
	if err != nil {
		t.Fatalf("add viewer failed: %v", err)
	}

	analytics, err := service.GetAnalytics(metadata.ID)
	if err != nil {
		t.Fatalf("get analytics failed: %v", err)
	}
	if analytics.ViewerCount != 2 {
		t.Fatalf("expected viewer count 2, got %d", analytics.ViewerCount)
	}
	if analytics.UniqueViewerCount != 2 {
		t.Fatalf("expected unique viewer count 2, got %d", analytics.UniqueViewerCount)
	}

	_, err = service.TrackEngagement(metadata.ID, EngagementTypeLike)
	if err != nil {
		t.Fatalf("track engagement failed: %v", err)
	}
	_, err = service.TrackEngagement(metadata.ID, EngagementTypeComment)
	if err != nil {
		t.Fatalf("track engagement failed: %v", err)
	}

	analytics, err = service.GetAnalytics(metadata.ID)
	if err != nil {
		t.Fatalf("get analytics failed: %v", err)
	}
	if analytics.EngagementCounts[EngagementTypeLike] != 1 {
		t.Fatalf("expected 1 like, got %d", analytics.EngagementCounts[EngagementTypeLike])
	}
	if analytics.EngagementCounts[EngagementTypeComment] != 1 {
		t.Fatalf("expected 1 comment, got %d", analytics.EngagementCounts[EngagementTypeComment])
	}

	_, err = service.RemoveViewer(metadata.ID, 100)
	if err != nil {
		t.Fatalf("remove viewer failed: %v", err)
	}
	analytics, err = service.GetAnalytics(metadata.ID)
	if err != nil {
		t.Fatalf("get analytics failed: %v", err)
	}
	if analytics.ViewerCount != 1 {
		t.Fatalf("expected viewer count 1 after leave, got %d", analytics.ViewerCount)
	}

	info, err = service.TransitionState(metadata.ID, StreamStateEnded)
	if err != nil {
		t.Fatalf("transition to ended failed: %v", err)
	}
	if info.State != StreamStateEnded {
		t.Fatalf("expected ended state, got %s", info.State)
	}

	info, err = service.TransitionState(metadata.ID, StreamStateArchived)
	if err != nil {
		t.Fatalf("transition to archived failed: %v", err)
	}
	if info.State != StreamStateArchived {
		t.Fatalf("expected archived state, got %s", info.State)
	}
}

func TestInvalidTransitions(t *testing.T) {
	service := NewService()
	metadata := StreamMetadata{ID: "stream-2", Title: "Test Event", CreatorID: 20}
	_, err := service.CreateStream(metadata, StreamStateLive)
	if err != nil {
		t.Fatalf("create stream failed: %v", err)
	}

	_, err = service.TransitionState(metadata.ID, StreamStateScheduled)
	if err == nil {
		t.Fatal("expected invalid transition error")
	}
}

func TestListStreams(t *testing.T) {
	service := NewService()
	metadata := StreamMetadata{ID: "stream-3", Title: "Talk Show", CreatorID: 30}
	_, _ = service.CreateStream(metadata, StreamStateScheduled)

	list := service.ListStreams()
	if len(list) != 1 {
		t.Fatalf("expected 1 stream, got %d", len(list))
	}
	if list[0].Metadata.ID != metadata.ID {
		t.Fatalf("expected stream ID %q, got %q", metadata.ID, list[0].Metadata.ID)
	}
}

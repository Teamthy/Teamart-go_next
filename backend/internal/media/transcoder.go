package media

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/teamart/commerce-api/config"
	"github.com/teamart/commerce-api/pkg/logger"
)

type Session struct {
	ID              string
	SourceURL       string
	OutputPath      string
	PlaylistName    string
	SegmentDuration int
	Profiles        []config.AdaptiveProfile
	BaseURL         string
}

type Transcoder interface {
	StartTranscoding(ctx context.Context, session Session) error
	StopTranscoding(sessionID string) error
	GetPlaybackManifestURL(sessionID string) string
}

type FFmpegTranscoder struct {
	cfg        config.StreamingConfig
	logger     *logger.Logger
	commands   map[string]*exec.Cmd
	commandsMu sync.Mutex
}

func NewFFmpegTranscoder(cfg config.StreamingConfig, logger *logger.Logger) *FFmpegTranscoder {
	return &FFmpegTranscoder{
		cfg:      cfg,
		logger:   logger,
		commands: make(map[string]*exec.Cmd),
	}
}

func (t *FFmpegTranscoder) StartTranscoding(ctx context.Context, session Session) error {
	if session.SourceURL == "" {
		return fmt.Errorf("source URL is required")
	}
	if len(session.Profiles) == 0 {
		return fmt.Errorf("at least one adaptive profile is required")
	}

	outputDir := filepath.Join(session.OutputPath, session.ID)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	var filterBuilder strings.Builder
	filterBuilder.WriteString(fmt.Sprintf("[0:v]split=%d", len(session.Profiles)))
	for i := range session.Profiles {
		filterBuilder.WriteString(fmt.Sprintf("[v%d]", i))
	}
	for i, profile := range session.Profiles {
		filterBuilder.WriteString(fmt.Sprintf(";[v%d]scale=%s[v%dout]", i, profile.Resolution, i))
	}

	var args []string
	args = append(args, "-hide_banner", "-loglevel", "info", "-i", session.SourceURL, "-filter_complex", filterBuilder.String())

	var streamMaps []string
	for idx, profile := range session.Profiles {
		args = append(args,
			"-map", fmt.Sprintf("[v%dout]", idx),
			"-map", "0:a",
			"-c:v:"+strconv.Itoa(idx), "libx264",
			"-b:v:"+strconv.Itoa(idx), profile.Bitrate,
			"-preset:"+strconv.Itoa(idx), "veryfast",
			"-c:a:"+strconv.Itoa(idx), "aac",
			"-b:a:"+strconv.Itoa(idx), profile.AudioBitrate,
		)
		streamMaps = append(streamMaps, fmt.Sprintf("v:%d,a:0", idx))
	}

	variantPattern := filepath.Join(outputDir, fmt.Sprintf("%s_%%v.m3u8", session.ID))
	args = append(args,
		"-f", "hls",
		"-hls_time", strconv.Itoa(session.SegmentDuration),
		"-hls_list_size", "5",
		"-hls_flags", "delete_segments+independent_segments+split_by_time",
		"-master_pl_name", session.PlaylistName,
		"-var_stream_map", strings.Join(streamMaps, " "),
		variantPattern,
	)

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch ffmpeg: %w", err)
	}

	t.commandsMu.Lock()
	defer t.commandsMu.Unlock()
	t.commands[session.ID] = cmd

	go func() {
		if err := cmd.Wait(); err != nil {
			t.logger.Errorf("ffmpeg process for session %s exited with error: %v", session.ID, err)
		}
		t.commandsMu.Lock()
		delete(t.commands, session.ID)
		t.commandsMu.Unlock()
	}()

	return nil
}

func (t *FFmpegTranscoder) StopTranscoding(sessionID string) error {
	t.commandsMu.Lock()
	cmd, ok := t.commands[sessionID]
	if !ok || cmd == nil || cmd.Process == nil {
		t.commandsMu.Unlock()
		return nil
	}
	delete(t.commands, sessionID)
	t.commandsMu.Unlock()

	if err := cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop ffmpeg process: %w", err)
	}
	return nil
}

func (t *FFmpegTranscoder) GetPlaybackManifestURL(sessionID string) string {
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(t.cfg.HLSBaseURL, "/"), sessionID, t.cfg.MasterPlaylistName)
}

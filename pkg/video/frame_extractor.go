package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type FrameExtractor struct {
	ffmpegPath string
}

type ExtractRequest struct {
	VideoPath  string
	OutputPath string
	Timestamp  string // Format: "00:00:05" or "5"
	Width      int
	Height     int
	Quality    int // 1-31, lower is better
}

func NewFrameExtractor() *FrameExtractor {
	ffmpegPath, _ := exec.LookPath("ffmpeg")
	return &FrameExtractor{ffmpegPath: ffmpegPath}
}

func (fe *FrameExtractor) ExtractFrame(req ExtractRequest) error {
	if fe.ffmpegPath == "" {
		return fmt.Errorf("ffmpeg not found. Please install FFmpeg for video processing")
	}

	args := []string{"-i", req.VideoPath}

	// Parse timestamp
	timestamp := req.Timestamp
	if !strings.Contains(timestamp, ":") {
		// Assume seconds, convert to HH:MM:SS
		seconds, err := strconv.Atoi(timestamp)
		if err == nil {
			hours := seconds / 3600
			minutes := (seconds % 3600) / 60
			secs := seconds % 60
			timestamp = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
		}
	}

	args = append(args, "-ss", timestamp, "-vframes", "1")

	// Add dimensions if specified
	if req.Width > 0 && req.Height > 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=%d:%d", req.Width, req.Height))
	}

	// Add quality if specified
	if req.Quality > 0 && req.Quality <= 31 {
		args = append(args, "-q:v", strconv.Itoa(req.Quality))
	}

	args = append(args, "-y", req.OutputPath)

	cmd := exec.Command(fe.ffmpegPath, args...)
	return cmd.Run()
}

func (fe *FrameExtractor) GetVideoInfo(videoPath string) (*VideoInfo, error) {
	if fe.ffmpegPath == "" {
		return nil, fmt.Errorf("ffmpeg not found")
	}

	args := []string{"-i", videoPath, "-hide_banner"}
	cmd := exec.Command(fe.ffmpegPath, args...)

	// FFmpeg outputs to stderr by default for info
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	return parseVideoInfo(string(output)), nil
}

type VideoInfo struct {
	Duration  string
	Width     int
	Height    int
	FrameRate string
	Bitrate   string
	Codec     string
}

func parseVideoInfo(output string) *VideoInfo {
	info := &VideoInfo{}
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "Duration:") {
			parts := strings.Split(line, "Duration:")
			if len(parts) > 1 {
				duration := strings.TrimSpace(parts[1])
				if idx := strings.Index(duration, ","); idx != -1 {
					info.Duration = duration[:idx]
				}
			}
		}

		if strings.Contains(line, "Video:") {
			// Extract dimensions
			if strings.Contains(line, "x") {
				fields := strings.Fields(line)
				for _, field := range fields {
					if strings.Contains(field, "x") && strings.Count(field, "x") == 1 {
						dims := strings.Split(field, "x")
						if len(dims) == 2 {
							if w, err := strconv.Atoi(dims[0]); err == nil {
								info.Width = w
							}
							if h, err := strconv.Atoi(dims[1]); err == nil {
								info.Height = h
							}
						}
						break
					}
				}
			}

			// Extract codec
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "Video:" && i+1 < len(fields) {
					info.Codec = fields[i+1]
					break
				}
			}
		}
	}

	return info
}

func (fe *FrameExtractor) ExtractMultipleFrames(videoPath, outputDir string, intervalSeconds int) ([]string, error) {
	if fe.ffmpegPath == "" {
		return nil, fmt.Errorf("ffmpeg not found")
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPattern := filepath.Join(outputDir, "frame_%04d.jpg")
	args := []string{
		"-i", videoPath,
		"-vf", fmt.Sprintf("fps=1/%d", intervalSeconds),
		"-y",
		outputPattern,
	}

	cmd := exec.Command(fe.ffmpegPath, args...)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to extract frames: %w", err)
	}

	// Get list of created files
	files, err := filepath.Glob(outputPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to get output files: %w", err)
	}

	return files, nil
}

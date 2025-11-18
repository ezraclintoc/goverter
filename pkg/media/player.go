package media

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Player struct {
	defaultPlayer string
}

func NewPlayer() *Player {
	return &Player{
		defaultPlayer: getDefaultPlayer(),
	}
}

func getDefaultPlayer() string {
	switch runtime.GOOS {
	case "windows":
		return "wmplayer"
	case "darwin":
		return "open"
	case "linux":
		players := []string{"vlc", "mpv", "mplayer", "totem", "gnome-mpv"}
		for _, player := range players {
			if _, err := exec.LookPath(player); err == nil {
				return player
			}
		}
		return "xdg-open"
	default:
		return "xdg-open"
	}
}

func (p *Player) Play(filePath string) error {
	if !p.isMediaFile(filePath) {
		return fmt.Errorf("not a supported media file: %s", filepath.Ext(filePath))
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", "/B", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	default:
		cmd = exec.Command(p.defaultPlayer, filePath)
	}

	return cmd.Start()
}

func (p *Player) isMediaFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	videoExts := []string{
		".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm",
		".m4v", ".3gp", ".ogv", ".ts", ".mts", ".m2ts",
	}

	audioExts := []string{
		".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".wma",
		".opus", ".aiff", ".au", ".ra", ".amr", ".ac3",
	}

	for _, videoExt := range videoExts {
		if ext == videoExt {
			return true
		}
	}

	for _, audioExt := range audioExts {
		if ext == audioExt {
			return true
		}
	}

	return false
}

func (p *Player) GetSupportedFormats() []string {
	return append(p.getVideoFormats(), p.getAudioFormats()...)
}

func (p *Player) getVideoFormats() []string {
	return []string{
		"mp4", "avi", "mkv", "mov", "wmv", "flv", "webm",
		"m4v", "3gp", "ogv", "ts", "mts", "m2ts",
	}
}

func (p *Player) getAudioFormats() []string {
	return []string{
		"mp3", "wav", "flac", "aac", "ogg", "m4a", "wma",
		"opus", "aiff", "au", "ra", "amr", "ac3",
	}
}

func (p *Player) IsPlayable(filePath string) bool {
	return p.isMediaFile(filePath)
}

func (p *Player) GetPlayerInfo() map[string]string {
	info := make(map[string]string)
	info["default_player"] = p.defaultPlayer
	info["os"] = runtime.GOOS
	info["supported_video"] = strings.Join(p.getVideoFormats(), ", ")
	info["supported_audio"] = strings.Join(p.getAudioFormats(), ", ")
	return info
}

type PreviewInfo struct {
	FilePath   string
	Thumbnail  string
	Duration   string
	Title      string
	Artist     string
	Resolution string
	FileSize   string
	Format     string
}

func (p *Player) GeneratePreview(filePath string) (*PreviewInfo, error) {
	if !p.isMediaFile(filePath) {
		return nil, fmt.Errorf("not a supported media file")
	}

	info := &PreviewInfo{
		FilePath: filePath,
		Format:   strings.TrimPrefix(filepath.Ext(filePath), "."),
	}

	if stat, err := os.Stat(filePath); err == nil {
		info.FileSize = formatFileSize(stat.Size())
	}

	if err := p.extractMediaInfo(filePath, info); err != nil {
		// Continue without detailed info
	}

	if p.isVideoFile(filePath) {
		thumbnailPath, err := p.generateThumbnail(filePath)
		if err == nil {
			info.Thumbnail = thumbnailPath
		}
	}

	return info, nil
}

func (p *Player) isVideoFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	videoExts := []string{
		".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm",
		".m4v", ".3gp", ".ogv", ".ts", ".mts", ".m2ts",
	}

	for _, videoExt := range videoExts {
		if ext == videoExt {
			return true
		}
	}
	return false
}

func (p *Player) extractMediaInfo(filePath string, info *PreviewInfo) error {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams", filePath)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	outputStr := string(output)

	if strings.Contains(outputStr, "duration") {
		parts := strings.Split(outputStr, "duration")
		if len(parts) > 1 {
			durationStr := strings.Split(strings.Split(parts[1], ",")[0], "\"")[1]
			info.Duration = formatDuration(durationStr)
		}
	}

	if strings.Contains(outputStr, "TAG:TITLE") {
		parts := strings.Split(outputStr, "TAG:TITLE")
		if len(parts) > 1 {
			info.Title = strings.Split(parts[1], "\"")[1]
		}
	}

	if strings.Contains(outputStr, "TAG:ARTIST") {
		parts := strings.Split(outputStr, "TAG:ARTIST")
		if len(parts) > 1 {
			info.Artist = strings.Split(parts[1], "\"")[1]
		}
	}

	if strings.Contains(outputStr, "width") && strings.Contains(outputStr, "height") {
		widthParts := strings.Split(outputStr, "width")
		heightParts := strings.Split(outputStr, "height")
		if len(widthParts) > 1 && len(heightParts) > 1 {
			width := strings.Split(strings.Split(widthParts[1], ",")[0], "\"")[1]
			height := strings.Split(strings.Split(heightParts[1], ",")[0], "\"")[1]
			info.Resolution = fmt.Sprintf("%sx%s", width, height)
		}
	}

	return nil
}

func (p *Player) generateThumbnail(filePath string) (string, error) {
	thumbnailPath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + "_thumb.jpg"

	cmd := exec.Command("ffmpeg", "-i", filePath, "-ss", "00:00:01", "-vframes", "1", "-q:v", "2", "-y", thumbnailPath)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return thumbnailPath, nil
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDuration(seconds string) string {
	var totalSeconds float64
	fmt.Sscanf(seconds, "%f", &totalSeconds)

	hours := int(totalSeconds) / 3600
	minutes := (int(totalSeconds) % 3600) / 60
	secs := int(totalSeconds) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

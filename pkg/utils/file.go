package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetFilesByExtension(dir string, extensions []string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		for _, validExt := range extensions {
			if ext == validExt {
				files = append(files, path)
				break
			}
		}

		return nil
	})

	return files, err
}

func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func GetUniqueFilename(basePath string) string {
	ext := filepath.Ext(basePath)
	base := basePath[:len(basePath)-len(ext)]

	counter := 1
	newPath := basePath

	for {
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}

		newPath = fmt.Sprintf("%s_%d%s", base, counter, ext)
		counter++
	}
}

func FormatFileSize(bytes int64) string {
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

func IsValidVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	videoExts := []string{".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm"}
	for _, validExt := range videoExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

func IsValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff"}
	for _, validExt := range imageExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

func IsValidAudioFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	audioExts := []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a"}
	for _, validExt := range audioExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

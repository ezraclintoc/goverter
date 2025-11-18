package converter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type FormatSupport struct {
	InputFormats  []string
	OutputFormats []string
	Category      string
}

var SupportedFormats = map[string]FormatSupport{
	"mp4": {
		InputFormats:  []string{".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm", ".m4v"},
		OutputFormats: []string{".mp4", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".webm", ".m4v", ".gif", ".mp3"},
		Category:      "video",
	},
	"jpg": {
		InputFormats:  []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".svg"},
		OutputFormats: []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".pdf"},
		Category:      "image",
	},
	"mp3": {
		InputFormats:  []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".wma"},
		OutputFormats: []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a"},
		Category:      "audio",
	},
	"pdf": {
		InputFormats:  []string{".pdf", ".doc", ".docx", ".txt", ".rtf", ".odt"},
		OutputFormats: []string{".pdf", ".txt", ".html", ".docx", ".jpg", ".png"},
		Category:      "document",
	},
}

type ConversionRequest struct {
	InputPath  string
	OutputPath string
	Options    map[string]string
	Progress   chan float64
	Error      chan error
}

type Converter struct {
	ffmpegPath string
	magickPath string
	pandocPath string
}

func NewConverter() *Converter {
	ffmpegPath, _ := exec.LookPath("ffmpeg")
	magickPath, _ := exec.LookPath("magick")
	pandocPath, _ := exec.LookPath("pandoc")

	return &Converter{
		ffmpegPath: ffmpegPath,
		magickPath: magickPath,
		pandocPath: pandocPath,
	}
}

func (c *Converter) Convert(req ConversionRequest) error {
	inputExt := strings.ToLower(filepath.Ext(req.InputPath))

	category := c.getCategory(inputExt)

	switch category {
	case "video":
		return c.convertVideo(req)
	case "image":
		return c.convertImage(req)
	case "audio":
		return c.convertAudio(req)
	case "document":
		return c.convertDocument(req)
	default:
		return fmt.Errorf("unsupported file type: %s", inputExt)
	}
}

func (c *Converter) getCategory(ext string) string {
	for _, format := range SupportedFormats {
		for _, inputExt := range format.InputFormats {
			if ext == inputExt {
				return format.Category
			}
		}
	}
	return ""
}

func (c *Converter) convertVideo(req ConversionRequest) error {
	if c.ffmpegPath == "" {
		return fmt.Errorf("ffmpeg not found. Please install FFmpeg for video conversions")
	}

	args := []string{"-i", req.InputPath}

	// Handle audio extraction (video to audio)
	if isAudioFormat(filepath.Ext(req.OutputPath)) {
		args = append(args, "-vn", "-acodec", "libmp3lame")
		if bitrate, ok := req.Options["bitrate"]; ok {
			args = append(args, "-ab", bitrate)
		} else {
			args = append(args, "-ab", "192k")
		}
	} else if filepath.Ext(req.OutputPath) == ".gif" {
		// Handle GIF creation (video to GIF)
		fps := "10"
		if fpsVal, ok := req.Options["fps"]; ok {
			fps = fpsVal
		}
		scale := "480:-1"
		if width, ok := req.Options["width"]; ok {
			if height, ok := req.Options["height"]; ok {
				scale = fmt.Sprintf("%s:%s", width, height)
			} else {
				scale = fmt.Sprintf("%s:-1", width)
			}
		}
		args = append(args, "-vf", fmt.Sprintf("fps=%s,scale=%s:flags=lanczos,split[s0][1],palettegen[p1][s0]", fps, scale))
		args = append(args, "-map", "[p1]", "-f", "gif")
	} else {
		// Regular video to video conversion
		// Add quality settings
		if quality, ok := req.Options["quality"]; ok {
			args = append(args, "-crf", quality)
		}

		// Add bitrate settings
		if bitrate, ok := req.Options["bitrate"]; ok {
			args = append(args, "-b:v", bitrate)
		}
	}

	args = append(args, "-y", req.OutputPath)

	cmd := exec.Command(c.ffmpegPath, args...)
	return cmd.Run()
}

func (c *Converter) convertImage(req ConversionRequest) error {
	if c.magickPath == "" {
		return fmt.Errorf("ImageMagick not found. Please install ImageMagick for image conversions")
	}

	args := []string{req.InputPath}

	// Add quality settings
	if quality, ok := req.Options["quality"]; ok {
		args = append(args, "-quality", quality)
	}

	// Add resize settings
	if width, ok := req.Options["width"]; ok {
		if height, ok := req.Options["height"]; ok {
			args = append(args, "-resize", fmt.Sprintf("%sx%s", width, height))
		} else {
			args = append(args, "-resize", width)
		}
	}

	args = append(args, req.OutputPath)

	cmd := exec.Command(c.magickPath, args...)
	return cmd.Run()
}

func (c *Converter) convertAudio(req ConversionRequest) error {
	if c.ffmpegPath == "" {
		return fmt.Errorf("ffmpeg not found. Please install FFmpeg for audio conversions")
	}

	args := []string{"-i", req.InputPath}

	// Add bitrate settings
	if bitrate, ok := req.Options["bitrate"]; ok {
		args = append(args, "-b:a", bitrate)
	}

	// Add sample rate settings
	if sampleRate, ok := req.Options["sample_rate"]; ok {
		args = append(args, "-ar", sampleRate)
	}

	args = append(args, "-y", req.OutputPath)

	cmd := exec.Command(c.ffmpegPath, args...)
	return cmd.Run()
}

func (c *Converter) convertDocument(req ConversionRequest) error {
	if c.pandocPath == "" {
		return fmt.Errorf("pandoc not found. Please install pandoc for document conversions")
	}

	args := []string{req.InputPath, "-o", req.OutputPath}

	// Add PDF-specific options
	if filepath.Ext(req.OutputPath) == ".pdf" {
		args = append(args, "--pdf-engine=pdflatex")
	}

	cmd := exec.Command(c.pandocPath, args...)
	return cmd.Run()
}

func (c *Converter) BatchConvert(requests []ConversionRequest) []error {
	errors := make([]error, len(requests))

	for i, req := range requests {
		if err := c.Convert(req); err != nil {
			errors[i] = fmt.Errorf("failed to convert %s: %w", req.InputPath, err)
		}
	}

	return errors
}

func (c *Converter) GetSupportedFormats() map[string]FormatSupport {
	return SupportedFormats
}

func (c *Converter) IsFormatSupported(ext string) bool {
	ext = strings.ToLower(ext)
	for _, format := range SupportedFormats {
		for _, inputExt := range format.InputFormats {
			if ext == inputExt {
				return true
			}
		}
	}
	return false
}

func (c *Converter) GetOutputFormats(inputExt string) []string {
	inputExt = strings.ToLower(inputExt)
	for _, format := range SupportedFormats {
		for _, ext := range format.InputFormats {
			if inputExt == ext {
				return format.OutputFormats
			}
		}
	}
	return []string{}
}

func ValidateTools() map[string]bool {
	return map[string]bool{
		"ffmpeg":      checkTool("ffmpeg"),
		"imagemagick": checkTool("magick"),
		"pandoc":      checkTool("pandoc"),
	}
}

func isAudioFormat(ext string) bool {
	audioExts := []string{".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a"}
	for _, audioExt := range audioExts {
		if ext == audioExt {
			return true
		}
	}
	return false
}

func checkTool(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

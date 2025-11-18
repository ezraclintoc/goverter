package image

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type CropRequest struct {
	InputPath     string
	OutputPath    string
	X, Y          int
	Width, Height int
	Quality       int // 1-100 for JPEG
}

type ResizeRequest struct {
	InputPath     string
	OutputPath    string
	Width, Height int
	Quality       int
}

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Crop(req CropRequest) error {
	img, err := imaging.Open(req.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	cropped := imaging.Crop(img, image.Rect(req.X, req.Y, req.X+req.Width, req.Y+req.Height))

	return p.saveImage(cropped, req.OutputPath, req.Quality)
}

func (p *Processor) Resize(req ResizeRequest) error {
	img, err := imaging.Open(req.InputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	resized := imaging.Resize(img, req.Width, req.Height, imaging.Lanczos)

	return p.saveImage(resized, req.OutputPath, req.Quality)
}

func (p *Processor) Rotate(inputPath, outputPath string, degrees float64, quality int) error {
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	var rotated *image.NRGBA
	switch degrees {
	case 90:
		rotated = imaging.Rotate90(img)
	case 180:
		rotated = imaging.Rotate180(img)
	case 270:
		rotated = imaging.Rotate270(img)
	default:
		rotated = imaging.Rotate(img, degrees, image.Transparent)
	}

	return p.saveImage(rotated, outputPath, quality)
}

func (p *Processor) Flip(inputPath, outputPath string, horizontal bool, quality int) error {
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	var flipped *image.NRGBA
	if horizontal {
		flipped = imaging.FlipH(img)
	} else {
		flipped = imaging.FlipV(img)
	}

	return p.saveImage(flipped, outputPath, quality)
}

func (p *Processor) GetImageInfo(imagePath string) (*ImageInfo, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	stat, err := os.Stat(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &ImageInfo{
		Width:    img.Width,
		Height:   img.Height,
		Format:   strings.TrimPrefix(filepath.Ext(imagePath), "."),
		FileSize: stat.Size(),
	}, nil
}

type ImageInfo struct {
	Width    int
	Height   int
	Format   string
	FileSize int64
}

func (p *Processor) saveImage(img image.Image, outputPath string, quality int) error {
	ext := strings.ToLower(filepath.Ext(outputPath))

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	switch ext {
	case ".jpg", ".jpeg":
		if quality <= 0 || quality > 100 {
			quality = 95
		}
		return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	case ".png":
		return png.Encode(file, img)
	default:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	}
}

func (p *Processor) BatchProcess(requests []interface{}) []error {
	errors := make([]error, len(requests))

	for i, req := range requests {
		switch r := req.(type) {
		case CropRequest:
			if err := p.Crop(r); err != nil {
				errors[i] = fmt.Errorf("failed to crop %s: %w", r.InputPath, err)
			}
		case ResizeRequest:
			if err := p.Resize(r); err != nil {
				errors[i] = fmt.Errorf("failed to resize %s: %w", r.InputPath, err)
			}
		default:
			errors[i] = fmt.Errorf("unsupported request type at index %d", i)
		}
	}

	return errors
}

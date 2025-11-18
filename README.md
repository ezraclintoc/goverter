# Goverter

A comprehensive file conversion tool inspired by VERT.sh, built in Go with both CLI and GUI interfaces.

## Features

- **Multi-format Support**: Convert between 250+ file formats
  - Video: MP4, AVI, MKV, MOV, WMV, FLV, WebM, etc.
  - Image: JPG, PNG, GIF, BMP, WebP, TIFF, SVG, etc.
  - Audio: MP3, WAV, FLAC, AAC, OGG, M4A, etc.
  - Document: PDF, DOC, DOCX, TXT, HTML, etc.

- **CLI Interface**: Command-line tool for automation and scripting
- **GUI Interface**: User-friendly desktop application (coming soon)
- **Bulk Processing**: Convert multiple files at once
- **Media Processing**: 
  - Extract frames from videos
  - Crop and resize images
  - Rotate and flip images
- **Local Processing**: All conversions happen on your machine
- **No File Limits**: No restrictions on file size or quantity

## Installation

### Prerequisites

Install required tools for full functionality:

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install ffmpeg imagemagick pandoc

# macOS
brew install ffmpeg imagemagick pandoc

# Windows
# Download and install from official websites:
# - FFmpeg: https://ffmpeg.org/download.html
# - ImageMagick: https://imagemagick.org/script/download.php
# - Pandoc: https://pandoc.org/installing.html
```

### Build from Source

```bash
git clone <repository-url>
cd goverter
go mod tidy
go build -o goverter-cli ./cmd/cli
go build -o goverter-gui ./cmd/gui
```

## Usage

### CLI Commands

#### Basic Conversion
```bash
# Convert video
./goverter-cli convert -i input.mp4 -o output.avi

# Convert image
./goverter-cli convert -i image.jpg -o image.png

# Convert audio
./goverter-cli convert -i audio.mp3 -o audio.wav

# Convert document
./goverter-cli convert -i document.pdf -o document.txt
```

#### Video Frame Extraction
```bash
# Extract frame at specific timestamp
./goverter-cli frame 00:00:05 -i video.mp4 -o frame.jpg

# Extract frame with custom dimensions
./goverter-cli frame 10 -i video.mp4 -o frame.jpg --width 800 --height 600
```

#### Image Processing
```bash
# Crop image
./goverter-cli crop 100 100 400 300 -i image.jpg -o cropped.jpg

# Resize image
./goverter-cli resize 800 600 -i image.jpg -o resized.jpg
```

#### File Information
```bash
# Get media file info
./goverter-cli info -i video.mp4
./goverter-cli info -i image.jpg
```

#### Quality Settings
```bash
# Set quality for video (CRF value, lower = better)
./goverter-cli convert -i input.mp4 -o output.avi --quality 23

# Set quality for images (1-100)
./goverter-cli convert -i image.jpg -o image.png --quality 95
```

#### Bulk Conversion
```bash
# Convert all files in directory
./goverter-cli convert --bulk /path/to/files --format mp4
```

### GUI Application

```bash
# Launch GUI (coming soon)
./goverter-gui
```

## Supported Formats

### Video Formats
- **Input**: MP4, AVI, MKV, MOV, WMV, FLV, WebM, M4V, 3GP, etc.
- **Output**: MP4, AVI, MKV, MOV, WMV, FLV, WebM, M4V, GIF, MP3

### Image Formats
- **Input**: JPG, JPEG, PNG, GIF, BMP, WebP, TIFF, SVG, ICO, etc.
- **Output**: JPG, JPEG, PNG, GIF, BMP, WebP, TIFF, PDF

### Audio Formats
- **Input**: MP3, WAV, FLAC, AAC, OGG, M4A, WMA, etc.
- **Output**: MP3, WAV, FLAC, AAC, OGG, M4A

### Document Formats
- **Input**: PDF, DOC, DOCX, TXT, RTF, ODT, HTML, etc.
- **Output**: PDF, TXT, HTML, DOCX, JPG, PNG

## Project Structure

```
goverter/
├── cmd/
│   ├── cli/          # CLI application
│   └── gui/          # GUI application
├── pkg/
│   ├── converter/     # Core conversion logic
│   ├── image/         # Image processing
│   ├── video/         # Video processing
│   └── utils/         # Utility functions
├── internal/
│   └── config/       # Configuration
├── go.mod
├── go.sum
└── README.md
```

## Dependencies

- **FFmpeg**: Video and audio conversion
- **ImageMagick**: Image conversion and processing
- **Pandoc**: Document conversion
- **Cobra**: CLI framework
- **Fyne**: GUI framework

## Configuration

The tool automatically detects installed tools and enables features accordingly:

```bash
# Check tool availability
./goverter-cli info
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Inspiration

This project is inspired by [VERT.sh](https://vert.sh), an excellent open-source file converter. Goverter aims to provide similar functionality with the added benefits of Go's performance and cross-platform support.

## Roadmap

- [ ] Complete GUI implementation with drag-and-drop
- [ ] Add more image processing filters
- [ ] Support for archive formats (ZIP, RAR, etc.)
- [ ] Web interface
- [ ] Plugin system for custom converters
- [ ] Batch processing with progress tracking
- [ ] Preset conversion profiles
- [ ] Integration with cloud storage
# Goverter

A comprehensive file conversion tool inspired by VERT.sh, built in Go with both CLI and GUI interfaces.

## ⚠️ Disclaimer

**This project was vibe coded** - built with passion, creativity, and a love for file conversion! While it aims to provide robust functionality, please note:

- This is an open-source project developed for educational and practical purposes
- Always backup your important files before conversion
- The software is provided "as-is" without warranties
- Contributions and feedback are welcome to improve the project

## Features

- **Multi-format Support**: Convert between 250+ file formats
  - Video: MP4, AVI, MKV, MOV, WMV, FLV, WebM, etc.
  - Image: JPG, PNG, GIF, BMP, WebP, TIFF, SVG, etc.
  - Audio: MP3, WAV, FLAC, AAC, OGG, M4A, etc.
  - Document: PDF, DOC, DOCX, TXT, HTML, etc.

- **CLI Interface**: Command-line tool for automation and scripting
- **GUI Interface**: User-friendly desktop application
- **Bulk Processing**: Convert multiple files at once
- **Media Processing**: 
  - Extract frames from videos
  - Crop and resize images
  - Rotate and flip images
- **Local Processing**: All conversions happen on your machine
- **No File Limits**: No restrictions on file size or quantity

## Installation

### Quick Install (Recommended)

```bash
# Download and run the installer
curl -fsSL https://raw.githubusercontent.com/ezraclintoc/goverter/main/install.sh | bash
```

### Manual Installation

#### Prerequisites

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
git clone https://github.com/ezraclintoc/goverter.git
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

#### Video to GIF
```bash
# Convert video to GIF with custom settings
./goverter-cli convert -i video.mp4 -o video.gif --quality 10
```

#### Video to Audio
```bash
# Extract audio from video
./goverter-cli convert -i video.mp4 -o audio.mp3 --bitrate 192k
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
# Launch GUI
./goverter-gui
```

The GUI provides:
- **Convert Tab**: Drag & drop files, select output format, adjust quality
- **Image Tools Tab**: Crop, resize, rotate images
- **Video Tools Tab**: Extract frames, convert to GIF, extract audio
- **Info Tab**: Check tool availability and supported formats

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
├── install.sh         # Cross-platform installer
├── build-release.sh  # Release build script
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

This was vibe coded with love, but contributions make it better! 

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Inspiration

This project is inspired by [VERT.sh](https://vert.sh), an excellent open-source file converter. Goverter aims to provide similar functionality with added benefits of:

- Go's performance and cross-platform support
- Enhanced GUI with drag-and-drop
- Video to GIF conversion
- Audio extraction from videos
- Comprehensive image processing tools

## Community

- 🐛 **Bug Reports**: [Issues](https://github.com/ezraclintoc/goverter/issues)
- 💡 **Feature Requests**: [Discussions](https://github.com/ezraclintoc/goverter/discussions)
- 🤝 **Contributions**: [Pull Requests](https://github.com/ezraclintoc/goverter/pulls)

---

**Built with passion for the open-source community! 🚀**

## Roadmap

- [ ] Complete GUI implementation with drag-and-drop ✅
- [ ] Add more image processing filters
- [ ] Support for archive formats (ZIP, RAR, etc.)
- [ ] Web interface
- [ ] Plugin system for custom converters
- [ ] Batch processing with progress tracking ✅
- [ ] Preset conversion profiles
- [ ] Integration with cloud storage
- [ ] Real-time conversion preview
- [ ] Advanced video editing features
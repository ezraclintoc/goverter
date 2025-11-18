package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"goverter/pkg/converter"
	"goverter/pkg/image"
	"goverter/pkg/media"
	"goverter/pkg/video"
)

type GUI struct {
	app            fyne.App
	window         fyne.Window
	converter      *converter.Converter
	imageProcessor *image.Processor
	frameExtractor *video.FrameExtractor
	mediaPlayer    *media.Player

	// UI Elements
	fileList      *widget.List
	files         []string
	statusLabel   *widget.Label
	progressBar   *widget.ProgressBar
	outputFormat  *widget.Select
	qualitySlider *widget.Slider
	qualityLabel  *widget.Label
	outputDir     *widget.Entry
	convertBtn    *widget.Button
	addFilesBtn   *widget.Button
	clearBtn      *widget.Button
	playBtn       *widget.Button

	// Current tab
	currentTab       string
	contentContainer *fyne.Container
	pages            map[string]fyne.CanvasObject
	previewArea      *fyne.Container
}

func main() {
	g := &GUI{
		app:            app.NewWithID("com.goverter.app"),
		converter:      converter.NewConverter(),
		imageProcessor: image.NewProcessor(),
		frameExtractor: video.NewFrameExtractor(),
		mediaPlayer:    media.NewPlayer(),
		files:          make([]string, 0),
		currentTab:     "upload",
	}

	g.createUI()
	g.window.ShowAndRun()
}

func (g *GUI) createUI() {
	g.window = g.app.NewWindow("🔄 Goverter - File Converter")
	g.window.Resize(fyne.NewSize(1200, 800))
	g.window.CenterOnScreen()

	// Initialize status label to prevent nil pointer
	g.statusLabel = widget.NewLabel("🔄 Ready")

	// Create pages
	uploadPage := g.createUploadPage()
	convertPage := g.createConvertPage()
	toolsPage := g.createToolsPage()
	settingsPage := g.createSettingsPage()

	// Store pages for switching
	g.pages = map[string]fyne.CanvasObject{
		"upload":   uploadPage,
		"convert":  convertPage,
		"tools":    toolsPage,
		"settings": settingsPage,
	}

	// Tab bar with navigation buttons (centered)
	tabBar := container.NewHBox(
		g.createTabButton("📤 Upload", g.currentTab == "upload", func() { g.switchTab("upload") }),
		g.createTabButton("🔄 Convert", g.currentTab == "convert", func() { g.switchTab("convert") }),
		g.createTabButton("🛠️ Tools", g.currentTab == "tools", func() { g.switchTab("tools") }),
		g.createTabButton("⚙️ Settings", g.currentTab == "settings", func() { g.switchTab("settings") }),
	)

	// Center the tab bar
	centeredTabBar := container.NewCenter(tabBar)

	// Main content area (centered)
	g.contentContainer = container.NewStack(uploadPage)

	// Layout with centered content
	mainContent := container.NewVBox(
		centeredTabBar,
		widget.NewSeparator(),
		container.NewCenter(g.contentContainer),
	)

	g.window.SetContent(mainContent)
}

func (g *GUI) createTabButton(text string, isActive bool, onClick func()) *widget.Button {
	btn := widget.NewButton(text, onClick)
	if isActive {
		btn.Importance = widget.HighImportance
	}
	return btn
}

func (g *GUI) switchTab(tabName string) {
	g.currentTab = tabName

	// Switch content based on tab
	if content, exists := g.pages[tabName]; exists {
		g.currentTab = tabName

		// Update content container
		if g.contentContainer != nil {
			g.contentContainer.Objects = []fyne.CanvasObject{content}
			g.contentContainer.Refresh()
		}
	}
}

func (g *GUI) createUploadPage() fyne.CanvasObject {
	// Drag and drop area
	dropArea := g.createDropArea()

	// Preview area
	g.previewArea = container.NewVBox()
	previewCard := widget.NewCard("📸 Media Preview", "", container.NewScroll(g.previewArea))

	// File list
	g.fileList = widget.NewList(
		func() int { return len(g.files) },
		func() fyne.CanvasObject {
			playBtn := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil)
			playBtn.Importance = widget.LowImportance
			previewBtn := widget.NewButtonWithIcon("", theme.SearchReplaceIcon(), nil)
			previewBtn.Importance = widget.LowImportance
			return container.NewHBox(
				playBtn,
				previewBtn,
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("File"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			fileName := filepath.Base(g.files[i])
			if hbox, ok := o.(*fyne.Container); ok && len(hbox.Objects) > 3 {
				hbox.Objects[3].(*widget.Label).SetText(fileName)

				// Configure play button
				playBtn := hbox.Objects[0].(*widget.Button)
				previewBtn := hbox.Objects[1].(*widget.Button)

				if g.mediaPlayer.IsPlayable(g.files[i]) {
					playBtn.Enable()
					playBtn.OnTapped = func() { g.playMedia(g.files[i]) }
					previewBtn.Enable()
					previewBtn.OnTapped = func() { g.showPreview(g.files[i]) }
				} else {
					playBtn.Disable()
					previewBtn.Disable()
				}
			}
		},
	)

	fileListContainer := container.NewVBox(
		widget.NewCard("📁 Selected Files", "", container.NewVScroll(
			container.NewVBox(g.fileList),
		)),
	)

	// Quick actions
	quickActions := container.NewVBox(
		widget.NewButton("🔄 Convert All", g.convertFiles),
		widget.NewButton("🗑️ Clear All", func() {
			g.files = make([]string, 0)
			g.fileList.Refresh()
			g.updateStatus("🗑️ Cleared all files")
		}),
	)

	// Main layout
	leftPanel := container.NewVBox(
		dropArea,
		widget.NewSeparator(),
		fileListContainer,
		widget.NewSeparator(),
		quickActions,
	)

	// Right panel - Quick info
	rightPanel := container.NewVBox(
		widget.NewCard("📊 Quick Stats", "", container.NewVBox(
			widget.NewLabel(fmt.Sprintf("📁 Files: %d", len(g.files))),
			widget.NewSeparator(),
			widget.NewLabel("💡 Tips:"),
			widget.NewLabel("• Drag & drop files here"),
			widget.NewLabel("• Use Convert tab for options"),
			widget.NewLabel("• Tools tab for advanced features"),
		),
		),
		widget.NewCard("🔧 Tool Status", "", g.createToolStatus()),
	)

	return container.NewHSplit(leftPanel, rightPanel)
}

func (g *GUI) createConvertPage() fyne.CanvasObject {
	// File selection area
	dropArea := g.createDropArea()

	// File list with preview buttons
	g.fileList = widget.NewList(
		func() int { return len(g.files) },
		func() fyne.CanvasObject {
			playBtn := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil)
			playBtn.Importance = widget.LowImportance
			previewBtn := widget.NewButtonWithIcon("", theme.SearchReplaceIcon(), nil)
			previewBtn.Importance = widget.LowImportance
			return container.NewHBox(
				playBtn,
				previewBtn,
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("File"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			fileName := filepath.Base(g.files[i])
			if hbox, ok := o.(*fyne.Container); ok && len(hbox.Objects) > 3 {
				hbox.Objects[3].(*widget.Label).SetText(fileName)

				// Configure play button
				playBtn := hbox.Objects[0].(*widget.Button)
				previewBtn := hbox.Objects[1].(*widget.Button)

				if g.mediaPlayer.IsPlayable(g.files[i]) {
					playBtn.Enable()
					playBtn.OnTapped = func() { g.playMedia(g.files[i]) }
					previewBtn.Enable()
					previewBtn.OnTapped = func() { g.showPreview(g.files[i]) }
				} else {
					playBtn.Disable()
					previewBtn.Disable()
				}
			}
		},
	)

	fileListContainer := container.NewVBox(
		widget.NewCard("📁 Selected Files", "", container.NewVScroll(
			container.NewVBox(g.fileList),
		)),
	)

	// Conversion options with dynamic format selection
	optionsContainer := g.createDynamicConversionOptions()

	// Status and progress
	g.statusLabel = widget.NewLabel("🔄 Ready to convert files")
	g.progressBar = widget.NewProgressBar()

	// Action buttons
	buttonContainer := container.NewHBox(
		g.addFilesBtn,
		g.clearBtn,
		g.convertBtn,
	)

	// Main layout (centered)
	leftPanel := container.NewVBox(
		dropArea,
		widget.NewSeparator(),
		fileListContainer,
	)

	rightPanel := container.NewVBox(
		widget.NewCard("⚙️ Conversion Options", "", optionsContainer),
		widget.NewSeparator(),
		widget.NewCard("📊 Status", "", container.NewVBox(
			g.statusLabel,
			g.progressBar,
		)),
		widget.NewSeparator(),
		buttonContainer,
	)

	// Center the content
	centeredContent := container.NewCenter(
		container.NewHSplit(leftPanel, rightPanel),
	)

	return centeredContent
}

func (g *GUI) createToolsPage() fyne.CanvasObject {
	// Create tool sections
	imageTools := g.createImageToolsSection()
	videoTools := g.createVideoToolsSection()
	audioTools := g.createAudioToolsSection()

	// Layout with tabs for tools
	toolTabs := container.NewAppTabs(
		container.NewTabItem("🖼️ Images", imageTools),
		container.NewTabItem("🎬 Videos", videoTools),
		container.NewTabItem("🎵 Audio", audioTools),
	)

	return widget.NewCard("", "", toolTabs)
}

func (g *GUI) createSettingsPage() fyne.CanvasObject {
	// Tool status
	tools := converter.ValidateTools()
	toolStatus := container.NewVBox()

	for tool, available := range tools {
		status := "❌ Not Available"
		if available {
			status = "✅ Available"
		}
		toolStatus.Add(widget.NewLabel(fmt.Sprintf("%s: %s", tool, status)))
	}

	// Supported formats
	supportedFormats := g.converter.GetSupportedFormats()
	formatContainer := container.NewVBox()

	for category, formats := range supportedFormats {
		categoryTitle := strings.ToUpper(string(category[0])) + strings.ToLower(string(category[1:]))
		categoryLabel := widget.NewLabel(fmt.Sprintf("%s Formats:", categoryTitle))
		categoryLabel.TextStyle = fyne.TextStyle{Bold: true}
		formatContainer.Add(categoryLabel)

		for _, format := range formats.OutputFormats {
			formatContainer.Add(widget.NewLabel(fmt.Sprintf("  • %s", format)))
		}
	}

	// App settings
	appSettings := container.NewVBox(
		widget.NewLabel("📱 Application Settings"),
		widget.NewSeparator(),
		widget.NewCheck("🔔 Enable notifications", func(checked bool) {
			fmt.Printf("Notifications: %v\n", checked)
		}),
		widget.NewCheck("🌙 Dark mode", func(checked bool) {
			fmt.Printf("Dark mode: %v\n", checked)
		}),
		widget.NewCheck("📁 Remember last directory", func(checked bool) {
			fmt.Printf("Remember directory: %v\n", checked)
		}),
	)

	// Create a container with all three cards
	rightPanel := container.NewVBox(
		widget.NewCard("📊 Supported Formats", "", container.NewScroll(formatContainer)),
		widget.NewCard("⚙️ App Settings", "", appSettings),
	)

	return container.NewHSplit(
		widget.NewCard("🔧 Tool Status", "", container.NewScroll(toolStatus)),
		rightPanel,
	)
}

func (g *GUI) createDropArea() *widget.Card {
	dropLabel := widget.NewLabel("📁 Drag & drop files here\nor click to select files")
	dropLabel.Alignment = fyne.TextAlignCenter

	g.addFilesBtn = widget.NewButton("📂 Add Files", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				g.addFile(reader.URI().Path())
				reader.Close()
			}
		}, g.window)
	})

	g.clearBtn = widget.NewButton("🗑️ Clear", func() {
		g.files = make([]string, 0)
		g.fileList.Refresh()
		g.updateStatus("🗑️ Cleared all files")
	})

	g.convertBtn = widget.NewButton("🔄 Convert Files", g.convertFiles)
	g.convertBtn.Disable()

	return widget.NewCard("", "", container.NewVBox(dropLabel, g.addFilesBtn, g.clearBtn))
}

func (g *GUI) createDynamicConversionOptions() *fyne.Container {
	// Get valid output formats based on current files
	validFormats := g.getValidOutputFormats()

	// Output format selection
	g.outputFormat = widget.NewSelect(validFormats, func(selected string) {
		g.updateConversionOptions(selected)
	})

	if len(validFormats) > 0 {
		g.outputFormat.SetSelected(validFormats[0])
		g.updateConversionOptions(validFormats[0])
	}

	// Quality slider
	g.qualitySlider = widget.NewSlider(1, 100)
	g.qualitySlider.SetValue(95)
	g.qualityLabel = widget.NewLabel("🎨 Quality: 95")
	g.qualitySlider.OnChanged = func(value float64) {
		g.qualityLabel.SetText(fmt.Sprintf("🎨 Quality: %.0f", value))
	}

	// Output directory
	g.outputDir = widget.NewEntry()
	g.outputDir.SetPlaceHolder("📁 Same as input directory")
	browseBtn := widget.NewButton("📂 Browse", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				g.outputDir.SetText(uri.Path())
			}
		}, g.window)
	})

	outputContainer := container.NewBorder(nil, nil, nil, browseBtn, g.outputDir)

	// Format info
	formatInfo := widget.NewLabel("📋 Select files to see available formats")
	formatInfo.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		widget.NewLabel("📂 Output Format:"),
		g.outputFormat,
		widget.NewSeparator(),
		widget.NewLabel("🎨 Quality:"),
		container.NewHBox(g.qualitySlider, g.qualityLabel),
		widget.NewSeparator(),
		widget.NewLabel("📁 Output Directory:"),
		outputContainer,
		widget.NewSeparator(),
		widget.NewCard("📋 Format Information", "", formatInfo),
	)
}

func (g *GUI) getValidOutputFormats() []string {
	if len(g.files) == 0 {
		return []string{"mp4", "jpg", "mp3", "pdf"} // Default formats
	}

	// Get supported formats from converter
	supportedFormats := g.converter.GetSupportedFormats()

	// Determine input file types
	inputTypes := make(map[string]bool)
	for _, file := range g.files {
		ext := strings.ToLower(filepath.Ext(file))
		if ext != "" {
			inputTypes[ext] = true
		}
	}

	// Collect all possible output formats
	allFormats := make(map[string]bool)

	for category, formats := range supportedFormats {
		for _, inputFormat := range formats.InputFormats {
			if inputTypes["."+inputFormat] {
				// Add all output formats for this input format
				for _, outputFormat := range formats.OutputFormats {
					allFormats[outputFormat] = true
				}
				break
			}
		}
	}

	// Convert to sorted slice
	result := make([]string, 0, len(allFormats))
	for format := range allFormats {
		result = append(result, format)
	}

	// Sort formats (basic alphabetical)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i] > result[j] {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if len(result) == 0 {
		return []string{"mp4", "jpg", "mp3", "pdf"} // Fallback
	}

	return result
}

func (g *GUI) updateConversionOptions(selectedFormat string) {
	// Update format information
	formatInfo := g.getFormatInfo(selectedFormat)

	// Find the format info label and update it
	// This is a simplified approach - in a real app you'd store the reference
	g.updateStatus(fmt.Sprintf("📂 Selected format: %s", selectedFormat))
}

func (g *GUI) getFormatInfo(format string) string {
	// Provide information about the selected format
	descriptions := map[string]string{
		"mp4":  "MP4 - Modern video format with good compression",
		"avi":  "AVI - Classic video format, large file sizes",
		"mkv":  "MKV - Flexible container for multiple tracks",
		"mov":  "MOV - Apple QuickTime video format",
		"jpg":  "JPG - Compressed image format for photos",
		"png":  "PNG - Lossless image format with transparency",
		"gif":  "GIF - Animated image format",
		"mp3":  "MP3 - Compressed audio format",
		"wav":  "WAV - Uncompressed audio format",
		"flac": "FLAC - Lossless audio compression",
		"pdf":  "PDF - Document format for sharing",
		"txt":  "TXT - Plain text document",
	}

	if desc, exists := descriptions[format]; exists {
		return desc
	}
	return fmt.Sprintf("%s - Media format", strings.ToUpper(format))
}

func (g *GUI) createConversionOptions() *fyne.Container {
	// Output format selection
	formats := []string{
		"mp4", "avi", "mkv", "mov", "wmv", "flv", "webm", "gif",
		"jpg", "jpeg", "png", "bmp", "webp", "tiff",
		"mp3", "wav", "flac", "aac", "ogg", "m4a",
		"pdf", "txt", "html", "docx",
	}
	g.outputFormat = widget.NewSelect(formats, nil)
	g.outputFormat.SetSelected("mp4")

	// Quality slider
	g.qualitySlider = widget.NewSlider(1, 100)
	g.qualitySlider.SetValue(95)
	g.qualityLabel = widget.NewLabel("🎨 Quality: 95")
	g.qualitySlider.OnChanged = func(value float64) {
		g.qualityLabel.SetText(fmt.Sprintf("🎨 Quality: %.0f", value))
	}

	// Output directory
	g.outputDir = widget.NewEntry()
	g.outputDir.SetPlaceHolder("📁 Same as input directory")
	browseBtn := widget.NewButton("📂 Browse", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				g.outputDir.SetText(uri.Path())
			}
		}, g.window)
	})

	outputContainer := container.NewBorder(nil, nil, nil, browseBtn, g.outputDir)

	return container.NewVBox(
		widget.NewLabel("📂 Output Format:"),
		g.outputFormat,
		widget.NewSeparator(),
		widget.NewLabel("🎨 Quality:"),
		container.NewHBox(g.qualitySlider, g.qualityLabel),
		widget.NewSeparator(),
		widget.NewLabel("📁 Output Directory:"),
		outputContainer,
	)
}

func (g *GUI) createImageToolsSection() fyne.CanvasObject {
	// Image selection
	imageEntry := widget.NewEntry()
	imageEntry.SetPlaceHolder("📷 Select an image file...")

	selectImageBtn := widget.NewButton("📂 Browse", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				imageEntry.SetText(reader.URI().Path())
				reader.Close()
			}
		}, g.window)
	})

	// Tool options
	cropContainer := g.createCropTool(imageEntry)
	resizeContainer := g.createResizeTool(imageEntry)
	rotateContainer := g.createRotateTool(imageEntry)

	return container.NewVBox(
		widget.NewCard("📷 Select Image", "", container.NewHBox(imageEntry, selectImageBtn)),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewCard("✂️ Crop", "", cropContainer),
			widget.NewCard("📏 Resize", "", resizeContainer),
		),
		widget.NewSeparator(),
		widget.NewCard("🔄 Rotate", "", rotateContainer),
	)
}

func (g *GUI) createVideoToolsSection() fyne.CanvasObject {
	// Video selection
	videoEntry := widget.NewEntry()
	videoEntry.SetPlaceHolder("🎬 Select a video file...")

	selectVideoBtn := widget.NewButton("📂 Browse", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				videoEntry.SetText(reader.URI().Path())
				reader.Close()
			}
		}, g.window)
	})

	// Tool options
	frameContainer := g.createFrameTool(videoEntry)
	gifContainer := g.createGifTool(videoEntry)
	audioContainer := g.createAudioExtractionTool(videoEntry)
	infoContainer := g.createVideoInfoTool(videoEntry)

	return container.NewVBox(
		widget.NewCard("🎬 Select Video", "", container.NewHBox(videoEntry, selectVideoBtn)),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewCard("📸 Extract Frame", "", frameContainer),
			widget.NewCard("🎨 Convert to GIF", "", gifContainer),
		),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewCard("🎵 Extract Audio", "", audioContainer),
			widget.NewCard("ℹ️ Video Info", "", infoContainer),
		),
	)
}

func (g *GUI) createAudioToolsSection() fyne.CanvasObject {
	// Audio selection
	audioEntry := widget.NewEntry()
	audioEntry.SetPlaceHolder("🎵 Select an audio file...")

	selectAudioBtn := widget.NewButton("📂 Browse", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				audioEntry.SetText(reader.URI().Path())
				reader.Close()
			}
		}, g.window)
	})

	// Tool options
	convertContainer := g.createAudioConversionTool(audioEntry)
	infoContainer := g.createAudioInfoTool(audioEntry)

	return container.NewVBox(
		widget.NewCard("🎵 Select Audio", "", container.NewHBox(audioEntry, selectAudioBtn)),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewCard("🔄 Convert", "", convertContainer),
			widget.NewCard("ℹ️ Audio Info", "", infoContainer),
		),
	)
}

func (g *GUI) createToolStatus() fyne.CanvasObject {
	tools := converter.ValidateTools()
	toolStatus := container.NewVBox()

	for tool, available := range tools {
		status := "❌ Not Available"
		if available {
			status = "✅ Available"
		}
		toolStatus.Add(widget.NewLabel(fmt.Sprintf("%s: %s", tool, status)))
	}

	return toolStatus
}

// Tool creation functions
func (g *GUI) createCropTool(imageEntry *widget.Entry) fyne.CanvasObject {
	cropX := widget.NewEntry()
	cropX.SetText("0")
	cropY := widget.NewEntry()
	cropWidth := widget.NewEntry()
	cropWidth.SetText("100")
	cropHeight := widget.NewEntry()
	cropHeight = widget.NewEntry()
	cropHeight.SetText("100")

	return container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("X:"), cropX,
			widget.NewLabel("Y:"), cropY,
			widget.NewLabel("Width:"), cropWidth,
			widget.NewLabel("Height:"), cropHeight,
		),
		widget.NewButton("✂️ Crop", func() {
			g.cropImage(imageEntry.Text, cropX.Text, cropY.Text, cropWidth.Text, cropHeight.Text)
		}),
	)
}

func (g *GUI) createResizeTool(imageEntry *widget.Entry) fyne.CanvasObject {
	resizeWidth := widget.NewEntry()
	resizeWidth.SetText("800")
	resizeHeight := widget.NewEntry()
	resizeHeight.SetText("600")

	return container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Width:"), resizeWidth,
			widget.NewLabel("Height:"), resizeHeight,
		),
		widget.NewButton("📏 Resize", func() {
			g.resizeImage(imageEntry.Text, resizeWidth.Text, resizeHeight.Text)
		}),
	)
}

func (g *GUI) createRotateTool(imageEntry *widget.Entry) fyne.CanvasObject {
	rotateSelect := widget.NewSelect([]string{"90°", "180°", "270°"}, nil)

	return container.NewVBox(
		rotateSelect,
		widget.NewButton("🔄 Rotate", func() {
			g.rotateImage(imageEntry.Text, rotateSelect.Selected)
		}),
	)
}

func (g *GUI) createFrameTool(videoEntry *widget.Entry) fyne.CanvasObject {
	timestampEntry := widget.NewEntry()
	timestampEntry.SetText("00:00:05")
	timestampEntry.SetPlaceHolder("HH:MM:SS or seconds")

	frameWidth := widget.NewEntry()
	frameWidth.SetText("1920")
	frameHeight := widget.NewEntry()
	frameHeight = widget.NewEntry()
	frameHeight.SetText("1080")

	return container.NewVBox(
		widget.NewLabel("⏰ Timestamp:"),
		timestampEntry,
		container.NewGridWithColumns(2,
			widget.NewLabel("Width:"), frameWidth,
			widget.NewLabel("Height:"), frameHeight,
		),
		widget.NewButton("📸 Extract", func() {
			g.extractFrame(videoEntry.Text, timestampEntry.Text, frameWidth.Text, frameHeight.Text)
		}),
	)
}

func (g *GUI) createGifTool(videoEntry *widget.Entry) fyne.CanvasObject {
	gifFps := widget.NewEntry()
	gifFps.SetText("10")
	gifWidth := widget.NewEntry()
	gifHeight := widget.NewEntry()
	gifHeight.SetText("-1")

	return container.NewVBox(
		container.NewGridWithColumns(3,
			widget.NewLabel("FPS:"), gifFps,
			widget.NewLabel("Width:"), gifWidth,
			widget.NewLabel("Height:"), gifHeight,
		),
		widget.NewButton("🎨 Convert to GIF", func() {
			g.convertToGif(videoEntry.Text, gifFps.Text, gifWidth.Text, gifHeight.Text)
		}),
	)
}

func (g *GUI) createAudioExtractionTool(videoEntry *widget.Entry) fyne.CanvasObject {
	audioBitrate := widget.NewSelect([]string{"128k", "192k", "256k", "320k"}, nil)
	audioBitrate.SetSelected("192k")

	return container.NewVBox(
		widget.NewLabel("🎵 Bitrate:"),
		audioBitrate,
		widget.NewButton("🎵 Extract Audio", func() {
			g.extractAudio(videoEntry.Text, audioBitrate.Selected)
		}),
	)
}

func (g *GUI) createVideoInfoTool(videoEntry *widget.Entry) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewButton("ℹ️ Get Video Info", func() {
			g.getVideoInfo(videoEntry.Text)
		}),
	)
}

func (g *GUI) createAudioConversionTool(audioEntry *widget.Entry) fyne.CanvasObject {
	audioFormat := widget.NewSelect([]string{"mp3", "wav", "flac", "aac", "ogg"}, nil)
	audioFormat.SetSelected("mp3")

	return container.NewVBox(
		widget.NewLabel("🎵 Output Format:"),
		audioFormat,
		widget.NewButton("🔄 Convert Audio", func() {
			g.convertAudio(audioEntry.Text, audioFormat.Selected)
		}),
	)
}

func (g *GUI) createAudioInfoTool(audioEntry *widget.Entry) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewButton("ℹ️ Get Audio Info", func() {
			g.getAudioInfo(audioEntry.Text)
		}),
	)
}

// File operations
func (g *GUI) addFile(filename string) {
	g.files = append(g.files, filename)
	g.fileList.Refresh()
	g.convertBtn.Enable()
	g.updateStatus(fmt.Sprintf("📁 Added: %s", filepath.Base(filename)))
}

func (g *GUI) convertFiles() {
	if len(g.files) == 0 {
		dialog.ShowError(fmt.Errorf("no files selected"), g.window)
		return
	}

	g.updateStatus("🔄 Converting files...")
	g.progressBar.SetValue(0)

	outputFormat := g.outputFormat.Selected
	quality := fmt.Sprintf("%.0f", g.qualitySlider.Value)
	outputDir := g.outputDir.Text

	for i, file := range g.files {
		outputPath := g.generateOutputPath(file, outputFormat, outputDir)

		req := converter.ConversionRequest{
			InputPath:  file,
			OutputPath: outputPath,
			Options:    map[string]string{"quality": quality},
		}

		err := g.converter.Convert(req)
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to convert %s: %w", file, err), g.window)
			g.updateStatus("❌ Conversion failed")
			return
		}

		// Update progress
		progress := float64(i+1) / float64(len(g.files))
		g.progressBar.SetValue(progress)
	}

	g.updateStatus("✅ Conversion completed!")
	dialog.ShowInformation("Success", "All files converted successfully!", g.window)
}

func (g *GUI) cropImage(imagePath, x, y, width, height string) {
	if imagePath == "" {
		dialog.ShowError(fmt.Errorf("please select an image file"), g.window)
		return
	}

	outputPath := strings.TrimSuffix(imagePath, filepath.Ext(imagePath)) + "_cropped.jpg"

	req := image.CropRequest{
		InputPath:  imagePath,
		OutputPath: outputPath,
		X:          parseInt(x),
		Y:          parseInt(y),
		Width:      parseInt(width),
		Height:     parseInt(height),
		Quality:    95,
	}

	err := g.imageProcessor.Crop(req)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to crop image: %w", err), g.window)
		return
	}

	dialog.ShowInformation("Success", fmt.Sprintf("✂️ Image cropped to: %s", outputPath), g.window)
}

func (g *GUI) resizeImage(imagePath, width, height string) {
	if imagePath == "" {
		dialog.ShowError(fmt.Errorf("please select an image file"), g.window)
		return
	}

	outputPath := strings.TrimSuffix(imagePath, filepath.Ext(imagePath)) + "_resized.jpg"

	req := image.ResizeRequest{
		InputPath:  imagePath,
		OutputPath: outputPath,
		Width:      parseInt(width),
		Height:     parseInt(height),
		Quality:    95,
	}

	err := g.imageProcessor.Resize(req)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to resize image: %w", err), g.window)
		return
	}

	dialog.ShowInformation("Success", fmt.Sprintf("📏 Image resized to: %s", outputPath), g.window)
}

func (g *GUI) rotateImage(imagePath, angle string) {
	if imagePath == "" {
		dialog.ShowError(fmt.Errorf("please select an image file"), g.window)
		return
	}

	outputPath := strings.TrimSuffix(imagePath, filepath.Ext(imagePath)) + "_rotated.jpg"

	var degrees float64
	switch angle {
	case "90°":
		degrees = 90
	case "180°":
		degrees = 180
	case "270°":
		degrees = 270
	default:
		degrees = 90
	}

	err := g.imageProcessor.Rotate(imagePath, outputPath, degrees, 95)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to rotate image: %w", err), g.window)
		return
	}

	dialog.ShowInformation("Success", fmt.Sprintf("🔄 Image rotated to: %s", outputPath), g.window)
}

func (g *GUI) extractFrame(videoPath, timestamp, width, height string) {
	if videoPath == "" {
		dialog.ShowError(fmt.Errorf("please select a video file"), g.window)
		return
	}

	outputPath := strings.TrimSuffix(videoPath, filepath.Ext(videoPath)) + "_frame.jpg"

	req := video.ExtractRequest{
		VideoPath:  videoPath,
		OutputPath: outputPath,
		Timestamp:  timestamp,
	}

	if width != "" && height != "" {
		fmt.Sscanf(width, "%d", &req.Width)
		fmt.Sscanf(height, "%d", &req.Height)
	}

	err := g.frameExtractor.ExtractFrame(req)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to extract frame: %w", err), g.window)
		return
	}

	dialog.ShowInformation("Success", fmt.Sprintf("📸 Frame extracted to: %s", outputPath), g.window)
}

func (g *GUI) convertToGif(videoPath, fps, width, height string) {
	if videoPath == "" {
		dialog.ShowError(fmt.Errorf("please select a video file"), g.window)
		return
	}

	outputPath := strings.TrimSuffix(videoPath, filepath.Ext(videoPath)) + ".gif"

	req := converter.ConversionRequest{
		InputPath:  videoPath,
		OutputPath: outputPath,
		Options: map[string]string{
			"fps":    fps,
			"width":  width,
			"height": height,
		},
	}

	err := g.converter.Convert(req)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to convert to GIF: %w", err), g.window)
		return
	}

	dialog.ShowInformation("Success", fmt.Sprintf("🎨 Video converted to GIF: %s", outputPath), g.window)
}

func (g *GUI) extractAudio(videoPath, bitrate string) {
	if videoPath == "" {
		dialog.ShowError(fmt.Errorf("please select a video file"), g.window)
		return
	}

	outputPath := strings.TrimSuffix(videoPath, filepath.Ext(videoPath)) + ".mp3"

	req := converter.ConversionRequest{
		InputPath:  videoPath,
		OutputPath: outputPath,
		Options: map[string]string{
			"bitrate": bitrate,
		},
	}

	err := g.converter.Convert(req)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to extract audio: %w", err), g.window)
		return
	}

	dialog.ShowInformation("Success", fmt.Sprintf("🎵 Audio extracted to: %s", outputPath), g.window)
}

func (g *GUI) convertAudio(audioPath, format string) {
	if audioPath == "" {
		dialog.ShowError(fmt.Errorf("please select an audio file"), g.window)
		return
	}

	outputPath := strings.TrimSuffix(audioPath, filepath.Ext(audioPath)) + "." + format

	req := converter.ConversionRequest{
		InputPath:  audioPath,
		OutputPath: outputPath,
		Options:    map[string]string{},
	}

	err := g.converter.Convert(req)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to convert audio: %w", err), g.window)
		return
	}

	dialog.ShowInformation("Success", fmt.Sprintf("🔄 Audio converted to: %s", outputPath), g.window)
}

func (g *GUI) getVideoInfo(videoPath string) {
	if videoPath == "" {
		dialog.ShowError(fmt.Errorf("please select a video file"), g.window)
		return
	}

	info, err := g.frameExtractor.GetVideoInfo(videoPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to get video info: %w", err), g.window)
		return
	}

	infoText := fmt.Sprintf(
		"⏱️ Duration: %s\n📐 Dimensions: %dx%d\n🎬 Codec: %s\n🎞 Frame Rate: %s",
		info.Duration, info.Width, info.Height, info.Codec, info.FrameRate,
	)

	dialog.ShowInformation("Video Information", infoText, g.window)
}

func (g *GUI) getAudioInfo(audioPath string) {
	if audioPath == "" {
		dialog.ShowError(fmt.Errorf("please select an audio file"), g.window)
		return
	}

	// This is a placeholder - you'd implement actual audio info extraction
	infoText := fmt.Sprintf(
		"🎵 Audio File: %s\n📊 Format: %s\n💾 Size: %s bytes",
		filepath.Base(audioPath),
		filepath.Ext(audioPath),
		"Unknown",
	)

	dialog.ShowInformation("Audio Information", infoText, g.window)
}

func (g *GUI) generateOutputPath(inputPath, outputFormat, outputDir string) string {
	ext := "." + outputFormat
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))

	if outputDir != "" {
		return filepath.Join(outputDir, baseName+ext)
	}

	return filepath.Join(filepath.Dir(inputPath), baseName+ext)
}

func (g *GUI) updateStatus(message string) {
	g.statusLabel.SetText(message)
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

func (g *GUI) playMedia(filePath string) {
	if filePath == "" {
		dialog.ShowError(fmt.Errorf("no file selected for playback"), g.window)
		return
	}

	err := g.mediaPlayer.Play(filePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to play media: %w", err), g.window)
		return
	}

	g.updateStatus(fmt.Sprintf("🎵 Playing: %s", filepath.Base(filePath)))
}

func (g *GUI) showPreview(filePath string) {
	if filePath == "" {
		dialog.ShowError(fmt.Errorf("no file selected for preview"), g.window)
		return
	}

	previewInfo, err := g.mediaPlayer.GeneratePreview(filePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to generate preview: %w", err), g.window)
		return
	}

	// Clear previous preview
	g.previewArea.Objects = []fyne.CanvasObject{}

	// Create preview content
	previewContent := container.NewVBox()

	// File info
	fileInfo := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("📁 File: %s", filepath.Base(filePath))),
		widget.NewLabel(fmt.Sprintf("📊 Size: %s", previewInfo.FileSize)),
		widget.NewLabel(fmt.Sprintf("🎯 Format: %s", previewInfo.Format)),
	)

	if previewInfo.Duration != "" {
		fileInfo.Add(widget.NewLabel(fmt.Sprintf("⏱️ Duration: %s", previewInfo.Duration)))
	}

	if previewInfo.Resolution != "" {
		fileInfo.Add(widget.NewLabel(fmt.Sprintf("📐 Resolution: %s", previewInfo.Resolution)))
	}

	if previewInfo.Title != "" {
		fileInfo.Add(widget.NewLabel(fmt.Sprintf("🎵 Title: %s", previewInfo.Title)))
	}

	if previewInfo.Artist != "" {
		fileInfo.Add(widget.NewLabel(fmt.Sprintf("👤 Artist: %s", previewInfo.Artist)))
	}

	previewContent.Add(widget.NewCard("📋 File Information", "", fileInfo))

	// Add thumbnail if available
	if previewInfo.Thumbnail != "" {
		thumbnailImg := widget.NewIcon(theme.FileIcon())
		thumbnailCard := widget.NewCard("🖼️ Thumbnail", "", thumbnailImg)
		previewContent.Add(thumbnailCard)
	}

	// Action buttons
	actions := container.NewHBox(
		widget.NewButton("▶️ Play", func() { g.playMedia(filePath) }),
		widget.NewButton("📁 Open Folder", func() { g.openFileFolder(filePath) }),
	)

	previewContent.Add(widget.NewCard("🎬 Actions", "", actions))

	// Add to preview area
	g.previewArea.Add(previewContent)
	g.previewArea.Refresh()

	g.updateStatus(fmt.Sprintf("📸 Preview: %s", filepath.Base(filePath)))
}

func (g *GUI) openFileFolder(filePath string) {
	folderPath := filepath.Dir(filePath)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", folderPath)
	case "darwin":
		cmd = exec.Command("open", folderPath)
	default:
		cmd = exec.Command("xdg-open", folderPath)
	}

	cmd.Start()
}

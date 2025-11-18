package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"goverter/pkg/converter"
	"goverter/pkg/image"
	"goverter/pkg/video"
)

var (
	inputFile    string
	outputFile   string
	quality      string
	timestamp    string
	width        int
	height       int
	cropX        int
	cropY        int
	cropWidth    int
	cropHeight   int
	bulkDir      string
	outputFormat string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "govert",
		Short: "A versatile file conversion and media processing tool",
		Long: `Goverter is a CLI tool for converting files between formats,
processing images and videos, and handling bulk operations.`,
	}

	// Convert command
	var convertCmd = &cobra.Command{
		Use:   "convert [format]",
		Short: "Convert files between formats",
		Args:  cobra.ExactArgs(1),
		Run:   runConvert,
	}
	convertCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path")
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path")
	convertCmd.Flags().StringVarP(&quality, "quality", "q", "", "Quality setting (CRF for video, 1-100 for images)")
	convertCmd.Flags().StringVarP(&bulkDir, "bulk", "b", "", "Bulk convert all files in directory")
	convertCmd.Flags().StringVarP(&outputFormat, "format", "f", "", "Output format for bulk conversion")

	// Frame command
	var frameCmd = &cobra.Command{
		Use:   "frame [timestamp]",
		Short: "Extract a frame from video at specific timestamp",
		Args:  cobra.ExactArgs(1),
		Run:   runFrame,
	}
	frameCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input video file path")
	frameCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output image file path")
	frameCmd.Flags().IntVar(&width, "width", 0, "Output width (optional)")
	frameCmd.Flags().IntVar(&height, "height", 0, "Output height (optional)")

	// Crop command
	var cropCmd = &cobra.Command{
		Use:   "crop [x] [y] [width] [height]",
		Short: "Crop an image",
		Args:  cobra.ExactArgs(4),
		Run:   runCrop,
	}
	cropCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input image file path")
	cropCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output image file path")
	cropCmd.Flags().StringVarP(&quality, "quality", "q", "95", "JPEG quality (1-100)")

	// Resize command
	var resizeCmd = &cobra.Command{
		Use:   "resize [width] [height]",
		Short: "Resize an image",
		Args:  cobra.ExactArgs(2),
		Run:   runResize,
	}
	resizeCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input image file path")
	resizeCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output image file path")
	resizeCmd.Flags().StringVarP(&quality, "quality", "q", "95", "JPEG quality (1-100)")

	// Info command
	var infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Get information about a media file",
		Run:   runInfo,
	}
	infoCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path")

	// Add subcommands
	rootCmd.AddCommand(convertCmd, frameCmd, cropCmd, resizeCmd, infoCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runConvert(cmd *cobra.Command, args []string) {
	if bulkDir != "" {
		runBulkConvert(args[0])
		return
	}

	if inputFile == "" || outputFile == "" {
		fmt.Println("Error: Both --input and --output flags are required")
		return
	}

	c := converter.NewConverter()

	options := make(map[string]string)
	if quality != "" {
		options["quality"] = quality
	}

	req := converter.ConversionRequest{
		InputPath:  inputFile,
		OutputPath: outputFile,
		Options:    options,
	}

	if err := c.Convert(req); err != nil {
		fmt.Printf("Error converting file: %v\n", err)
		return
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
}

func runBulkConvert(targetFormat string) {
	if bulkDir == "" || outputFormat == "" {
		fmt.Println("Error: Both --bulk and --format flags are required for bulk conversion")
		return
	}

	// Implementation for bulk conversion
	fmt.Printf("Bulk conversion not yet implemented for directory: %s\n", bulkDir)
}

func runFrame(cmd *cobra.Command, args []string) {
	timestamp := args[0]

	if inputFile == "" || outputFile == "" {
		fmt.Println("Error: Both --input and --output flags are required")
		return
	}

	fe := video.NewFrameExtractor()
	req := video.ExtractRequest{
		VideoPath:  inputFile,
		OutputPath: outputFile,
		Timestamp:  timestamp,
		Width:      width,
		Height:     height,
	}

	if err := fe.ExtractFrame(req); err != nil {
		fmt.Printf("Error extracting frame: %v\n", err)
		return
	}

	fmt.Printf("Successfully extracted frame at %s to %s\n", timestamp, outputFile)
}

func runCrop(cmd *cobra.Command, args []string) {
	if inputFile == "" || outputFile == "" {
		fmt.Println("Error: Both --input and --output flags are required")
		return
	}

	processor := image.NewProcessor()
	req := image.CropRequest{
		InputPath:  inputFile,
		OutputPath: outputFile,
		X:          parseInt(args[0]),
		Y:          parseInt(args[1]),
		Width:      parseInt(args[2]),
		Height:     parseInt(args[3]),
		Quality:    parseInt(quality),
	}

	if err := processor.Crop(req); err != nil {
		fmt.Printf("Error cropping image: %v\n", err)
		return
	}

	fmt.Printf("Successfully cropped %s to %s\n", inputFile, outputFile)
}

func runResize(cmd *cobra.Command, args []string) {
	if inputFile == "" || outputFile == "" {
		fmt.Println("Error: Both --input and --output flags are required")
		return
	}

	processor := image.NewProcessor()
	req := image.ResizeRequest{
		InputPath:  inputFile,
		OutputPath: outputFile,
		Width:      parseInt(args[0]),
		Height:     parseInt(args[1]),
		Quality:    parseInt(quality),
	}

	if err := processor.Resize(req); err != nil {
		fmt.Printf("Error resizing image: %v\n", err)
		return
	}

	fmt.Printf("Successfully resized %s to %s\n", inputFile, outputFile)
}

func runInfo(cmd *cobra.Command, args []string) {
	if inputFile == "" {
		fmt.Println("Error: --input flag is required")
		return
	}

	ext := getExt(inputFile)

	switch ext {
	case ".mp4", ".avi", ".mkv", ".mov", ".wmv":
		fe := video.NewFrameExtractor()
		info, err := fe.GetVideoInfo(inputFile)
		if err != nil {
			fmt.Printf("Error getting video info: %v\n", err)
			return
		}
		fmt.Printf("Video Information:\n")
		fmt.Printf("  Duration: %s\n", info.Duration)
		fmt.Printf("  Dimensions: %dx%d\n", info.Width, info.Height)
		fmt.Printf("  Codec: %s\n", info.Codec)
		fmt.Printf("  Frame Rate: %s\n", info.FrameRate)

	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		processor := image.NewProcessor()
		info, err := processor.GetImageInfo(inputFile)
		if err != nil {
			fmt.Printf("Error getting image info: %v\n", err)
			return
		}
		fmt.Printf("Image Information:\n")
		fmt.Printf("  Dimensions: %dx%d\n", info.Width, info.Height)
		fmt.Printf("  Format: %s\n", info.Format)
		fmt.Printf("  File Size: %d bytes\n", info.FileSize)

	default:
		fmt.Printf("Unsupported file type for info: %s\n", ext)
	}
}

func getExt(filename string) string {
	for i := len(filename) - 1; i >= 0 && filename[i] != '/'; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ""
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}

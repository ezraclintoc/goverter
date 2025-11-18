package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"goverter/pkg/converter"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Goverter - File Converter")
	myWindow.Resize(fyne.NewSize(400, 300))

	// Check dependencies
	tools := converter.ValidateTools()
	statusText := "Goverter - File Converter\n\nDependencies Status:\n"

	for tool, available := range tools {
		status := "❌ Not Available"
		if available {
			status = "✅ Available"
		}
		statusText += fmt.Sprintf("%s: %s\n", tool, status)
	}

	statusText += "\nUse CLI for full functionality:\n./goverter-cli --help"

	label := widget.NewLabel(statusText)
	label.Wrapping = fyne.TextWrapWord

	myWindow.SetContent(label)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

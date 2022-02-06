package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID("nl.ganzeinfach.apps.bbtrader")
	mainWindow := app.NewWindow("Trade Manager")
	mainContent := makeMainContent()
	width := app.Preferences().FloatWithFallback("width", 800)
	height := app.Preferences().FloatWithFallback("height", 900)
	mainWindow.Resize(fyne.Size{Width: float32(width), Height: float32(height)})
	mainWindow.SetCloseIntercept(func() {
        log.Printf("Window W: %.0f x H: %.0f", mainWindow.Canvas().Size().Width, mainWindow.Canvas().Size().Height)
		app.Preferences().SetFloat("width", float64(mainWindow.Canvas().Size().Width))
		app.Preferences().SetFloat("height", float64(mainWindow.Canvas().Size().Height))
        mainWindow.Close()
	})

	mainWindow.CenterOnScreen()
	mainWindow.SetContent(mainContent)
	mainWindow.ShowAndRun()
}

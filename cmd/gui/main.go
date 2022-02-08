package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID("nl.ganzeinfach.apps.bbtrader")
    app.Settings().SetTheme(&myTheme{})
	mainWindow := app.NewWindow("Trade Manager")
	mainContent := makeMainContent()
	width := app.Preferences().FloatWithFallback("width", 800)
	height := app.Preferences().FloatWithFallback("height", 900)
	mainWindow.Resize(fyne.Size{Width: float32(width), Height: float32(height)})
	mainWindow.SetCloseIntercept(func() {
		app.Preferences().SetFloat("width", float64(mainWindow.Canvas().Size().Width))
		app.Preferences().SetFloat("height", float64(mainWindow.Canvas().Size().Height))
		mainWindow.Close()
	})

    // TODO: also remember position
	mainWindow.CenterOnScreen()
    
	mainWindow.SetContent(mainContent)
	mainWindow.ShowAndRun()
}

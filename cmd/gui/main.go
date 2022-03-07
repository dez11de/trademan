package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var mainWindow fyne.Window

func main() {
	app := app.NewWithID("nl.ganzeinfach.apps.bbtrader")
	app.Settings().SetTheme(&myTheme{})
	dbName, err := getDBName()
	if err != nil {
		log.Panicf("unable to connect to server: %s", err)
	}
	mainWindow = app.NewWindow(fmt.Sprintf("Trade Manager (%s)", dbName))
	mainContent := makeMainContent()
	width := app.Preferences().FloatWithFallback("width", 850)
	height := app.Preferences().FloatWithFallback("height", 1000)
	mainWindow.Resize(fyne.Size{Width: float32(width), Height: float32(height)})
	mainWindow.SetCloseIntercept(func() {
		app.Preferences().SetFloat("width", float64(mainWindow.Canvas().Size().Width))
		app.Preferences().SetFloat("height", float64(mainWindow.Canvas().Size().Height))
		mainWindow.Close()
	})

	mainWindow.SetContent(mainContent)
	mainWindow.CenterOnScreen() // TODO: also remember position
	mainWindow.ShowAndRun()
}

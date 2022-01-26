package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// TODO: come up with a better name that doesn't focus on ByBit cause you never know
	app := app.NewWithID("bbtrader")
	mainWindow := app.NewWindow("ByBit Trade Manager")
	// TODO: would it hurt if these are global?
	mainContent := makeMainContent()
	// TODO: store and restore size and position settings
	mainWindow.Resize(fyne.Size{Width: 800, Height: 600})
	mainWindow.CenterOnScreen()
	mainWindow.SetContent(mainContent)
	mainWindow.ShowAndRun()
}

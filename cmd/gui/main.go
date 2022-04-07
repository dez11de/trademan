package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var trademanCfg trademanConfig
var BaseURL string
var mainWindow fyne.Window
var a fyne.App

func main() {
	err := readConfig(&trademanCfg)
    BaseURL = fmt.Sprintf("http://%s:%s/api/v1/", trademanCfg.RESTServer.Host, trademanCfg.RESTServer.Port)
	if err != nil {
		log.Fatalf("Error reading configuration: %s", err)
	}

	a = app.NewWithID("nl.ganzeinfach.apps.bbtrader")
	a.Settings().SetTheme(&myTheme{})
	dbName, err := getDBName()
	if err != nil {
		log.Panicf("unable to connect to server: %s", err)
	}
	mainWindow = a.NewWindow(fmt.Sprintf("Trade Manager (%s)", dbName))
	mainContent := makeMainContent()
	width := a.Preferences().FloatWithFallback("width", 850)
	height := a.Preferences().FloatWithFallback("height", 1000)
	mainWindow.Resize(fyne.Size{Width: float32(width), Height: float32(height)})
	mainWindow.SetCloseIntercept(func() {
		a.Preferences().SetFloat("width", float64(mainWindow.Canvas().Size().Width))
		a.Preferences().SetFloat("height", float64(mainWindow.Canvas().Size().Height))
		mainWindow.Close()
	})

	mainWindow.SetContent(mainContent)
	mainWindow.CenterOnScreen() // TODO: also remember position
	mainWindow.ShowAndRun()
}

package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

var trademanCfg trademanConfig
var BaseURL string
var mainWindow fyne.Window
var a fyne.App

type tradeMan struct {
	pairs             []cryptodb.Pair
    plans             []cryptodb.Plan
	assessmentOptions map[string][]string
}

var tm tradeMan

type active struct {
	pair       cryptodb.Pair
	plan       cryptodb.Plan
	orders     []cryptodb.Order
	assessment cryptodb.Assessment

	List                *widget.List
	statisticsContainer *fyne.Container
}

var act active

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

	tm.pairs, err = getPairs()
	if err != nil {
        log.Panicf("unable to load pairs: %s", err)
	}
	tm.plans, err = getPlans()
	if err != nil {
        log.Panicf("unable to load plans: %s", err)
	}

	tm.assessmentOptions, err = getAssessmentOptions()
	if err != nil {
        log.Panicf("unable to load assessment options: %s", err)
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

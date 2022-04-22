package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

var trademanCfg trademanConfig
var BaseURL string

type tradeMan struct {
	pairs         []cryptodb.Pair
	plans         []cryptodb.Plan
	reviewOptions map[string][]string

	pair   cryptodb.Pair
	plan   cryptodb.Plan
	orders []cryptodb.Order
	review cryptodb.Review
}

var tm tradeMan

var ui struct {
	app          fyne.App
	mainWindow   fyne.Window
	showArchived bool

	noPlanSelectedContainer *fyne.Container
	planListSplit           *container.Split
	planList                *widget.List
}

func main() {
	err := readConfig(&trademanCfg)
	BaseURL = fmt.Sprintf("http://%s:%s/api/v1/", trademanCfg.RESTServer.Host, trademanCfg.RESTServer.Port)
	if err != nil {
		log.Fatalf("Error reading configuration: %s", err)
	}
	dbName, err := getDBName()
	if err != nil {
		log.Panicf("unable to connect to server: %s", err)
	}
	ui.app = app.NewWithID("nl.ganzeinfach.apps.bbtrader")
	ui.app.Settings().SetTheme(&myTheme{})
	ui.mainWindow = ui.app.NewWindow(fmt.Sprintf("Trade Manager (%s)", dbName))
	ui.showArchived = ui.app.Preferences().BoolWithFallback("showArchived", false)

	tm.pairs, err = getPairs()
	if err != nil {
		log.Panicf("unable to load pairs: %s", err)
	}
	tm.plans, err = getPlans()
	if err != nil {
		log.Panicf("unable to load plans: %s", err)
	}

	tm.reviewOptions, err = getReviewOptions()
	if err != nil {
		log.Panicf("unable to load review options: %s", err)
	}

	mainContent := makeMainContent()
	width := ui.app.Preferences().FloatWithFallback("main-width", 815.0)
	height := ui.app.Preferences().FloatWithFallback("main-height", 610.0)
	ui.mainWindow.Resize(fyne.Size{Width: float32(width), Height: float32(height)})
	ui.mainWindow.SetCloseIntercept(func() {
		ui.app.Preferences().SetBool("showArchived", ui.showArchived)
		ui.app.Preferences().SetFloat("main-width", float64(ui.mainWindow.Canvas().Size().Width))
		ui.app.Preferences().SetFloat("main-height", float64(ui.mainWindow.Canvas().Size().Height))
		ui.mainWindow.Close()
	})

	ui.mainWindow.SetContent(mainContent)
	ui.mainWindow.CenterOnScreen() // TODO: also remember position

	ui.mainWindow.SetMainMenu(makeMainMenu())

	ui.mainWindow.ShowAndRun()
}

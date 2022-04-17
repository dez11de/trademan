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

type tradeMan struct {
	pairs         []cryptodb.Pair
	plans         []cryptodb.Plan
	reviewOptions map[string][]string
}

var tm tradeMan

type active struct {
	pair   cryptodb.Pair
	plan   cryptodb.Plan
	orders []cryptodb.Order
	review cryptodb.Review
}

var act active

var application struct {
	fa       fyne.App
	mw       fyne.Window
	planList *widget.List
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

	application.fa = app.NewWithID("nl.ganzeinfach.apps.bbtrader")
	application.fa.Settings().SetTheme(&myTheme{})
	application.mw = application.fa.NewWindow(fmt.Sprintf("Trade Manager (%s)", dbName))
	mainContent := makeMainContent()
	width := application.fa.Preferences().FloatWithFallback("main-width", 815.0)
	height := application.fa.Preferences().FloatWithFallback("main-height", 610.0)
	application.mw.Resize(fyne.Size{Width: float32(width), Height: float32(height)})
	application.mw.SetCloseIntercept(func() {
		application.fa.Preferences().SetFloat("main-width", float64(application.mw.Canvas().Size().Width))
		application.fa.Preferences().SetFloat("main-height", float64(application.mw.Canvas().Size().Height))
		application.mw.Close()
	})

	application.mw.SetContent(mainContent)
	application.mw.CenterOnScreen() // TODO: also remember position
	application.mw.ShowAndRun()
}

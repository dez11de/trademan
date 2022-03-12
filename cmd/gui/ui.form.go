package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
)

func (pf *planForm) gatherSetup() cryptodb.Setup {
	// TODO: check for errors
	ui.activePlan.PairID = ui.activePair.ID
	ui.activePlan.Direction.Scan(pf.directionItem.Widget.(*widget.RadioGroup).Selected)
	ui.activePlan.Risk = decimal.RequireFromString(pf.riskItem.Widget.(*widget.Entry).Text)
	ui.activeOrders[cryptodb.MarketStopLoss].Price = decimal.RequireFromString(pf.stopLossItem.Widget.(*widget.Entry).Text)
	ui.activeOrders[cryptodb.Entry].Price = decimal.RequireFromString(pf.entryItem.Widget.(*widget.Entry).Text)
	ui.activePlan.TakeProfitStrategy.Scan(pf.TPStratItem.Widget.(*widget.Select).Selected)

	for i := 0; i < cryptodb.MaxTakeProfits; i++ {
		tempPrice, err := decimal.NewFromString(pf.takeProfitItems[i].Widget.(*widget.Entry).Text)
		if err == nil {
			ui.activeOrders[3+i].Price = tempPrice
		} else {
			ui.activeOrders[3+i].Price = decimal.Zero
		}
	}
	if pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL != nil {
		ui.activePlan.TradingViewPlan = pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL.String()
	}

	return cryptodb.Setup{Plan: ui.activePlan, Orders: ui.activeOrders}
}

func (pf *planForm) okAction() {
	setup := pf.gatherSetup()
	setup, err := storeSetup(setup)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}
	ui.activePlan = setup.Plan
	ui.activeOrders = setup.Orders
	ui.Plans, _ = getPlans()
	ui.List.Refresh()
}

func (pf *planForm) cancelAction() {

}

func (pf *planForm) executeAction() {
	setup := pf.gatherSetup()
	setup, err := storeSetup(setup)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}

	plan, err := executePlan(setup.Plan.ID)
	ui.activePlan = plan
	ui.activeOrders = setup.Orders
	ui.Plans, _ = getPlans()
    pf.form.Refresh()
	ui.List.Refresh()
}

func (pf *planForm) logAction() {
	entries, err := getLogs(ui.activePlan.ID)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}

	logFile := widget.NewRichText()
	for _, e := range entries {
		seg := &widget.TextSegment{
			Style: widget.RichTextStyle{},
			Text:  fmt.Sprintf("%s  %s", e.CreatedAt.Format("2006-01-02 15:04:05"), e.Text),
		}
		logFile.Segments = append(logFile.Segments, seg)
	}

	logWindow := widget.NewPopUp(logFile, mainWindow.Canvas())
	logAnimation := canvas.NewSizeAnimation(
		fyne.NewSize(mainWindow.Canvas().Size().Width-2*50.0, 0),
		fyne.NewSize(mainWindow.Canvas().Size().Width-2*50.0, mainWindow.Canvas().Size().Height-1*50.0),
		50*time.Millisecond,
		func(s fyne.Size) {
			logWindow.Resize(s)
		})

	logWindow.Resize(fyne.NewSize(0, 0))
	logWindow.ShowAtPosition(fyne.Position{X: 50, Y: 0})
	logAnimation.Start()
}

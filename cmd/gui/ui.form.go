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

func (pf *planForm) gatherPlan() cryptodb.Plan {
	// TODO: check for errors?
	ui.activePlan.PairID = ui.activePair.ID
	ui.activePlan.Direction.Scan(pf.directionItem.Widget.(*widget.RadioGroup).Selected)
	ui.activePlan.Risk = decimal.RequireFromString(pf.riskItem.Widget.(*FloatEntry).Text)
	ui.activePlan.TakeProfitStrategy.Scan(pf.TPStratItem.Widget.(*widget.Select).Selected)
	if pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL != nil {
		ui.activePlan.TradingViewPlan = pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL.String()
	}
	return ui.activePlan
}

func (pf *planForm) gatherOrders() []cryptodb.Order {
	ui.activeOrders[cryptodb.MarketStopLoss].Price = decimal.RequireFromString(pf.stopLossItem.Widget.(*FloatEntry).Text)
	ui.activeOrders[cryptodb.Entry].Price = decimal.RequireFromString(pf.entryItem.Widget.(*FloatEntry).Text)
	for i := 0; i < cryptodb.MaxTakeProfits; i++ {
		tempPrice, err := decimal.NewFromString(pf.takeProfitItems[i].Widget.(*FloatEntry).Text)
		if err == nil {
			ui.activeOrders[3+i].Price = tempPrice
		} else {
			ui.activeOrders[3+i].Price = decimal.Zero
		}
	}
	return ui.activeOrders
}

func (pf *planForm) saveSetup() {
	plan := pf.gatherPlan()
	updatedPlan, err := savePlan(plan)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}

	orders := pf.gatherOrders()
	if orders[cryptodb.MarketStopLoss].PlanID == 0 {
		for i := range orders {
			orders[i].PlanID = updatedPlan.ID
		}
	}

	_, err = saveOrders(orders)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}

    ui.activeAssessment.PlanID = updatedPlan.ID
    _, err = saveAssessment(ui.activeAssessment)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}

	pf.FillForm(updatedPlan)
}

func (pf *planForm) okAction() {
	pf.saveSetup()

	ui.Plans, _ = getPlans()
	ui.List.Refresh()
}

func (pf *planForm) undoAction() {
	reloadedPlan, _ := getPlan(ui.activePlan.ID)
	pf.FillForm(reloadedPlan)
}

func (pf *planForm) executeAction() {
	pf.saveSetup()
	executePlan(ui.activePlan.ID)

	ui.Plans, _ = getPlans()
	ui.List.Refresh()
}

func (pf *planForm) historyAction() {
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

func (pf *planForm) assessAction() {
    pf.assessmentForm = NewAssessmentForm(ui.activePlan)
    pf.assessmentForm.window = a.NewWindow("Assessment")
    pf.assessmentForm.window.SetContent(pf.assessmentForm.container)
    pf.assessmentForm.window.Show()
}

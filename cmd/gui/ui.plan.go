package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
)

type planForm struct {
	leftForm  *widget.Form
	rightForm *widget.Form

	pairItem            *widget.FormItem
	directionItem       *widget.FormItem
	riskItem            *widget.FormItem
	TPStratItem         *widget.FormItem
	stopLossItem        *widget.FormItem
	entryItem           *widget.FormItem
	takeProfitItems     [cryptodb.MaxTakeProfits]*widget.FormItem
	notesItem           *widget.FormItem
	tradingViewPlanItem *widget.FormItem
}

var pf planForm

func NewPlanContainer() *fyne.Container {
	pf.leftForm = widget.NewForm()
	pf.rightForm = widget.NewForm()
	bottomForm := widget.NewForm()

	pf.leftForm.AppendItem(pf.makePairItem())
	pf.leftForm.AppendItem(pf.makeDirectionItem())
	pf.leftForm.AppendItem(pf.makeRiskItem())
	pf.leftForm.AppendItem(pf.makeTradingViewItem())
	pf.leftForm.AppendItem(pf.makeTakeProfitStrategyItem())
	pf.rightForm.AppendItem(pf.makeStopLossItem())
	pf.rightForm.AppendItem(pf.makeEntryItem())

	for i := 0; i <= cryptodb.MaxTakeProfits-1; i++ {
		pf.rightForm.AppendItem(pf.makeTakeProfitItem(i, tm.pair.PriceScale, tm.pair.Price.Tick))
	}
	pf.setQuoteCurrency()
	pf.setPriceScale()

	bottomForm.AppendItem(pf.makeNotesItem())
	totalContainer := container.NewVBox(pf.makeStatContainer(),
		container.New(layout.NewGridLayoutWithColumns(2), pf.leftForm, pf.rightForm),
        bottomForm,
		pf.makeToolBar())

	return totalContainer
}

func (pf *planForm) gatherPlan() {
	tm.plan.PairID = tm.pair.ID
	tm.plan.Direction.Scan(pf.directionItem.Widget.(*widget.RadioGroup).Selected)
	tm.plan.Risk = decimal.RequireFromString(pf.riskItem.Widget.(*FloatEntry).Text)
	tm.plan.TakeProfitStrategy.Scan(pf.TPStratItem.Widget.(*widget.Select).Selected)
	if pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL != nil {
		tm.plan.TradingViewPlan = pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL.String()
	}
	tm.plan.Notes = pf.notesItem.Widget.(*widget.Entry).Text
}

func (pf *planForm) gatherOrders() {
	tm.orders[cryptodb.MarketStopLoss].Price = decimal.RequireFromString(pf.stopLossItem.Widget.(*FloatEntry).Text)
	tm.orders[cryptodb.Entry].Price = decimal.RequireFromString(pf.entryItem.Widget.(*FloatEntry).Text)
	for i := 0; i < cryptodb.MaxTakeProfits; i++ {
		tempPrice, err := decimal.NewFromString(pf.takeProfitItems[i].Widget.(*FloatEntry).Text)
		if err == nil {
			tm.orders[3+i].Price = tempPrice
		} else {
			tm.orders[3+i].Price = decimal.Zero
		}
	}
}

func (pf *planForm) saveSetup() {
	pf.gatherPlan()
	var err error
	tm.plan, err = savePlan(tm.plan)
	if err != nil {
		dialog.ShowError(err, ui.mainWindow)
	}

	pf.gatherOrders()
	if tm.orders[cryptodb.MarketStopLoss].PlanID == 0 {
		for i := range tm.orders {
			tm.orders[i].PlanID = tm.plan.ID
		}
	}

	_, err = saveOrders(tm.orders)
	if err != nil {
		dialog.ShowError(err, ui.mainWindow)
	}

	tm.review.PlanID = tm.plan.ID
	_, err = saveReview(tm.review)
	if err != nil {
		dialog.ShowError(err, ui.mainWindow)
	}
}

func (pf *planForm) okAction() {
	pf.saveSetup()

	tm.plans, _ = getPlans()
	ui.planList.Refresh()
	// makePlanForm()
}

func (pf *planForm) undoAction() {
	tm.plan, _ = getPlan(tm.plan.ID)
	tm.orders, _ = getOrders(tm.plan.ID)
	// makePlanForm()
	// pf.leftForm.Refresh()
	// pf.rightForm.Refresh()
}

func (pf *planForm) executeAction() {
	pf.saveSetup()
	executePlan(tm.plan.ID)

	tm.plans, _ = getPlans()
	ui.planList.Refresh()
}

func (pf *planForm) historyAction() {
	logEntries, err := getLogs(tm.plan.ID)
	if err != nil {
		dialog.ShowError(err, ui.mainWindow)
	}

	logFile := widget.NewRichText()
	for _, entry := range logEntries {
		logSegment := &widget.TextSegment{
			Style: widget.RichTextStyle{},
			Text:  fmt.Sprintf("%s - %s", entry.CreatedAt.Format("2006-01-02 15:04:05"), entry.Text),
		}
		logFile.Segments = append(logFile.Segments, logSegment)
	}

	logWindow := widget.NewPopUp(container.NewVScroll(logFile), ui.mainWindow.Canvas())
	logAnimation := canvas.NewSizeAnimation(
		fyne.NewSize(ui.mainWindow.Canvas().Size().Width-2*50.0, 0),
		fyne.NewSize(ui.mainWindow.Canvas().Size().Width-2*50.0, ui.mainWindow.Canvas().Size().Height-1*50.0),
		150*time.Millisecond,
		func(s fyne.Size) {
			logWindow.Resize(s)
		})

	logWindow.ShowAtPosition(fyne.Position{X: 50, Y: 0})
	logAnimation.Start()
}

func (pf *planForm) reviewAction() {
	reviewContainer := makeReviewForm()
	af.parentWindow = ui.app.NewWindow("Review")
	af.parentWindow.SetContent(reviewContainer)
	af.parentWindow.Show()
}

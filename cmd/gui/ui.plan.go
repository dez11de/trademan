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

type planContainer struct {
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

	toolbar *widget.Toolbar
}

var pc planContainer

func NewPlanContainer() *fyne.Container {
	pc.leftForm = widget.NewForm()
	pc.rightForm = widget.NewForm()
	bottomForm := widget.NewForm()

	pc.leftForm.AppendItem(pc.makePairItem())
	pc.leftForm.AppendItem(pc.makeDirectionItem())
	pc.leftForm.AppendItem(pc.makeRiskItem())
	pc.leftForm.AppendItem(pc.makeTradingViewItem())
	pc.leftForm.AppendItem(pc.makeTakeProfitStrategyItem())
	pc.rightForm.AppendItem(pc.makeStopLossItem())
	pc.rightForm.AppendItem(pc.makeEntryItem())

	for i := 0; i <= cryptodb.MaxTakeProfits-1; i++ {
		pc.rightForm.AppendItem(pc.makeTakeProfitItem(i, tm.pair.PriceScale, tm.pair.Price.Tick))
	}
	pc.setQuoteCurrency()
	pc.setPriceScale()

	bottomForm.AppendItem(pc.makeNotesItem())
	pc.toolbar = pc.makeToolbar()
	totalContainer := container.NewVBox(pc.makeStatContainer(),
		container.New(layout.NewGridLayoutWithColumns(2), pc.leftForm, pc.rightForm),
		bottomForm,
		pc.toolbar)

	return totalContainer
}

func (pf *planContainer) gatherPlan() {
	tm.plan.PairID = tm.pair.ID
	tm.plan.Direction.Scan(pf.directionItem.Widget.(*widget.RadioGroup).Selected)
	tm.plan.Risk = decimal.RequireFromString(pf.riskItem.Widget.(*FloatEntry).Text)
	tm.plan.TakeProfitStrategy.Scan(pf.TPStratItem.Widget.(*widget.Select).Selected)
	if pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL != nil {
		tm.plan.TradingViewPlan = pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL.String()
	}
	tm.plan.Notes = pf.notesItem.Widget.(*widget.Entry).Text
}

func (pf *planContainer) gatherOrders() {
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

func (pf *planContainer) saveSetup() {
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
}

func (pf *planContainer) okAction() {
	pf.saveSetup()

	tm.plans, _ = getPlans()
	ui.planList.Refresh()
}

func (pf *planContainer) undoAction() {
	tm.plan, _ = getPlan(tm.plan.ID)
	tm.orders, _ = getOrders(tm.plan.ID)
	NewPlanContainer()
	pf.leftForm.Refresh()
	pf.rightForm.Refresh()
}

func (pf *planContainer) executeAction() {
	pf.saveSetup()
	executePlan(tm.plan.ID)

	tm.plans, _ = getPlans()
	ui.planList.Refresh()
}

func (pf *planContainer) historyAction() {
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

func (pf *planContainer) reviewAction() {
	pf.toolbar.Hide()
	ui.planListSplit.Refresh()
	reviewContainer := makeReviewForm()
	af.parentWindow = ui.app.NewWindow("Review")
	af.parentWindow.SetContent(reviewContainer)
	af.parentWindow.Show()
	pf.toolbar.Show()
	ui.planListSplit.Refresh()
}

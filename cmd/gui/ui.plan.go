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
	// pf.leftForm.Refresh()
	pf.rightForm.AppendItem(pf.makeStopLossItem())
	pf.rightForm.AppendItem(pf.makeEntryItem())

	for i := 0; i <= cryptodb.MaxTakeProfits-1; i++ {
		pf.rightForm.AppendItem(pf.makeTakeProfitItem(i, act.pair.PriceScale, act.pair.Price.Tick))
	}
	pf.setQuoteCurrency()
	pf.setPriceScale()
	// pf.rightForm.Refresh()

	bottomForm.AppendItem(pf.makeNotesItem())

	totalContainer := container.NewVBox(pf.makeStatContainer(),
		container.New(layout.NewGridLayoutWithColumns(2), pf.leftForm, pf.rightForm),
        bottomForm,
		pf.makeToolBar())

	return totalContainer
}

func (pf *planForm) gatherPlan() {
	act.plan.PairID = act.pair.ID
	act.plan.Direction.Scan(pf.directionItem.Widget.(*widget.RadioGroup).Selected)
	act.plan.Risk = decimal.RequireFromString(pf.riskItem.Widget.(*FloatEntry).Text)
	act.plan.TakeProfitStrategy.Scan(pf.TPStratItem.Widget.(*widget.Select).Selected)
	if pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL != nil {
		act.plan.TradingViewPlan = pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).URL.String()
	}
	act.plan.Notes = pf.notesItem.Widget.(*widget.Entry).Text
}

func (pf *planForm) gatherOrders() {
	act.orders[cryptodb.MarketStopLoss].Price = decimal.RequireFromString(pf.stopLossItem.Widget.(*FloatEntry).Text)
	act.orders[cryptodb.Entry].Price = decimal.RequireFromString(pf.entryItem.Widget.(*FloatEntry).Text)
	for i := 0; i < cryptodb.MaxTakeProfits; i++ {
		tempPrice, err := decimal.NewFromString(pf.takeProfitItems[i].Widget.(*FloatEntry).Text)
		if err == nil {
			act.orders[3+i].Price = tempPrice
		} else {
			act.orders[3+i].Price = decimal.Zero
		}
	}
}

func (pf *planForm) saveSetup() {
	pf.gatherPlan()
	var err error
	act.plan, err = savePlan(act.plan)
	if err != nil {
		dialog.ShowError(err, application.mw)
	}

	pf.gatherOrders()
	if act.orders[cryptodb.MarketStopLoss].PlanID == 0 {
		for i := range act.orders {
			act.orders[i].PlanID = act.plan.ID
		}
	}

	_, err = saveOrders(act.orders)
	if err != nil {
		dialog.ShowError(err, application.mw)
	}

	act.review.PlanID = act.plan.ID
	_, err = saveReview(act.review)
	if err != nil {
		dialog.ShowError(err, application.mw)
	}
}

func (pf *planForm) okAction() {
	pf.saveSetup()

	tm.plans, _ = getPlans()
	application.planList.Refresh()
	// makePlanForm()
}

func (pf *planForm) undoAction() {
	act.plan, _ = getPlan(act.plan.ID)
	act.orders, _ = getOrders(act.plan.ID)
	// makePlanForm()
	// pf.leftForm.Refresh()
	// pf.rightForm.Refresh()
}

func (pf *planForm) executeAction() {
	pf.saveSetup()
	executePlan(act.plan.ID)

	tm.plans, _ = getPlans()
	application.planList.Refresh()
}

func (pf *planForm) historyAction() {
	logEntries, err := getLogs(act.plan.ID)
	if err != nil {
		dialog.ShowError(err, application.mw)
	}

	logFile := widget.NewRichText()
	for _, entry := range logEntries {
		logSegment := &widget.TextSegment{
			Style: widget.RichTextStyle{},
			Text:  fmt.Sprintf("%s - %s", entry.CreatedAt.Format("2006-01-02 15:04:05"), entry.Text),
		}
		logFile.Segments = append(logFile.Segments, logSegment)
	}

	logWindow := widget.NewPopUp(container.NewVScroll(logFile), application.mw.Canvas())
	logAnimation := canvas.NewSizeAnimation(
		fyne.NewSize(application.mw.Canvas().Size().Width-2*50.0, 0),
		fyne.NewSize(application.mw.Canvas().Size().Width-2*50.0, application.mw.Canvas().Size().Height-1*50.0),
		150*time.Millisecond,
		func(s fyne.Size) {
			logWindow.Resize(s)
		})

	logWindow.ShowAtPosition(fyne.Position{X: 50, Y: 0})
	logAnimation.Start()
}

func (pf *planForm) reviewAction() {
	reviewContainer := makeReviewForm()
	af.parentWindow = application.fa.NewWindow("Review")
	af.parentWindow.SetContent(reviewContainer)
	af.parentWindow.Show()
}

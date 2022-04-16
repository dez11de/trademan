package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
)

type planForm struct {
	form          *widget.Form
	formContainer *fyne.Container

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

func makePlanForm() *fyne.Container {
	pf.form = widget.NewForm()

	pf.form.AppendItem(pf.makePairItem())
	pf.form.AppendItem(pf.makeDirectionItem())
	pf.form.AppendItem(pf.makeRiskItem())
	pf.form.AppendItem(pf.makeStopLossItem())
	pf.form.AppendItem(pf.makeEntryItem())
	pf.form.AppendItem(pf.makeTakeProfitStrategyItem())

	for i := 0; i <= cryptodb.MaxTakeProfits-1; i++ {
		pf.form.AppendItem(pf.makeTakeProfitItem(i, act.pair.PriceScale, act.pair.Price.Tick))
	}

	pf.form.AppendItem(pf.makeTradingViewItem())
	pf.form.AppendItem(pf.makeNotesItem())

	pf.setQuoteCurrency()
	pf.setPriceScale()

	toolBar := pf.makeToolBar()

	act.statisticsContainer = pf.makeStatContainer()

	pf.form.Refresh()
	pf.formContainer = container.NewBorder(act.statisticsContainer, toolBar, nil, nil, pf.form)

	return pf.formContainer
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
		dialog.ShowError(err, mainWindow)
	}

	pf.gatherOrders()
	if act.orders[cryptodb.MarketStopLoss].PlanID == 0 {
		for i := range act.orders {
			act.orders[i].PlanID = act.plan.ID
		}
	}

	_, err = saveOrders(act.orders)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}

	act.assessment.PlanID = act.plan.ID
	_, err = saveAssessment(act.assessment)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}

	// pf.FillForm(updatedPlan)
}

func (pf *planForm) okAction() {
	pf.saveSetup()

	tm.plans, _ = getPlans()
	act.List.Refresh()
}

func (pf *planForm) undoAction() {
	// reloadedPlan, _ := getPlan(act.plan.ID)
	// pf.FillForm(reloadedPlan)
}

func (pf *planForm) executeAction() {
	pf.saveSetup()
	executePlan(act.plan.ID)

	tm.plans, _ = getPlans()
	act.List.Refresh()
}

func (pf *planForm) historyAction() {
	entries, err := getLogs(act.plan.ID)
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
	container := makeAssessmentForm()
	af.window = a.NewWindow("Assessment")
	af.window.SetContent(container)
	af.window.Show()
}

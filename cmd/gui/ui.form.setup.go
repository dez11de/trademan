package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
)

type planForm struct {
	formContainer           *fyne.Container
	form                    *widget.Form
	pairItem                *widget.FormItem
	directionItem           *widget.FormItem
	riskItem                *widget.FormItem
	TPStratItem             *widget.FormItem
	stopLossItem            *widget.FormItem
	entryItem               *widget.FormItem
	takeProfitItems         [cryptodb.MaxTakeProfits]*widget.FormItem
	notesMultilineEntryItem *widget.FormItem
	tradingViewPlanItem     *widget.FormItem
}

func NewForm() *planForm {
	pf := new(planForm)

	pf.form = widget.NewForm()

	pf.pairItem = pf.makePairItem()
	pf.form.AppendItem(pf.pairItem)
	pf.directionItem = pf.makeDirectionItem()
	pf.form.AppendItem(pf.directionItem)
	pf.riskItem = pf.makeRiskItem()
	pf.form.AppendItem(pf.riskItem)
	pf.stopLossItem = pf.makeStopLossItem(0, decimal.Zero)
	pf.form.AppendItem(pf.stopLossItem)
	pf.entryItem = pf.makeEntryItem(0, decimal.Zero)
	pf.form.AppendItem(pf.entryItem)
	pf.TPStratItem = pf.makeTakeProfitStrategyItem()
	pf.form.AppendItem(pf.TPStratItem)

	for i := 0; i <= cryptodb.MaxTakeProfits-1; i++ {
		pf.takeProfitItems[i] = pf.makeTakeProfitItem(i, 0, decimal.Zero)
		pf.form.AppendItem(pf.takeProfitItems[i])
	}

	pf.tradingViewPlanItem = pf.makeTradingViewLinkItem()
	pf.form.AppendItem(pf.tradingViewPlanItem)
	pf.notesMultilineEntryItem = pf.makeNotesMultilineItem()
	pf.form.AppendItem(pf.notesMultilineEntryItem)

	pf.setQuoteCurrency(" . . . ")
	pf.setPriceScale(0, decimal.Zero)

	toolBar := pf.makeToolBar()

	ui.statisticsContainer = pf.makeStatContainer()

	pf.formContainer = container.NewBorder(ui.statisticsContainer, toolBar, nil, nil, pf.form)

	pf.form.Refresh()
	return pf
}

func (pf *planForm) FillForm(p cryptodb.Plan) {
	var err error
	ui.activePlan = p

	if ui.activePlan.PairID == 0 {
		ui.activeOrders = cryptodb.NewOrders(0)
		ui.activeAssessment= cryptodb.NewAssessment(0)
	} else {
		ui.activePair = ui.Pairs[ui.activePlan.PairID-1]
		ui.activeOrders, err = getOrders(ui.activePlan.ID)
		if err != nil {
			dialog.ShowError(err, mainWindow)
		}
        ui.activeAssessment, err = getAssessment(ui.activePlan.ID)
		if err != nil {
			dialog.ShowError(err, mainWindow)
		}
	}

	if ui.activePlan.ID == 0 {
		pf.pairItem.Widget.(*xwidget.CompletionEntry).Enable()
	} else {
		pf.pairItem.Widget.(*xwidget.CompletionEntry).SetText(ui.activePair.Name)
		pf.pairItem.Widget.(*xwidget.CompletionEntry).Disable()
		pf.directionItem.Widget.(*widget.RadioGroup).SetSelected(ui.activePlan.Direction.String())
		pf.directionItem.Widget.(*widget.RadioGroup).Disable()

		pf.riskItem.Widget.(*FloatEntry).SetText(ui.activePlan.Risk.StringFixed(1))
		pf.stopLossItem.Widget.(*FloatEntry).SetText(ui.activeOrders[cryptodb.MarketStopLoss].Price.StringFixed(ui.activePair.PriceScale))
		pf.entryItem.Widget.(*FloatEntry).SetText(ui.activeOrders[cryptodb.Entry].Price.StringFixed(ui.activePair.PriceScale))
		pf.TPStratItem.Widget.(*widget.Select).SetSelected(ui.activePlan.TakeProfitStrategy.String())
		pf.TPStratItem.Widget.(*widget.Select).Disable()

		if ui.activeOrders[cryptodb.Entry].Status >= cryptodb.New {
			pf.riskItem.Widget.(*FloatEntry).Disable()
			pf.stopLossItem.Widget.(*FloatEntry).Disable()
			pf.entryItem.Widget.(*FloatEntry).Disable()
			pf.TPStratItem.Widget.(*widget.Select).Disable()
		}

		takeProfitCount := 0
		for _, o := range ui.activeOrders {
			if o.OrderKind == cryptodb.TakeProfit {
				pf.takeProfitItems[takeProfitCount].Widget.(*FloatEntry).SetText(o.Price.StringFixed(ui.activePair.PriceScale))
				if o.Status <= cryptodb.Ordered {
					pf.takeProfitItems[takeProfitCount].Widget.(*FloatEntry).Enable()
				}
				takeProfitCount++
			}
		}

		pf.setQuoteCurrency(ui.activePair.QuoteCurrency)
		pf.setPriceScale(ui.activePair.PriceScale, ui.activePair.Price.Tick)

		if ui.activePlan.TradingViewPlan == "" {
			pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Entry).Show()
			pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Hyperlink).Hide()
			pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button).Show()
			pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button).Hide()
		}
		pf.notesMultilineEntryItem.Widget.(*widget.Entry).SetText(ui.activePlan.Notes)
	}
	pf.form.Refresh()
}

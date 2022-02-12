package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
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
	stopLossItem            *widget.FormItem
	entryItem               *widget.FormItem
	takeProfitItems         [cryptodb.MaxTakeProfits]*widget.FormItem
	notesMultilineEntryItem *widget.FormItem
	tradingViewPlanItem     *widget.FormItem
}

func NewForm() *planForm {
	pf := new(planForm)
	pf.form = widget.NewForm()

	// ui.activeOrders = cryptodb.NewOrders(0)

	pf.pairItem = pf.makePairItem()
	pf.form.AppendItem(pf.pairItem)
	pf.directionItem = pf.makeDirectionItem()
	pf.form.AppendItem(pf.directionItem)
	pf.riskItem = pf.makeRiskItem()
	pf.form.AppendItem(pf.riskItem)
	pf.stopLossItem = pf.makeStopLossItem()
	pf.form.AppendItem(pf.stopLossItem)
	pf.entryItem = pf.makeEntryItem()
	pf.form.AppendItem(pf.entryItem)

	for i := 0; i < cryptodb.MaxTakeProfits; i++ {
		pf.takeProfitItems[i] = pf.makeTakeProfitItem(i)
		pf.form.AppendItem(pf.takeProfitItems[i])
	}

	pf.tradingViewPlanItem = pf.makeTradingViewLinkItem()
	pf.form.AppendItem(pf.tradingViewPlanItem)
	pf.notesMultilineEntryItem = pf.makeNotesMultilineItem()
	pf.form.AppendItem(pf.notesMultilineEntryItem)

	pf.setQuoteCurrency(" . . . ")
	pf.setPriceScale(0)

	toolBar := pf.makeToolBar()

	ui.statisticsContainer = pf.makeStatContainer()

	pf.formContainer = container.New(layout.NewBorderLayout(ui.statisticsContainer, toolBar, nil, nil), ui.statisticsContainer, toolBar, pf.form)

	pf.form.Refresh()
	return pf
}

func (pf *planForm) FillForm(p cryptodb.Plan) {
	var err error
	ui.activePlan = p
	if ui.activePlan.PairID != 0 {
		ui.activePair = ui.Pairs[ui.activePlan.PairID-1]
	}
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}
	ui.activeOrders, err = getOrders(ui.activePlan.ID)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}

	if ui.activePlan.ID != 0 {
		pf.pairItem.Widget.(*xwidget.CompletionEntry).Disable()
		pf.pairItem.Widget.(*xwidget.CompletionEntry).SetText(ui.activePair.Name)
		pf.setQuoteCurrency(ui.activePair.QuoteCurrency)
		pf.setPriceScale(ui.activePair.PriceScale)
	}

	if ui.activePlan.ID != 0 {
		switch ui.activePlan.Direction {
		case cryptodb.Direction(cryptodb.DirectionLong):
			pf.directionItem.Widget.(*widget.RadioGroup).SetSelected("Long")
		case cryptodb.Direction(cryptodb.DirectionShort):
			pf.directionItem.Widget.(*widget.RadioGroup).SetSelected("Short")
		}
		pf.directionItem.Widget.(*widget.RadioGroup).Disable()
	}

	// TODO: think about in which statusses changing is allowed
	if ui.activePlan.Risk.Cmp(decimal.Zero) != 0 {
		pf.riskItem.Widget.(*widget.Entry).SetText(ui.activePlan.Risk.StringFixed(1))
		if ui.activePlan.Status != cryptodb.StatusPlanned {
			pf.riskItem.Widget.(*widget.Entry).Disable()
		}
	}

	// TODO: think about in which statusses changing is allowed
	if !ui.activeOrders[cryptodb.KindMarketStopLoss].Price.Equal(decimal.Zero) {
		pf.stopLossItem.Widget.(*widget.Entry).SetText(ui.activeOrders[cryptodb.KindMarketStopLoss].Price.StringFixed(ui.activePair.PriceScale))
	}

	// TODO: think about in which statusses changing is allowed, disable editting if required
	if ui.activeOrders[cryptodb.KindEntry].Price.Cmp(decimal.Zero) != 0 {
		pf.entryItem.Widget.(*widget.Entry).SetText(ui.activeOrders[cryptodb.KindEntry].Price.StringFixed(ui.activePair.PriceScale))
	}

	// TODO: think about in which statusses changing is allowed, disable editting if required
	takeProfitCount := 0
	for _, o := range ui.activeOrders {
		if o.OrderKind == cryptodb.KindTakeProfit && o.Price.Cmp(decimal.Zero) != 0 {
			pf.takeProfitItems[takeProfitCount].Widget.(*widget.Entry).SetText(o.Price.StringFixed(ui.activePair.PriceScale))
			takeProfitCount++
		}
	}

	// TODO: think about in which statusses changing is allowed
	if p.Notes != "" {
		pf.notesMultilineEntryItem.Widget.(*widget.Entry).SetText(ui.activePlan.Notes)
	}

	/*
		if pf.plan.TradingViewPlan == "" {
			pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[1].Show()
		} else {
			pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[0].Show()
		}
	*/
}

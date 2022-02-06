package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"github.com/dez11de/cryptodb"
	"github.com/shopspring/decimal"
)

type planForm struct {
	PairCache  []cryptodb.Pair
	activePair cryptodb.Pair

	plan   cryptodb.Plan
	orders []cryptodb.Order

	formContainer *fyne.Container
	statisticsContainer *fyne.Container

	form                    *widget.Form
	pairItem                *widget.FormItem
	sideItem                *widget.FormItem
	riskItem                *widget.FormItem
	stopLossItem            *widget.FormItem
	entryItem               *widget.FormItem
	takeProfitItems         [5]*widget.FormItem // TODO: restore MaxTakeProfits to it's former glory
	notesMultilineEntryItem *widget.FormItem
	tradingViewPlanItem     *widget.FormItem
}

func NewForm() *planForm {
	pf := &planForm{}
	var err error
	pf.PairCache, err = getPairs()
	if err != nil {
		return pf
	}
	pf.orders = cryptodb.NewOrders(0)

	pf.form = widget.NewForm()

	pf.pairItem = pf.makePairItem()
	pf.form.AppendItem(pf.pairItem)
	pf.sideItem = pf.makeSideItem()
	pf.form.AppendItem(pf.sideItem)
	pf.riskItem = pf.makeRiskItem()
	pf.form.AppendItem(pf.riskItem)
	pf.stopLossItem = pf.makeStopLossItem()
	pf.form.AppendItem(pf.stopLossItem)
	pf.entryItem = pf.makeEntryItem()
	pf.form.AppendItem(pf.entryItem)

	takeProfitCount := 0
	for _, order := range pf.orders {
		if order.OrderType == cryptodb.TypeTakeProfit {
			pf.takeProfitItems[takeProfitCount] = pf.makeTakeProfitItem(takeProfitCount)
			pf.form.AppendItem(pf.takeProfitItems[takeProfitCount])
			takeProfitCount++
		}
	}

	pf.tradingViewPlanItem = pf.makeTradingViewLinkItem()
	pf.form.AppendItem(pf.tradingViewPlanItem)
	pf.notesMultilineEntryItem = pf.makeNotesMultilineItem()
	pf.form.AppendItem(pf.notesMultilineEntryItem)

	pf.setQuoteCurrency(" . . . ")
	pf.setPriceScale(1)

	pf.form.OnSubmit = pf.formSubmit
	pf.form.OnCancel = pf.formCancel
	pf.form.SubmitText = "OK"
	pf.form.CancelText = "Cancel"
	pf.form.Disable()

	pf.statisticsContainer = pf.makeStatContainer()

	pf.formContainer = container.New(layout.NewBorderLayout(pf.statisticsContainer, nil, nil, nil), pf.statisticsContainer, pf.form)
	return pf
}

func (pf *planForm) FillForm(p cryptodb.Plan) {
	var err error
	pf.plan = p
	pf.activePair, err = getPair(pf.plan.PairID)
	if err != nil {
	}
	pf.orders, err = getOrders(pf.plan.ID)
	if err != nil {
	}

	if pf.plan.ID != 0 {
		pf.pairItem.Widget.(*xwidget.CompletionEntry).Disable()
		pf.pairItem.Widget.(*xwidget.CompletionEntry).SetText(pf.activePair.Name)
		pf.setQuoteCurrency(pf.activePair.Name)
		pf.setPriceScale(int32(pf.activePair.PriceScale))
	}

	if pf.plan.ID != 0 {
		switch pf.plan.Side {
		case cryptodb.SideLong:
			pf.sideItem.Widget.(*widget.RadioGroup).SetSelected("Long")
		case cryptodb.SideShort:
			pf.sideItem.Widget.(*widget.RadioGroup).SetSelected("Short")
		}
		pf.sideItem.Widget.(*widget.RadioGroup).Disable()
	}

	// TODO: think about in which statusses changing is allowed
	if pf.plan.Risk.Cmp(decimal.Zero) != 0 {
		pf.riskItem.Widget.(*widget.Entry).SetText(pf.plan.Risk.StringFixed(2))
		if pf.plan.Status != cryptodb.StatusPlanned {
			pf.riskItem.Widget.(*widget.Entry).Disable()
		}
	}

	// TODO: think about in which statusses changing is allowed
	if !pf.orders[cryptodb.TypeHardStopLoss].Price.Equal(decimal.Zero) {
		pf.stopLossItem.Widget.(*widget.Entry).SetText(pf.orders[cryptodb.TypeHardStopLoss].Price.StringFixed(int32(pf.activePair.PriceScale)))
	}

	// TODO: think about in which statusses changing is allowed, disable editting if required
	if pf.orders[cryptodb.TypeEntry].Price.Cmp(decimal.Zero) != 0 {
		pf.entryItem.Widget.(*widget.Entry).SetText(pf.orders[cryptodb.TypeEntry].Price.StringFixed(int32(pf.activePair.PriceScale)))
	}

	// TODO: think about in which statusses changing is allowed, disable editting if required
	takeProfitCount := 0
	for _, o := range pf.orders {
		if o.OrderType == cryptodb.TypeTakeProfit && o.Price.Cmp(decimal.Zero) != 0 {
			pf.takeProfitItems[takeProfitCount].Widget.(*widget.Entry).SetText(o.Price.StringFixed(int32(pf.activePair.PriceScale)))
			takeProfitCount++
		}
	}

	// TODO: think about in which statusses changing is allowed
	if p.Notes != "" {
		pf.notesMultilineEntryItem.Widget.(*widget.Entry).SetText(pf.plan.Notes)
	}

	//pf.tradingViewPlanItem.Widget.(*widget.Hyperlink).SetURLFromString(pf.plan.TradingViewPlan)
	pf.tradingViewPlanItem.Widget.(*fyne.Container).Objects[1].Show()
}

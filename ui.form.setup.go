package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"github.com/shopspring/decimal"
)

type planForm struct {
	db *Database

	plan       Plan
	activePair Pair
	orders     Orders

	form *widget.Form

	pairItem                *widget.FormItem
	sideItem                *widget.FormItem
	riskItem                *widget.FormItem
	stopLossItem            *widget.FormItem
	entryItem               *widget.FormItem
	takeProfitItems         [MaxTakeProfits]*widget.FormItem
	notesMultilineEntryItem *widget.FormItem
	tradingViewPlanItem     *widget.FormItem
}

func (pf *planForm) makePairItem() *widget.FormItem {
	CompletionEntry := xwidget.NewCompletionEntry([]string{})

	CompletionEntry.OnChanged = func(s string) {
		s = strings.ToUpper(s)
		possiblePairs, _ := pf.db.SearchPairs(s)
		if len(possiblePairs) > 1 {
			CompletionEntry.SetOptions(possiblePairs)
			CompletionEntry.ShowCompletion()
		} else {
			CompletionEntry.SetText(possiblePairs[0])
			CompletionEntry.HideCompletion()
		}
	}

	CompletionEntry.OnSubmitted = func(s string) {
		var ok bool
		pf.activePair, ok = pf.db.PairCache[s]
		if ok {
			pf.stopLossItem.Text = fmt.Sprintf("Stop Loss (%s)", pf.activePair.QuoteCurrency)
			pf.stopLossItem.Widget.(*widget.Entry).SetPlaceHolder(decimal.Zero.StringFixed(pf.activePair.PriceScale))
			pf.entryItem.Text = fmt.Sprintf("Entry (%s)", pf.activePair.QuoteCurrency)
			pf.entryItem.Widget.(*widget.Entry).SetPlaceHolder(decimal.Zero.StringFixed(pf.activePair.PriceScale))
			for i, takeProfitItem := range pf.takeProfitItems {
				takeProfitItem.Text = fmt.Sprintf("Take Profit #%d (%s)", i+1, pf.activePair.QuoteCurrency)
				takeProfitItem.Widget.(*widget.Entry).SetPlaceHolder(decimal.Zero.StringFixed(pf.activePair.PriceScale))
			}
			pf.sideItem.Widget.(*widget.RadioGroup).Enable()
			pf.form.Refresh()
		}
	}
	return widget.NewFormItem("Pair", CompletionEntry)
}

func (pf *planForm) makeSideItem() *widget.FormItem {
	sideRadio := widget.NewRadioGroup([]string{"Long", "Short"},
		func(s string) {
			pf.riskItem.Widget.(*widget.Entry).Enable()
			pf.form.Refresh()
		})
	sideRadio.Horizontal = true
	sideRadio.Disable()
	return widget.NewFormItem("Side", sideRadio)
}

func (pf *planForm) makeRiskItem() *widget.FormItem {
	// TODO: validate input
	riskEntry := widget.NewEntry()
	riskEntry.Disable()
	riskEntry.SetPlaceHolder(pf.plan.Risk.StringFixed(2))
	riskEntry.OnChanged = func(s string) {
		// TODO: catch error
		tempRisk := decimal.RequireFromString(s)
		if (tempRisk.Cmp(decimal.NewFromFloat(5)) == -1) && (tempRisk.Cmp(decimal.NewFromFloat(0.499)) == 1) {
			pf.stopLossItem.Widget.(*widget.Entry).Enable()
			pf.form.Refresh()
		}
	}

	return widget.NewFormItem("Risk (%)", riskEntry)
}

func (pf *planForm) makeStopLossItem() *widget.FormItem {
	StopLossEntry := widget.NewEntry()
	StopLossEntry.Disable()
	StopLossEntry.OnChanged = func(s string) {
		// TODO: properly validate input
		pf.entryItem.Widget.(*widget.Entry).Enable()
		pf.form.Refresh()
	}
	return widget.NewFormItem(fmt.Sprintf("Stop Loss (%s)", pf.activePair.QuoteCurrency), StopLossEntry)
}

func (pf *planForm) makeEntryItem() *widget.FormItem {
	entryEntry := widget.NewEntry()
	entryEntry.Disable()
	entryEntry.OnChanged = func(s string) {
		// TODO: properly validate input
		pf.takeProfitItems[0].Widget.(*widget.Entry).Enable()
		pf.form.Refresh()
	}
	//entryEntry.SetPlaceHolder(decimal.Zero.StringFixed(pf.activePair.PriceScale))
	return widget.NewFormItem(fmt.Sprintf("Entry (%s)", pf.activePair.QuoteCurrency), entryEntry)
}

func (pf *planForm) makeTakeProfitItem(n int) *widget.FormItem {
	takeProfitEntry := widget.NewEntry()
	takeProfitEntry.Disable()
	takeProfitEntry.OnChanged = func(s string) {
		// TODO: properly validate input
		// TODO: show % difference with entry and previous take profit
		if !decimal.RequireFromString(s).IsZero() && s != "" {
			pf.takeProfitItems[n+1].Widget.(*widget.Entry).Enable()
			pf.form.Refresh()
		}
	}
	return widget.NewFormItem(fmt.Sprintf("Take profit #%d (%s)", n+1, pf.activePair.QuoteCurrency), takeProfitEntry)
}

func (pf *planForm) makeTradingViewLinkItem() *widget.FormItem {
	var tvurl url.URL
	log.Printf("TV Url: %v", pf.plan.TradingViewPlan)
	tvurl.Parse(pf.plan.TradingViewPlan)
	log.Printf("TV Url: %v", tvurl)
	tradingViewLink := widget.NewHyperlink("Open", &tvurl)

	return widget.NewFormItem("TradingView", tradingViewLink)
}

func (pf *planForm) makeNotesMultilineItem() *widget.FormItem {
	notesMultiLineEntry := widget.NewMultiLineEntry()
	notesMultiLineEntry.SetPlaceHolder("Enter notes...")
	notesMultiLineEntry.SetText(pf.plan.Notes)
	notesMultiLineEntry.Wrapping = fyne.TextWrapWord

	return widget.NewFormItem("Notes", notesMultiLineEntry)
}

func NewForm(d *Database, bb *ByBit) *planForm {
	pf := &planForm{db: d}
	pf.orders = NewOrders()

	pf.form = widget.NewForm()
	pf.form.SubmitText = "OK"
	pf.form.CancelText = "Cancel"

	pf.form.OnSubmit = func() {
		log.Printf("Submit button pressed.")
		pf.plan.PairID = pf.activePair.PairID
		// pf.plan.Side.Scan = pf.sideItem.Widget.(*widget.RadioGroup).Selected
		pf.plan.Risk = decimal.RequireFromString(pf.riskItem.Widget.(*widget.Entry).Text)
		pf.orders[typeHardStopLoss].Price = decimal.RequireFromString(pf.stopLossItem.Widget.(*widget.Entry).Text)
		pf.orders[typeEntry].Price = decimal.RequireFromString(pf.entryItem.Widget.(*widget.Entry).Text)
		for i := 0; i < MaxTakeProfits-1; i++ {
			tempPrice, err := decimal.NewFromString(pf.takeProfitItems[i].Widget.(*widget.Entry).Text)
			if err == nil {
				pf.orders[3+i].Price = tempPrice
			} else {
				pf.orders[3+i].Price = decimal.Zero
			}
		}

		/*
			pf.plan.SetEntrySize(pf.activePair, d.WalletCache[pf.activePair.QuoteCurrency].Equity, &pf.orders)
			pf.plan.SetTakeProfitSizes(pf.activePair, &pf.orders)
			pf.plan.SetRewardRiskRatio(pf.orders)
			log.Printf("[Form] Storing to database...")
			pf.db.StorePlanAndOrders(pf.plan, pf.orders)
		*/
		log.Printf("Sending to exchange...")
		bb.PlaceOrders(pf.plan, pf.activePair, pf.orders)

	}

	pf.form.OnCancel = func() {
		log.Printf("Cancel button pressed.")
	}
	//pf.form.Disable()

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
		if order.OrderType == typeTakeProfit {
			pf.takeProfitItems[takeProfitCount] = pf.makeTakeProfitItem(takeProfitCount)
			pf.form.AppendItem(pf.takeProfitItems[takeProfitCount])
			takeProfitCount++
		}
	}

	pf.tradingViewPlanItem = pf.makeTradingViewLinkItem()
	pf.form.AppendItem(pf.tradingViewPlanItem)
	pf.notesMultilineEntryItem = pf.makeNotesMultilineItem()
	pf.form.AppendItem(pf.notesMultilineEntryItem)

	return pf
}

func (pf *planForm) FillForm(p Plan) {
	pf.plan = p
	pf.orders, _ = pf.db.GetOrders(pf.plan.PlanID)
	for i, o := range pf.orders {
		log.Printf("[%d] %v", i, o)
	}
	pf.activePair, _ = pf.db.GetPairFromID(pf.plan.PairID)

	if pf.plan.PlanID != 0 {
		pf.pairItem.Widget.(*xwidget.CompletionEntry).SetText(pf.activePair.Pair)
		pf.pairItem.Widget.(*xwidget.CompletionEntry).Disable()
	}

	// TODO: consider adding a .Help line

	// TODO: make better use of enumer options, see TODO.txt
	if pf.plan.PlanID != 0 {
		switch pf.plan.Side {
		case sideLong:
			pf.sideItem.Widget.(*widget.RadioGroup).SetSelected("Long")
		case sideShort:
			pf.sideItem.Widget.(*widget.RadioGroup).SetSelected("Short")
		}
		pf.sideItem.Widget.(*widget.RadioGroup).Disable()
	}

	// TODO: think about in which statusses changing is allowed
	if pf.plan.Risk.Cmp(decimal.Zero) != 0 {
		pf.riskItem.Widget.(*widget.Entry).SetText(pf.plan.Risk.StringFixed(2))
	}

	// TODO: think about in which statusses changing is allowed
	if pf.orders[typeHardStopLoss].Price.Cmp(decimal.Zero) != 0 {
		pf.stopLossItem.Widget.(*widget.Entry).SetText(pf.orders[typeHardStopLoss].Price.StringFixed(pf.activePair.PriceScale))
	}

	// TODO: think about in which statusses changing is allowed, disable editting if required
	if pf.orders[typeEntry].Price.Cmp(decimal.Zero) != 0 {
		pf.entryItem.Widget.(*widget.Entry).SetText(pf.orders[typeEntry].Price.StringFixed(pf.activePair.PriceScale))
	}

	// TODO: think about in which statusses changing is allowed, disable editting if required
	takeProfitCount := 0
	for _, o := range pf.orders {
		if o.OrderType == typeTakeProfit && o.Price.Cmp(decimal.Zero) != 0 {
			pf.takeProfitItems[takeProfitCount].Widget.(*widget.Entry).SetText(o.Price.StringFixed(pf.activePair.PriceScale))
			takeProfitCount++
		}
	}

	// TODO: think about in which statusses changing is allowed
	if p.Notes != "" {
		pf.notesMultilineEntryItem.Widget.(*widget.Entry).SetText(pf.plan.Notes)
	}

	// TODO: think about in which statusses changing is allowed
	// TODO: Set to widget.Entry if non given and/or can be modified.
	// TODO: Show edit button if can be modified
	pf.tradingViewPlanItem.Widget.(*widget.Hyperlink).SetURLFromString(pf.plan.TradingViewPlan)
}

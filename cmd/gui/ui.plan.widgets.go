package main

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func (pf *planForm) setQuoteCurrency() {
	pf.stopLossItem.Text = fmt.Sprintf("Stop Loss (%s)", act.pair.QuoteCurrency)

	pf.entryItem.Text = fmt.Sprintf("Entry (%s)", act.pair.QuoteCurrency)

	for i, takeProfitItem := range pf.takeProfitItems {
		takeProfitItem.Text = fmt.Sprintf("TP #%d (%s)", i+1, act.pair.QuoteCurrency)
	}
	pf.rightForm.Refresh()
}

func (pf *planForm) setPriceScale() {
	pf.stopLossItem.Widget.(*FloatEntry).decimals = act.pair.PriceScale
	pf.stopLossItem.Widget.(*FloatEntry).tick = act.pair.Price.Tick
	pf.stopLossItem.Widget.(*FloatEntry).SetPlaceHolder(decimal.Zero.StringFixed(act.pair.PriceScale))

	pf.entryItem.Widget.(*FloatEntry).decimals = act.pair.PriceScale
	pf.entryItem.Widget.(*FloatEntry).tick = act.pair.Price.Tick
	pf.entryItem.Widget.(*FloatEntry).SetPlaceHolder(decimal.Zero.StringFixed(act.pair.PriceScale))

	for _, takeProfitItem := range pf.takeProfitItems {
		takeProfitItem.Widget.(*FloatEntry).decimals = act.pair.PriceScale
		takeProfitItem.Widget.(*FloatEntry).tick = act.pair.Price.Tick
		takeProfitItem.Widget.(*FloatEntry).SetPlaceHolder(decimal.Zero.StringFixed(act.pair.PriceScale))
	}
}

type FloatEntry struct {
	widget.Entry
	decimals int32
	tick     decimal.Decimal
}

func NewFloatEntry(decimals int32, t decimal.Decimal) *FloatEntry {
	priceEntry := &FloatEntry{decimals: decimals, tick: t}
	priceEntry.ExtendBaseWidget(priceEntry)

	return priceEntry
}

func (fe *FloatEntry) TypedRune(r rune) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		fe.Entry.TypedRune(r)
	}
}

func (fe *FloatEntry) FocusLost() {
	v, err := decimal.NewFromString(fe.Entry.Text)
	if err == nil {
		fe.Entry.SetText(v.RoundStep(fe.tick, false).StringFixed(fe.decimals))
	}
	fe.Entry.FocusLost()
}

func (ui *active) fzfPairs(s string) (possiblePairs []string) {
	var pairNames []string
	for _, p := range tm.pairs {
		pairNames = append(pairNames, p.Name)
	}

	matches := fuzzy.RankFind(s, pairNames)
	sort.Sort(matches)

	for i, p := range matches {
		if i > 10 {
			break
		}
		possiblePairs = append(possiblePairs, p.Target)
	}

	return possiblePairs
}

func (pf *planForm) makeStatContainer() *fyne.Container {
	// TODO: make distinction between start RRR and evolved RRR
	startRewardRiskRatioLabel := widget.NewLabel("Start RRR: ")
	startRewardRiskRatioValue := widget.NewLabel(fmt.Sprintf("%.1f", 0.0))
	evolvedRewardRiskRatioLabel := widget.NewLabel("Current RRR: ")
	evolvedRewardRiskRatioValue := widget.NewLabel(fmt.Sprintf("%.1f", 0.0))

	currentPnLLabel := widget.NewLabel("PnL: ")
	currentPnLValue := widget.NewLabel(fmt.Sprintf("%s%%", act.plan.Profit.StringFixed(1))) // TODO: should be relative to entrySize.
	// TODO: figure out what this even means, see CryptoCred.
	breakEvenLabel := widget.NewLabel("B/E: ")
	breakEvenValue := widget.NewLabel(fmt.Sprintf("%.0f%%", 0.0))
	container := container.NewHBox(
		layout.NewSpacer(),
		startRewardRiskRatioLabel, startRewardRiskRatioValue,
		evolvedRewardRiskRatioLabel, evolvedRewardRiskRatioValue,
		currentPnLLabel, currentPnLValue,
		breakEvenLabel, breakEvenValue,
		layout.NewSpacer())

	return container
}

func (pf *planForm) makePairItem() *widget.FormItem {
	pairEntry := xwidget.NewCompletionEntry([]string{})
	pairEntry.SetPlaceHolder("Select pair from list")
	if act.plan.PairID != 0 {
		pairEntry.SetText(tm.pairs[act.plan.PairID-1].Name)
		pairEntry.Disable()
	}

	pf.pairItem = widget.NewFormItem("Pair", pairEntry)
	pf.pairItem.HintText = " "

	pairEntry.OnChanged = func(s string) {
		pairEntry.SetText(strings.ToUpper(s))
		if len(s) >= 2 {
			possiblePairs := act.fzfPairs(strings.ToUpper(s))
			if len(possiblePairs) == 1 {
				pairEntry.SetText(possiblePairs[0])
				pairEntry.HideCompletion()
			} else {
				pairEntry.SetOptions(possiblePairs)
				pairEntry.ShowCompletion()
			}
		}
	}

	pairEntry.OnSubmitted = func(s string) {
		s = strings.ToUpper(s)
		for _, p := range tm.pairs {
			if s == p.Name {
				act.pair = p
				act.plan.PairID = act.pair.ID
				pf.setQuoteCurrency()
				pf.setPriceScale()

				pf.directionItem.Widget.(*widget.RadioGroup).Enable()
				pf.leftForm.Refresh()
				break
			}
		}
	}

	return pf.pairItem
}

func (pf *planForm) makeDirectionItem() *widget.FormItem {
	directionRadio := widget.NewRadioGroup(nil, nil)
	directionRadio.Disable()
	directionRadio.Horizontal = true
	directionRadio.Options = cryptodb.DirectionStrings()
	pf.directionItem = widget.NewFormItem("Direction", directionRadio)
	pf.directionItem.HintText = " "

	if act.plan.ID != 0 {
		directionRadio.SetSelected(act.plan.Direction.String())
	}

	directionRadio.OnChanged =
		func(s string) {
			pf.riskItem.Widget.(*FloatEntry).Enable()
			pf.leftForm.Refresh()
		}

	return pf.directionItem
}

func (pf *planForm) makeRiskItem() *widget.FormItem {
	riskEntry := NewFloatEntry(1, decimal.NewFromFloat(0.1))
	riskEntry.Disable()
	riskEntry.SetPlaceHolder("0.0")
	pf.riskItem = widget.NewFormItem("Risk (%)", riskEntry)
	pf.riskItem.HintText = " "

	if act.plan.ID != 0 {
		riskEntry.SetText(act.plan.Risk.StringFixed(1))
		if act.plan.Status <= cryptodb.PartiallyFilled {
			riskEntry.Enable()
		}
	}
	riskEntry.OnChanged = func(s string) {
		tempRisk, err := decimal.NewFromString(s)
		if err != nil || tempRisk.GreaterThan(decimal.NewFromFloat(5.0)) || tempRisk.LessThan(decimal.NewFromFloat(0.5)) {
			pf.riskItem.HintText = "enter a 0.5 > risk < 5.0"
			pf.TPStratItem.Widget.(*widget.Select).Disable()
		} else {
			pf.riskItem.HintText = " "
			pf.TPStratItem.Widget.(*widget.Select).Enable()
		}
		pf.leftForm.Refresh()
	}

	return pf.riskItem
}

func (pf *planForm) makeTakeProfitStrategyItem() *widget.FormItem {
	tPStratSelect := widget.NewSelect(nil, nil)
	tPStratSelect.Options = cryptodb.TakeProfitStrategyStrings()
	tPStratSelect.Disable()

	if act.plan.ID != 0 {
		tPStratSelect.SetSelected(act.plan.TakeProfitStrategy.String())
		if act.orders[3].Status <= cryptodb.PartiallyFilled {
			tPStratSelect.Enable()
		}
	}

	tPStratSelect.OnChanged = func(s string) {
		pf.stopLossItem.Widget.(*FloatEntry).Enable()
		pf.rightForm.Refresh()
	}

	pf.TPStratItem = widget.NewFormItem("TP Strategy", tPStratSelect)
	pf.TPStratItem.HintText = " "

	return pf.TPStratItem
}

func (pf *planForm) makeStopLossItem() *widget.FormItem {
	StopLossFloatEntry := NewFloatEntry(act.pair.PriceScale, act.pair.Price.Tick)
	StopLossFloatEntry.Disable()
	pf.stopLossItem = widget.NewFormItem(fmt.Sprintf("Stop Loss (%s)", act.pair.QuoteCurrency), StopLossFloatEntry)
	pf.stopLossItem.HintText = " "

	if act.plan.ID != 0 {
		StopLossFloatEntry.SetText(act.orders[cryptodb.MarketStopLoss].Price.StringFixed(act.pair.PriceScale))
		if act.orders[cryptodb.MarketStopLoss].Status <= cryptodb.PartiallyFilled {
			StopLossFloatEntry.Enable()
		}
	}

	StopLossFloatEntry.OnChanged = func(s string) {
		_, err := decimal.NewFromString(s)
		if err != nil {
			pf.entryItem.Widget.(*FloatEntry).Disable()
			pf.stopLossItem.HintText = fmt.Sprintf("Enter a valid price in %s", act.pair.QuoteCurrency)
		} else {
			pf.entryItem.Widget.(*FloatEntry).Enable()
			pf.stopLossItem.HintText = " "
		}
		pf.rightForm.Refresh()
	}

	return pf.stopLossItem
}

func (pf *planForm) makeEntryItem() *widget.FormItem {
	entryFloatEntry := NewFloatEntry(act.pair.PriceScale, act.pair.Price.Tick)
	entryFloatEntry.Disable()
	pf.entryItem = widget.NewFormItem(fmt.Sprintf("Entry (%s)", act.pair.QuoteCurrency), entryFloatEntry)
	pf.entryItem.HintText = " "

	if act.plan.ID != 0 {
		entryFloatEntry.SetText(act.orders[cryptodb.Entry].Price.StringFixed(act.pair.PriceScale))
		if act.orders[cryptodb.Entry].Status <= cryptodb.PartiallyFilled {
			entryFloatEntry.Enable()
		}
	}

	entryFloatEntry.OnChanged = func(s string) {
		marketStopLossPrice := decimal.RequireFromString(pf.stopLossItem.Widget.(*FloatEntry).Text)
		entryPrice, err := decimal.NewFromString(s)
		switch {
		case err != nil:
			pf.TPStratItem.Widget.(*fyne.Container).Objects[0].(*widget.Select).Disable()
			pf.entryItem.HintText = fmt.Sprintf("enter a valid price in %s", act.pair.QuoteCurrency)
		case entryPrice.IsZero():
			pf.entryItem.HintText = " "
		case pf.directionItem.Widget.(*widget.RadioGroup).Selected == cryptodb.Long.String() && marketStopLossPrice.GreaterThanOrEqual(entryPrice):
			pf.entryItem.HintText = "must be higher than stoploss"
			pf.TPStratItem.Widget.(*widget.Select).Disable()
		case pf.directionItem.Widget.(*widget.RadioGroup).Selected == cryptodb.Short.String() && marketStopLossPrice.LessThanOrEqual(entryPrice):
			pf.entryItem.HintText = "must be lower than stoploss"
			pf.takeProfitItems[0].Widget.(*FloatEntry).Disable()
		default:
			pf.takeProfitItems[0].Widget.(*FloatEntry).Enable()
			pf.entryItem.HintText = " "
		}
		pf.rightForm.Refresh()
	}

	return pf.entryItem
}

func (pf *planForm) makeTakeProfitItem(n int, decimals int32, tick decimal.Decimal) *widget.FormItem {
	takeProfitFloatEntry := NewFloatEntry(decimals, tick)
	takeProfitFloatEntry.Disable()
	pf.takeProfitItems[n] = widget.NewFormItem(fmt.Sprintf("TP #%d (%s)", n+1, act.pair.QuoteCurrency), takeProfitFloatEntry)
	pf.takeProfitItems[n].HintText = " "

	if act.plan.ID != 0 {
		takeProfitFloatEntry.SetText(act.orders[3+n].Price.StringFixed(decimals))
		if act.orders[3+n].Status <= cryptodb.PartiallyFilled {
			takeProfitFloatEntry.Enable()
		}
	}

	takeProfitFloatEntry.OnChanged = func(s string) {
		var prevPrice decimal.Decimal
		var prevName string
		entryPrice := decimal.RequireFromString(pf.entryItem.Widget.(*FloatEntry).Text)
		if n == 0 {
			prevPrice = entryPrice
			prevName = "entry"
		} else {
			prevPrice = decimal.RequireFromString(pf.takeProfitItems[n-1].Widget.(*FloatEntry).Text)
			prevName = fmt.Sprintf("take profit #%d", n)
		}
		takeProfitPrice, err := decimal.NewFromString(s)
		switch {
		case err != nil:
			pf.takeProfitItems[n].HintText = fmt.Sprintf("enter a valid price in %s", act.pair.QuoteCurrency)
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Disable()
			}
		case takeProfitPrice.IsZero():
			pf.takeProfitItems[n].HintText = " "
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Disable()
			}
		case pf.directionItem.Widget.(*widget.RadioGroup).Selected == cryptodb.Long.String() && !takeProfitPrice.GreaterThanOrEqual(prevPrice):
			pf.takeProfitItems[n].HintText = "must be higher than " + prevName
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Disable()
			}
		case pf.directionItem.Widget.(*widget.RadioGroup).Selected == cryptodb.Short.String() && !takeProfitPrice.LessThanOrEqual(prevPrice):
			pf.takeProfitItems[n].HintText = "must be lower than " + prevName
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Disable()
			}
		default:
			pf.takeProfitItems[n].HintText = fmt.Sprintf("%.1f%% / %.1f%%",
				takeProfitPrice.Sub(prevPrice).Abs().Div(prevPrice).Mul(decimal.NewFromInt(100)).InexactFloat64(),
				takeProfitPrice.Sub(entryPrice).Abs().Div(entryPrice).Mul(decimal.NewFromInt(100)).InexactFloat64())
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Enable()
			}
		}

		pf.rightForm.Refresh()
	}

	return pf.takeProfitItems[n]
}

func (pf *planForm) makeTradingViewItem() *widget.FormItem {
	editEntry := widget.NewEntry()
	saveButton := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), nil)

	tradingViewLink := widget.NewHyperlink("", nil)
	createButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)

	if act.plan.TradingViewPlan != "" {
		tradingViewLink.URL, _ = url.Parse(act.plan.TradingViewPlan)
		tradingViewLink.SetText("Open")
		editEntry.Hide()
		saveButton.Hide()
	} else {
		tradingViewLink.Hide()
		createButton.Hide()
	}

	saveButton.OnTapped = func() {
		tvurl, err := url.Parse(editEntry.Text)
		if err == nil {
			saveButton.Hide()
			editEntry.Hide()
			createButton.Show()
			tradingViewLink.SetText("Open")
			tradingViewLink.Show()
			tradingViewLink.URL = tvurl
		}
	}

	createButton.OnTapped = func() {
		editEntry.SetText(tradingViewLink.URL.String())
		saveButton.Show()
		editEntry.Show()
		createButton.Hide()
		tradingViewLink.Hide()
	}

	buttonContainer := container.NewVBox(createButton, saveButton)
	urlContainer := container.NewVBox(editEntry, tradingViewLink)
	holdingContainer := container.NewBorder(nil, nil, nil, buttonContainer, urlContainer)

	pf.tradingViewPlanItem = widget.NewFormItem("TradingView", holdingContainer)
	pf.tradingViewPlanItem.HintText = " "

	return pf.tradingViewPlanItem
}

func (pf *planForm) makeNotesItem() *widget.FormItem {
	notesMultiLineEntry := widget.NewMultiLineEntry()
	notesMultiLineEntry.SetPlaceHolder("Enter notes...")
	if act.plan.Notes != "" {
		notesMultiLineEntry.SetText(act.plan.Notes)
	}
	if act.plan.Status == cryptodb.Logged {
		notesMultiLineEntry.Disable()
	}
	notesMultiLineEntry.Wrapping = fyne.TextWrapWord

	pf.notesItem = widget.NewFormItem("Notes", notesMultiLineEntry)

	return pf.notesItem
}

func (pf *planForm) makeToolBar() *widget.Toolbar {
	reviewAction := widget.NewToolbarAction(theme.ListIcon(), pf.reviewAction)
	historyAction := widget.NewToolbarAction(theme.HistoryIcon(), pf.historyAction)
	executeAction := widget.NewToolbarAction(theme.UploadIcon(), pf.executeAction)
	undoAction := widget.NewToolbarAction(theme.ContentUndoIcon(), pf.undoAction)
	okAction := widget.NewToolbarAction(theme.ConfirmIcon(), pf.okAction)

	return widget.NewToolbar(widget.NewToolbarSpacer(), reviewAction, historyAction, executeAction, undoAction, okAction)
}

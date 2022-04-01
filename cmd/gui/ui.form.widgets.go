package main

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func (pf *planForm) setQuoteCurrency(s string) {
	pf.stopLossItem.Text = fmt.Sprintf("Stop Loss (%s)", s)

	pf.entryItem.Text = fmt.Sprintf("Entry (%s)", s)

	for i, takeProfitItem := range pf.takeProfitItems {
		takeProfitItem.Text = fmt.Sprintf("Take Profit #%d (%s)", i+1, s)
	}
}

func (pf *planForm) setPriceScale(i int32, tick decimal.Decimal) {
	pf.stopLossItem.Widget.(*FloatEntry).decimals = i
	pf.stopLossItem.Widget.(*FloatEntry).tick = tick
	pf.stopLossItem.Widget.(*FloatEntry).SetPlaceHolder(decimal.Zero.StringFixed(i))

	pf.entryItem.Widget.(*FloatEntry).decimals = i
	pf.entryItem.Widget.(*FloatEntry).tick = tick
	pf.entryItem.Widget.(*FloatEntry).SetPlaceHolder(decimal.Zero.StringFixed(i))

	for _, takeProfitItem := range pf.takeProfitItems {
		takeProfitItem.Widget.(*FloatEntry).decimals = i
		takeProfitItem.Widget.(*FloatEntry).tick = tick
		takeProfitItem.Widget.(*FloatEntry).SetPlaceHolder(decimal.Zero.StringFixed(i))
	}
}

func (ui *UI) fzfPairs(s string) (possiblePairs []string) {
	var pairNames []string
	for _, p := range ui.Pairs {
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

func (pf *planForm) makePairItem() *widget.FormItem {
	CompletionEntry := xwidget.NewCompletionEntry([]string{})
	CompletionEntry.SetPlaceHolder("Select pair from list")
	CompletionEntry.Disable()
	item := widget.NewFormItem("Pair", CompletionEntry)
	item.HintText = " "

	CompletionEntry.OnChanged = func(s string) {
		CompletionEntry.SetText(strings.ToUpper(s))
		if len(s) >= 2 {
			possiblePairs := ui.fzfPairs(strings.ToUpper(s))
			if len(possiblePairs) == 1 {
				CompletionEntry.SetText(possiblePairs[0])
				CompletionEntry.HideCompletion()
			} else {
				CompletionEntry.SetOptions(possiblePairs)
				CompletionEntry.ShowCompletion()
			}
		}
	}

	CompletionEntry.OnSubmitted = func(s string) {
		s = strings.ToUpper(s)
		for _, p := range ui.Pairs {
			if s == p.Name {
				ui.activePair = p
				ui.activePlan.PairID = ui.activePair.ID
				pf.setQuoteCurrency(ui.activePair.QuoteCurrency)
				pf.setPriceScale(ui.activePair.PriceScale, ui.activePair.Price.Tick)

				pf.directionItem.Widget.(*widget.RadioGroup).Enable()
				pf.form.Refresh()
				break
			}
		}
	}

	return item
}

func (pf *planForm) makeDirectionItem() *widget.FormItem {
	directionRadio := widget.NewRadioGroup(nil, nil)
	item := widget.NewFormItem("Direction", directionRadio)
	item.HintText = " "
	directionRadio.Horizontal = true
	directionRadio.Disable()

	directionRadio.Options = cryptodb.DirectionStrings()
	directionRadio.OnChanged =
		func(s string) {
			pf.riskItem.Widget.(*FloatEntry).Enable()
			pf.form.Refresh()
		}

	return item
}

func (pf *planForm) makeRiskItem() *widget.FormItem {
	riskEntry := NewFloatEntry(1, decimal.NewFromFloat(0.1))
	riskEntry.Disable()
	riskEntry.SetPlaceHolder("0.0")
	item := widget.NewFormItem("Risk (%)", riskEntry)
	item.HintText = " "

	riskEntry.OnChanged = func(s string) {
		tempRisk, err := decimal.NewFromString(s)
		if err != nil || tempRisk.GreaterThan(decimal.NewFromFloat(5.0)) || tempRisk.LessThan(decimal.NewFromFloat(0.5)) {
			item.HintText = "enter a 0.5 > risk < 5.0"
			pf.stopLossItem.Widget.(*FloatEntry).Disable()
		} else {
			item.HintText = " "
			pf.stopLossItem.Widget.(*FloatEntry).Enable()
		}
		pf.form.Refresh()
	}

	return item
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

func (pf *planForm) makeStopLossItem(decimals int32, tick decimal.Decimal) *widget.FormItem {
	StopLossFloatEntry := NewFloatEntry(decimals, tick)
	StopLossFloatEntry.Disable()
	item := widget.NewFormItem("Stop Loss", StopLossFloatEntry)
	item.HintText = " "

	StopLossFloatEntry.OnChanged = func(s string) {
		_, err := decimal.NewFromString(s)
		if err != nil {
			pf.entryItem.Widget.(*FloatEntry).Disable()
			item.HintText = fmt.Sprintf("Enter a valid price in %s", ui.activePair.QuoteCurrency)
		} else {
			pf.entryItem.Widget.(*FloatEntry).Enable()
			item.HintText = " "
		}
		pf.form.Refresh()
	}

	return item
}

func (pf *planForm) makeEntryItem(decimals int32, tick decimal.Decimal) *widget.FormItem {
	entryFloatEntry := NewFloatEntry(decimals, tick)
	entryFloatEntry.Disable()
	item := widget.NewFormItem("Entry (%s)", entryFloatEntry)
	item.HintText = " "

	entryFloatEntry.OnChanged = func(s string) {
		marketStopLossPrice := decimal.RequireFromString(pf.stopLossItem.Widget.(*FloatEntry).Text)
		entryPrice, err := decimal.NewFromString(s)
		switch {
		case err != nil:
			pf.TPStratItem.Widget.(*fyne.Container).Objects[0].(*widget.Select).Disable()
			item.HintText = fmt.Sprintf("enter a valid price in %s", ui.activePair.QuoteCurrency)
		case entryPrice.IsZero():
			item.HintText = " "
		case pf.directionItem.Widget.(*widget.RadioGroup).Selected == cryptodb.Long.String() && marketStopLossPrice.GreaterThanOrEqual(entryPrice):
			item.HintText = "must be higher than stoploss"
			pf.TPStratItem.Widget.(*widget.Select).Disable()
		case pf.directionItem.Widget.(*widget.RadioGroup).Selected == cryptodb.Short.String() && marketStopLossPrice.LessThanOrEqual(entryPrice):
			item.HintText = "must be lower than stoploss"
			pf.TPStratItem.Widget.(*widget.Select).Disable()
		default:
			pf.TPStratItem.Widget.(*widget.Select).Enable()
			item.HintText = " "
		}
		pf.form.Refresh()
	}

	return item
}

func (pf *planForm) makeTakeProfitStrategyItem() *widget.FormItem {
	tPStratSelect := widget.NewSelect(nil, nil)
	tPStratSelect.Disable()

	tPStratSelect.Options = cryptodb.TakeProfitStrategyStrings()

	tPStratSelect.OnChanged = func(s string) {
		pf.takeProfitItems[0].Widget.(*FloatEntry).Enable()
		pf.form.Refresh()
	}

	item := widget.NewFormItem("TP Strategy", tPStratSelect)
	item.HintText = " "

	return item
}

func (pf *planForm) makeTakeProfitItem(n int, decimals int32, tick decimal.Decimal) *widget.FormItem {
	takeProfitFloatEntry := NewFloatEntry(decimals, tick)
	takeProfitFloatEntry.Disable()
	item := widget.NewFormItem("Take profit #", takeProfitFloatEntry)
	item.HintText = " "

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
			item.HintText = fmt.Sprintf("enter a valid price in %s", ui.activePair.QuoteCurrency)
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Disable()
			}
		case takeProfitPrice.IsZero():
			item.HintText = " "
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Disable()
			}
		case pf.directionItem.Widget.(*widget.RadioGroup).Selected == cryptodb.Long.String() && !takeProfitPrice.GreaterThanOrEqual(prevPrice):
			item.HintText = "must be higher than " + prevName
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Disable()
			}
		case pf.directionItem.Widget.(*widget.RadioGroup).Selected == cryptodb.Short.String() && !takeProfitPrice.LessThanOrEqual(prevPrice):
			item.HintText = "must be lower than " + prevName
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Disable()
			}
		default:
			item.HintText = fmt.Sprintf("%.1f%% / %.1f%%",
				takeProfitPrice.Sub(prevPrice).Abs().Div(prevPrice).Mul(decimal.NewFromInt(100)).InexactFloat64(),
				takeProfitPrice.Sub(entryPrice).Abs().Div(entryPrice).Mul(decimal.NewFromInt(100)).InexactFloat64())
			if n != cryptodb.MaxTakeProfits-1 {
				pf.takeProfitItems[n+1].Widget.(*FloatEntry).Enable()
			}
		}

		pf.form.Refresh()
	}

	return item
}

func (pf *planForm) makeTradingViewLinkItem() *widget.FormItem {
	editEntry := widget.NewEntry()
	saveButton := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), nil)

	tradingViewLink := widget.NewHyperlink("", nil)
	createButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)

	saveButton.Hide()
	editEntry.Hide()
	createButton.Show()
	tradingViewLink.Show()

	saveButton.OnTapped = func() {
		tvurl, err := url.Parse(editEntry.Text)
		if err == nil {
			saveButton.Hide()
			editEntry.Hide()
			createButton.Show()
			tradingViewLink.SetText(tvurl.String())
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

	item := widget.NewFormItem("TradingView", holdingContainer)
	item.HintText = " "

	return item
}

func (pf *planForm) makeNotesMultilineItem() *widget.FormItem {
	notesMultiLineEntry := widget.NewMultiLineEntry()
	notesMultiLineEntry.SetPlaceHolder("Enter notes...")
	notesMultiLineEntry.SetText(ui.activePlan.Notes)
	notesMultiLineEntry.Wrapping = fyne.TextWrapWord

	item := widget.NewFormItem("Notes", notesMultiLineEntry)
	item.HintText = " "

	return item
}

func (pf *planForm) makeToolBar() *widget.Toolbar {
	logAction := widget.NewToolbarAction(theme.DocumentIcon(), pf.logAction)
	executeAction := widget.NewToolbarAction(theme.UploadIcon(), pf.executeAction)
	cancelAction := widget.NewToolbarAction(theme.CancelIcon(), pf.cancelAction)
	okAction := widget.NewToolbarAction(theme.ConfirmIcon(), pf.okAction)

	return widget.NewToolbar(widget.NewToolbarSpacer(), logAction, executeAction, cancelAction, okAction)
}

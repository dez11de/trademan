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

func (pf *planForm) setPriceScale(i int32) {
	pf.stopLossItem.Widget.(*widget.Entry).SetPlaceHolder(decimal.Zero.StringFixed(i))

	pf.entryItem.Widget.(*widget.Entry).SetPlaceHolder(decimal.Zero.StringFixed(i))

	for _, takeProfitItem := range pf.takeProfitItems {
		takeProfitItem.Widget.(*widget.Entry).SetPlaceHolder(decimal.Zero.StringFixed(i))
	}
}

func (ui *UI) fzfPairs(s string) (possiblePairs []string) {
	var pairNames []string
	for _, p := range ui.Pairs {
		pairNames = append(pairNames, p.Name)
	}

	matches := fuzzy.RankFind(s, pairNames)
	sort.Sort(matches)

	for _, p := range matches {
		possiblePairs = append(possiblePairs, p.Target)
	}

	return possiblePairs
}

func (pf *planForm) makePairItem() *widget.FormItem {
	CompletionEntry := xwidget.NewCompletionEntry([]string{})
	CompletionEntry.SetPlaceHolder("Select pair from list")
	item := widget.NewFormItem("Pair", CompletionEntry)
	item.HintText = " "

	CompletionEntry.OnChanged = func(s string) {
		CompletionEntry.SetText(strings.ToUpper(s))
		possiblePairs := ui.fzfPairs(strings.ToUpper(s))

		if len(possiblePairs) == 1 {
			CompletionEntry.SetText(possiblePairs[0])
			CompletionEntry.HideCompletion()
		} else if len(s) >= 2 {
			CompletionEntry.SetOptions(possiblePairs)
			CompletionEntry.ShowCompletion()
		}
	}

	CompletionEntry.OnSubmitted = func(s string) {
		s = strings.ToUpper(s)
		for _, p := range ui.Pairs {
			if s == p.Name {
				ui.activePair = p
				ui.activePlan.PairID = ui.activePair.ID
				pf.setQuoteCurrency(ui.activePair.QuoteCurrency)
				pf.setPriceScale(ui.activePair.PriceScale)

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
			pf.riskItem.Widget.(*widget.Entry).Enable()
			pf.form.Refresh()
		}

	return item
}

func (pf *planForm) makeRiskItem() *widget.FormItem {
	riskEntry := widget.NewEntry()
	riskEntry.Disable()
	riskEntry.SetPlaceHolder("0.0")
	item := widget.NewFormItem("Risk (%)", riskEntry)
	item.HintText = " "

	riskEntry.OnChanged = func(s string) {
		tempRisk, err := decimal.NewFromString(s)
		if err != nil || tempRisk.GreaterThan(decimal.NewFromFloat(5.0)) || tempRisk.LessThan(decimal.NewFromFloat(0.5)) {
			item.HintText = "enter a 0.5 > risk < 5.0"
			pf.stopLossItem.Widget.(*widget.Entry).Disable()
		} else {
			item.HintText = " "
			pf.stopLossItem.Widget.(*widget.Entry).Enable()
		}
		pf.form.Refresh()
	}

	return item
}

func (pf *planForm) makeStopLossItem() *widget.FormItem {
	StopLossEntry := widget.NewEntry()
	StopLossEntry.Disable()
	item := widget.NewFormItem("Stop Loss", StopLossEntry)
	item.HintText = " "

	StopLossEntry.OnChanged = func(s string) {
		_, err := decimal.NewFromString(s)
		if err != nil {
			pf.entryItem.Widget.(*widget.Entry).Disable()
			item.HintText = fmt.Sprintf("Enter a valid price in %s", ui.activePair.QuoteCurrency)
		} else {
			pf.entryItem.Widget.(*widget.Entry).Enable()
			item.HintText = " "
		}

		pf.form.Refresh()
	}

	return item
}

func (pf *planForm) makeEntryItem() *widget.FormItem {
	entryEntry := widget.NewEntry()
	entryEntry.Disable()
	item := widget.NewFormItem("Entry (%s)", entryEntry)
	item.HintText = " "

	entryEntry.OnChanged = func(s string) {
		_, err := decimal.NewFromString(s)
		if err != nil {
			pf.TPStratItem.Widget.(*widget.Select).Enable()
			item.HintText = fmt.Sprintf("Enter a valid price in %s", ui.activePair.QuoteCurrency)
		} else {
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
	item := widget.NewFormItem("TP Strategy", tPStratSelect)
	item.HintText = " "

	tPStratSelect.Options = cryptodb.TakeProfitStrategyStrings()

	tPStratSelect.OnChanged = func(s string) {
		pf.takeProfitItems[0].Widget.(*widget.Entry).Enable()
		pf.form.Refresh()
	}

	return item
}

func (pf *planForm) makeTakeProfitItem(n int) *widget.FormItem {
	takeProfitEntry := widget.NewEntry()
	takeProfitEntry.Disable()
	item := widget.NewFormItem("Take profit #", takeProfitEntry)
	item.HintText = " "

	takeProfitEntry.OnChanged = func(s string) {
		if n != cryptodb.MaxTakeProfits-1 {
			v, err := decimal.NewFromString(s)
			if err != nil || v.IsZero() {
				pf.takeProfitItems[n+1].Widget.(*widget.Entry).Disable()
			} else {
				pf.takeProfitItems[n+1].Widget.(*widget.Entry).Enable()
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
	executeAction := widget.NewToolbarAction(theme.UploadIcon(), pf.executeAction)
	cancelAction := widget.NewToolbarAction(theme.CancelIcon(), pf.cancelAction)
	okAction := widget.NewToolbarAction(theme.ConfirmIcon(), pf.okAction)

	return widget.NewToolbar(widget.NewToolbarSpacer(), executeAction, cancelAction, okAction)
}

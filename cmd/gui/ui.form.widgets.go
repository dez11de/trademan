package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"github.com/bart613/decimal"
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

func (pf *planForm) makeStatContainer() *fyne.Container {
	// TODO: make distinction between start RRR and evolved RRR
	startRewardRiskRatioLabel := widget.NewLabel("Start RRR: ")
	startRewardRiskRatioValue := widget.NewLabel(fmt.Sprintf("%.1f", 0.0))
	evolvedRewardRiskRatioLabel := widget.NewLabel("Current RRR: ")
	evolvedRewardRiskRatioValue := widget.NewLabel(fmt.Sprintf("%.1f", 0.0))

	currentPnLLabel := widget.NewLabel("PnL: ")
	currentPnLValue := widget.NewLabel(fmt.Sprintf("%s%%", pf.plan.Profit.StringFixed(1))) // TODO: should be relative to entrySize.
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
	CompletionEntry := xwidget.NewCompletionEntry([]string{})
	CompletionEntry.SetPlaceHolder("Select pair from list")
	CompletionEntry.OnChanged = func(s string) {
		if len(s) < 1 {
			CompletionEntry.HideCompletion()
		}
		possiblePairs, err := searchPairs(strings.ToUpper(s))
		if err != nil {
			return
		}
		if len(possiblePairs) == 1 {
			CompletionEntry.SetText(possiblePairs[0])
			CompletionEntry.HideCompletion()
		} else {
			CompletionEntry.SetOptions(possiblePairs)
			CompletionEntry.ShowCompletion()
		}
	}

	CompletionEntry.OnSubmitted = func(s string) {
		ok := false
		s = strings.ToUpper(s)
		for _, p := range pf.PairCache {
			if s == p.Name {
				pf.activePair = p
				ok = true
				break // No need to look further
			}
		}
		if ok {
			pf.plan.PairID = pf.activePair.ID
			pf.setQuoteCurrency(pf.activePair.QuoteCurrency)
			pf.setPriceScale(int32(pf.activePair.PriceScale))
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
	riskEntry := widget.NewEntry()
	riskEntry.Disable()
	riskEntry.SetPlaceHolder(pf.plan.Risk.StringFixed(2))

	riskEntry.OnChanged = func(s string) {
		tempRisk, err := decimal.NewFromString(s)
		if err != nil || tempRisk.Cmp(decimal.NewFromFloat(5)) != -1 || tempRisk.Cmp(decimal.NewFromFloat(0.499)) != 1 {
			log.Printf("Invalid risk %s", s)
			pf.riskItem.HintText = "enter a valid risk"
			pf.stopLossItem.Widget.(*widget.Entry).Disable()
			pf.form.Refresh()
		} else {
			pf.riskItem.HintText = ""
			pf.stopLossItem.Widget.(*widget.Entry).Enable()
			pf.form.Refresh()
		}
	}

	item := widget.NewFormItem("Risk (%)", riskEntry)
	item.HintText = " "

	return item
}

func (pf *planForm) makeStopLossItem() *widget.FormItem {
	StopLossEntry := widget.NewEntry()
	StopLossEntry.Disable()
	StopLossEntry.OnChanged = func(s string) {
		// TODO: properly validate input
		pf.entryItem.Widget.(*widget.Entry).Enable()
		pf.form.Refresh()
	}
	item := widget.NewFormItem(fmt.Sprintf("Stop Loss (%s)", pf.activePair.QuoteCurrency), StopLossEntry)
	item.HintText = " "
	return item
}

func (pf *planForm) makeEntryItem() *widget.FormItem {
	entryEntry := widget.NewEntry()
	entryEntry.Disable()
	entryEntry.OnChanged = func(s string) {
		// TODO: properly validate input
		pf.takeProfitItems[0].Widget.(*widget.Entry).Enable()
		pf.form.Refresh()
	}
	item := widget.NewFormItem(fmt.Sprintf("Entry (%s)", pf.activePair.QuoteCurrency), entryEntry)
	item.HintText = " "
	return item
}

func (pf *planForm) makeTakeProfitItem(n int) *widget.FormItem {
	takeProfitEntry := widget.NewEntry()
	takeProfitEntry.Disable()
	takeProfitEntry.OnChanged = func(s string) {
		// TODO: properly validate input
		// TODO: show % difference with entry and previous take profit
		if !decimal.RequireFromString(s).IsZero() && s != "" {
			pf.takeProfitItems[n+1].Widget.(*widget.Entry).Enable()
			if n >= 1 {
				pf.form.Enable()
			}
			pf.form.Refresh()
		}
	}
	item := widget.NewFormItem(fmt.Sprintf("Take profit #%d (%s)", n+1, pf.activePair.QuoteCurrency), takeProfitEntry)
	item.HintText = " "
	return item
}

// TODO: think about in which statusses changing is allowed
func (pf *planForm) makeTradingViewLinkItem() *widget.FormItem {
	var tvurl url.URL

	editEntry := widget.NewEntry()
	saveButton := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		log.Printf("Save button clicked.")
	})
	editContainer := container.NewHBox(editEntry, saveButton)
	editContainer.Hide()

	tvurl.Parse(pf.plan.TradingViewPlan)
	tradingViewLink := widget.NewHyperlink("Open", &tvurl)
	createButton := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		log.Printf("Create button clicked.")
	})
	showContainer := container.NewHBox(tradingViewLink, createButton)
	showContainer.Hide()

	holdingContainer := container.NewWithoutLayout(editContainer, showContainer)

	return widget.NewFormItem("TradingView", holdingContainer)
}

func (pf *planForm) makeNotesMultilineItem() *widget.FormItem {
	notesMultiLineEntry := widget.NewMultiLineEntry()
	notesMultiLineEntry.SetPlaceHolder("Enter notes...")
	notesMultiLineEntry.SetText(pf.plan.Notes)
	notesMultiLineEntry.Wrapping = fyne.TextWrapWord

	return widget.NewFormItem("Notes", notesMultiLineEntry)
}

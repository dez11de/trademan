package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makePairSelect(d *Database, p Plan, a Pair) *widget.FormItem {
	var pairOptions []string
	for _, p := range d.PairCache {
		pairOptions = append(pairOptions, p.Pair)
	}
	pairSelector := widget.NewSelect(pairOptions,
		func(s string) {
		})
	pairSelector.Selected = a.Pair
	return widget.NewFormItem("Pair", pairSelector)
}

func makeSideRadio(p Plan) *widget.FormItem {
	sideSelection := widget.NewRadioGroup([]string{"Long", "Short"},
		func(s string) {
		})
	sideSelection.SetSelected(p.Side.String())
	sideSelection.Required = true
	return widget.NewFormItem("Side", sideSelection)
}

func makeRiskEntry(p Plan) *widget.FormItem {
	riskField := widget.NewEntry()
	return widget.NewFormItem("Risk (%)", riskField)
}

func makeStopLossEntry(p Plan, a Pair, sl Order) *widget.FormItem {
	hardStopLossEntry := widget.NewEntry()
	hardStopLossEntry.SetText(sl.Price.StringFixed(a.PriceScale))
	return widget.NewFormItem(fmt.Sprintf("Stop Loss (%s)", a.QuoteCurrency), hardStopLossEntry)
}

func makeEntryEntry(p Plan, a Pair, e Order) *widget.FormItem {
	entryEntry := widget.NewEntry()
	entryEntry.SetText(e.Price.StringFixed(a.PriceScale))
	return widget.NewFormItem(fmt.Sprintf("Entry (%s)", a.QuoteCurrency), entryEntry)
}

func makeTakeProfitEntry(p Plan, a Pair, takeProfit Order, n int) *widget.FormItem {
	takeProfitEntry := widget.NewEntry()
	takeProfitEntry.Text = takeProfit.Price.StringFixed(a.PriceScale)
	return widget.NewFormItem(fmt.Sprintf("Take profit #%d (%s)", n+1, a.QuoteCurrency), takeProfitEntry)
}

func makeTakeProfitEntries(p Plan, a Pair, takeProfits Orders) (takeProfitEntries []*widget.FormItem) {
	for i, takeProfit := range takeProfits {
		takeProfitEntries = append(takeProfitEntries, makeTakeProfitEntry(p, a, takeProfit, i))
	}
	return takeProfitEntries
}

func makePlanForm(d *Database, p Plan) fyne.CanvasObject {
	a, _ := d.GetPairFromID(p.PairID)
	orders, _ := d.GetPlanOrders(p.PlanID)
	sl, _ := orders.GetHardStopLoss()
	e, _ := orders.GetEntry()
	tps, _ := orders.GetTakeProfits()

	pairSelect := makePairSelect(d, p, a)
	sideRadio := makeSideRadio(p)
	riskEntry := makeRiskEntry(p)
	stopLossEntry := makeStopLossEntry(p, a, sl)
	entryEntry := makeEntryEntry(p, a, e)
	if len(tps) == 0 {
		newTakeProfitOrder := Order{PlanID: p.PlanID, OrderType: TakeProfit}
		orders = append(orders, newTakeProfitOrder)
		tps = append(tps, newTakeProfitOrder)
	}
	takeProfitEntries := makeTakeProfitEntries(p, a, tps)

	form := widget.NewForm(pairSelect, sideRadio, riskEntry, stopLossEntry, entryEntry)
	for _, tpi := range takeProfitEntries {
		form.AppendItem(tpi)
	}

	addTakeProfitButton := widget.NewButton("Add TP", func() {
		log.Printf("Add TP button pressed")
		newTakeProfitOrder := Order{PlanID: p.PlanID, OrderType: TakeProfit}
		orders = append(orders, newTakeProfitOrder)
		tps = append(tps, newTakeProfitOrder)
		form.AppendItem(makeTakeProfitEntry(p, a, newTakeProfitOrder, len(tps)-1))
	})

	okButton := widget.NewButton("Ok", func() {
		log.Printf("Ok button pressed")
	})
	delButton := widget.NewButton("Delete", func() {
		log.Printf("Delete button pressed")
	})

	buttons := container.NewHBox(okButton, delButton, addTakeProfitButton)
	formPage := container.NewVBox(form, buttons)

	return formPage
}

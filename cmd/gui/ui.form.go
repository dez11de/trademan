package main

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
)

func (pf *planForm) gatherSetup() cryptodb.Setup {
	pf.plan.PairID = pf.activePair.ID
	pf.plan.Direction.Scan(pf.sideItem.Widget.(*widget.RadioGroup).Selected)
	pf.plan.Risk = decimal.RequireFromString(pf.riskItem.Widget.(*widget.Entry).Text)
	pf.orders[cryptodb.KindMarketStopLoss].Price = decimal.RequireFromString(pf.stopLossItem.Widget.(*widget.Entry).Text)
	pf.orders[cryptodb.KindEntry].Price = decimal.RequireFromString(pf.entryItem.Widget.(*widget.Entry).Text)
	for i := 0; i < cryptodb.MaxTakeProfits; i++ {
		tempPrice, err := decimal.NewFromString(pf.takeProfitItems[i].Widget.(*widget.Entry).Text)
		if err == nil {
			pf.orders[3+i].Price = tempPrice
		} else {
			pf.orders[3+i].Price = decimal.Zero
		}
	}

	return cryptodb.Setup{Plan: pf.plan, Orders: pf.orders}
}

func (pf *planForm) okAction() {
	setup := pf.gatherSetup()
	setup, err := storeSetup(setup)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}
	pf.plan = setup.Plan
	pf.orders = setup.Orders
	pf.form.Refresh()
}

func (pf *planForm) cancelAction() {

}

func (pf *planForm) executeAction() {
	setup := pf.gatherSetup()
	setup, err := storeSetup(setup)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}

	plan, err := executePlan(setup.Plan.ID)
	pf.plan = plan
	pf.orders = setup.Orders
	pf.form.Refresh()
}

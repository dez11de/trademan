package main

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
)

func (pf *planForm) gatherSetup() cryptodb.Setup {
	ui.activePlan.PairID = ui.activePair.ID
	ui.activePlan.Direction.Scan(pf.directionItem.Widget.(*widget.RadioGroup).Selected)
	ui.activePlan.Risk = decimal.RequireFromString(pf.riskItem.Widget.(*widget.Entry).Text)
	ui.activeOrders[cryptodb.KindMarketStopLoss].Price = decimal.RequireFromString(pf.stopLossItem.Widget.(*widget.Entry).Text)
	ui.activeOrders[cryptodb.KindEntry].Price = decimal.RequireFromString(pf.entryItem.Widget.(*widget.Entry).Text)
	for i := 0; i < cryptodb.MaxTakeProfits; i++ {
		// TODO: make this more robust
		tempPrice, err := decimal.NewFromString(pf.takeProfitItems[i].Widget.(*widget.Entry).Text)
		if err == nil {
			ui.activeOrders[3+i].Price = tempPrice
		} else {
			ui.activeOrders[3+i].Price = decimal.Zero
		}
	}

	return cryptodb.Setup{Plan: ui.activePlan, Orders: ui.activeOrders}
}

func (pf *planForm) okAction() {
	setup := pf.gatherSetup()
	setup, err := storeSetup(setup)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}
	ui.activePlan = setup.Plan
	ui.activeOrders = setup.Orders
	ui.Plans, _ = getPlans()
	ui.List.Refresh()
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
	ui.activePlan = plan
	ui.activeOrders = setup.Orders
	ui.Plans, _ = getPlans()
	ui.List.Refresh()
}

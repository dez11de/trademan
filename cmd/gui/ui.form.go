package main

import (
	"log"

	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
	"github.com/bart613/decimal"
)

func (pf *planForm) gatherSetup() cryptodb.Setup {
	pf.plan.PairID = pf.activePair.ID
	pf.plan.Direction.Scan(pf.sideItem.Widget.(*widget.RadioGroup).Selected)
	pf.plan.Risk = decimal.RequireFromString(pf.riskItem.Widget.(*widget.Entry).Text)
	pf.orders[cryptodb.KindMarketStopLoss].Price = decimal.RequireFromString(pf.stopLossItem.Widget.(*widget.Entry).Text)
	pf.orders[cryptodb.KindEntry].Price = decimal.RequireFromString(pf.entryItem.Widget.(*widget.Entry).Text)
	for i := 0; i < 5-1; i++ { // TODO: restore MaxTakeProfits to it's former glory
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
    log.Printf("Submit button pressed.")
    setup := pf.gatherSetup()
	storeSetup(setup)
}

func (pf *planForm) cancelAction() {
    log.Print("not implemented")
}

func (pf *planForm) executeAction() {
    log.Printf("Execute button pressed.")
    setup := pf.gatherSetup()
	storeSetup(setup)
    executePlan(setup.Plan.ID)
}

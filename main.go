package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type UI struct {
	app           *tview.Application
	loggingWindow *tview.TextView
	positionTable *tview.Table
	positionForm  *tview.Form
	flex          *tview.Flex
	pages         *tview.Pages
}

func main() {
	ui := SetupMain()

	db := NewDB()
	dberr := db.Connect()
	if dberr != nil {
		fmt.Fprintf(ui.loggingWindow, "Database error: %+v\n", dberr)
	}

	t := Position{
		Symbol:          "ADAUSDT",
		Side:            planned,
		Risk:            0.5,
		Size:            123.4,
		EntryPrice:      1.23,
		HardStopLoss:    1.22,
		Notes:           "Dit is maar gewoon een test.\nIk heb nog niet echt een plan.\n",
		TradingViewPlan: "https://tradingview.com/1234asfdf",
		RewardRiskRatio: 9.06,
	}

	t.TradeID, dberr = db.AddPosition(t)

	o := Order{
		TradeID:           t.TradeID,
		ExchangeOrderID:   "somestring",
		Status:            planned,
		OrderType:         TakeProfit,
		Quantity:          33.3,
		TakeProfitTrigger: 1.29,
		Price:             1.30,
	}

	o.OrderID, dberr = db.AddOrder(o)

	if dberr != nil {
		fmt.Printf("Error adding order %v", dberr)
	}

	db.AddLog(t.TradeID, user, "First test log.")

	if err := ui.app.Run(); err != nil {
		panic(err)
	}
}

package main

import (
	"fmt"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

func matchExchangeOrder(eID string) (o cryptodb.Order, err error) {
	// TODO: marketStopLossOrder ExchangeOrderID is the same as entryOrder ExchangeOrderID. Figure out how to deal with this
	result := db.Where("exchange_order_id = ?", eID).Find(&o)

	return o, result.Error
}

func processOpenOrder(o exchange.Order) (err error) {
	openOrder, err := matchExchangeOrder(o.ExchangeOrderID)
	if err != nil {
		return err
	}

	var marketStopLossOrder cryptodb.Order
	db.Where("order_kind = ?", cryptodb.MarketStopLoss).Where("plan_id = ?", openOrder.PlanID).Find(&marketStopLossOrder)

	openOrder.ExchangeOrderID = o.ExchangeOrderID
	marketStopLossOrder.ExchangeOrderID = o.ExchangeOrderID
	openOrder.Status.Scan(o.OrderStatus)

	var stopLossSetMsg string
	if marketStopLossOrder.Price.Equal(o.StopLoss) {
		stopLossSetMsg = "and"
	} else {
		stopLossSetMsg = "but DID NOT" // TODO: this should NEVER happen.
	}

	marketStopLossOrder.Status.Scan(o.OrderStatus)

	var tmLog cryptodb.Log
	tmLog.PlanID = openOrder.PlanID
	tmLog.Source = cryptodb.Exchange

	tmLog.Text = fmt.Sprintf("Exchange processed %s order %d %s set stoploss as OrderID: %s, at %s.", openOrder.Status.String(), openOrder.ID, stopLossSetMsg, openOrder.ExchangeOrderID, o.CreatedAt.Format("2006-01-02 15:04:05.000"))

	tx := db.Begin()

	result := tx.Save(&openOrder)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Save(&marketStopLossOrder)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Create(&tmLog)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Commit()

	return result.Error
}

func processCloseOrder(o exchange.Order) (err error) {
	closeOrder, err := matchExchangeOrder(o.ExchangeOrderID)
	if err != nil {
		return err
	}

	tx := db.Begin()

	var tmLog cryptodb.Log
	closeOrder.ExchangeOrderID = o.ExchangeOrderID
	closeOrder.Status.Scan(o.OrderStatus)
	tmLog.PlanID = closeOrder.PlanID
	tmLog.Source = cryptodb.Exchange
	tmLog.Text = fmt.Sprintf("Exchange processed %s order %d as OrderID: %s, at %s.", closeOrder.Status.String(), closeOrder.ID, closeOrder.ExchangeOrderID, o.CreatedAt.Format("2006-01-02 15:04:05.000"))

	result := tx.Save(&closeOrder)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Create(&tmLog)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Commit()

	return result.Error
}

func processOrder(o exchange.Order) (err error) {
	// TODO: check if this works for Short orders the same way or, in reverse
	if o.Side == "Sell" {
		err := processCloseOrder(o)
		return err
	}

	if o.Side == "Buy" {
		err := processOpenOrder(o)
		return err
	}

	return nil
}

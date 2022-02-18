package main

import (
	"fmt"
	"log"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

func matchExchangeOrder(exchangeID string) (o cryptodb.Order, err error) {
	result := db.Where("exchange_order_id = ?", exchangeID).Last(&o)

	return o, result.Error
}

func processEntryOrder(o exchange.Order) (err error) {
	entryOrder, err := matchExchangeOrder(o.ExchangeOrderID)
	if err != nil {
		return err
	}

	var marketStopLossOrder cryptodb.Order
	db.Where("order_kind = ?", cryptodb.MarketStopLoss).Where("plan_id = ?", entryOrder.PlanID).Find(&marketStopLossOrder)

	entryOrder.ExchangeOrderID = o.ExchangeOrderID
	marketStopLossOrder.ExchangeOrderID = o.ExchangeOrderID
	entryOrder.Status.Scan(o.OrderStatus)

	var logEntry cryptodb.Log
	logEntry.PlanID = entryOrder.PlanID
	logEntry.Source = cryptodb.Exchange

	if entryOrder.Status.String() == "New" {
		var stopLossSetMsg string
		if marketStopLossOrder.Price.Equal(o.StopLoss) {
			stopLossSetMsg = "and"
		} else {
			stopLossSetMsg = "but DID NOT" // TODO: this should NEVER happen.
		}
		marketStopLossOrder.Status.Scan(o.OrderStatus)

		logEntry.Text = fmt.Sprintf("Exchange processed %s order %d %s set stoploss as OrderID: %s, at %s.", entryOrder.Status.String(), entryOrder.ID, stopLossSetMsg, entryOrder.ExchangeOrderID, o.CreatedAt.Format("2006-01-02 15:04:05.000"))
	}

	tx := db.Begin()

	result := tx.Save(&entryOrder)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Save(&marketStopLossOrder)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Create(&logEntry)
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

	var logEntry cryptodb.Log
	closeOrder.ExchangeOrderID = o.ExchangeOrderID
	closeOrder.Status.Scan(o.OrderStatus)
	logEntry.PlanID = closeOrder.PlanID
	logEntry.Source = cryptodb.Exchange
	logEntry.Text = fmt.Sprintf("Exchange processed %s order %d as OrderID: %s, at %s.", closeOrder.Status.String(), closeOrder.ID, closeOrder.ExchangeOrderID, o.CreatedAt.Format("2006-01-02 15:04:05.000"))

	result := tx.Save(&closeOrder)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Create(&logEntry)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.Commit()

	return result.Error
}

func processOrder(o exchange.Order) (err error) {
	// TODO: check if this works for Short orders the same way or, in reverse
	// this is Entry with MarketStopLoss
	if !o.StopLoss.IsZero() {
		log.Printf("Processing Entry w/ MarketStopLoss")
	}
	if o.Side == "Sell" {
		err := processCloseOrder(o)
		return err
	}

	if o.Side == "Buy" {
		err := processEntryOrder(o)
		return err
	}

	return nil
}

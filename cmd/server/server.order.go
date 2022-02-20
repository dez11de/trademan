package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

func PlaceOrders(p cryptodb.Plan, activePair cryptodb.Pair, ticker exchange.Ticker, o []cryptodb.Order) (err error) {

	// TODO: move to it's own routine and log it
	switch p.Direction {
	case cryptodb.Long:
		if p.Leverage.GreaterThan(activePair.Leverage.Buy) {
			log.Printf("Long leverage updated to %s", p.Leverage.String())
			SetLeverage(p, &activePair)
		}
	case cryptodb.Short:
		if p.Leverage.GreaterThan(activePair.Leverage.Sell) {
			log.Printf("Short leverage updated to %s", p.Leverage.String())
			SetLeverage(p, &activePair)
		}
	}

	err = placeEntry(p, activePair, &o[cryptodb.MarketStopLoss], o[cryptodb.LimitStopLoss], &o[cryptodb.Entry])
	if err != nil {
		return err
	}

	db.Save(o[cryptodb.MarketStopLoss])
	db.Save(o[cryptodb.Entry])

	err = placeLimitStopLoss(p, activePair, ticker, o[cryptodb.MarketStopLoss], &o[cryptodb.LimitStopLoss], o[cryptodb.Entry])
	if err != nil {
		return err
	}
	db.Save(o[cryptodb.LimitStopLoss])

	for i := 3; i < 3+cryptodb.MaxTakeProfits; i++ {
		if !o[i].Price.IsZero() {
			err = placeTakeProfit(p, activePair, ticker, o[cryptodb.MarketStopLoss], o[cryptodb.Entry], &o[i])
			if err != nil {
				return err
			}
			db.Save(o[i])
		}
	}

	return nil
}

func SetLeverage(p cryptodb.Plan, pair *cryptodb.Pair) (err error) {

	var result exchange.LeverageResponse
	levParams := make(exchange.RequestParameters)

	levParams["symbol"] = pair
	switch p.Direction {
	case cryptodb.Long:
		levParams["buy_leverage"] = p.Leverage.InexactFloat64()
		pair.Leverage.Buy = p.Leverage
	case cryptodb.Short:
		levParams["sell"] = p.Leverage.InexactFloat64()
		pair.Leverage.Sell = p.Leverage
	}

	_, _, err = e.SignedRequest(http.MethodPost, "/private/linear/position/set-leverage", levParams, &result)
	log.Printf("Result of setting leverage: %+v", result)
	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		return errors.New(result.ExtendedCode)
	}

	logEntry := &cryptodb.Log{
		PlanID: p.ID,
		Source: cryptodb.Server,
		Text:   fmt.Sprintf("Set %s leverage to %s", p.Direction.String(), p.Leverage.String()),
	}
	// TODO: should check if it was accepted by exchange
	db.Save(pair)
	db.Create(logEntry)

	return nil
}

func placeEntry(plan cryptodb.Plan, pair cryptodb.Pair, marketStopLoss *cryptodb.Order, limitStopLoss cryptodb.Order, entry *cryptodb.Order) (err error) {

	log.Printf("Placing entry")
	var result exchange.OrderResponseRest
	entryParams := make(exchange.RequestParameters)

	entryParams["order_link_id"] = entry.LinkOrderID
	entryParams["symbol"] = pair.Name
	if plan.Direction == cryptodb.Long {
		entryParams["side"] = "Buy"
	} else {
		entryParams["side"] = "Sell"
	}
	entryParams["order_type"] = "Limit"
	entryParams["qty"] = entry.Size.InexactFloat64()
	entryParams["price"] = entry.Price.InexactFloat64()
	entryParams["close_on_trigger"] = false
	entryParams["reduce_only"] = false
	entryParams["time_in_force"] = "GoodTillCancel"
	entryParams["stop_loss"] = marketStopLoss.Price.InexactFloat64()
	entryParams["sl_trigger_by"] = "LastPrice"

	var response exchange.OrderResponseRest
	_, responseBuffer, err := e.SignedRequest(http.MethodPost, "/private/linear/order/create", entryParams, &result)

	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		log.Printf("Order not accepted. ReturnCode: %d, ReturnMessage: %s", result.ReturnCode, result.ReturnMessage)
		return errors.New(result.ReturnMessage)
	}

	// log.Printf("Response buffer: %s", string(responseBuffer))
	err = json.Unmarshal(responseBuffer, &response)
	if err != nil {
		log.Printf("error unmarshalling responseBuffer %s", err)
		return err
	}

	log.Printf("Response to setting Entry: %v", response)
	log.Printf("Setting SystemOrderID to: %s", response.Result.OrderID)
	entry.SystemOrderID = response.Result.OrderID
	entry.Status.Scan(response.Result.OrderStatus)
	log.Printf("Placing entry succesfull")

	return nil
}

func placeLimitStopLoss(plan cryptodb.Plan, pair cryptodb.Pair, ticker exchange.Ticker, marketStopLoss cryptodb.Order, limitStopLoss *cryptodb.Order, entry cryptodb.Order) (err error) {
	log.Printf("Placing limit stoploss...")
	var result exchange.OrderResponseRest

	lslParams := make(exchange.RequestParameters)
	lslParams["order_link_id"] = limitStopLoss.LinkOrderID
	lslParams["symbol"] = pair.Name
	lslParams["order_type"] = "Limit"
	lslParams["qty"] = limitStopLoss.Size.InexactFloat64()
	if plan.Direction == cryptodb.Long {
		lslParams["side"] = "Sell"
	} else {
		lslParams["side"] = "Buy"
	}

	lslParams["trigger_by"] = "LastPrice"
	lslParams["price"] = limitStopLoss.Price.InexactFloat64()
	lslParams["stop_px"] = limitStopLoss.Price.InexactFloat64()
	lslParams["base_price"] = ticker.LastPrice.InexactFloat64()

	lslParams["close_on_trigger"] = false
	lslParams["reduce_only"] = true
	lslParams["time_in_force"] = "GoodTillCancel"

	var response exchange.OrderResponseRest
	_, responseBuffer, err := e.SignedRequest(http.MethodPost, "/private/linear/stop-order/create", lslParams, &result)
	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		log.Printf("Order not accepted. ReturnCode: %d, ReturnMessage: %s", result.ReturnCode, result.ReturnMessage)
		return errors.New(result.ReturnMessage)
	}

	log.Printf("Response: %s", string(responseBuffer))
	err = json.Unmarshal(responseBuffer, &response)
	if err != nil {
		log.Printf("error unmarshalling responseBuffer %s", err)
		return err
	}

	limitStopLoss.SystemOrderID = response.Result.StopOrderID
	limitStopLoss.Status.Scan(response.Result.OrderStatus)

	log.Printf("Placing limit stoploss succesfull")
	return nil
}

func placeTakeProfit(plan cryptodb.Plan, pair cryptodb.Pair, ticker exchange.Ticker, marketStopLoss, entry cryptodb.Order, takeProfit *cryptodb.Order) (err error) {
	log.Printf("Placing take profit...")
	var result exchange.OrderResponseRest
	tpParams := make(exchange.RequestParameters)

	tpParams["order_link_id"] = takeProfit.LinkOrderID
	tpParams["symbol"] = pair.Name
	tpParams["order_type"] = "Limit"
	tpParams["qty"] = takeProfit.Size.InexactFloat64()
	if plan.Direction == cryptodb.Long {
		tpParams["side"] = "Sell"
	} else {
		tpParams["side"] = "Buy"
	}
	tpParams["trigger_by"] = "LastPrice"
	tpParams["price"] = takeProfit.Price.InexactFloat64()
	// TODO: also implement short situation
	priceDifference := takeProfit.Price.Sub(entry.Price)
	triggerPrice := entry.Price.Add(priceDifference.Mul(decimal.NewFromFloat(0.95)))
	log.Printf("Calculatied triggerPrice as: %s", triggerPrice.String())
	tpParams["stop_px"] = triggerPrice.InexactFloat64()
	tpParams["base_price"] = entry.Price.InexactFloat64()

	tpParams["close_on_trigger"] = false
	tpParams["reduce_only"] = true
	tpParams["time_in_force"] = "GoodTillCancel"

	var response exchange.OrderResponseRest
	_, responseBuffer, err := e.SignedRequest(http.MethodPost, "/private/linear/stop-order/create", tpParams, &result)
	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		log.Printf("Order not accepted. ReturnCode: %d, ReturnMessage: %s", result.ReturnCode, result.ReturnMessage)
		return errors.New(result.ReturnMessage)
	}

	log.Printf("Response: %s", string(responseBuffer))
	err = json.Unmarshal(responseBuffer, &response)
	if err != nil {
		log.Printf("error unmarshalling responseBuffer %s", err)
		return err
	}

	takeProfit.SystemOrderID = response.Result.StopOrderID

	log.Printf("Placing takeprofit succesfull.")
	return nil
}

func matchExchangeOrder(SystemOrderID string) (o cryptodb.Order, err error) {
	result := db.Where("system_order_id = ?", SystemOrderID).Last(&o)

	if result.RowsAffected > 1 {
		log.Printf("Multiple results found??? WTH.")
	}

	return o, result.Error
}

func processEntryOrder(entryOrder cryptodb.Order, o exchange.Order) (err error) {
	log.Printf("Processing Entry w/ MarketStopLoss started")

	var plan cryptodb.Plan

	var marketStopLossOrder cryptodb.Order
	result := db.Where("order_kind = ?", cryptodb.MarketStopLoss).Where("plan_id = ?", entryOrder.PlanID).First(&marketStopLossOrder)
	if result.Error != nil {
		log.Printf("Entry w. MarketStoploss order for MarketStopLoss not found")
		return result.Error
	}

	result = db.Where("id = ?", entryOrder.PlanID).First(&plan)
	if result.Error != nil {
		log.Printf("Entry Plan not found")
		return result.Error
	}

	entryOrder.Status.Scan(o.OrderStatus)
	// TODO: also handle change price if user moves price from exchange website or app, although that wouldn't be the prefered way
	plan.Status = entryOrder.Status

	tx := db.Begin()
	result = tx.Save(plan)
	if result.Error != nil {
		log.Printf("This is weird AF.")
		tx.Rollback()
		return result.Error
	}

	var logEntry cryptodb.Log
	logEntry.PlanID = entryOrder.PlanID
	logEntry.Source = cryptodb.Exchange

	switch o.OrderStatus {
	case "New":
		var stopLossSetMsg string
		if marketStopLossOrder.Price.Equal(o.StopLoss) {
			marketStopLossOrder.Status = entryOrder.Status
			stopLossSetMsg = "and"
		} else {
			stopLossSetMsg = "but NOT" // TODO: this should NEVER happen.
		}
		marketStopLossOrder.Status = entryOrder.Status

		logEntry.Text = fmt.Sprintf("Exchange processed Entry Order %d %s stoploss, and set status to %s.", entryOrder.ID, stopLossSetMsg, entryOrder.Status.String())
	case "Filled":
		logEntry.Text = fmt.Sprintf("Entry completely filled at %s. Changing Plan Status to Filled.", o.CreatedAt.Format("2006-01-02 15:04:05.000"))
	case "Cancelled":
		logEntry.Text = fmt.Sprintf("Entry and market stoploss cancelled. Changing Plan Status to cancelled.")
	case "PartiallyFilled":
		logEntry.Text = fmt.Sprintf("Entry partially filled: %s/%s", o.Leaves.String(), entryOrder.Size.String())
	default:
		log.Printf("Handling of OrderStatus: %s not implemented yet.", o.OrderStatus)
		return errors.New("Unhandled OrderStates")
	}

	result = tx.Save(&entryOrder)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Saving Entry errored %s", result.Error)
		return result.Error
	}

	result = tx.Save(&marketStopLossOrder)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Processing Entry w/ MarketStopLoss errored")
		return result.Error
	}

	result = tx.Create(&logEntry)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("logging for entry order errored %s", result.Error)
		return result.Error
	}

	result = tx.Commit()

	log.Printf("Processing Entry finished")
	return result.Error
}

func processMarketStoploss(marketStopLossOrder cryptodb.Order, o exchange.Order) (err error) {
	log.Printf("Processing MarketStopLoss started")

	var logEntry cryptodb.Log
	logEntry.PlanID = marketStopLossOrder.PlanID
	logEntry.Source = cryptodb.Exchange

	switch o.OrderStatus {
	case "New":
		var stopLossSetMsg string
		if marketStopLossOrder.Price.Equal(o.StopLoss) {
			marketStopLossOrder.Status.Scan(o.OrderStatus)
			stopLossSetMsg = "and"
		} else {
			stopLossSetMsg = "but NOT" // TODO: this should NEVER happen.
		}

		logEntry.Text = fmt.Sprintf("Exchange processed Market StopLoss order %d Market stoploss, %s and set status to %s.", marketStopLossOrder.ID, stopLossSetMsg, marketStopLossOrder.Status.String())
	case "Filled":
		logEntry.Text = fmt.Sprintf("Entry completely filled at %s. Changing Plan Status to Filled.", o.CreatedAt.Format("2006-01-02 15:04:05.000"))
	case "Cancelled":
		logEntry.Text = fmt.Sprintf("Entry and market stoploss cancelled. Changing Plan Status to cancelled.")
	case "PartiallyFilled":
		logEntry.Text = fmt.Sprintf("Entry partially filled: /%s", marketStopLossOrder.Price.String())
	default:
		log.Printf("Handling of OrderStatus: %s not implemented yet.", o.OrderStatus)
		return errors.New("Unhandled OrderStatus")
	}

	tx := db.Begin()
	result := tx.Save(&marketStopLossOrder)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Saving MarketStopLoss errored")
		return result.Error
	}

	result = tx.Create(&logEntry)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Saving log MarketStopLoss errored")
		return result.Error
	}

	result = tx.Commit()

	log.Printf("Processing MarketStopLoss finished")
	return result.Error
}

func processTakeProfit(takeProfit cryptodb.Order, o exchange.Order) (err error) {
	log.Printf("Processing takeProfit order started")

	tx := db.Begin()

	var logEntry cryptodb.Log
	takeProfit.Status.Scan(o.OrderStatus)
	logEntry.PlanID = takeProfit.PlanID
	logEntry.Source = cryptodb.Exchange
	logEntry.Text = fmt.Sprintf("Exchange processed %s order %d, and set status to %s.", takeProfit.OrderKind.String(), takeProfit.ID, o.OrderStatus)

	result := tx.Save(&takeProfit)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Processing takeProfit errored on saving: %s", result.Error)
		return result.Error
	}

	result = tx.Create(&logEntry)
	if result.Error != nil {
		tx.Rollback()
		log.Printf("Processing takeProfit errored on logCreation: %s", result.Error)
		return result.Error
	}

	result = tx.Commit()

	log.Printf("Processing takeProfit completed")
	return result.Error
}

func processOrder(o exchange.Order) error {

	var matchOrder string

	if o.OrderID == "" && o.StopOrderID == "" {
		log.Printf("Unknown order_id AND stop_order_id. Impossible to find order?")
	}
	if o.OrderID != "" {
		matchOrder = o.OrderID
	}
	if o.StopOrderID != "" {
		matchOrder = o.StopOrderID
	}
	order, err := matchExchangeOrder(matchOrder)
	if err != nil {
		log.Printf("Couldn't find matching order for %s. Weird. Aborting.", matchOrder)
		return err
	}
	switch order.OrderKind {
	case cryptodb.Entry:
		err := processEntryOrder(order, o)
		if err != nil {
			log.Printf("Error processing entryOrder: %s", err)
		}
	case cryptodb.MarketStopLoss:
		err := processMarketStoploss(order, o)
		if err != nil {
			log.Printf("Error processing marketStopLossOrder: %s", err)
		}
	case cryptodb.TakeProfit, cryptodb.LimitStopLoss:
		err := processTakeProfit(order, o)
		if err != nil {
			log.Printf("Error processing limitStopLossOrder: %s", err)
		}
	}
	return nil
}

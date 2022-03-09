package main

import (
	"errors"
	"fmt"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

func placeOrders(p cryptodb.Plan, pair cryptodb.Pair, o []cryptodb.Order) (err error) {
	switch p.Direction {
	case cryptodb.Long:
		if p.Leverage.GreaterThan(pair.Leverage.Long) {
			setLeverage(p, &pair)
		}
	case cryptodb.Short:
		if p.Leverage.GreaterThan(pair.Leverage.Short) {
			setLeverage(p, &pair)
		}
	}

	err = setEntry(p, pair, &o[cryptodb.MarketStopLoss], &o[cryptodb.Entry])
	if err != nil {
		return err
	}

	// err = setLimitStopLoss(&p, pair, o[cryptodb.MarketStopLoss], &o[cryptodb.LimitStopLoss], o[cryptodb.Entry])
	// if err != nil {
	// 	return err
	// }

	for i := 3; i < 3+cryptodb.MaxTakeProfits; i++ {
		if !o[i].Price.IsZero() {
			err = setTakeProfit(p, pair, o[cryptodb.MarketStopLoss], o[cryptodb.Entry], &o[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func setLeverage(plan cryptodb.Plan, pair *cryptodb.Pair) (err error) {
	switch plan.Direction {
	case cryptodb.Long:
		err = e.SendLeverage(pair.Name, pair.QuoteCurrency, plan.Leverage, pair.Leverage.Short)
	case cryptodb.Short:
		err = e.SendLeverage(pair.Name, pair.QuoteCurrency, pair.Leverage.Long, plan.Leverage)
	}
	if err != nil {
		logEntry := &cryptodb.Log{
			PlanID: plan.ID,
			Source: cryptodb.Server,
			Text:   fmt.Sprintf("Error sending leverage: %s", err),
		}
		db.Create(logEntry)
		return err
	}

	switch plan.Direction {
	case cryptodb.Long:
		pair.Leverage.Long = plan.Leverage
	case cryptodb.Short:
		pair.Leverage.Short = plan.Leverage
	}

	logEntry := &cryptodb.Log{
		PlanID: plan.ID,
		Source: cryptodb.Server,
		Text:   fmt.Sprintf("Sending set %s leverage to %s successful.", plan.Direction.String(), plan.Leverage.String()),
	}

	db.Save(pair)
	db.Create(logEntry)
	return nil
}

func setEntry(p cryptodb.Plan, pair cryptodb.Pair, marketStopLoss *cryptodb.Order, entry *cryptodb.Order) (err error) {
	err = e.SendEntry(p, pair, marketStopLoss, entry)
	if err != nil {
		logEntry := &cryptodb.Log{
			PlanID: p.ID,
			Source: cryptodb.Server,
			Text:   fmt.Sprintf("Error sending entry and market stoploss: %s", err),
		}
		result := db.Create(logEntry)
		if result.Error != nil {
			return result.Error
		}
		return err
	}

	logEntry := &cryptodb.Log{
		PlanID: p.ID,
		Source: cryptodb.Server,
		Text: fmt.Sprintf("Sending set entry (%s %s@%s) and market stoploss (@%s) successful",
			p.Direction.String(), entry.Size.String(), entry.Price.String(), marketStopLoss.Price.String()),
	}

	db.Save(entry)
	db.Create(logEntry)

	return nil
}

func setLimitStopLoss(p *cryptodb.Plan, pair cryptodb.Pair, marketStopLoss cryptodb.Order, limitStopLoss *cryptodb.Order, entry cryptodb.Order) (err error) {
	err = e.SendLimitStopLoss(p, pair, marketStopLoss, limitStopLoss, entry)
	if err != nil {
		logEntry := &cryptodb.Log{
			PlanID: p.ID,
			Source: cryptodb.Server,
			Text:   fmt.Sprintf("Error sending Limit StopLoss (%s@%s %s): %s", limitStopLoss.Size.String(), limitStopLoss.Price.String(), limitStopLoss.TriggerPrice.String(), err),
		}
		result := db.Create(logEntry)
		if result.Error != nil {
			return result.Error
		}
		p.Status = cryptodb.Error
		db.Save(p)
		if result.Error != nil {
			return result.Error
		}
		return err
	}

	logEntry := &cryptodb.Log{
		PlanID: p.ID,
		Source: cryptodb.Server,
		Text:   fmt.Sprintf("Sending set Limit StopLoss (%s@%s %s) successful.", limitStopLoss.Size.String(), limitStopLoss.Price.String(), limitStopLoss.TriggerPrice.String()),
	}

	db.Save(limitStopLoss)
	db.Create(logEntry)

	return nil
}

func setTakeProfit(p cryptodb.Plan, pair cryptodb.Pair, marketStopLoss, entry cryptodb.Order, takeProfit *cryptodb.Order) (err error) {
	err = e.SendTakeProfit(p, pair, marketStopLoss, entry, takeProfit)
	if err != nil {
		logEntry := &cryptodb.Log{
			PlanID: p.ID,
			Source: cryptodb.Server,
			Text:   fmt.Sprintf("Error sending take profit (%s@%s %s): %s", takeProfit.Size.String(), takeProfit.Price.String(), takeProfit.TriggerPrice.String(), err),
		}
		result := db.Create(logEntry)
		if result.Error != nil {
			return result.Error
		}
		p.Status = cryptodb.Error
		db.Save(p)
		if result.Error != nil {
			return result.Error
		}
		return err
	}

	logEntry := &cryptodb.Log{
		PlanID: p.ID,
		Source: cryptodb.Server,
		Text:   fmt.Sprintf("Sending set Take Profit (%s@%s %s) successful.", takeProfit.Size.String(), takeProfit.Price.String(), takeProfit.TriggerPrice.String()),
	}

	db.Save(takeProfit)
	db.Create(logEntry)

	return nil
}

func processOrder(incomingOrder exchange.Order) error {
	var matchOrder string
	var marketStopLossOrder cryptodb.Order
	var plan cryptodb.Plan

	if incomingOrder.OrderType == "Market" {
		result := db.
			// TODO: make more sure that it's the correct order, also check Pair and Price?
			// and first check if the SystemOrderID is known.
			Where("system_order_id = ? AND order_kind = ?", incomingOrder.StopOrderID, cryptodb.MarketStopLoss).
			First(&marketStopLossOrder)
		if result.Error != nil {
			result := db.
				Joins("JOIN plans ON orders.plan_id = plans.id").
				Joins("JOIN pairs ON pairs.id = plans.pair_id").
				Where("order_kind = ? AND pairs.name = ?", cryptodb.MarketStopLoss, incomingOrder.Symbol).
				First(&marketStopLossOrder)
			if result.Error != nil {
				return result.Error
			} else {
				if incomingOrder.StopOrderID != "" {
					marketStopLossOrder.SystemOrderID = incomingOrder.StopOrderID
					db.Save(marketStopLossOrder)
					db.Create(&cryptodb.Log{
						PlanID: marketStopLossOrder.PlanID,
						Source: cryptodb.Server,
						Text:   fmt.Sprintf("Set SystemOrderID (%s) for Market StopLoss Order.", incomingOrder.StopOrderID),
					})
				}
				result = db.Where("id = ?", marketStopLossOrder.PlanID).First(&plan)
				processMarketStoploss(plan, marketStopLossOrder, incomingOrder)
				return nil
			}
		} else {
			result = db.Where("id = ?", marketStopLossOrder.PlanID).First(&plan)
			processMarketStoploss(plan, marketStopLossOrder, incomingOrder)
			return nil
		}
	} else {
		if incomingOrder.OrderID != "" {
			matchOrder = incomingOrder.OrderID
		} else if incomingOrder.StopOrderID != "" {
			matchOrder = incomingOrder.StopOrderID
		} else {
			return errors.New("both order_id and stop_order_id empty")
		}
	}

	var order cryptodb.Order
	var pair cryptodb.Pair
	var entryOrder cryptodb.Order

	result := db.Where("system_order_id = ?", matchOrder).First(&order)
	if result.Error != nil {
		return result.Error
	}

	result = db.Where("id = ?", order.PlanID).First(&plan)
	if result.Error != nil {
		return result.Error
	}

	result = db.Where("id = ?", plan.PairID).First(&pair)
	if result.Error != nil {
		return result.Error
	}

	result = db.Where("order_kind = ?", cryptodb.Entry).Where("plan_id = ?", order.PlanID).First(&entryOrder)
	if result.Error != nil {
		return result.Error
	}

	result = db.Where("order_kind = ?", cryptodb.MarketStopLoss).Where("plan_id = ?", order.PlanID).First(&marketStopLossOrder)
	if result.Error != nil {
		return result.Error
	}
	switch order.OrderKind {
	case cryptodb.Entry:
		err := processEntryOrder(plan, pair, marketStopLossOrder, order, incomingOrder)
		if err != nil {
			return err
		}

	case cryptodb.TakeProfit:
		err := processTakeProfit(plan, pair, entryOrder, order, incomingOrder)
		if err != nil {
			return err
		}
	// case cryptodb.LimitStopLoss:
	// 	err := processLimitStoploss(plan, order, incomingOrder)
	// 	if err != nil {
	// 		return err
	// 	}
	}

	return nil
}

func processEntryOrder(plan cryptodb.Plan, pair cryptodb.Pair, marketStopLossOrder, entryOrder cryptodb.Order, o exchange.Order) (err error) {
	// TODO: also handle change price if user moves price from exchange website or app, although that wouldn't be the prefered way

	var exchangeLogEntry cryptodb.Log
	var planUpdateLogEntry cryptodb.Log
	exchangeLogEntry.PlanID = entryOrder.PlanID
	exchangeLogEntry.Source = cryptodb.Exchange
	planUpdateLogEntry.PlanID = entryOrder.PlanID
	planUpdateLogEntry.Source = cryptodb.Server

	switch o.OrderStatus {
	case "New":
		entryOrder.Status.Scan(o.OrderStatus)
		db.Save(&entryOrder)

		var stopLossSetMsg string
		if marketStopLossOrder.Price.Equal(o.StopLoss) {
			marketStopLossOrder.Status = entryOrder.Status
			stopLossSetMsg = "and"
		} else {
			stopLossSetMsg = "but NOT" // TODO: this should NEVER happen.
		}
		exchangeLogEntry.Text = fmt.Sprintf("Processed Entry Order %s  Market StopLoss, and set status to %s.", stopLossSetMsg, entryOrder.Status.String())
		db.Create(&exchangeLogEntry)

		plan.Status = entryOrder.Status
		db.Save(&plan)
		planUpdateLogEntry.Text = fmt.Sprintf("Changing plan status to %s.", entryOrder.Status.String())
		db.Create(&planUpdateLogEntry)

	case "Filled", "Cancelled", "PartiallyFilled":
		if entryOrder.Status.String() != o.OrderStatus {
			exchangeLogEntry.Text = fmt.Sprintf("Entry %s.", o.OrderStatus)
			entryOrder.Status.Scan(o.OrderStatus)
			db.Save(&entryOrder)
			db.Create(&exchangeLogEntry)

			planUpdateLogEntry.Text = fmt.Sprintf("Changing plan status to %s.", o.OrderStatus)
			plan.Status = entryOrder.Status
			db.Save(&plan)
			db.Create(&planUpdateLogEntry)
		}

	default:
		exchangeLogEntry.Text = fmt.Sprintf("Unhandled OrderStatus (%s) for Entry", o.OrderStatus)
		db.Create(&exchangeLogEntry)
		return errors.New(fmt.Sprintf("Unhandled OrderStatus (%s) for entry", o.OrderStatus))
	}

	return nil
}

func processMarketStoploss(plan cryptodb.Plan, marketStopLossOrder cryptodb.Order, o exchange.Order) (err error) {

	var exchangeLogEntry cryptodb.Log
	var planUpdateLogEntry cryptodb.Log
	exchangeLogEntry.PlanID = marketStopLossOrder.PlanID
	exchangeLogEntry.Source = cryptodb.Exchange
	planUpdateLogEntry.PlanID = marketStopLossOrder.PlanID
	planUpdateLogEntry.Source = cryptodb.Server

	switch o.OrderStatus {
	case "New", "Untriggered", "PartiallyFilled":
		if marketStopLossOrder.Status.String() != o.OrderStatus {
			marketStopLossOrder.Status.Scan(o.OrderStatus)
			db.Save(&marketStopLossOrder)
			exchangeLogEntry.Text = fmt.Sprintf("Processed Market StopLoss set status to %s.", marketStopLossOrder.Status.String())
			db.Create(&exchangeLogEntry)
		}

	case "Cancelled", "Filled":
		marketStopLossOrder.Status.Scan(o.OrderStatus)
		exchangeLogEntry.Text = fmt.Sprintf("Market stoploss %s", marketStopLossOrder.Status.String())
		db.Create(&exchangeLogEntry)
		marketStopLossOrder.Status.Scan(o.OrderStatus)
		db.Save(&marketStopLossOrder)
		planUpdateLogEntry.Text = fmt.Sprintf("Changing plan status to %s", marketStopLossOrder.Status.String())
		db.Create(&planUpdateLogEntry)
		plan.Status = marketStopLossOrder.Status
		db.Save(&plan)

	default:
		exchangeLogEntry.Text = fmt.Sprintf("Unhandled OrderStatus (%s) for Market StopLoss", o.OrderStatus)
		db.Create(&exchangeLogEntry)
		return errors.New(fmt.Sprintf("Unhandled OrderStatus (%s) for marketstoploss", o.OrderStatus))
	}

	return nil
}

func processLimitStoploss(plan cryptodb.Plan, limitStopLossOrder cryptodb.Order, o exchange.Order) (err error) {
	var exchangeLogEntry cryptodb.Log
	var planUpdateLogEntry cryptodb.Log
	exchangeLogEntry.PlanID = limitStopLossOrder.PlanID
	exchangeLogEntry.Source = cryptodb.Exchange
	planUpdateLogEntry.PlanID = limitStopLossOrder.PlanID
	planUpdateLogEntry.Source = cryptodb.Server

	switch o.OrderStatus {
	case "New", "Untriggered", "PartiallyFilled", "Triggered", "Deactivated":
		if limitStopLossOrder.Status.String() != o.OrderStatus {
			limitStopLossOrder.Status.Scan(o.OrderStatus)
			db.Save(&limitStopLossOrder)
			exchangeLogEntry.Text = fmt.Sprintf("Processed Limit StopLoss set status to %s.", limitStopLossOrder.Status.String())
			db.Create(&exchangeLogEntry)
		}

	case "Cancelled", "Filled", "Rejected":
		limitStopLossOrder.Status.Scan(o.OrderStatus)
		db.Save(&limitStopLossOrder)
		exchangeLogEntry.Text = fmt.Sprintf("Limit StopLoss status set to %s.", limitStopLossOrder.Status.String())
		db.Create(&exchangeLogEntry)
		planUpdateLogEntry.Text = fmt.Sprintf("Changing plan status to %s", limitStopLossOrder.Status.String())
		db.Create(&planUpdateLogEntry)
		plan.Status = limitStopLossOrder.Status
		db.Save(&plan)

		// TODO: cancel remaining open (takeProfit) orders
		// TODO: mark it as a win or loss. calculate performance?
		// Lower leverage

	default:
		exchangeLogEntry.Text = fmt.Sprintf("Unhandled OrderStatus (%s) for Limit StopLoss", o.OrderStatus)
		db.Create(&exchangeLogEntry)
		return errors.New(fmt.Sprintf("Unhandled OrderStatus (%s) for Limit StopLoss", o.OrderStatus))
	}

	return nil
}

func processTakeProfit(plan cryptodb.Plan, pair cryptodb.Pair, entryOrder, takeProfit cryptodb.Order, o exchange.Order) (err error) {
	var exchangeLogEntry cryptodb.Log
	var planUpdateLogEntry cryptodb.Log
	exchangeLogEntry.PlanID = takeProfit.PlanID
	exchangeLogEntry.Source = cryptodb.Exchange

	switch o.OrderStatus {
	case "New", "Untriggered", "PartiallyFilled", "Filled":
		if takeProfit.Status.String() != o.OrderStatus {
			takeProfit.Status.Scan(o.OrderStatus)
			db.Save(&takeProfit)
			exchangeLogEntry.Text = fmt.Sprintf("Take Profit status set to %s.", takeProfit.Status.String())
			db.Create(&exchangeLogEntry)
		}

	case "Rejected", "Deactivated":
		takeProfit.Status.Scan(o.OrderStatus)
		db.Save(&takeProfit)
		exchangeLogEntry.Text = fmt.Sprintf("Take Profit status set to %s.", takeProfit.Status.String())
		db.Create(&exchangeLogEntry)
		planUpdateLogEntry.Text = fmt.Sprintf("Changing plan status to %s", takeProfit.Status.String())
		db.Create(&planUpdateLogEntry)
		plan.Status = takeProfit.Status
		db.Save(&plan)

	default:
		exchangeLogEntry.Text = fmt.Sprintf("Unhandled OrderStatus (%s) for takeProfit", o.OrderStatus)
		db.Create(&exchangeLogEntry)
		return errors.New(fmt.Sprintf("Unhandled OrderStatus (%s) for takeProfit", o.OrderStatus))
	}

	return nil
}

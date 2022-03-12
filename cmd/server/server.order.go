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
	db.Save(marketStopLoss)
	db.Create(logEntry)

	return nil
}

func setLimitStopLoss(p cryptodb.Plan, pair cryptodb.Pair, marketStopLoss cryptodb.Order, limitStopLoss *cryptodb.Order, entry cryptodb.Order) (err error) {
	err = e.SendLimitOrder(p, pair, entry, limitStopLoss)
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
	err = e.SendLimitOrder(p, pair, entry, takeProfit)
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
	var order cryptodb.Order
	var plan cryptodb.Plan

	if incomingOrder.OrderID != "" {
		matchOrder = incomingOrder.OrderID
	} else if incomingOrder.StopOrderID != "" {
		matchOrder = incomingOrder.StopOrderID
	} else {
		return errors.New("both order_id and stop_order_id empty")
	}

	result := db.Where("system_order_id = ?", matchOrder).First(&order)
	if result.Error != nil && incomingOrder.OrderType == "Market" {
		// TODO: make sure the plan is still open and in the same direction?
		result := db.Joins("JOIN plans ON orders.plan_id = plans.id").
			Joins("JOIN pairs ON pairs.id = plans.pair_id").
			Where("order_kind = ? AND pairs.name = ?", cryptodb.MarketStopLoss, incomingOrder.Symbol).
			First(&order)
		if result.Error != nil {
			return result.Error
		} else {
			order.SystemOrderID = incomingOrder.StopOrderID
		}
	} else {
		return result.Error
	}

	result = db.Where("id = ?", order.PlanID).First(&plan)
	if result.Error != nil {
		return result.Error
	}

	updateStatus(plan, order, incomingOrder)

	return nil
}

func updateStatus(plan cryptodb.Plan, dbOrder cryptodb.Order, exchangeOrder exchange.Order) {
	var newStatus cryptodb.Status
	newStatus.Scan(exchangeOrder.OrderStatus)

	updateOrderStatus(&plan, &dbOrder, newStatus)
	updatePlanStatus(&plan, dbOrder.Status)
}

func updatePlanStatus(plan *cryptodb.Plan, newStatus cryptodb.Status) {
	if plan.Status >= newStatus {
		return
	}

	db.Create(&cryptodb.Log{
		PlanID: plan.ID,
		Source: cryptodb.Server,
		Text:   fmt.Sprintf("Plan status updated to %s.", newStatus),
	})
	plan.Status = newStatus
	db.Save(plan)
}

func updateOrderStatus(plan *cryptodb.Plan, order *cryptodb.Order, newStatus cryptodb.Status) {
	if order.Status >= newStatus {
		return
	}

	db.Create(&cryptodb.Log{
		PlanID: plan.ID,
		Source: cryptodb.Exchange,
		Text:   fmt.Sprintf("%s@%s status updated to %s.", order.OrderKind.String(), order.Price.String(), newStatus),
	})
	order.Status = newStatus
	db.Save(order)
}

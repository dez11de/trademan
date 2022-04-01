package main

import (
	"errors"
	"fmt"
	"log"

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
	err = e.SendOrder(p, pair, entry, marketStopLoss)
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

	db.Save(entry)
	db.Save(marketStopLoss)
	db.Create(&cryptodb.Log{PlanID: p.ID,
		Source: cryptodb.Server,
		Text: fmt.Sprintf("Sending set entry (%s %s@%s) and market stoploss (@%s) successful.",
			p.Direction.String(), entry.Size.String(), entry.Price.String(), marketStopLoss.Price.String()),
	})

	p.Status = entry.Status
	db.Save(p)
	db.Create(&cryptodb.Log{PlanID: p.ID,
		Source: cryptodb.Server,
		Text:   fmt.Sprintf("Setting plan status to %s", p.Status.String()),
	})

	return nil
}

func setTakeProfit(p cryptodb.Plan, pair cryptodb.Pair, entry, takeProfit *cryptodb.Order) (err error) {
	err = e.SendOrder(p, pair, entry, takeProfit)
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

	result := db.Save(takeProfit)
	if result.Error != nil {
		return result.Error
	}

	logEntry := &cryptodb.Log{
		PlanID: p.ID,
		Source: cryptodb.Server,
		Text:   fmt.Sprintf("Sending set Take Profit (%s@%s %s) successful.", takeProfit.Size.String(), takeProfit.Price.String(), takeProfit.TriggerPrice.String()),
	}
	result = db.Create(logEntry)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func matchOrder(incomingOrder *exchange.Order) (o cryptodb.Order, err error) {
	log.Printf("Matching order: %+v", incomingOrder)
	var matchOrderID string
	var order cryptodb.Order

	if incomingOrder.OrderID != "" {
		matchOrderID = incomingOrder.OrderID
	} else if incomingOrder.StopOrderID != "" {
		matchOrderID = incomingOrder.StopOrderID
	} else {
		return o, errors.New("both order_id and stop_order_id empty")
	}

	log.Printf("Finding order: %+v", matchOrderID)

	result := db.Where("system_order_id = ?", matchOrderID).Take(&order)
	if result.Error != nil && incomingOrder.OrderType == "Market" {
		// TODO: make sure the plan is still open and in the same direction?
		result := db.Joins("JOIN plans ON orders.plan_id = plans.id").
			Joins("JOIN pairs ON pairs.id = plans.pair_id").
			Where("order_kind = ? AND pairs.name = ?", cryptodb.MarketStopLoss, incomingOrder.Symbol).
			Take(&order)
		if result.Error != nil {
			log.Printf("Error finding market stoploss order: %s", result.Error)
			return o, result.Error
		} else {
			// First sign of stoploss order
			log.Printf("First sign of marketstoploss order: %+v", order)
			order.SystemOrderID = incomingOrder.StopOrderID
		}
	} else {
		log.Printf("Found order: %+v", order)
		return order, result.Error
	}

	return order, nil
}

func processOrder(incomingOrder exchange.Order) error {
	var plan cryptodb.Plan
	order, err := matchOrder(&incomingOrder)
	if err != nil {
		log.Printf("Error matching order: %s", err)
		return err
	}
	log.Printf("Finding plan for order: %+v", order)
	result := db.Debug().Where("id = ?", order.PlanID).Take(&plan)
	if result.Error != nil {
		log.Printf("Error finding plan: %s", result.Error)
		return result.Error
	}

	updateStatus(&plan, &order, incomingOrder)

	switch {
	case order.OrderKind == cryptodb.Entry && order.Status == cryptodb.Filled:
		sendTakeProfits(plan)

		// TODO: or LimitStopLoss
	case order.OrderKind == cryptodb.MarketStopLoss && order.Status == cryptodb.Filled:
		cancelTakeProfits(plan)
	}

	return nil
}

func updateStatus(plan *cryptodb.Plan, dbOrder *cryptodb.Order, exchangeOrder exchange.Order) {
	var newStatus cryptodb.Status
	newStatus.Scan(exchangeOrder.OrderStatus)

	updateOrderStatus(plan, dbOrder, newStatus)
	if dbOrder.OrderKind != cryptodb.MarketStopLoss && newStatus != cryptodb.Untriggered {
		updatePlanStatus(plan, newStatus)
	}
}

func updatePlanStatus(plan *cryptodb.Plan, newStatus cryptodb.Status) {
	log.Printf("Updating plan status: %s(%d) -> %s(%d)", plan.Status.String(), plan.Status, newStatus.String(), newStatus)
	if plan.Status >= newStatus {
		return
	}

	db.Create(&cryptodb.Log{
		PlanID: plan.ID,
		Source: cryptodb.Server,
		Text:   fmt.Sprintf("Plan status updated to %s.", newStatus),
	})
	plan.Status = newStatus
	db.Save(&plan)
}

func updateOrderStatus(plan *cryptodb.Plan, order *cryptodb.Order, newStatus cryptodb.Status) {
	log.Printf("Updating order status: %s(%d) -> %s(%d)", order.Status.String(), order.Status, newStatus.String(), newStatus)
	if order.Status >= newStatus {
		return
	}

	db.Create(&cryptodb.Log{
		PlanID: plan.ID,
		Source: cryptodb.Exchange,
		Text:   fmt.Sprintf("%s@%s status updated to %s.", order.OrderKind.String(), order.Price.String(), newStatus),
	})
	order.Status = newStatus
	db.Save(&order)
}

func sendTakeProfits(p cryptodb.Plan) (err error) {
	db.Create(&cryptodb.Log{
		PlanID: p.ID,
		Source: cryptodb.Server,
		Text:   "Entry Filled. Sending Take Profits...",
	})

	var takeProfits []cryptodb.Order
	result := db.Where("plan_id = ? AND order_kind = ?", p.ID, cryptodb.TakeProfit).Find(&takeProfits)
	if result.Error != nil {
		log.Print("No matching TakeProfits found....")
		return result.Error
	}
	var entry cryptodb.Order
	result = db.Where("plan_id = ? AND order_kind = ?", p.ID, cryptodb.Entry).Take(&entry)
	if result.Error != nil {
		log.Print("No matching TakeProfits found....")
		return result.Error
	}
	var pair cryptodb.Pair
	result = db.Where("id = ?", p.PairID).Find(&pair)
	if result.Error != nil {
		return result.Error
	}

	for _, o := range takeProfits {
		if !o.Price.IsZero() {
			log.Printf("setting take profit - entry price: %s, trigger price: %s, order price: %s", entry.Price.String(), o.TriggerPrice.String(), o.Price.String())
			err := setTakeProfit(p, pair, &entry, &o)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func cancelTakeProfits(p cryptodb.Plan) (err error) {
	db.Create(&cryptodb.Log{
		PlanID: p.ID,
		Source: cryptodb.Server,
		Text:   "StopLoss. Canceling Take Profits...",
	})

	var takeProfits []cryptodb.Order
	result := db.Where("plan_id = ? AND order_kind = ?", p.ID, cryptodb.TakeProfit).Find(&takeProfits)
	if result.Error != nil {
		log.Print("No matching TakeProfits found....")
		return result.Error
	}
	var entry cryptodb.Order
	result = db.Where("plan_id = ? AND order_kind = ?", p.ID, cryptodb.Entry).Take(&entry)
	if result.Error != nil {
		log.Print("No matching TakeProfits found....")
		return result.Error
	}
	var pair cryptodb.Pair
	result = db.Where("id = ?", p.PairID).Find(&pair)
	if result.Error != nil {
		return result.Error
	}

	for _, o := range takeProfits {
		if !o.Price.IsZero() {
			log.Printf("Cancelling take profit (@%s)", o.Price.String())
			err := e.CancelOrder(pair.Name, o.SystemOrderID)
			if err != nil {
				db.Create(&cryptodb.Log{
					PlanID: p.ID,
					Source: cryptodb.Server,
					Text:   "Canceling Take Profits failed",
				})
				return err
			}
			db.Create(&cryptodb.Log{
				PlanID: p.ID,
				Source: cryptodb.Server,
				Text:   "Canceling Take Profits success",
			})
		}
	}
	return nil
}

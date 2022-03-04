package main

import (
	"fmt"
	"log"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

func processPosition(position exchange.Position) (err error) {

	var plan cryptodb.Plan
	result := db.Joins("JOIN pairs ON pairs.id = plans.pair_id").Where("pairs.name = ?", position.Pair).Last(&plan)

	if result.Error != nil {
		log.Printf("Error finding plan: %s", result.Error)
		return result.Error
	}

	if (plan.Direction == cryptodb.Long && position.Side == cryptodb.Sell.String()) ||
		(plan.Direction == cryptodb.Short && position.Side == cryptodb.Buy.String()) {
		return
	}

	if !plan.AverageEntryPrice.Equal(position.Price) && !position.Price.IsZero() {
		plan.AverageEntryPrice = position.Price
		logEntry := &cryptodb.Log{
			Source: cryptodb.Server,
			PlanID: plan.ID,
			Text:   fmt.Sprintf("Updated Average Entry Price to %s", plan.AverageEntryPrice.String()), // TODO: format to pair
		}
		db.Create(&logEntry)
		db.Save(&plan)
	}

	return nil
}

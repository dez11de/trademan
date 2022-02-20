package main

import (
	"fmt"
	"log"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

func processPosition(p exchange.Position) (err error) {

	log.Printf("Processing position %+v", p)
	var plan cryptodb.Plan
	result := db.Joins("JOIN pairs ON pairs.id = plans.pair_id").Where("pairs.name = ?", p.Pair).Last(&plan)

	if result.Error != nil {
		log.Printf("Error finding plan: %s", result.Error)
		return result.Error
	}

	log.Printf("Found plan: %v", plan)

	if (plan.Direction == cryptodb.Long && p.Side == cryptodb.Sell.String()) ||
		(plan.Direction == cryptodb.Short && p.Side == cryptodb.Buy.String()) {
		log.Printf("Useless position update... skipping processing.")
		return
	}

	tx := db.Begin()

	if !plan.AverageEntryPrice.Equal(p.EntryPrice) {
		plan.AverageEntryPrice = p.EntryPrice

		logEntry := &cryptodb.Log{
			Source: cryptodb.Exchange,
			PlanID: plan.ID,
			Text:   fmt.Sprintf("Updated Average Entry Price to %s", plan.AverageEntryPrice.String()), // TODO: format to pair
		}
		result = tx.Create(&logEntry)
		if result.Error != nil {
			log.Printf("Error writing log to db: %s", result.Error)
			tx.Rollback()
			return result.Error
		}
	}

	// TODO: if Size > 0 Plan.Status = partially if Size = EntrySize Plan Status = Filled

	result = tx.Save(&plan)
	if result.Error != nil {
		log.Printf("Error saving plan: %s", err)
		tx.Rollback()
		return result.Error
	}

	result = tx.Commit()

	log.Print("Processing position completed")
	return result.Error
}

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

	tx := db.Begin()

	if !plan.AverageEntryPrice.Equal(p.EntryPrice) {
		plan.AverageEntryPrice = p.EntryPrice

		var logEntry cryptodb.Log
		logEntry.Source = cryptodb.Exchange
		logEntry.PlanID = plan.ID
		logEntry.Text = fmt.Sprintf("Updated Average Entry Price to %s", plan.AverageEntryPrice.String()) // TODO: format to pair

		result = tx.Create(&logEntry)
		if result.Error != nil {
			log.Printf("Error writing log to db: %s", result.Error)
			tx.Rollback()
			return result.Error
		}
	}

	if !plan.LatestValue.Equal(p.PositionValue) {
		plan.LatestValue = p.PositionValue

		var logEntry cryptodb.Log
		logEntry.Source = cryptodb.Exchange
		logEntry.PlanID = plan.ID
		logEntry.Text = fmt.Sprintf("Updated Position Value to %s", plan.LatestValue.String()) // TODO: format to pair

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

	return result.Error
}

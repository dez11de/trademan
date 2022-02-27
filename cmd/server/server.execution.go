package main

import (
	"log"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

func processExecution(e exchange.Execution) (err error) {
	log.Printf("Processing execution %+v", e)

	var o cryptodb.Order
	var p cryptodb.Plan

	result := db.Where("system_order_id = ?", e.OrderID).First(&o)

	if result.Error != nil {
		log.Printf("No matching order found. %s", result.Error)
		return result.Error
	}

	result = db.Where("id = ?", o.PlanID).First(&p)
	if result.Error != nil {
		log.Printf("No matching plan found. %s", result.Error)
		return result.Error
	}

	if !e.ExecFee.IsZero() {
		log.Printf("Adding %s to plan fee", e.ExecFee.String())
		p.Fee.Add(e.ExecFee)
		result = db.Save(&p)
		if result.Error != nil {
			log.Printf("[execution] an error occured saving plan: %s", result.Error)
			return result.Error
		}
	}

	return err
}

package main

import (
	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

func processExecution(e exchange.Execution) (err error) {
	var o cryptodb.Order
	var p cryptodb.Plan

	result := db.Where("system_order_id = ?", e.OrderID).First(&o)
	if result.Error != nil {
		return result.Error
	}

	result = db.Where("id = ?", o.PlanID).First(&p)
	if result.Error != nil {
		return result.Error
	}

	if !e.ExecFee.IsZero() {
		p.Fee = p.Fee.Add(e.ExecFee)
		result = db.Save(&p)
		if result.Error != nil {
			return result.Error
		}
	}

	return err
}

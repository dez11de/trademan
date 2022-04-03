package main

import (
	"fmt"

	"github.com/dez11de/cryptodb"
)

func processFundings() (err error) {
	var plans []cryptodb.Plan
	result := db.Where("status BETWEEN ? AND ?", cryptodb.Created, cryptodb.Logged).Find(&plans)
	if result.Error != nil {
		return result.Error
	}

	for _, p := range plans {
		var pair cryptodb.Pair
		result := db.Where("id = ?", p.PairID).First(&pair)
		if result.Error != nil {
			return result.Error
		}

		f, err := e.GetRecentFunding(pair.Name)
		if err != nil {
			db.Create(&cryptodb.Log{
				PlanID: p.ID,
				Source: cryptodb.Server,
				Text:   fmt.Sprintf("Unable to process funding fee for %s: %s", pair.Name, err.Error()),
			})
			return err
		}

        p.Fee.Add(f.ExecFee)
        db.Save(&p)

        db.Create(&cryptodb.Log{
			PlanID: p.ID,
			Source: cryptodb.Server,
			Text:   fmt.Sprintf("Added %s %s funding to costs.", f.ExecFee.StringFixed(pair.PriceScale), pair.Name),
		})
	}

	return nil
}

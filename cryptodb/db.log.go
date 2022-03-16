package cryptodb

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func logPlanDifferences(tx *gorm.DB, logSource LogSource, pair Pair, oldPlan, newPlan Plan) error {
	var logEntry Log
	logEntry.PlanID = newPlan.ID
	logEntry.Source = logSource

	log.Print("Starting new transaction for difference logging")
	subTx := tx.Debug().Begin()
	// if oldPlan.Status != newPlan.Status {
	// 	logEntry.Text = fmt.Sprintf("\tStatus changed from %s to %s.", oldPlan.Status.String(), newPlan.Status.String())
	// 	result := subTx.Create(&logEntry)
	// 	if result.Error != nil {
	// 		subTx.Rollback()
	// 		return
	// 	}
	// }

	if !oldPlan.Risk.Equal(newPlan.Risk) {
		logEntry.Text = fmt.Sprintf("Risk changed from %s to %s.", oldPlan.Risk.StringFixed(2), newPlan.Risk.StringFixed(2))
		result := subTx.Create(logEntry)
		if result.Error != nil {
			log.Printf("Error logging risk change: %s", result.Error)
			subTx.Rollback()
			return result.Error
		}
	}

	if oldPlan.Status != newPlan.Status {
		logEntry.Text = fmt.Sprintf("Tradingview plan changed from %s to %s.", oldPlan.TradingViewPlan, newPlan.TradingViewPlan)
		result := subTx.Create(logEntry)
		if result.Error != nil {
			log.Printf("Error logging tvlink change: %s", result.Error)
			subTx.Rollback()
			return result.Error
		}
	}

	if oldPlan.RewardRiskRatio != newPlan.RewardRiskRatio {
		logEntry.Text = fmt.Sprintf("RRR changed from %.2f to %.2f.", oldPlan.RewardRiskRatio, newPlan.RewardRiskRatio)
		result := subTx.Create(logEntry)
		if result.Error != nil {
			log.Printf("Error logging RRR change: %s", result.Error)
			subTx.Rollback()
			return result.Error
		}
	}

	if !oldPlan.Profit.Equal(newPlan.Profit) {
		logEntry.Text = fmt.Sprintf("Profit changed from %s to %s.",
			oldPlan.Profit.StringFixed(pair.PriceScale),
			newPlan.Profit.StringFixed(2))
		result := subTx.Create(logEntry)
		if result.Error != nil {
			log.Printf("Error logging RRR change: %s", result.Error)
			subTx.Rollback()
			return result.Error
		}
	}
    return subTx.Debug().Commit().Error
}

func logOrderDifferences(tx *gorm.DB, logSource LogSource, pair Pair, oldOrders, newOrders []Order) {
	var logEntry Log
	logEntry.PlanID = oldOrders[0].PlanID
	logEntry.Source = logSource

	subTx := tx.Begin()
	for i := 0; i <= len(oldOrders)-1; i++ {
		var orderName string
		switch oldOrders[i].OrderKind {
		case MarketStopLoss:
			orderName = "(market)StopLoss"
		case LimitStopLoss:
			orderName = "(limit)StopLoss"
		case Entry:
			orderName = "Entry"
		case TakeProfit:
			orderName = fmt.Sprintf("Take profit #%d", i-2)
		}
		if oldOrders[i].Status != newOrders[i].Status {
			logEntry.Text = fmt.Sprintf("Status of %s changed from %s to %s.", orderName, oldOrders[i].Status.String(), newOrders[i].Status.String())
			result := subTx.Create(&logEntry)
			if result.Error != nil {
				subTx.Rollback()
				return
			}
		}
		if !oldOrders[i].Price.Equal(newOrders[i].Price) {
			logEntry.Text = fmt.Sprintf("Price of %s changed from %s to %s.", orderName, oldOrders[i].Price.StringFixed(pair.PriceScale), newOrders[i].Price.StringFixed(pair.PriceScale))
			result := subTx.Create(&logEntry)
			if result.Error != nil {
				subTx.Rollback()
				return
			}
		}
	}
	subTx.Commit()
}

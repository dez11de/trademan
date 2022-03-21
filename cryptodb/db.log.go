package cryptodb

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func logPlanDifferences(tx *gorm.DB, logSource LogSource, pair Pair, oldPlan, newPlan Plan) error {
	if !oldPlan.Risk.Equal(newPlan.Risk) {
		result := tx.Create(&Log{
			PlanID: oldPlan.ID,
			Source: logSource,
			Text:   fmt.Sprintf("Risk changed from %s to %s.", oldPlan.Risk.StringFixed(2), newPlan.Risk.StringFixed(2)),
		})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	if oldPlan.Status != newPlan.Status {
		result := tx.Create(&Log{
			PlanID: oldPlan.ID,
			Source: logSource,
			Text:   fmt.Sprintf("Tradingview plan changed from %s to %s.", oldPlan.TradingViewPlan, newPlan.TradingViewPlan),
		})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	if oldPlan.RewardRiskRatio != newPlan.RewardRiskRatio {
		result := tx.Create(&Log{
			PlanID: oldPlan.ID,
			Source: logSource,
			Text:   fmt.Sprintf("RRR changed from %.2f to %.2f.", oldPlan.RewardRiskRatio, newPlan.RewardRiskRatio),
		})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	if !oldPlan.Profit.Equal(newPlan.Profit) {
		result := tx.Create(&Log{
			PlanID: oldPlan.ID,
			Source: logSource,
			Text:   fmt.Sprintf("Profit changed from %s to %s.", oldPlan.Profit.StringFixed(pair.PriceScale), newPlan.Profit.StringFixed(2)),
		})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}
    return nil
}

func logOrderDifferences(tx *gorm.DB, logSource LogSource, pair Pair, plan Plan, oldOrders, newOrders []Order) error {
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
		if !oldOrders[i].Price.Equal(newOrders[i].Price) {
			result := tx.Create(&Log{
				PlanID: plan.ID,
				Source: logSource,
				Text:   fmt.Sprintf("Price of %s changed from %s to %s.", orderName, oldOrders[i].Price.StringFixed(pair.PriceScale), newOrders[i].Price.StringFixed(pair.PriceScale)),
			})
			if result.Error != nil {
				tx.Rollback()
				return result.Error
			}
		}
	}
    return nil 
}

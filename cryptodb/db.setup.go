package cryptodb

import (
	"fmt"
)

func (db *Database) CreateSetup(s *Setup) (err error) {
	tx := db.Begin()

	err = db.CreatePlan(&s.Plan)
	if err != nil {
		tx.Rollback()
		return err
	}

	for i := range s.Orders {
		s.Orders[i].PlanID = s.Plan.ID
        s.Orders[i].ExchangeOrderID = fmt.Sprintf("TM-%d-%d-%d", s.Plan.ID, s.Orders[i].OrderKind, s.Plan.CreatedAt.Unix())
	}

	err = db.CreateOrders(s.Orders)
	if err != nil {
		tx.Rollback()
		return err
	}

	var logEntry Log
	logEntry.PlanID = s.Plan.ID
	logEntry.Source = SourceUser
	logEntry.Text = "Plan created."
	db.CreateLog(&logEntry)

	tx.Commit()

	return err
}

func (db *Database) SaveSetup(logSource LogSource, newSetup *Setup) (err error) {
	oldPlan, err := db.GetPlan(newSetup.Plan.ID)
	if err != nil {
		return err
	}

	tx := db.Begin()

	err = db.SavePlan(&newSetup.Plan)
	if err != nil {
		tx.Rollback()
		return err
	}

	db.logPlanDifferences(logSource, oldPlan, newSetup.Plan)

	oldOrders, err := db.GetOrders(newSetup.Plan.ID)

	err = db.SaveOrders(newSetup.Orders)
	if err != nil {
		tx.Rollback()
		return err
	}

	db.logOrderDifferences(logSource, newSetup.Plan.PairID, oldOrders, newSetup.Orders)

	tx.Commit()

	return err
}

func (db *Database) logPlanDifferences(logSource LogSource, oldPlan, newPlan Plan) {
	// Not comparing pair and direction. Those shouldn' change.
	var logEntry Log
	logEntry.PlanID = oldPlan.ID
	logEntry.Source = logSource

	if oldPlan.Status != newPlan.Status {
		logEntry.Text = fmt.Sprintf("\tStatus changed from %s to %s.", oldPlan.Status.String(), newPlan.Status.String())
		db.CreateLog(&logEntry)
	}

	if !oldPlan.Risk.Equal(newPlan.Risk) {
		logEntry.Text = fmt.Sprintf("Risk changed from %s to %s.", oldPlan.Risk.StringFixed(2), newPlan.Risk.StringFixed(2))
		db.CreateLog(&logEntry)
	}

	if oldPlan.Status != newPlan.Status {
		logEntry.Text = fmt.Sprintf("Tradingview plan changed from %s to %s.", oldPlan.TradingViewPlan, newPlan.TradingViewPlan)
		db.CreateLog(&logEntry)
	}

	if oldPlan.RewardRiskRatio != newPlan.RewardRiskRatio {
		logEntry.Text = fmt.Sprintf("RRR changed from %.2f to %.2f.", oldPlan.RewardRiskRatio, newPlan.RewardRiskRatio)
		db.CreateLog(&logEntry)
	}

	if !oldPlan.Profit.Equal(newPlan.Profit) {
		logEntry.Text = fmt.Sprintf("Profit changed from %s to %s.",
			oldPlan.Profit.StringFixed(2), // TODO: get number of decimals from pair, multiple places in this file
			newPlan.Profit.StringFixed(2))
		db.CreateLog(&logEntry)
	}
}

func (db *Database) logOrderDifferences(logSource LogSource, pairID uint, oldOrders, newOrders []Order) {
	pair, _ := db.GetPair(pairID)
	var logEntry Log
	logEntry.PlanID = oldOrders[0].PlanID
	logEntry.Source = logSource

	for i := 0; i <= len(oldOrders)-1; i++ {
		var orderName string
		switch oldOrders[i].OrderKind {
		case KindMarketStopLoss:
			orderName = "(market)StopLoss"
		case KindLimitStopLoss:
			orderName = "(limit)StopLoss"
		case KindEntry:
			orderName = "Entry"
		case KindTakeProfit:
			orderName = fmt.Sprintf("Take profit #%d", i-2)
		}
		if oldOrders[i].Status != newOrders[i].Status {
			logEntry.Text = fmt.Sprintf("Status of %s changed from %s to %s.", orderName, oldOrders[i].Status.String(), newOrders[i].Status.String())
			db.CreateLog(&logEntry)
		}
		if !oldOrders[i].Price.Equal(newOrders[i].Price) {
			logEntry.Text = fmt.Sprintf("Price of %s changed from %s to %s.", orderName, oldOrders[i].Price.StringFixed(int32(pair.PriceScale)), newOrders[i].Price.StringFixed(int32(pair.PriceScale)))
			db.CreateLog(&logEntry)
		}
	}
}

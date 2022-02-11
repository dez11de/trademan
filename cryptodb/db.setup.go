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
	}

	err = db.CreateOrders(s.Orders)
	if err != nil {
		tx.Rollback()
		return err
	}

	for i := range s.Orders {
		s.Orders[i].ExchangeOrderID = fmt.Sprintf("TM-%04d-%05d-%d", s.Plan.ID, s.Orders[i].ID, s.Plan.CreatedAt.Unix())
	}
	err = db.SaveOrders(s.Orders)
	if err != nil {
		tx.Rollback()
		return err
	}

	var logEntry Log
	logEntry.PlanID = s.Plan.ID
	logEntry.Source = SourceUser
	logEntry.Text = "Plan created."
	err = db.CreateLog(&logEntry)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return err
}

func (db *Database) SaveSetup(logSource LogSource, newSetup *Setup) (err error) {
	pair, err := db.GetPair(newSetup.Plan.PairID)
	if err != nil {
		return err
	}
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

	db.logPlanDifferences(logSource, pair, oldPlan, newSetup.Plan)

	oldOrders, err := db.GetOrders(newSetup.Plan.ID)

	err = db.SaveOrders(newSetup.Orders)
	if err != nil {
		tx.Rollback()
		return err
	}

	db.logOrderDifferences(logSource, pair, oldOrders, newSetup.Orders)

	tx.Commit()

	return err
}

func (db *Database) logPlanDifferences(logSource LogSource, pair Pair, oldPlan, newPlan Plan) {
	var logEntry Log
	logEntry.PlanID = oldPlan.ID
	logEntry.Source = logSource

	tx := db.Begin()
	if oldPlan.Status != newPlan.Status {
		logEntry.Text = fmt.Sprintf("\tStatus changed from %s to %s.", oldPlan.Status.String(), newPlan.Status.String())
		err := db.CreateLog(&logEntry)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	if !oldPlan.Risk.Equal(newPlan.Risk) {
		logEntry.Text = fmt.Sprintf("Risk changed from %s to %s.", oldPlan.Risk.StringFixed(2), newPlan.Risk.StringFixed(2))
		err := db.CreateLog(&logEntry)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	if oldPlan.Status != newPlan.Status {
		logEntry.Text = fmt.Sprintf("Tradingview plan changed from %s to %s.", oldPlan.TradingViewPlan, newPlan.TradingViewPlan)
		err := db.CreateLog(&logEntry)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	if oldPlan.RewardRiskRatio != newPlan.RewardRiskRatio {
		logEntry.Text = fmt.Sprintf("RRR changed from %.2f to %.2f.", oldPlan.RewardRiskRatio, newPlan.RewardRiskRatio)
		err := db.CreateLog(&logEntry)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	if !oldPlan.Profit.Equal(newPlan.Profit) {
		logEntry.Text = fmt.Sprintf("Profit changed from %s to %s.",
			oldPlan.Profit.StringFixed(pair.PriceScale),
			newPlan.Profit.StringFixed(2))
		err := db.CreateLog(&logEntry)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
}

func (db *Database) logOrderDifferences(logSource LogSource, pair Pair, oldOrders, newOrders []Order) {
	var logEntry Log
	logEntry.PlanID = oldOrders[0].PlanID
	logEntry.Source = logSource

	tx := db.Begin()
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
			err := db.CreateLog(&logEntry)
			if err != nil {
				tx.Rollback()
				return
			}
		}
		if !oldOrders[i].Price.Equal(newOrders[i].Price) {
			logEntry.Text = fmt.Sprintf("Price of %s changed from %s to %s.", orderName, oldOrders[i].Price.StringFixed(pair.PriceScale), newOrders[i].Price.StringFixed(pair.PriceScale))
			err := db.CreateLog(&logEntry)
			if err != nil {
				tx.Rollback()
				return
			}
		}
	}
	tx.Commit()
}

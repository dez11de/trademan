package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddPlanStatement() (err error) {
	db.addPlanStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET PairID=?, Status=?, Side=?, Risk=?, Notes=?, TradingViewPlan=?, RewardRiskRatio=?, Profit=?", db.planTableName))
	return err
}

func (db *Database) AddPlan(p Plan) (TradeID int64, err error) {
	result, err := db.addPlanStatement.Exec(p.PlanID, p.Status.String(), p.Side.String(), p.Risk, p.Notes, p.TradingViewPlan, p.RewardRiskRatio, p.Profit)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *Database) GetPlans() (p []Plan, err error) {
	rows, err := db.database.Query("SELECT * FROM `PLAN`;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var plan Plan
		if err := rows.Scan(&plan.PlanID, &plan.PairID, &plan.Status, &plan.Side, &plan.Risk, &plan.Notes, &plan.TradingViewPlan, &plan.RewardRiskRatio, &plan.Profit, &plan.EntryTime, &plan.ModifyTime); err != nil {
			return nil, err
		}
		p = append(p, plan)
	}
	return p, nil
}

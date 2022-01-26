package cryptodb

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (db *api) AddPlan(p Plan) (planID int64, err error) {
	result, err := db.database.Exec(
		`INSERT INTO 'PLAN' (PairID, Status, Side, Risk, Notes, TradingViewPlan, RewardRiskRatio, Profit) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.PairID, p.Status, p.Side, p.Risk, p.Notes, p.TradingViewPlan, p.RewardRiskRatio, p.Profit)
	if err != nil {
		log.Printf("[AddPlan] error occured executing statement: %v", err)
	}

	return result.LastInsertId()
}

func (db *api) GetPlans() (p []Plan, err error) {

	rows, err := db.database.Query("SELECT * FROM `PLAN`;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plan Plan
	for rows.Next() {
		if err := rows.Scan(&plan.PlanID, &plan.PairID, &plan.Status, &plan.Side, &plan.Risk, &plan.Notes, &plan.TradingViewPlan, &plan.RewardRiskRatio, &plan.Profit, &plan.EntryTime, &plan.ModifyTime); err != nil {
			return nil, err
		}
		p = append(p, plan)
	}

	return p, nil
}

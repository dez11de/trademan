package cryptodb

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (db *api) AddPlan(p Plan) (planID int64, err error) {
    // TODO: also add/update orders as a transaction, and rollback if anything fails
	result, err := db.database.Exec(
		`INSERT INTO PLAN (PairID, Status, Side, Risk, Notes, TradingViewPlan, RewardRiskRatio, Profit) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.PairID, p.Status, p.Side, p.Risk, p.Notes, p.TradingViewPlan, p.RewardRiskRatio, p.Profit)
	if err != nil {
		log.Printf("[AddPlan] error occured executing statement: %v", err)
	}

	return result.LastInsertId()
}

func (db *api) GetPlans() (p []Plan, err error) {
    // TODO: also read orders
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

func (db *api) GetPlan(id int64) (p Plan, err error) {
    // TODO: also read orders
	row := db.database.QueryRow("SELECT * FROM `PLAN` WHERE PlanID=?;", id)

	err = row.Scan(&p.PlanID, &p.PairID, &p.Status, &p.Side, &p.Risk, &p.Notes, &p.TradingViewPlan, &p.RewardRiskRatio, &p.Profit, &p.EntryTime, &p.ModifyTime)

    if err != nil {
		return Plan{}, err
	}

	return p, nil
}

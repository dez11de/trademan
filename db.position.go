package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddPositionStatement() (err error) {
	db.addPositionStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET PairID=?, Status=?, Side=?, Risk=?, Notes=?, TradingViewPlan=?, RewardRiskRatio=?, Profit=?", db.positionTableName))
	return err
}

func (db *Database) AddPosition(p Position) (TradeID int64, err error) {
	result, err := db.addPositionStatement.Exec(p.PositionID, p.Status.String(), p.Side.String(), p.Risk, p.Notes, p.TradingViewPlan, p.RewardRiskRatio, p.Profit)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *Database) GetPositions() (p []Position, err error) {
	rows, err := db.database.Query("SELECT * FROM `POSITION`;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pos Position
		if err := rows.Scan(&pos.PositionID, &pos.PairID, &pos.Status, &pos.Side, &pos.Risk, &pos.Notes, &pos.TradingViewPlan, &pos.RewardRiskRatio, &pos.Profit, &pos.EntryTime, &pos.ModifyTime); err != nil {
			return nil, err
		}
		p = append(p, pos)
	}
	return p, nil
}

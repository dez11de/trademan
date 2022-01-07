package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddPositionStatement() (err error) {
	db.addPositionStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET Symbol=?, Status=?, Risk=?, `Size`=?, EntryPrice=?, HardStopLoss=?, Notes=?, TradingViewPlan=?, RewardRiskRatio=?, Profit=?", db.positionTableName))
	return err
}

func (db *Database) AddPosition(p Position) (TradeID int64, err error) {
	result, err := db.addPositionStatement.Exec(p.Symbol, p.Side.String(), p.Risk, p.Size, p.EntryPrice, p.HardStopLoss, p.Notes, p.TradingViewPlan, p.RewardRiskRatio, p.Profit)
	if err != nil {
		log.Printf("Error adding symbol %v", err)
		return 0, err
	}
	return result.LastInsertId()
}

func (db *Database) GetPositions() (p []Position, err error) {
	rows, err := db.database.Query("SELECT Symbol, Status, Size FROM `POSITION`;")
	if err != nil {
		log.Printf("Error querying database %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pos Position
		if err := rows.Scan(&pos.Symbol, &pos.Status, &pos.Size); err != nil {
			log.Printf("Error querying database %v", err)
			return nil, err
		}
		p = append(p, pos)
	}
	return p, nil
}

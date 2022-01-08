package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddPairStatement() (err error) {
	db.addPairStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET Pair=?, BaseCurrency=?, QuoteCurrency=?, PriceScale=?, TakerFee=?, MakerFee=?, MinLeverage=?, MaxLeverage=?, LeverageStep=?, MinPrice=?, MaxPrice=?, TickSize=?, MinOrderSize=?, MaxOrderSize=?, StepOrderSize=?", db.pairTableName))
	return err
}

func (db *Database) AddPair(p pair) (PairID int64, err error) {
	// TODO: handle updates of existing symbols
	result, err := db.addPairStatement.Exec(p.Pair, p.BaseCurrency, p.QuoteCurrency, p.PriceScale, p.TakerFee, p.MakerFee, p.Leverage.Min, p.Leverage.Max, p.Leverage.Step, p.Price.Min, p.Price.Max, p.Price.Tick, p.OrderSize.Min, p.OrderSize.Max, p.OrderSize.Step)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (db *Database) GetSymbols() (map[string]pair, error) {
	rows, err := db.database.Query("SELECT * FROM `PAIR`;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var s pair
	symbolMap := make(map[string]pair)
	for rows.Next() {
		err = rows.Scan(&s.PairID, &s.Pair, &s.BaseCurrency, &s.QuoteCurrency, &s.PriceScale, &s.TakerFee, &s.MakerFee, &s.Leverage.Min, &s.Leverage.Max, &s.Leverage.Step, &s.Price.Min, &s.Price.Max, &s.Price.Tick, &s.OrderSize.Min, &s.OrderSize.Max, &s.OrderSize.Step)
		if err != nil {
		}
		symbolMap[s.Pair] = s
	}
	return symbolMap, nil
}

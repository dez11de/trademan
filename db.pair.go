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

func (db *Database) GetPairs() (pairs map[string]pair, err error) {
	rows, err := db.database.Query("SELECT * FROM `PAIR`;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var p pair
	for rows.Next() {
		err = rows.Scan(&p.PairID, &p.Pair, &p.BaseCurrency, &p.QuoteCurrency, &p.PriceScale, &p.TakerFee, &p.MakerFee, &p.Leverage.Min, &p.Leverage.Max, &p.Leverage.Step, &p.Price.Min, &p.Price.Max, &p.Price.Tick, &p.OrderSize.Min, &p.OrderSize.Max, &p.OrderSize.Step)
		if err != nil {
		}
		pairs[p.Pair] = p
	}
	return pairs, nil
}

func (db *Database) GetPairID(p string) int64 {
	for symbol, pair := range db.PairCache {
		if symbol == p {
			return pair.PairID
		}
	}
	return -1
}

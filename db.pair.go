package main

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddPairStatement() (err error) {
	db.addPairStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET Pair=?, BaseCurrency=?, QuoteCurrency=?, PriceScale=?, TakerFee=?, MakerFee=?, MinLeverage=?, MaxLeverage=?, LeverageStep=?, MinPrice=?, MaxPrice=?, TickSize=?, MinOrderSize=?, MaxOrderSize=?, StepOrderSize=?", db.pairTableName))
	return err
}

func (db *Database) AddPair(p Pair) (PairID int64, err error) {
	// TODO: handle updates of existing symbols
	result, err := db.addPairStatement.Exec(p.Pair, p.BaseCurrency, p.QuoteCurrency, p.PriceScale, p.TakerFee, p.MakerFee, p.Leverage.Min, p.Leverage.Max, p.Leverage.Step, p.Price.Min, p.Price.Max, p.Price.Tick, p.OrderSize.Min, p.OrderSize.Max, p.OrderSize.Step)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (db *Database) GetPairs() (pairs map[string]Pair, err error) {
	pairs = make(map[string]Pair)
	rows, err := db.database.Query("SELECT * FROM `PAIR` ORDER BY Pair;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var p Pair
	for rows.Next() {
		err = rows.Scan(&p.PairID, &p.Pair, &p.BaseCurrency, &p.QuoteCurrency, &p.PriceScale, &p.TakerFee, &p.MakerFee, &p.Leverage.Min, &p.Leverage.Max, &p.Leverage.Step, &p.Price.Min, &p.Price.Max, &p.Price.Tick, &p.OrderSize.Min, &p.OrderSize.Max, &p.OrderSize.Step)
		if err != nil {
		}
		pairs[p.Pair] = p
	}
	return pairs, nil
}

func (db *Database) GetPairID(p string) int64 {
	if id, ok := db.PairCache[p]; ok {
		return id.PairID
	}
	return -1
}

func (db *Database) GetPairString(id int64) string {
	for _, pair := range db.PairCache {
		if pair.PairID == id {
			return pair.Pair
		}
	}
	return ""
}

func (db *Database) GetPairFromString(p string) (pair Pair, err error) {
	pair, ok := db.PairCache[p]
	if !ok {
		return Pair{}, errors.New("pair not found in cache")
	}
	return pair, nil
}

func (db *Database) GetPairFromID(i int64) (pair Pair, err error) {
	return db.GetPairFromString(db.GetPairString(i))
}

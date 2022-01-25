package cryptodb

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddPairStatement() (err error) {
	db.addPairStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET Pair=?, BaseCurrency=?, QuoteCurrency=?, PriceScale=?, TakerFee=?, MakerFee=?, MinLeverage=?, MaxLeverage=?, LeverageStep=?, MinPrice=?, MaxPrice=?, TickSize=?, MinOrderSize=?, MaxOrderSize=?, StepOrderSize=?", db.config.pairTableName))
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
			// TODO: shouldn't i be doing something?
			log.Print(err)
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
		log.Print(err)
		return Pair{}, err
	}
	return pair, nil
}

func (db *Database) GetPairFromID(i int64) (pair Pair, err error) {
	return db.GetPairFromString(db.GetPairString(i))
}

func (db *Database) SearchPairs(s string) (pairs []string, err error) {
	rows, err := db.database.Query(fmt.Sprintf("SELECT Pair FROM `PAIR` WHERE Pair LIKE '%%%s%%' ORDER BY Pair", s))
	if err != nil {
		return nil, err
	}
	var pair string
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&pair)
		if err != nil {
			// TODO: shouldn't i be doing something?
			log.Print(err)
		}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}

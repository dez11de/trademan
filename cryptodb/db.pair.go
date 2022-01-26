package cryptodb

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)


func (db *api) AddPair(p Pair) (PairID int64, err error) {
	// TODO: handle updates of existing symbols
    addStmt, err := db.database.Prepare("INSERT INTO `PAIR` (Pair, BaseCurrency, QuoteCurrency, PriceScale, TakerFee, MakerFee, MinLeverage, MaxLeverage, LeverageStep, MinPrice, MaxPrice, TickSize, MinOrderSize, MaxOrderSize, StepOrderSize) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

    result, err := addStmt.Exec(p.Pair, p.BaseCurrency, p.QuoteCurrency, p.PriceScale, p.TakerFee, p.MakerFee, p.Leverage.Min, p.Leverage.Max, p.Leverage.Step, p.Price.Min, p.Price.Max, p.Price.Tick, p.OrderSize.Min, p.OrderSize.Max, p.OrderSize.Step)

	return result.LastInsertId()
}

func (db *api) GetPairs() (pairs map[string]Pair, err error) {
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

func (db *api) GetPairID(p string) int64 {
    //TODO: reimplement as database query
	return -1
}

func (db *api) GetPairString(id int64) string {
    //TODO: reimplement as database query
	return ""
}

func (db *api) GetPairFromString(p string) (pair Pair, err error) {
    //TODO: reimplement as database query
	return pair, nil
}

func (db *api) GetPairFromID(i int64) (pair Pair, err error) {
    //TODO: reimplement as database query
	return db.GetPairFromString(db.GetPairString(i))
}

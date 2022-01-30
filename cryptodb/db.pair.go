package cryptodb

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (db *api) WritePair(pair Pair) (RowsAffected int64, err error) {
    existingPair, err := db.GetPair(pair.Pair)
    if err != nil {
        log.Printf("Error getting pair: %v", err)
        return 0, err
    }
    if existingPair.Pair != pair.Pair {
        log.Printf("Adding NEW pair")
		return db.AddPair(pair)
	} else {
        log.Printf("Updating EXISTING pair")
		return db.UpdatePair(pair)
	}
}

func (db *api) AddPair(p Pair) (PairID int64, err error) {
	result, err := db.database.Exec(
		`INSERT INTO PAIR (Pair, BaseCurrency, QuoteCurrency, PriceScale, TakerFee, MakerFee, MinLeverage, MaxLeverage, LeverageStep, MinPrice, MaxPrice, TickSize, MinOrderSize, MaxOrderSize, StepOrderSize) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Pair, p.BaseCurrency, p.QuoteCurrency, p.PriceScale, p.TakerFee, p.MakerFee, p.Leverage.Min, p.Leverage.Max, p.Leverage.Step, p.Price.Min, p.Price.Max, p.Price.Tick, p.OrderSize.Min, p.OrderSize.Max, p.OrderSize.Step)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (db *api) UpdatePair(p Pair) (RowsAffected int64, err error) {
	result, err := db.database.Exec(
		`UPDATE PAIR SET BaseCurrency=?, QuoteCurrency=?, PriceScale=?, TakerFee=?, MakerFee=?, MinLeverage=?, MaxLeverage=?, LeverageStep=?, MinPrice=?, MaxPrice=?, TickSize=?, MinOrderSize=?, MaxOrderSize=?, StepOrderSize=? WHERE PairID=?`,
		p.BaseCurrency, p.QuoteCurrency, p.PriceScale, p.TakerFee, p.MakerFee, p.Leverage.Min, p.Leverage.Max, p.Leverage.Step, p.Price.Min, p.Price.Max, p.Price.Tick, p.OrderSize.Min, p.OrderSize.Max, p.OrderSize.Step, p.PairID)
	if err != nil {
		log.Printf("error updating: %v", err)
		return 0, err
	}

	return result.RowsAffected()
}

func (db *api) GetPairs() (pairs map[string]Pair, err error) {
	pairs = make(map[string]Pair)
	rows, err := db.database.Query("SELECT * FROM PAIR ORDER BY Pair")
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

func (db *api) GetPair(p string) (pair Pair, err error) {
	rows, err := db.database.Query(`SELECT * FROM PAIR WHERE Pair=? LIMIT 1;`, p)
	if err != nil {
		log.Printf("error executing query %v", err)
		return Pair{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&pair.PairID, &pair.Pair, &pair.BaseCurrency, &pair.QuoteCurrency, &pair.PriceScale, &pair.TakerFee, &pair.MakerFee, &pair.Leverage.Min, &pair.Leverage.Max, &pair.Leverage.Step, &pair.Price.Min, &pair.Price.Max, &pair.Price.Tick, &pair.OrderSize.Min, &pair.OrderSize.Max, &pair.OrderSize.Step)
		if err != nil {
			// TODO: shouldn't i be doing something?
			log.Print(err)
		}
	}

	return pair, err
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

func (db *api) SearchPairs(s string) (pairs []string, err error) {
    rows, err := db.database.Query("SELECT Pair FROM PAIR WHERE Pair LIKE ? ORDER BY Pair ASC", "%"+s+"%")
    if err != nil {
        log.Printf("No pairs found? %v", err)
    }
	defer rows.Close()

	var p Pair
	for rows.Next() {
		err = rows.Scan(&p.Pair)
		if err != nil {
			// TODO: shouldn't i be doing something?
			log.Print(err)
		}
        log.Printf("Returning pair %s", p.Pair)
        pairs = append(pairs, p.Pair)
	}
	return pairs, nil
}

package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddSymbolStatement() (err error) {
	db.addSymbolStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET Symbol=?, BaseCurrency=?, QuoteCurrency=?, PriceScale=?, TakerFee=?, MakerFee=?, MinLeverage=?, MaxLeverage=?, LeverageStep=?, MinPrice=?, MaxPrice=?, TickSize=?, MinOrderSize=?, MaxOrderSize=?, StepOrderSize=?", db.symbolTableName))
	return err
}

func (db *Database) AddSymbol(s symbol) (SymbolID int64, err error) {
	// TODO: handle updates of existing symbols
	result, err := db.addSymbolStatement.Exec(s.Symbol, s.BaseCurrency, s.QuoteCurrency, s.PriceScale, s.TakerFee, s.MakerFee, s.Leverage.Min, s.Leverage.Max, s.Leverage.Step, s.Price.Min, s.Price.Max, s.Price.Tick, s.OrderSize.Min, s.OrderSize.Max, s.OrderSize.Step)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func (db *Database) GetSymbols() (symbols []symbol, err error) {
	rows, err := db.database.Query("SELECT * FROM `SYMBOL`;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s symbol
		if err := rows.Scan(&s.Symbol, &s.BaseCurrency, &s.QuoteCurrency, &s.PriceScale, &s.TakerFee, &s.MakerFee, &s.Leverage.Min, &s.Leverage.Max, &s.Leverage.Step, &s.Price.Min, &s.Price.Max, &s.Price.Tick, &s.OrderSize.Min, &s.OrderSize.Max, &s.OrderSize.Step); err != nil {
			return nil, err
		}
		symbols = append(symbols, s)
	}
	return symbols, nil
}

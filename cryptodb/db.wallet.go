package cryptodb

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func (db *api) AddBalance(b Balance) (err error) {
	// TODO: should be posible to not expect a result right?
	_, err = db.database.Exec(
		`INSERT INTO WALLET (Symbol, Equity, Available, UsedMargin, OrderMargin, PositionMargin, OCCClosingFee, OCCFundingFee, WalletBalance, DailyPnL, UnrealisedPnL, TotalPnL) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		b.Symbol, b.Equity, b.Available, b.UsedMargin, b.OrderMargin, b.PositionMargin, b.OCCClosingFee, b.OCCFundingFee, b.WalletBalance, b.DailyPnL, b.UnrealisedPnL, b.TotalPnL)
	if err != nil {
		log.Printf("[AddWallet] error occured executing statement: %v", err)
	}

	return err
}

// TODO: rewrite this into a GetBalance (symbol string) Balance, err function
func (db *api) GetRecentWallet() (wallet map[string]Balance, err error) {
	// Get most recent TimeStamp
	// TODO: this can probably be simplified if we assume fixed number of Symbols in wallet
	rows, err := db.database.Query("SELECT EntryTime FROM WALLET ORDER BY EntryTime DESC LIMIT 1;")
	if err != nil {
		return nil, err
	}

	var lastBalance Balance
	for rows.Next() {
		if err := rows.Scan(&lastBalance.EntryTime); err != nil {
			return nil, err
		}
	}
	rows.Close()

	// TODO: bonus points for doing this in one query, see GetPerformance
	rows, err = db.database.Query("SELECT * FROM WALLET ORDER BY EntryTime DESC LIMIT 1;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallet = make(map[string]Balance)
	for rows.Next() {
		if err := rows.Scan(&lastBalance.Symbol, &lastBalance.Equity, &lastBalance.Available, &lastBalance.UsedMargin, &lastBalance.OrderMargin, &lastBalance.PositionMargin, &lastBalance.OCCClosingFee, &lastBalance.OCCFundingFee, &lastBalance.WalletBalance, &lastBalance.DailyPnL, &lastBalance.UnrealisedPnL, &lastBalance.TotalPnL, &lastBalance.EntryTime); err != nil {
			return nil, err
		}
		wallet[lastBalance.Symbol] = lastBalance
	}

	return wallet, nil
}

func (db *api) GetPerformance(symbol string, periodStart time.Time) (performance float64, err error) {
	// TODO: Seems to work, but need more data. Maybe it's better to use TIMESTAMPDIFF? I don't really understand the query anyway.
	// TODO: reformat query so it's a bit easier on the eyes
	row := db.database.QueryRow("SELECT (RecentEquity - PreviousEquity) / PreviousEquity * 100 AS Performance "+
		"FROM ( SELECT ( SELECT Equity "+
		"                FROM WALLET p1 "+
		"                WHERE p1.EntryTime = x.PreviousTimestamp AND Symbol = x.Symbol "+
		"            ) AS PreviousEquity, "+
		"            ( SELECT Equity "+
		"                FROM WALLET p1 "+
		"                WHERE p1.EntryTime = x.RecentTimestamp AND Symbol = x.Symbol "+
		"            ) AS RecentEquity "+
		"        FROM ( SELECT Symbol, MIN(EntryTime) AS PreviousTimestamp, MAX(EntryTime) AS RecentTimestamp "+
		"                FROM WALLET "+
		"                WHERE Symbol = ? AND EntryTime BETWEEN ? AND NOW() "+
		"            ) x) x2;", symbol, periodStart.Format(MySQLTimestampFormat))

	if row == nil {
		log.Printf("no results found: %v", err)
        return 0, err
	}

	err = row.Scan(&performance)

	if err != nil {
		log.Printf("error: %v", err)
        return 0, err
	}

	return performance, err
}

package cryptodb

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
)

func (db *api) AddWallet(b Balance) (err error) {

    addStmt, err := db.database.Prepare("INSERT INTO `WALLET` (Symbol, Equity, Available, UsedMargin, OrderMargin, PositionMargin, OCCClosingFee, OCCFundingFee, WalletBalance, DailyPnL, UnrealisedPnL, TotalPnL) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?")
	if err != nil {
		return err
	}
    _, err = addStmt.Exec(b.Symbol, b.Equity, b.Available, b.UsedMargin, b.OrderMargin, b.PositionMargin, b.OCCClosingFee, b.OCCFundingFee, b.WalletBalance, b.DailyPnL, b.UnrealisedPnL, b.TotalPnL)

	return err
}

func (db *api) GetRecentWallet() (wallet map[string]Balance, err error) {
    // Get most recent TimeStamp TODO: this can probably be simplified if we assume fixed number of Symbols in wallet
	rows, err := db.database.Query("SELECT EntryTime FROM `WALLET` ORDER BY EntryTime DESC LIMIT 1;")
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

	rows, err = db.database.Query("SELECT * FROM `WALLET` ORDER BY EntryTime DESC LIMIT 1;")
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

func (db *api) GetPerformance(p time.Duration) decimal.Decimal {
	periodStart := time.Now().Add(-p)
	result, err := db.database.Query(fmt.Sprintf("SELECT Equity FROM WALLET WHERE Symbol='USDT' ORDER BY abs(TIMESTAMPDIFF(second, EntryTime, '%s')) LIMIT 1", periodStart.Format("2006-01-02 15:04:05")))
	if err != nil {
		log.Print(err)
	}
	result.Next()
	var balanceAtPeriodStart decimal.Decimal
	result.Scan(&balanceAtPeriodStart)
    // TODO: see GetRecentWallet()
	// currentBalance := db.WalletCache["USDT"].Equity
    // return (currentBalance.Sub(balanceAtPeriodStart).Div(currentBalance)).Mul(decimal.NewFromInt(100))
    return decimal.Zero
}

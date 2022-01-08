package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddWalletStatement() (err error) {
	db.addWalletStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET Symbol=?, Equity=?, Available=?, UsedMargin=?, OrderMargin=?, PositionMargin=?, OCCClosingFee=?, OCCFundingFee=?, WalletBalance=?, DailyPnL=?, UnrealisedPnL=?, TotalPnL=?, EntryTime=?", db.walletTableName))
	return err
}

func (db *Database) AddWallet(b balance) (err error) {
	_, err = db.addWalletStatement.Exec(b.Symbol, b.Equity, b.Available, b.UsedMargin, b.OrderMargin, b.PositionMargin, b.OCCClosingFee, b.OCCFundingFee, b.WalletBalance, b.DailyPnL, b.UnrealisedPnL, b.TotalPnL, b.EntryTime.Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}
	return err
}

func (db *Database) GetRecentWallet() (map[string]balance, error) {
	// Get most recent TimeStamp
	rows, err := db.database.Query("SELECT EntryTime FROM `WALLET` ORDER BY EntryTime DESC LIMIT 1;")
	if err != nil {
		return nil, err
	}

	var lastBalance balance
	for rows.Next() {
		if err := rows.Scan(&lastBalance.EntryTime); err != nil {
			return nil, err
		}
	}
	rows.Close()

	rows, err = db.database.Query(fmt.Sprintf("SELECT * FROM `WALLET` WHERE EntryTime='%s';", lastBalance.EntryTime.Format("2006-01-02 15:04:05")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallet := make(map[string]balance)
	for rows.Next() {
		if err := rows.Scan(&lastBalance.Symbol, &lastBalance.Equity, &lastBalance.Available, &lastBalance.UsedMargin, &lastBalance.OrderMargin, &lastBalance.PositionMargin, &lastBalance.OCCClosingFee, &lastBalance.OCCFundingFee, &lastBalance.WalletBalance, &lastBalance.DailyPnL, &lastBalance.UnrealisedPnL, &lastBalance.TotalPnL, &lastBalance.EntryTime); err != nil {
			return nil, err
		}
		wallet[lastBalance.Symbol] = lastBalance
	}

	return wallet, nil
}

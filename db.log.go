package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddLogStatement() (err error) {
	db.addLogStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET PositionID=?, Source=?, Text=?", db.logTableName))
	return err
}

func (db *Database) AddLog(tradeID int64, source LogSource, text string) (err error) {
	_, err = db.addLogStatement.Exec(tradeID, source, text)
	return err
}

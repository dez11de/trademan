package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddOrderStatement() (err error) {
	db.addOrderStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET PositionID=?, ExchangeOrderID=?, Status=?, OrderType=?, `Size`=?, TriggerPrice=?, Price=?", db.orderTableName))
	return err
}

func (db *Database) AddOrder(o Order) (OrderID int64, err error) {
	result, err := db.addOrderStatement.Exec(o.PositionID, o.ExchangeOrderID, o.Status.String(), o.OrderType.String(), o.Size, o.TriggerPrice, o.Price)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

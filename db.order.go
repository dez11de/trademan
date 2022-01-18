package main

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *Database) PrepareAddOrderStatement() (err error) {
	db.addOrderStatement, err = db.database.Prepare(fmt.Sprintf("INSERT %s SET PlanID=?, ExchangeOrderID=?, Status=?, OrderType=?, `Size`=?, TriggerPrice=?, Price=?", db.orderTableName))
	return err
}

func (db *Database) AddOrder(o Order) (OrderID int64, err error) {
	result, err := db.addOrderStatement.Exec(o.PlanID, o.ExchangeOrderID, o.Status.String(), o.OrderType.String(), o.Size, o.TriggerPrice, o.Price)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (db *Database) GetPlanOrders(id int64) (o Orders, err error) {
	rows, err := db.database.Query(fmt.Sprintf("SELECT * FROM `ORDER` where PlanID=%d", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.OrderID, &order.PlanID, &order.ExchangeOrderID, &order.Status, &order.OrderType, &order.Size, &order.TriggerPrice, &order.Price, &order.EntryTime, &order.ModifyTime); err != nil {
			return nil, err
		}
		o = append(o, order)
	}
	return o, nil
}

func (ol Orders) GetHardStopLoss() (o Order, err error) {
	for _, order := range ol {
		if order.OrderType == HardStopLoss {
			return order, nil
		}
	}
	return Order{}, errors.New("HardStopLoss order not found")
}

func (ol Orders) GetEntry() (o Order, err error) {
	for _, order := range ol {
		if order.OrderType == Entry {
			return order, nil
		}
	}
	return Order{}, errors.New("Entry order not found")
}

func (ol Orders) GetTakeProfits() (tps Orders, err error) {
	for _, o := range ol {
		if o.OrderType == TakeProfit {
			tps = append(tps, o)
		}
	}
	if len(tps) != 0 {
		return tps, nil
	}
	return Orders{}, errors.New("Take profits not found")
}

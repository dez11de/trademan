package main

import (
	"fmt"
	"log"

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

func (db *Database) GetOrders(PlanID int64) (orders Orders, err error) {
	log.Printf("[db.order.GetOrders] getting orders belonging to PlanID: %d", PlanID)
	orders = NewOrders()
	rows, err := db.database.Query(fmt.Sprintf("SELECT * FROM `ORDER` where PlanID=%d", PlanID))
	if err != nil {
		log.Print(err)
		return NewOrders(), err
	}
	defer rows.Close()

	takeProfitCount := 0
	for rows.Next() {
		var order Order
		err = rows.Scan(&order.OrderID, &order.PlanID, &order.Status, &order.ExchangeOrderID, &order.OrderType, &order.Size, &order.TriggerPrice, &order.Price, &order.EntryTime, &order.ModifyTime)
		if err != nil {
			log.Print(err)
			return Orders{}, err
		}
		switch order.OrderType {
		case typeHardStopLoss:
			orders[typeHardStopLoss] = order
		case typeSoftStopLoss:
			orders[typeSoftStopLoss] = order
		case typeEntry:
			orders[typeEntry] = order
		case typeTakeProfit:
			orders[3+takeProfitCount] = order
			takeProfitCount++
		}
	}

	for ; takeProfitCount < MaxTakeProfits; takeProfitCount++ {
		orders[3+takeProfitCount] = Order{OrderType: typeTakeProfit, PlanID: PlanID}
	}
	return orders, nil
}

func (db *Database) StoreOrders(PlanID int64) (err error) {
	log.Printf("[StoreOrders] storing orders...")
	return err
}

package cryptodb

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func (db *api) AddOrder(o Order) (OrderID int64, err error) {
    result, err := db.database.Exec(
        "INSERT INTO `ORDER` (PlanID, ExchangeOrderID, Status, OrderType, `Size`, TriggerPrice, Price) VALUES (?, ?, ?, ?, ?, ?, ?)",
        o.PlanID, o.ExchangeOrderID, o.Status, o.OrderType, o.Size, o.TriggerPrice, o.Price)
	if err != nil {
		log.Printf("[AddOrder] error occured executing statement: %v", err)
	}

	return result.LastInsertId()
}

func (db *api) GetOrders(PlanID int64) (orders Orders, err error) {
	orders = NewOrders()

	rows, err := db.database.Query(fmt.Sprintf("SELECT * FROM 'ORDER' where PlanID=%d", PlanID))
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
		case TypeHardStopLoss:
			orders[TypeHardStopLoss] = order
		case TypeSoftStopLoss:
			orders[TypeSoftStopLoss] = order
		case TypeEntry:
			orders[TypeEntry] = order
		case TypeTakeProfit:
			orders[3+takeProfitCount] = order
			takeProfitCount++
		}
	}

	for ; takeProfitCount < MaxTakeProfits; takeProfitCount++ {
		orders[3+takeProfitCount] = Order{OrderType: TypeTakeProfit, PlanID: PlanID}
	}
	return orders, nil
}

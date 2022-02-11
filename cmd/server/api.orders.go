package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
	"github.com/julienschmidt/httprouter"
)

func getOrdersHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	planID, err := strconv.Atoi(p.ByName("PlanID"))
	orders, err := db.GetOrders(uint(planID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	jsonResp, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonResp)
}

func processOrder(o exchange.Order) (err error) {
	var marketStopLossOrder cryptodb.Order
	tradeManOrder, err := db.MatchExchangeOrder(o.ExchangeOrderID)
	if err != nil {
		return err
	}
	if tradeManOrder.Status == cryptodb.StatusPlanned && (o.OrderStatus == "New" || o.OrderStatus == "Untriggered") {
		tradeManOrder.Status = cryptodb.StatusOrdered
		if tradeManOrder.OrderKind == cryptodb.KindEntry {
			db.Where("order_kind = ?", cryptodb.KindMarketStopLoss).Where("plan_id = ?", tradeManOrder.PlanID).Find(&marketStopLossOrder)
			if o.StopLoss.Equal(marketStopLossOrder.Price) {
				marketStopLossOrder.Status = cryptodb.StatusOrdered
			} else {
				// TODO: think about raising hell or just quitely logging this situation. This should never ever happen.
			}
		}
		tx := db.Begin()
		err = db.SaveOrder(&tradeManOrder)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = db.SaveOrder(&marketStopLossOrder)
		if err != nil {
			tx.Rollback()
			return err
		}

		var tmLog cryptodb.Log
		tmLog.PlanID = tradeManOrder.PlanID
		tmLog.Source = cryptodb.SourceExchange

		var ExchangeOrderID string
		if o.OrderID != "" {
			ExchangeOrderID = o.OrderID
		} else if o.StopOrderID != "" {
			ExchangeOrderID = o.StopOrderID
		}
		var stopLossSetMsg string
		if marketStopLossOrder.Status == cryptodb.StatusOrdered {
			stopLossSetMsg = "and"
		} else {
			stopLossSetMsg = "but DID NOT"
		}

        tmLog.Text = fmt.Sprintf("Exchange accepted order %d %s set stoploss as OrderID: %s, at %s.", tradeManOrder.ID, stopLossSetMsg, ExchangeOrderID, o.CreatedAt.Format("2006-01-02 15:04:05.000"))
		err = db.CreateLog(&tmLog)
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}
	return nil
}

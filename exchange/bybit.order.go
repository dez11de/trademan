package exchange

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dez11de/cryptodb"
)

// Also sends Market StopLoss
func (e *Exchange) SendEntry(plan cryptodb.Plan, pair cryptodb.Pair, marketStopLoss *cryptodb.Order, entry *cryptodb.Order) (err error) {
	entryParams := make(RequestParameters)
	var endPoint string
	var response OrderResponseRest
	var result OrderResponseRest

	if entry.SystemOrderID == "" {
		endPoint = "/private/linear/order/create"
	} else {
		endPoint = "/private/linear/order/replace"
	}

	if plan.Direction == cryptodb.Long {
		entryParams["side"] = "Buy"
	} else {
		entryParams["side"] = "Sell"
	}
	entryParams["symbol"] = pair.Name
	entryParams["order_type"] = "Limit"
	entryParams["qty"] = entry.Size.InexactFloat64()
	entryParams["close_on_trigger"] = false
	entryParams["reduce_only"] = false
	entryParams["time_in_force"] = "GoodTillCancel"
	entryParams["stop_loss"] = marketStopLoss.Price.InexactFloat64()
	entryParams["sl_trigger_by"] = "LastPrice"

	entryParams["price"] = entry.Price.InexactFloat64()

	_, responseBuffer, err := e.SignedRequest(http.MethodPost, endPoint, entryParams, &result)
	if err != nil {
		return err
	}

	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		return errors.New(result.ReturnMessage)
	}

	err = json.Unmarshal(responseBuffer, &response)
	if err != nil {
		return err
	}

	entry.SystemOrderID = response.Result.OrderID
	entry.Status.Scan(response.Result.OrderStatus)
	return marketStopLoss.Status.Scan(response.Result.OrderStatus)
}

// TODO: basicly the same as function above. Rewrite into 1 function.
func (e *Exchange) SendLimitOrder(plan cryptodb.Plan, pair cryptodb.Pair, entry cryptodb.Order, limitOrder *cryptodb.Order) (err error) {
	tpParams := make(RequestParameters)
	var endPoint string
	var result OrderResponseRest
	var response OrderResponseRest

	if entry.SystemOrderID == "" {
		endPoint = "/private/linear/stop-order/create"
	} else {
		endPoint = "/private/linear/stop-order/replace"
	}

	if plan.Direction == cryptodb.Long {
		tpParams["side"] = "Sell"
	} else {
		tpParams["side"] = "Buy"
	}

	tpParams["symbol"] = pair.Name
	tpParams["order_type"] = "Limit"
	tpParams["qty"] = limitOrder.Size.InexactFloat64()
	tpParams["trigger_by"] = "LastPrice"
	tpParams["close_on_trigger"] = false
	tpParams["reduce_only"] = true
	tpParams["time_in_force"] = "GoodTillCancel"

	tpParams["price"] = limitOrder.Price.InexactFloat64()
	tpParams["stop_px"] = limitOrder.TriggerPrice.InexactFloat64()
	tpParams["base_price"] = entry.Price.InexactFloat64()

	_, responseBuffer, err := e.SignedRequest(http.MethodPost, endPoint, tpParams, &result)
	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		return errors.New(result.ReturnMessage)
	}

	err = json.Unmarshal(responseBuffer, &response)
	if err != nil {
		return err
	}

	limitOrder.SystemOrderID = response.Result.StopOrderID
	return limitOrder.Status.Scan(response.Result.OrderStatus)
}

package exchange

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dez11de/cryptodb"
)

func (e *Exchange) SendOrder(plan cryptodb.Plan, pair cryptodb.Pair, entry, order *cryptodb.Order) (err error) {
	var endPoint string
	params := make(RequestParameters)
	var response, result OrderResponseRest

	switch {
	case entry.SystemOrderID == "":
		endPoint = "/private/linear/order/create"
		params["symbol"] = pair.Name
		if plan.Direction == cryptodb.Long {
			params["side"] = "Buy"
		} else {
			params["side"] = "Sell"
		}
		params["order_type"] = "Limit"
		params["qty"] = entry.Size.InexactFloat64()
		params["price"] = entry.Price.InexactFloat64()
		params["close_on_trigger"] = false
		params["reduce_only"] = false
		params["time_in_force"] = "GoodTillCancel"
		params["sl_trigger_by"] = "LastPrice"

		params["stop_loss"] = order.Price.InexactFloat64()

	case order.SystemOrderID == "":
		endPoint = "/private/linear/stop-order/create"
		params["symbol"] = pair.Name
		if plan.Direction == cryptodb.Long {
			params["side"] = "Sell"
		} else {
			params["side"] = "Buy"
		}
		params["order_type"] = "Limit"
		params["qty"] = order.Size.InexactFloat64()
		params["price"] = order.Price.InexactFloat64()
		params["close_on_trigger"] = false
		params["time_in_force"] = "GoodTillCancel"

		params["stop_px"] = order.TriggerPrice.InexactFloat64()
		params["base_price"] = entry.Price.InexactFloat64()
		params["trigger_by"] = "LastPrice"
		params["reduce_only"] = true

	case entry.SystemOrderID != "" && order.OrderKind == cryptodb.MarketStopLoss:
		endPoint = "/private/linear/order/replace"
		params["order_id"] = entry.SystemOrderID
		params["p_r_price"] = entry.Price.InexactFloat64()
		params["stop_px"] = order.TriggerPrice.InexactFloat64()

	case order.SystemOrderID != "":
		endPoint = "/private/linear/stop-order/replace"
		params["order_id"] = order.SystemOrderID
		params["p_r_price"] = order.Price.InexactFloat64()
	}

	_, responseBuffer, err := e.SignedRequest(http.MethodPost, endPoint, params, &result)
	if err != nil {
		return err
	}

	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		return errors.New(fmt.Sprintf("(%d) %s", result.ReturnCode, result.ReturnMessage))
	}

    // TODO: this might be redundant, can just use result
	err = json.Unmarshal(responseBuffer, &response)
	if err != nil {
		return err
	}

	if entry.SystemOrderID == "" {
		entry.SystemOrderID = response.Result.OrderID
		err = entry.Status.Scan(response.Result.OrderStatus)
	} else {
		order.SystemOrderID = response.Result.StopOrderID
		err = order.Status.Scan(response.Result.OrderStatus)
	}

	return err
}

func (e *Exchange) CancelOrder(symbol, SystemOrderID string) (err error) {
	var endPoint string
	params := make(RequestParameters)
	var response, result OrderResponseRest
    
		endPoint = "/private/linear/stop-order/cancel"
		params["symbol"] = symbol
        params["stop_order_id"] = SystemOrderID
	_, responseBuffer, err := e.SignedRequest(http.MethodPost, endPoint, params, &result)
	if err != nil {
		return err
	}

	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		return errors.New(fmt.Sprintf("(%d) %s", result.ReturnCode, result.ReturnMessage))
	}

    // TODO: is this really necessary?
	err = json.Unmarshal(responseBuffer, &response)
	if err != nil {
		return err
	}

    return nil
}

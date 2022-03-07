package exchange

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/dez11de/cryptodb"
)

// Also sends Market StopLoss
func (e *Exchange) SendEntry(plan cryptodb.Plan, pair cryptodb.Pair, marketStopLoss *cryptodb.Order, entry *cryptodb.Order) (err error) {
	var result OrderResponseRest
	entryParams := make(RequestParameters)

	if plan.Direction == cryptodb.Long {
		entryParams["side"] = "Buy"
	} else {
		entryParams["side"] = "Sell"
	}
	entryParams["order_link_id"] = entry.LinkOrderID
	entryParams["symbol"] = pair.Name
	entryParams["order_type"] = "Limit"
	entryParams["qty"] = entry.Size.InexactFloat64()
	entryParams["price"] = entry.Price.InexactFloat64()
	entryParams["close_on_trigger"] = false
	entryParams["reduce_only"] = false
	entryParams["time_in_force"] = "GoodTillCancel"
	entryParams["stop_loss"] = marketStopLoss.Price.InexactFloat64()
	entryParams["sl_trigger_by"] = "LastPrice"

	var response OrderResponseRest
	_, responseBuffer, err := e.SignedRequest(http.MethodPost, "/private/linear/order/create", entryParams, &result)
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

	return nil
}

func (e *Exchange) SendLimitStopLoss(plan cryptodb.Plan, pair cryptodb.Pair, marketStopLoss cryptodb.Order, limitStopLoss *cryptodb.Order, entry cryptodb.Order) (err error) {
	var result OrderResponseRest
	lslParams := make(RequestParameters)

	if plan.Direction == cryptodb.Long {
		lslParams["side"] = "Sell"
	} else {
		lslParams["side"] = "Buy"
	}
	lslParams["order_link_id"] = limitStopLoss.LinkOrderID
	lslParams["symbol"] = pair.Name
	lslParams["order_type"] = "Limit"
	lslParams["qty"] = limitStopLoss.Size.InexactFloat64()

	lslParams["trigger_by"] = "LastPrice"
	lslParams["price"] = limitStopLoss.Price.InexactFloat64()
	lslParams["stop_px"] = limitStopLoss.TriggerPrice.InexactFloat64()
	lslParams["base_price"] = entry.Price.InexactFloat64()

	lslParams["close_on_trigger"] = false
	lslParams["reduce_only"] = true
	lslParams["time_in_force"] = "GoodTillCancel"

	var response OrderResponseRest
	_, responseBuffer, err := e.SignedRequest(http.MethodPost, "/private/linear/stop-order/create", lslParams, &result)
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

	limitStopLoss.SystemOrderID = response.Result.StopOrderID
	limitStopLoss.Status.Scan(response.Result.OrderStatus)

	return nil
}

func (e *Exchange) SendTakeProfit(plan cryptodb.Plan, pair cryptodb.Pair, marketStopLoss, entry cryptodb.Order, takeProfit *cryptodb.Order) (err error) {
	var result OrderResponseRest
	tpParams := make(RequestParameters)

	if plan.Direction == cryptodb.Long {
		tpParams["side"] = "Sell"
	} else {
		tpParams["side"] = "Buy"
	}
	tpParams["order_link_id"] = takeProfit.LinkOrderID
	tpParams["symbol"] = pair.Name
	tpParams["order_type"] = "Limit"
	tpParams["qty"] = takeProfit.Size.InexactFloat64()
	tpParams["trigger_by"] = "LastPrice"
	tpParams["price"] = takeProfit.Price.InexactFloat64()
	tpParams["stop_px"] = takeProfit.TriggerPrice.InexactFloat64() 
	tpParams["base_price"] = entry.Price.InexactFloat64()

	tpParams["close_on_trigger"] = false
	tpParams["reduce_only"] = true
	tpParams["time_in_force"] = "GoodTillCancel"

	var response OrderResponseRest
	_, responseBuffer, err := e.SignedRequest(http.MethodPost, "/private/linear/stop-order/create", tpParams, &result)
	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		log.Printf("[SendTakeProfit] error: %s", result.ReturnMessage)
		return errors.New(result.ReturnMessage)
	}

	err = json.Unmarshal(responseBuffer, &response)
	if err != nil {
		return err
	}

	takeProfit.SystemOrderID = response.Result.StopOrderID
	takeProfit.Status.Scan(response.Result.OrderStatus)

	return nil
}

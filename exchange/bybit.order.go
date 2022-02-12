package exchange

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dez11de/cryptodb"
)

func (e *Exchange) PlaceOrders(p cryptodb.Plan, activePair cryptodb.Pair, o []cryptodb.Order) (err error) {
	err = e.placeEntry(p, activePair, o[cryptodb.KindMarketStopLoss], o[cryptodb.KindLimitStopLoss], o[cryptodb.KindEntry])
	if err != nil {
		return err
	}

	for i := 3; i < 3+cryptodb.MaxTakeProfits; i++ {
		if !o[i].Price.IsZero() {
			err = e.placeTakeProfit(p, activePair, o[cryptodb.KindEntry], o[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Exchange) placeEntry(plan cryptodb.Plan, pair cryptodb.Pair, marketStopLoss, limitStopLoss, entry cryptodb.Order) (err error) {
	var result OrderResponse
	var orderResult OrderResponse
	entryParams := make(RequestParameters)

	// Set entry and HardStopLoss
	entryParams["order_link_id"] = entry.ExchangeOrderID
	entryParams["symbol"] = pair.Name
	if plan.Direction == cryptodb.Long {
		entryParams["side"] = "Buy"
	} else {
		entryParams["side"] = "Sell"
	}
	entryParams["order_type"] = "Limit"
	entryParams["qty"] = entry.Size.InexactFloat64()    // TODO: check if this is has been properly 'steprounded'
	entryParams["price"] = entry.Price.InexactFloat64() // TODO: check if this is has been properly 'steprounded'
	entryParams["close_on_trigger"] = false             // TODO: figure out what exactly this means
	entryParams["reduce_only"] = false                  // TODO: figure out what exactly this means
	entryParams["time_in_force"] = "GoodTillCancel"
	entryParams["stop_loss"] = marketStopLoss.Price.InexactFloat64()

	_, resp, err := e.SignedRequest(http.MethodPost, "/private/linear/order/create", entryParams, &result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &orderResult)
	if err != nil {
		return err
	}
	if result.ReturnMessage != "OK" {
		return errors.New(result.ReturnMessage)
	}

	// Set LimitStopLoss
	sslParams := make(RequestParameters)
	sslParams["order_link_id"] = limitStopLoss.ExchangeOrderID
	sslParams["symbol"] = pair.Name
	if plan.Direction == cryptodb.Long {
		sslParams["side"] = "Sell"
	} else {
		sslParams["side"] = "Buy"
	}
	sslParams["order_type"] = "Limit"
	sslParams["qty"] = limitStopLoss.Size.InexactFloat64()    // TODO: check if this is has been properly 'steprounded'
	sslParams["price"] = limitStopLoss.Price.InexactFloat64() // TODO: check if this is has been properly 'steprounded'
	sslParams["base_price"] = entry.Price.InexactFloat64()    // TODO: check if this is has been properly 'steprounded'
	sslParams["stop_px"] = entry.Price.InexactFloat64()       // TODO: check if this is has been properly 'steprounded'
	sslParams["close_on_trigger"] = false                     // TODO: i have no idea
	sslParams["trigger_by"] = "LastPrice"                     // TODO: i have some sort of idea, api documentation seems incorrect
	sslParams["reduce_only"] = true                           // TODO: i have no idea
	sslParams["time_in_force"] = "GoodTillCancel"

	_, resp, err = e.SignedRequest(http.MethodPost, "/private/linear/stop-order/create", sslParams, &result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return err
	}
	if result.ReturnMessage != "OK" {
		return errors.New(result.ReturnMessage)
	}

	return nil
}

func (e *Exchange) placeTakeProfit(plan cryptodb.Plan, pair cryptodb.Pair, entry, takeProfit cryptodb.Order) (err error) {
	var result OrderResponse
	var orderResult OrderResponse
	takeProfitParams := make(RequestParameters)

	takeProfitParams["order_link_id"] = takeProfit.ExchangeOrderID
	takeProfitParams["symbol"] = pair.Name
	if plan.Direction == cryptodb.Long {
		takeProfitParams["side"] = "Sell"
	} else {
		takeProfitParams["side"] = "Buy"
	}
	takeProfitParams["order_type"] = "Limit"
	takeProfitParams["qty"] = takeProfit.Size.InexactFloat64()    // TODO: check if this is has been properly 'steprounded'
	takeProfitParams["price"] = takeProfit.Price.InexactFloat64() // TODO: check if this is has been properly 'steprounded'
	takeProfitParams["stop_px"] = entry.Price.InexactFloat64()    // TODO: check if this is has been properly 'steprounded'
	takeProfitParams["base_price"] = entry.Price.InexactFloat64() // TODO: check if this is has been properly 'steprounded'
	takeProfitParams["trigger_by"] = "LastPrice"                  // TODO: i have some sort of idea, api documentation seems incorrect
	takeProfitParams["close_on_trigger"] = false                  // TODO: figure out what exactly this means
	takeProfitParams["reduce_only"] = true                        // TODO: figure out what exactly this means
	takeProfitParams["time_in_force"] = "GoodTillCancel"

	_, resp, err := e.SignedRequest(http.MethodPost, "/private/linear/stop-order/create", takeProfitParams, &result)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, &orderResult)
	if err != nil {
		return err
	}
	if result.ReturnMessage != "OK" {
		return errors.New(result.ReturnMessage)
	}
	return nil
}

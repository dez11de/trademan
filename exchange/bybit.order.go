package exchange

import (
	"errors"
	"log"
	"net/http"

	"github.com/dez11de/cryptodb"
)

func (e *Exchange) PlaceOrders(p cryptodb.Plan, activePair cryptodb.Pair, ticker Ticker, o []cryptodb.Order) (err error) {
	err = e.placeEntry(p, activePair, o[cryptodb.MarketStopLoss], o[cryptodb.LimitStopLoss], o[cryptodb.Entry])
	if err != nil {
		return err
	}

	err = e.placeLimitStopLoss(p, activePair, ticker, o[cryptodb.MarketStopLoss], o[cryptodb.LimitStopLoss], o[cryptodb.Entry])
	if err != nil {
		return err
	}

	for i := 3; i < 3+cryptodb.MaxTakeProfits; i++ {
		if !o[i].Price.IsZero() {
			err = e.placeTakeProfit(p, activePair, ticker, o[cryptodb.MarketStopLoss], o[cryptodb.Entry], o[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Exchange) placeEntry(plan cryptodb.Plan, pair cryptodb.Pair, marketStopLoss, limitStopLoss, entry cryptodb.Order) (err error) {
	var result OrderResponse
	entryParams := make(RequestParameters)

	switch plan.Direction {
	case cryptodb.Long:
		if plan.Leverage.GreaterThan(pair.Leverage.Buy) {
			e.SetLeverage(pair.Name, plan.Leverage, pair.Leverage.Sell)
		}
	case cryptodb.Short:
		if plan.Leverage.GreaterThan(pair.Leverage.Sell) {
			e.SetLeverage(pair.Name, pair.Leverage.Buy, plan.Leverage)
		}
	}

	// Set entry and marketStopLoss
	entryParams["order_link_id"] = entry.ExchangeOrderID
	entryParams["symbol"] = pair.Name
	if plan.Direction == cryptodb.Long {
		entryParams["side"] = "Buy"
	} else {
		entryParams["side"] = "Sell"
	}
	entryParams["order_type"] = "Limit"
	entryParams["qty"] = entry.Size.InexactFloat64()
	entryParams["price"] = entry.Price.InexactFloat64()
	entryParams["close_on_trigger"] = false
	entryParams["reduce_only"] = false
	entryParams["time_in_force"] = "GoodTillCancel"
	entryParams["stop_loss"] = marketStopLoss.Price.InexactFloat64()

	_, _, err = e.SignedRequest(http.MethodPost, "/private/linear/order/create", entryParams, &result)

	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		log.Printf("Order not accepted. ret_code: %d, ext_code: %s", result.ReturnCode, result.ExtendedCode)
		log.Printf("------------\nResult: %v\n------------", result)
		// log.Printf("Response: %s", string(response))
		return errors.New(result.ExtendedCode)
	}

	return nil
}

func (e *Exchange) placeLimitStopLoss(plan cryptodb.Plan, pair cryptodb.Pair, ticker Ticker, marketStopLoss, limitStopLoss, entry cryptodb.Order) (err error) {
	var result OrderResponse
	// Set LimitStopLoss
	lslParams := make(RequestParameters)
	lslParams["order_link_id"] = limitStopLoss.ExchangeOrderID
	lslParams["symbol"] = pair.Name
	lslParams["order_type"] = "Limit"
	lslParams["qty"] = limitStopLoss.Size.InexactFloat64()
	if plan.Direction == cryptodb.Long {
		lslParams["side"] = "Sell"
	} else {
		lslParams["side"] = "Buy"
	}

	lslParams["price"] = limitStopLoss.Price.InexactFloat64()
	lslParams["stop_px"] = entry.Price.InexactFloat64()
	lslParams["base_price"] = ticker.LastPrice.InexactFloat64()

	lslParams["trigger_by"] = "LastPrice"
	lslParams["close_on_trigger"] = false
	lslParams["reduce_only"] = true
	lslParams["time_in_force"] = "GoodTillCancel"

	_, _, err = e.SignedRequest(http.MethodPost, "/private/linear/stop-order/create", lslParams, &result)
	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		log.Printf("Order not accepted. ret_code: %d, ext_code: %s", result.ReturnCode, result.ExtendedCode)
		log.Printf("------------\nResult: %v\n------------", result)
		// log.Printf("Response: %s", string(response))
		return errors.New(result.ExtendedCode)
	}

	return nil
}

func (e *Exchange) placeTakeProfit(plan cryptodb.Plan, pair cryptodb.Pair, ticker Ticker, marketStopLoss, entry, takeProfit cryptodb.Order) (err error) {
	var result OrderResponse
	tpParams := make(RequestParameters)

	tpParams["order_link_id"] = takeProfit.ExchangeOrderID
	tpParams["symbol"] = pair.Name
	tpParams["order_type"] = "Limit"
	tpParams["qty"] = takeProfit.Size.InexactFloat64()
	if plan.Direction == cryptodb.Long {
		tpParams["side"] = "Sell"
	} else {
		tpParams["side"] = "Buy"
	}
	tpParams["price"] = takeProfit.Price.InexactFloat64()
	tpParams["stop_px"] = marketStopLoss.Price.InexactFloat64()
	tpParams["base_price"] = ticker.LastPrice.InexactFloat64()

	tpParams["trigger_by"] = "LastPrice"
	tpParams["close_on_trigger"] = false
	tpParams["reduce_only"] = true
	tpParams["time_in_force"] = "GoodTillCancel"

	_, _, err = e.SignedRequest(http.MethodPost, "/private/linear/stop-order/create", tpParams, &result)
	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		log.Printf("Order not accepted. ret_code: %d, ext_code: %s", result.ReturnCode, result.ExtendedCode)
		log.Printf("------------\nResult: %v\n------------", result)
		// log.Printf("Response: %s", string(response))
		return errors.New(result.ExtendedCode)
	}

	return nil
}

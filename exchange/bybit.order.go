package exchange

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dez11de/cryptodb"
)

func (e *Exchange) PlaceOrders(p cryptodb.Plan, activePair cryptodb.Pair, o []cryptodb.Order) (err error) {
	err = e.placeEntry(p, activePair, o[cryptodb.KindMarketStopLoss], o[cryptodb.KindLimitStopLoss], o[cryptodb.KindEntry])
	return err
}

func (e *Exchange) placeEntry(plan cryptodb.Plan, pair cryptodb.Pair, marketStopLoss, limitStopLoss, entry cryptodb.Order) (err error) {
	var result OrderResponse
	entryParams := make(RequestParameters)

	// Set entry and HardStopLoss
	entryParams["order_link_id"] = entry.ExchangeOrderID
	entryParams["symbol"] = pair.Name
	if plan.Direction == cryptodb.DirectionLong {
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

	fullURL, resp, err := e.SignedRequest(http.MethodPost, "/private/linear/order/create", entryParams, &result)
	json.Unmarshal(resp, &result)
	if err != nil || result.ReturnMessage != "OK" {
		log.Printf("Entry not accepted: %s", err)
		log.Printf("URL: %s", fullURL)
		log.Printf("Response: %v", result)
		return err // TODO: if no error but ReturnMessage not "OK" return that
	}

	entry.Status = cryptodb.StatusOrdered
	marketStopLoss.Status = cryptodb.StatusOrdered

	// Set LimitStopLoss
	sslParams := make(RequestParameters)
	sslParams["order_link_id"] = limitStopLoss.ExchangeOrderID
	sslParams["symbol"] = pair.Name
	if plan.Direction == cryptodb.DirectionLong {
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
	sslParams["reduce_only"] = false                          // TODO: i have no idea
	sslParams["time_in_force"] = "GoodTillCancel"

	fullURL, resp, err = e.SignedRequest(http.MethodPost, "/private/linear/stop-order/create", sslParams, &result)
	json.Unmarshal(resp, &result)
	if err != nil || result.ReturnMessage != "OK" {
		log.Printf("Soft stoploss not accepted: %s", err)
		log.Printf("URL: %s", fullURL)
		log.Printf("Response: %s", string(resp))
		return err // TODO: if no error but ReturnMessage not "OK" return that
	}
	limitStopLoss.Status = cryptodb.StatusOrdered

	return nil
}

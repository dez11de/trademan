package main

import (
	"log"
	"net/http"
)

func (bb *ByBit) PlaceOrders(p Plan, activePair Pair, o Orders) (err error) {
	log.Print("Placing orders...")

	var result OrderResponse

	params := map[string]interface{}{}
	// TODO: create new enum Side, not to be confused with direction?
	params["side"] = "Buy"
	params["symbol"] = activePair.Pair
	// TODO: create new enum for order_type, not the be confused with OrderType
	params["order_type"] = "Limit"
	params["qty"] = o[typeEntry].Size.InexactFloat64()
	params["price"] = o[typeEntry].Price
	params["time_in_force"] = "GoodTillCancel"
	// TODO: figure out what exactly this means
	params["close_on_trigger"] = false
	// TODO: figure out what exactly this means
	params["reduce_only"] = false
	// TODO: there is a better solution for this, either wait for RoundStep to get integrated or use the fork https://github.com/bart613/decimal
	log.Printf("Order size step: %s", activePair.OrderSize.Step.String())
	log.Printf("Order size adjusted for step: %s", o[typeEntry].Size.Div(activePair.OrderSize.Step).Floor().Mul(activePair.OrderSize.Step).String())
	params["stop_loss"] = o[typeHardStopLoss].Price.InexactFloat64()

	fullUrl, response, err := bb.SignedRequest(http.MethodPost, "/private/linear/order/create", params, &result)
	log.Printf("Full URL: %s", fullUrl)
	log.Printf("Result: %v", result)
	log.Printf("Response: %v", string(response))
	return err
}

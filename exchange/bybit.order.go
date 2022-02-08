package exchange

import (
	"log"
	"net/http"

	"github.com/dez11de/cryptodb"
)

func (e *Exchange) PlaceOrders(p cryptodb.Plan, activePair cryptodb.Pair, o []cryptodb.Order) (err error) {
	log.Print("Placing orders...")

	var result OrderResponse

	params := map[string]interface{}{}
	params["side"] = "Buy"
	params["symbol"] = activePair.Name
	params["order_type"] = "Limit"
	params["qty"] = o[cryptodb.KindEntry].Size.InexactFloat64()
	params["price"] = o[cryptodb.KindEntry].Price
	params["time_in_force"] = "GoodTillCancel"
	// TODO: figure out what exactly this means
	params["close_on_trigger"] = false
	// TODO: figure out what exactly this means
	params["reduce_only"] = false
	// TODO: there is a better solution for this, either wait for RoundStep to get integrated or use the fork https://github.com/bart613/decimal
	log.Printf("Order size step: %s", activePair.Order.Step.String())
	log.Printf("Order size adjusted for step: %s", o[cryptodb.KindEntry].Size.Div(activePair.Order.Step).Floor().Mul(activePair.Order.Step).String())
	params["stop_loss"] = o[cryptodb.KindHardStopLoss].Price.InexactFloat64()

	fullUrl, response, err := e.SignedRequest(http.MethodPost, "/private/linear/order/create", params, &result)
	log.Printf("Full URL: %s", fullUrl)
	log.Printf("Result: %v", result)
	log.Printf("Response: %v", string(response))
	return err
}

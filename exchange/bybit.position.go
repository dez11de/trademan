package exchange

import (
	"errors"
	"net/http"

	"github.com/bart613/decimal"
)

/*
func (e *Exchange) GetPosition(pair string) (position Position, err error) {
	log.Printf("Getting position for symbol: %s", pair)
	var result PositionResponse
	params := make(RequestParameters)

	params["symbol"] = pair

	fullUrl, response, err := e.SignedRequest(http.MethodPost, "/private/linear/position/list", params, &result)

    log.Printf("Request to: %s\nReponse: %+v", fullUrl, response)

	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		return position, errors.New(result.ExtendedCode)
	}

    if len (result.Positions) == 0 {
        log.Printf("Position not found?")
        return position, nil
    }

	return result.Positions[0], nil
}
*/

// TODO: implement switch Cross/Isolated switch function?

// Assumes position leverage is already set to isolated
func (e *Exchange) SetLeverage(pair string, buyLeverage, sellLeverage decimal.Decimal) (err error) {

	var result PositionResponse
	params := make(RequestParameters)

	params["symbol"] = pair
	params["buy_leverage"] = buyLeverage.String()
	params["sell_leverage"] = sellLeverage.String()

	_, _, err = e.SignedRequest(http.MethodPost, "/private/linear/position/set-leverage", params, &result)
	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		return errors.New(result.ExtendedCode)
	}

	return nil
}

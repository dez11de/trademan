package exchange

import (
	"github.com/dez11de/cryptodb"
)

func (e *Exchange) GetPairs() (pairs []cryptodb.Pair, err error) {

	var pr PairResponse
	_, err = e.PublicRequest("GET", "/v2/public/symbols", nil, &pr)

	return pr.Pairs, err
}

// Assumes position leverage is already set to isolated
// func (e *Exchange) SetLeverage(pair string, buyLeverage, sellLeverage decimal.Decimal) (err error) {
// 
// 	var result PositionResponse
// 	params := make(RequestParameters)
// 
// 	params["symbol"] = pair
// 	params["buy_leverage"] = buyLeverage.String()
// 	params["sell_leverage"] = sellLeverage.String()
// 
// 	_, _, err = e.SignedRequest(http.MethodPost, "/private/linear/position/set-leverage", params, &result)
// 	if result.ReturnCode != 0 || result.ExtendedCode != "" {
// 		return errors.New(result.ExtendedCode)
// 	}
// 
// 	return nil
// }

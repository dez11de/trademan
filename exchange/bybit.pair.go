package exchange

import (
	"errors"
	"log"
	"net/http"

	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
)

func (e *Exchange) GetPairs() (pairs []cryptodb.Pair, err error) {

	var pr PairResponse
	_, err = e.PublicRequest("GET", "/v2/public/symbols", nil, &pr)

	return pr.Pairs, err
}

func (e *Exchange) SendLeverage(n string, q string, l, s decimal.Decimal) (err error) {
	var url string
	levParams := make(RequestParameters)
	var result LeverageResponse

	switch q {
	case "USDT":
		levParams["symbol"] = n
		levParams["buy_leverage"] = l.InexactFloat64()
		levParams["sell_leverage"] = s.InexactFloat64()
		url = "/private/linear/position/set-leverage"
	default:
		log.Printf("Setting margin not supported for this pair.")
	}

	_, _, err = e.SignedRequest(http.MethodPost, url, levParams, &result)
	log.Printf("SendLeverage result: %+v", result)
	if result.ReturnCode != 0 || result.ExtendedCode != "" {
		return errors.New(result.ReturnMessage)
	}

	return nil
}

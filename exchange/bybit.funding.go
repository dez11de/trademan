package exchange

import (
	"log"

	"github.com/dez11de/cryptodb"
)

func (e *Exchange) GetRecentFunding(p string) (f cryptodb.Funding, err error) {
	var fr FundingResponse
	params := make(RequestParameters)

    params["symbol"] = p
	_, err = e.PrivateRequest("GET", "/private/linear/funding/prev-funding", params, &fr)
    log.Printf("GetRecentFunding: %v", fr)
		return fr.Funding, err
}

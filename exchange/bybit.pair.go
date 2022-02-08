package exchange

import (
	"github.com/dez11de/cryptodb"
)

func (e *Exchange) GetPairs() (pairs []cryptodb.Pair, err error) {
	var pr PairResponse
	_, err = e.PublicRequest("GET", "/v2/public/symbols", nil, &pr)

    return pr.Pairs, err
}

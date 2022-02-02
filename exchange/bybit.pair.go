package exchange

import (
	"github.com/dez11de/cryptodb"
)

func (b *ByBit) GetPairs() (pairs []cryptodb.Pair, err error) {
	var pr PairResponse
	_, err = b.PublicRequest("GET", "/v2/public/symbols", nil, &pr)
	if err != nil {
		return nil, err
	}
    /*
	for _, pair := range pr.Results {
        pairs = append(pairs, pair)
	}
	return pairs, nil
    */ 

    return pr.Results, err
}

package exchange

import (
	"log"

	"github.com/dez11de/cryptodb"
)

func (b *ByBit) GetPairs() (pairs map[string]cryptodb.Pair, err error) {
    pairs = make(map[string]cryptodb.Pair)
	var pr PairResponse
	_, err = b.PublicRequest("GET", "/v2/public/symbols", nil, &pr)
	if err != nil {
		return nil, err
	}
	for _, pair := range pr.Results {
        log.Printf("Storing pair %v", pair)
		pairs[pair.Pair] = pair
	}
	return pairs, nil
}

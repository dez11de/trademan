package exchange

import "github.com/dez11de/cryptodb"

func (b *ByBit) GetPairs() map[string]cryptodb.Pair {
	var pr pairResponse
	s := make(map[string]cryptodb.Pair)
	b.PublicRequest("GET", "/v2/public/symbols", nil, &pr)
	for _, pair := range pr.Results {
		s[pair.Pair] = pair
	}
	return s
}

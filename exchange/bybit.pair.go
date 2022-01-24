package exchange

import "github.com/dez11de/cryptodb"

func (b *ByBit) GetPairs() map[string]cryptoDB.Pair {
	var pr pairResponse
	s := make(map[string]cryptoDB.Pair)
	b.PublicRequest("GET", "/v2/public/symbols", nil, &pr)
	for _, pair := range pr.Results {
		s[pair.Pair] = pair
	}
	return s
}

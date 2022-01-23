package main

func (b *ByBit) GetPairs() map[string]Pair {
	var pr pairResponse
	s := make(map[string]Pair)
	b.PublicRequest("GET", "/v2/public/symbols", nil, &pr)
	for _, pair := range pr.Results {
		s[pair.Pair] = pair
	}
	return s
}

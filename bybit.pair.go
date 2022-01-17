package main

type pairResponse struct {
	ReturnCode       int    `json:"ret_code"`
	ReturnMessage    string `json:"ret_msg"`
	ExtendedCode     string `json:"ext_code"`
	Results          []Pair `json:"result"`
	ExtendedInfo     string `json:"ext_info"`
	ServerTime       string `json:"time_now,string"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMS int    `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}

func (b *ByBit) GetPairs() map[string]Pair {
	var pr pairResponse
	s := make(map[string]Pair)
	b.PublicRequest("GET", "/v2/public/symbols", nil, &pr)
	for _, pair := range pr.Results {
		s[pair.Pair] = pair
	}
	return s
}

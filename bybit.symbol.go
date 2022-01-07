package main

type SymbolResponse struct {
	ReturnCode       int      `json:"ret_code"`
	ReturnMessage    string   `json:"ret_msg"`
	ExtendedCode     string   `json:"ext_code"`
	Results          []symbol `json:"result"`
	ExtendedInfo     string   `json:"ext_info"`
	ServerTime       string   `json:"time_now,string"`
	RateLimitStatus  int      `json:"rate_limit_status"`
	RateLimitResetMS int      `json:"rate_limit_reset_ms"`
	RateLimit        int      `json:"rate_limit"`
}

func (b *ByBit) GetSymbols() []symbol {
	var sr SymbolResponse
	b.PublicRequest("GET", "/v2/public/symbols", nil, &sr)
	return sr.Results
}

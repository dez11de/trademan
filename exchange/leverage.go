package exchange

// TODO: bybit considers this part of position, according to documentation?
type LeverageResponse struct {
	ReturnCode       int    `json:"ret_code"`
	ReturnMessage    string `json:"ret_msg"`
	ExtendedCode     string `json:"ext_code"`
	Result           string  `json:"result"`
	ExtendedInfo     string `json:"ext_info"`
	ServerTime       string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMS int    `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}


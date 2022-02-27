package exchange

import "github.com/bart613/decimal"

type Position struct {
	Pair          string          `json:"symbol"`
	Side          string          `json:"side"`
	Size          decimal.Decimal `json:"size"`
	Price         decimal.Decimal `json:"entry_price,string"` // Average Entry Price
	OCCClosingFee decimal.Decimal `json:"occ_closing_fee,string"`
}

type PositionResponse struct {
	ReturnCode       int64      `json:"ret_code"`
	ReturnMessage    string     `json:"ret_msg"`
	ExtendedCode     string     `json:"ext_code"`
	Positions        []Position `json:"result"`
	ExtendedInfo     string     `json:"ext_info"`
	ServerTime       string     `json:"time_now,string"`
	RateLimitStatus  int64      `json:"rate_limit_status"`
	RateLimitResetMS int64      `json:"rate_limit_reset_ms"`
	RateLimit        int64      `json:"rate_limit"`
}

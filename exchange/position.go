package exchange

import "github.com/bart613/decimal"

type Position struct {
	Pair          string          `json:"symbol"`
	Side          string          `json:"side"`
	Size          decimal.Decimal `json:"size"`
	EntryPrice    decimal.Decimal `json:"entry_price,string"` // Average Entry Price
}

type PositionResponse struct {
	ReturnCode       int        `json:"ret_code"`
	ReturnMessage    string     `json:"ret_msg"`
	ExtendedCode     string     `json:"ext_code"`
	Positions        []Position `json:"result"`
	ExtendedInfo     string     `json:"ext_info"`
	ServerTime       string     `json:"time_now,string"`
	RateLimitStatus  int        `json:"rate_limit_status"`
	RateLimitResetMS int        `json:"rate_limit_reset_ms"`
	RateLimit        int        `json:"rate_limit"`
}

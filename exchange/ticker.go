package exchange

import (
	"github.com/bart613/decimal"
)

type Ticker struct {
	Symbol    string          `json:"symbol"`
	LastPrice decimal.Decimal `json:"last_price,string"`
}

type TickerResponse struct {
	ReturnCode       int      `json:"ret_code"`
	ReturnMessage    string   `json:"ret_msg"`
	ExtendedCode     string   `json:"ext_code"`
	Results          []Ticker `json:"result"`
	ExtendedInfo     string   `json:"ext_info"`
	ServerTime       string   `json:"time_now"`
	RateLimitStatus  int      `json:"rate_limit_status"`
	RateLimitResetMS int      `json:"rate_limit_reset_ms"`
	RateLimit        int      `json:"rate_limit"`
}

package exchange

import (
	"github.com/dez11de/cryptodb"
)

type OrderPage struct {
	CurrentPage int              `json:"current_page"`
	LastPage    int              `json:"last_page"`
	Orders      []cryptodb.Order `json:"data"`
}

type OrderResponse struct {
	ReturnCode       int       `json:"ret_code"`
	ReturnMessage    string    `json:"ret_msg"`
	ExtendedCode     string    `json:"ext_code"`
	Results          OrderPage `json:"result"`
	ExtendedInfo     string    `json:"ext_info"`
	ServerTime       string    `json:"time_now,string"`
	RateLimitStatus  int       `json:"rate_limit_status"`
	RateLimitResetMS int       `json:"rate_limit_reset_ms"`
	RateLimit        int       `json:"rate_limit"`
}

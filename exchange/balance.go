package exchange

import (
	"github.com/dez11de/cryptodb"
)

type WalletResponse struct {
	ReturnCode       int64                       `json:"ret_code"`
	ReturnMessage    string                      `json:"ret_msg"`
	ExtendedCode     string                      `json:"ext_code"`
	Results          map[string]cryptodb.Balance `json:"result"`
	ExtendedInfo     string                      `json:"ext_info"`
	ServerTime       string                      `json:"time_now"`
	RateLimitStatus  int64                       `json:"rate_limit_status"`
	RateLimitResetMS int64                       `json:"rate_limit_reset_ms"`
	RateLimit        int64                       `json:"rate_limit"`
}

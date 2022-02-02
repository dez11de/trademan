package exchange

import (
	"github.com/dez11de/cryptodb"
)

type WalletResponse struct {
	ReturnCode       int                         `json:"ret_code"`
	ReturnMessage    string                      `json:"ret_msg"`
	ExtendedCode     string                      `json:"ext_code"`
	Results          map[string]cryptodb.Balance `json:"result"`
	ExtendedInfo     string                      `json:"ext_info"`
	ServerTime       string                      `json:"time_now,string"`
	RateLimitStatus  int                         `json:"rate_limit_status"`
	RateLimitResetMS int                         `json:"rate_limit_reset_ms"`
	RateLimit        int                         `json:"rate_limit"`
}

func (b *ByBit) GetCurrentWallet() (map[string]cryptodb.Balance, error) {
	var wr WalletResponse
	params := make(map[string]interface{})
	b.PrivateRequest("GET", "/v2/private/wallet/balance", params, &wr)

	wallet := make(map[string]cryptodb.Balance)
	for s, b := range wr.Results {
		b.Symbol = s
		wallet[s] = b
	}
	return wallet, nil // TODO: return an actual error on all the things that can go wrong
}

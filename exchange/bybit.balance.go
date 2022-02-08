package exchange

import (
	"github.com/dez11de/cryptodb"
)

// TODO: return an actual error on all the things that can go wrong
func (b *ByBit) GetCurrentWallet() (map[string]cryptodb.Balance, error) {
	var wr WalletResponse
	params := make(map[string]interface{})
	b.PrivateRequest("GET", "/v2/private/wallet/balance", params, &wr)

	wallet := make(map[string]cryptodb.Balance)
	for s, b := range wr.Results {
		b.Symbol = s
		wallet[s] = b
	}
	return wallet, nil
}
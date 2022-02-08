package exchange

import (
	"github.com/dez11de/cryptodb"
)

// TODO: return an actual error on all the things that can go wrong
func (b *Exchange) GetCurrentWallet() (balances map[string]cryptodb.Balance, err error) {
	var wr WalletResponse
	params := make(map[string]interface{})
	_, err = b.PrivateRequest("GET", "/v2/private/wallet/balance", params, &wr)
    if err != nil {
        return balances, err
    }
	balances = make(map[string]cryptodb.Balance)
	for s, b := range wr.Results {
		b.Symbol = s
		balances[s] = b
	}
	return balances, nil
}

package main

import (
	"fmt"
	"log"

    "github.com/dez11de/cryptodb"
    "github.com/dez11de/exchange"
)

func main() {
	db := cryptoDB.NewDB()
	err := db.Connect()
	if err != nil {
		log.Printf("Error connecting to database: %s", err)
	}

	exchange := exchange.NewByBit("https://api-testnet.bybit.com",
		"wss://stream-testnet.bybit.com/realtime_private",
		"rLot58Xxaj4Kdb3pog",
		"0a3GihYe3CfFkLbYsE41wWoNTtofwY2WPkwi",
		false)
	if exchange == nil {
		fmt.Println("Error creating ByBit object")
	}

	err = exchange.Connect()
	if err != nil {
		fmt.Println(err)
	}

	// Might be needed to reload stuff in database
	currentPairs := exchange.GetPairs()
	for _, p := range currentPairs {
		db.AddPair(p)
	}
	currentBalances, _ := exchange.GetCurrentWallet()
	for _, b := range currentBalances {
		db.AddWallet(b)
	}

	for _, pair := range db.PairCache {
		fmt.Printf("Pair: %s, ID: %d, Quote Currency: %s\n", pair.Pair, pair.PairID, pair.QuoteCurrency)
	}

	for currency, balance := range db.WalletCache {
		fmt.Printf("Currency: %s, Available balance: %s\n", currency, balance.Available.String())
	}
}

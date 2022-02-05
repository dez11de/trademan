package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

var db *cryptodb.Database

func main() {
	// TODO: read and pass config from env/commandline/configfile
    var err error
	db, err = cryptodb.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	// TODO: read these from config or envvar
	exchange := exchange.NewByBit("https://api-testnet.bybit.com",
		"wss://stream-testnet.bybit.com/realtime_private",
		"rLot58Xxaj4Kdb3pog",
		"0a3GihYe3CfFkLbYsE41wWoNTtofwY2WPkwi",
		false)
	// TODO: read and pass config from env/commandline/configfile
	// TODO: return err if connecting failed
	if exchange == nil {
		log.Fatalf("Error creating ByBit object")
	}

	err = exchange.Connect()
	if err != nil {
		fmt.Println(err)
	}

	// TODO: add exchange.Ping() to keep connection alive
	pingExchangeTicker := time.NewTicker(1 * time.Minute)
	refreshWalletTicker := time.NewTicker(2 * time.Hour)
	refreshPairsTicker := time.NewTicker(24 * time.Hour)
	quit := make(chan struct{})

	go HandleRequests()

	for {
		select {
		case <-refreshWalletTicker.C:
			log.Print("Getting balance from exchange")
			currentBalances, err := exchange.GetCurrentWallet()
			if err != nil {
				log.Printf("error getting current wallet from exchange %v", err)
			} else {
				for _, b := range currentBalances {
					err = db.CreateBalance(&b)
					if err != nil {
						log.Printf("error writing wallet to database %v", err)
					}
				}
			}
		case <-refreshPairsTicker.C:
			currentPairs, err := exchange.GetPairs()
			if err != nil {
				log.Printf("error getting current pairs %v", err)
			} else {
				for _, p := range currentPairs {
					err = db.CreatePair(&p)
					if err != nil {
						log.Printf("error writing pair to database %v", err)
					}
				}
			}
		case <-pingExchangeTicker.C:
			log.Print("Pinging exchange...")
			exchange.Ping()
		case <-quit:
			refreshWalletTicker.Stop()
			return
		}
	}
}

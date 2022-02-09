package main

import (
	"log"
	"time"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

var db *cryptodb.Database
var e *exchange.Exchange

func main() {
	var trademanCfg trademanConfig
	err := readConfig(&trademanCfg)
	if err != nil {
		log.Fatalf("Error reading configuration: %s", err)
	}

	db, err = cryptodb.Connect(trademanCfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

    e, err = exchange.Connect(trademanCfg.Exchange)
	if err != nil {
		log.Fatalf("Error connecting to exchange: %s", err)
	}

    if trademanCfg.Database.ResetTables {
        // database tables will be recreated by cryptodb
        // fill the `pairs` table again
        exchangePairs, err := e.GetPairs()
        if err != nil {
            log.Fatalf("unable to reload pairs from exchange: %s", err)
        }
        for _, p := range exchangePairs {
            db.CrupdatePair(&p) // eventhough tables have just been reset 
        }
        exchangeWallet, err := e.GetCurrentWallet()
        if err != nil {
            log.Fatalf("unable to get current wallet from exchange: %s", err)
        }
        for _, b := range exchangeWallet {
            db.CreateBalance(&b)
        }
    }

	pingExchangeTicker := time.NewTicker(1 * time.Minute)
	refreshWalletTicker := time.NewTicker(2 * time.Hour)
	refreshPairsTicker := time.NewTicker(24 * time.Hour)
	quit := make(chan struct{})

	// TODO: what if it can't open the port?
	go HandleRequests(trademanCfg.RESTServer)
	// go exchange.ProcessMessages()

	for {
		select {
		case <-refreshWalletTicker.C:
			log.Print("Getting balance from exchange")
			currentBalances, err := e.GetCurrentWallet()
			if err != nil {
				log.Printf("error getting current wallet from exchange %v", err)
			} else {
				for _, b := range currentBalances {
					err = db.CreateBalance(&b)
					if err != nil {
						log.Printf("error writing balance to database %v", err)
					}
				}
			}
		case <-refreshPairsTicker.C:
			currentPairs, err := e.GetPairs()
			if err != nil {
				log.Printf("error getting current pairs %v", err)
			} else {
				for _, p := range currentPairs {
					err = db.CrupdatePair(&p)
					if err != nil {
						log.Printf("error writing pair to database %v", err)
					}
				}
			}
		case <-pingExchangeTicker.C:
			log.Print("Pinging exchange...")
			e.Ping()
		case <-quit:
			refreshWalletTicker.Stop()
			return
		}
	}
}

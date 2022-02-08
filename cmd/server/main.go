package main

import (
	"log"
	"time"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

var db *cryptodb.Database

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

	exchange, err := exchange.Connect(trademanCfg.Exchange)
	if err != nil {
		log.Fatalf("Error connecting to exchange: %s", err)
	}

    if trademanCfg.Database.ResetTables {
        // database tables will be recreated by cryptodb
        // fill the `pairs` table again
        exchangePairs, err := exchange.GetPairs()
        if err != nil {
            log.Fatalf("unable to reload pairs from exchange: %s", err)
        }
        for _, p := range exchangePairs {
            db.CreatePair(&p)
            log.Printf("Pair %s assigned ID %d", p.Name, p.ID)
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
			currentBalances, err := exchange.GetCurrentWallet()
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
			currentPairs, err := exchange.GetPairs()
			if err != nil {
				log.Printf("error getting current pairs %v", err)
			} else {
				for _, p := range currentPairs {
                    // TODO: should check if it should update or create
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

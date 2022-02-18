package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

var db *cryptodb.Database
var e *exchange.Exchange

func main() {
	// TODO: respond to os.Signal messages in the exepected way. See https://pace.dev/blog/2020/02/17/repond-to-ctrl-c-interrupt-signals-gracefully-with-context-in-golang-by-mat-ryer.html
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

	if trademanCfg.Database.CreateTables {
		db.CreateTables()
		exchangePairs, err := e.GetPairs()
		if err != nil {
			log.Fatalf("unable to reload pairs from exchange: %s", err)
		}
		for _, p := range exchangePairs {
			db.CrupdatePair(&p) // eventhough tables have just been reset
			time.Sleep(1543 * time.Millisecond)
			p.Leverage.Buy = decimal.NewFromInt(1)
			p.Leverage.Sell = decimal.NewFromInt(1)
			e.SetLeverage(p.Name, p.Leverage.Buy, p.Leverage.Sell)
			db.CrupdatePair(&p)
		}

		exchangeWallet, err := e.GetCurrentWallet()
		if err != nil {
			log.Fatalf("unable to get current wallet from exchange: %s", err)
		}
		for _, b := range exchangeWallet {
			db.CreateBalance(&b)
		}
	}

	if trademanCfg.Database.TruncTables {
        log.Printf("Truncating most tables.")
		db.TruncTables()
	}

	err = e.Subscribe("position")
	if err != nil {
		fmt.Println(err)
	}
	positionUpdate := make(chan exchange.Position)

	err = e.Subscribe("execution")
	if err != nil {
		fmt.Println(err)
	}
	executionUpdate := make(chan exchange.Execution)

	err = e.Subscribe("order")
	if err != nil {
		fmt.Println(err)
	}
	err = e.Subscribe("stop_order")
	if err != nil {
		fmt.Println(err)
	}
	orderUpdate := make(chan exchange.Order)

	// TODO: also subscribe to wallet socket?

	errorMessage := make(chan error)

	pingExchangeTicker := time.NewTicker(1 * time.Minute)
	refreshWalletTicker := time.NewTicker(2 * time.Hour)
	refreshPairsTicker := time.NewTicker(24 * time.Hour)
	quit := make(chan struct{})

	// TODO: what if it can't open the port?
	go HandleRequests(trademanCfg.RESTServer)

	go e.ProcessMessages(positionUpdate, executionUpdate, orderUpdate, errorMessage)

	for {
		select {
		case <-refreshWalletTicker.C:
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
			e.Ping()

		case <-quit:
			refreshWalletTicker.Stop()
			return

		case p := <-positionUpdate:
			err := processPosition(p)
			if err != nil {
				errorMessage <- err
			}

		case e := <-executionUpdate:
			err := processExecution(e)
			if err != nil {
				errorMessage <- err
			}

		case o := <-orderUpdate:
			err := processOrder(o)
			if err != nil {
				errorMessage <- err
			}

		case e := <-errorMessage:
			// TODO: should probably log these to file
			log.Printf("Received error: %s", e)
		}
	}
}

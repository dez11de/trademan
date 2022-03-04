package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bart613/decimal"
	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
)

var db *cryptodb.Database
var e *exchange.Exchange

var play = make(chan struct{})
var pause = make(chan struct{})
var wg sync.WaitGroup

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
			time.Sleep(1543 * time.Millisecond)
			p.Leverage.Long = decimal.NewFromInt(1)
			p.Leverage.Short = decimal.NewFromInt(1)
			e.SendLeverage(p.Name, p.QuoteCurrency, p.Leverage.Long.RoundStep(p.Leverage.Step, false), p.Leverage.Short.RoundStep(p.Leverage.Step, false))
			db.Create(&p)
		}

		exchangeWallet, err := e.GetCurrentWallet()
		if err != nil {
			log.Fatalf("unable to get current wallet from exchange: %s", err)
		}
		for _, b := range exchangeWallet {
			db.Create(&b)
		}
	}

	if trademanCfg.Database.TruncTables {
		log.Printf("Truncating most tables.")
		db.TruncTables()

		exchangeWallet, err := e.GetCurrentWallet()
		if err != nil {
			log.Fatalf("unable to get current wallet from exchange: %s", err)
		}
		for _, b := range exchangeWallet {
			db.Create(&b)
		}
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

	pingExchangeTicker := time.NewTicker(1 * time.Minute)
	refreshWalletTicker := time.NewTicker(2 * time.Hour)
	refreshPairsTicker := time.NewTicker(24 * time.Hour)
	quit := make(chan struct{})

	// TODO: what if it can't open the port?
	go HandleRequests(trademanCfg.RESTServer)

	go e.ProcessMessages(positionUpdate, executionUpdate, orderUpdate)

	wg.Add(1)

	for {
		select {
		case <-refreshWalletTicker.C:
			currentBalances, err := e.GetCurrentWallet()
			if err != nil {
				log.Printf("error getting current wallet from exchange %v", err)
			} else {
				for _, b := range currentBalances {
					result := db.Create(&b)
					if result.Error != nil {
						log.Printf("error writing balance to database %v", result.Error)
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
			// log.Print("Im still pinging...")
			e.Ping()

		case <-quit:
			refreshWalletTicker.Stop()
			return

		case p := <-positionUpdate:
			err := processPosition(p)
			if err != nil {
				log.Printf("Error occured processing position: %s", err)
			}

		case e := <-executionUpdate:
			err := processExecution(e)
			if err != nil {
				log.Printf("Error occured processing execution: %s", err)
			}

		case o := <-orderUpdate:
			err := processOrder(o)
			if err != nil {
				log.Printf("Error occured processing order: %s", err)
			}
		case <-pause:
			select {
			case <-play:
			}
		}

	}
}

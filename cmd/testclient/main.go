package main

import (
	"fmt"
	"log"

	"github.com/dez11de/cryptodb"
	"github.com/dez11de/exchange"
	"github.com/shopspring/decimal"
)

func main() {
	db, err := cryptodb.Connect()
	if err != nil {
		log.Printf("Error connecting to database: %s", err)
	}

	// TODO: read and pass config from env/commandline/configfile
	// TODO: return err if connecting failed
	exchange := exchange.NewByBit("https://api-testnet.bybit.com",
		"wss://stream-testnet.bybit.com/realtime_private",
		"rLot58Xxaj4Kdb3pog",
		"0a3GihYe3CfFkLbYsE41wWoNTtofwY2WPkwi",
		false)
	if exchange == nil {
		log.Fatalf("Error creating ByBit object")
	}

	err = exchange.Connect()
	if err != nil {
		fmt.Println(err)
	}

	db.RecreateTables()

	exchangePairs, err := exchange.GetPairs()
	for _, exchangePair := range exchangePairs {
		db.CreatePair(&exchangePair)
	}

	currentPair, err := db.GetPairByName("ADAUSDT")
	if err != nil {
		log.Printf("error GetPairByName: %s", err)
	}
	fmt.Printf("Pair name: %s\nPair status: %s\n\n", currentPair.Name, currentPair.Status)

    currentPlan := cryptodb.Plan{}
    currentPlan.PairID = currentPair.ID
    currentPlan.Notes = "This is just a test plan."
    db.CreatePlan(&currentPlan)
    fmt.Printf("Plan #: %d\nPair name: %s\nPlan status: %s\n", currentPlan.ID, currentPair.Name, currentPlan.Notes)

    currentPlan.Notes = currentPlan.Notes + " And this is an update."
    db.SavePlan(&currentPlan)
    fmt.Printf("Plan #: %d\nPair name: %s\nPlan status: %s\n", currentPlan.ID, currentPair.Name, currentPlan.Notes)

    orders := cryptodb.NewOrders(currentPlan.ID)
    orders[cryptodb.TypeHardStopLoss].Price = decimal.NewFromFloat(1.0003)
    orders[cryptodb.TypeEntry].Price = decimal.NewFromFloat(1.0263)
    orders[cryptodb.TypeTakeProfit+0].Price = decimal.NewFromFloat(1.202)
    orders[cryptodb.TypeTakeProfit+1].Price = decimal.NewFromFloat(1.3166)
    db.CreateOrders(&orders)
}

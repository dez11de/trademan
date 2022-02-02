package main

import (
	"fmt"
	"log"
	"time"

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
		db.SavePair(&exchangePair)
	}

    var orders  []cryptodb.Order
    var dbPlan cryptodb.Plan

	dbPair, err := db.GetPairByName("ADAUSDT")
    dbPlan.PairID = dbPair.ID

	orders = append(orders, cryptodb.Order{OrderType: cryptodb.TypeHardStopLoss, Price: decimal.NewFromFloat(1.0003)})
	orders = append(orders, cryptodb.Order{OrderType: cryptodb.TypeSoftStopLoss, Price: decimal.NewFromFloat(1.0004)})
	orders = append(orders, cryptodb.Order{OrderType: cryptodb.TypeEntry, Price: decimal.NewFromFloat(1.0263)})
	orders = append(orders, cryptodb.Order{OrderType: cryptodb.TypeTakeProfit, Price: decimal.NewFromFloat(1.202)})
	orders = append(orders, cryptodb.Order{OrderType: cryptodb.TypeTakeProfit, Price: decimal.NewFromFloat(1.3166)})

	db.SavePlan(&dbPlan)
    log.Printf("Saved plan: %v", dbPlan)
    for _, o := range orders {
        o.PlanID = dbPlan.ID
        db.SaveOrder(&o)
    }

	testPlan, err := db.GetPlan(dbPlan.ID)
    testPair, err := db.GetPair(dbPlan.PairID)
    testOrders, err := db.GetOrders(dbPlan.ID)

	fmt.Printf("Initial plan #%d\nPair name: %s\nEntry price: %s\n\n", testPlan.ID, testPair.Name, testOrders[cryptodb.TypeEntry].Price.StringFixed(4))

	orders[cryptodb.TypeEntry].Price = decimal.NewFromFloat(1.2)
	db.SavePlan(&testPlan)

    time.Sleep(2* time.Second)

	newPlan, err := db.GetPlan(dbPlan.ID)
    orders, err = db.GetOrders(dbPlan.ID)
	fmt.Printf("Updated plan #%d\nPair name: %s\nEntry price: %s\n\n", newPlan.ID, testPair.Name, orders[cryptodb.TypeEntry].Price.StringFixed(4))
}

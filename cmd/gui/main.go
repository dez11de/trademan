package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

    "github.com/dez11de/cryptodb"
)

func main() {
	db := cryptodb.NewDB()
	err := db.Connect()
	if err != nil {
		log.Printf("Error connecting to database: %s", err)
	}

	for _, pair := range db.PairCache {
		fmt.Printf("Pair: %s, ID: %d, Quote Currency: %s\n", pair.Pair, pair.PairID, pair.QuoteCurrency)
	}

	for currency, balance := range db.WalletCache {
		fmt.Printf("Currency: %s, Available balance: %s\n", currency, balance.Available.String())
	}

	app := app.NewWithID("bbtrader")
	// TODO: come up with a better name that doesn't focus on ByBit cause you never know
	mainWindow := app.NewWindow("ByBit Trade Manager")
	// TODO: would it hurt if these are global?
	mainContent := makeMainContent(db)
	// TODO: store and restore size and position settings
	mainWindow.Resize(fyne.Size{Width: 800, Height: 600})
	mainWindow.CenterOnScreen()
	mainWindow.SetContent(mainContent)
	mainWindow.ShowAndRun()
}

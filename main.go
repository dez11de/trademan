package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	db := NewDB()
	err := db.Connect()
	if err != nil {
		log.Printf("Error connecting to database: %s", err)
	}

	exchange := NewByBit("https://api-testnet.bybit.com",
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

	if err := tea.NewProgram(newMainUIModel(db), tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

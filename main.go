package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	db := NewDB()
	db.Connect()

	if err := tea.NewProgram(newModel(db), tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

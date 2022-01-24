package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
    "github.com/dez11de/cryptodb"
)

func MakePerformanceContainer(db *cryptodb.Database) *fyne.Container {
    // TODO: also show winrate and average rrr over time, maybe as a toggle?
	dailyPerformance := canvas.NewText(fmt.Sprintf("Daily: %s%%", db.GetPerformance(1*24*time.Hour).StringFixed(2)), nil)
	dailyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	dailyPerformance.TextSize = 10
	weeklyPerformance := canvas.NewText(fmt.Sprintf("Weekly: %s%%", db.GetPerformance(7*24*time.Hour).StringFixed(2)), nil)
	weeklyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	weeklyPerformance.TextSize = 10
	monthlyPerformance := canvas.NewText(fmt.Sprintf("Monthly: %s%%", db.GetPerformance(30*24*time.Hour).StringFixed(2)), nil)
	monthlyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	monthlyPerformance.TextSize = 10
	quarterlyPerformance := canvas.NewText(fmt.Sprintf("Quarterly: %s%%", db.GetPerformance(91*24*time.Hour).StringFixed(2)), nil)
	quarterlyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	quarterlyPerformance.TextSize = 10
	yearlyPerformance := canvas.NewText(fmt.Sprintf("Yearly: %s%%", db.GetPerformance(365*24*time.Hour).StringFixed(2)), nil)
	yearlyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	yearlyPerformance.TextSize = 10

	performancePane := container.New(layout.NewGridLayout(5), dailyPerformance, weeklyPerformance, monthlyPerformance, quarterlyPerformance, yearlyPerformance)

	return performancePane
}

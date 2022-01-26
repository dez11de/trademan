package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func MakePerformanceContainer() *fyne.Container {
    // TODO: also show winrate and average rrr over time, maybe as a toggle?
	dailyPerformance := canvas.NewText(fmt.Sprintf("Daily: %s%%", "1.1"), nil)
	dailyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	dailyPerformance.TextSize = 10
	weeklyPerformance := canvas.NewText(fmt.Sprintf("Weekly: %s%%", "1.2"), nil)
	weeklyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	weeklyPerformance.TextSize = 10
	monthlyPerformance := canvas.NewText(fmt.Sprintf("Monthly: %s%%", "1.3"), nil)
	monthlyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	monthlyPerformance.TextSize = 10
	quarterlyPerformance := canvas.NewText(fmt.Sprintf("Quarterly: %s%%", "1.4"), nil)
	quarterlyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	quarterlyPerformance.TextSize = 10
	yearlyPerformance := canvas.NewText(fmt.Sprintf("Yearly: %s%%", "1.5"), nil)
	yearlyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	yearlyPerformance.TextSize = 10

	performancePane := container.New(layout.NewGridLayout(5), dailyPerformance, weeklyPerformance, monthlyPerformance, quarterlyPerformance, yearlyPerformance)

	return performancePane
}

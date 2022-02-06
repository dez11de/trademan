package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func getPerformance(symbol string, period time.Duration) float64 {
    /*
	jsonRequest, err := json.Marshal(cryptodb.Performance{Symbol: symbol,
		Since: time.Now().Add(-period),
	})
	if err != nil {
		log.Printf("error marshalling request %v", err)
	}

	client := http.Client{Timeout: time.Second * 2}
	// TODO: make host configurable in env/param/file
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8888/performance", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Printf("error requesting: %v", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := client.Do(req)
	if err != nil {
		log.Printf("error doing request %v", err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("error reading response %v", err)
	}
	defer response.Body.Close()
	if err != nil {
		log.Printf("Error reading response.Body")
	}
	var performance cryptodb.Performance
	err = json.Unmarshal(body, &performance)
	if err != nil {
		log.Printf("Error unmarshalling performance %v", err)
	}
    */

	return 0.0 
}

func MakePerformanceContainer() *fyne.Container {
	// TODO: also show winrate and average rrr over time, maybe as a toggle?
	dailyPerformance := canvas.NewText(fmt.Sprintf("Daily: %.1f%%", getPerformance("USDT", 1*24*time.Hour)), nil)
	dailyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	dailyPerformance.TextSize = 10
	weeklyPerformance := canvas.NewText(fmt.Sprintf("Weekly: %.1f%%", getPerformance("USDT", 7*24*time.Hour)), nil)
	weeklyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	weeklyPerformance.TextSize = 10
	monthlyPerformance := canvas.NewText(fmt.Sprintf("Monthly: %.1f%%", getPerformance("USDT", 30*24*time.Hour)), nil)
	monthlyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	monthlyPerformance.TextSize = 10
	quarterlyPerformance := canvas.NewText(fmt.Sprintf("Quarterly: %.1f%%", getPerformance("USDT", 91*24*time.Hour)), nil)
	quarterlyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	quarterlyPerformance.TextSize = 10
	yearlyPerformance := canvas.NewText(fmt.Sprintf("Yearly: %.1f%%", getPerformance("USDT", 365*24*time.Hour)), nil)
	yearlyPerformance.TextStyle = fyne.TextStyle{Monospace: true}
	yearlyPerformance.TextSize = 10

	performancePane := container.NewHBox(layout.NewSpacer(), dailyPerformance, weeklyPerformance, monthlyPerformance, quarterlyPerformance, yearlyPerformance, layout.NewSpacer())

	return performancePane
}

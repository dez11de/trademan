package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (pf *planForm) makeStatContainer() *fyne.Container {
	// TODO: make distinction between start RRR and evolved RRR
	startRewardRiskRatioLabel := widget.NewLabel("Start RRR: ")
	startRewardRiskRatioValue := widget.NewLabel(fmt.Sprintf("%.1f", 0.0))
	evolvedRewardRiskRatioLabel := widget.NewLabel("Current RRR: ")
	evolvedRewardRiskRatioValue := widget.NewLabel(fmt.Sprintf("%.1f", 0.0))

	currentPnLLabel := widget.NewLabel("PnL: ")
	currentPnLValue := widget.NewLabel(fmt.Sprintf("%s%%", act.plan.Profit.StringFixed(1))) // TODO: should be relative to entrySize.
	// TODO: figure out what this even means, see CryptoCred.
	breakEvenLabel := widget.NewLabel("B/E: ")
	breakEvenValue := widget.NewLabel(fmt.Sprintf("%.0f%%", 0.0))
	container := container.NewHBox(
		layout.NewSpacer(),
		startRewardRiskRatioLabel, startRewardRiskRatioValue,
		evolvedRewardRiskRatioLabel, evolvedRewardRiskRatioValue,
		currentPnLLabel, currentPnLValue,
		breakEvenLabel, breakEvenValue,
		layout.NewSpacer())

	return container
}


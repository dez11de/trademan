package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func makeMainContent(db *Database, bb *ByBit) *fyne.Container {
	performanceContainer := MakePerformanceContainer(db)
	planListSplitContainer := MakePlanListSplit(db, bb)

	mainContent := container.New(
		layout.NewBorderLayout(
			nil, performanceContainer, nil, nil),
		performanceContainer, planListSplitContainer)

	planListSplitContainer.Refresh()
	mainContent.Refresh()
	return mainContent
}

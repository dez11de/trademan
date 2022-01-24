package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
    "github.com/dez11de/cryptodb"
)

func makeMainContent(db *cryptodb.Database) *fyne.Container {
	performanceContainer := MakePerformanceContainer(db)
	planListSplitContainer := MakePlanListSplit(db)

	mainContent := container.New(
		layout.NewBorderLayout(
			nil, performanceContainer, nil, nil),
		performanceContainer, planListSplitContainer)

	planListSplitContainer.Refresh()
	mainContent.Refresh()
	return mainContent
}

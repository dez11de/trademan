package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func makeMainContent() *fyne.Container {
	performanceContainer := MakePerformanceContainer()
	planListSplitContainer := MakePlanListSplit()

	mainContent := container.New(
		layout.NewBorderLayout(
			nil, performanceContainer, nil, nil),
		performanceContainer, planListSplitContainer)

	planListSplitContainer.Refresh()
	mainContent.Refresh()
	return mainContent
}

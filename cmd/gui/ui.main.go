package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func makeMainContent() *fyne.Container {
    mainContent := container.NewBorder(nil, MakePerformanceContainer(), nil, nil, MakePlanListSplit())

	return mainContent
}

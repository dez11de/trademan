package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func makeMainContent() *fyne.Container {
	return container.NewBorder(nil, MakePerformanceContainer(), nil, nil, MakePlanListSplit())
}

func makeMainMenu() *fyne.MainMenu {
	showArchivedItem := fyne.NewMenuItem("Show Archived", nil)
	showArchivedItem.Action = func() {
		ui.showArchived = !ui.showArchived
		showArchivedItem.Checked = ui.showArchived
		tm.plans, _ = getPlans()
		ui.planListSplit.Trailing = ui.noPlanSelectedContainer
		ui.planListSplit.Refresh()
		ui.planList.Refresh()
	}

	showArchivedItem.Checked = ui.showArchived

	viewMenu := fyne.NewMenu("View", showArchivedItem)

	return fyne.NewMainMenu(viewMenu)
}

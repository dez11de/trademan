package main

import (
	"fyne.io/fyne/v2"
)

func SetupMainWindow(db *Database, w fyne.Window) {
	w.Resize(fyne.Size{Width: 800, Height: 600})
	w.CenterOnScreen()
	SetupPlanList(db, w)
}

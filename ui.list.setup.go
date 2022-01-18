package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

func SetupPlanList(d *Database, w fyne.Window) {
	plans, _ := d.GetPlans()
	list := widget.NewList(
		func() int {
			return len(plans)
		},
		func() fyne.CanvasObject {
			pairText := canvas.NewText("PAIRCUR", colornames.White)
			pairText.TextStyle = fyne.TextStyle{Bold: true}
			statusText := canvas.NewText("STATUS", colornames.White)
			directionText := canvas.NewText("Long", colornames.Green)
			directionText.Alignment = fyne.TextAlignTrailing
			return container.NewVBox(
				pairText,
				container.NewHBox(statusText, directionText),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0].(*canvas.Text).Text = d.GetPairString(plans[i].PairID)

			var statusColor color.Color
			// TODO: give all posible statuses a different color
			switch plans[i].Status {
			case Planned:
				statusColor = colornames.Blue
			case Ordered:
				statusColor = colornames.Green
			case Filled:
				statusColor = colornames.Purple
			default:
				statusColor = colornames.White
			}
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Text = plans[i].Status.String()
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Color = statusColor

			var directionColor color.Color
			switch plans[i].Side {
			case Long:
				directionColor = colornames.Green
			case Short:
				directionColor = colornames.Red
			}
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Text).Color = directionColor
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*canvas.Text).Text = plans[i].Side.String()

		})

	selectedPlanLabel := widget.NewLabel("Select a position from the list.")
	formBox := container.NewHSplit(list, selectedPlanLabel)
	formBox.SetOffset(0.15)
	w.SetContent(formBox)

	list.OnSelected = func(id widget.ListItemID) {
		log.Printf("Loading form with plan %v", plans[id])
		formBox.Trailing = makePlanForm(d, plans[id])
		formBox.Refresh()
	}
}

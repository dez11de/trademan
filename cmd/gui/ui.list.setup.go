package main

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
    "github.com/dez11de/cryptodb"
    "github.com/dez11de/exchange"
)

type planListUI struct {
	Plans []cryptoDB.Plan
	List  *widget.List

	addPlanAction    widget.ToolbarItem
	removePlanAction widget.ToolbarItem
	actionBar        *widget.Toolbar
}

func MakePlanListSplit(d *cryptoDB.Database, bb *exchange.ByBit) *container.Split {
	planList := &planListUI{}
	planList.Plans, _ = d.GetPlans()
	planList.List = widget.NewList(
		func() int {
			return len(planList.Plans)
		},
		func() fyne.CanvasObject {
			// TODO: change this to widget.RichText
			pairText := canvas.NewText("PAIRCUR", colornames.White)
			pairText.TextStyle = fyne.TextStyle{Bold: true}
			statusText := canvas.NewText("STATUS", colornames.White)
			directionText := canvas.NewText("Long", colornames.Green)
			return container.NewVBox(
				container.NewHBox(pairText, layout.NewSpacer(), directionText),
				container.New(layout.NewCenterLayout(), statusText),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*canvas.Text).Text = d.GetPairString(planList.Plans[i].PairID)

			// TODO: use theme colors
			var directionColor color.Color
			switch planList.Plans[i].Side {
			case cryptoDB.SideLong:
				directionColor = colornames.Green
			case cryptoDB.SideShort:
				directionColor = colornames.Red
			}
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*canvas.Text).Color = directionColor
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*canvas.Text).Text = planList.Plans[i].Side.String()

			var statusColor color.Color
			// TODO: give all posible statuses a different color
			switch planList.Plans[i].Status {
			case cryptoDB.StatusPlanned:
				statusColor = colornames.Blue
			case cryptoDB.StatusOrdered:
				statusColor = colornames.Green
			case cryptoDB.StatusFilled:
				statusColor = colornames.Purple
			default:
				statusColor = colornames.White
			}
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Text = planList.Plans[i].Status.String()
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Color = statusColor
		})

	selectPlanLabel := container.New(layout.NewCenterLayout(), canvas.NewText("Select a plan from the list, or press + to make a new plan.", nil))

	listAndButtons := container.NewWithoutLayout(widget.NewLabel("nothing to see here"))
	planListSplit := container.NewHSplit(listAndButtons, container.NewMax(selectPlanLabel))
	planListSplit.SetOffset(0.20)

	planList.addPlanAction = widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		log.Print("Add button pressed")
		f := NewForm(d, bb)
		f.FillForm(cryptoDB.Plan{})
		planListSplit.Trailing = f.form

		planListSplit.Refresh()
	})
	planList.removePlanAction = widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {
		log.Print("Remove button pressed")
	})

	planList.actionBar = widget.NewToolbar(widget.NewToolbarSpacer(), planList.addPlanAction, planList.removePlanAction)
	planList.actionBar.Refresh()
	listAndButtons = container.New(layout.NewBorderLayout(nil, planList.actionBar, nil, nil), container.NewMax(planList.List), planList.actionBar)
	planListSplit.Leading = listAndButtons
	planList.List.Refresh()
	listAndButtons.Refresh()
	planListSplit.Leading.Refresh()
	planListSplit.Refresh()

	planList.List.OnSelected = func(id widget.ListItemID) {
		f := NewForm(d, bb)
		f.FillForm(planList.Plans[id])
		planListSplit.Trailing = f.form
		planListSplit.Refresh()
	}

	return planListSplit
}

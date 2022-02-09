package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
	"golang.org/x/image/colornames"
)

type planListUI struct {
	Plans []cryptodb.Plan
    PairCache []cryptodb.Pair
	List  *widget.List

	addPlanAction    widget.ToolbarItem
	removePlanAction widget.ToolbarItem
	actionBar        *widget.Toolbar
}

func MakePlanListSplit() *container.Split {
	planList := &planListUI{}
    planList.Plans = make([]cryptodb.Plan,0)
    var err error
    planList.Plans, err = getPlans()
    if err != nil {
        dialog.ShowError(err, mainWindow)
    }
    planList.PairCache, err = getPairs()
    if err != nil {
        dialog.ShowError(err, mainWindow)
    }
	planList.List = widget.NewList(
		func() int {
			return len(planList.Plans)
		},
		func() fyne.CanvasObject {
			// TODO: change this to widget.RichText
			pairText := canvas.NewText("Pair", theme.ForegroundColor())
			pairText.TextStyle = fyne.TextStyle{Bold: true}
			statusText := canvas.NewText("Status", colornames.White)
			directionText := canvas.NewText("Direction", colornames.Green)
			return container.NewVBox(
				container.NewHBox(pairText, layout.NewSpacer(), directionText),
				container.New(layout.NewCenterLayout(), statusText),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*canvas.Text).Text = planList.PairCache[int64(planList.Plans[i].PairID-1)].Name

			// TODO: use theme colors
			var directionColor color.Color
            var directionName string
			switch planList.Plans[i].Direction {
			case cryptodb.DirectionLong:
                directionName = "Long"
				directionColor = theme.PrimaryColorNamed("Green")
			case cryptodb.DirectionShort:
                directionName = "Short"
				directionColor = theme.PrimaryColorNamed("Red")
			}
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*canvas.Text).Text = directionName 
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*canvas.Text).Color = directionColor 

			var statusColor color.Color
            var statusName string
			// TODO: give all posible statuses a different color
			switch planList.Plans[i].Status {
			case cryptodb.StatusPlanned:
                statusName = "Planned"
				statusColor = theme.PrimaryColorNamed("Blue")
			case cryptodb.StatusOrdered:
                statusName = "Ordered"
				statusColor = theme.PrimaryColorNamed("Green")
			case cryptodb.StatusFilled:
                statusName = "Filled"
				statusColor = theme.PrimaryColorNamed("Magenta")
			default:
				statusColor = colornames.White
			}
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Text = statusName
            o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Color = statusColor
		})

	selectPlanLabel := container.New(layout.NewCenterLayout(), canvas.NewText("Select a plan from the list, or press + to make a new plan.", nil))

	listAndButtons := container.NewWithoutLayout(widget.NewLabel("Error loading plans.\nCheck internet connection."))
	planListSplit := container.NewHSplit(listAndButtons, container.NewMax(selectPlanLabel))
	planListSplit.SetOffset(0.20)

	planList.addPlanAction = widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		f := NewForm()
		f.FillForm(cryptodb.Plan{})
		planListSplit.Trailing = f.formContainer

		planListSplit.Refresh()
	})

	planList.actionBar = widget.NewToolbar(widget.NewToolbarSpacer(), planList.addPlanAction)
	planList.actionBar.Refresh()
	listAndButtons = container.New(layout.NewBorderLayout(nil, planList.actionBar, nil, nil), container.NewMax(planList.List), planList.actionBar)
	planListSplit.Leading = listAndButtons
	planList.List.Refresh()
	listAndButtons.Refresh()
	planListSplit.Leading.Refresh()
	planListSplit.Refresh()

	planList.List.OnSelected = func(id widget.ListItemID) {
		f := NewForm()
		f.FillForm(planList.Plans[id])
		planListSplit.Trailing = f.formContainer
		planListSplit.Refresh()
	}

	return planListSplit
}

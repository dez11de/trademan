package main

import (
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

type UI struct {
	Pairs            []cryptodb.Pair
	activePair       cryptodb.Pair
	Plans            []cryptodb.Plan
	activePlan       cryptodb.Plan
	activeOrders     []cryptodb.Order
	activeAssessment cryptodb.Assessment

	List                *widget.List
	statisticsContainer *fyne.Container
}

var ui UI

func MakePlanListSplit() *container.Split {
	var err error
	ui.Plans, err = getPlans()
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}
	ui.Pairs, err = getPairs()
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}
	ui.List = widget.NewList(
		func() int {
			return len(ui.Plans)
		},
		func() fyne.CanvasObject {
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
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*canvas.Text).Text = ui.Pairs[int64(ui.Plans[i].PairID-1)].Name

			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*canvas.Text).Text = ui.Plans[i].Direction.String()
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*canvas.Text).Color = DirectionColor(ui.Plans[i].Direction)

			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Text = ui.Plans[i].Status.String()
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Color = StatusColor(ui.Plans[i].Status)
		})

	selectPlanLabel := container.New(layout.NewCenterLayout(), canvas.NewText("Select a plan from the list, or press + to make a new plan.", nil))

	ListAndButtons := container.NewWithoutLayout(widget.NewLabel("Error loading plans. Check network connection."))
	planListSplit := container.NewHSplit(ListAndButtons, container.NewMax(selectPlanLabel))
	planListSplit.SetOffset(0.22)

	addPlanAction := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		f := NewForm()
		planListSplit.Trailing = f.formContainer

		f.FillForm(cryptodb.Plan{})
		planListSplit.Refresh()
	})

	refreshListAction := widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
		ui.Plans, err = getPlans()
		if err != nil {
			dialog.ShowError(err, mainWindow)
		}
		ui.List.Refresh()
	})

	actionBar := widget.NewToolbar(widget.NewToolbarSpacer(), refreshListAction, addPlanAction)
	ListAndButtons = container.NewBorder(nil, actionBar, nil, nil, ui.List)
	planListSplit.Leading = ListAndButtons
	planListSplit.Refresh()

	ui.List.OnSelected = func(id widget.ListItemID) {
		f := NewForm()
		planListSplit.Trailing = f.formContainer

		f.FillForm(ui.Plans[id])
		planListSplit.Refresh()
	}

	return planListSplit
}

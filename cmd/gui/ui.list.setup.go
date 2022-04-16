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

func MakePlanListSplit() *container.Split {
	var err error
	act.List = widget.NewList(
		func() int {
			return len(tm.plans)
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
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*canvas.Text).Text = tm.pairs[int64(tm.plans[i].PairID-1)].Name

			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*canvas.Text).Text = tm.plans[i].Direction.String()
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*canvas.Text).Color = DirectionColor(tm.plans[i].Direction)

			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Text = tm.plans[i].Status.String()
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*canvas.Text).Color = StatusColor(tm.plans[i].Status)
		})

	selectPlanLabel := container.New(layout.NewCenterLayout(), canvas.NewText("Select a plan from the list, or press + to make a new plan.", nil))

	ListAndButtons := container.NewWithoutLayout(widget.NewLabel("Error loading plans. Check database connection."))
	planListSplit := container.NewHSplit(ListAndButtons, container.NewMax(selectPlanLabel))
	planListSplit.SetOffset(0.22)

	act.List.OnSelected = func(id widget.ListItemID) {
		act.plan, _ = getPlan(tm.plans[id].ID)
		act.pair = tm.pairs[int64(act.plan.PairID-1)]
		act.orders, _ = getOrders(act.plan.ID)
		act.assessment, _ = getAssessment(act.plan.ID)
		planListSplit.Trailing = makePlanForm()

		planListSplit.Refresh()
	}

	addPlanAction := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		act.plan = cryptodb.Plan{}
		act.pair = cryptodb.Pair{}
		act.pair.QuoteCurrency = " . . . "
		act.orders = cryptodb.NewOrders(0)
		act.assessment = cryptodb.Assessment{}

		planListSplit.Trailing = makePlanForm()

		planListSplit.Refresh()
	})

	refreshListAction := widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
		tm.plans, err = getPlans()
		if err != nil {
			dialog.ShowError(err, mainWindow)
		}
		act.List.Refresh()
	})

	actionBar := widget.NewToolbar(widget.NewToolbarSpacer(), refreshListAction, addPlanAction)
	ListAndButtons = container.NewBorder(nil, actionBar, nil, nil, act.List)
	planListSplit.Leading = ListAndButtons
	planListSplit.Refresh()

	return planListSplit
}

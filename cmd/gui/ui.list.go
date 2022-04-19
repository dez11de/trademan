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
	ui.planList = widget.NewList(
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

	ListAndButtons := container.NewCenter(widget.NewLabel("Error loading plans. Check database connection."))
	planListSplit := container.NewHSplit(ListAndButtons, container.NewCenter(widget.NewLabel("Select a plan from the list, or press + to make a new plan.")))
	planListSplit.SetOffset(0.22)

	ui.planList.OnSelected = func(id widget.ListItemID) {
		tm.plan, _ = getPlan(tm.plans[id].ID)
		tm.pair = tm.pairs[int64(tm.plan.PairID-1)]
		tm.orders, _ = getOrders(tm.plan.ID)
		tm.review, _ = getReview(tm.plan.ID)
		planListSplit.Trailing = NewPlanContainer()

		planListSplit.Refresh()
	}

	addPlanAction := widget.NewToolbarAction(theme.ContentAddIcon(), func() {
		tm.plan = cryptodb.Plan{}
		tm.pair = cryptodb.Pair{}
		tm.pair.QuoteCurrency = " . . . "
		tm.orders = cryptodb.NewOrders(0)
		tm.review = cryptodb.NewReview(0)

		planListSplit.Trailing = NewPlanContainer()

		planListSplit.Refresh()
	})

	refreshListAction := widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
		tm.plans, err = getPlans()
		if err != nil {
			dialog.ShowError(err, ui.mainWindow)
		}
		ui.planList.Refresh()
	})

	actionBar := widget.NewToolbar(widget.NewToolbarSpacer(), refreshListAction, addPlanAction)
	ListAndButtons = container.NewBorder(nil, actionBar, nil, nil, ui.planList)
	planListSplit.Leading = ListAndButtons
	planListSplit.Refresh()

	return planListSplit
}

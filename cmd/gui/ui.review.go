package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

type reviewForm struct {
	parentWindow fyne.Window

	TimingSelectEntry               *widget.SelectEntry
	RiskSelectEntry                 *widget.SelectEntry
	RewardRiskSelectEntry           *widget.SelectEntry
	StopLossSelectEntry             *widget.SelectEntry
	EntrySelectEntry                *widget.SelectEntry
	EmotionSelectEntry              *widget.SelectEntry
	FollowPlanSelectEntry           *widget.SelectEntry
	OrderManagementSelectEntry      *widget.SelectEntry
	MoveStopLossInProfitSelectEntry *widget.SelectEntry
	TakeProfitStrategySelectEntry   *widget.SelectEntry
	TakeProfitCountSelectEntry      *widget.SelectEntry
	FeeSelectEntry                  *widget.SelectEntry
	ProfitSelectEntry               *widget.SelectEntry
	NotesEntry                      *widget.Entry
}

var af reviewForm

func makeReviewForm() *fyne.Container {
	leftColumn := widget.NewForm()
	rightColumn := widget.NewForm()
	bottom := widget.NewForm()

	leftColumn.AppendItem(af.makeRiskItem())
	leftColumn.AppendItem(af.makeRewardRiskItem())
	leftColumn.AppendItem(af.makeEmotionItem())
	leftColumn.AppendItem(af.makeTimingItem())
	leftColumn.AppendItem(af.makeFollowPlanItem())
	leftColumn.AppendItem(af.makeOrderManagementItem())
	leftColumn.AppendItem(af.makeMoveStopLossInProfitItem())

	rightColumn.AppendItem(af.makeStopLossItem())
	rightColumn.AppendItem(af.makeEntryItem())
	rightColumn.AppendItem(af.makeTakeProfitStrategyItem())
	rightColumn.AppendItem(af.makeTakeProfitCountItem())
	rightColumn.AppendItem(af.makeFeeItem())
	rightColumn.AppendItem(af.makeProfitItem())

	bottom.AppendItem(af.makeNotesItem())

	totalContainer := container.NewVBox(container.New(layout.NewGridLayoutWithColumns(2), leftColumn, rightColumn),
		bottom,
		af.makeToolBar())

	return totalContainer
}

func (af *reviewForm) gatherReview() {
	tm.review.Risk = af.RiskSelectEntry.Text
    tm.review.RewardRisk = af.RewardRiskSelectEntry.Text
	tm.review.Emotion = af.EmotionSelectEntry.Text
	tm.review.Timing = af.TimingSelectEntry.Text
	tm.review.FollowPlan = af.FollowPlanSelectEntry.Text
	tm.review.OrderManagement = af.OrderManagementSelectEntry.Text
    tm.review.MoveStopLossInProfit = af.MoveStopLossInProfitSelectEntry.Text

	tm.review.StopLoss = af.StopLossSelectEntry.Text
	tm.review.Entry = af.EntrySelectEntry.Text
	tm.review.TakeProfitStrategy = af.TakeProfitStrategySelectEntry.Text
	tm.review.TakeProfitCount = af.TakeProfitCountSelectEntry.Text
	tm.review.Fee = af.FeeSelectEntry.Text
	tm.review.Profit = af.ProfitSelectEntry.Text
	tm.review.Notes = af.NotesEntry.Text
}

func (af *reviewForm) saveReviewAction() {
	af.gatherReview()
	tm.review, _ = saveReview(tm.review)
	af.parentWindow.Close()
}

func (af *reviewForm) loadReviewAction() {
	tm.review, _ = getReview(tm.plan.ID)
	af.parentWindow.Close()
}

func (af *reviewForm) archiveAction() {
	af.gatherReview()
	tm.review, _ = saveReview(tm.review)
	tm.plan.Status = cryptodb.Archived
	savePlan(tm.plan)
	af.parentWindow.Close()
	tm.plans, _ = getPlans()
    ui.planListSplit.Trailing = ui.noPlanSelectedContainer
    ui.planListSplit.Refresh()
	ui.planList.Refresh()
}

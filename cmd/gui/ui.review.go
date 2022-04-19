package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type reviewForm struct {
	parentWindow fyne.Window

	RiskSelectEntry                 *widget.SelectEntry
	TimingSelectEntry               *widget.SelectEntry
	StopLossSelectEntry             *widget.SelectEntry
	EntrySelectEntry                *widget.SelectEntry
	EmotionSelectEntry              *widget.SelectEntry
	FollowPlanSelectEntry           *widget.SelectEntry
	OrderManagementSelectEntry      *widget.SelectEntry
	MoveStopLossInProfitSelectEntry *widget.SelectEntry
	TakeProfitStrategySelectEntry   *widget.SelectEntry
	TakeProfitCountSelectEntry      *widget.SelectEntry
	NotesEntry                      *widget.Entry
}

var af reviewForm

func makeReviewForm() *fyne.Container {
	leftForm := widget.NewForm()
	rightForm := widget.NewForm()
	bottomForm := widget.NewForm()

	leftForm.AppendItem(af.makeRiskItem())
	leftForm.AppendItem(af.makeTimingItem())
	leftForm.AppendItem(af.makeStopLossItem())
	leftForm.AppendItem(af.makeEntryItem())
	leftForm.AppendItem(af.makeEmotionItem())
	rightForm.AppendItem(af.makeFollowPlanItem())
	rightForm.AppendItem(af.makeOrderManagementItem())
	rightForm.AppendItem(af.makeMoveStopLossInProfitItem())
	rightForm.AppendItem(af.makeTakeProfitStrategyItem())
	rightForm.AppendItem(af.makeTakeProfitCountItem())

	bottomForm.AppendItem(af.makeNotesItem())

	totalContainer := container.NewVBox(container.New(layout.NewGridLayoutWithColumns(2), leftForm, rightForm),
		bottomForm,
		af.makeToolBar())

	return totalContainer
}

func (af *reviewForm) gatherReview() {
	tm.review.Risk = af.RiskSelectEntry.Text
	tm.review.Timing = af.TimingSelectEntry.Text
	tm.review.StopLoss = af.StopLossSelectEntry.Text
	tm.review.Entry = af.EntrySelectEntry.Text
	tm.review.Emotion = af.EmotionSelectEntry.Text
	tm.review.FollowPlan = af.FollowPlanSelectEntry.Text
	tm.review.OrderManagement = af.OrderManagementSelectEntry.Text
	tm.review.TakeProfitStrategy = af.TakeProfitStrategySelectEntry.Text
	tm.review.TakeProfitCount = af.TakeProfitCountSelectEntry.Text
	tm.review.Notes = af.NotesEntry.Text
}

func (af *reviewForm) saveReviewAction() {
	af.gatherReview()
	saveReview(tm.review)
	af.parentWindow.Close()
}

func (af *reviewForm) loadReviewAction() {
	tm.review, _ = getReview(tm.plan.ID)
	af.parentWindow.Close()
}

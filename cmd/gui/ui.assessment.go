package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type assessmentForm struct {
	window fyne.Window

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

var af assessmentForm

func makeAssessmentForm() *fyne.Container {
	leftForm := widget.NewForm()
	rightForm := widget.NewForm()
	notesForm := widget.NewForm()

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

	notesForm.AppendItem(af.makeNotesItem())

	formContainer := container.NewHBox(leftForm, rightForm)
	totalContainer := container.NewVBox(formContainer, notesForm)

	return container.NewBorder(nil, af.makeToolBar(), nil, nil, totalContainer)
}

func (af *assessmentForm) gatherAssessment() {
	act.assessment.Risk = af.RiskSelectEntry.Text
	act.assessment.Timing = af.TimingSelectEntry.Text
	act.assessment.StopLoss = af.StopLossSelectEntry.Text
	act.assessment.Entry = af.EntrySelectEntry.Text
	act.assessment.Emotion = af.EmotionSelectEntry.Text
	act.assessment.FollowPlan = af.FollowPlanSelectEntry.Text
	act.assessment.OrderManagement = af.OrderManagementSelectEntry.Text
	act.assessment.TakeProfitStrategy = af.TakeProfitStrategySelectEntry.Text
	act.assessment.TakeProfitCount = af.TakeProfitCountSelectEntry.Text
	act.assessment.Notes = af.NotesEntry.Text
}

func (af *assessmentForm) saveAssessmentAction() {
	af.gatherAssessment()
	saveAssessment(act.assessment)
	af.window.Close()
}

func (af *assessmentForm) loadAssessmentAction() {
	act.assessment, _ = getAssessment(act.plan.ID)
	af.window.Close()
}

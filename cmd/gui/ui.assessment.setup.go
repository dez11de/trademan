package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

type assessmentForm struct {
	form      *widget.Form
	container *fyne.Container
	window    fyne.Window

	assessmentOptions map[string][]string

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
}

func NewAssessmentForm(p cryptodb.Plan) *assessmentForm {
	var af assessmentForm

	var err error
	af.assessmentOptions, err = getAssessmentOptions()
	if err != nil {
        dialog.ShowError(err, af.window)
	}

	af.form = widget.NewForm()

	af.form.AppendItem(af.makeRiskItem())
	af.form.AppendItem(af.makeTimingItem())
	af.form.AppendItem(af.makeStopLossItem())
	af.form.AppendItem(af.makeEntryItem())
	af.form.AppendItem(af.makeEmotionItem())
	af.form.AppendItem(af.makeFollowPlanItem())
	af.form.AppendItem(af.makeOrderManagementItem())
	af.form.AppendItem(af.makeMoveStopLossInProfitItem())
	af.form.AppendItem(af.makeTakeProfitStrategyItem())
	af.form.AppendItem(af.makeTakeProfitCountItem())

	af.container = container.NewBorder(nil, af.makeToolBar(), nil, nil, af.form)

	return &af
}

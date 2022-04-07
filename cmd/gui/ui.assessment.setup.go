package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

type assessmentForm struct {
    form *widget.Form
    assessmentFormContainer *fyne.Container
}

func NewAssessmentForm(p cryptodb.Plan) *fyne.Container {
    var af assessmentForm

    af.form = widget.NewForm()

	RiskAssessmentItem := af.makeRiskItem()
	af.form.AppendItem(RiskAssessmentItem)

	TimingAssessmentItem := af.makeTimingItem()
	af.form.AppendItem(TimingAssessmentItem)

	c := container.NewBorder(nil, af.makeToolBar(), nil, nil, af.form)

	return c
}

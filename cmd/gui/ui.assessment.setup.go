package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

type assessmentForm struct {
	form      *widget.Form
	container *fyne.Container
	window    fyne.Window

	RiskSelect   *widget.Select
	TimingSelect *widget.Select
}

func NewAssessmentForm(p cryptodb.Plan) *assessmentForm {
	var af assessmentForm

	af.form = widget.NewForm()

	af.form.AppendItem(af.makeRiskItem())
	af.form.AppendItem(af.makeTimingItem())

	af.container = container.NewBorder(nil, af.makeToolBar(), nil, nil, af.form)

	return &af
}

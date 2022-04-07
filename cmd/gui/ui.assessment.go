package main

import (
	"fyne.io/fyne/v2/dialog"
)

func (af *assessmentForm) saveAssessmentAction() {
    var err error
	ui.activeAssessment, err = saveAssessment(ui.activeAssessment)
	if err != nil {
		dialog.ShowError(err, mainWindow)
	}
}

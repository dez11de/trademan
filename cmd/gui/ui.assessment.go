package main

func (af *assessmentForm) gatherAssessment() {
    ui.activeAssessment.Risk = af.RiskSelect.Selected
    ui.activeAssessment.Timing = af.TimingSelect.Selected
}

func (af *assessmentForm) saveAssessmentAction() {
    af.gatherAssessment()
    saveAssessment(ui.activeAssessment)
    af.window.Close()
}

package main

func (af *assessmentForm) gatherAssessment() {
    ui.activeAssessment.Risk = af.RiskSelectEntry.Text
    ui.activeAssessment.Timing = af.TimingSelectEntry.Text
    ui.activeAssessment.StopLoss = af.StopLossSelectEntry.Text
    ui.activeAssessment.Entry = af.EntrySelectEntry.Text
    ui.activeAssessment.Emotion = af.EmotionSelectEntry.Text
    ui.activeAssessment.FollowPlan = af.FollowPlanSelectEntry.Text
    ui.activeAssessment.OrderManagement = af.OrderManagementSelectEntry.Text
    ui.activeAssessment.TakeProfitStrategy = af.TakeProfitStrategySelectEntry.Text
    ui.activeAssessment.TakeProfitCount = af.TakeProfitCountSelectEntry.Text
}

func (af *assessmentForm) saveAssessmentAction() {
    af.gatherAssessment()
    saveAssessment(ui.activeAssessment)
    af.window.Close()
}

func (af *assessmentForm) loadAssessmentAction() {
    ui.activeAssessment, _ = getAssessment(ui.activePlan.ID)
    af.window.Close()
}

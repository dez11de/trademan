package main

import (
	"fmt"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

func (af *assessmentForm) makeRiskItem() *widget.FormItem {
	af.RiskSelectEntry = widget.NewSelectEntry(af.assessmentOptions["risk"])
	af.RiskSelectEntry.SetText(ui.activeAssessment.Risk)
	i := widget.NewFormItem("Risk", af.RiskSelectEntry)
	i.HintText = fmt.Sprintf("Risk was: %s%%", ui.activePlan.Risk.StringFixed(2))

	return i
}

func (af *assessmentForm) makeTimingItem() *widget.FormItem {
	af.TimingSelectEntry = widget.NewSelectEntry(af.assessmentOptions["timing"])
	af.TimingSelectEntry.SetText(ui.activeAssessment.Timing)
	i := widget.NewFormItem("Timing", af.TimingSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeStopLossItem() *widget.FormItem {
	af.StopLossSelectEntry = widget.NewSelectEntry(af.assessmentOptions["stop_loss"])
	af.StopLossSelectEntry.SetText(ui.activeAssessment.StopLoss)
	i := widget.NewFormItem("StopLoss", af.StopLossSelectEntry)
	i.HintText = fmt.Sprintf("StopLoss was: %s %s", ui.activeOrders[cryptodb.MarketStopLoss].Price.StringFixed(ui.activePair.PriceScale), ui.activePair.QuoteCurrency)

	return i
}

func (af *assessmentForm) makeEntryItem() *widget.FormItem {
	af.EntrySelectEntry = widget.NewSelectEntry(af.assessmentOptions["entry"])
	af.EntrySelectEntry.SetText(ui.activeAssessment.Entry)
	i := widget.NewFormItem("Entry", af.EntrySelectEntry)
	i.HintText = fmt.Sprintf("Entry was: %s %s", ui.activeOrders[cryptodb.Entry].Price.StringFixed(ui.activePair.PriceScale), ui.activePair.QuoteCurrency)

	return i
}

func (af *assessmentForm) makeEmotionItem() *widget.FormItem {
	af.EmotionSelectEntry = widget.NewSelectEntry(af.assessmentOptions["emotion"])
	af.EmotionSelectEntry.SetText(ui.activeAssessment.Emotion)
	i := widget.NewFormItem("Emotion", af.EmotionSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeFollowPlanItem() *widget.FormItem {
	af.FollowPlanSelectEntry = widget.NewSelectEntry(af.assessmentOptions["follow_plan"])
	af.FollowPlanSelectEntry.SetText(ui.activeAssessment.FollowPlan)
	i := widget.NewFormItem("Follow plan", af.FollowPlanSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeOrderManagementItem() *widget.FormItem {
	af.OrderManagementSelectEntry = widget.NewSelectEntry(af.assessmentOptions["order_management"])
	af.OrderManagementSelectEntry.SetText(ui.activeAssessment.OrderManagement)
	i := widget.NewFormItem("Order management", af.OrderManagementSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeMoveStopLossInProfitItem() *widget.FormItem {
	af.MoveStopLossInProfitSelectEntry = widget.NewSelectEntry(af.assessmentOptions["move_stop_loss_in_profit"])
	af.MoveStopLossInProfitSelectEntry.SetText(ui.activeAssessment.MoveStopLossInProfit)
	i := widget.NewFormItem("StopLoss in Profit", af.MoveStopLossInProfitSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeTakeProfitStrategyItem() *widget.FormItem {
	af.TakeProfitStrategySelectEntry = widget.NewSelectEntry(af.assessmentOptions["take_profit_strategy"])
	af.TakeProfitStrategySelectEntry.SetText(ui.activeAssessment.TakeProfitStrategy)
	i := widget.NewFormItem("Take Profit Strategy", af.TakeProfitStrategySelectEntry)
	i.HintText = fmt.Sprintf("Strategy was: %s", ui.activePlan.TakeProfitStrategy.String())

	return i
}

func (af *assessmentForm) makeTakeProfitCountItem() *widget.FormItem {
	af.TakeProfitCountSelectEntry = widget.NewSelectEntry(af.assessmentOptions["take_profit_count"])
	af.TakeProfitCountSelectEntry.SetText(ui.activeAssessment.TakeProfitCount)
	i := widget.NewFormItem("Take Profit Count", af.TakeProfitCountSelectEntry)
	i.HintText = "too diffucult to calculate"

	return i
}

func (af *assessmentForm) makeToolBar() *widget.Toolbar {
	saveAssassmentButton := widget.NewToolbarAction(theme.ConfirmIcon(), af.saveAssessmentAction)
	cancelAssessmentButton := widget.NewToolbarAction(theme.ContentUndoIcon(), nil)
	toolbar := widget.NewToolbar(widget.NewToolbarSpacer(), cancelAssessmentButton, saveAssassmentButton)
	return toolbar
}

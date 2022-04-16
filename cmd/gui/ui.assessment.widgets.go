package main

import (
	"fmt"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

func (af *assessmentForm) makeRiskItem() *widget.FormItem {
    af.RiskSelectEntry = widget.NewSelectEntry(tm.assessmentOptions["risk"])
	af.RiskSelectEntry.SetText(act.assessment.Risk)
	i := widget.NewFormItem("Risk", af.RiskSelectEntry)
	i.HintText = fmt.Sprintf("Risk was: %s%%", act.plan.Risk.StringFixed(2))

	return i
}

func (af *assessmentForm) makeTimingItem() *widget.FormItem {
	af.TimingSelectEntry = widget.NewSelectEntry(tm.assessmentOptions["timing"])
	af.TimingSelectEntry.SetText(act.assessment.Timing)
	i := widget.NewFormItem("Timing", af.TimingSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeStopLossItem() *widget.FormItem {
	af.StopLossSelectEntry = widget.NewSelectEntry(tm.assessmentOptions["stop_loss"])
	af.StopLossSelectEntry.SetText(act.assessment.StopLoss)
	i := widget.NewFormItem("StopLoss", af.StopLossSelectEntry)
	i.HintText = fmt.Sprintf("StopLoss was: %s %s", act.orders[cryptodb.MarketStopLoss].Price.StringFixed(act.pair.PriceScale), act.pair.QuoteCurrency)

	return i
}

func (af *assessmentForm) makeEntryItem() *widget.FormItem {
	af.EntrySelectEntry = widget.NewSelectEntry(tm.assessmentOptions["entry"])
	af.EntrySelectEntry.SetText(act.assessment.Entry)
	i := widget.NewFormItem("Entry", af.EntrySelectEntry)
	i.HintText = fmt.Sprintf("Entry was: %s %s", act.orders[cryptodb.Entry].Price.StringFixed(act.pair.PriceScale), act.pair.QuoteCurrency)

	return i
}

func (af *assessmentForm) makeEmotionItem() *widget.FormItem {
	af.EmotionSelectEntry = widget.NewSelectEntry(tm.assessmentOptions["emotion"])
	af.EmotionSelectEntry.SetText(act.assessment.Emotion)
	i := widget.NewFormItem("Emotion", af.EmotionSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeFollowPlanItem() *widget.FormItem {
	af.FollowPlanSelectEntry = widget.NewSelectEntry(tm.assessmentOptions["follow_plan"])
	af.FollowPlanSelectEntry.SetText(act.assessment.FollowPlan)
	i := widget.NewFormItem("Follow plan", af.FollowPlanSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeOrderManagementItem() *widget.FormItem {
	af.OrderManagementSelectEntry = widget.NewSelectEntry(tm.assessmentOptions["order_management"])
	af.OrderManagementSelectEntry.SetText(act.assessment.OrderManagement)
	i := widget.NewFormItem("Order management", af.OrderManagementSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeMoveStopLossInProfitItem() *widget.FormItem {
	af.MoveStopLossInProfitSelectEntry = widget.NewSelectEntry(tm.assessmentOptions["move_stop_loss_in_profit"])
	af.MoveStopLossInProfitSelectEntry.SetText(act.assessment.MoveStopLossInProfit)
	i := widget.NewFormItem("StopLoss in Profit", af.MoveStopLossInProfitSelectEntry)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeTakeProfitStrategyItem() *widget.FormItem {
	af.TakeProfitStrategySelectEntry = widget.NewSelectEntry(tm.assessmentOptions["take_profit_strategy"])
	af.TakeProfitStrategySelectEntry.SetText(act.assessment.TakeProfitStrategy)
	i := widget.NewFormItem("Take Profit Strategy", af.TakeProfitStrategySelectEntry)
	i.HintText = fmt.Sprintf("Strategy was: %s", act.plan.TakeProfitStrategy.String())

	return i
}

func (af *assessmentForm) makeTakeProfitCountItem() *widget.FormItem {
	af.TakeProfitCountSelectEntry = widget.NewSelectEntry(tm.assessmentOptions["take_profit_count"])
	af.TakeProfitCountSelectEntry.SetText(act.assessment.TakeProfitCount)
	i := widget.NewFormItem("Take Profit Count", af.TakeProfitCountSelectEntry)
	i.HintText = "too diffucult to calculate"

	return i
}

func (af *assessmentForm) makeNotesItem() *widget.FormItem {
	af.NotesEntry = widget.NewMultiLineEntry()
	af.NotesEntry.SetText(act.assessment.Notes)
	i := widget.NewFormItem("Notes", af.NotesEntry)

	return i
}

func (af *assessmentForm) makeToolBar() *widget.Toolbar {
	saveAssassmentButton := widget.NewToolbarAction(theme.ConfirmIcon(), af.saveAssessmentAction)
	cancelAssessmentButton := widget.NewToolbarAction(theme.ContentUndoIcon(), nil)
	toolbar := widget.NewToolbar(widget.NewToolbarSpacer(), cancelAssessmentButton, saveAssassmentButton)
	return toolbar
}

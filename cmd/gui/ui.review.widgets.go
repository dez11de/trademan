package main

import (
	"fmt"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

func (af *reviewForm) makeRiskItem() *widget.FormItem {
    af.RiskSelectEntry = widget.NewSelectEntry(tm.reviewOptions["risk"])
	af.RiskSelectEntry.SetText(act.review.Risk)
	i := widget.NewFormItem("Risk", af.RiskSelectEntry)
	i.HintText = fmt.Sprintf("Risk was: %s%%", act.plan.Risk.StringFixed(2))

	return i
}

func (af *reviewForm) makeTimingItem() *widget.FormItem {
	af.TimingSelectEntry = widget.NewSelectEntry(tm.reviewOptions["timing"])
	af.TimingSelectEntry.SetText(act.review.Timing)
	i := widget.NewFormItem("Timing", af.TimingSelectEntry)
	i.HintText = " "

	return i
}

func (af *reviewForm) makeStopLossItem() *widget.FormItem {
	af.StopLossSelectEntry = widget.NewSelectEntry(tm.reviewOptions["stop_loss"])
	af.StopLossSelectEntry.SetText(act.review.StopLoss)
	i := widget.NewFormItem("StopLoss", af.StopLossSelectEntry)
	i.HintText = fmt.Sprintf("StopLoss was: %s %s", act.orders[cryptodb.MarketStopLoss].Price.StringFixed(act.pair.PriceScale), act.pair.QuoteCurrency)

	return i
}

func (af *reviewForm) makeEntryItem() *widget.FormItem {
	af.EntrySelectEntry = widget.NewSelectEntry(tm.reviewOptions["entry"])
	af.EntrySelectEntry.SetText(act.review.Entry)
	i := widget.NewFormItem("Entry", af.EntrySelectEntry)
	i.HintText = fmt.Sprintf("Entry was: %s %s", act.orders[cryptodb.Entry].Price.StringFixed(act.pair.PriceScale), act.pair.QuoteCurrency)

	return i
}

func (af *reviewForm) makeEmotionItem() *widget.FormItem {
	af.EmotionSelectEntry = widget.NewSelectEntry(tm.reviewOptions["emotion"])
	af.EmotionSelectEntry.SetText(act.review.Emotion)
	i := widget.NewFormItem("Emotion", af.EmotionSelectEntry)
	i.HintText = " "

	return i
}

func (af *reviewForm) makeFollowPlanItem() *widget.FormItem {
	af.FollowPlanSelectEntry = widget.NewSelectEntry(tm.reviewOptions["follow_plan"])
	af.FollowPlanSelectEntry.SetText(act.review.FollowPlan)
	i := widget.NewFormItem("Follow plan", af.FollowPlanSelectEntry)
	i.HintText = " "

	return i
}

func (af *reviewForm) makeOrderManagementItem() *widget.FormItem {
	af.OrderManagementSelectEntry = widget.NewSelectEntry(tm.reviewOptions["order_management"])
	af.OrderManagementSelectEntry.SetText(act.review.OrderManagement)
	i := widget.NewFormItem("Order management", af.OrderManagementSelectEntry)
	i.HintText = " "

	return i
}

func (af *reviewForm) makeMoveStopLossInProfitItem() *widget.FormItem {
	af.MoveStopLossInProfitSelectEntry = widget.NewSelectEntry(tm.reviewOptions["move_stop_loss_in_profit"])
	af.MoveStopLossInProfitSelectEntry.SetText(act.review.MoveStopLossInProfit)
	i := widget.NewFormItem("StopLoss in Profit", af.MoveStopLossInProfitSelectEntry)
	i.HintText = " "

	return i
}

func (af *reviewForm) makeTakeProfitStrategyItem() *widget.FormItem {
	af.TakeProfitStrategySelectEntry = widget.NewSelectEntry(tm.reviewOptions["take_profit_strategy"])
	af.TakeProfitStrategySelectEntry.SetText(act.review.TakeProfitStrategy)
	i := widget.NewFormItem("Take Profit Strategy", af.TakeProfitStrategySelectEntry)
	i.HintText = fmt.Sprintf("Strategy was: %s", act.plan.TakeProfitStrategy.String())

	return i
}

func (af *reviewForm) makeTakeProfitCountItem() *widget.FormItem {
	af.TakeProfitCountSelectEntry = widget.NewSelectEntry(tm.reviewOptions["take_profit_count"])
	af.TakeProfitCountSelectEntry.SetText(act.review.TakeProfitCount)
	i := widget.NewFormItem("Take Profit Count", af.TakeProfitCountSelectEntry)
	i.HintText = "too diffucult to calculate"

	return i
}

func (af *reviewForm) makeNotesItem() *widget.FormItem {
	af.NotesEntry = widget.NewMultiLineEntry()
	af.NotesEntry.SetText(act.review.Notes)
	i := widget.NewFormItem("Notes", af.NotesEntry)

	return i
}

func (af *reviewForm) makeToolBar() *widget.Toolbar {
	saveAssassmentButton := widget.NewToolbarAction(theme.ConfirmIcon(), af.saveReviewAction)
	cancelAssessmentButton := widget.NewToolbarAction(theme.ContentUndoIcon(), nil)
	toolbar := widget.NewToolbar(widget.NewToolbarSpacer(), cancelAssessmentButton, saveAssassmentButton)
	return toolbar
}

package main

import (
	"fmt"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dez11de/cryptodb"
)

func (af *assessmentForm) makeRiskItem() *widget.FormItem {
	af.RiskSelect = widget.NewSelect([]string{"Too Low", "Good", "Too High"}, nil)
	af.RiskSelect.SetSelected(ui.activeAssessment.Risk)
	i := widget.NewFormItem("Risk", af.RiskSelect)
	i.HintText = fmt.Sprintf("Risk was: %s%%", ui.activePlan.Risk.StringFixed(2))

	return i
}

func (af *assessmentForm) makeTimingItem() *widget.FormItem {
	af.TimingSelect = widget.NewSelect([]string{"Early", "On-time", "Late"}, nil)
	af.TimingSelect.SetSelected(ui.activeAssessment.Timing)
	i := widget.NewFormItem("Timing", af.TimingSelect)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeStopLossItem() *widget.FormItem {
	StopLossSelect := widget.NewSelect([]string{"Too tight", "Good", "Too wide"}, nil)
	StopLossSelect.SetSelected(ui.activeAssessment.StopLossPosition)
	i := widget.NewFormItem("StopLoss", StopLossSelect)
	i.HintText = fmt.Sprintf("StopLoss was: %s %s", ui.activeOrders[cryptodb.MarketStopLoss].Price.StringFixed(ui.activePair.PriceScale), ui.activePair.QuoteCurrency)

	return i
}

func (af *assessmentForm) makeEntryItem() *widget.FormItem {
	EntrySelect := widget.NewSelect([]string{"Too tight", "Good", "Too wide"}, nil)
	EntrySelect.SetSelected(ui.activeAssessment.EntryPosition)
	i := widget.NewFormItem("Entry", EntrySelect)
	i.HintText = fmt.Sprintf("Entry was: %s %s", ui.activeOrders[cryptodb.Entry].Price.StringFixed(ui.activePair.PriceScale), ui.activePair.QuoteCurrency)

	return i
}

func (af *assessmentForm) makeEmotionItem() *widget.FormItem {
	EmotionSelect := widget.NewSelect([]string{"FOMO", "In control", "FOJI"}, nil)
	EmotionSelect.SetSelected(ui.activeAssessment.Emotion)
	i := widget.NewFormItem("Emotion", EmotionSelect)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeFollowPlanItem() *widget.FormItem {
	FolowPlanSelect := widget.NewSelect([]string{"No", "Neutral", "Yes"}, nil)
	FolowPlanSelect.SetSelected(ui.activeAssessment.FollowPlan)
	i := widget.NewFormItem("Follow plan", FolowPlanSelect)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeOrderManagementItem() *widget.FormItem {
	OrderManagementSelect := widget.NewSelect([]string{"Bad", "Neutral", "Good"}, nil)
	OrderManagementSelect.SetSelected(ui.activeAssessment.OrderManagement)
	i := widget.NewFormItem("Order management", OrderManagementSelect)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeMoveStopLossInProfitItem() *widget.FormItem {
	MoveStopLossInProfitSelect := widget.NewSelect([]string{"Early", "On-time", "Late"}, nil)
	MoveStopLossInProfitSelect.SetSelected(ui.activeAssessment.MoveStopLossInProfit)
	i := widget.NewFormItem("StopLoss in Profit", MoveStopLossInProfitSelect)
	i.HintText = " "

	return i
}

func (af *assessmentForm) makeTakeProfitStrategyItem() *widget.FormItem {
	TakeProfitStrategySelect := widget.NewSelect([]string{"Bad", "Neutral", "Good"}, nil)
	TakeProfitStrategySelect.SetSelected(ui.activeAssessment.TakeProfitStrategy)
	i := widget.NewFormItem("Take Profit Strategy", TakeProfitStrategySelect)
	i.HintText = fmt.Sprintf("Strategy was: %s", ui.activePlan.TakeProfitStrategy.String())

	return i
}

func (af *assessmentForm) makeTakeProfitCountItem() *widget.FormItem {
	TakeProfitCountSelect := widget.NewSelect([]string{"Too few", "Good", "Too many"}, nil)
	TakeProfitCountSelect.SetSelected(ui.activeAssessment.TakeProfitCount)
	i := widget.NewFormItem("Take Profit Strategy", TakeProfitCountSelect)
	i.HintText = "too diffucult to calculate"

	return i
}

func (af *assessmentForm) makeToolBar() *widget.Toolbar {
	saveAssassmentButton := widget.NewToolbarAction(theme.ConfirmIcon(), af.saveAssessmentAction)
	cancelAssessmentButton := widget.NewToolbarAction(theme.ContentUndoIcon(), nil)
	toolbar := widget.NewToolbar(widget.NewToolbarSpacer(), cancelAssessmentButton, saveAssassmentButton)
	return toolbar
}

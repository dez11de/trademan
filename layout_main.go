package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetupMain() (ui UI) {
	ui.app = tview.NewApplication()

	ui.loggingWindow = tview.NewTextView().SetText("----- logging started\n")
	ui.loggingWindow.SetBorder(true).SetTitle(" Log ")

	ui.positionTable = tview.NewTable().SetBorders(false).SetSelectable(true, false).SetFixed(1, 1).Select(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyESC {
			ui.app.Stop()
		}
	})

	ui.positionTable.SetSelectionChangedFunc(func(row, column int) {
		if row < 1 {
			ui.positionTable.Select(1, column)
		}
	})
	ui.positionTable.SetCell(0, 0, tview.NewTableCell("Status").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.positionTable.SetCell(0, 1, tview.NewTableCell("Symbol").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.positionTable.SetCell(0, 2, tview.NewTableCell("Side").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.positionTable.SetCell(0, 3, tview.NewTableCell("Current").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.positionTable.SetCell(0, 4, tview.NewTableCell("Stop Loss").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	ui.positionTable.SetCell(0, 5, tview.NewTableCell("Entry").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

	ui.positionTable.SetCell(1, 0, tview.NewTableCell("Order").SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(1, 1, tview.NewTableCell("BTCUSDT").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(1, 2, tview.NewTableCell("Long").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(1, 3, tview.NewTableCell("62123.0").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(1, 4, tview.NewTableCell("61000.0").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(1, 5, tview.NewTableCell("62000.0").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))

	ui.positionTable.SetCell(2, 0, tview.NewTableCell("Position").SetTextColor(tcell.ColorBlue).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(2, 1, tview.NewTableCell("ETHUSDT").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(2, 2, tview.NewTableCell("Long").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(2, 3, tview.NewTableCell("4321").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(2, 4, tview.NewTableCell("4200").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))
	ui.positionTable.SetCell(2, 5, tview.NewTableCell("4100").SetTextColor(tcell.ColorWhite).SetAlign(tview.AlignLeft))

	ui.positionForm = tview.NewForm().
		AddInputField("Symbol", "", 10, nil, nil).
		AddDropDown("Side", []string{"Long", "Short"}, 0, nil).
		AddInputField("Entry Price", "", 6, nil, nil).
		AddInputField("Stop Loss", "", 6, nil, nil).
		AddInputField("Take Profit #1", "", 6, nil, nil).
		AddButton("Cancel", nil).
		AddButton("Save as plan", nil).
		AddButton("Execute", nil)

	ui.positionForm.SetBorder(true)
	ui.positionForm.AddButton("Add TP", func() {
		ui.positionForm.AddInputField("Another TP", "", 6, nil, nil)
	})
	ui.pages = tview.NewPages().
		AddPage("Positions", ui.positionTable, true, true).
		AddPage("Position edit", ui.positionForm, true, false)

	ui.flex = tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(ui.pages, 0, 20, true).
			AddItem(ui.loggingWindow, 0, 5, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle(" Performance "), 0, 2, false), 0, 20, false)

	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyESC:
			fmt.Fprintf(ui.loggingWindow, "---- <%s> pressed\n", "---- <ESC> pressed")
			ui.pages.SwitchToPage("Positions")
			return nil
		}
		switch event.Rune() {
		case 'E':
			fmt.Fprintf(ui.loggingWindow, "%s", "---- <E> pressed\n")
			ui.pages.SwitchToPage("Position edit")
			return event
		default:
			fmt.Fprintf(ui.loggingWindow, "---- <%v> pressed\n", event.Rune())
			return event
		}
	})

	ui.app.SetRoot(ui.flex, true).SetFocus(ui.pages)
	return ui
}

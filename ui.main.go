package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (ui mainUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		ui.log.Add(Debug, fmt.Sprintf("Key pressed %v", msg))
		switch {
		case key.Matches(msg, ui.paneKeys.quitApp):
			return ui, tea.Quit
		case key.Matches(msg, ui.paneKeys.nextPane):
			ui.activePane++
		case key.Matches(msg, ui.paneKeys.previousPane):
			ui.activePane--
		case key.Matches(msg, ui.paneKeys.toggleHelp):
			ui.log.Add(Debug, "Help toggled")
			ui.help.ShowAll = !ui.help.ShowAll
		}
	}

	if ui.activePane < 0 {
		ui.activePane = 2
	}
	if ui.activePane > 2 {
		ui.activePane = 0
	}

	for i := 0; i < len(ui.panes)-1; i++ {
		ui.panes[i].Border(lipgloss.NormalBorder())
	}
	ui.panes[ui.activePane].Border(lipgloss.DoubleBorder())

	newListModel, cmd := ui.positionList.list.Update(msg)
	ui.positionList.list = newListModel

	cmds = append(cmds, cmd)
	return ui, tea.Batch(cmds...)
}

func (ui mainUI) View() string {
	// TODO: Full Help isn't working yet.
	var helpView string
	switch ui.activePane {
	case int(PositionPane):
		helpView = ui.help.View(ui.paneKeys)
	case int(LogPane):
		helpView = ui.help.View(ui.logKeys)
	default:
		helpView = ""
	}
	helpHeight := strings.Count(helpView, "\n") + 1
	panesStr := fmt.Sprint(lipgloss.JoinVertical(lipgloss.Top,
		ui.panes[PositionPane].Render(ui.positionList.View()),
		ui.panes[LogPane].Render(ui.log.View()),
		ui.panes[PerformancePane].Render(fmt.Sprintf("Performance WxH %dx%d", ui.panes[PerformancePane].GetWidth(), ui.panes[PerformancePane].GetHeight()))))
	return fmt.Sprint(panesStr + "\n" + lipgloss.PlaceVertical(helpHeight, lipgloss.Top, helpView))
}

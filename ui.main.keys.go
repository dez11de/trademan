package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type panesKeyMap struct {
	nextPane     key.Binding
	previousPane key.Binding
	toggleHelp   key.Binding
	quitApp      key.Binding
}

func newPanesKeyMap() panesKeyMap {
	return panesKeyMap{
		nextPane: key.NewBinding(
			key.WithKeys(tea.KeyTab.String()),
			key.WithHelp("<TAB>", "Next pane"),
		),
		previousPane: key.NewBinding(
			key.WithKeys(tea.KeyShiftTab.String()),
			key.WithHelp("<SHIFT><TAB>", "Previous pane"),
		),
		toggleHelp: key.NewBinding(
			key.WithKeys("h", "?"),
			key.WithHelp("h/?", "Toggle help"),
		),
		quitApp: key.NewBinding(
			key.WithKeys("q", tea.KeyCtrlC.String()),
			key.WithHelp("q", "Quit"),
		),
	}
}

func (k panesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.nextPane,
		k.previousPane,
		k.toggleHelp,
		k.quitApp,
	}
}

func (k panesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.nextPane,
			k.previousPane,
		},
		{
			k.toggleHelp,
			k.quitApp,
		},
	}
}

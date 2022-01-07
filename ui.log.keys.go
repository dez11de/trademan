package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type logKeyMap struct {
	toggleTimestamp key.Binding
	incLevel        key.Binding
	decLevel        key.Binding
}

func newLogKeyMap() logKeyMap {
	return logKeyMap{
		toggleTimestamp: key.NewBinding(
			key.WithKeys("d", "D"),
			key.WithHelp("d", "Toggle timestamp"),
		),
		incLevel: key.NewBinding(
			key.WithKeys("+", "="),
			key.WithHelp("+", "Increase level"),
		),
		decLevel: key.NewBinding(
			key.WithKeys("-", "_"),
			key.WithHelp("-", "Decrease level"),
		),
	}
}

func (k logKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.toggleTimestamp,
		k.incLevel,
		k.decLevel,
	}
}

func (k logKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.toggleTimestamp,
		},
		{
			k.incLevel,
			k.decLevel,
		},
	}
}

package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m positionListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m positionListModel) View() string {
	return m.list.View()
}

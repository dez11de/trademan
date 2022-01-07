package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type positionListModel struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	db           *Database
}

func newPositionListModel(d *Database, height int) positionListModel {
	var (
		delegateKeys = newDelegateKeyMap()
		//		listKeys     = newListKeyMap()
	)

	positionStatusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

	p, _ := d.GetPositions()
	items := make([]list.Item, len(p))
	for i := 0; i < len(p); i++ {
		items[i] = item{title: p[i].Symbol + " " + positionStatusStyle.Render(p[i].Status.String()),
			description: "Size: " + strconv.FormatFloat(p[i].Size, 'f', 4, 64),
		}
	}

	delegate := newItemDelegate(delegateKeys)
	delegate.Styles.SelectedTitle.UnsetForeground()
	delegate.Styles.SelectedDesc.UnsetForeground()
	delegate.Styles.NormalTitle.UnsetForeground()
	delegate.Styles.NormalDesc.UnsetForeground()
	delegate.Styles.DimmedTitle.UnsetForeground()
	delegate.Styles.DimmedDesc.UnsetForeground()
	delegate.Styles.SelectedTitle.UnsetBorderStyle()
	delegate.Styles.SelectedDesc.UnsetBorderStyle()
	delegate.Styles.SelectedTitle.Background(lipgloss.Color("#121433"))
	delegate.Styles.SelectedDesc.Background(lipgloss.Color("#121433"))
	delegate.Styles.SelectedTitle.Width(155)
	delegate.Styles.SelectedDesc.Width(155)
	delegate.Styles.SelectedTitle.Padding(0, 0, 0, 2)
	delegate.Styles.SelectedDesc.Padding(0, 0, 0, 2)

	positionList := list.NewModel(items, delegate, 0, height)
	/*
		positionList.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{
				listKeys.insertItem,
				listKeys.toggleHelpMenu,
			}
		}
	*/
	positionList.SetShowTitle(false)
	positionList.Styles.TitleBar.UnsetPadding()
	positionList.SetShowStatusBar(false)
	positionList.SetShowHelp(false)

	return positionListModel{
		list: positionList,
		//		keys:         listKeys,
		//		delegateKeys: delegateKeys,
		db: d,
		//		itemGenerator: &itemGenerator,
	}
}

func (m positionListModel) Init() tea.Cmd {
	return nil
}

package main

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listModel struct {
	list list.Model
	//	itemGenerator *randomItemGenerator
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	db           *Database
}

func newModel(d *Database) listModel {
	var (
		//		itemGenerator randomItemGenerator
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	// Make initial list of items
	p, _ := d.GetPositions()
	items := make([]list.Item, len(p))
	for i := 0; i < len(p); i++ {
		items[i] = item{title: p[i].Symbol,
			description: p[i].Status.String() + "\t" + strconv.FormatFloat(p[i].Size, 'f', 4, 64),
		}
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	positionList := list.NewModel(items, delegate, 0, 0)
	positionList.Title = "Positions"
	positionList.Styles.Title = titleStyle
	positionList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.insertItem,
			listKeys.toggleStatusBar,
			listKeys.toggleHelpMenu,
		}
	}

	return listModel{
		list:         positionList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
		db:           d,
		//		itemGenerator: &itemGenerator,
	}
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		topGap, rightGap, bottomGap, leftGap := appStyle.GetPadding()
		m.list.SetSize(msg.Width-leftGap-rightGap, msg.Height-topGap-bottomGap)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.keys.insertItem):
			m.delegateKeys.remove.SetEnabled(true)
			// newItem := m.itemGenerator.next()
			// insCmd := m.list.InsertItem(0, newItem)
			// statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
			// return m, tea.Batch(insCmd, statusCmd)
			return m, nil
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m listModel) View() string {
	return appStyle.Render(m.list.View())
}

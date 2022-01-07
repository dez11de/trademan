package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	PositionPane int = iota
	LogPane
	PerformancePane
)

const (
	// think of these as %
	relPositionPaneHeight    = 36
	relLogPaneHeight         = 55
	relPerformancePaneHeight = 5

	verticalBorder = 1
	verticalMargin = 1

	horizontalBorder = 1
	horizontalMargin = 1

	shortHelpHeight = 1

	totalHeight = verticalMargin +
		verticalBorder + relPositionPaneHeight + verticalBorder +
		verticalBorder + relLogPaneHeight + verticalBorder +
		verticalBorder + relPerformancePaneHeight + verticalBorder +
		shortHelpHeight +
		verticalMargin
)

type mainUI struct {
	panes      []lipgloss.Style
	activePane int

	log logModel

	paneKeys panesKeyMap
	logKeys  logKeyMap

	db           *Database
	positionList positionListModel
	help         help.Model
}

var (
	paneStyle = lipgloss.NewStyle().
		Margin(0, horizontalMargin, 0)
)

func newMainUIModel(d *Database) (ui mainUI) {
	ui.panes = make([]lipgloss.Style, 3)
	screenWidth, screenHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Panic(err)
	}

	ui.panes[PositionPane] = paneStyle.Copy().
		Width(screenWidth - horizontalMargin - horizontalBorder - horizontalBorder - horizontalMargin).
		Height(screenHeight * relPositionPaneHeight / totalHeight)
	ui.panes[LogPane] = paneStyle.Copy().
		Width(screenWidth - horizontalMargin - horizontalBorder - horizontalBorder - horizontalMargin).
		Height(screenHeight * relLogPaneHeight / totalHeight)
	ui.panes[PerformancePane] = paneStyle.Copy().
		Width(screenWidth - horizontalMargin - horizontalBorder - horizontalBorder - horizontalMargin).
		Height(screenHeight * relPerformancePaneHeight / totalHeight)
	ui.activePane = PositionPane

	ui.help = help.NewModel()

	ui.paneKeys = newPanesKeyMap()
	ui.logKeys = newLogKeyMap()

	ui.db = d
	ui.positionList = newPositionListModel(ui.db,
		ui.panes[PositionPane].GetHeight()-horizontalBorder-horizontalMargin)

	ui.log = newLogModel()
	ui.log.maxHeight = ui.panes[LogPane].GetHeight() - horizontalBorder - horizontalMargin
	ui.log.Add(Debug, "Logging started")

	return ui
}

func (ui mainUI) Init() tea.Cmd {
	return nil
}

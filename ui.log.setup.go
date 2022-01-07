package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type logLevel int

const ( // TODO: move to enums
	Debug logLevel = iota
	Info
	Warning
	Error
)

type logEntry struct {
	timestamp time.Time
	level     logLevel
	entry     string
}

type logModel struct {
	logContent    []logEntry
	minShowLevel  logLevel
	showTimestamp bool
	maxHeight     int
}

func newLogModel() logModel {
	return logModel{
		showTimestamp: true,
		minShowLevel:  Debug,
	}
}

func (l logModel) Init() tea.Cmd {
	return nil
}

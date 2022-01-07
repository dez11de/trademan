package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (l logModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return l, nil
}

func (l logModel) View() string {
	var tempStr string
	shownLines := 0
	for i := len(l.logContent) - 1; (shownLines <= l.maxHeight) && (i >= 0); i-- {
		if l.logContent[i].level >= l.minShowLevel {
			tempStr = tempStr + "\n"
			if l.showTimestamp {
				tempStr = tempStr + l.logContent[i].timestamp.Format("02-01-06 15:04:05") + " "
			}
			tempStr = tempStr + l.logContent[i].entry
			shownLines++
		}
	}
	return tempStr
}

func (l *logModel) Add(level logLevel, entry string) {
	l.logContent = append(l.logContent, logEntry{
		timestamp: time.Now(),
		level:     level,
		entry:     entry,
	})
}

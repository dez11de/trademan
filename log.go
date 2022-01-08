package main

import "time"

type Log struct {
	LogID      int64
	PositionID int64
	Source     LogSource
	EntryTime  time.Time
	Text       string
}
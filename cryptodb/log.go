package cryptodb

import "time"

type Log struct {
	LogID     int64
	PlanID    int64
	Source    LogSource
	EntryTime time.Time
	Text      string
}

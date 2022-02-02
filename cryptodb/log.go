package cryptodb

import "time"

type Log struct {
	ID        uint `gorm:"autoIncrement;primaryKey"`
    PlanID    uint `gorm:"index"`
	Source    LogSource
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
	Text      string
}

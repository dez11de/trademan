package cryptodb

import "time"

type Log struct {
	ID        uint64
	PlanID    uint64    `gorm:"index"`
	Source    LogSource `gorm:"type:varchar(25)"`
	Text      string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}

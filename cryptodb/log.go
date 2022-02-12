package cryptodb

import "time"

type Log struct {
    ID        uint 
	PlanID    uint      `gorm:"index"`
	Source    LogSource `gorm:"type:varchar(25)"`
	Text      string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime;index"`
}

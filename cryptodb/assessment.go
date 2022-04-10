package cryptodb

import (
	"time"
)

type Assessment struct {
	ID                   uint64
	PlanID               uint64
	Status               string `gorm:"type:varchar(25)"`
	Risk                 string `gorm:"type:varchar(25);index"`
	Timing               string `gorm:"type:varchar(25);index"`
	StopLoss             string `gorm:"type:varchar(25);index"`
	Entry                string `gorm:"type:varchar(25);index"`
	Emotion              string `gorm:"type:varchar(25);index"`
	FollowPlan           string `gorm:"type:varchar(25);index"`
	OrderManagement      string `gorm:"type:varchar(25);index"`
	MoveStopLossInProfit string `gorm:"type:varchar(25);index"`
	TakeProfitStrategy   string `gorm:"type:varchar(25);index"`
	TakeProfitCount      string `gorm:"type:varchar(25);index"`
	Notes                string
	CreatedAt            time.Time `gorm:"index"`
	UpdatedAt            time.Time `gorm:"index"`
}

func NewAssessment(planID uint64) Assessment {
	return Assessment{
		PlanID: planID,
	}
}

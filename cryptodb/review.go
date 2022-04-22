package cryptodb

import (
	"time"
)

type Review struct {
	ID                   uint64
	PlanID               uint64
	Risk                 string `gorm:"type:varchar(25)"`
	RiskReward           string `gorm:"type:varchar(25)"`
	Timing               string `gorm:"type:varchar(25)"`
	StopLoss             string `gorm:"type:varchar(25)"`
	Entry                string `gorm:"type:varchar(25)"`
	Emotion              string `gorm:"type:varchar(25)"`
	FollowPlan           string `gorm:"type:varchar(25)"`
	OrderManagement      string `gorm:"type:varchar(25)"`
	MoveStopLossInProfit string `gorm:"type:varchar(25)"`
	TakeProfitStrategy   string `gorm:"type:varchar(25)"`
	TakeProfitCount      string `gorm:"type:varchar(25)"`
	Fee                  string `gorm:"type:varchar(25)"`
	Profit               string `gorm:"type:varchar(25)"`
	Notes                string
	CreatedAt            time.Time `gorm:"index"`
	UpdatedAt            time.Time `gorm:"index"`
}

func NewReview(planID uint64) Review {
	return Review{
		PlanID: planID,
	}
}

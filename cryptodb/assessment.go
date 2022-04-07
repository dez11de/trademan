package cryptodb

import (
	"time"
)

type Assessment struct {
	ID                   uint64
	PlanID               uint64
	Status               string `gorm:"type:varchar(25)"` // In progress, Completed
	Risk                 string `gorm:"type:varchar(25)"` // Too high, Too low, Good
	Timing               string `gorm:"type:varchar(25)"` // Late, Early, On-Time
	StopLossPosition     string `gorm:"type:varchar(25)"` // TooTight, Good, TooWide
	EntryPosition        string `gorm:"type:varchar(25)"` // TooTight, Good, TooWide
	Emotion              string `gorm:"type:varchar(25)"` // Happy, Disappointed, Neutral
	FollowPlan           string `gorm:"type:varchar(25)"` // Yes, No, Neutral
	OrderManagement      string `gorm:"type:varchar(25)"` // Good, Bad, Neutral
	MoveStopLossInProfit string `gorm:"type:varchar(25)"` // Late, Early, On-Time
	TakeProfitStrategy   string `gorm:"type:varchar(25)"` // Good, Bad, Neutral
	TakeProfitCount      string `gorm:"type:varchar(25)"` // Too High, Too Low, good
	Notes                string
	CreatedAt            time.Time `gorm:"index"`
	UpdatedAt            time.Time `gorm:"index"`
}

func NewAssessment(planID uint64) Assessment {
	return Assessment{
		PlanID:               planID,
		Status:               "In Progress",
		Risk:                 "Good",
		Timing:               "On-Time",
		StopLossPosition:     "Good",
		EntryPosition:        "Good",
		Emotion:              "Neutral",
		FollowPlan:           "Neutral",
		OrderManagement:      "Good",
		MoveStopLossInProfit: "On-Time",
		TakeProfitStrategy:   "Good",
		TakeProfitCount:      "Good",
		Notes:                "",
	}
}

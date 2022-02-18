package cryptodb

import (
	"time"

	"github.com/bart613/decimal"
)

// TODO: consider including a seperate struct for statistics such as RRR @ start, Evolved RRR,
// Break/Even, relative PnL.
type Plan struct {
	ID                 uint
	PairID             uint
	Status             Status             `gorm:"type:varchar(25);index"`
	Direction          Direction          `gorm:"type:varchar(25)"`
	Risk               decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	TakeProfitStrategy TakeProfitStrategy `gorm:"type:varchar(25)"`
	Notes              string             `gorm:"type:text"`
	TradingViewPlan    string             `gorm:"type:tinytext"`
	RewardRiskRatio    float64            `gorm:"type:float"`
	Leverage           decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	AverageEntryPrice  decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	LatestValue        decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	Profit             decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	CreatedAt          time.Time          `gorm:"index"`
	UpdatedAt          time.Time          `gorm:"index"`
}

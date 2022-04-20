package cryptodb

import (
	"time"

	"github.com/bart613/decimal"
)

type Plan struct {
	ID                 uint64
	PairID             uint64
    Status             Status             `gorm:"type:varchar(25);index"`
	Direction          Direction          `gorm:"type:varchar(25)"`
	Risk               decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	TakeProfitStrategy TakeProfitStrategy `gorm:"type:varchar(25)"`
	Notes              string             `gorm:"type:text"`
	TradingViewPlan    string             `gorm:"type:tinytext"`
	RewardRiskRatio    float64            `gorm:"type:float"`
	Leverage           decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	AverageEntryPrice  decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	Profit             decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	Fee                decimal.Decimal    `gorm:"type:decimal(20, 8)"`
	CreatedAt          time.Time          `gorm:"index"`
	UpdatedAt          time.Time          `gorm:"index"`
}

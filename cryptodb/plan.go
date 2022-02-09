package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
)

// TODO: consider including a seperate struct for statastics such as RRR @ start, Evolved RRR,
// Break/Even, relative PnL.
type Plan struct {
	ID              uint
	PairID          uint
	Status          Status          `gorm:"type:varchar(25);index"`
	Direction       Direction       `gorm:"type:varchar(25)"`
	Risk            decimal.Decimal `gorm:"type:decimal(20, 8)"`
	Notes           string          `gorm:"type:text"`
	TradingViewPlan string          `gorm:"type:tinytext"`
	RewardRiskRatio float64
	Profit          decimal.Decimal
	CreatedAt       time.Time `gorm:"index"`
	UpdatedAt       time.Time `gorm:"index"`
}

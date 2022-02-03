package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
)

type Balance struct {
	Symbol         string          `gorm:"index"`
	Equity         decimal.Decimal `gorm:"type:decimal(20, 8) "json:"equity"`
	Available      decimal.Decimal `gorm:"type:decimal(20, 8)" json:"available_balance"`
	UsedMargin     decimal.Decimal `gorm:"type:decimal(20, 8)" json:"used_margin"`
	OrderMargin    decimal.Decimal `gorm:"type:decimal(20, 8)" json:"order_margin"`
	PositionMargin decimal.Decimal `gorm:"type:decimal(20, 8)" json:"position_margin"`
	OCCClosingFee  decimal.Decimal `gorm:"type:decimal(20, 8)" json:"occ_closing_fee"`
	OCCFundingFee  decimal.Decimal `gorm:"type:decimal(20, 8)" json:"occ_funding_fee"`
	WalletBalance  decimal.Decimal `gorm:"type:decimal(20, 8)" json:"wallet_balance"`
	DailyPnL       decimal.Decimal `gorm:"type:decimal(20, 8)" gorm:"column:daily_pnl" json:"realised_pnl"`
	UnrealisedPnL  decimal.Decimal `gorm:"column:unrealised_pnl;type:decimal(20, 8)" json:"unrealised_pnl"`
	TotalPnL       decimal.Decimal `gorm:"column:total_pnl;type:decimal(20, 8)" json:"cum_realised_pnl"`
	CreatedAt      time.Time       `gorm:"index"`
}

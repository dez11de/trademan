package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
)

type Balance struct {
	Symbol         string
	Equity         decimal.Decimal `json:"equity"`
	Available      decimal.Decimal `json:"available_balance"`
	UsedMargin     decimal.Decimal `json:"used_margin"`
	OrderMargin    decimal.Decimal `json:"order_margin"`
	PositionMargin decimal.Decimal `json:"position_margin"`
	OCCClosingFee  decimal.Decimal `json:"occ_closing_fee"`
	OCCFundingFee  decimal.Decimal `json:"occ_funding_fee"`
	WalletBalance  decimal.Decimal `json:"wallet_balance"`
	DailyPnL       decimal.Decimal `json:"realised_pnl"`
	UnrealisedPnL  decimal.Decimal `json:"unrealised_pnl"`
	TotalPnL       decimal.Decimal `json:"cum_realised_pnl"`
	EntryTime      time.Time
}

package main

import "time"

type balance struct {
	Symbol         string
	Equity         float64 `json:"equity"`
	Available      float64 `json:"available_balance"`
	UsedMargin     float64 `json:"used_margin"`
	OrderMargin    float64 `json:"order_margin"`
	PositionMargin float64 `json:"position_margin"`
	OCCClosingFee  float64 `json:"occ_closing_fee"`
	OCCFundingFee  float64 `json:"occ_funding_fee"`
	WalletBalance  float64 `json:"wallet_balance"`
	DailyPnL       float64 `json:"realised_pnl"`
	UnrealisedPnL  float64 `json:"unrealised_pnl"`
	TotalPnL       float64 `json:"cum_realised_pnl"`
	EntryTime      time.Time
}

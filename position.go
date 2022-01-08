package main

import "time"

type Position struct {
	PositionID      int64
	Status          Status
	PairID          int64
	Pair            string
	Size            float64 `json:"size"`
	Side            Status  `json:"side"`
	Risk            float64
	EntryPrice      float64 `json:"entry_price"`
	HardStopLoss    float64 `json:"stop_loss"`
	Notes           string
	TradingViewPlan string
	RewardRiskRatio float64
	Profit          float64 `json:"realised_pnl"`
	EntryTime       time.Time
	ModifyTime      time.Time
}

package main

import (
	"time"

	"github.com/shopspring/decimal"
)

type Plan struct {
	PlanID          int64
	PairID          int64
	Status          Status
	Side            Side `json:"side"`
	Risk            float64
	Notes           string
	TradingViewPlan string
	RewardRiskRatio float64
	Profit          decimal.Decimal `json:"realised_pnl"`
	EntryTime       time.Time
	ModifyTime      time.Time
}

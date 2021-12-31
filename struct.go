package main

import (
	"time"
)

type Position struct {
	PositionID      int64
	Status          Status
	Symbol          string  `json:"symbol"`
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

type Order struct {
	OrderID         int64
	PositionID      int64
	ExchangeOrderID string `json:"order_link_id"`
	Status          Status
	OrderType       OrderType `json:"order_type"`
	Size            float64   `json:"qty"`
	TriggerPrice    float64   `json:"tp_trigger"`
	Price           float64   `json:"price"`
	EntryTime       time.Time
	ModifyTime      time.Time
}

type Log struct {
	LogID      int64
	PositionID int64
	Source     LogSource
	EntryTime  time.Time
	Text       string
}

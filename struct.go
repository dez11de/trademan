package main

import (
	"fmt"
	"time"
)

type TradeDirection int

const (
	long TradeDirection = iota
	short
)

// TODO: apparently this sort of stuff can be done using the stringer package.
// Or apparently even better with https://github.com/alvaroloes/enumer
func (td TradeDirection) String() string {
	switch td {
	case long:
		return "Long"
	case short:
		return "Short"
	default:
		return fmt.Sprintf("UNDEFINED OrderDirection(%d)", td)
	}
}

type Status int

const (
	planned Status = iota
	ordered
	position
	stopped
	closed
	cancelled
	liquidated
	logged
)

func (cs Status) String() string {
	switch cs {
	case planned:
		return "Planned"
	case ordered:
		return "Ordered"
	case position:
		return "Position"
	case stopped:
		return "Stopped"
	case closed:
		return "Closed"
	case cancelled:
		return "Cancelled"
	case liquidated:
		return "Liquidated"
	case logged:
		return "Logged"
	default:
		return fmt.Sprintf("UNDEFINED TradeState(%d)", cs)
	}
}

type Position struct {
	TradeID         int64
	Symbol          string `json:"symbol"`
	Side            Status `json:"side"`
	Risk            float64
	Size            float64 `json:"size"`
	EntryPrice      float64 `json:"entry_price"`
	HardStopLoss    float64 `json:"stop_loss"`
	Notes           string
	TradingViewPlan string
	RewardRiskRatio float64
	Profit          float64 `json:"realised_pnl"`
	EntryTime       time.Time
	ModifyTime      time.Time
}

type OrderType int

const (
	SoftStopLoss OrderType = iota
	TakeProfit
)

func (ot OrderType) String() string {
	switch ot {
	case SoftStopLoss:
		return "Soft StopLoss"
	case TakeProfit:
		return "Take Profit"
	default:
		return fmt.Sprintf("UNDEFINED OrderType(%d)", ot)
	}
}

type Order struct {
	OrderID           int64
	TradeID           int64  // foreign key
	ExchangeOrderID   string `json:"order_link_id"`
	Status            Status
	OrderType         OrderType `json:"order_type"`
	Quantity          float64   `json:"qty"`
	TakeProfitTrigger float64   `json:"tp_trigger"`
	Price             float64   `json:"price"`
	EntryTime         time.Time
	ModifyTime        time.Time
}

type LogSource int

const (
	trigger LogSource = iota
	software
	user
)

func (s LogSource) String() string {
	switch s {
	case trigger:
		return "Trigger"
	case software:
		return "Software"
	case user:
		return "User"
	default:
		return fmt.Sprintf("UNDEFINED LogSource(%d)", s)
	}
}

type ExecutionLog struct {
	LogID     int64
	TradeID   int64
	Source    LogSource
	EntryTime time.Time
	Text      string
}

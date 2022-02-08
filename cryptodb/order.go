package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID              uint
	PlanID          uint `gorm:"index"`
	Status          Status
	ExchangeOrderID string          `json:"order_link_id"`
    OrderKind       OrderKind
	OrderType       OrderType       `json:"order_type"`
	Size            decimal.Decimal `gorm:"type:decimal(20, 8)" json:"qty"`
	TriggerPrice    decimal.Decimal `gorm:"type:decimal(20, 8)" json:"tp_trigger"`
	Price           decimal.Decimal `gorm:"type:decimal(20, 8)" json:"price"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// const MaxTakeProfits = 5

func NewOrders(PlanID uint) []Order {
	return []Order{
		{PlanID: PlanID, OrderKind: KindHardStopLoss},
		{PlanID: PlanID, OrderKind: KindSoftStopLoss},
		{PlanID: PlanID, OrderKind: KindEntry},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
	}
}

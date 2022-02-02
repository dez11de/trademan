package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID              uint
	PlanID          uint
	Status          Status
	ExchangeOrderID string          `json:"order_link_id"`
	OrderType       OrderType       `json:"order_type"`
	Size            decimal.Decimal `gorm:"type:decimal(20, 8)" json:"qty"`
	TriggerPrice    decimal.Decimal `gorm:"type:decimal(20, 8)" json:"tp_trigger"`
	Price           decimal.Decimal `gorm:"type:decimal(20, 8)" json:"price"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

const MaxTakeProfits = 5

func NewOrders(PlanID uint) []Order {
	return []Order{
        {PlanID: PlanID, OrderType: TypeHardStopLoss},
        {PlanID: PlanID, OrderType: TypeSoftStopLoss},
        {PlanID: PlanID, OrderType: TypeEntry},
        {PlanID: PlanID, OrderType: TypeTakeProfit},
        {PlanID: PlanID, OrderType: TypeTakeProfit},
        {PlanID: PlanID, OrderType: TypeTakeProfit},
        {PlanID: PlanID, OrderType: TypeTakeProfit},
        {PlanID: PlanID, OrderType: TypeTakeProfit},
	}
}

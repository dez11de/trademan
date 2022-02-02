package cryptodb

import (
	//	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	//	ID              uint `gorm:"primaryKey"`
	PlanID          uint //`gorm:"foreignKey:PlanID"`
	Status          Status
	ExchangeOrderID string          `json:"order_link_id"`
	OrderType       OrderType       `json:"order_type"`
	Size            decimal.Decimal `gorm:"type:decimal(20, 8)" json:"qty"`
	TriggerPrice    decimal.Decimal `gorm:"type:decimal(20, 8)" json:"tp_trigger"`
	Price           decimal.Decimal `gorm:"type:decimal(20, 8)" json:"price"`
	//	CreatedAt       time.Time
	//	UpdatedAt       time.Time
}

/*
const MaxTakeProfits = 5

func NewOrders() []Order {
	return []Order{
		{OrderType: TypeHardStopLoss},
		{OrderType: TypeSoftStopLoss},
		{OrderType: TypeEntry},
		{OrderType: TypeTakeProfit},
		{OrderType: TypeTakeProfit},
		{OrderType: TypeTakeProfit},
		{OrderType: TypeTakeProfit},
		{OrderType: TypeTakeProfit},
	}
}
*/

package cryptodb

import (
	"fmt"
	"time"

	"github.com/bart613/decimal"
)

type Order struct {
	ID              uint
	PlanID          uint            `gorm:"index"`
	Status          Status          `gorm:"type:varchar(25)"`
    ExchangeOrderID string          `gorm:"type:varchar(36)" json:"order_link_id"`
	OrderKind       OrderKind       `gorm:"type:varchar(25)"`
	OrderType       OrderType       `gorm:"-" json:"order_type"`
	Size            decimal.Decimal `gorm:"type:decimal(20, 8)" json:"qty"`
	// TriggerPrice    decimal.Decimal `gorm:"-" json:"tp_trigger"`
	Price           decimal.Decimal `gorm:"type:decimal(20, 8)" json:"price"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// const MaxTakeProfits = 5

func NewOrders(PlanID uint) []Order {
	return []Order{
		{PlanID: PlanID, OrderKind: KindMarketStopLoss},
		{PlanID: PlanID, OrderKind: KindLimitStopLoss},
		{PlanID: PlanID, OrderKind: KindEntry},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
		{PlanID: PlanID, OrderKind: KindTakeProfit},
	}
}

func debugPrintOrders(orders []Order) {
    for i, o := range orders {
        fmt.Printf("[%d] %d %s %s %s\n", i, o.ID, o.OrderKind.String(), o.Size.StringFixed(5), o.Price.StringFixed(5))
    }
}


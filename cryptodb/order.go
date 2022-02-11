package cryptodb

import (
	"fmt"
	"time"

	"github.com/bart613/decimal"
)

type Order struct {
	ID              uint
	PlanID          uint            `gorm:"index"`
	Status          Status          `gorm:"type:varchar(25)"` // TODO: investigate option to define enums, see: https://github.com/go-gorm/gorm/issues/1978#issuecomment-476673540
	ExchangeOrderID string          `gorm:"type:varchar(36);index" json:"order_link_id"`
	OrderKind       OrderKind       `gorm:"type:varchar(25)"`
	Size            decimal.Decimal `gorm:"type:decimal(20, 8)" json:"qty"`
	Price           decimal.Decimal `gorm:"type:decimal(20, 8)" json:"price"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

const MaxTakeProfits = 5

func NewOrders(PlanID uint) []Order {
	return []Order{
		{PlanID: PlanID, Status: StatusPlanned, OrderKind: KindMarketStopLoss},
		{PlanID: PlanID, Status: StatusPlanned, OrderKind: KindLimitStopLoss},
		{PlanID: PlanID, Status: StatusPlanned, OrderKind: KindEntry},
		{PlanID: PlanID, Status: StatusPlanned, OrderKind: KindTakeProfit},
		{PlanID: PlanID, Status: StatusPlanned, OrderKind: KindTakeProfit},
		{PlanID: PlanID, Status: StatusPlanned, OrderKind: KindTakeProfit},
		{PlanID: PlanID, Status: StatusPlanned, OrderKind: KindTakeProfit},
		{PlanID: PlanID, Status: StatusPlanned, OrderKind: KindTakeProfit},
	}
}

func debugPrintOrders(orders []Order) {
	for i, o := range orders {
		fmt.Printf("[%d] %d %s %s %s\n", i, o.ID, o.OrderKind.String(), o.Size.StringFixed(5), o.Price.StringFixed(5))
	}
}

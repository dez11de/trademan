package cryptodb

import (
	"time"

	"github.com/bart613/decimal"
)

type Order struct {
	ID            uint
	PlanID        uint            `gorm:"index;not null"`
	Status        Status          `gorm:"type:varchar(25)"` // TODO: investigate option to define enums, see: https://github.com/go-gorm/gorm/issues/1978#issuecomment-476673540
	OrderKind     OrderKind       `gorm:"type:varchar(25)"`
	Size          decimal.Decimal `gorm:"type:decimal(20, 8)" json:"qty"`
	Price         decimal.Decimal `gorm:"type:decimal(20, 8)" json:"price"`
	LinkOrderID   string          `gorm:"type:varchar(36);index" json:"order_link_id"` // Links trademan to exchange
	SystemOrderID string          `gorm:"type:varchar(36);index" json:"order_id"`      // Stores either ByBits order_id, or stop_order_id
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

const MaxTakeProfits = 5

func NewOrders(PlanID uint) []Order {
	return []Order{
		{PlanID: PlanID, Status: Unplanned, OrderKind: MarketStopLoss},
		{PlanID: PlanID, Status: Unplanned, OrderKind: LimitStopLoss},
		{PlanID: PlanID, Status: Unplanned, OrderKind: Entry},
		{PlanID: PlanID, Status: Unplanned, OrderKind: TakeProfit},
		{PlanID: PlanID, Status: Unplanned, OrderKind: TakeProfit},
		{PlanID: PlanID, Status: Unplanned, OrderKind: TakeProfit},
		{PlanID: PlanID, Status: Unplanned, OrderKind: TakeProfit},
		{PlanID: PlanID, Status: Unplanned, OrderKind: TakeProfit},
	}
}

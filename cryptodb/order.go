package cryptodb

import (
	"time"

	"github.com/bart613/decimal"
)

type Order struct {
	ID            uint64
	PlanID        uint64          `gorm:"index;not null"`
	Status        Status          `gorm:"type:varchar(25)"` // TODO: investigate option to define enums, see: https://github.com/go-gorm/gorm/issues/1978#issuecomment-476673540
	OrderKind     OrderKind       `gorm:"type:varchar(25)"`
	Size          decimal.Decimal `gorm:"type:decimal(20, 8)" json:"qty"`
	Price         decimal.Decimal `gorm:"type:decimal(20, 8)" json:"price"`
	TriggerPrice  decimal.Decimal `gorm:"type:decimal(20, 8)"`
	// LinkOrderID   string          `gorm:"type:varchar(36);index" json:"order_link_id"`
	SystemOrderID string          `gorm:"type:varchar(36);index" json:"order_id"`    
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

const MaxTakeProfits = 5

func NewOrders(PlanID uint64) []Order {
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

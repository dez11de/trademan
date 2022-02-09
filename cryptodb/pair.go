package cryptodb

import (
	"time"

	"github.com/bart613/decimal"
)

type Pair struct {
	ID            uint
	Name          string          `gorm:"uniqueIndex; size:15; unique; not null" json:"name"`
	Alias         string          `gorm:"size:15" json:"alias"`
	Status        string          `gorm:"size:15" json:"status"` // TODO: make enum
	BaseCurrency  string          `gorm:"size:10" json:"base_currency"`
	QuoteCurrency string          `gorm:"size:10" json:"quote_currency"`
	PriceScale    int32           `json:"price_scale,number"`
	TakerFee      decimal.Decimal `gorm:"type:decimal(20, 8)" json:"taker_fee,string"` //TODO: investigate right size
	MakerFee      decimal.Decimal `gorm:"type:decimal(20, 8)" json:"maker_fee,string"`

	Leverage struct {
		Min  decimal.Decimal `gorm:"type:decimal(20, 8)" json:"min_leverage"`
		Max  decimal.Decimal `gorm:"type:decimal(20, 8)" json:"max_leverage"`
		Step decimal.Decimal `gorm:"type:decimal(20, 8)" json:"leverage_step,string"`
    } `gorm:"embedded;embeddedPrefix:leverage_" json:"leverage_filter"` // TODO: is embedding still necesarry?

	Price struct {
		Min  decimal.Decimal `gorm:"type:decimal(20, 8)" json:"min_price,string"`
		Max  decimal.Decimal `gorm:"type:decimal(20, 8)" json:"max_price,string"`
		Tick decimal.Decimal `gorm:"type:decimal(20, 8)" json:"tick_size,string"`
	} `gorm:"embedded;embeddedPrefix:price_" json:"price_filter"`

	Order struct {
		Min  decimal.Decimal `gorm:"type:decimal(20, 8)" json:"min_trading_qty,number"`
		Max  decimal.Decimal `gorm:"type:decimal(20, 8)" json:"max_trading_qty,number"`
		Step decimal.Decimal `gorm:"type:decimal(20, 8)" json:"qty_step,number"`
	} `gorm:"embedded;embeddedPrefix:order_" json:"lot_size_filter"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

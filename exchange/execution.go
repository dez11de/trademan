package exchange

import (
	"time"

	"github.com/dez11de/cryptodb"
	"github.com/shopspring/decimal"
)

type Execution struct {
	Symbol         string          `json:"symbol"`
	Side           cryptodb.Side   `json:"side"`
	OrderID        uint            `json:"order_id,string"`
	ExecID         uint            `json:"exec_id,string"`
	OrderLinkID    string          `json:"order_link_id"`
	Price          decimal.Decimal `json:"price"`
	OrderQuantity  decimal.Decimal `json:"order_qty"`
	ExecType       string          `json:"exec_type"`
	ExecQuantity   decimal.Decimal `json:"exec_qty"`
	ExecFee        decimal.Decimal `json:"exec_fee,string"`
	LeavesQuantity decimal.Decimal `json:"leaves_qty,string"`
	IsMaker        bool            `json:"is_maker"`
	TradeTimestamp time.Time       `json:"trade_time"`
}

package exchange

import "time"

type Execution struct {
	Symbol         string    `json:"symbol"`
	Side           string    `json:"side"`
	OrderID        string    `json:"order_id"`
	ExecID         string    `json:"exec_id"`
	OrderLinkID    string    `json:"order_link_id"`
	Price          float64   `json:"price"`
	OrderQuantity  float64   `json:"order_qty"`
	ExecType       string    `json:"exec_type"`
	ExecQuantity   float64   `json:"exec_qty"`
	ExecFee        float64   `json:"exec_fee"`
	LeavesQuantity float64   `json:"leaves_qty"`
	IsMaker        bool      `json:"is_maker"`
	TradeTimestamp time.Time `json:"trade_time"`
}

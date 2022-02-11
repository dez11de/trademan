package exchange

import (
	"time"

	"github.com/bart613/decimal"
)

type Order struct {
	OrderID         string          `json:"order_id"` // order_id and stop_order_id are the same but somehow different
	StopOrderID     string          `json:"stop_order_id"`
	ExchangeOrderID string          `json:"order_link_id"`
	Symbol          string          `json:"symbol"`
	Side            string          `json:"side"`       // Buy/Sell
	OrderType       string          `json:"order_type"` // Limit/Market
	Price           decimal.Decimal `json:"price"`
	Size            decimal.Decimal `json:"qty"`
	StopLoss        decimal.Decimal `json:"stop_loss"` // bybit api-documentation is ambigious
	Fees            decimal.Decimal `json:"cum_exec_fee"`
	OrderStatus     string          `json:"order_status"`        // Created/Rejected/New/PartiallyFilled/Filled/Cancelled/PendingCancel/Untriggered (Untriggered is missing from api-documentation)
	PositionIndex   int             `json:"position_idx,string"` // 0:One-Way Mode, 1:Buy side of both side mode, 2: Sell side of both side mode
	CreatedAt       time.Time       `json:"create_time,string"`
	UpdatedAt       time.Time       `json:"update_time,string"`
}

type OrderPage struct {
	CurrentPage int     `json:"current_page"`
	LastPage    int     `json:"last_page"`
	Orders      []Order `json:"data"`
}

type OrderResponse struct {
	ReturnCode       int       `json:"ret_code"`
	ReturnMessage    string    `json:"ret_msg"`
	ExtendedCode     string    `json:"ext_code"`
	Results          OrderPage `json:"result"`
	ExtendedInfo     string    `json:"ext_info"`
	ServerTime       string    `json:"time_now"`
	RateLimitStatus  int       `json:"rate_limit_status"`
	RateLimitResetMS int       `json:"rate_limit_reset_ms"`
	RateLimit        int       `json:"rate_limit"`
}

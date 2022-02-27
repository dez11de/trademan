package exchange

import (
	"time"

	"github.com/bart613/decimal"
)

type Order struct {
	OrderID     string          `json:"order_id"`
	StopOrderID string          `json:"stop_order_id"`
	LinkOrderID string          `json:"order_link_id"`
	Symbol      string          `json:"symbol"`
	Side        string          `json:"side"`         // Buy/Sell
	OrderType   string          `json:"order_type"`   // Limit/Market
	OrderStatus string          `json:"order_status"` // Created/Rejected/New/PartiallyFilled/Filled/Cancelled/PendingCancel/Untriggered
	Price       decimal.Decimal `json:"price"`
	Size        decimal.Decimal `json:"qty"`
	Leaves      decimal.Decimal `json:"leaves_qty"`
	StopLoss    decimal.Decimal `json:"stop_loss"`
	Fees        decimal.Decimal `json:"cum_exec_fee"`
	CreatedAt   time.Time       `json:"create_time,string"`
	UpdatedAt   time.Time       `json:"update_time,string"`
}

type OrderResponseRest struct {
	ReturnCode       int64  `json:"ret_code"`
	ReturnMessage    string `json:"ret_msg"`
	ExtendedCode     string `json:"ext_code"`
	Result           Order  `json:"result"`
	ExtendedInfo     string `json:"ext_info"`
	ServerTime       string `json:"time_now"`
	RateLimitStatus  int64  `json:"rate_limit_status"`
	RateLimitResetMS int64  `json:"rate_limit_reset_ms"`
	RateLimit        int64  `json:"rate_limit"`
}

type OrderPage struct {
	CurrentPage int64   `json:"current_page"`
	LastPage    int64   `json:"last_page"`
	Orders      []Order `json:"data"`
}

type OrderResponseWS struct {
	ReturnCode       int64     `json:"ret_code"`
	ReturnMessage    string    `json:"ret_msg"`
	ExtendedCode     string    `json:"ext_code"`
	Results          OrderPage `json:"result"`
	ExtendedInfo     string    `json:"ext_info"`
	ServerTime       string    `json:"time_now"`
	RateLimitStatus  int64     `json:"rate_limit_status"`
	RateLimitResetMS int64     `json:"rate_limit_reset_ms"`
	RateLimit        int64     `json:"rate_limit"`
}

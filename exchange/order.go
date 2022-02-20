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
	Side        string          `json:"side"`       // Buy/Sell
	OrderType   string          `json:"order_type"` // Limit/Market
	Price       decimal.Decimal `json:"price"`
	Size        decimal.Decimal `json:"qty"`
	Leaves      decimal.Decimal `json:"leaves_qty"`
	StopLoss    decimal.Decimal `json:"stop_loss"`
	Fees        decimal.Decimal `json:"cum_exec_fee"`
	OrderStatus string          `json:"order_status"` // Created/Rejected/New/PartiallyFilled/Filled/Cancelled/PendingCancel/Untriggered
	// PositionIndex int             `json:"position_idx"` // sometimes it's a string, sometimes it's a number, sometimes an integer. Bybit Documentation :') '0:One-Way Mode, 1:Buy side of both side mode, 2: Sell side of both side mode
	CreatedAt time.Time `json:"create_time,string"`
	UpdatedAt time.Time `json:"update_time,string"`
}

type OrderResponseRest struct {
	ReturnCode       int    `json:"ret_code"`
	ReturnMessage    string `json:"ret_msg"`
	ExtendedCode     string `json:"ext_code"`
	Result           Order  `json:"result"`
	ExtendedInfo     string `json:"ext_info"`
	ServerTime       string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMS int    `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}

type OrderPage struct {
	CurrentPage int     `json:"current_page"`
	LastPage    int     `json:"last_page"`
	Orders      []Order `json:"data"`
}

type OrderResponseWS struct {
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

package exchange

import "github.com/bart613/decimal"

type Position struct {
	UserID              int             `json:"user_id,string"`
	Symbol              string          `json:"symbol"`
	Side                string          `json:"side"`
	Size                decimal.Decimal `json:"size"`
	PositionValue       decimal.Decimal `json:"position_value"`
	EntryPrice          decimal.Decimal `json:"entry_price"`
	LiquidationPrice    decimal.Decimal `json:"liq_price"`
	BustPrice           decimal.Decimal `json:"bust_price"`
	Leverage            decimal.Decimal `json:"leverage"`
	Isolated            bool            `json:"is_isolated"`
	AutoAddMargin       decimal.Decimal `json:"auto_add_margin,string"`
	PositionMargin      decimal.Decimal `json:"position_margin"`
	OCCClosingFee       decimal.Decimal `json:"occ_closing_fee"`
	RealisedPnL         decimal.Decimal `json:"realised_pnl"`
	CumalitivePnL       decimal.Decimal `json:"cum_realised_pnl"`
	FreeQuantity        decimal.Decimal `json:"free_qty"`
	TPSLMode            string          `json:"tp_sl_mode"`
	UnrealisedPnL       decimal.Decimal `json:"unrealised_pnl"`
	DeleverageIndicator decimal.Decimal `json:"deleverage_indicator"`
	RiskID              int             `json:"risk_id,string"`
	StopLoss            decimal.Decimal `json:"stop_loss,string"`
	TakeProfit          decimal.Decimal `json:"take_profit,string"`
	TrailingStop        decimal.Decimal `json:"trailing_stop,string"`
}

type PositionResponse struct {
	ReturnCode       int        `json:"ret_code"`
	ReturnMessage    string     `json:"ret_msg"`
	ExtendedCode     string     `json:"ext_code"`
	Results          []Position `json:"result"`
	ExtendedInfo     string     `json:"ext_info"`
	ServerTime       string     `json:"time_now,string"`
	RateLimitStatus  int        `json:"rate_limit_status"`
	RateLimitResetMS int        `json:"rate_limit_reset_ms"`
	RateLimit        int        `json:"rate_limit"`
}

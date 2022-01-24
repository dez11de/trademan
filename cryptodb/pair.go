package cryptodb

import "github.com/shopspring/decimal"

type pairLeverage struct {
	Min  decimal.Decimal `json:"min_leverage"`
	Max  decimal.Decimal `json:"max_leverage"`
	Step decimal.Decimal `json:"leverage_step,string"`
}

type pairPrice struct {
	Min  decimal.Decimal `json:"min_price,string"`
	Max  decimal.Decimal `json:"max_price,string"`
	Tick decimal.Decimal `json:"tick_size,string"`
}

type pairOrderSize struct {
	Min  decimal.Decimal `json:"min_trading_qty"`
	Max  decimal.Decimal `json:"max_trading_qty"`
	Step decimal.Decimal `json:"qty_step"`
}

type Pair struct {
	PairID        int64
	Pair          string          `json:"name"`
	BaseCurrency  string          `json:"base_currency"`
	QuoteCurrency string          `json:"quote_currency"`
	PriceScale    int32           `json:"price_scale"`
	TakerFee      decimal.Decimal `json:"taker_fee,string"`
	MakerFee      decimal.Decimal `json:"maker_fee,string"`
	Leverage      pairLeverage    `json:"leverage_filter"`
	Price         pairPrice       `json:"price_filter"`
	OrderSize     pairOrderSize   `json:"lot_size_filter"`
}

type pairResponse struct {
	ReturnCode       int    `json:"ret_code"`
	ReturnMessage    string `json:"ret_msg"`
	ExtendedCode     string `json:"ext_code"`
	Results          []Pair `json:"result"`
	ExtendedInfo     string `json:"ext_info"`
	ServerTime       string `json:"time_now,string"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMS int    `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}


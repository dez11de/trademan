package main

type symbolLeverage struct {
	Min  float64 `json:"min_leverage"`
	Max  float64 `json:"max_leverage"`
	Step float64 `json:"leverage_step,string"`
}

type symbolPrice struct {
	Min  float64 `json:"min_price,string"`
	Max  float64 `json:"max_price,string"`
	Tick float64 `json:"tick_size,string"`
}

type symbolOrderSize struct {
	Min  float64 `json:"min_trading_qty"`
	Max  float64 `json:"max_trading_qty"`
	Step float64 `json:"qty_step"`
}

type symbol struct {
	SymbolID      int64
	Symbol        string          `json:"name"`
	BaseCurrency  string          `json:"base_currency"`
	QuoteCurrency string          `json:"quote_currency"`
	PriceScale    int             `json:"price_scale"`
	TakerFee      float64         `json:"taker_fee,string"`
	MakerFee      float64         `json:"maker_fee,string"`
	Leverage      symbolLeverage  `json:"leverage_filter"`
	Price         symbolPrice     `json:"price_filter"`
	OrderSize     symbolOrderSize `json:"lot_size_filter"`
}

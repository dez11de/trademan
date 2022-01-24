package cryptoDB

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	OrderID         int64
	PlanID          int64
	Status          Status
	ExchangeOrderID string          `json:"order_link_id"`
	OrderType       OrderType       `json:"order_type"`
	Size            decimal.Decimal `json:"qty"`
	TriggerPrice    decimal.Decimal `json:"tp_trigger"`
	Price           decimal.Decimal `json:"price"`
	EntryTime       time.Time
	ModifyTime      time.Time
}

const MaxTakeProfits = 5

type Orders [3 + MaxTakeProfits]Order

type OrderPage struct {
	CurrentPage int     `json:"current_page"`
	LastPage    int     `json:"last_page"`
	OrderList   []Order `json:"data"`
}

type OrderResponse struct {
	ReturnCode       int       `json:"ret_code"`
	ReturnMessage    string    `json:"ret_msg"`
	ExtendedCode     string    `json:"ext_code"`
	Results          OrderPage `json:"result"`
	ExtendedInfo     string    `json:"ext_info"`
	ServerTime       string    `json:"time_now,string"`
	RateLimitStatus  int       `json:"rate_limit_status"`
	RateLimitResetMS int       `json:"rate_limit_reset_ms"`
	RateLimit        int       `json:"rate_limit"`
}

func NewOrders() Orders {
	return Orders{
        {OrderType: TypeHardStopLoss},
		{OrderType: TypeSoftStopLoss},
		{OrderType: TypeEntry},
		{OrderType: TypeTakeProfit},
		{OrderType: TypeTakeProfit},
		{OrderType: TypeTakeProfit},
		{OrderType: TypeTakeProfit},
		{OrderType: TypeTakeProfit},
	}
}

/*
func (b *ByBit) GetOrder(symbol string, status string) Order {
	apiURL := "/private/linear/order/list"
	params := make(map[string]interface{})
	params["symbol"] = symbol
	if status != "" {
		params["order_status"] = status
	}
	var or OrderResponse
	b.PrivateRequest("GET", apiURL, params, &or)
	if len(or.Results.OrderList) > 1 || len(or.Results.OrderList) == 0 {
		log.Printf("More orders then expected.")
		return Order{}
	} else {
		return or.Results.OrderList[0]
	}
}
*/

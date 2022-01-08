package main

import "time"

type Order struct {
	OrderID         int64
	PositionID      int64
	ExchangeOrderID string `json:"order_link_id"`
	Status          Status
	OrderType       OrderType `json:"order_type"`
	Size            float64   `json:"qty"`
	TriggerPrice    float64   `json:"tp_trigger"`
	Price           float64   `json:"price"`
	EntryTime       time.Time
	ModifyTime      time.Time
}

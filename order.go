package main

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	OrderID         int64
	PlanID          int64
	ExchangeOrderID string `json:"order_link_id"`
	Status          Status
	OrderType       OrderType       `json:"order_type"`
	Size            decimal.Decimal `json:"qty"`
	TriggerPrice    decimal.Decimal `json:"tp_trigger"`
	Price           decimal.Decimal `json:"price"`
	EntryTime       time.Time
	ModifyTime      time.Time
}

type Orders []Order

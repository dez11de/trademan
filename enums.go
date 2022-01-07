package main

type Side int

const (
	Long Side = iota
	Short
)

type Status int

const (
	Planned Status = iota
	Ordered
	Filled
	Stopped
	Closed
	Cancelled
	Liquidated
	Logged
)

type OrderType int

const (
	SoftStopLoss OrderType = iota
	TakeProfit
)

type LogSource int

const (
	Trigger LogSource = iota
	Software
	User
)

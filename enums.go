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
	HardStopLoss OrderType = iota
	SoftStopLoss
	Entry
	TakeProfit
)

type LogSource int

const (
	Trigger LogSource = iota
	Software
	User
)

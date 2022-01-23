package main

type Side int

const (
	sideLong Side = iota
	sideShort
)

type Status int

const (
	statusPlanned Status = iota
	statusOrdered
	statusFilled
	statusStopped
	statusClosed
	statusCancelled
	statusLiquidated
	statusLogged
)

type OrderType int

const (
	typeHardStopLoss OrderType = iota
	typeSoftStopLoss
	typeEntry
	typeTakeProfit
)

type LogSource int

const (
	sourceTrigger LogSource = iota
	sourceSoftware
	sourceUser
)

package cryptodb

//go:generate enumer -json -sql -type Direction,Side,Status,OrderType,OrderKind,LogSource -output enums_helpers.go

type Direction int

const (
	DirectionLong Direction = iota
	DirectionShort
)

type Side int

const (
	SideBuy Side = iota
	SideSell
)

type Status int

const (
	StatusPlanned Status = iota
	StatusOrdered
    StatusUntriggered
	StatusFilled
	StatusStopped
	StatusClosed
	StatusCancelled
	StatusLiquidated
	StatusLogged
)

type OrderType int

const (
	Market OrderType = iota
	Limit
)

type OrderKind int

const (
	KindMarketStopLoss OrderKind = iota
	KindLimitStopLoss
	KindEntry
	KindTakeProfit
)

type LogSource int

const (
	SourceExchange LogSource = iota
	SourceSoftware
	SourceUser
)

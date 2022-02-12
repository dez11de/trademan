package cryptodb

//go:generate enumer -json -sql -type Direction,Side,TakeProfitStrategy,Status,OrderType,OrderKind,LogSource -output enums_helpers.go

type Direction int

const (
	Long Direction = iota
	Short
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

type TakeProfitStrategy int

const (
	Manual TakeProfitStrategy = iota // NOT implemented yet
	AutoLinear
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

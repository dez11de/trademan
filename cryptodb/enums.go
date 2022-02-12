package cryptodb

//go:generate enumer -json -sql -type Direction,Side,TakeProfitStrategy,Status,OrderType,OrderKind,LogSource -output enums_helpers.go

type Direction int

const (
	Long Direction = iota
	Short
)

type Side int

const (
	Buy Side = iota
	Sell
)

type Status int

const (
	Planned Status = iota
	Ordered
	Untriggered
	Filled
	Stopped
	Cancelled
	Closed
	Liquidated
	Logged
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
	MarketStopLoss OrderKind = iota
	LimitStopLoss
	Entry
	TakeProfit
)

type LogSource int

const (
	Exchange LogSource = iota
	Server
	User
)

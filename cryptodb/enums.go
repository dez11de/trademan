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
    Unplanned Status = iota
	Planned
    Created
	Ordered
	New         
	Untriggered
    PartiallyFilled
	Filled
    Rejected
    Error
	Stopped
    PendingCancel
	Cancelled
	Deactivated
	Closed 
	Liquidated
	Logged
)

type TakeProfitStrategy int

const (
	Manual TakeProfitStrategy = iota // NOT implemented yet, this requires an extra field in orders, don't reuse Size
	AutoLinear
    // TODO: for auto-rejection trading use Fibonacci retracement, see https://www.investopedia.com/terms/f/fibonacciretracement.asp for values
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

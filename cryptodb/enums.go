package cryptoDB

type Side int

const (
	SideLong Side = iota
	SideShort
)

type Status int

const (
	StatusPlanned Status = iota
	StatusOrdered
	StatusFilled
	StatusStopped
	StatusClosed
	StatusCancelled
	StatusLiquidated
	StatusLogged
)

type OrderType int

const (
	TypeHardStopLoss OrderType = iota
	TypeSoftStopLoss
	TypeEntry
	TypeTakeProfit
)

type LogSource int

const (
	SourceTrigger LogSource = iota
	SourceSoftware
	SourceUser
)

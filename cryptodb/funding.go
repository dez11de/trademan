package cryptodb

import (
	"time"

	"github.com/bart613/decimal"
)

type Funding struct {
	Pair        string          `json:"symbol"`
	Side        string          `json:"side"`
	Size        decimal.Decimal `json:"size"`
	FundingRate decimal.Decimal `json:"funding_rate"`
	ExecFee     decimal.Decimal `json:"exec_fee"`
	Timestamp   time.Time       `json:"exec_time,string"`
}

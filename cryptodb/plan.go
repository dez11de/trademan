package cryptodb

import (
	"time"

	"github.com/shopspring/decimal"
)

type Plan struct {
	PlanID          int64
	PairID          int64
	Status          Status
	Side            Side `json:"side"`
	Risk            decimal.Decimal
	Notes           string
	TradingViewPlan string
	RewardRiskRatio float64
	Profit          decimal.Decimal `json:"realised_pnl"`
	EntryTime       time.Time
	ModifyTime      time.Time
	Orders          Orders
}

func (p *Plan) SetEntrySize(activePair Pair, equity decimal.Decimal, o *Orders) (positionSize decimal.Decimal) {
	// TODO take into account transaction fees, see ByBit site for formula
	maxRisk := equity.Mul(p.Risk.Div(decimal.NewFromInt(100)))
	entryStopLossDistance := o[TypeHardStopLoss].Price.Sub(o[TypeEntry].Price).Abs()
	positionSize = maxRisk.Div(entryStopLossDistance).Round(activePair.PriceScale)
	o[TypeHardStopLoss].Size = positionSize
	o[TypeSoftStopLoss].Size = positionSize
	o[TypeEntry].Size = positionSize
	return positionSize
}

func (p *Plan) SetTakeProfitSizes(a Pair) {
	totalSize := p.Orders[TypeEntry].Size
	takeProfitsCount := 0
	for i := 1; i < MaxTakeProfits; i++ {
		if !p.Orders[2+i].Price.IsZero() {
			takeProfitsCount++
		}
	}

	takeProfitsCountDec := decimal.NewFromInt(int64(takeProfitsCount))
	// TODO: calculate Size remaining excluding orders that already been executed
	remainingSize := totalSize
	// TODO: implement other ways to divide takeprofit sizes
	takeProfitSize := totalSize.DivRound(takeProfitsCountDec, a.PriceScale)
	var i int
	// TODO: it feels like this could be much simpler, but i'm tired and it works so....
	// TODO: special case if there is only one TakeProfit
	if takeProfitsCount == 1 {
		p.Orders[3].Size = takeProfitSize
	} else {
		for i := 1; i <= takeProfitsCount-1; i++ {
			p.Orders[2+i].Size = takeProfitSize
			remainingSize = remainingSize.Sub(takeProfitSize)
		}
		// Set last take profit to the remainder
		p.Orders[2+i+2].Size = remainingSize.Round(a.PriceScale)
	}
}

func (p *Plan) SetRewardRiskRatio() (rrr float64) {
	maxRisk := (p.Orders[TypeEntry].Price.Mul(p.Orders[TypeEntry].Size)).Sub(p.Orders[TypeHardStopLoss].Price.Mul(p.Orders[TypeHardStopLoss].Size))
	maxProfit := decimal.Zero
	for _, order := range p.Orders {
		if order.OrderType == TypeTakeProfit {
			maxProfit = maxProfit.Add(order.Price.Sub(p.Orders[TypeEntry].Price).Mul(order.Size))
		}
	}
	rrr = maxProfit.Div(maxRisk).InexactFloat64()
	p.RewardRiskRatio = rrr
	return rrr
}

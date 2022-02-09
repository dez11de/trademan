package cryptodb

import (

	"github.com/bart613/decimal"
)

func (p *Plan) FinalizeOrders(available decimal.Decimal, activePair Pair, o []Order) {
	// TODO take into account transaction fees, see ByBit site for formula
	debugPrintOrders(o)
	maxRisk := available.Mul(p.Risk.Div(decimal.NewFromInt(100)))
	entryStopLossDistance := o[KindMarketStopLoss].Price.Sub(o[KindEntry].Price).Abs()
	positionSize := maxRisk.Div(entryStopLossDistance).RoundStep(activePair.Order.Step, false)
	o[KindMarketStopLoss].Size = positionSize
	o[KindLimitStopLoss].Size = positionSize
	o[KindEntry].Size = positionSize
	takeProfitsCount := int64(0)
	for i := 1; i < 5; i++ {
		if !o[2+i].Price.IsZero() {
			takeProfitsCount++
		}
	}

	// TODO: calculate Size remaining excluding orders that already been executed
	remainingSize := positionSize
	// TODO: implement other ways to divide takeprofit sizes instead of 100/TPCount
	takeProfitSize := positionSize.Div(decimal.NewFromInt(takeProfitsCount)).RoundStep(activePair.Order.Step, true)
	i := int64(1)
	for ; i <= takeProfitsCount-1; i++ {
		o[2+i].Size = takeProfitSize
		remainingSize = remainingSize.Sub(takeProfitSize)
	}
	o[2+i].Size = remainingSize.RoundStep(activePair.Order.Step, false)

	o[KindLimitStopLoss].Price = o[KindMarketStopLoss].Price.Add(entryStopLossDistance.Div(decimal.NewFromInt(100)).Mul(decimal.NewFromInt(5))).RoundStep(activePair.Price.Tick, false)
	debugPrintOrders(o)
}

// TODO: this should be done with "virtual" positionSizes since they are unknown at time of planning, and be set with real ordersizes at time of execution, although there should be no difference
func (p *Plan) SetRewardRiskRatio(o []Order) (rrr float64) {
	maxRisk := (o[KindEntry].Price.Mul(o[KindEntry].Size)).Sub(o[KindMarketStopLoss].Price.Mul(o[KindMarketStopLoss].Size))
	maxProfit := decimal.Zero
	for _, order := range o {
		if order.OrderType == OrderType(KindTakeProfit) {
			maxProfit = maxProfit.Add(order.Price.Sub(o[KindEntry].Price).Mul(order.Size))
		}
	}
	rrr = maxProfit.Div(maxRisk).InexactFloat64()
	p.RewardRiskRatio = rrr
	return rrr
}

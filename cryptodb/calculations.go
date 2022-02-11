package cryptodb

import (
	"github.com/bart613/decimal"
)

func (p *Plan) FinalizeOrders(available decimal.Decimal, activePair Pair, o []Order) {
    // positionSize = (entryPrice - stopLossPrice) * a vailable_balance * riskPerc * (1 - (pair.TakerFee * 2))
    maxFee := decimal.NewFromFloat(1.0).Sub(activePair.TakerFee.Mul(decimal.NewFromInt(2)))
	maxRisk := available.Mul(p.Risk.Div(decimal.NewFromInt(100))).Mul(maxFee)
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

	// TODO: calculate Size remaining excluding orders that already been filled
	remainingSize := positionSize
	// TODO: implement other ways to divide takeprofit sizes instead of 100/TPCount
	takeProfitSize := positionSize.Div(decimal.NewFromInt(takeProfitsCount)).RoundStep(activePair.Order.Step, false)
	i := int64(1)
	for ; i <= takeProfitsCount-1; i++ {
		o[2+i].Size = takeProfitSize
		remainingSize = remainingSize.Sub(takeProfitSize)
	}
	o[2+i].Size = remainingSize.RoundStep(activePair.Order.Step, false)

	o[KindLimitStopLoss].Price = o[KindMarketStopLoss].Price.Add(entryStopLossDistance.Div(decimal.NewFromInt(100)).Mul(decimal.NewFromInt(5))).RoundStep(activePair.Price.Tick, false)
}

// TODO: this should be done with "virtual" positionSizes since they are unknown at time of planning,
// and be set with real ordersizes at time of execution, although there should be no difference
func (p *Plan) SetRewardRiskRatio(o []Order) (rrr float64) {
    // TODO: see above for better calculation
	maxRisk := (o[KindEntry].Price.Mul(o[KindEntry].Size)).Sub(o[KindMarketStopLoss].Price.Mul(o[KindMarketStopLoss].Size))
	maxProfit := decimal.Zero
	for _, order := range o {
		if OrderType(order.OrderKind) == OrderType(KindTakeProfit) {
			maxProfit = maxProfit.Add(order.Price.Sub(o[KindEntry].Price).Mul(order.Size))
		}
	}
	rrr = maxProfit.Div(maxRisk).InexactFloat64()
	p.RewardRiskRatio = rrr
	return rrr
}

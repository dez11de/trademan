package cryptodb

import (
	"github.com/bart613/decimal"
)

func (p *Plan) FinalizeOrders(available decimal.Decimal,activePair Pair, o []Order) {
	// TODO take into account transaction fees, see ByBit site for formula
	maxRisk := available.Mul(p.Risk.Div(decimal.NewFromInt(100)))
	entryStopLossDistance := o[KindMarketStopLoss].Price.Sub(o[KindEntry].Price).Abs()
    positionSize := maxRisk.Div(entryStopLossDistance).Round(int32(activePair.PriceScale)) // TODO: should stepround to Pair.Order.Step
	o[KindMarketStopLoss].Size = positionSize
	o[KindLimitStopLoss].Size = positionSize
	o[KindEntry].Size = positionSize
	takeProfitsCount := 0
	for i := 1; i < 5; i++ {
		if !o[2+i].Price.IsZero() {
			takeProfitsCount++
		}
	}

	takeProfitsCountDec := decimal.NewFromInt(int64(takeProfitsCount))
	// TODO: calculate Size remaining excluding orders that already been executed
	remainingSize := positionSize
	// TODO: implement other ways to divide takeprofit sizes
	takeProfitSize := positionSize.DivRound(takeProfitsCountDec, int32(activePair.PriceScale)) // TODO: should stepround to Pair.Order.Step
	var i int
	// TODO: it feels like this could be much simpler, but i'm tired and it works so....
	// TODO: special case if there is only one TakeProfit
	if takeProfitsCount == 1 {
		o[3].Size = takeProfitSize 
	} else {
		for i := 1; i <= takeProfitsCount-1; i++ {
			o[2+i].Size = takeProfitSize
			remainingSize = remainingSize.Sub(takeProfitSize)
		}
		// Set last take profit to the remainder
		o[2+i+2].Size = remainingSize.Round(int32(activePair.PriceScale))// TODO: should stepround to Pair.Order.Step
	}
    o[KindLimitStopLoss].Price = o[KindMarketStopLoss].Price.Add(entryStopLossDistance.Div(decimal.NewFromInt(100)).Mul(decimal.NewFromInt(5)))
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

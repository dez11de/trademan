package cryptodb

import (
	"log"

	"github.com/bart613/decimal"
)

func triggerPrice(dir Direction, base, price decimal.Decimal) (trigger decimal.Decimal) {
	diff := price.Sub(base).Abs()
	triggerDiff := diff.Mul(decimal.NewFromFloat(0.95))
	if dir == Long {
		trigger = base.Add(triggerDiff)
	} else {
		trigger = base.Sub(triggerDiff)
	}

	return trigger
}

func (p *Plan) FinalizeOrders(available decimal.Decimal, activePair Pair, o []Order) {
	// positionSize = (entryPrice - stopLossPrice) * a vailable_balance * riskPerc * (1 - (pair.TakerFee * 2))
	maxFee := decimal.NewFromFloat(1.0).Sub(activePair.TakerFee.Mul(decimal.NewFromInt(2)))
	maxRisk := available.Mul(p.Risk.Div(decimal.NewFromInt(100))).Mul(maxFee)
	entryStopLossDistance := o[MarketStopLoss].Price.Sub(o[Entry].Price)
	positionSize := maxRisk.Div(entryStopLossDistance.Abs()).RoundStep(activePair.Order.Step, false)
	positionValue := positionSize.Mul(o[Entry].Price)
	o[Entry].Size = positionSize
	o[MarketStopLoss].Size = positionSize

	// o[LimitStopLoss].Size = positionSize
	// if p.Direction == Long {
	// 	o[LimitStopLoss].Price = triggerPrice(Short, o[Entry].Price, o[MarketStopLoss].Price).RoundStep(activePair.Price.Tick, false)
	// } else {
	// 	o[LimitStopLoss].Price = triggerPrice(Long, o[Entry].Price, o[MarketStopLoss].Price).RoundStep(activePair.Price.Tick, false)
	// }
	// o[LimitStopLoss].TriggerPrice = o[LimitStopLoss].Price

	if available.LessThan(positionSize.Mul(o[Entry].Price)) {
		// TODO: find out how to included costs and risk bybit calculates
		p.Leverage = positionValue.Add(maxRisk.Mul(decimal.NewFromInt(2))).Div(available).RoundStep(activePair.Leverage.Step, false)
	} else {
		p.Leverage = decimal.NewFromInt(1)
	}

	takeProfitsCount := int64(0)
	for i := 1; i <= MaxTakeProfits; i++ {
		if !o[2+i].Price.IsZero() {
			takeProfitsCount++
		}
	}

	// TODO: calculate Size remaining excluding orders that already been filled
	remainingSize := positionSize
	switch p.TakeProfitStrategy {
	case AutoLinear:
		takeProfitSize := positionSize.Div(decimal.NewFromInt(takeProfitsCount)).RoundStep(activePair.Order.Step, false)
		i := int64(1)
		for ; i <= takeProfitsCount-1; i++ {
			o[2+i].TriggerPrice = triggerPrice(p.Direction, o[Entry].Price, o[2+i].Price).RoundStep(activePair.Price.Tick, false)
			o[2+i].Size = takeProfitSize
			remainingSize = remainingSize.Sub(takeProfitSize)
		}
		o[2+i].TriggerPrice = triggerPrice(p.Direction, o[Entry].Price, o[2+i].Price).RoundStep(activePair.Price.Tick, false)
		o[2+i].Size = remainingSize.RoundStep(activePair.Order.Step, false)
	default:
		log.Printf("Take profit strategy %s is not (yet) implemented, sizes not set!", p.TakeProfitStrategy.String())
	}
}

// TODO: this should be done with "virtual" positionSizes since they are unknown at time of planning,
// and be set with real ordersizes at time of execution, although there should be no difference
func (p *Plan) SetRewardRiskRatio(o []Order) (rrr float64) {
	// TODO: see above for better calculation
	maxRisk := (o[Entry].Price.Mul(o[Entry].Size)).Sub(o[MarketStopLoss].Price.Mul(o[MarketStopLoss].Size))
	maxProfit := decimal.Zero
	for _, order := range o {
		if OrderType(order.OrderKind) == OrderType(TakeProfit) {
			maxProfit = maxProfit.Add(order.Price.Sub(o[Entry].Price).Mul(order.Size))
		}
	}
	rrr = maxProfit.Div(maxRisk).InexactFloat64()
	p.RewardRiskRatio = rrr
	return rrr
}

package indicators

import (
	"math"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func ATR(candles []market.Candle, period int) []market.IndicatorPoint {
	if period <= 0 || len(candles) < period {
		return nil
	}

	trueRanges := make([]float64, 0, len(candles))
	for i := range candles {
		if i == 0 {
			trueRanges = append(trueRanges, candles[i].High-candles[i].Low)
			continue
		}
		trueRanges = append(trueRanges, trueRange(candles[i], candles[i-1].Close))
	}

	var initialSum float64
	for i := 0; i < period; i++ {
		initialSum += trueRanges[i]
	}
	atr := initialSum / float64(period)

	points := []market.IndicatorPoint{{
		Timestamp: candles[period-1].Timestamp,
		Value:     atr,
	}}

	for i := period; i < len(candles); i++ {
		atr = ((atr * float64(period-1)) + trueRanges[i]) / float64(period)
		points = append(points, market.IndicatorPoint{
			Timestamp: candles[i].Timestamp,
			Value:     atr,
		})
	}

	return points
}

func trueRange(current market.Candle, previousClose float64) float64 {
	hl := current.High - current.Low
	hc := math.Abs(current.High - previousClose)
	lc := math.Abs(current.Low - previousClose)
	return math.Max(hl, math.Max(hc, lc))
}

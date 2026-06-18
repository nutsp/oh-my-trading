package indicators

import "github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"

func EMA(candles []market.Candle, period int) []market.IndicatorPoint {
	if period <= 0 || len(candles) < period {
		return nil
	}

	var sum float64
	for i := 0; i < period; i++ {
		sum += candles[i].Close
	}

	multiplier := 2.0 / float64(period+1)
	ema := sum / float64(period)

	points := []market.IndicatorPoint{{
		Timestamp: candles[period-1].Timestamp,
		Value:     ema,
	}}

	for i := period; i < len(candles); i++ {
		ema = ((candles[i].Close - ema) * multiplier) + ema
		points = append(points, market.IndicatorPoint{
			Timestamp: candles[i].Timestamp,
			Value:     ema,
		})
	}

	return points
}

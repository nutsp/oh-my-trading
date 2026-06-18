package indicators

import "github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"

func RSI(candles []market.Candle, period int) []market.IndicatorPoint {
	if period <= 0 || len(candles) <= period {
		return nil
	}

	var gains float64
	var losses float64

	for i := 1; i <= period; i++ {
		change := candles[i].Close - candles[i-1].Close
		if change >= 0 {
			gains += change
		} else {
			losses += -change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	points := []market.IndicatorPoint{{
		Timestamp: candles[period].Timestamp,
		Value:     rsiValue(avgGain, avgLoss),
	}}

	for i := period + 1; i < len(candles); i++ {
		change := candles[i].Close - candles[i-1].Close
		gain := 0.0
		loss := 0.0
		if change >= 0 {
			gain = change
		} else {
			loss = -change
		}

		avgGain = ((avgGain * float64(period-1)) + gain) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + loss) / float64(period)

		points = append(points, market.IndicatorPoint{
			Timestamp: candles[i].Timestamp,
			Value:     rsiValue(avgGain, avgLoss),
		})
	}

	return points
}

func rsiValue(avgGain, avgLoss float64) float64 {
	if avgLoss == 0 {
		return 100
	}
	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

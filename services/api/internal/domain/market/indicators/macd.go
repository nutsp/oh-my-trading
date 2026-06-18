package indicators

import "github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"

func MACD(candles []market.Candle, fastPeriod, slowPeriod, signalPeriod int) []market.MACDPoint {
	if fastPeriod <= 0 || slowPeriod <= 0 || signalPeriod <= 0 {
		return nil
	}

	fast := EMA(candles, fastPeriod)
	slow := EMA(candles, slowPeriod)
	if len(fast) == 0 || len(slow) == 0 {
		return nil
	}

	fastByTimestamp := make(map[int64]float64, len(fast))
	for _, point := range fast {
		fastByTimestamp[point.Timestamp.Unix()] = point.Value
	}

	macdCandidates := make([]market.IndicatorPoint, 0, len(slow))
	for _, slowPoint := range slow {
		fastValue, ok := fastByTimestamp[slowPoint.Timestamp.Unix()]
		if !ok {
			continue
		}
		macdCandidates = append(macdCandidates, market.IndicatorPoint{
			Timestamp: slowPoint.Timestamp,
			Value:     fastValue - slowPoint.Value,
		})
	}

	if len(macdCandidates) < signalPeriod {
		return nil
	}

	signalLine := emaFromPoints(macdCandidates, signalPeriod)
	if len(signalLine) == 0 {
		return nil
	}

	macdByTimestamp := make(map[int64]float64, len(macdCandidates))
	for _, point := range macdCandidates {
		macdByTimestamp[point.Timestamp.Unix()] = point.Value
	}

	points := make([]market.MACDPoint, 0, len(signalLine))
	for _, signalPoint := range signalLine {
		macdValue := macdByTimestamp[signalPoint.Timestamp.Unix()]
		points = append(points, market.MACDPoint{
			Timestamp: signalPoint.Timestamp,
			MACD:      macdValue,
			Signal:    signalPoint.Value,
			Histogram: macdValue - signalPoint.Value,
		})
	}

	return points
}

func emaFromPoints(points []market.IndicatorPoint, period int) []market.IndicatorPoint {
	if period <= 0 || len(points) < period {
		return nil
	}

	var sum float64
	for i := 0; i < period; i++ {
		sum += points[i].Value
	}

	multiplier := 2.0 / float64(period+1)
	ema := sum / float64(period)

	result := []market.IndicatorPoint{{
		Timestamp: points[period-1].Timestamp,
		Value:     ema,
	}}

	for i := period; i < len(points); i++ {
		ema = ((points[i].Value - ema) * multiplier) + ema
		result = append(result, market.IndicatorPoint{
			Timestamp: points[i].Timestamp,
			Value:     ema,
		})
	}

	return result
}

package marketdataadapter

import (
	"context"
	"time"

	app "github.com/sutad-p/oh-my-trading/services/api/internal/application/marketdata"
	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

type SyntheticProvider struct{}

func NewSyntheticProvider() *SyntheticProvider {
	return &SyntheticProvider{}
}

func (p *SyntheticProvider) FetchCandles(ctx context.Context, request app.FetchCandlesRequest) ([]market.Candle, error) {
	step := timeframeDuration(request.Timeframe)
	if step <= 0 {
		step = time.Hour
	}

	var candles []market.Candle
	for ts := request.From; !ts.After(request.To); ts = ts.Add(step) {
		offset := float64(len(candles))
		open := 100 + offset
		candles = append(candles, market.Candle{
			SymbolID:  request.SymbolID,
			Timeframe: request.Timeframe,
			Timestamp: ts,
			Open:      open,
			High:      open + 2,
			Low:       open - 1,
			Close:     open + 1,
			Volume:    1000 + offset,
		})
	}
	return candles, nil
}

func timeframeDuration(timeframe string) time.Duration {
	switch timeframe {
	case "1m":
		return time.Minute
	case "5m":
		return 5 * time.Minute
	case "15m":
		return 15 * time.Minute
	case "1h":
		return time.Hour
	case "4h":
		return 4 * time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return 0
	}
}

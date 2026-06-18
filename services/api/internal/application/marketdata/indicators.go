package marketdata

import (
	"context"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market/indicators"
)

type IndicatorService struct {
	repo  market.CandleRepository
	nowFn func() time.Time
}

func NewIndicatorService(repo market.CandleRepository, nowFn func() time.Time) *IndicatorService {
	return &IndicatorService{
		repo:  repo,
		nowFn: nowFn,
	}
}

func (s *IndicatorService) ListIndicators(ctx context.Context, query market.IndicatorQuery) (market.IndicatorSeries, error) {
	to := s.nowFn().UTC()
	from := to.Add(-lookbackByTimeframe(query.Timeframe))

	candles, err := s.repo.ListCandles(ctx, market.CandleQuery{
		SymbolCode: query.SymbolCode,
		Timeframe:  query.Timeframe,
		From:       from,
		To:         to,
	})
	if err != nil {
		return market.IndicatorSeries{}, err
	}

	return market.IndicatorSeries{
		SymbolCode: query.SymbolCode,
		Timeframe:  query.Timeframe,
		EMA20:      indicators.EMA(candles, 20),
		EMA50:      indicators.EMA(candles, 50),
		RSI14:      indicators.RSI(candles, 14),
		MACD:       indicators.MACD(candles, 12, 26, 9),
		ATR14:      indicators.ATR(candles, 14),
	}, nil
}

func lookbackByTimeframe(timeframe string) time.Duration {
	switch timeframe {
	case "15m":
		return 45 * 24 * time.Hour
	case "1h":
		return 180 * 24 * time.Hour
	case "4h":
		return 540 * 24 * time.Hour
	case "1d":
		return 5 * 365 * 24 * time.Hour
	default:
		return 365 * 24 * time.Hour
	}
}

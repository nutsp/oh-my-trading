package marketdata

import (
	"context"
	"testing"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestIndicatorServiceComputesSeriesFromCandles(t *testing.T) {
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	repo := &memoryIndicatorCandleRepository{
		candles: indicatorFixtureCandles(now.Add(-40 * 24 * time.Hour)),
	}

	service := NewIndicatorService(repo, func() time.Time { return now })
	series, err := service.ListIndicators(context.Background(), market.IndicatorQuery{
		SymbolCode: "XAUUSD",
		Timeframe:  "1h",
	})
	if err != nil {
		t.Fatalf("ListIndicators returned error: %v", err)
	}

	if series.SymbolCode != "XAUUSD" {
		t.Fatalf("SymbolCode = %q, want XAUUSD", series.SymbolCode)
	}
	if series.Timeframe != "1h" {
		t.Fatalf("Timeframe = %q, want 1h", series.Timeframe)
	}
	if len(series.EMA20) == 0 || len(series.EMA50) == 0 || len(series.RSI14) == 0 || len(series.ATR14) == 0 || len(series.MACD) == 0 {
		t.Fatalf("expected non-empty indicator series, got ema20=%d ema50=%d rsi14=%d atr14=%d macd=%d",
			len(series.EMA20), len(series.EMA50), len(series.RSI14), len(series.ATR14), len(series.MACD))
	}

	if repo.query.SymbolCode != "XAUUSD" || repo.query.Timeframe != "1h" {
		t.Fatalf("query = %+v", repo.query)
	}
	if !repo.query.To.Equal(now) {
		t.Fatalf("query.To = %s, want %s", repo.query.To, now)
	}
}

type memoryIndicatorCandleRepository struct {
	query   market.CandleQuery
	candles []market.Candle
}

func (r *memoryIndicatorCandleRepository) UpsertCandles(ctx context.Context, candles []market.Candle) error {
	r.candles = append([]market.Candle(nil), candles...)
	return nil
}

func (r *memoryIndicatorCandleRepository) ListCandles(ctx context.Context, query market.CandleQuery) ([]market.Candle, error) {
	r.query = query
	return append([]market.Candle(nil), r.candles...), nil
}

func indicatorFixtureCandles(start time.Time) []market.Candle {
	candles := make([]market.Candle, 0, 120)
	close := 2300.0

	for i := 0; i < 120; i++ {
		drift := 0.4
		if i%12 == 0 || i%12 == 1 {
			drift = -0.7
		}
		close += drift + float64((i%4)-2)*0.25
		open := close - 0.35

		candles = append(candles, market.Candle{
			SymbolID:  "fixture",
			Timeframe: "1h",
			Timestamp: start.Add(time.Duration(i) * time.Hour),
			Open:      open,
			High:      close + 0.9,
			Low:       open - 0.8,
			Close:     close,
			Volume:    1000 + float64(i),
		})
	}

	return candles
}

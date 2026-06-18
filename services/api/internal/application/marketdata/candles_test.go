package marketdata

import (
	"context"
	"testing"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestCandleServiceStoresAndListsCandles(t *testing.T) {
	repo := &memoryCandleRepository{}
	service := NewCandleService(repo)

	ts := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	err := service.UpsertCandles(context.Background(), []market.Candle{{
		SymbolID:  "018f4f8a-0000-7000-9000-000000000001",
		Timeframe: "1h",
		Timestamp: ts,
		Open:      2320.1,
		High:      2328.3,
		Low:       2318.6,
		Close:     2325.5,
		Volume:    12345,
	}})
	if err != nil {
		t.Fatalf("UpsertCandles returned error: %v", err)
	}

	candles, err := service.ListCandles(context.Background(), market.CandleQuery{
		SymbolCode: "XAUUSD",
		Timeframe:  "1h",
		From:       ts.Add(-time.Hour),
		To:         ts.Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("ListCandles returned error: %v", err)
	}
	if len(candles) != 1 {
		t.Fatalf("len(candles) = %d, want 1", len(candles))
	}
	if candles[0].Close != 2325.5 {
		t.Fatalf("Close = %f, want 2325.5", candles[0].Close)
	}
}

type memoryCandleRepository struct {
	candles []market.Candle
}

func (r *memoryCandleRepository) UpsertCandles(ctx context.Context, candles []market.Candle) error {
	r.candles = append([]market.Candle(nil), candles...)
	return nil
}

func (r *memoryCandleRepository) ListCandles(ctx context.Context, query market.CandleQuery) ([]market.Candle, error) {
	return append([]market.Candle(nil), r.candles...), nil
}

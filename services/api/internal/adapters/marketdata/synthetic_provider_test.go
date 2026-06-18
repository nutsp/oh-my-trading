package marketdataadapter

import (
	"context"
	"testing"
	"time"

	app "github.com/sutad-p/oh-my-trading/services/api/internal/application/marketdata"
)

func TestSyntheticProviderReturnsCandlesForRequestedRange(t *testing.T) {
	provider := NewSyntheticProvider()
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(2 * time.Hour)

	candles, err := provider.FetchCandles(context.Background(), app.FetchCandlesRequest{
		SymbolID:   "018f4f8a-0000-7000-9000-000000000401",
		SymbolCode: "XAUUSD",
		Timeframe:  "1h",
		From:       from,
		To:         to,
	})
	if err != nil {
		t.Fatalf("FetchCandles returned error: %v", err)
	}
	if len(candles) != 3 {
		t.Fatalf("len(candles) = %d, want 3", len(candles))
	}
	if candles[0].Timestamp != from {
		t.Fatalf("first timestamp = %s, want %s", candles[0].Timestamp, from)
	}
	if candles[2].Timestamp != to {
		t.Fatalf("last timestamp = %s, want %s", candles[2].Timestamp, to)
	}
}

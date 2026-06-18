package marketdata

import (
	"context"
	"testing"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestSyncServiceFetchesAndStoresCandlesForEachTimeframe(t *testing.T) {
	repo := &recordingCandleRepository{}
	provider := &fakeMarketDataProvider{}
	service := NewSyncService(provider, repo)

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := from.Add(2 * time.Hour)
	err := service.SyncCandles(context.Background(), SyncRequest{
		SymbolID:   "018f4f8a-0000-7000-9000-000000000401",
		SymbolCode: "XAUUSD",
		Timeframes: []string{"1h", "4h"},
		From:       from,
		To:         to,
	})
	if err != nil {
		t.Fatalf("SyncCandles returned error: %v", err)
	}

	if len(provider.requests) != 2 {
		t.Fatalf("provider requests = %d, want 2", len(provider.requests))
	}
	if len(repo.upserted) != 2 {
		t.Fatalf("upserted = %d, want 2", len(repo.upserted))
	}
	if repo.upserted[0].SymbolID != "018f4f8a-0000-7000-9000-000000000401" {
		t.Fatalf("SymbolID = %q", repo.upserted[0].SymbolID)
	}
}

type fakeMarketDataProvider struct {
	requests []FetchCandlesRequest
}

func (p *fakeMarketDataProvider) FetchCandles(ctx context.Context, request FetchCandlesRequest) ([]market.Candle, error) {
	p.requests = append(p.requests, request)
	return []market.Candle{{
		SymbolID:  request.SymbolID,
		Timeframe: request.Timeframe,
		Timestamp: request.From,
		Open:      1,
		High:      2,
		Low:       0.5,
		Close:     1.5,
		Volume:    10,
	}}, nil
}

type recordingCandleRepository struct {
	upserted []market.Candle
}

func (r *recordingCandleRepository) UpsertCandles(ctx context.Context, candles []market.Candle) error {
	r.upserted = append(r.upserted, candles...)
	return nil
}

func (r *recordingCandleRepository) ListCandles(ctx context.Context, query market.CandleQuery) ([]market.Candle, error) {
	return nil, nil
}

package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestListCandlesRouteReturnsCandles(t *testing.T) {
	ts := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	service := &fakeCandleService{
		candles: []market.Candle{{
			SymbolID:  "018f4f8a-0000-7000-9000-000000000001",
			Timeframe: "1h",
			Timestamp: ts,
			Open:      2320.1,
			High:      2328.3,
			Low:       2318.6,
			Close:     2325.5,
			Volume:    12345,
		}},
	}
	router := NewRouter(WithCandleService(service))

	req := httptest.NewRequest(http.MethodGet, "/api/candles?symbol=XAUUSD&timeframe=1h&from=2026-01-01T00:00:00Z&to=2026-01-02T00:00:00Z", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if service.query.SymbolCode != "XAUUSD" {
		t.Fatalf("SymbolCode = %q, want XAUUSD", service.query.SymbolCode)
	}
	if service.query.Timeframe != "1h" {
		t.Fatalf("Timeframe = %q, want 1h", service.query.Timeframe)
	}

	var response []candleResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("len(response) = %d, want 1", len(response))
	}
	if response[0].Timestamp != "2026-01-01T12:00:00Z" {
		t.Fatalf("Timestamp = %q", response[0].Timestamp)
	}
	if response[0].Close != 2325.5 {
		t.Fatalf("Close = %f, want 2325.5", response[0].Close)
	}
}

func TestListCandlesRouteRequiresQueryParams(t *testing.T) {
	router := NewRouter(WithCandleService(&fakeCandleService{}))

	req := httptest.NewRequest(http.MethodGet, "/api/candles?symbol=XAUUSD", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

type fakeCandleService struct {
	query   market.CandleQuery
	candles []market.Candle
}

func (s *fakeCandleService) ListCandles(ctx context.Context, query market.CandleQuery) ([]market.Candle, error) {
	s.query = query
	return s.candles, nil
}

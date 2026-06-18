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

func TestListIndicatorsRouteReturnsIndicatorSeries(t *testing.T) {
	ts := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	service := &fakeIndicatorService{
		series: market.IndicatorSeries{
			SymbolCode: "XAUUSD",
			Timeframe:  "1h",
			EMA20:      []market.IndicatorPoint{{Timestamp: ts, Value: 2321.2}},
			EMA50:      []market.IndicatorPoint{{Timestamp: ts, Value: 2315.8}},
			RSI14:      []market.IndicatorPoint{{Timestamp: ts, Value: 58.4}},
			ATR14:      []market.IndicatorPoint{{Timestamp: ts, Value: 11.3}},
			MACD:       []market.MACDPoint{{Timestamp: ts, MACD: 5.1, Signal: 4.9, Histogram: 0.2}},
		},
	}
	router := NewRouter(WithIndicatorService(service))

	req := httptest.NewRequest(http.MethodGet, "/api/indicators?symbol=XAUUSD&timeframe=1h", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if service.query.SymbolCode != "XAUUSD" || service.query.Timeframe != "1h" {
		t.Fatalf("query = %+v", service.query)
	}

	var response indicatorResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Symbol != "XAUUSD" || response.Timeframe != "1h" {
		t.Fatalf("response symbol/timeframe = %q/%q", response.Symbol, response.Timeframe)
	}
	if len(response.Series.EMA20) != 1 || response.Series.EMA20[0].Value != 2321.2 {
		t.Fatalf("unexpected EMA20 payload: %+v", response.Series.EMA20)
	}
	if len(response.Series.MACD) != 1 || response.Series.MACD[0].Histogram != 0.2 {
		t.Fatalf("unexpected MACD payload: %+v", response.Series.MACD)
	}
}

func TestListIndicatorsRouteRequiresQueryParams(t *testing.T) {
	router := NewRouter(WithIndicatorService(&fakeIndicatorService{}))

	req := httptest.NewRequest(http.MethodGet, "/api/indicators?symbol=XAUUSD", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

type fakeIndicatorService struct {
	query  market.IndicatorQuery
	series market.IndicatorSeries
}

func (s *fakeIndicatorService) ListIndicators(ctx context.Context, query market.IndicatorQuery) (market.IndicatorSeries, error) {
	s.query = query
	return s.series, nil
}

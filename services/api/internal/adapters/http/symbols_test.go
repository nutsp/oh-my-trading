package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestListSymbolsRouteReturnsSymbols(t *testing.T) {
	router := NewRouter(WithSymbolService(&fakeSymbolService{
		symbols: []market.Symbol{{
			ID:         "018f4f8a-0000-7000-9000-000000000001",
			Code:       "XAUUSD",
			Market:     "forex",
			BaseAsset:  "XAU",
			QuoteAsset: "USD",
			Enabled:    true,
		}},
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/symbols", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response []symbolResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(response) != 1 {
		t.Fatalf("len(response) = %d, want 1", len(response))
	}
	if response[0].Code != "XAUUSD" {
		t.Fatalf("Code = %q, want XAUUSD", response[0].Code)
	}
}

func TestCreateSymbolRouteCreatesSymbol(t *testing.T) {
	service := &fakeSymbolService{}
	router := NewRouter(WithSymbolService(service))

	req := httptest.NewRequest(http.MethodPost, "/api/symbols", strings.NewReader(`{
		"code": "BTCUSD",
		"market": "crypto",
		"baseAsset": "BTC",
		"quoteAsset": "USD"
	}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if service.created.Code != "BTCUSD" {
		t.Fatalf("created code = %q, want BTCUSD", service.created.Code)
	}
}

type fakeSymbolService struct {
	symbols []market.Symbol
	created market.CreateSymbolInput
}

func (s *fakeSymbolService) ListSymbols(ctx context.Context) ([]market.Symbol, error) {
	return s.symbols, nil
}

func (s *fakeSymbolService) CreateSymbol(ctx context.Context, input market.CreateSymbolInput) (market.Symbol, error) {
	s.created = input
	return market.Symbol{
		ID:         "018f4f8a-0000-7000-9000-000000000002",
		Code:       input.Code,
		Market:     input.Market,
		BaseAsset:  input.BaseAsset,
		QuoteAsset: input.QuoteAsset,
		Enabled:    true,
	}, nil
}

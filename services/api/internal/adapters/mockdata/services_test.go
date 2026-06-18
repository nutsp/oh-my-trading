package mockdata

import (
	"context"
	"testing"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestSymbolServiceListsDefaultSymbols(t *testing.T) {
	service := NewSymbolService(func() string { return "generated-id" })

	symbols, err := service.ListSymbols(context.Background())
	if err != nil {
		t.Fatalf("ListSymbols returned error: %v", err)
	}
	if len(symbols) < 3 {
		t.Fatalf("len(symbols) = %d, want at least 3", len(symbols))
	}
}

func TestSymbolServiceCreateSymbolNormalizesInput(t *testing.T) {
	service := NewSymbolService(func() string { return "generated-id" })

	symbol, err := service.CreateSymbol(context.Background(), market.CreateSymbolInput{
		Code:       "btcusd",
		Market:     "crypto",
		BaseAsset:  "btc",
		QuoteAsset: "usd",
	})
	if err != nil {
		t.Fatalf("CreateSymbol returned error: %v", err)
	}
	if symbol.ID != "generated-id" {
		t.Fatalf("ID = %q, want generated-id", symbol.ID)
	}
	if symbol.Code != "BTCUSD" || symbol.BaseAsset != "BTC" || symbol.QuoteAsset != "USD" {
		t.Fatalf("symbol normalization failed: %+v", symbol)
	}
}

func TestCandleServiceReturnsMockCandles(t *testing.T) {
	service := NewCandleService()

	candles, err := service.ListCandles(context.Background(), market.CandleQuery{
		SymbolCode: "XAUUSD",
		Timeframe:  "1h",
		From:       time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		To:         time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("ListCandles returned error: %v", err)
	}
	if len(candles) == 0 {
		t.Fatal("expected candles, got empty result")
	}
}

func TestIndicatorServiceReturnsCalculatedIndicators(t *testing.T) {
	service := NewIndicatorService()

	series, err := service.ListIndicators(context.Background(), market.IndicatorQuery{
		SymbolCode: "XAUUSD",
		Timeframe:  "1h",
	})
	if err != nil {
		t.Fatalf("ListIndicators returned error: %v", err)
	}

	if len(series.EMA20) == 0 || len(series.EMA50) == 0 || len(series.RSI14) == 0 || len(series.ATR14) == 0 || len(series.MACD) == 0 {
		t.Fatalf("indicator series should not be empty: %+v", series)
	}
}

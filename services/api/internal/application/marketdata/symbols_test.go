package marketdata

import (
	"context"
	"testing"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestCreateSymbolAssignsIDAndEnablesSymbol(t *testing.T) {
	repo := &memorySymbolRepository{}
	service := NewSymbolService(repo, func() string {
		return "018f4f8a-0000-7000-9000-000000000001"
	})

	symbol, err := service.CreateSymbol(context.Background(), market.CreateSymbolInput{
		Code:       "XAUUSD",
		Market:     "forex",
		BaseAsset:  "XAU",
		QuoteAsset: "USD",
	})
	if err != nil {
		t.Fatalf("CreateSymbol returned error: %v", err)
	}

	if symbol.ID != "018f4f8a-0000-7000-9000-000000000001" {
		t.Fatalf("ID = %q", symbol.ID)
	}
	if !symbol.Enabled {
		t.Fatal("expected symbol to be enabled")
	}

	symbols, err := service.ListSymbols(context.Background())
	if err != nil {
		t.Fatalf("ListSymbols returned error: %v", err)
	}
	if len(symbols) != 1 {
		t.Fatalf("len(symbols) = %d, want 1", len(symbols))
	}
}

type memorySymbolRepository struct {
	symbols []market.Symbol
}

func (r *memorySymbolRepository) ListSymbols(ctx context.Context) ([]market.Symbol, error) {
	return append([]market.Symbol(nil), r.symbols...), nil
}

func (r *memorySymbolRepository) CreateSymbol(ctx context.Context, symbol market.Symbol) (market.Symbol, error) {
	r.symbols = append(r.symbols, symbol)
	return symbol, nil
}

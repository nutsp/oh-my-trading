package market

import "context"

type Symbol struct {
	ID         string
	Code       string
	Market     string
	BaseAsset  string
	QuoteAsset string
	Enabled    bool
}

type CreateSymbolInput struct {
	Code       string
	Market     string
	BaseAsset  string
	QuoteAsset string
}

type SymbolRepository interface {
	ListSymbols(ctx context.Context) ([]Symbol, error)
	CreateSymbol(ctx context.Context, symbol Symbol) (Symbol, error)
}

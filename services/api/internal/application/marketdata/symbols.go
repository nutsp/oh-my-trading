package marketdata

import (
	"context"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

type IDGenerator func() string

type SymbolService struct {
	repo       market.SymbolRepository
	generateID IDGenerator
}

func NewSymbolService(repo market.SymbolRepository, generateID IDGenerator) *SymbolService {
	return &SymbolService{
		repo:       repo,
		generateID: generateID,
	}
}

func (s *SymbolService) ListSymbols(ctx context.Context) ([]market.Symbol, error) {
	return s.repo.ListSymbols(ctx)
}

func (s *SymbolService) CreateSymbol(ctx context.Context, input market.CreateSymbolInput) (market.Symbol, error) {
	symbol := market.Symbol{
		ID:         s.generateID(),
		Code:       input.Code,
		Market:     input.Market,
		BaseAsset:  input.BaseAsset,
		QuoteAsset: input.QuoteAsset,
		Enabled:    true,
	}
	return s.repo.CreateSymbol(ctx, symbol)
}

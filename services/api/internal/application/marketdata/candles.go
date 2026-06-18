package marketdata

import (
	"context"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

type CandleService struct {
	repo market.CandleRepository
}

func NewCandleService(repo market.CandleRepository) *CandleService {
	return &CandleService{repo: repo}
}

func (s *CandleService) UpsertCandles(ctx context.Context, candles []market.Candle) error {
	return s.repo.UpsertCandles(ctx, candles)
}

func (s *CandleService) ListCandles(ctx context.Context, query market.CandleQuery) ([]market.Candle, error) {
	return s.repo.ListCandles(ctx, query)
}

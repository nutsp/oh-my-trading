package marketdata

import (
	"context"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

type SyncRequest struct {
	SymbolID   string    `json:"symbolId"`
	SymbolCode string    `json:"symbolCode"`
	Timeframes []string  `json:"timeframes"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
}

type FetchCandlesRequest struct {
	SymbolID   string
	SymbolCode string
	Timeframe  string
	From       time.Time
	To         time.Time
}

type MarketDataProvider interface {
	FetchCandles(ctx context.Context, request FetchCandlesRequest) ([]market.Candle, error)
}

type SyncPublisher interface {
	PublishSyncRequest(ctx context.Context, request SyncRequest) error
}

type SyncService struct {
	provider MarketDataProvider
	candles  market.CandleRepository
}

func NewSyncService(provider MarketDataProvider, candles market.CandleRepository) *SyncService {
	return &SyncService{
		provider: provider,
		candles:  candles,
	}
}

func (s *SyncService) SyncCandles(ctx context.Context, request SyncRequest) error {
	for _, timeframe := range request.Timeframes {
		candles, err := s.provider.FetchCandles(ctx, FetchCandlesRequest{
			SymbolID:   request.SymbolID,
			SymbolCode: request.SymbolCode,
			Timeframe:  timeframe,
			From:       request.From,
			To:         request.To,
		})
		if err != nil {
			return err
		}
		if err := s.candles.UpsertCandles(ctx, candles); err != nil {
			return err
		}
	}
	return nil
}

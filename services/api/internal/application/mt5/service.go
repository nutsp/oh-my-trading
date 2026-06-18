package mt5

import (
	"context"
	"errors"
	"fmt"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
	domainmt5 "github.com/sutad-p/oh-my-trading/services/api/internal/domain/mt5"
)

var (
	ErrUnsupportedSymbol   = errors.New("unsupported symbol")
	ErrSymbolNotConfigured = errors.New("symbol is not configured")
)

type Service struct {
	mt5Repo    domainmt5.Repository
	symbolRepo market.SymbolRepository
	candleRepo market.CandleRepository
}

func NewService(mt5Repo domainmt5.Repository, symbolRepo market.SymbolRepository, candleRepo market.CandleRepository) *Service {
	return &Service{
		mt5Repo:    mt5Repo,
		symbolRepo: symbolRepo,
		candleRepo: candleRepo,
	}
}

func (s *Service) IngestHeartbeat(ctx context.Context, heartbeat domainmt5.Heartbeat) error {
	if err := s.mt5Repo.SaveHeartbeat(ctx, heartbeat); err != nil {
		return fmt.Errorf("ingest mt5 heartbeat: %w", err)
	}
	return nil
}

func (s *Service) IngestTicks(ctx context.Context, ticks []domainmt5.Tick) error {
	for _, tick := range ticks {
		if err := validateXAUUSD(tick.Symbol); err != nil {
			return err
		}
	}
	if err := s.mt5Repo.SaveTicks(ctx, ticks); err != nil {
		return fmt.Errorf("ingest mt5 ticks: %w", err)
	}
	return nil
}

func (s *Service) IngestCandles(ctx context.Context, batch domainmt5.CandleBatch) error {
	if err := validateXAUUSD(batch.Symbol); err != nil {
		return err
	}
	if len(batch.Candles) == 0 {
		return nil
	}

	symbol, err := s.findEnabledXAUUSD(ctx)
	if err != nil {
		return err
	}

	candles := make([]market.Candle, 0, len(batch.Candles))
	for _, candle := range batch.Candles {
		candles = append(candles, market.Candle{
			SymbolID:  symbol.ID,
			Timeframe: batch.Timeframe,
			Timestamp: candle.Timestamp,
			Open:      candle.Open,
			High:      candle.High,
			Low:       candle.Low,
			Close:     candle.Close,
			Volume:    candle.Volume,
		})
	}

	if err := s.candleRepo.UpsertCandles(ctx, candles); err != nil {
		return fmt.Errorf("ingest mt5 candles: %w", err)
	}
	return nil
}

func (s *Service) IngestAccountSnapshot(ctx context.Context, snapshot domainmt5.AccountSnapshot) error {
	if err := s.mt5Repo.SaveAccountSnapshot(ctx, snapshot); err != nil {
		return fmt.Errorf("ingest mt5 account snapshot: %w", err)
	}
	return nil
}

func (s *Service) IngestPositionSnapshots(ctx context.Context, positions []domainmt5.PositionSnapshot) error {
	for _, position := range positions {
		if err := validateXAUUSD(position.Symbol); err != nil {
			return err
		}
	}
	if err := s.mt5Repo.SavePositionSnapshots(ctx, positions); err != nil {
		return fmt.Errorf("ingest mt5 position snapshots: %w", err)
	}
	return nil
}

func (s *Service) LatestHeartbeat(ctx context.Context, bridgeID string) (domainmt5.Heartbeat, error) {
	return s.mt5Repo.LatestHeartbeat(ctx, bridgeID)
}

func (s *Service) LatestTick(ctx context.Context, symbol string) (domainmt5.Tick, error) {
	if err := validateXAUUSD(symbol); err != nil {
		return domainmt5.Tick{}, err
	}
	return s.mt5Repo.LatestTick(ctx, symbol)
}

func (s *Service) LatestAccountSnapshot(ctx context.Context, accountLogin string) (domainmt5.AccountSnapshot, error) {
	return s.mt5Repo.LatestAccountSnapshot(ctx, accountLogin)
}

func (s *Service) LatestPositionSnapshots(ctx context.Context, accountLogin string) ([]domainmt5.PositionSnapshot, error) {
	return s.mt5Repo.LatestPositionSnapshots(ctx, accountLogin)
}

func (s *Service) findEnabledXAUUSD(ctx context.Context) (market.Symbol, error) {
	symbols, err := s.symbolRepo.ListSymbols(ctx)
	if err != nil {
		return market.Symbol{}, fmt.Errorf("list symbols: %w", err)
	}
	for _, symbol := range symbols {
		if symbol.Code == domainmt5.XAUUSDSymbol && symbol.Enabled {
			return symbol, nil
		}
	}
	return market.Symbol{}, ErrSymbolNotConfigured
}

func validateXAUUSD(symbol string) error {
	if symbol != domainmt5.XAUUSDSymbol {
		return fmt.Errorf("%w: %s", ErrUnsupportedSymbol, symbol)
	}
	return nil
}

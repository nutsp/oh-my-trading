package mt5

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
	domainmt5 "github.com/sutad-p/oh-my-trading/services/api/internal/domain/mt5"
)

func TestServiceIngestsXAUUSDReadOnlyData(t *testing.T) {
	mt5Repo := &memoryMT5Repository{}
	symbolRepo := &memorySymbolRepository{
		symbols: []market.Symbol{{
			ID:         "018f4f8a-0000-7000-9000-000000000301",
			Code:       "XAUUSD",
			Market:     "forex",
			BaseAsset:  "XAU",
			QuoteAsset: "USD",
			Enabled:    true,
		}},
	}
	candleRepo := &memoryCandleRepository{}
	service := NewService(mt5Repo, symbolRepo, candleRepo)

	now := time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC)
	if err := service.IngestHeartbeat(context.Background(), domainmt5.Heartbeat{
		BridgeID:     "local-mt5",
		Terminal:     "MetaTrader 5",
		AccountLogin: "12345678",
		Server:       "Broker-Demo",
		Status:       "healthy",
		SentAt:       now,
	}); err != nil {
		t.Fatalf("IngestHeartbeat returned error: %v", err)
	}

	if err := service.IngestTicks(context.Background(), []domainmt5.Tick{{
		Symbol: "XAUUSD",
		Bid:    2325.42,
		Ask:    2325.62,
		Last:   2325.52,
		Volume: 12,
		Time:   now,
	}}); err != nil {
		t.Fatalf("IngestTicks returned error: %v", err)
	}

	if err := service.IngestCandles(context.Background(), domainmt5.CandleBatch{
		Symbol:    "XAUUSD",
		Timeframe: "1m",
		Source:    "mt5-python-bridge",
		Candles: []domainmt5.Candle{{
			Timestamp: now,
			Open:      2320.1,
			High:      2328.3,
			Low:       2318.6,
			Close:     2325.5,
			Volume:    12345,
		}},
	}); err != nil {
		t.Fatalf("IngestCandles returned error: %v", err)
	}

	if err := service.IngestAccountSnapshot(context.Background(), domainmt5.AccountSnapshot{
		AccountLogin: "12345678",
		Currency:     "USD",
		Balance:      10000,
		Equity:       10080,
		Margin:       400,
		FreeMargin:   9680,
		MarginLevel:  2520,
		Time:         now,
	}); err != nil {
		t.Fatalf("IngestAccountSnapshot returned error: %v", err)
	}

	if err := service.IngestPositionSnapshots(context.Background(), []domainmt5.PositionSnapshot{{
		AccountLogin: "12345678",
		Ticket:       "987654321",
		Symbol:       "XAUUSD",
		Side:         "buy",
		Volume:       0.1,
		OpenPrice:    2320.1,
		StopLoss:     2310,
		TakeProfit:   2340,
		Profit:       55,
		OpenedAt:     now.Add(-time.Hour),
		SnapshotTime: now,
	}}); err != nil {
		t.Fatalf("IngestPositionSnapshots returned error: %v", err)
	}

	if mt5Repo.heartbeat.BridgeID != "local-mt5" {
		t.Fatalf("heartbeat BridgeID = %q, want local-mt5", mt5Repo.heartbeat.BridgeID)
	}
	if len(mt5Repo.ticks) != 1 {
		t.Fatalf("len(ticks) = %d, want 1", len(mt5Repo.ticks))
	}
	if len(candleRepo.candles) != 1 {
		t.Fatalf("len(candles) = %d, want 1", len(candleRepo.candles))
	}
	if candleRepo.candles[0].SymbolID != "018f4f8a-0000-7000-9000-000000000301" {
		t.Fatalf("candle SymbolID = %q, want configured XAUUSD symbol id", candleRepo.candles[0].SymbolID)
	}
	if mt5Repo.account.Equity != 10080 {
		t.Fatalf("account Equity = %f, want 10080", mt5Repo.account.Equity)
	}
	if len(mt5Repo.positions) != 1 {
		t.Fatalf("len(positions) = %d, want 1", len(mt5Repo.positions))
	}
}

func TestServiceRejectsNonXAUUSDPayloads(t *testing.T) {
	service := NewService(&memoryMT5Repository{}, &memorySymbolRepository{}, &memoryCandleRepository{})

	err := service.IngestTicks(context.Background(), []domainmt5.Tick{{Symbol: "EURUSD"}})
	if !errors.Is(err, ErrUnsupportedSymbol) {
		t.Fatalf("IngestTicks error = %v, want ErrUnsupportedSymbol", err)
	}

	err = service.IngestCandles(context.Background(), domainmt5.CandleBatch{Symbol: "BTCUSD"})
	if !errors.Is(err, ErrUnsupportedSymbol) {
		t.Fatalf("IngestCandles error = %v, want ErrUnsupportedSymbol", err)
	}

	err = service.IngestPositionSnapshots(context.Background(), []domainmt5.PositionSnapshot{{Symbol: "ETHUSD"}})
	if !errors.Is(err, ErrUnsupportedSymbol) {
		t.Fatalf("IngestPositionSnapshots error = %v, want ErrUnsupportedSymbol", err)
	}
}

func TestServiceRequiresConfiguredXAUUSDSymbolForCandles(t *testing.T) {
	service := NewService(&memoryMT5Repository{}, &memorySymbolRepository{}, &memoryCandleRepository{})

	err := service.IngestCandles(context.Background(), domainmt5.CandleBatch{
		Symbol:    "XAUUSD",
		Timeframe: "1m",
		Candles: []domainmt5.Candle{{
			Timestamp: time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC),
			Open:      2320.1,
			High:      2328.3,
			Low:       2318.6,
			Close:     2325.5,
		}},
	})
	if !errors.Is(err, ErrSymbolNotConfigured) {
		t.Fatalf("IngestCandles error = %v, want ErrSymbolNotConfigured", err)
	}
}

type memoryMT5Repository struct {
	heartbeat domainmt5.Heartbeat
	ticks     []domainmt5.Tick
	account   domainmt5.AccountSnapshot
	positions []domainmt5.PositionSnapshot
}

func (r *memoryMT5Repository) SaveHeartbeat(ctx context.Context, heartbeat domainmt5.Heartbeat) error {
	r.heartbeat = heartbeat
	return nil
}

func (r *memoryMT5Repository) LatestHeartbeat(ctx context.Context, bridgeID string) (domainmt5.Heartbeat, error) {
	return r.heartbeat, nil
}

func (r *memoryMT5Repository) SaveTicks(ctx context.Context, ticks []domainmt5.Tick) error {
	r.ticks = append([]domainmt5.Tick(nil), ticks...)
	return nil
}

func (r *memoryMT5Repository) LatestTick(ctx context.Context, symbol string) (domainmt5.Tick, error) {
	if len(r.ticks) == 0 {
		return domainmt5.Tick{}, nil
	}
	return r.ticks[len(r.ticks)-1], nil
}

func (r *memoryMT5Repository) SaveAccountSnapshot(ctx context.Context, snapshot domainmt5.AccountSnapshot) error {
	r.account = snapshot
	return nil
}

func (r *memoryMT5Repository) LatestAccountSnapshot(ctx context.Context, accountLogin string) (domainmt5.AccountSnapshot, error) {
	return r.account, nil
}

func (r *memoryMT5Repository) SavePositionSnapshots(ctx context.Context, positions []domainmt5.PositionSnapshot) error {
	r.positions = append([]domainmt5.PositionSnapshot(nil), positions...)
	return nil
}

func (r *memoryMT5Repository) LatestPositionSnapshots(ctx context.Context, accountLogin string) ([]domainmt5.PositionSnapshot, error) {
	return append([]domainmt5.PositionSnapshot(nil), r.positions...), nil
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

type memoryCandleRepository struct {
	candles []market.Candle
}

func (r *memoryCandleRepository) UpsertCandles(ctx context.Context, candles []market.Candle) error {
	r.candles = append([]market.Candle(nil), candles...)
	return nil
}

func (r *memoryCandleRepository) ListCandles(ctx context.Context, query market.CandleQuery) ([]market.Candle, error) {
	return append([]market.Candle(nil), r.candles...), nil
}

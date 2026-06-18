package mt5

import (
	"context"
	"time"
)

const XAUUSDSymbol = "XAUUSD"

type Heartbeat struct {
	BridgeID     string
	Terminal     string
	AccountLogin string
	Server       string
	Status       string
	LastError    string
	SentAt       time.Time
}

type Tick struct {
	Symbol string
	Bid    float64
	Ask    float64
	Last   float64
	Volume float64
	Time   time.Time
}

type CandleBatch struct {
	Symbol    string
	Timeframe string
	Source    string
	Candles   []Candle
}

type Candle struct {
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type AccountSnapshot struct {
	AccountLogin string
	Currency     string
	Balance      float64
	Equity       float64
	Margin       float64
	FreeMargin   float64
	MarginLevel  float64
	Time         time.Time
}

type PositionSnapshot struct {
	AccountLogin string
	Ticket       string
	Symbol       string
	Side         string
	Volume       float64
	OpenPrice    float64
	StopLoss     float64
	TakeProfit   float64
	Profit       float64
	OpenedAt     time.Time
	SnapshotTime time.Time
}

type Repository interface {
	SaveHeartbeat(ctx context.Context, heartbeat Heartbeat) error
	LatestHeartbeat(ctx context.Context, bridgeID string) (Heartbeat, error)
	SaveTicks(ctx context.Context, ticks []Tick) error
	LatestTick(ctx context.Context, symbol string) (Tick, error)
	SaveAccountSnapshot(ctx context.Context, snapshot AccountSnapshot) error
	LatestAccountSnapshot(ctx context.Context, accountLogin string) (AccountSnapshot, error)
	SavePositionSnapshots(ctx context.Context, positions []PositionSnapshot) error
	LatestPositionSnapshots(ctx context.Context, accountLogin string) ([]PositionSnapshot, error)
}

package market

import (
	"context"
	"time"
)

type Candle struct {
	SymbolID  string
	Timeframe string
	Timestamp time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type CandleQuery struct {
	SymbolCode string
	Timeframe  string
	From       time.Time
	To         time.Time
}

type CandleRepository interface {
	UpsertCandles(ctx context.Context, candles []Candle) error
	ListCandles(ctx context.Context, query CandleQuery) ([]Candle, error)
}

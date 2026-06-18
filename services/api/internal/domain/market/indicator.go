package market

import (
	"context"
	"time"
)

type IndicatorQuery struct {
	SymbolCode string
	Timeframe  string
}

type IndicatorPoint struct {
	Timestamp time.Time
	Value     float64
}

type MACDPoint struct {
	Timestamp time.Time
	MACD      float64
	Signal    float64
	Histogram float64
}

type IndicatorSeries struct {
	SymbolCode string
	Timeframe  string
	EMA20      []IndicatorPoint
	EMA50      []IndicatorPoint
	RSI14      []IndicatorPoint
	MACD       []MACDPoint
	ATR14      []IndicatorPoint
}

type IndicatorService interface {
	ListIndicators(ctx context.Context, query IndicatorQuery) (IndicatorSeries, error)
}

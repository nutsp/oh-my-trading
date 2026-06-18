package mockdata

import (
	"context"
	"fmt"
	"hash/fnv"
	"math"
	"strings"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market/indicators"
)

type IDGenerator func() string

type SymbolService struct {
	generateID IDGenerator
}

func NewSymbolService(generateID IDGenerator) *SymbolService {
	return &SymbolService{generateID: generateID}
}

func (s *SymbolService) ListSymbols(ctx context.Context) ([]market.Symbol, error) {
	return []market.Symbol{
		{
			ID:         "mock-xauusd",
			Code:       "XAUUSD",
			Market:     "forex",
			BaseAsset:  "XAU",
			QuoteAsset: "USD",
			Enabled:    true,
		},
		{
			ID:         "mock-btcusd",
			Code:       "BTCUSD",
			Market:     "crypto",
			BaseAsset:  "BTC",
			QuoteAsset: "USD",
			Enabled:    true,
		},
		{
			ID:         "mock-eurusd",
			Code:       "EURUSD",
			Market:     "forex",
			BaseAsset:  "EUR",
			QuoteAsset: "USD",
			Enabled:    true,
		},
	}, nil
}

func (s *SymbolService) CreateSymbol(ctx context.Context, input market.CreateSymbolInput) (market.Symbol, error) {
	id := "mock-symbol"
	if s.generateID != nil {
		id = s.generateID()
	}

	return market.Symbol{
		ID:         id,
		Code:       strings.ToUpper(input.Code),
		Market:     input.Market,
		BaseAsset:  strings.ToUpper(input.BaseAsset),
		QuoteAsset: strings.ToUpper(input.QuoteAsset),
		Enabled:    true,
	}, nil
}

type CandleService struct{}

func NewCandleService() *CandleService {
	return &CandleService{}
}

func (s *CandleService) ListCandles(ctx context.Context, query market.CandleQuery) ([]market.Candle, error) {
	return generateCandles(query), nil
}

type IndicatorService struct{}

func NewIndicatorService() *IndicatorService {
	return &IndicatorService{}
}

func (s *IndicatorService) ListIndicators(ctx context.Context, query market.IndicatorQuery) (market.IndicatorSeries, error) {
	now := time.Now().UTC()
	candles := generateCandles(market.CandleQuery{
		SymbolCode: query.SymbolCode,
		Timeframe:  query.Timeframe,
		From:       now.Add(-90 * 24 * time.Hour),
		To:         now,
	})

	return market.IndicatorSeries{
		SymbolCode: query.SymbolCode,
		Timeframe:  query.Timeframe,
		EMA20:      indicators.EMA(candles, 20),
		EMA50:      indicators.EMA(candles, 50),
		RSI14:      indicators.RSI(candles, 14),
		MACD:       indicators.MACD(candles, 12, 26, 9),
		ATR14:      indicators.ATR(candles, 14),
	}, nil
}

func generateCandles(query market.CandleQuery) []market.Candle {
	step := timeframeDuration(query.Timeframe)
	if !query.To.After(query.From) || step <= 0 {
		return nil
	}

	maxPoints := 600
	duration := query.To.Sub(query.From)
	points := int(duration / step)
	if points <= 0 {
		points = 1
	}
	if points > maxPoints {
		points = maxPoints
	}

	base := basePriceBySymbol(query.SymbolCode)
	phase := symbolPhase(query.SymbolCode)
	candles := make([]market.Candle, 0, points)

	current := query.From.UTC()
	for i := 0; i < points; i++ {
		// Deterministic trend + cycle so charts feel alive but reproducible.
		trend := float64(i) * 0.12
		wave := math.Sin(float64(i)/6.0+phase) * 6.5
		center := base + trend + wave
		body := math.Sin(float64(i)/3.0+phase) * 1.6

		open := center - body
		close := center + body
		high := math.Max(open, close) + 2.4 + math.Abs(math.Cos(float64(i)/5.0))*0.8
		low := math.Min(open, close) - 2.1 - math.Abs(math.Sin(float64(i)/5.0))*0.7
		volume := 1000 + 250*math.Abs(math.Sin(float64(i)/4.0+phase))

		candles = append(candles, market.Candle{
			SymbolID:  fmt.Sprintf("mock-%s", strings.ToLower(query.SymbolCode)),
			Timeframe: query.Timeframe,
			Timestamp: current,
			Open:      round2(open),
			High:      round2(high),
			Low:       round2(low),
			Close:     round2(close),
			Volume:    round2(volume),
		})
		current = current.Add(step)
	}

	return candles
}

func timeframeDuration(timeframe string) time.Duration {
	switch timeframe {
	case "1m":
		return time.Minute
	case "5m":
		return 5 * time.Minute
	case "15m":
		return 15 * time.Minute
	case "1h":
		return time.Hour
	case "4h":
		return 4 * time.Hour
	case "1d":
		return 24 * time.Hour
	default:
		return time.Hour
	}
}

func basePriceBySymbol(symbol string) float64 {
	switch strings.ToUpper(symbol) {
	case "XAUUSD":
		return 2320
	case "BTCUSD":
		return 66000
	case "EURUSD":
		return 1.09
	default:
		return 100
	}
}

func symbolPhase(symbol string) float64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(strings.ToUpper(symbol)))
	return float64(h.Sum32()%360) * (math.Pi / 180.0)
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

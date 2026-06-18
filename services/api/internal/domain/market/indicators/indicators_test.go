package indicators

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

func TestIndicatorsMatchGoldenFixture(t *testing.T) {
	golden := loadGolden(t)
	candles := fixtureCandles()

	ema20 := EMA(candles, 20)
	assertIndicatorPoints(t, "ema20", ema20, golden.EMA20.Length, golden.EMA20.Last3)

	ema50 := EMA(candles, 50)
	assertIndicatorPoints(t, "ema50", ema50, golden.EMA50.Length, golden.EMA50.Last3)

	rsi14 := RSI(candles, 14)
	assertIndicatorPoints(t, "rsi14", rsi14, golden.RSI14.Length, golden.RSI14.Last3)

	atr14 := ATR(candles, 14)
	assertIndicatorPoints(t, "atr14", atr14, golden.ATR14.Length, golden.ATR14.Last3)

	macd := MACD(candles, 12, 26, 9)
	if len(macd) != golden.MACD.Length {
		t.Fatalf("macd length = %d, want %d", len(macd), golden.MACD.Length)
	}
	for i, expected := range golden.MACD.Last3 {
		actual := macd[len(macd)-3+i]
		assertAlmostEqual(t, "macd", actual.MACD, expected.MACD)
		assertAlmostEqual(t, "macd signal", actual.Signal, expected.Signal)
		assertAlmostEqual(t, "macd histogram", actual.Histogram, expected.Histogram)
	}
}

type indicatorGolden struct {
	EMA20 pointGolden `json:"ema20"`
	EMA50 pointGolden `json:"ema50"`
	RSI14 pointGolden `json:"rsi14"`
	ATR14 pointGolden `json:"atr14"`
	MACD  macdGolden  `json:"macd"`
}

type pointGolden struct {
	Length int       `json:"length"`
	Last3  []float64 `json:"last3"`
}

type macdGolden struct {
	Length int               `json:"length"`
	Last3  []macdPointGolden `json:"last3"`
}

type macdPointGolden struct {
	MACD      float64 `json:"macd"`
	Signal    float64 `json:"signal"`
	Histogram float64 `json:"histogram"`
}

func loadGolden(t *testing.T) indicatorGolden {
	t.Helper()

	path := filepath.Join("testdata", "golden.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden fixture: %v", err)
	}

	var golden indicatorGolden
	if err := json.Unmarshal(raw, &golden); err != nil {
		t.Fatalf("decode golden fixture: %v", err)
	}
	return golden
}

func fixtureCandles() []market.Candle {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	close := 2300.0

	candles := make([]market.Candle, 0, 80)
	for i := 0; i < 80; i++ {
		drift := 0.9
		if i%7 == 0 || i%7 == 1 {
			drift = -0.4
		}

		wave := float64((i%5)-2) * 0.3
		close += drift + wave

		open := close - 0.5
		high := math.Max(open, close) + 1.2
		low := math.Min(open, close) - 1.0

		candles = append(candles, market.Candle{
			Timeframe: "1h",
			Timestamp: start.Add(time.Duration(i) * time.Hour),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
		})
	}
	return candles
}

func assertIndicatorPoints(t *testing.T, name string, points []market.IndicatorPoint, expectedLength int, expectedLast3 []float64) {
	t.Helper()

	if len(points) != expectedLength {
		t.Fatalf("%s length = %d, want %d", name, len(points), expectedLength)
	}
	for i, expected := range expectedLast3 {
		actual := points[len(points)-3+i].Value
		assertAlmostEqual(t, name, actual, expected)
	}
}

func assertAlmostEqual(t *testing.T, metric string, actual, expected float64) {
	t.Helper()
	const epsilon = 0.000001
	if math.Abs(actual-expected) > epsilon {
		t.Fatalf("%s = %.6f, want %.6f", metric, actual, expected)
	}
}

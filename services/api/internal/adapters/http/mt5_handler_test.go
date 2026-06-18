package httpadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appmt5 "github.com/sutad-p/oh-my-trading/services/api/internal/application/mt5"
	domainmt5 "github.com/sutad-p/oh-my-trading/services/api/internal/domain/mt5"
)

func TestMT5RoutesIngestAndReadStatus(t *testing.T) {
	service := &fakeMT5Service{
		heartbeat: domainmt5.Heartbeat{
			BridgeID:     "local-mt5",
			Terminal:     "MetaTrader 5",
			AccountLogin: "12345678",
			Server:       "Broker-Demo",
			Status:       "healthy",
			SentAt:       time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC),
		},
		tick: domainmt5.Tick{
			Symbol: "XAUUSD",
			Bid:    2325.42,
			Ask:    2325.62,
			Last:   2325.52,
			Volume: 12,
			Time:   time.Date(2026, 6, 18, 10, 0, 1, 0, time.UTC),
		},
	}
	router := NewRouter(WithMT5Service(service))

	rec := performJSON(router, http.MethodPost, "/api/mt5/heartbeat", `{
		"bridgeId":"local-mt5",
		"terminal":"MetaTrader 5",
		"accountLogin":"12345678",
		"server":"Broker-Demo",
		"status":"healthy",
		"sentAt":"2026-06-18T10:00:00Z"
	}`)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("heartbeat status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if service.ingestedHeartbeat.BridgeID != "local-mt5" {
		t.Fatalf("ingested heartbeat BridgeID = %q", service.ingestedHeartbeat.BridgeID)
	}

	rec = performJSON(router, http.MethodPost, "/api/mt5/ticks", `{
		"ticks":[{"symbol":"XAUUSD","bid":2325.42,"ask":2325.62,"last":2325.52,"volume":12,"time":"2026-06-18T10:00:01Z"}]
	}`)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("ticks status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if len(service.ingestedTicks) != 1 {
		t.Fatalf("len(ingestedTicks) = %d, want 1", len(service.ingestedTicks))
	}

	req := httptest.NewRequest(http.MethodGet, "/api/mt5/status", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status route code = %d, want %d", rec.Code, http.StatusOK)
	}
	var response mt5StatusResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode status response: %v", err)
	}
	if response.Heartbeat.Status != "healthy" {
		t.Fatalf("heartbeat status = %q, want healthy", response.Heartbeat.Status)
	}
	if response.LatestTick.Ask != 2325.62 {
		t.Fatalf("latest tick ask = %f, want 2325.62", response.LatestTick.Ask)
	}
}

func TestMT5RoutesIngestCandlesAccountAndPositions(t *testing.T) {
	service := &fakeMT5Service{
		account: domainmt5.AccountSnapshot{
			AccountLogin: "12345678",
			Currency:     "USD",
			Balance:      10000,
			Equity:       10080,
			Time:         time.Date(2026, 6, 18, 10, 2, 0, 0, time.UTC),
		},
		positions: []domainmt5.PositionSnapshot{{
			AccountLogin: "12345678",
			Ticket:       "987654321",
			Symbol:       "XAUUSD",
			Side:         "buy",
			Volume:       0.1,
			OpenPrice:    2320.1,
			Profit:       55,
			OpenedAt:     time.Date(2026, 6, 18, 9, 0, 0, 0, time.UTC),
			SnapshotTime: time.Date(2026, 6, 18, 10, 3, 0, 0, time.UTC),
		}},
	}
	router := NewRouter(WithMT5Service(service))

	rec := performJSON(router, http.MethodPost, "/api/mt5/candles", `{
		"symbol":"XAUUSD",
		"timeframe":"1m",
		"source":"mt5-python-bridge",
		"candles":[{"timestamp":"2026-06-18T10:01:00Z","open":2320.1,"high":2328.3,"low":2318.6,"close":2325.5,"volume":12345}]
	}`)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("candles status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if service.ingestedCandles.Symbol != "XAUUSD" {
		t.Fatalf("ingested candle symbol = %q, want XAUUSD", service.ingestedCandles.Symbol)
	}

	rec = performJSON(router, http.MethodPost, "/api/mt5/account-snapshot", `{
		"accountLogin":"12345678",
		"currency":"USD",
		"balance":10000,
		"equity":10080,
		"margin":400,
		"freeMargin":9680,
		"marginLevel":2520,
		"time":"2026-06-18T10:02:00Z"
	}`)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("account status = %d, want %d", rec.Code, http.StatusAccepted)
	}

	rec = performJSON(router, http.MethodPost, "/api/mt5/positions", `{
		"accountLogin":"12345678",
		"time":"2026-06-18T10:03:00Z",
		"positions":[{"ticket":"987654321","symbol":"XAUUSD","side":"buy","volume":0.1,"openPrice":2320.1,"stopLoss":2310,"takeProfit":2340,"profit":55,"openedAt":"2026-06-18T09:00:00Z"}]
	}`)
	if rec.Code != http.StatusAccepted {
		t.Fatalf("positions status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if len(service.ingestedPositions) != 1 {
		t.Fatalf("len(ingestedPositions) = %d, want 1", len(service.ingestedPositions))
	}

	req := httptest.NewRequest(http.MethodGet, "/api/mt5/account/latest?accountLogin=12345678", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("account latest status = %d, want %d", rec.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/mt5/positions/latest?accountLogin=12345678", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("positions latest status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestMT5RoutesReturnBadRequestForUnsupportedSymbol(t *testing.T) {
	router := NewRouter(WithMT5Service(&fakeMT5Service{err: appmt5.ErrUnsupportedSymbol}))

	rec := performJSON(router, http.MethodPost, "/api/mt5/ticks", `{
		"ticks":[{"symbol":"EURUSD","bid":1.1,"ask":1.2,"time":"2026-06-18T10:00:01Z"}]
	}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func performJSON(handler http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

type fakeMT5Service struct {
	err               error
	heartbeat         domainmt5.Heartbeat
	tick              domainmt5.Tick
	account           domainmt5.AccountSnapshot
	positions         []domainmt5.PositionSnapshot
	ingestedHeartbeat domainmt5.Heartbeat
	ingestedTicks     []domainmt5.Tick
	ingestedCandles   domainmt5.CandleBatch
	ingestedAccount   domainmt5.AccountSnapshot
	ingestedPositions []domainmt5.PositionSnapshot
}

func (s *fakeMT5Service) IngestHeartbeat(ctx context.Context, heartbeat domainmt5.Heartbeat) error {
	s.ingestedHeartbeat = heartbeat
	return s.err
}

func (s *fakeMT5Service) IngestTicks(ctx context.Context, ticks []domainmt5.Tick) error {
	s.ingestedTicks = append([]domainmt5.Tick(nil), ticks...)
	return s.err
}

func (s *fakeMT5Service) IngestCandles(ctx context.Context, batch domainmt5.CandleBatch) error {
	s.ingestedCandles = batch
	return s.err
}

func (s *fakeMT5Service) IngestAccountSnapshot(ctx context.Context, snapshot domainmt5.AccountSnapshot) error {
	s.ingestedAccount = snapshot
	return s.err
}

func (s *fakeMT5Service) IngestPositionSnapshots(ctx context.Context, positions []domainmt5.PositionSnapshot) error {
	s.ingestedPositions = append([]domainmt5.PositionSnapshot(nil), positions...)
	return s.err
}

func (s *fakeMT5Service) LatestHeartbeat(ctx context.Context, bridgeID string) (domainmt5.Heartbeat, error) {
	return s.heartbeat, s.err
}

func (s *fakeMT5Service) LatestTick(ctx context.Context, symbol string) (domainmt5.Tick, error) {
	return s.tick, s.err
}

func (s *fakeMT5Service) LatestAccountSnapshot(ctx context.Context, accountLogin string) (domainmt5.AccountSnapshot, error) {
	return s.account, s.err
}

func (s *fakeMT5Service) LatestPositionSnapshots(ctx context.Context, accountLogin string) ([]domainmt5.PositionSnapshot, error) {
	return append([]domainmt5.PositionSnapshot(nil), s.positions...), s.err
}

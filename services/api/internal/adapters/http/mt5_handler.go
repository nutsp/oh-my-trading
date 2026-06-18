package httpadapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	appmt5 "github.com/sutad-p/oh-my-trading/services/api/internal/application/mt5"
	domainmt5 "github.com/sutad-p/oh-my-trading/services/api/internal/domain/mt5"
)

const (
	defaultMT5BridgeID = "local-mt5"
	defaultMT5Symbol   = domainmt5.XAUUSDSymbol
)

type mt5Service interface {
	IngestHeartbeat(ctx context.Context, heartbeat domainmt5.Heartbeat) error
	IngestTicks(ctx context.Context, ticks []domainmt5.Tick) error
	IngestCandles(ctx context.Context, batch domainmt5.CandleBatch) error
	IngestAccountSnapshot(ctx context.Context, snapshot domainmt5.AccountSnapshot) error
	IngestPositionSnapshots(ctx context.Context, positions []domainmt5.PositionSnapshot) error
	LatestHeartbeat(ctx context.Context, bridgeID string) (domainmt5.Heartbeat, error)
	LatestTick(ctx context.Context, symbol string) (domainmt5.Tick, error)
	LatestAccountSnapshot(ctx context.Context, accountLogin string) (domainmt5.AccountSnapshot, error)
	LatestPositionSnapshots(ctx context.Context, accountLogin string) ([]domainmt5.PositionSnapshot, error)
}

type mt5HeartbeatRequest struct {
	BridgeID     string    `json:"bridgeId"`
	Terminal     string    `json:"terminal"`
	AccountLogin string    `json:"accountLogin"`
	Server       string    `json:"server"`
	Status       string    `json:"status"`
	LastError    string    `json:"lastError"`
	SentAt       time.Time `json:"sentAt"`
}

type mt5TickRequest struct {
	Symbol string    `json:"symbol"`
	Bid    float64   `json:"bid"`
	Ask    float64   `json:"ask"`
	Last   float64   `json:"last"`
	Volume float64   `json:"volume"`
	Time   time.Time `json:"time"`
}

type mt5TicksRequest struct {
	Ticks []mt5TickRequest `json:"ticks"`
}

type mt5CandleRequest struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
}

type mt5CandlesRequest struct {
	Symbol    string             `json:"symbol"`
	Timeframe string             `json:"timeframe"`
	Source    string             `json:"source"`
	Candles   []mt5CandleRequest `json:"candles"`
}

type mt5AccountSnapshotRequest struct {
	AccountLogin string    `json:"accountLogin"`
	Currency     string    `json:"currency"`
	Balance      float64   `json:"balance"`
	Equity       float64   `json:"equity"`
	Margin       float64   `json:"margin"`
	FreeMargin   float64   `json:"freeMargin"`
	MarginLevel  float64   `json:"marginLevel"`
	Time         time.Time `json:"time"`
}

type mt5PositionRequest struct {
	Ticket     string    `json:"ticket"`
	Symbol     string    `json:"symbol"`
	Side       string    `json:"side"`
	Volume     float64   `json:"volume"`
	OpenPrice  float64   `json:"openPrice"`
	StopLoss   float64   `json:"stopLoss"`
	TakeProfit float64   `json:"takeProfit"`
	Profit     float64   `json:"profit"`
	OpenedAt   time.Time `json:"openedAt"`
}

type mt5PositionsRequest struct {
	AccountLogin string               `json:"accountLogin"`
	Time         time.Time            `json:"time"`
	Positions    []mt5PositionRequest `json:"positions"`
}

type mt5HeartbeatResponse struct {
	BridgeID     string `json:"bridgeId"`
	Terminal     string `json:"terminal"`
	AccountLogin string `json:"accountLogin"`
	Server       string `json:"server"`
	Status       string `json:"status"`
	LastError    string `json:"lastError,omitempty"`
	SentAt       string `json:"sentAt"`
}

type mt5TickResponse struct {
	Symbol string  `json:"symbol"`
	Bid    float64 `json:"bid"`
	Ask    float64 `json:"ask"`
	Last   float64 `json:"last"`
	Volume float64 `json:"volume"`
	Time   string  `json:"time"`
}

type mt5StatusResponse struct {
	State      string               `json:"state"`
	Heartbeat  mt5HeartbeatResponse `json:"heartbeat"`
	LatestTick mt5TickResponse      `json:"latestTick"`
}

type mt5AccountSnapshotResponse struct {
	AccountLogin string  `json:"accountLogin"`
	Currency     string  `json:"currency"`
	Balance      float64 `json:"balance"`
	Equity       float64 `json:"equity"`
	Margin       float64 `json:"margin"`
	FreeMargin   float64 `json:"freeMargin"`
	MarginLevel  float64 `json:"marginLevel"`
	Time         string  `json:"time"`
}

type mt5PositionResponse struct {
	AccountLogin string  `json:"accountLogin"`
	Ticket       string  `json:"ticket"`
	Symbol       string  `json:"symbol"`
	Side         string  `json:"side"`
	Volume       float64 `json:"volume"`
	OpenPrice    float64 `json:"openPrice"`
	StopLoss     float64 `json:"stopLoss"`
	TakeProfit   float64 `json:"takeProfit"`
	Profit       float64 `json:"profit"`
	OpenedAt     string  `json:"openedAt"`
	SnapshotTime string  `json:"snapshotTime"`
}

func mt5HeartbeatHandler(service mt5Service) http.HandlerFunc {
	return mt5PostHandler(service, func(w http.ResponseWriter, r *http.Request, service mt5Service) error {
		var request mt5HeartbeatRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return nil
		}
		err := service.IngestHeartbeat(r.Context(), domainmt5.Heartbeat{
			BridgeID:     request.BridgeID,
			Terminal:     request.Terminal,
			AccountLogin: request.AccountLogin,
			Server:       request.Server,
			Status:       request.Status,
			LastError:    request.LastError,
			SentAt:       request.SentAt,
		})
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusAccepted)
		return nil
	})
}

func mt5TicksHandler(service mt5Service) http.HandlerFunc {
	return mt5PostHandler(service, func(w http.ResponseWriter, r *http.Request, service mt5Service) error {
		var request mt5TicksRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return nil
		}
		err := service.IngestTicks(r.Context(), mapMT5TickRequests(request.Ticks))
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusAccepted)
		return nil
	})
}

func mt5CandlesHandler(service mt5Service) http.HandlerFunc {
	return mt5PostHandler(service, func(w http.ResponseWriter, r *http.Request, service mt5Service) error {
		var request mt5CandlesRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return nil
		}
		err := service.IngestCandles(r.Context(), domainmt5.CandleBatch{
			Symbol:    request.Symbol,
			Timeframe: request.Timeframe,
			Source:    request.Source,
			Candles:   mapMT5CandleRequests(request.Candles),
		})
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusAccepted)
		return nil
	})
}

func mt5AccountSnapshotHandler(service mt5Service) http.HandlerFunc {
	return mt5PostHandler(service, func(w http.ResponseWriter, r *http.Request, service mt5Service) error {
		var request mt5AccountSnapshotRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return nil
		}
		err := service.IngestAccountSnapshot(r.Context(), domainmt5.AccountSnapshot{
			AccountLogin: request.AccountLogin,
			Currency:     request.Currency,
			Balance:      request.Balance,
			Equity:       request.Equity,
			Margin:       request.Margin,
			FreeMargin:   request.FreeMargin,
			MarginLevel:  request.MarginLevel,
			Time:         request.Time,
		})
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusAccepted)
		return nil
	})
}

func mt5PositionsHandler(service mt5Service) http.HandlerFunc {
	return mt5PostHandler(service, func(w http.ResponseWriter, r *http.Request, service mt5Service) error {
		var request mt5PositionsRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return nil
		}
		err := service.IngestPositionSnapshots(r.Context(), mapMT5PositionRequests(request))
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusAccepted)
		return nil
	})
}

func mt5StatusHandler(service mt5Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "mt5 service is not configured", http.StatusNotImplemented)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		heartbeat, err := service.LatestHeartbeat(r.Context(), defaultMT5BridgeID)
		if err != nil {
			if isMT5WaitingForBridge(err) {
				writeJSON(w, http.StatusOK, mt5WaitingStatusResponse())
				return
			}
			http.Error(w, "latest mt5 heartbeat", http.StatusInternalServerError)
			return
		}
		tick, err := service.LatestTick(r.Context(), defaultMT5Symbol)
		if err != nil {
			if isMT5WaitingForBridge(err) {
				writeJSON(w, http.StatusOK, mt5WaitingStatusResponse())
				return
			}
			http.Error(w, "latest mt5 tick", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, mt5StatusResponse{
			State:      "connected",
			Heartbeat:  mapMT5Heartbeat(heartbeat),
			LatestTick: mapMT5Tick(tick),
		})
	}
}

func isMT5WaitingForBridge(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func mt5WaitingStatusResponse() mt5StatusResponse {
	return mt5StatusResponse{
		State: "waiting_for_bridge",
		Heartbeat: mt5HeartbeatResponse{
			BridgeID: defaultMT5BridgeID,
			Status:   "disconnected",
		},
		LatestTick: mt5TickResponse{
			Symbol: defaultMT5Symbol,
		},
	}
}

func mt5LatestAccountHandler(service mt5Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "mt5 service is not configured", http.StatusNotImplemented)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		accountLogin := r.URL.Query().Get("accountLogin")
		if accountLogin == "" {
			http.Error(w, "accountLogin is required", http.StatusBadRequest)
			return
		}
		account, err := service.LatestAccountSnapshot(r.Context(), accountLogin)
		if err != nil {
			http.Error(w, "latest mt5 account snapshot", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, mapMT5AccountSnapshot(account))
	}
}

func mt5LatestPositionsHandler(service mt5Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "mt5 service is not configured", http.StatusNotImplemented)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		accountLogin := r.URL.Query().Get("accountLogin")
		if accountLogin == "" {
			http.Error(w, "accountLogin is required", http.StatusBadRequest)
			return
		}
		positions, err := service.LatestPositionSnapshots(r.Context(), accountLogin)
		if err != nil {
			http.Error(w, "latest mt5 position snapshots", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, mapMT5Positions(positions))
	}
}

func mt5PostHandler(service mt5Service, handle func(http.ResponseWriter, *http.Request, mt5Service) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "mt5 service is not configured", http.StatusNotImplemented)
			return
		}
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := handle(w, r, service); err != nil {
			writeMT5Error(w, err)
		}
	}
}

func writeMT5Error(w http.ResponseWriter, err error) {
	if errors.Is(err, appmt5.ErrUnsupportedSymbol) || errors.Is(err, appmt5.ErrSymbolNotConfigured) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Error(w, "mt5 ingest failed", http.StatusInternalServerError)
}

func mapMT5TickRequests(requests []mt5TickRequest) []domainmt5.Tick {
	ticks := make([]domainmt5.Tick, 0, len(requests))
	for _, request := range requests {
		ticks = append(ticks, domainmt5.Tick{
			Symbol: request.Symbol,
			Bid:    request.Bid,
			Ask:    request.Ask,
			Last:   request.Last,
			Volume: request.Volume,
			Time:   request.Time,
		})
	}
	return ticks
}

func mapMT5CandleRequests(requests []mt5CandleRequest) []domainmt5.Candle {
	candles := make([]domainmt5.Candle, 0, len(requests))
	for _, request := range requests {
		candles = append(candles, domainmt5.Candle{
			Timestamp: request.Timestamp,
			Open:      request.Open,
			High:      request.High,
			Low:       request.Low,
			Close:     request.Close,
			Volume:    request.Volume,
		})
	}
	return candles
}

func mapMT5PositionRequests(request mt5PositionsRequest) []domainmt5.PositionSnapshot {
	positions := make([]domainmt5.PositionSnapshot, 0, len(request.Positions))
	for _, position := range request.Positions {
		positions = append(positions, domainmt5.PositionSnapshot{
			AccountLogin: request.AccountLogin,
			Ticket:       position.Ticket,
			Symbol:       position.Symbol,
			Side:         position.Side,
			Volume:       position.Volume,
			OpenPrice:    position.OpenPrice,
			StopLoss:     position.StopLoss,
			TakeProfit:   position.TakeProfit,
			Profit:       position.Profit,
			OpenedAt:     position.OpenedAt,
			SnapshotTime: request.Time,
		})
	}
	return positions
}

func mapMT5Heartbeat(heartbeat domainmt5.Heartbeat) mt5HeartbeatResponse {
	return mt5HeartbeatResponse{
		BridgeID:     heartbeat.BridgeID,
		Terminal:     heartbeat.Terminal,
		AccountLogin: heartbeat.AccountLogin,
		Server:       heartbeat.Server,
		Status:       heartbeat.Status,
		LastError:    heartbeat.LastError,
		SentAt:       formatRFC3339(heartbeat.SentAt),
	}
}

func mapMT5Tick(tick domainmt5.Tick) mt5TickResponse {
	return mt5TickResponse{
		Symbol: tick.Symbol,
		Bid:    tick.Bid,
		Ask:    tick.Ask,
		Last:   tick.Last,
		Volume: tick.Volume,
		Time:   formatRFC3339(tick.Time),
	}
}

func mapMT5AccountSnapshot(snapshot domainmt5.AccountSnapshot) mt5AccountSnapshotResponse {
	return mt5AccountSnapshotResponse{
		AccountLogin: snapshot.AccountLogin,
		Currency:     snapshot.Currency,
		Balance:      snapshot.Balance,
		Equity:       snapshot.Equity,
		Margin:       snapshot.Margin,
		FreeMargin:   snapshot.FreeMargin,
		MarginLevel:  snapshot.MarginLevel,
		Time:         formatRFC3339(snapshot.Time),
	}
}

func mapMT5Positions(positions []domainmt5.PositionSnapshot) []mt5PositionResponse {
	response := make([]mt5PositionResponse, 0, len(positions))
	for _, position := range positions {
		response = append(response, mt5PositionResponse{
			AccountLogin: position.AccountLogin,
			Ticket:       position.Ticket,
			Symbol:       position.Symbol,
			Side:         position.Side,
			Volume:       position.Volume,
			OpenPrice:    position.OpenPrice,
			StopLoss:     position.StopLoss,
			TakeProfit:   position.TakeProfit,
			Profit:       position.Profit,
			OpenedAt:     formatRFC3339(position.OpenedAt),
			SnapshotTime: formatRFC3339(position.SnapshotTime),
		})
	}
	return response
}

func formatRFC3339(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

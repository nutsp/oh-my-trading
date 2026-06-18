package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

type symbolService interface {
	ListSymbols(ctx context.Context) ([]market.Symbol, error)
	CreateSymbol(ctx context.Context, input market.CreateSymbolInput) (market.Symbol, error)
}

type candleService interface {
	ListCandles(ctx context.Context, query market.CandleQuery) ([]market.Candle, error)
}

type routerConfig struct {
	symbols symbolService
	candles candleService
}

type Option func(*routerConfig)

func WithSymbolService(service symbolService) Option {
	return func(cfg *routerConfig) {
		cfg.symbols = service
	}
}

func WithCandleService(service candleService) Option {
	return func(cfg *routerConfig) {
		cfg.candles = service
	}
}

func NewRouter(options ...Option) http.Handler {
	var cfg routerConfig
	for _, option := range options {
		option(&cfg)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/symbols", symbolsHandler(cfg.symbols))
	mux.HandleFunc("/api/candles", candlesHandler(cfg.candles))
	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

type symbolResponse struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	Market     string `json:"market"`
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
	Enabled    bool   `json:"enabled"`
}

type createSymbolRequest struct {
	Code       string `json:"code"`
	Market     string `json:"market"`
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
}

func symbolsHandler(service symbolService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "symbol service is not configured", http.StatusNotImplemented)
			return
		}

		switch r.Method {
		case http.MethodGet:
			listSymbols(w, r, service)
		case http.MethodPost:
			createSymbol(w, r, service)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func listSymbols(w http.ResponseWriter, r *http.Request, service symbolService) {
	symbols, err := service.ListSymbols(r.Context())
	if err != nil {
		http.Error(w, "list symbols", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, mapSymbols(symbols))
}

func createSymbol(w http.ResponseWriter, r *http.Request, service symbolService) {
	var request createSymbolRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	symbol, err := service.CreateSymbol(r.Context(), market.CreateSymbolInput{
		Code:       request.Code,
		Market:     request.Market,
		BaseAsset:  request.BaseAsset,
		QuoteAsset: request.QuoteAsset,
	})
	if err != nil {
		http.Error(w, "create symbol", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, mapSymbol(symbol))
}

func mapSymbols(symbols []market.Symbol) []symbolResponse {
	response := make([]symbolResponse, 0, len(symbols))
	for _, symbol := range symbols {
		response = append(response, mapSymbol(symbol))
	}
	return response
}

func mapSymbol(symbol market.Symbol) symbolResponse {
	return symbolResponse{
		ID:         symbol.ID,
		Code:       symbol.Code,
		Market:     symbol.Market,
		BaseAsset:  symbol.BaseAsset,
		QuoteAsset: symbol.QuoteAsset,
		Enabled:    symbol.Enabled,
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

type candleResponse struct {
	Timestamp string  `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
}

func candlesHandler(service candleService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "candle service is not configured", http.StatusNotImplemented)
			return
		}
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query, ok := parseCandleQuery(w, r)
		if !ok {
			return
		}

		candles, err := service.ListCandles(r.Context(), query)
		if err != nil {
			http.Error(w, "list candles", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, mapCandles(candles))
	}
}

func parseCandleQuery(w http.ResponseWriter, r *http.Request) (market.CandleQuery, bool) {
	values := r.URL.Query()
	symbol := values.Get("symbol")
	timeframe := values.Get("timeframe")
	fromRaw := values.Get("from")
	toRaw := values.Get("to")
	if symbol == "" || timeframe == "" || fromRaw == "" || toRaw == "" {
		http.Error(w, "symbol, timeframe, from, and to are required", http.StatusBadRequest)
		return market.CandleQuery{}, false
	}

	from, err := time.Parse(time.RFC3339, fromRaw)
	if err != nil {
		http.Error(w, "from must be RFC3339", http.StatusBadRequest)
		return market.CandleQuery{}, false
	}
	to, err := time.Parse(time.RFC3339, toRaw)
	if err != nil {
		http.Error(w, "to must be RFC3339", http.StatusBadRequest)
		return market.CandleQuery{}, false
	}

	return market.CandleQuery{
		SymbolCode: symbol,
		Timeframe:  timeframe,
		From:       from,
		To:         to,
	}, true
}

func mapCandles(candles []market.Candle) []candleResponse {
	response := make([]candleResponse, 0, len(candles))
	for _, candle := range candles {
		response = append(response, candleResponse{
			Timestamp: candle.Timestamp.UTC().Format(time.RFC3339),
			Open:      candle.Open,
			High:      candle.High,
			Low:       candle.Low,
			Close:     candle.Close,
			Volume:    candle.Volume,
		})
	}
	return response
}

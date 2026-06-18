package httpadapter

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

type symbolService interface {
	ListSymbols(ctx context.Context) ([]market.Symbol, error)
	CreateSymbol(ctx context.Context, input market.CreateSymbolInput) (market.Symbol, error)
}

type routerConfig struct {
	symbols symbolService
}

type Option func(*routerConfig)

func WithSymbolService(service symbolService) Option {
	return func(cfg *routerConfig) {
		cfg.symbols = service
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

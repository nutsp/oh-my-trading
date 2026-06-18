package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	appsignals "github.com/sutad-p/oh-my-trading/services/api/internal/application/signals"
	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/signal"
)

type signalService interface {
	CreateSignal(ctx context.Context, input signal.CreateInput) (signal.Signal, error)
	ListSignals(ctx context.Context) ([]signal.Signal, error)
	UpdateStatus(ctx context.Context, id string, status signal.Status) (signal.Signal, error)
}

type createPaperSignalRequest struct {
	Symbol     string  `json:"symbol"`
	Timeframe  string  `json:"timeframe"`
	Side       string  `json:"side"`
	Confidence float64 `json:"confidence"`
	EntryPrice float64 `json:"entryPrice"`
	StopLoss   float64 `json:"stopLoss"`
	TakeProfit float64 `json:"takeProfit"`
	Thesis     string  `json:"thesis"`
}

type updatePaperSignalStatusRequest struct {
	Status signal.Status `json:"status"`
}

type paperSignalResponse struct {
	ID         string  `json:"id"`
	Symbol     string  `json:"symbol"`
	Timeframe  string  `json:"timeframe"`
	Side       string  `json:"side"`
	Status     string  `json:"status"`
	Confidence float64 `json:"confidence"`
	EntryPrice float64 `json:"entryPrice"`
	StopLoss   float64 `json:"stopLoss"`
	TakeProfit float64 `json:"takeProfit"`
	Thesis     string  `json:"thesis"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

func paperSignalsHandler(service signalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "signal service is not configured", http.StatusNotImplemented)
			return
		}

		switch r.Method {
		case http.MethodGet:
			signals, err := service.ListSignals(r.Context())
			if err != nil {
				http.Error(w, "list paper signals", http.StatusInternalServerError)
				return
			}
			writeJSON(w, http.StatusOK, mapPaperSignals(signals))
		case http.MethodPost:
			var request createPaperSignalRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			created, err := service.CreateSignal(r.Context(), signal.CreateInput{
				Symbol:     request.Symbol,
				Timeframe:  request.Timeframe,
				Side:       request.Side,
				Confidence: request.Confidence,
				EntryPrice: request.EntryPrice,
				StopLoss:   request.StopLoss,
				TakeProfit: request.TakeProfit,
				Thesis:     request.Thesis,
			})
			if err != nil {
				writeSignalError(w, err)
				return
			}
			writeJSON(w, http.StatusCreated, mapPaperSignal(created))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func paperSignalStatusHandler(service signalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if service == nil {
			http.Error(w, "signal service is not configured", http.StatusNotImplemented)
			return
		}
		if r.Method != http.MethodPatch {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id, ok := parsePaperSignalStatusPath(r.URL.Path)
		if !ok {
			http.NotFound(w, r)
			return
		}

		var request updatePaperSignalStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		updated, err := service.UpdateStatus(r.Context(), id, request.Status)
		if err != nil {
			writeSignalError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, mapPaperSignal(updated))
	}
}

func parsePaperSignalStatusPath(path string) (string, bool) {
	trimmed := strings.TrimPrefix(path, "/api/paper-signals/")
	id, suffix, ok := strings.Cut(trimmed, "/")
	return id, ok && id != "" && suffix == "status"
}

func writeSignalError(w http.ResponseWriter, err error) {
	if errors.Is(err, appsignals.ErrUnsupportedSymbol) || errors.Is(err, appsignals.ErrInvalidStatus) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, signal.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Error(w, "paper signal request failed", http.StatusInternalServerError)
}

func mapPaperSignals(signals []signal.Signal) []paperSignalResponse {
	response := make([]paperSignalResponse, 0, len(signals))
	for _, item := range signals {
		response = append(response, mapPaperSignal(item))
	}
	return response
}

func mapPaperSignal(item signal.Signal) paperSignalResponse {
	return paperSignalResponse{
		ID:         item.ID,
		Symbol:     item.Symbol,
		Timeframe:  item.Timeframe,
		Side:       item.Side,
		Status:     string(item.Status),
		Confidence: item.Confidence,
		EntryPrice: item.EntryPrice,
		StopLoss:   item.StopLoss,
		TakeProfit: item.TakeProfit,
		Thesis:     item.Thesis,
		CreatedAt:  formatTime(item.CreatedAt),
		UpdatedAt:  formatTime(item.UpdatedAt),
	}
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

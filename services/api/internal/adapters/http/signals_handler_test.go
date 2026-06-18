package httpadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/signal"
)

func TestPaperSignalRoutesCreateListAndUpdateStatus(t *testing.T) {
	service := &fakeSignalService{
		created: signal.Signal{
			ID:         "018f4f8a-0000-7000-9000-000000000401",
			Symbol:     "XAUUSD",
			Timeframe:  "1h",
			Side:       "long",
			Status:     signal.StatusPendingReview,
			Confidence: 0.72,
			EntryPrice: 2325.5,
			StopLoss:   2310,
			TakeProfit: 2340,
			Thesis:     "Bullish continuation.",
			CreatedAt:  time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC),
			UpdatedAt:  time.Date(2026, 6, 18, 10, 0, 0, 0, time.UTC),
		},
	}
	service.listed = []signal.Signal{service.created}
	service.updated = service.created
	service.updated.Status = signal.StatusApprovedPaper
	router := NewRouter(WithSignalService(service))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/paper-signals", bytes.NewBufferString(`{
		"symbol":"XAUUSD",
		"timeframe":"1h",
		"side":"long",
		"confidence":0.72,
		"entryPrice":2325.5,
		"stopLoss":2310,
		"takeProfit":2340,
		"thesis":"Bullish continuation."
	}`))
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if service.input.Symbol != "XAUUSD" {
		t.Fatalf("input.Symbol = %q, want XAUUSD", service.input.Symbol)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/paper-signals", nil)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", rec.Code, http.StatusOK)
	}
	var listed []paperSignalResponse
	if err := json.NewDecoder(rec.Body).Decode(&listed); err != nil {
		t.Fatalf("decode listed: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("len(listed) = %d, want 1", len(listed))
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPatch, "/api/paper-signals/018f4f8a-0000-7000-9000-000000000401/status", bytes.NewBufferString(`{
		"status":"approved_paper"
	}`))
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", rec.Code, http.StatusOK)
	}
	if service.status != signal.StatusApprovedPaper {
		t.Fatalf("status = %q, want approved_paper", service.status)
	}
}

type fakeSignalService struct {
	input   signal.CreateInput
	status  signal.Status
	created signal.Signal
	listed  []signal.Signal
	updated signal.Signal
}

func (s *fakeSignalService) CreateSignal(ctx context.Context, input signal.CreateInput) (signal.Signal, error) {
	s.input = input
	return s.created, nil
}

func (s *fakeSignalService) ListSignals(ctx context.Context) ([]signal.Signal, error) {
	return append([]signal.Signal(nil), s.listed...), nil
}

func (s *fakeSignalService) UpdateStatus(ctx context.Context, id string, status signal.Status) (signal.Signal, error) {
	s.status = status
	return s.updated, nil
}

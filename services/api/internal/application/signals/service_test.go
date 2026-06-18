package signals

import (
	"context"
	"errors"
	"testing"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/signal"
)

func TestServiceCreatesListsAndUpdatesXAUUSDPaperSignals(t *testing.T) {
	repo := &memorySignalRepository{}
	service := NewService(repo, func() string { return "sig-1" })

	created, err := service.CreateSignal(context.Background(), signal.CreateInput{
		Symbol:     "XAUUSD",
		Timeframe:  "1h",
		Side:       "long",
		Confidence: 0.72,
		EntryPrice: 2325.5,
		StopLoss:   2310,
		TakeProfit: 2340,
		Thesis:     "Top-down bullish continuation with EMA alignment.",
	})
	if err != nil {
		t.Fatalf("CreateSignal returned error: %v", err)
	}
	if created.ID != "sig-1" {
		t.Fatalf("created.ID = %q, want sig-1", created.ID)
	}
	if created.Status != signal.StatusPendingReview {
		t.Fatalf("created.Status = %q, want pending_review", created.Status)
	}

	listed, err := service.ListSignals(context.Background())
	if err != nil {
		t.Fatalf("ListSignals returned error: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("len(listed) = %d, want 1", len(listed))
	}

	updated, err := service.UpdateStatus(context.Background(), "sig-1", signal.StatusApprovedPaper)
	if err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
	}
	if updated.Status != signal.StatusApprovedPaper {
		t.Fatalf("updated.Status = %q, want approved_paper", updated.Status)
	}
}

func TestServiceRejectsUnsupportedSignalSymbol(t *testing.T) {
	service := NewService(&memorySignalRepository{}, func() string { return "sig-1" })

	_, err := service.CreateSignal(context.Background(), signal.CreateInput{Symbol: "BTCUSD"})
	if !errors.Is(err, ErrUnsupportedSymbol) {
		t.Fatalf("CreateSignal error = %v, want ErrUnsupportedSymbol", err)
	}
}

func TestServiceRejectsInvalidStatus(t *testing.T) {
	service := NewService(&memorySignalRepository{}, func() string { return "sig-1" })

	_, err := service.UpdateStatus(context.Background(), "sig-1", signal.Status("executed"))
	if !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("UpdateStatus error = %v, want ErrInvalidStatus", err)
	}
}

type memorySignalRepository struct {
	signals []signal.Signal
}

func (r *memorySignalRepository) CreateSignal(ctx context.Context, created signal.Signal) (signal.Signal, error) {
	r.signals = append(r.signals, created)
	return created, nil
}

func (r *memorySignalRepository) ListSignals(ctx context.Context) ([]signal.Signal, error) {
	return append([]signal.Signal(nil), r.signals...), nil
}

func (r *memorySignalRepository) UpdateStatus(ctx context.Context, id string, status signal.Status) (signal.Signal, error) {
	for index, item := range r.signals {
		if item.ID == id {
			r.signals[index].Status = status
			return r.signals[index], nil
		}
	}
	return signal.Signal{}, signal.ErrNotFound
}

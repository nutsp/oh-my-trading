package signals

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/signal"
)

var (
	ErrUnsupportedSymbol = errors.New("unsupported signal symbol")
	ErrInvalidStatus     = errors.New("invalid signal status")
)

type Service struct {
	repo  signal.Repository
	newID func() string
	now   func() time.Time
}

func NewService(repo signal.Repository, newID func() string) *Service {
	return &Service{
		repo:  repo,
		newID: newID,
		now:   time.Now,
	}
}

func (s *Service) CreateSignal(ctx context.Context, input signal.CreateInput) (signal.Signal, error) {
	if input.Symbol != signal.XAUUSDSymbol {
		return signal.Signal{}, fmt.Errorf("%w: %s", ErrUnsupportedSymbol, input.Symbol)
	}

	now := s.now().UTC()
	created := signal.Signal{
		ID:         s.newID(),
		Symbol:     input.Symbol,
		Timeframe:  input.Timeframe,
		Side:       input.Side,
		Status:     signal.StatusPendingReview,
		Confidence: input.Confidence,
		EntryPrice: input.EntryPrice,
		StopLoss:   input.StopLoss,
		TakeProfit: input.TakeProfit,
		Thesis:     input.Thesis,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	result, err := s.repo.CreateSignal(ctx, created)
	if err != nil {
		return signal.Signal{}, fmt.Errorf("create paper signal: %w", err)
	}
	return result, nil
}

func (s *Service) ListSignals(ctx context.Context) ([]signal.Signal, error) {
	signals, err := s.repo.ListSignals(ctx)
	if err != nil {
		return nil, fmt.Errorf("list paper signals: %w", err)
	}
	return signals, nil
}

func (s *Service) UpdateStatus(ctx context.Context, id string, status signal.Status) (signal.Signal, error) {
	if !signal.IsValidStatus(status) || status == signal.StatusPendingReview {
		return signal.Signal{}, fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}

	updated, err := s.repo.UpdateStatus(ctx, id, status)
	if err != nil {
		return signal.Signal{}, fmt.Errorf("update paper signal status: %w", err)
	}
	return updated, nil
}

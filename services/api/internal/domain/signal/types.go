package signal

import (
	"context"
	"errors"
	"time"
)

const XAUUSDSymbol = "XAUUSD"

var ErrNotFound = errors.New("signal not found")

type Status string

const (
	StatusPendingReview Status = "pending_review"
	StatusApprovedPaper Status = "approved_paper"
	StatusRejected      Status = "rejected"
	StatusExpired       Status = "expired"
)

type Signal struct {
	ID         string
	Symbol     string
	Timeframe  string
	Side       string
	Status     Status
	Confidence float64
	EntryPrice float64
	StopLoss   float64
	TakeProfit float64
	Thesis     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateInput struct {
	Symbol     string
	Timeframe  string
	Side       string
	Confidence float64
	EntryPrice float64
	StopLoss   float64
	TakeProfit float64
	Thesis     string
}

type Repository interface {
	CreateSignal(ctx context.Context, created Signal) (Signal, error)
	ListSignals(ctx context.Context) ([]Signal, error)
	UpdateStatus(ctx context.Context, id string, status Status) (Signal, error)
}

func IsValidStatus(status Status) bool {
	switch status {
	case StatusPendingReview, StatusApprovedPaper, StatusRejected, StatusExpired:
		return true
	default:
		return false
	}
}

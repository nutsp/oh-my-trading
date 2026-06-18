package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/signal"
)

type PaperSignalRepository struct {
	db *sql.DB
}

func NewPaperSignalRepository(db *sql.DB) *PaperSignalRepository {
	return &PaperSignalRepository{db: db}
}

func (r *PaperSignalRepository) CreateSignal(ctx context.Context, created signal.Signal) (signal.Signal, error) {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO paper_signals (
		  id, symbol, timeframe, side, status, confidence, entry_price, stop_loss, take_profit, thesis, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id::text, symbol, timeframe, side, status, confidence::float8, entry_price::float8,
		          stop_loss::float8, take_profit::float8, thesis, created_at, updated_at
	`,
		created.ID,
		created.Symbol,
		created.Timeframe,
		created.Side,
		created.Status,
		created.Confidence,
		created.EntryPrice,
		created.StopLoss,
		created.TakeProfit,
		created.Thesis,
		created.CreatedAt,
		created.UpdatedAt,
	).Scan(
		&created.ID,
		&created.Symbol,
		&created.Timeframe,
		&created.Side,
		&created.Status,
		&created.Confidence,
		&created.EntryPrice,
		&created.StopLoss,
		&created.TakeProfit,
		&created.Thesis,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return signal.Signal{}, fmt.Errorf("create paper signal: %w", err)
	}
	return created, nil
}

func (r *PaperSignalRepository) ListSignals(ctx context.Context) ([]signal.Signal, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, symbol, timeframe, side, status, confidence::float8, entry_price::float8,
		       stop_loss::float8, take_profit::float8, thesis, created_at, updated_at
		FROM paper_signals
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list paper signals: %w", err)
	}
	defer rows.Close()

	var signals []signal.Signal
	for rows.Next() {
		item, err := scanSignal(rows)
		if err != nil {
			return nil, err
		}
		signals = append(signals, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate paper signals: %w", err)
	}
	return signals, nil
}

func (r *PaperSignalRepository) UpdateStatus(ctx context.Context, id string, status signal.Status) (signal.Signal, error) {
	var updated signal.Signal
	err := r.db.QueryRowContext(ctx, `
		UPDATE paper_signals
		SET status = $2, updated_at = now()
		WHERE id = $1
		RETURNING id::text, symbol, timeframe, side, status, confidence::float8, entry_price::float8,
		          stop_loss::float8, take_profit::float8, thesis, created_at, updated_at
	`, id, status).Scan(
		&updated.ID,
		&updated.Symbol,
		&updated.Timeframe,
		&updated.Side,
		&updated.Status,
		&updated.Confidence,
		&updated.EntryPrice,
		&updated.StopLoss,
		&updated.TakeProfit,
		&updated.Thesis,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return signal.Signal{}, signal.ErrNotFound
	}
	if err != nil {
		return signal.Signal{}, fmt.Errorf("update paper signal status: %w", err)
	}
	return updated, nil
}

type signalScanner interface {
	Scan(dest ...any) error
}

func scanSignal(scanner signalScanner) (signal.Signal, error) {
	var item signal.Signal
	err := scanner.Scan(
		&item.ID,
		&item.Symbol,
		&item.Timeframe,
		&item.Side,
		&item.Status,
		&item.Confidence,
		&item.EntryPrice,
		&item.StopLoss,
		&item.TakeProfit,
		&item.Thesis,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return signal.Signal{}, fmt.Errorf("scan paper signal: %w", err)
	}
	return item, nil
}

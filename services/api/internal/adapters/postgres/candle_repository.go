package postgres

import (
	"context"
	"database/sql"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

type CandleRepository struct {
	db *sql.DB
}

func NewCandleRepository(db *sql.DB) *CandleRepository {
	return &CandleRepository{db: db}
}

func (r *CandleRepository) UpsertCandles(ctx context.Context, candles []market.Candle) error {
	if len(candles) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO candles (symbol_id, timeframe, ts, open, high, low, close, volume)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (symbol_id, timeframe, ts)
		DO UPDATE SET
		  open = EXCLUDED.open,
		  high = EXCLUDED.high,
		  low = EXCLUDED.low,
		  close = EXCLUDED.close,
		  volume = EXCLUDED.volume
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, candle := range candles {
		if _, err := stmt.ExecContext(ctx,
			candle.SymbolID,
			candle.Timeframe,
			candle.Timestamp,
			candle.Open,
			candle.High,
			candle.Low,
			candle.Close,
			candle.Volume,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *CandleRepository) ListCandles(ctx context.Context, query market.CandleQuery) ([]market.Candle, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.symbol_id::text, c.timeframe, c.ts, c.open::float8, c.high::float8, c.low::float8, c.close::float8, COALESCE(c.volume, 0)::float8
		FROM candles c
		JOIN symbols s ON s.id = c.symbol_id
		WHERE s.code = $1
		  AND c.timeframe = $2
		  AND c.ts >= $3
		  AND c.ts <= $4
		ORDER BY c.ts ASC
	`, query.SymbolCode, query.Timeframe, query.From, query.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candles []market.Candle
	for rows.Next() {
		var candle market.Candle
		if err := rows.Scan(
			&candle.SymbolID,
			&candle.Timeframe,
			&candle.Timestamp,
			&candle.Open,
			&candle.High,
			&candle.Low,
			&candle.Close,
			&candle.Volume,
		); err != nil {
			return nil, err
		}
		candles = append(candles, candle)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return candles, nil
}

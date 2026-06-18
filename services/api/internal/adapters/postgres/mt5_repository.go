package postgres

import (
	"context"
	"database/sql"
	"fmt"

	domainmt5 "github.com/sutad-p/oh-my-trading/services/api/internal/domain/mt5"
)

type MT5Repository struct {
	db *sql.DB
}

func NewMT5Repository(db *sql.DB) *MT5Repository {
	return &MT5Repository{db: db}
}

type MT5Heartbeat = domainmt5.Heartbeat
type MT5Tick = domainmt5.Tick
type MT5AccountSnapshot = domainmt5.AccountSnapshot
type MT5PositionSnapshot = domainmt5.PositionSnapshot

func (r *MT5Repository) SaveHeartbeat(ctx context.Context, heartbeat MT5Heartbeat) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO mt5_bridge_heartbeats (
		  bridge_id, terminal, account_login, server, status, last_error, sent_at
		)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7)
		ON CONFLICT (bridge_id)
		DO UPDATE SET
		  terminal = EXCLUDED.terminal,
		  account_login = EXCLUDED.account_login,
		  server = EXCLUDED.server,
		  status = EXCLUDED.status,
		  last_error = EXCLUDED.last_error,
		  sent_at = EXCLUDED.sent_at,
		  received_at = now()
	`,
		heartbeat.BridgeID,
		heartbeat.Terminal,
		heartbeat.AccountLogin,
		heartbeat.Server,
		heartbeat.Status,
		heartbeat.LastError,
		heartbeat.SentAt,
	)
	if err != nil {
		return fmt.Errorf("save mt5 heartbeat: %w", err)
	}
	return nil
}

func (r *MT5Repository) LatestHeartbeat(ctx context.Context, bridgeID string) (MT5Heartbeat, error) {
	var heartbeat MT5Heartbeat
	var lastError sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT bridge_id, terminal, account_login, server, status, last_error, sent_at
		FROM mt5_bridge_heartbeats
		WHERE bridge_id = $1
	`, bridgeID).Scan(
		&heartbeat.BridgeID,
		&heartbeat.Terminal,
		&heartbeat.AccountLogin,
		&heartbeat.Server,
		&heartbeat.Status,
		&lastError,
		&heartbeat.SentAt,
	)
	if err != nil {
		return MT5Heartbeat{}, fmt.Errorf("latest mt5 heartbeat: %w", err)
	}
	heartbeat.LastError = lastError.String
	return heartbeat, nil
}

func (r *MT5Repository) SaveTicks(ctx context.Context, ticks []MT5Tick) error {
	if len(ticks) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO mt5_ticks (symbol, bid, ask, last, volume, ts)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (symbol, ts, bid, ask) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("prepare mt5 tick insert: %w", err)
	}
	defer stmt.Close()

	for _, tick := range ticks {
		if _, err := stmt.ExecContext(ctx, tick.Symbol, tick.Bid, tick.Ask, tick.Last, tick.Volume, tick.Time); err != nil {
			return fmt.Errorf("save mt5 tick: %w", err)
		}
	}

	return tx.Commit()
}

func (r *MT5Repository) LatestTick(ctx context.Context, symbol string) (MT5Tick, error) {
	var tick MT5Tick
	err := r.db.QueryRowContext(ctx, `
		SELECT symbol, bid::float8, ask::float8, COALESCE(last, 0)::float8, COALESCE(volume, 0)::float8, ts
		FROM mt5_ticks
		WHERE symbol = $1
		ORDER BY ts DESC, id DESC
		LIMIT 1
	`, symbol).Scan(&tick.Symbol, &tick.Bid, &tick.Ask, &tick.Last, &tick.Volume, &tick.Time)
	if err != nil {
		return MT5Tick{}, fmt.Errorf("latest mt5 tick: %w", err)
	}
	return tick, nil
}

func (r *MT5Repository) SaveAccountSnapshot(ctx context.Context, snapshot MT5AccountSnapshot) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO mt5_account_snapshots (
		  account_login, currency, balance, equity, margin, free_margin, margin_level, ts
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (account_login, ts)
		DO UPDATE SET
		  currency = EXCLUDED.currency,
		  balance = EXCLUDED.balance,
		  equity = EXCLUDED.equity,
		  margin = EXCLUDED.margin,
		  free_margin = EXCLUDED.free_margin,
		  margin_level = EXCLUDED.margin_level,
		  received_at = now()
	`,
		snapshot.AccountLogin,
		snapshot.Currency,
		snapshot.Balance,
		snapshot.Equity,
		snapshot.Margin,
		snapshot.FreeMargin,
		snapshot.MarginLevel,
		snapshot.Time,
	)
	if err != nil {
		return fmt.Errorf("save mt5 account snapshot: %w", err)
	}
	return nil
}

func (r *MT5Repository) LatestAccountSnapshot(ctx context.Context, accountLogin string) (MT5AccountSnapshot, error) {
	var snapshot MT5AccountSnapshot
	err := r.db.QueryRowContext(ctx, `
		SELECT account_login, currency, balance::float8, equity::float8, margin::float8, free_margin::float8, COALESCE(margin_level, 0)::float8, ts
		FROM mt5_account_snapshots
		WHERE account_login = $1
		ORDER BY ts DESC, id DESC
		LIMIT 1
	`, accountLogin).Scan(
		&snapshot.AccountLogin,
		&snapshot.Currency,
		&snapshot.Balance,
		&snapshot.Equity,
		&snapshot.Margin,
		&snapshot.FreeMargin,
		&snapshot.MarginLevel,
		&snapshot.Time,
	)
	if err != nil {
		return MT5AccountSnapshot{}, fmt.Errorf("latest mt5 account snapshot: %w", err)
	}
	return snapshot, nil
}

func (r *MT5Repository) SavePositionSnapshots(ctx context.Context, positions []MT5PositionSnapshot) error {
	if len(positions) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO mt5_position_snapshots (
		  account_login, ticket, symbol, side, volume, open_price, stop_loss, take_profit, profit, opened_at, snapshot_ts
		)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, 0), NULLIF($8, 0), $9, $10, $11)
		ON CONFLICT (account_login, ticket, snapshot_ts)
		DO UPDATE SET
		  symbol = EXCLUDED.symbol,
		  side = EXCLUDED.side,
		  volume = EXCLUDED.volume,
		  open_price = EXCLUDED.open_price,
		  stop_loss = EXCLUDED.stop_loss,
		  take_profit = EXCLUDED.take_profit,
		  profit = EXCLUDED.profit,
		  opened_at = EXCLUDED.opened_at,
		  received_at = now()
	`)
	if err != nil {
		return fmt.Errorf("prepare mt5 position insert: %w", err)
	}
	defer stmt.Close()

	for _, position := range positions {
		if _, err := stmt.ExecContext(ctx,
			position.AccountLogin,
			position.Ticket,
			position.Symbol,
			position.Side,
			position.Volume,
			position.OpenPrice,
			position.StopLoss,
			position.TakeProfit,
			position.Profit,
			position.OpenedAt,
			position.SnapshotTime,
		); err != nil {
			return fmt.Errorf("save mt5 position snapshot: %w", err)
		}
	}

	return tx.Commit()
}

func (r *MT5Repository) LatestPositionSnapshots(ctx context.Context, accountLogin string) ([]MT5PositionSnapshot, error) {
	rows, err := r.db.QueryContext(ctx, `
		WITH latest AS (
		  SELECT max(snapshot_ts) AS snapshot_ts
		  FROM mt5_position_snapshots
		  WHERE account_login = $1
		)
		SELECT p.account_login, p.ticket, p.symbol, p.side, p.volume::float8, p.open_price::float8,
		       COALESCE(p.stop_loss, 0)::float8, COALESCE(p.take_profit, 0)::float8, p.profit::float8,
		       p.opened_at, p.snapshot_ts
		FROM mt5_position_snapshots p
		JOIN latest ON latest.snapshot_ts = p.snapshot_ts
		WHERE p.account_login = $1
		ORDER BY p.ticket
	`, accountLogin)
	if err != nil {
		return nil, fmt.Errorf("latest mt5 position snapshots: %w", err)
	}
	defer rows.Close()

	var positions []MT5PositionSnapshot
	for rows.Next() {
		var position MT5PositionSnapshot
		if err := rows.Scan(
			&position.AccountLogin,
			&position.Ticket,
			&position.Symbol,
			&position.Side,
			&position.Volume,
			&position.OpenPrice,
			&position.StopLoss,
			&position.TakeProfit,
			&position.Profit,
			&position.OpenedAt,
			&position.SnapshotTime,
		); err != nil {
			return nil, fmt.Errorf("scan mt5 position snapshot: %w", err)
		}
		positions = append(positions, position)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate mt5 position snapshots: %w", err)
	}
	return positions, nil
}

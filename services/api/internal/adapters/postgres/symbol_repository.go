package postgres

import (
	"context"
	"database/sql"

	"github.com/sutad-p/oh-my-trading/services/api/internal/domain/market"
)

type SymbolRepository struct {
	db *sql.DB
}

func NewSymbolRepository(db *sql.DB) *SymbolRepository {
	return &SymbolRepository{db: db}
}

func (r *SymbolRepository) ListSymbols(ctx context.Context) ([]market.Symbol, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id::text, code, market, COALESCE(base_asset, ''), COALESCE(quote_asset, ''), enabled
		FROM symbols
		ORDER BY code
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var symbols []market.Symbol
	for rows.Next() {
		var symbol market.Symbol
		if err := rows.Scan(
			&symbol.ID,
			&symbol.Code,
			&symbol.Market,
			&symbol.BaseAsset,
			&symbol.QuoteAsset,
			&symbol.Enabled,
		); err != nil {
			return nil, err
		}
		symbols = append(symbols, symbol)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return symbols, nil
}

func (r *SymbolRepository) CreateSymbol(ctx context.Context, symbol market.Symbol) (market.Symbol, error) {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO symbols (id, code, market, base_asset, quote_asset, enabled)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, code, market, COALESCE(base_asset, ''), COALESCE(quote_asset, ''), enabled
	`,
		symbol.ID,
		symbol.Code,
		symbol.Market,
		symbol.BaseAsset,
		symbol.QuoteAsset,
		symbol.Enabled,
	).Scan(
		&symbol.ID,
		&symbol.Code,
		&symbol.Market,
		&symbol.BaseAsset,
		&symbol.QuoteAsset,
		&symbol.Enabled,
	)
	if err != nil {
		return market.Symbol{}, err
	}
	return symbol, nil
}

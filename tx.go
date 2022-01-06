package pge

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

type tx struct {
	pgx.Tx
}

func (t tx) Execute(ctx context.Context, query Query, args ...interface{}) (pgconn.CommandTag, error) {
	return pgExec(ctx, t.Tx, query, args...)
}

func (t tx) Get(ctx context.Context, dst interface{}, query Query, args ...interface{}) error {
	return pgGet(ctx, t.Tx, dst, query, args...)
}

func (t tx) Select(ctx context.Context, dst interface{}, query Query, args ...interface{}) error {
	return pgSelect(ctx, t.Tx, dst, query, args...)
}

func (t tx) PaginatedSelect(ctx context.Context, dst interface{}, query Query, args ...interface{}) (Cursors, error) {
	return pgPaginatedSelect(ctx, t.Tx, dst, query, args...)
}

type TxOption func(*TxInfo)

type TxInfo struct {
	opts pgx.TxOptions
}

func WithIsoLevel(level pgx.TxIsoLevel) TxOption {
	return func(info *TxInfo) {
		info.opts.IsoLevel = level
	}
}

func WithAccessMode(mode pgx.TxAccessMode) TxOption {
	return func(info *TxInfo) {
		info.opts.AccessMode = mode
	}
}

func WithDeferrableMode(mode pgx.TxDeferrableMode) TxOption {
	return func(info *TxInfo) {
		info.opts.DeferrableMode = mode
	}
}

type queryable interface {
	pgxscan.Querier

	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}

package pge

import (
	"context"
	"io"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

type Store interface {
	Conn
	io.Closer

	Migrate(ctx context.Context, migrations []Migration) error

	Tx(ctx context.Context, fn func(tx Tx) error, opts ...TxOption) error
}

type Tx interface {
	pgx.Tx
	Conn
}

type Conn interface {
	Execute(ctx context.Context, query Query, args ...interface{}) (pgconn.CommandTag, error)

	Get(ctx context.Context, dst interface{}, query Query, args ...interface{}) error

	Select(ctx context.Context, dst interface{}, query Query, args ...interface{}) error

	PaginatedSelect(ctx context.Context, dst interface{}, query Query, args ...interface{}) (Cursors, error)
}

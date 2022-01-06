package pge

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type store struct {
	pool *pgxpool.Pool
}

func NewStore(ctx context.Context, cfg *pgxpool.Config) (Store, error) {
	conn, err := pgxpool.ConnectConfig(ctx, cfg)
	return &store{pool: conn}, err
}

func (s *store) Close() error {
	s.pool.Close()
	return nil
}

func (s *store) Migrate(ctx context.Context, migrations []Migration) error {
	return s.Tx(ctx, func(tx Tx) error {
		// 1. Acquire advisory lock governing schema migrations
		_, err := tx.Execute(ctx, LockMigrations)
		if err != nil {
			return err
		}
		fmt.Println("Acquired advisory lock")

		// 2. Query schema version.
		_, err = tx.Execute(ctx, CreateTableSchemaVersion)
		if err != nil {
			return err
		}
		fmt.Println("Query schema version")

		version := 0
		err = tx.Get(ctx, &version, SelectLatestSchemaVersion)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		fmt.Println("Query latest version")

		// 3. Execute all migrations from current schema version.
		for i := version; i < len(migrations); i++ {
			migration := migrations[i]
			for _, query := range migration.Queries {
				_, err := tx.Execute(ctx, query)
				if err != nil {
					return err
				}
			}
		}
		fmt.Println("Finished migrations")

		// Up to date.
		if version == len(migrations) {
			return nil
		}

		// Update schema version.
		_, err = tx.Execute(ctx, InsertSchemaVersion, len(migrations))
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *store) Tx(ctx context.Context, fn func(tx Tx) error, opts ...TxOption) error {
	var info TxInfo
	for _, opt := range opts {
		opt(&info)
	}

	sqlTx, err := s.pool.BeginTx(ctx, info.opts)
	if err != nil {
		return err
	}
	defer sqlTx.Rollback(ctx)

	err = fn(tx{Tx: sqlTx})
	if err != nil {
		return err
	}

	return sqlTx.Commit(ctx)
}

func (s *store) Execute(ctx context.Context, query Query, args ...interface{}) (pgconn.CommandTag, error) {
	return pgExec(ctx, s.pool, query, args...)
}

func (s *store) Get(ctx context.Context, dst interface{}, query Query, args ...interface{}) error {
	return pgGet(ctx, s.pool, dst, query, args...)
}

func (s *store) Select(ctx context.Context, dst interface{}, query Query, args ...interface{}) error {
	return pgSelect(ctx, s.pool, dst, query, args...)
}

func (s *store) PaginatedSelect(ctx context.Context, dst interface{}, query Query, args ...interface{}) (Cursors, error) {
	return pgPaginatedSelect(ctx, s.pool, dst, query, args...)
}

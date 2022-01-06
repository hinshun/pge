package pge

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
)

// Listen takes a *pgx.Conn as an argument because we want the LISTEN to be
// effective for a specific connection.
func Listen(ctx context.Context, conn *pgx.Conn, channel string) error {
	return sqlQuery(ctx, NewQuery("LISTEN", ""), func(ctx context.Context) error {
		_, err := conn.Exec(ctx, "LISTEN "+channel)
		return err
	})
}

// Unlisten takes a *pgx.Conn as an argument because we want the UNLISTEN to be
// effective for a specific connection.
func Unlisten(ctx context.Context, conn *pgx.Conn, channel string) error {
	return sqlQuery(ctx, NewQuery("UNLISTEN", ""), func(ctx context.Context) error {
		_, err := conn.Exec(ctx, "UNLISTEN "+channel)
		return err
	})
}

func pgExec(ctx context.Context, q queryable, query Query, args ...interface{}) (pgconn.CommandTag, error) {
	var tag pgconn.CommandTag
	err := sqlQuery(ctx, query, func(ctx context.Context) error {
		var err error
		tag, err = q.Exec(ctx, query.String(), args...)
		return err
	})
	return tag, err
}

func pgGet(ctx context.Context, q queryable, dst interface{}, query Query, args ...interface{}) error {
	return sqlQuery(ctx, query, func(ctx context.Context) error {
		return pgxscan.Get(ctx, q, dst, query.String(), args...)
	})
}

func pgSelect(ctx context.Context, q queryable, dst interface{}, query Query, args ...interface{}) error {
	return sqlQuery(ctx, query, func(ctx context.Context) error {
		return pgxscan.Select(ctx, q, dst, query.String(), args...)
	})
}

func pgPaginatedSelect(ctx context.Context, q queryable, dst interface{}, query Query, args ...interface{}) (Cursors, error) {
	if query.paginator != nil {
		if query.paginator.AfterCursor != "" {
			cursor, err := query.CursorFromString(query.paginator.AfterCursor)
			if err != nil {
				return Cursors{}, err
			}
			args = append(args, cursor)
		} else if query.paginator.BeforeCursor != "" {
			cursor, err := query.CursorFromString(query.paginator.BeforeCursor)
			if err != nil {
				return Cursors{}, err
			}
			args = append(args, cursor)
		}
	}

	err := sqlQuery(ctx, query, func(ctx context.Context) error {
		return pgxscan.Select(ctx, q, dst, query.String(), args...)
	})
	if err != nil {
		return Cursors{}, err
	}
	return extractCursors(query.paginator, dst)
}

func BulkInsert(ctx context.Context, numRows int, q queryable, query Query, fn func(int) []interface{}) (results []interface{}, err error) {
	results = make([]interface{}, numRows)
	n, err := Bulk(numRows, 100, func(i, j int) error {
		var args []interface{}
		for k := i; k < j; k++ {
			args = append(args, fn(k)...)
		}

		var batch []interface{}
		err := pgSelect(ctx, q, &batch, query.WithValues(j-i), args...)
		if err != nil {
			return err
		}

		for k, item := range batch {
			results[i+k] = item
		}
		return nil
	})
	return results[:n], err
}

func Bulk(total, batch int, fn func(i, j int) error) (n int, err error) {
	if total == 0 {
		return
	}

	i := 0
	for ; (i+1)*batch < total; i++ {
		err := fn(i*batch, (i+1)*batch)
		if err != nil {
			return i * batch, err
		}
	}

	err = fn(i*batch, total)
	if err != nil {
		return i * batch, err
	}

	return total, nil
}

type QueryError struct {
	Query Query
	Err   error
}

func (e *QueryError) Error() string {
	return e.Query.Name + ": " + e.Err.Error()
}

func (e *QueryError) Unwrap() error {
	return e.Err
}

func sqlQuery(ctx context.Context, query Query, f func(context.Context) error) error {
	err := f(ctx)
	if err != nil {
		return &QueryError{query, err}
	}
	return nil
}

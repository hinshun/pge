package main

import (
	"context"

	"github.com/hinshun/pge"
	"github.com/hinshun/pge/example/sql"
)

func SelectCustomers(ctx context.Context, c pge.Conn) (cts []Customer, err error) {
	err = c.Select(ctx, &cts, sql.SelectCustomers)
	return
}

func SelectCustomerByID(ctx context.Context, c pge.Conn, customerID int64) (ct Customer, err error) {
	err = c.Get(ctx, &ct, sql.SelectCustomerByID, customerID)
	return
}

func SelectTransactionsByProductID(ctx context.Context, c pge.Conn, productID int64) (txs []Transaction, err error) {
	err = c.Select(ctx, &txs, sql.SelectTransactionsByProductID, productID)
	return
}

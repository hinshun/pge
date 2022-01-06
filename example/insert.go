package main

import (
	"context"

	"github.com/hinshun/pge"
	"github.com/hinshun/pge/example/sql"
)

func InsertCustomer(ctx context.Context, c pge.Conn, ct Customer) (id int64, err error) {
	err = c.Get(ctx, &id, sql.InsertCustomer, ct.PostalCode)
	return
}

func InsertProduct(ctx context.Context, c pge.Conn, p Product) (id int64, err error) {
	err = c.Get(ctx, &id, sql.InsertProduct, p.ProductName, p.Price, p.Quantity)
	return
}

func InsertTransaction(ctx context.Context, c pge.Conn, t Transaction) (id int64, err error) {
	err = c.Get(ctx, &id, sql.InsertTransaction, t.CustomerID)
	return
}

func InsertProductTransaction(ctx context.Context, c pge.Conn, pt ProductTransaction) error {
	_, err := c.Execute(ctx, sql.InsertProductTransaction, pt.ProductID, pt.TransactionID)
	return err
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hinshun/pge"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("example expects exactly 1 arg <postgres-uri>")
	}

	err := run(context.Background(), os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, uri string) error {
	cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return err
	}

	store, err := pge.NewStore(ctx, cfg)
	if err != nil {
		return err
	}
	defer store.Close()
	fmt.Println("Store connected")

	err = store.Migrate(ctx, SchemaMigrations)
	if err != nil {
		return err
	}
	fmt.Println("Store migrated")

	customers, err := SelectCustomers(ctx, store)
	if err != nil {
		return err
	}
	fmt.Println("Num customers", len(customers))

	var productID int64
	err = store.Tx(ctx, func(tx pge.Tx) error {
		customerID, err := InsertCustomer(ctx, tx, Customer{PostalCode: "123456"})
		if err != nil {
			return err
		}

		productID, err = InsertProduct(ctx, tx, Product{
			ProductName: "Apple",
			Price:       300, // 3 dollars in cents.
			Quantity:    5,
		})
		if err != nil {
			return err
		}

		txID, err := InsertTransaction(ctx, tx, Transaction{CustomerID: customerID})
		if err != nil {
			return err
		}

		return InsertProductTransaction(ctx, tx, ProductTransaction{
			ProductID:     productID,
			TransactionID: txID,
		})
	})
	if err != nil {
		return err
	}

	txs, err := SelectTransactionsByProductID(ctx, store, productID)
	if err != nil {
		return err
	}

	fmt.Println("Num transactions", len(txs))
	if len(txs) > 0 {
		fmt.Println("Customer ID", txs[0].CustomerID)
		customer, err := SelectCustomerByID(ctx, store, txs[0].CustomerID)
		if err != nil {
			return err
		}
		fmt.Println("Customer Postal Code", customer.PostalCode)
	}
	return nil
}

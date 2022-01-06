package sql

import "github.com/hinshun/pge"

var (
	CreateTableCustomers = pge.NewQuery("create table customers", `
		CREATE TABLE IF NOT EXISTS customers (
			id bigserial PRIMARY KEY,
			postal_code text NOT NULL
		);
	`)

	CreateTableProducts = pge.NewQuery("create table products", `
		CREATE TABLE IF NOT EXISTS products (
			id bigserial PRIMARY KEY,
			product_name text NOT NULL,
			price bigint NOT NULL,
			quantity int NOT NULL
		);
	`)

	CreateTableTransactions = pge.NewQuery("create table transactions", `
		CREATE TABLE IF NOT EXISTS transactions (
			id bigserial PRIMARY KEY,
			customer_id bigint NOT NULL,
			created timestamptz NOT NULL DEFAULT NOW(),
			FOREIGN KEY (customer_id) REFERENCES customers (id)
		);
	`)

	CreateTableProductTransactions = pge.NewQuery("create table product_transactions", `
		CREATE TABLE IF NOT EXISTS product_transactions (
			product_id bigint,
			transaction_id bigint,
			PRIMARY KEY (product_id, transaction_id),
			FOREIGN KEY (product_id) REFERENCES products (id),
			FOREIGN KEY (transaction_id) REFERENCES transactions (id)
		);
	`)
)

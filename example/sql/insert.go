package sql

import "github.com/hinshun/pge"

var (
	InsertCustomer = pge.NewQuery("insert customer", `
		INSERT INTO customers(postal_code)
	`, pge.WithInsertColumns(1), pge.WithSuffix(`
		RETURNING id
	`))

	InsertProduct = pge.NewQuery("insert product", `
		INSERT INTO products(product_name, price, quantity)
	`, pge.WithInsertColumns(3), pge.WithSuffix(`
		RETURNING id
	`))

	InsertTransaction = pge.NewQuery("insert transaction", `
		INSERT INTO transactions(customer_id)
	`, pge.WithInsertColumns(1), pge.WithSuffix(`
		RETURNING id
	`))

	InsertProductTransaction = pge.NewQuery("insert product transaction", `
		INSERT INTO product_transactions(product_id, transaction_id)
	`, pge.WithInsertColumns(2))
)

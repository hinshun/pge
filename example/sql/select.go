package sql

import "github.com/hinshun/pge"

var (
	SelectCustomers = pge.NewQuery("select customers", `
		SELECT *
		FROM customers;
	`)

	SelectCustomerByID = pge.NewQuery("select customer by id", `
		SELECT *
		FROM customers
		WHERE id = $1;
	`)

	SelectTransactionsByProductID = pge.NewQuery("select transactions by product id", `
		SELECT t.id, t.customer_id, t.created
		FROM products p
		JOIN product_transactions pt ON pt.product_id = p.id
		JOIN transactions t ON t.id = pt.transaction_id
		WHERE p.id = $1;
	`)
)

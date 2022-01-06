package main

import "time"

type Customer struct {
	ID         int64
	PostalCode string
}

type Product struct {
	ID          int64
	ProductName string
	Price       int64
	Quantity    int
}

type Transaction struct {
	ID         int64
	CustomerID int64
	Created    time.Time
}

type ProductTransaction struct {
	ProductID     int64
	TransactionID int64
}
